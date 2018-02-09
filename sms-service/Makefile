GOPATH=$(shell pwd)
BINARY=sms
DEPENDENCIES := github.com/golang/dep/cmd/dep

all: build test
deploy: build test

build: deps format
	cd src/sms && $(GOPATH)/bin/dep ensure
	go build -o target/$(BINARY) -v src/sms/sms.go

clean:
	go clean
	rm -f target/$(BINARY)

test:
	go test -v ./src/sms/...

format:
	go fmt ./src/sms/...

deps:
	go get -u $(DEPENDENCIES)