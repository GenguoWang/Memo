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

echo "Put Notes"
curl -X PUT http://localhost:8080/api/note \
   -H 'Content-Type: plain/text' \
   -d 'Note1
Note11

Note2
Note22
   

Note3
Note33


Note44

Note55





'

echo "Fetch Notes"
curl -X GET http://localhost:8080/api/note
