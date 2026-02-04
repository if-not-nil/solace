# register invalid - name too short

POST http://localhost:8080/api/auth/register
Content-Type: application/json

{
  "name": "ab",
  "password": "password123"
}

HTTP 400