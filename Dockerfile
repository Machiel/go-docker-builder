FROM golang:latest

VOLUME /output

ADD . /go/src/go-docker-builder

WORKDIR /go/src/go-docker-builder

RUN go get github.com/fsouza/go-dockerclient

ENTRYPOINT ./build.sh
