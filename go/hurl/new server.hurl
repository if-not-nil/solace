# new server

POST http://localhost:8080/api/server
Authorization: Bearer {{token}}
Content-Type: application/json

{
  "name": "My Awesome Server"
}

HTTP 201
[Asserts]
jsonpath "$.server.name" == "My Awesome Server"
jsonpath "$.server.id" exists