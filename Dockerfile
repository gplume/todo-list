FROM golang:1.10.3-alpine3.8
EXPOSE 8000

# copy the app to a proper build path to import vendor directory
RUN mkdir -p $GOPATH/src/gplume/todo-list
COPY . $GOPATH/src/gplume/todo-list
WORKDIR $GOPATH/src/gplume/todo-list

RUN go test -v
RUN go build -o todo-list .
CMD ["./todo-list"]