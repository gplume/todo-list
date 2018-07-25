#!/bin/bash
CT="Content-Type: application/json"
CONTENT1=$(cat <<EOF
{
    "deadlineDate": "2018-07-27T16:00:00+00:00",
    "title": "First task title",
    "description": "Description of the first task",
    "priority": 1
}
EOF
)
CONTENT2=$(cat <<EOF
{
    "deadlineDate": "2018-07-27T17:00:00+00:00",
    "title": "Second task title",
    "description": "Description of the Second task",
    "priority": 2
}
EOF
)
CONTENT3=$(cat <<EOF
{
    "deadlineDate": "2018-07-28T19:00:00+00:00",
    "title": "Third task title",
    "description": "Description of the Third task",
    "priority": 3
}
EOF
)
CONTENT4=$(cat <<EOF
{
    "deadlineDate": "2018-07-27T18:00:00+00:00",
    "title": "Fourth task title",
    "description": "Description of the Fourth task",
    "priority": 1
}
EOF
)

UPDATED=$(cat <<EOF
{
    "id": 2,
    "deadlineDate": "2018-07-27T18:00:00+00:00",
    "title": "Second task updated title",
    "description": "Description of the updated Second task",
    "priority": 2
}
EOF
)

HOST=http://localhost
if [ $1 = "win" ]; then
    HOST=http://`docker-machine ip default`
fi
PORT=8000

curl -v -X POST $HOST:$PORT/todo -H "$CT" -d "$CONTENT1"
curl -v -X POST $HOST:$PORT/todo -H "$CT" -d "$CONTENT2"
curl -v -X POST $HOST:$PORT/todo -H "$CT" -d "$CONTENT3"
curl -v -X POST $HOST:$PORT/todo -H "$CT" -d "$CONTENT4"

sleep 1

curl -v $HOST:$PORT/todo -H "$CT"
curl -v $HOST:$PORT/todo/1 -H "$CT"
curl -v $HOST:$PORT/todo/2 -H "$CT"
curl -v $HOST:$PORT/todo/3 -H "$CT"
curl -v $HOST:$PORT/todo/4 -H "$CT"



curl -v -X PUT $HOST:$PORT/todo -H "$CT" -d "$UPDATED"

sleep 1

# if following run multiple time produce errors because some ids gets deleted:
# curl -v -X DELETE $HOST:$PORT/todo/4 -H "$CT"
