# get current user

GET http://localhost:8080/api/user/me
Authorization: Bearer {{token}}

HTTP 200
[Asserts]
jsonpath "$.user.name" exists
jsonpath "$.user.id" exists