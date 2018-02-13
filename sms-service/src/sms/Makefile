GOPATH := $(shell realpath "$(PWD)/../../")
BINARY := sms
DEPENDENCIES := github.com/golang/dep/cmd/dep

export GOPATH ...

all: test build
deploy: test build

build: deps format
	go build -o $(GOPATH)/target/$(BINARY) -v sms.go

clean:
	go clean
	rm -f $(GOPATH)/target/$(BINARY)

test: deps
	go test -cover ./...

format:
	go fmt ./...

deps:
	go get -u $(DEPENDENCIES)
	$(GOPATH)/bin/dep ensure

.PHONY: test
