FROM golang:latest as builder
WORKDIR /go/src/
COPY simple_server.go .
RUN CGO_ENABLED=0 GOOS=linux go build -o simple_server

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /go/src/simple_server .
CMD ["./simple_server"]
