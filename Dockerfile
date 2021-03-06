FROM golang:alpine3.7 as builder

LABEL maintainer='Hugo Fonseca <https://github.com/hugomcfonseca>'

WORKDIR /go/src/github.com/hugomcfonseca/go-simple-webserver/

COPY /app/server.go .

RUN apk add --update --no-cache git \
    && go get -d -v \
    && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o webserver .

FROM alpine:3.7

LABEL maintainer='Hugo Fonseca <https://github.com/hugomcfonseca>'

ENV LISTEN_PORT='10000' \
    \
    CONFD_VERSION='0.15.0' \
    CONFD_OPTS='-backend=env'

COPY --from=builder /go/src/github.com/hugomcfonseca/go-simple-webserver/webserver /usr/local/bin/

RUN chmod u+x /usr/local/bin/webserver

VOLUME /confs

ENTRYPOINT [ "webserver" ]
