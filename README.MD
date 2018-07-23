

Test manually: (192.168.99.100 insted of local host is only for Windows-toolbox docker...)

```
curl -v POST \
http://192.168.99.100:8000/todo/ \
-H 'content-type: application/json' \
-d '{ "deadlineDate": "2018-07-27T16:00:00+00:00", "title": "First task title", "description": "First task description", "priority": 1}'

curl -v POST \
http://192.168.99.100:8000/todo/ \
-H 'content-type: application/json' \
-d '{ "deadlineDate": "2018-07-27T17:00:00+00:00", "title": "Second task title", "description": "Second task description", "priority": 2}'

curl -v POST \
http://192.168.99.100:8000/todo/ \
-H 'content-type: application/json' \
-d '{ "deadlineDate": "2018-07-30T16:00:00+00:00", "title": "Third task title", "description": "Third task description", "priority": 3}'

curl -v  \
http://192.168.99.100:8000/todo/ \
-H 'content-type: application/json'

curl -v  \
http://192.168.99.100:8000/todo/2 \
-H 'content-type: application/json'

curl -v PUT \
http://192.168.99.100:8000/todo/ \
-H 'content-type: application/json' \
-d '{ "id": 2, "deadlineDate": "2018-07-27T17:00:00+00:00", "title": "Second task title", "description": "Second task description", "priority": 2}'
```