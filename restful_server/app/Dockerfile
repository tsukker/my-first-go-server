FROM golang:latest as builder
WORKDIR /go/src/
COPY ./restful_server.go .
RUN set -x && \
    go get -u github.com/go-chi/chi && \
    go get github.com/lib/pq && \
    CGO_ENABLED=0 GOOS=linux go build -o restful_server

FROM alpine:latest
WORKDIR /root/
RUN apk --update add tzdata && \
    cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime && \
    apk del tzdata && \
    rm -rf /var/cache/apk/*
COPY --from=builder /go/src/restful_server .
CMD ["./restful_server"]
