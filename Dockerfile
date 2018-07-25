# multi-stage build
FROM golang:1.10.3 as builder
RUN mkdir -p $GOPATH/src/gplume/todo-list
RUN echo $GOPATH
WORKDIR $GOPATH/src/gplume/todo-list
COPY . .
RUN go test -v
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o todo-list .

# stage 2
FROM alpine:3.8
WORKDIR /root/
RUN echo $GOPATH
COPY --from=builder /go/src/gplume/todo-list .
CMD ["./todo-list"]