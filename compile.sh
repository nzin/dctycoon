#!/bin/bash

go install github.com/shuLhan/go-bindata/v4/cmd/go-bindata@master
$(go env GOPATH)/bin/go-bindata -o global/assets.go -pkg global -ignore *.sh assets/...
# on windows:
# go build -ldflags "-H windowsgui" -o dctycoon ./main

# for macosx, see https://github.com/Xeoncross/macappshell
go build -o dctycoon ./main
