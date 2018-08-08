# multi-stage build
FROM golang:1.10.3 as builder
RUN mkdir -p $GOPATH/src/github.com/gplume/todo-list
WORKDIR $GOPATH/src/github.com/gplume/todo-list
COPY . .
RUN go test -v
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o todo-list .

# stage 2
FROM alpine:3.8
WORKDIR /root/
RUN echo $GOPATH
COPY --from=builder /go/src/github.com/gplume/todo-list .
CMD ["./todo-list"]