# list invite links

GET http://localhost:8080/api/server/{{server_id}}/links
Authorization: Bearer {{token}}

HTTP 200
[Asserts]
jsonpath "$.links" exists