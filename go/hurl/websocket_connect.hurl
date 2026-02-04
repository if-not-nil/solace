# websocket connection

GET http://localhost:8080/api/ws/connect
Authorization: Bearer {{token}}
Connection: Upgrade
Upgrade: websocket
Sec-WebSocket-Key: {{websocket_key}}
Sec-WebSocket-Version: 13

HTTP 101