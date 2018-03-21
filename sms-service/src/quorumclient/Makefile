GOPATH := $(shell realpath "$(CURDIR)/../../")
BINARY := quorumclient
PLATFORM := linux
DEPENDENCIES := github.com/golang/dep/cmd/dep

export GOPATH ...

all: test build
deploy: test build

build: deps format
	CGO_ENABLED=0 GOOS=$(PLATFORM) go build -a \
	-ldflags '-extldflags "-static"' \
	-o $(GOPATH)/target/$(BINARY) -v quorumclient.go

clean:
	go clean
	rm -f $(GOPATH)/target/$(BINARY)

test:
	@echo "Yet to Implement"

format:
	go fmt ./...

deps:
	go get -u $(DEPENDENCIES)
	$(GOPATH)/bin/dep ensure

.PHONY: test
