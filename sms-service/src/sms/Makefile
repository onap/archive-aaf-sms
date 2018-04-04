GOPATH := $(shell realpath "$(CURDIR)/../../")
BINARY := sms
PLATFORM := linux
DEPENDENCIES := github.com/golang/dep/cmd/dep

export GOPATH ...

all: test build
deploy: test build

build: deps format
	CGO_ENABLED=0 GOOS=$(PLATFORM) go build -a \
	-ldflags '-extldflags "-static"' \
	-o $(GOPATH)/target/$(BINARY) -v sms.go

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
