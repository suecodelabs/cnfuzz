SRCS = $(shell git ls-files '*.go' | grep -v '^vendor/')

APP_NAME := ghcr.io/suecodelabs/cnfuzz
TAG_NAME := $(shell git tag -l --contains HEAD)
SHA := $(shell git rev-parse HEAD)
VERSION_GIT := $(if $(TAG_NAME),$(TAG_NAME),$(SHA))
VERSION := $(if $(VERSION),$(VERSION),$(VERSION_GIT))

BIN_NAME ?= cnfuzz
BIN_DIR ?= dist

GIT_BRANCH := $(subst heads/,,$(shell git rev-parse --abbrev-ref HEAD 2>/dev/null))
GIT_COMMIT := $(subst heads/,,$(shell git rev-parse --short HEAD 2>/dev/null))
DEV_IMAGE := cnfuzz-debug$(if $(GIT_BRANCH),:$(subst /,-,$(GIT_BRANCH)))
KIND_IMAGE := $(APP_NAME)$(if $(GIT_COMMIT),:$(subst /,-,$(GIT_COMMIT)))
KIND_EXAMPLE_IMAGE := $(APP_NAME)$(if $(GIT_COMMIT),todo-api:$(subst /,-,$(GIT_COMMIT)))
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

kind: build
    cd example && docker build -t $(KIND_EXAMPLE_IMAGE) -f Dockerfile . && cd ..
	docker build -t $(KIND_IMAGE) -f local.Dockerfile .
	kind load docker-image $(KIND_IMAGE)
	kind load docker-image $(KIND_EXAMPLE_IMAGE)
	kubectl apply -f example/deployment.yaml
	kubectl set image deployment/todo-api todoapi=$(KIND_EXAMPLE_IMAGE)
	helm install --wait --timeout 10m0s $(if $(GIT_COMMIT),cnfuzz-$(subst /,-,$(GIT_COMMIT))) charts/cnfuzz $(if $(GIT_COMMIT),--set image.tag=$(subst /,-,$(GIT_COMMIT)))

kind-clean:
	helm delete $(if $(GIT_COMMIT),cnfuzz-$(subst /,-,$(GIT_COMMIT)))

kill-jobs:
	# Kill running jobs
	JOBS=$(shell kubectl get jobs.batch --all-namespaces --no-headers | awk '{if ($$2 ~ "cnfuzz-") print $$2}')
	@if [ $(JOBS) ]; then\
        kubectl delete jobs.batch $$($(JOBS));\
    fi

.PHONY : clean
