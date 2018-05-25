GOPATH     ?= $(shell go env GOPATH)
DEP        ?= $(GOPATH)/bin/dep
GORELEASER ?= $(GOPATH)/bin/goreleaser
VERSION    := v$(shell cat VERSION)

.PHONY: setup test build build-snapshot sync-tag release
all: setup test build build-snapshot sync-tag release

setup: $(DEP)
	@echo '>> setup'
	@$(DEP) ensure -v

test:
	@echo '>> unit test'
	@go test ./...

build:
	@echo '>> build'
	@go build -ldflags='\
	-X github.com/kobtea/iapetus/vendor/github.com/prometheus/common/version.Version=$(shell cat VERSION) \
	-X github.com/kobtea/iapetus/vendor/github.com/prometheus/common/version.Revision=$(shell git rev-parse HEAD) \
	-X github.com/kobtea/iapetus/vendor/github.com/prometheus/common/version.Branch=$(shell git symbolic-ref --short HEAD) \
	-X github.com/kobtea/iapetus/vendor/github.com/prometheus/common/version.BuildUser=$(shell whoami)@$(shell hostname) \
	-X github.com/kobtea/iapetus/vendor/github.com/prometheus/common/version.BuildDate=$(shell date +%Y%m%d-%H:%M:%S)' \
	./cmd/iapetus

build-snapshot: $(GORELEASER)
	@echo '>> cross-build for testing'
	BUILD_BRANCH=$(shell git symbolic-ref --short HEAD) \
	BUILD_USER=$(shell whoami) \
	BUILD_HOST=$(shell hostname) \
	BUILD_DATE=$(shell date +%Y%m%d-%H:%M:%S) \
	$(GORELEASER) release --snapshot --rm-dist --debug

sync-tag:
	@git rev-parse $(VERSION) > /dev/null 2>&1 || \
	(git tag -a $(VERSION) -m "release $(VERSION)" && git push origin $(VERSION))

release: $(GORELEASER)
	@echo '>> release'
	BUILD_BRANCH=$(shell git symbolic-ref --short HEAD) \
	BUILD_USER=$(shell whoami) \
	BUILD_HOST=$(shell hostname) \
	BUILD_DATE=$(shell date +%Y%m%d-%H:%M:%S) \
	$(GORELEASER) release --rm-dist --debug

$(DEP):
	go get -u github.com/golang/dep/cmd/dep

$(GORELEASER): $(DEP)
	go get golang.org/x/tools/cmd/stringer
	go get -d github.com/goreleaser/goreleaser
	cd $(GOPATH)/src/github.com/goreleaser/goreleaser && \
	$(DEP) ensure -vendor-only && \
	make build && \
	mv ./goreleaser $(GOPATH)/bin
