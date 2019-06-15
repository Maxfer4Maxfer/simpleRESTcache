FROM golang:alpine

ADD ./ /simpleRestCache
WORKDIR /simpleRestCache


RUN apk add git
RUN go mod download
RUN go install -v ./cmd/simplerestcache

ENTRYPOINT ["simplerestcache"]


