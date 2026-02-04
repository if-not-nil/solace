package server

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"solace/orm"
	"solace/util"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type WSMessage struct {
	Type    string `json:"type"`
	Payload any    `json:"payload"`
}

type SwapChannelPayload struct {
	ChannelID string `json:"channel_id"`
}

type Client struct {
	conn      *websocket.Conn
	send      chan []byte
	userID    string
	channelID string
}

type Hub struct {
	clients    map[string]*Client // userID -> client
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
	mu         sync.RWMutex
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var hub = &Hub{
	clients:    make(map[string]*Client),
	register:   make(chan *Client),
	unregister: make(chan *Client),
	broadcast:  make(chan []byte),
}

func init() {
	go hub.run()
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.userID] = client
			h.mu.Unlock()
			log.Printf("websocket connected: %s", client.userID)

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.userID]; ok {
				delete(h.clients, client.userID)
				close(client.send)
				log.Printf("websocket disconnected: %s", client.userID)
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.mu.RLock()
			for _, client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client.userID)
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (c *Client) readPump() {
	defer func() {
		hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadDeadline(time.Now().Add(10 * time.Second))

	var authMsg WSMessage
	err := c.conn.ReadJSON(&authMsg)
	if err != nil || authMsg.Type != "auth" {
		log.Printf("[ws] failed authentication: %v", err)
		return
	}

	payload, ok := authMsg.Payload.(map[string]any)
	if !ok {
		log.Println("[ws] auth payload is not a map")
		return
	}

	token, ok := payload["token"].(string)
	if !ok {
		log.Println("[ws] auth payload is missing 'token'")
		return
	}

	decoded, err := util.DecodeJWT(token)
	if err != nil {
		log.Printf("[ws] failed to decode JWT: %v", err)
		return
	}

	userID, err := decoded.GetSubject()
	if err != nil {
		log.Printf("[ws] invalid JWT subject: %v", err)
		return
	}

	channelID, ok := payload["channel_id"].(string)
	if !ok || channelID == "" {
		log.Printf("[ws] no initial channel_id provided for user %s", userID)
	}

	c.userID = userID
	c.channelID = channelID
	hub.register <- c

	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		var msg WSMessage
		err := c.conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("[ws] error reading JSON: %v (for %s)", c.userID, err)
			}
			break
		}

		switch msg.Type {
		case "swap_channel":
			if payload, ok := msg.Payload.(map[string]any); ok {
				if newChannelID, exists := payload["channel_id"].(string); exists {
					hub.mu.Lock()
					c.channelID = newChannelID
					hub.mu.Unlock()
					log.Printf("[ws] user %s swapped to channel %s", c.userID, newChannelID)
				}
			}
		case "ping":
			response, _ := json.Marshal(WSMessage{Type: "pong", Payload: nil})
			select {
			case c.send <- response:
			default:
			}
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("[ws] write error for %s in %s: %v", c.userID, c.channelID, err)
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("[ws] ping error for %s in %s: %v", c.userID, c.channelID, err)
				return
			}
		}
	}
}

func (s *Server) WSConnect(c echo.Context) error {
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return err
	}

	client := &Client{
		conn:      conn,
		send:      make(chan []byte, 256),
		userID:    "",
		channelID: "",
	}

	go client.writePump()
	go client.readPump()

	return nil
}

func SendToChannel(channelID string, msg orm.MessageResponse) {
	asJSON, _ := json.Marshal(WSMessage{
		Type:    "message",
		Payload: msg,
	})

	hub.mu.RLock()
	defer hub.mu.RUnlock()

	for _, client := range hub.clients {
		if client.channelID == channelID {
			select {
			case client.send <- asJSON:
			default:
				// client's send channel is full, close it
				close(client.send)
				delete(hub.clients, client.userID)
			}
		}
	}
}
func SendMessageEditedToChannel(channelID string, messageID string, content string) {
	asJSON, _ := json.Marshal(WSMessage{
		Type:    "message_edited",
		Payload: map[string]string{"message_id": messageID, "new_content": content},
	})

	hub.mu.RLock()
	defer hub.mu.RUnlock()

	for _, client := range hub.clients {
		if client.channelID == channelID {
			select {
			case client.send <- asJSON:
			default:
				// client's send channel is full, close it
				close(client.send)
				delete(hub.clients, client.userID)
			}
		}
	}
}
func SendMessageDeletedToChannel(channelID string, messageID string) {
	asJSON, _ := json.Marshal(WSMessage{
		Type:    "message_deleted",
		Payload: map[string]string{"message_id": messageID},
	})

	hub.mu.RLock()
	defer hub.mu.RUnlock()

	for _, client := range hub.clients {
		if client.channelID == channelID {
			select {
			case client.send <- asJSON:
			default:
				// client's send channel is full, close it
				close(client.send)
				delete(hub.clients, client.userID)
			}
		}
	}
}
