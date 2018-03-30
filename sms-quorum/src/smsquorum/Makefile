GOPATH := $(shell realpath "$(PWD)/../../")
BINARY := quorumclient
PLATFORM := linux
#DEPENDENCIES := github.com/golang/dep/cmd/dep

export GOPATH ...

all: test build
deploy: test build

build: format
	CGO_ENABLED=0 GOOS=$(PLATFORM) go build -a \
	-ldflags '-extldflags "-static"' \
	-o $(GOPATH)/target/$(BINARY) -v quorumclient.go

clean:
	go clean
	rm -f $(GOPATH)/target/$(BINARY)

test: #deps
	@echo "Yet to Implement"

format:
	go fmt ./...

.PHONY: test
