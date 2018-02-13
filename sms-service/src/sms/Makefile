GOPATH := $(shell realpath "$(PWD)/../../")
BINARY := sms
DEPENDENCIES := github.com/golang/dep/cmd/dep

export GOPATH ...

all: build test
deploy: test build

build: deps format test
	$(GOPATH)/bin/dep ensure
	go build -o $(GOPATH)/target/$(BINARY) -v sms.go

clean:
	go clean
	rm -f $(GOPATH)/target/$(BINARY)

test:
	go test -cover ./...

format:
	go fmt ./...

deps:
	go get -u $(DEPENDENCIES)

.PHONY: test
