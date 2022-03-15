SRCS = $(shell git ls-files '*.go' | grep -v '^vendor/')

APP_NAME := github.com/suecodelabs/cnfuzz
TAG_NAME := $(shell git tag -l --contains HEAD)
SHA := $(shell git rev-parse HEAD)
VERSION_GIT := $(if $(TAG_NAME),$(TAG_NAME),$(SHA))
VERSION := $(if $(VERSION),$(VERSION),$(VERSION_GIT))

BIN_NAME ?= cnfuzz
BIN_DIR ?= dist

GIT_BRANCH := $(subst heads/,,$(shell git rev-parse --abbrev-ref HEAD 2>/dev/null))
DEV_IMAGE := cnfuzz-debug$(if $(GIT_BRANCH),:$(subst /,-,$(GIT_BRANCH)))
IMAGE ?= "cnfuzz"

init:
	mkdir -p $(BIN_DIR)

run:
	go run src/main.go $(RUN_ARGS)

build: init
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(BIN_DIR)/$(BIN_NAME) src/main.go

build-debug:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -gcflags "all=-N -l" -o dist/cnfuzz-debug src/main.go

test:
	go test ./...

clean:
	go clean
	rm -rf $(BIN_DIR)

fmt: format
format:
	gofmt -s -l -w $(SRCS)

image:
	docker build -t $(IMAGE) .

image.local: build
	docker build -t $(IMAGE) -f local.Dockerfile .

image-debug:
	docker build -t $(DEV_IMAGE) -f Dockerfile .

kill-jobs:
	# Kill running jobs
	JOBS=$(shell kubectl get jobs.batch --all-namespaces --no-headers | awk '{if ($$2 ~ "cnfuzz-") print $$2}')
	@if [ $(JOBS) ]; then\
        kubectl delete jobs.batch $$($(JOBS));\
    fi

.PHONY : clean
