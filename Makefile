GOPATH ?= $(shell go env GOPATH)
PROMU ?= $(GOPATH)/bin/promu

.PHONY: test build
all: test build

test:
	@echo '>> unit test'
	@go test ./...

build:
	@echo '>> build'
	@$(PROMU) build

