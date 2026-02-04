# join server with invite link

POST http://localhost:8080/api/server/join/{{invite_id}}
Authorization: Bearer {{token}}

HTTP 200
[Asserts]
jsonpath "$.server.id" exists