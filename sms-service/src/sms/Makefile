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
	go clean -modcache
	go get github.com/hashicorp/go-rootcerts@6bb64b370b90e7ef1fa532be9e591a81c3493e00
	go clean
	rm -f $(GOPATH)/target/$(BINARY)

test:
	go test -cover ./...

format:
	go fmt ./...

.PHONY: test
