# refresh token

GET http://localhost:8080/api/auth/refresh
Authorization: Bearer {{token}}

HTTP 200
[Asserts]
jsonpath "$.token" exists