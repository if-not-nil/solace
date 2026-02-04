# update role

PUT http://localhost:8080/api/server/{{server_id}}/roles/{{role_id}}
Authorization: Bearer {{token}}
Content-Type: application/json

{
  "name": "Updated Role",
  "permissions": ["send_messages", "manage_channels"]
}

HTTP 200
[Asserts]
jsonpath "$.role.name" == "Updated Role"