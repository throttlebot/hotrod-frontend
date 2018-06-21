FROM golang:1.8
MAINTAINER Hantao Wang

EXPOSE 8080

RUN mkdir -p /go/src/gitlab.com/will.wang1
RUN mkdir -p /go/bin

ADD . /go/src/gitlab.com/will.wang1/hotrod-frontend
WORKDIR /go/src/gitlab.com/will.wang1/hotrod-frontend

RUN go build -o hotrod main.go
RUN mv hotrod /go/bin/

ENTRYPOINT ["/go/bin/hotrod", "frontend"]
