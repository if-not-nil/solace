# assign role to user

PUT http://localhost:8080/api/server/{{server_id}}/roles/assign/{{user_id}}/{{role_id}}
Authorization: Bearer {{token}}

HTTP 200