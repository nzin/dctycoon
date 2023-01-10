FROM golang:1.18-alpine
MAINTAINER Jordi Riera <kender.jr@gmail.com>

RUN apk add --no-cache \
    curl \
    git \
    gcc \
    cmake \
    build-base \
    libx11-dev \
    pkgconf \
    sdl2-dev \
    sdl2_ttf-dev \
    sdl2_image-dev \
    libjpeg


WORKDIR /go/src/github.com/nzin/dctycoon/
COPY . .
RUN go install github.com/shuLhan/go-bindata/v4/cmd/go-bindata@master && \
    go get -u github.com/stretchr/testify/assert && \
    go get github.com/axw/gocov/gocov && \
    curl -L https://codeclimate.com/downloads/test-reporter/test-reporter-latest-linux-amd64 > ./cc-test-reporter && \
    chmod +x ./cc-test-reporter

ENV CC_TEST_REPORTER_ID=bacedd92b18dc3389348470c6d536f7250336e28f571f389a302c67b68a3096d

RUN "$(go env GOPATH)/bin/go-bindata" -o global/assets.go -pkg global assets/... && \
    go get ./... && \
    go build ./...
