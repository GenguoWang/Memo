#! /bin/bash
echo "Post Word"
curl -X POST http://localhost:8080/api/word \
   -H 'Content-Type: application/json' \
   -d '{"name":"test_wor"}'

echo "Fetch Words"
curl -X GET http://localhost:8080/api/word

echo "Delete Word"
curl -X DELETE http://localhost:8080/api/word \
   -H 'Content-Type: application/json' \
   -d '{"name":"test_wor"}'