GO ?= go
BINARY_NAME = inca
# GOFLAGS := -mod=vendor
STATIC := 1
VERSION := $(shell git describe --abbrev=0 --tags)
GIT_REPO := registry.ekbl.net/sistemi/inca
DOCKER ?= docker
DOCKER_TAG = $(GIT_REPO):$(VERSION)
LDFLAGS = -X main.version=$(VERSION)

# ifeq ($(STATIC), 1)
# LDFLAGS += -s -w -extldflags "-static"
# endif

.PHONY: all

all: build run

clean:
	rm $(BINARY_NAME)

mod-init:
	$(GO) mod init $(GIT_REPO)

mod-deps:
	$(GO) mod tidy
#	$(GO) mod vendor

init: mod-init mod-deps

build: mod-deps
	GOOS=linux GOARCH=amd64 $(GO) build $(GOFLAGS) -ldflags '$(LDFLAGS)' -v -o $(BINARY_NAME) .

docker-build:
	$(DOCKER) build -t $(DOCKER_TAG) .

run: build
	./$(BINARY_NAME)

docker-run:
	$(DOCKER) run -it -v --network host ${PWD}/inca.yml:/etc/inca:ro $(DOCKER_TAG)

docker-push:
	$(DOCKER) push $(DOCKER_TAG)
