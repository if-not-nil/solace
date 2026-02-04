# revoke invite link

DELETE http://localhost:8080/api/server/{{server_id}}/links/{{invite_id}}
Authorization: Bearer {{token}}

HTTP 200