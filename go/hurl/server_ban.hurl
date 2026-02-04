# ban user from server

POST http://localhost:8080/api/server/{{server_id}}/ban
Authorization: Bearer {{token}}
Content-Type: application/json

{
  "user_id": "{{user_id}}",
  "reason": "Violation of server rules"
}

HTTP 200