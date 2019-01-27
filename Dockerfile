FROM golang:1.9.1
MAINTAINER sabatm144@gmail.com

COPY ./server /go/src/sample/server
RUN go get -d -v ./... 