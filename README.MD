
# TodoList API

Welcome to my todo API built as per your instructions.

## Language and Frameworks used
This API is written in Go which is my preferred backend language.
The Router and Web Framework to make things easier to read and program is [gin-gonic](https://github.com/gin-gonic/gin).
The routing part of this framework is the very fast [httprouter](https://github.com/julienschmidt/httprouter). On your request I can write the same API without those helpers but who doesn't use any router which have proven faster than the default mux of Go's net/http package itself?

Go Version needs to be >= 1.8

### Database
I have included a simple key/value embedded database system ([bolt](https://github.com/boltdb/bolt)).
One idea is to wrap an interface (`datamapper`) that would make very easy to use another database system without rewriting the logical part of the API.
Don't forget to link an external volume for the database file(s) when running production container!

## EndPoints
All endpoints returns JSON Format.
RestFull methods used:

`GET /todo` : retreive todos list sorted ascending (closer to deadline date first) by default.
There is an optional `sort` parameter implemented:

    `http://locahost:8000/todo?sort=desc` (descending or `asc' for ascending but it is default)

    `http://locahost:8000/todo?sort=priority` (sorted by priority value)
	Returns `200` Status Code


`GET /todo/:key` : retreive specific todo object defined by key. Returns `200` Status Code

`POST /todo` : Post a correctly structured todo object. Returns `201` Status Code

`PUT /todo` : Update existing todo object. Returns `200` Status Code

`DELETE /todo/:key` : Persistently delete todo defined by `:key`. Returns `200` Status Code



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

## Errors

All errors are logged both to Stdout and a local `server.log` file. Of course it would be better to bind this file externally at container start.

## Unity Tests
Launched as usual by `go test -v`

Run Dockerized tests only: (will stop at stage #1 build)
```
docker build --target builder -t test .
```


## Live Tests
`live-tests.sh` script is included to lauch a series of 'live' test with `curl` that fill the database with four tasks, add one and modify it.

example of deletion:
```
curl -v -X DELETE -H "Content-Type: application/json" http://localhost:8000/todo/1
```

Windows users using the legacy docker-toolbox should add the `win` parameter to the script: `./live-tests.sh win` to retreive the IP of the server (which is not localhost but more something like: `192.168.99.100`).

## Instrumentation
[Prometheus](https://github.com/prometheus/client_golang) client library is used and `/metrics` should be used by a Prometheus server to consume data. Displays nicely with [Grafana](https://grafana.com/).

Requests counters and duration are implemented.

## DOCKER Commands
linux/osx:
```
sudo docker build -t todo-list-api . && sudo docker run --rm -it -p 8000:8000 todo-list-api

```
windows (git terminal):
```
docker build -t todo-list-api .; docker run --rm -it -p 8000:8000 todo-list-api
```
