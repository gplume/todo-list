
# TodoList API

Welcome to my todo API attempt built as per your instructions.

### Language and Frameworks used
This API is written in Go which is my preferred backend language.
The Router and Web Framework to make things easier to read and program is [gin-gonic].
The routing part of this framework is the very fast [httprouter] thanks to his(https://github.com/julienschmidt/httprouter). On your request I can write the same API without those helpers but who doesn't use any router which have proven faster than the default mux of Go's net/http package itself?

### Database
I have included a very simple key/value embedded database system ([BoltDB](https://github.com/boltdb/bolt)).
The idea is to use an interface (datamapper) that would make very easy to use another database system without rewriting the logical part of the API.

### EndPoints
All returns endpoints JSON Format.
RestFull methods used:

`GET /todo` : retreive todos list sorted ascending (closer to deadline date first) by default.
There is an optional `sort` parameter implemented:

    `http:locahost/todo?sort=desc` (descending ar `asc' for ascending but is default)

    `http:locahost/todo?sort=priority` (sorted by priority value)
Returns `200` Code


`GET /todo/:key` : retreive specific todo object defined by key (defined by ID in struct and automatically incremented by database). Returns `200` Code

`POST /todo` : Post a correctly structured todo object. Returns `201` Code

`PUT /todo` : Update existing todo object. Returns `200` Code

`DELETE /todo/:key` : Persistently delete todo defined by `:key`. Returns `200` Code



### Todo struct
```
type todo struct {
	ID          int        `json:"id"`
	Creation    time.Time  `json:"creationDate"`
	Deadline    time.Time  `json:"deadlineDate"` // mandatory at insert
	Title       string     `json:"title"`        // mandatory at insert
	Description string     `json:"description"`
	Priority    priorities `json:"priority"`     // mandatory at insert
}

```

### Unity Tests
Lanched as usual by `go test -v`

## Live Tests


```
curl -v -X POST \
http://192.168.99.100:8000/todo \
-H 'content-type: application/json' \
-d '{ "deadlineDate": "2018-07-27T16:00:00+00:00", "title": "First task title", "description": "First task description", "priority": 1}'

curl -v -X POST \
http://192.168.99.100:8000/todo \
-H 'content-type: application/json' \
-d '{ "deadlineDate": "2018-07-27T17:00:00+00:00", "title": "Second task title", "description": "Second task description", "priority": 2}'

curl -v -X POST \
http://192.168.99.100:8000/todo \
-H 'content-type: application/json' \
-d '{ "deadlineDate": "2018-07-30T16:00:00+00:00", "title": "Third task title", "description": "Third task description", "priority": 3}'

curl -v  \
http://192.168.99.100:8000/todo \
-H 'content-type: application/json'

curl -v  \
http://192.168.99.100:8000/todo/2 \
-H 'content-type: application/json'

curl -v PUT \
http://192.168.99.100:8000/todo \
-H 'content-type: application/json' \
-d '{ "id": 2, "deadlineDate": "2018-07-27T17:00:00+00:00", "title": "Second task title", "description": "Second task description", "priority": 2}'

```

windows:
```
docker build -t test .; docker run --rm -it -p 8000:8000 test
docker build -t test .; docker run --mount type=bind,source="/GoProjects/todo-api.log,target=/go/src/gplume/server.log" --rm -it -p 8000:8000 test
```
linux/osx:
```
sudo docker build -t test .; sudo docker run --rm -it -p 8000:8000 test
sudo docker build -t test .; sudo docker run --rm -it -p 8000:8000 -v /appropriate/path/todo-api.log:/go/src/gplume/todo-list/server.log:ro test
```