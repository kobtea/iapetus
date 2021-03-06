GOPATH     ?= $(shell go env GOPATH)
GORELEASER ?= $(GOPATH)/bin/goreleaser
VERSION    := v$(shell cat VERSION)

.PHONY: test build build-snapshot sync-tag release docker-build docker-release
all: test build build-snapshot sync-tag release docker-build docker-release

test:
	@echo '>> unit test'
	@go test ./...

build:
	@echo '>> build'
	@go build -ldflags='\
	-X github.com/kobtea/iapetus/vendor/github.com/prometheus/common/version.Version=$(shell cat VERSION) \
	-X github.com/kobtea/iapetus/vendor/github.com/prometheus/common/version.Revision=$(shell git rev-parse HEAD) \
	-X github.com/kobtea/iapetus/vendor/github.com/prometheus/common/version.Branch=$(shell git rev-parse --abbrev-ref HEAD) \
	-X github.com/kobtea/iapetus/vendor/github.com/prometheus/common/version.BuildUser=$(shell whoami)@$(shell hostname) \
	-X github.com/kobtea/iapetus/vendor/github.com/prometheus/common/version.BuildDate=$(shell date +%Y%m%d-%H:%M:%S)' \
	./cmd/iapetus

build-snapshot: $(GORELEASER)
	@echo '>> cross-build for testing'
	BUILD_BRANCH=$(shell git rev-parse --abbrev-ref HEAD) \
	BUILD_USER=$(shell whoami) \
	BUILD_HOST=$(shell hostname) \
	BUILD_DATE=$(shell date +%Y%m%d-%H:%M:%S) \
	$(GORELEASER) release --snapshot --rm-dist --debug

sync-tag:
	@git config user.name  || git config --local user.name  "circleci-job"
	@git config user.email || git config --local user.email "kobtea9696@gmail.com"
	@git rev-parse $(VERSION) > /dev/null 2>&1 || \
	(git tag -a $(VERSION) -m "release $(VERSION)" && git push origin $(VERSION))

release: $(GORELEASER)
	@echo '>> release'
	BUILD_BRANCH=$(shell git rev-parse --abbrev-ref HEAD) \
	BUILD_USER=$(shell whoami) \
	BUILD_HOST=$(shell hostname) \
	BUILD_DATE=$(shell date +%Y%m%d-%H:%M:%S) \
	$(GORELEASER) release --rm-dist --debug

docker-build:
	@echo '>> build docker image'
	@docker build -t kobtea/iapetus:$(shell cat VERSION) .
	@docker build -t kobtea/iapetus:latest .

docker-release: docker-build
	@echo '>> release docker image'
	@docker login -u ${DOCKERHUB_USER} -p ${DOCKERHUB_PASS}
	@docker push kobtea/iapetus:$(shell cat VERSION)
	@docker push kobtea/iapetus:latest

$(GORELEASER):
	@wget -O - "https://github.com/goreleaser/goreleaser/releases/download/v0.172.1/goreleaser_$(shell uname -o | cut -d'/' -f2)_$(shell uname -m).tar.gz" | tar xvzf - -C /tmp
	@mv /tmp/goreleaser $(GOPATH)/bin
