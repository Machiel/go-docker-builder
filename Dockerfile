FROM golang:latest

VOLUME /output

ADD /Users/Machiel/Documents/Projects/go/src/go-docker-builder /go/src/go-docker-builder

WORKDIR /go/src/go-docker-builder

RUN go get github.com/fsouza/go-dockerclient

ENTRYPOINT ./build.sh
