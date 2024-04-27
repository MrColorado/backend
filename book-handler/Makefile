SHELL := /bin/bash
 
GOTAGS ?=
GOFILES ?= $(shell go list ./... | grep -v /vendor/)

GOOS=$(shell go env GOOS)
GOARCH=$(shell go env GOARCH)
GOPATH=$(shell go env GOPATH)

TARGET := $(shell basename ${PWD})
.DEFAULT_GOAL: $(TARGET)

# Use linker flags to provide version/build settings to the target
GOLDFLAGS=-X=main.Version=$(GIT_VERSION) -X=main.Commit=$(GIT_COMMIT)

.PHONY: all bin clean install

all: bin

bin:
	CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags '$(GOLDFLAGS)' -tags '$(GOTAGS)' -o _bin/$(TARGET)

test: test-unit

test-unit:
	@echo '> Run unit tests'
	go test ./...

clean:
	@rm -fr _bin

install:
	CGO_ENABLED=0 go install -a -installsuffix cgo -ldflags '$(GOLDFLAGS)' -tags '$(GOTAGS)'