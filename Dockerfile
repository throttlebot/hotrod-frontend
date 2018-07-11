FROM golang:1.8
MAINTAINER Hantao Wang

EXPOSE 8080

RUN mkdir -p /go/src/github.com/kelda-inc
RUN mkdir -p /go/bin

ADD . /go/src/github.com/kelda-inc/hotrod-frontend
WORKDIR /go/src/github.com/kelda-inc/hotrod-frontend

RUN go build -o hotrod main.go
RUN mv hotrod /go/bin/

ENTRYPOINT ["/go/bin/hotrod", "frontend"]
