# subscribe to channel sse

GET http://localhost:8080/api/channel/{{channel_id}}/sse
Authorization: Bearer {{token}}
Accept: text/event-stream

HTTP 200