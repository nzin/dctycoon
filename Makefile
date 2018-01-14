ifndef $(GOPATH)
    GOPATH=$(shell go env GOPATH)
    export GOPATH
endif

all:
	go get -u github.com/jteeuwen/go-bindata/...
	$(GOPATH)/bin/go-bindata -o global/assets.go -pkg global -ignore *.sh assets/...
	go get ./...
	go build -o dctycoon ./main
