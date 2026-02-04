# register invalid - empty fields

POST http://localhost:8080/api/auth/register
Content-Type: application/json

{
  "name": "",
  "password": ""
}

HTTP 400