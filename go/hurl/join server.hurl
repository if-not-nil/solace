# join server

POST http://localhost:8080/api/server/join
Authorization: Bearer {{token}}
Content-Type: application/json

{
  "invite_code": "{{invite_code}}"
}

HTTP 200