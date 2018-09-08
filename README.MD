
# TodoList API

Welcome to my todo API built as per your instructions.

## Language and Frameworks used
This API is written in Go which is my preferred backend language.
The Router and Web Framework used here to make things easier to read and program is [gin-gonic](https://github.com/gin-gonic/gin).
The routing part of this framework is the very fast [httprouter](https://github.com/julienschmidt/httprouter), especially great for API specifications. On your request I can write the same API without those helpers but who doesn't use any custom router which have proven to be faster than the default mux of Go's net/http package itself?

Go Version needs to be >= 1.8

### Database
Is included a simple key/value embedded database system ([bolt](https://github.com/boltdb/bolt)).
One idea is to wrap an interface (`datamapper`) that make very easy to use another database system without rewriting the logical part of the API.
User should, of course, think about linking an external volume for the database file(s) when running the production container.

## EndPoints
All endpoints returns JSON and eventual errors have the "error" key first.
RestFull methods used:

	GET /todo	 : retreive full todos list sorted ascending (closer to deadline date first) by default.
There is an optional `sort` parameter implemented:

    http://locahost:8000/todo?sort=desc (descending or 'asc' for ascending but it is default)

    http://locahost:8000/todo?sort=priority (sorted by priority value)

...

	GET /todo/:id : retreive specific todo object defined by id. Returns 200 Status Code

	POST /todo : Post a correctly json structured todo object. Returns 201 Status Code

	PUT /todo : Update existing todo object. Returns 200 Status Code

	DELETE /todo/:id : Persistently delete todo defined by `:id`. Returns 204 Status Code



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

## Logs

All actions and errors are logged both to Stdout and a local `server.log` file.

## Unity Tests

Coverage is not 100% yet.

Run Dockerized tests only: (will stop at stage #1 build)
```
docker build --target builder -t test .
```


## Live Tests
A `live-tests.sh` script is included to launch a series of 'live' test with `curl` that fill the database with four tasks, add one and modify it.

example of deletion:
```
curl -v -X DELETE http://localhost:8000/todo/1
```

Windows users using the legacy docker-toolbox should add the `win` parameter to the script: `./live-tests.sh win` to retreive the IP of the server (which is not localhost but more something like: `192.168.99.100`).

## Instrumentation
[Prometheus](https://github.com/prometheus/client_golang) client library is used and `/metrics` should be used by a Prometheus server to consume data. Displays nicely with [Grafana](https://grafana.com/).

Requests counters and durations are implemented.

## DOCKER Commands
linux/osx:
```
sudo docker build -t todo-list-api . && sudo docker run --rm -it -p 8000:8000 todo-list-api

```
windows (git terminal):
```
docker build -t todo-list-api .
docker run --rm -it -p 8000:8000 todo-list-api
```
Of course the -d flag should be added if you wish to 'deamonize' or 'detach' the application and let it run.

## Minikube Commands

```
minikube start --vm-driver="hyperv" --insecure-registry="$REG_IP":80
minikube start --vm-driver="hyperv" --hyperv-virtual-switch "Primary Virtual Switch"

eval $(minikube docker-env)

docker build -t todo-list-api:1.0 .

kubectl create --validate -f kube-api.yml

check:
minikube dashboard
or:
kubectl get pods
kubectl get services -o wide

```
Play with the API:
```
curl -v http://$(minikube ip):31317/todo
...
...
````
Cleanup:
```
kubectl delete -f kube-api.yml
minikube stop
eval $(minikube docker-env -u)
minikube delete
````