VERSION := $(shell git describe --tags)
REVISION := $(shell git rev-parse --short HEAD)

BINARY_NAME := cryptorious

LDFLAGS := -X github.com/malnick/cryptorious/config.VERSION=$(VERSION) -X github.com/malnick/cryptorious/config.REVISION=$(REVISION) 

FILES := $(shell go list ./... | grep -v vendor)

all: test install

test:
	@echo "+$@"
	go test $(FILES)  -cover

build: 
	@echo "+$@"
	go build -v -o cryptorious_$(VERSION) -ldflags '$(LDFLAGS)' cryptorious.go

install:
	@echo "+$@"
	go install -v -ldflags '$(LDFLAGS)' $(FILES)
