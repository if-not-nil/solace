//go:generate go run github.com/tkrajina/typescriptify-golang-structs/typescriptify -package=orm -input=types.go -output=../../web/src/lib/orm_types.ts -overwrite

package server

import "time"

// basic types for typescript generation
type UserResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	AvatarID string `json:"avatar"`
}

type ServerResponse struct {
	ID       string            `json:"id"`
	Name     string            `json:"name"`
	OwnerID  string            `json:"owner_id"`
	AvatarID string            `json:"avatar"`
	Channels []ChannelResponse `json:"channels"`
}

type ChannelResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	ServerID string `json:"server_id"`
	Type     string `json:"type"`
}

type MessageResponse struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	UserID    string    `json:"user_id"`
	ChannelID string    `json:"channel_id"`
	CreatedAt time.Time `json:"created_at"`
}

type MessageHistoryResponse struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// request types
type LoginRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type ChannelSendRequest struct {
	Content string `json:"content"`
}

type ChannelHistoryRequest struct {
	Until *int64 `json:"until,omitempty"`
	Count int    `json:"count,omitempty"`
}

type ServerNewRequest struct {
	Name string `json:"name"`
}

type ChannelNewRequest struct {
	Name string `json:"name"`
}

type InviteLinksNewRequest struct {
	MaxUsers uint `json:"max_users"`
}

type RoleAssignRequest struct {
	UserID string `json:"user_id"`
	RoleID string `json:"role_id"`
}

type ServerKickRequest struct {
	UserID string `json:"id"`
}

type RoleCreateRequest struct {
	Name        string   `json:"name"`
	Permissions []string `json:"permissions"`
}

type RoleUpdateRequest struct {
	Administrator   bool `json:"administrator"`
	ManageServer    bool `json:"manage_server"`
	ManageRoles     bool `json:"manage_roles"`
	ManageChannels  bool `json:"manage_channels"`
	KickUsers       bool `json:"kick_users"`
	BanUsers        bool `json:"ban_users"`
	CreateInvites   bool `json:"create_invites"`
	ManageInvites   bool `json:"manage_invites"`
	SendMessages    bool `json:"send_messages"`
	ManageMessages  bool `json:"manage_messages"`
	EmbedLinks      bool `json:"embed_links"`
	AttachFiles     bool `json:"attach_files"`
	MentionEveryone bool `json:"mention_everyone"`
}

// response types
type LoginResponse struct {
	Message string       `json:"message"`
	Token   string       `json:"token"`
	User    UserResponse `json:"user"`
}

type RegisterResponse struct {
	Message string       `json:"message"`
	Token   string       `json:"token"`
	User    UserResponse `json:"user"`
}

type UserMeResponse struct {
	User    UserResponse     `json:"user"`
	Servers []ServerResponse `json:"servers"`
}

type MessageSendResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

type ChannelHistoryResponse struct {
	Messages []MessageHistoryResponse `json:"messages"`
	Users    []UserResponse           `json:"users"`
}

type InviteLinkResponse struct {
	ID        string `json:"id"`
	ServerID  string `json:"server_id"`
	JoinsLeft uint   `json:"joins_left"`
}

type RoleResponse struct {
	ID          string `json:"id"`
	ServerID    string `json:"server_id"`
	Name        string `json:"name"`
	Permissions uint64 `json:"permissions"`
}

// note: these types are defined in orm/types.go and re-exported here for convenience
// in handlers. the actual definitions are in the orm package.
