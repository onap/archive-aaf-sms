GOPATH := $(shell realpath "$(CURDIR)/../../")
BINARY := sms
PLATFORM := linux

export GOPATH ...
export GO111MODULE=on

all: test build
deploy: test build

build: clean
	CGO_ENABLED=0 GOOS=$(PLATFORM) go build -a \
	-ldflags '-extldflags "-static"' \
	-o $(GOPATH)/target/$(BINARY) -v sms.go

clean:
	go clean
	rm -f $(GOPATH)/target/$(BINARY)

test:
	go test -cover ./...

format:
	go fmt ./...

.PHONY: test
