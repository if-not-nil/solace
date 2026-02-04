# register success

POST http://localhost:8080/api/auth/register
Content-Type: application/json

{
  "name": "testuser",
  "password": "password123"
}

HTTP 201
[Asserts]
jsonpath "$.message" == "registration successful"
jsonpath "$.token" exists
jsonpath "$.user.name" == "testuser"
jsonpath "$.user.id" exists