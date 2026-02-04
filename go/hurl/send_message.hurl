# send message to channel

POST http://localhost:8080/api/channel/{{channel_id}}/send
Authorization: Bearer {{token}}
Content-Type: application/json

{
  "content": "Hello, world!"
}

HTTP 200