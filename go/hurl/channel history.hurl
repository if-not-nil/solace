# get channel history

POST http://localhost:8080/api/channel/{{channel_id}}/history
Authorization: Bearer {{token}}
Content-Type: application/json

HTTP 200
[Asserts]
jsonpath "$.messages" exists
jsonpath "$.users" exists