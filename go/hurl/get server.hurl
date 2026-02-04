# get server

GET http://localhost:8080/api/server/{{server_id}}
Authorization: Bearer {{token}}

HTTP 200
[Asserts]
jsonpath "$.server" exists
jsonpath "$.channels" exists
jsonpath "$.members" exists