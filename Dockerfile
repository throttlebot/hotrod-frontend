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
ARG build_time
RUN git clone https://user:$git_pass@gitlab.com/kelda-hotrod/hotrod-base
RUN git clone https://user:$git_pass@gitlab.com/kelda-hotrod/hotrod-route
RUN git clone https://user:$git_pass@gitlab.com/kelda-hotrod/hotrod-frontend
RUN git clone https://user:$git_pass@gitlab.com/kelda-hotrod/hotrod-customer
RUN git clone https://user:$git_pass@gitlab.com/kelda-hotrod/hotrod-driver

WORKDIR /go/src/gitlab.com/kelda-hotrod/hotrod-frontend

RUN go build -o hotrod main.go
RUN mv hotrod /go/bin/

ENTRYPOINT ["/go/bin/hotrod", "frontend"]
