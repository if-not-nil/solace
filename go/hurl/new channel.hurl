# new channel

POST http://localhost:8080/api/server/{{server_id}}/new_channel
Authorization: Bearer {{token}}
Content-Type: application/json

{
  "name": "new-channel",
  "type": "text"
}

HTTP 201