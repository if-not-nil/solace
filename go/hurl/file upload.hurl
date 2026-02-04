# upload file

POST http://localhost:8080/api/file/upload
Authorization: Bearer {{token}}
Content-Type: multipart/form-data

--boundary
Content-Disposition: form-data; name="file"; filename="test.txt"
Content-Type: text/plain

test file content
--boundary--

HTTP 201
[Asserts]
jsonpath "$.file.id" exists
jsonpath "$.file.name" == "test.txt"