SHELL          := /bin/bash
PROGRAM        := utsusemi
VERSION        := v0.1.0
GOOS           := $(shell go env GOOS)
GOARCH         := $(shell go env GOARCH)
RUNTIME_GOPATH := $(GOPATH):$(shell pwd)
TEST_SRC       := $(wildcard src/*/*_test.go) $(wildcard src/*/test_*.go)
SRC            := $(filter-out $(TEST_SRC),$(wildcard src/*/*.go))

UBUNTU_IMAGE          := docker-go-pkg-build-ubuntu
UBUNTU_CONTAINER_NAME := docker-go-pkg-build-ubuntu-$(shell date +%s)

.PHONY: all
all: $(PROGRAM)

.PHONY: go-get
go-get:
	go get github.com/BurntSushi/toml

$(PROGRAM): $(SRC)
ifeq ($(GOOS),linux)
	GOPATH=$(RUNTIME_GOPATH) CGO_ENABLED=0 go build -ldflags "-X utsusemi.version=$(VERSION)" -a -tags netgo -installsuffix netgo -o $(PROGRAM)
else
	GOPATH=$(RUNTIME_GOPATH) CGO_ENABLED=0 go build -ldflags "-X utsusemi.version=$(VERSION)" -o $(PROGRAM)
endif

.PHONY: test
test: $(TEST_SRC)
	GOPATH=$(RUNTIME_GOPATH) go test -v $(TEST_SRC)

.PHONY: clean
clean: $(TEST_SRC)
	rm -f $(PROGRAM)

.PHONY: package
package: clean test $(PROGRAM)
	gzip -c $(PROGRAM) > pkg/$(PROGRAM)-$(VERSION)-$(GOOS)-$(GOARCH).gz
	rm -f $(PROGRAM)

.PHONY: package/linux
package/linux:
	docker run \
	  --name $(UBUNTU_CONTAINER_NAME) \
	  -v $(shell pwd):/tmp/src $(UBUNTU_IMAGE) \
	  make -C /tmp/src go-get package
	docker rm $(UBUNTU_CONTAINER_NAME)

.PHONY: docker/build/ubuntu
docker/build/ubuntu: etc/Dockerfile.ubuntu
	docker build -f etc/Dockerfile.ubuntu -t $(UBUNTU_IMAGE) .
