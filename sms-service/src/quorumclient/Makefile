GOPATH := $(shell realpath "$(CURDIR)/../../")
BINARY := quorumclient
PLATFORM := linux

export GOPATH ...
export GO111MODULE=on

all: test build
deploy: test build

build: clean
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

.PHONY: test
