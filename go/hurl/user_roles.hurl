# get user roles

GET http://localhost:8080/api/server/{{server_id}}/user/roles
Authorization: Bearer {{token}}

HTTP 200
[Asserts]
jsonpath "$.roles" exists