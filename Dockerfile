FROM golang:alpine as build-env

ENV GO111MODULE=on

RUN apk update && apk add bash ca-certificates git gcc g++ libc-dev

RUN mkdir /chatApp
RUN mkdir -p /chatApp/proto

WORKDIR /chatApp

COPY ./proto/service.pb.go /chatApp/proto
COPY ./main.go /chatApp/

COPY go.mod .
COPY go.sum .

RUN go mod download

RUN go build -o chatApp .

CMD ./chatApp
