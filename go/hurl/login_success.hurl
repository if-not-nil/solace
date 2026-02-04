# login success

POST http://localhost:8080/api/auth/login
Content-Type: application/json

{
  "name": "testuser",
  "password": "password123"
}

HTTP 200
[Asserts]
jsonpath "$.message" == "login successful"
jsonpath "$.token" exists
jsonpath "$.user.name" == "testuser"