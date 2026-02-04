# create invite link

POST http://localhost:8080/api/server/{{server_id}}/links
Authorization: Bearer {{token}}
Content-Type: application/json

{
  "max_uses": 10
}

HTTP 201
[Asserts]
jsonpath "$.link.id" exists
jsonpath "$.link.code" exists
jsonpath "$.link.max_uses" == 10