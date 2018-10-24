GOPATH := $(shell realpath "$(CURDIR)/../../")
BINARY := preload
PLATFORM := linux

export GOPATH ...

all: test build
deploy: test build

build:
	CGO_ENABLED=0 GOOS=$(PLATFORM) go build -a \
	-ldflags '-extldflags "-static"' \
	-o $(GOPATH)/target/$(BINARY) -v $(BINARY).go

clean:
	go clean
	rm -f $(GOPATH)/target/$(BINARY)

test:
	@echo "Yet to Implement"

format:
	go fmt ./...

.PHONY: test
