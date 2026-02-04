# kick user from server

POST http://localhost:8080/api/server/{{server_id}}/kick
Authorization: Bearer {{token}}
Content-Type: application/json

{
  "user_id": "{{user_id}}"
}

HTTP 200