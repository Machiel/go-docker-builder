FROM golang:latest

VOLUME /output

ADD . /go/src/go-docker-builder

WORKDIR /go/src/go-docker-builder

ENTRYPOINT ./build.sh
