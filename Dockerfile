FROM golang:1.8
MAINTAINER Hantao Wang

EXPOSE 8080

RUN mkdir -p /go/src/github.com/kelda-inc
RUN mkdir -p /go/bin

ADD . /go/src/github.com/kelda-inc/hotrod-frontend
WORKDIR /go/src/github.com/kelda-inc/hotrod-frontend

ENTRYPOINT ["go", "run", "main.go", "frontend"]
