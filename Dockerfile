FROM golang:1.8
MAINTAINER Hantao Wang

EXPOSE 8080

RUN mkdir -p /go/src/gitlab.com/kelda-hotrod
RUN mkdir -p /go/bin
RUN go get github.com/go-redis/redis
RUN go get github.com/lib/pq
RUN go get github.com/sirupsen/logrus

WORKDIR /go/src/gitlab.com/kelda-hotrod

ARG git_pass
ARG build_time=1

RUN git clone https://user:$git_pass@gitlab.com/will.wang1/hotrod-base
RUN git clone https://user:$git_pass@gitlab.com/will.wang1/hotrod-route
RUN git clone https://user:$git_pass@gitlab.com/will.wang1/hotrod-frontend
RUN git clone https://user:$git_pass@gitlab.com/will.wang1/hotrod-customer
RUN git clone https://user:$git_pass@gitlab.com/will.wang1/hotrod-driver

WORKDIR /go/src/gitlab.com/kelda-hotrod/hotrod-frontend

RUN go build -o hotrod main.go
RUN mv hotrod /go/bin/

ENTRYPOINT ["/go/bin/hotrod", "frontend"]
