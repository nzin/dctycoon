FROM golang:1.9-alpine3.7
MAINTAINER Jordi Riera <kender.jr@gmail.com>

RUN apk add --no-cache \
    git \
    gcc \
    cmake \
    build-base \
    libx11-dev \
    pkgconf \
    sdl2-dev \
    sdl2_ttf-dev \
    sdl2_image-dev


WORKDIR /go/src/github.com/nzin/dctycoon/
COPY . .

RUN go get -u github.com/jteeuwen/go-bindata/...
RUN "$(go env GOPATH)/bin/go-bindata" -o global/assets.go -pkg global assets/...

RUN go get ./...
RUN go build ./...
