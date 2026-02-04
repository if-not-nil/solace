# edit message

PUT http://localhost:8080/api/message/{{message_id}}
Authorization: Bearer {{token}}
Content-Type: application/json

{
  "content": "Updated message content"
}

HTTP 200
[Asserts]
jsonpath "$.message.content" == "Updated message content"