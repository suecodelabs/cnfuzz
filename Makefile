SRCS = $(shell git ls-files '*.go' | grep -v '^vendor/')

APP_NAME := ghcr.io/suecodelabs/cnfuzz
TAG_NAME := $(shell git tag -l --contains HEAD)
SHA := $(shell git rev-parse HEAD)
VERSION_GIT := $(if $(TAG_NAME),$(TAG_NAME),$(SHA))
VERSION := $(if $(VERSION),$(VERSION),$(VERSION_GIT))

BIN_DIR ?= dist

GIT_BRANCH := $(subst heads/,,$(shell git rev-parse --abbrev-ref HEAD 2>/dev/null))
GIT_COMMIT := $(subst heads/,,$(shell git rev-parse --short HEAD 2>/dev/null))
DEV_IMAGE := cnfuzz-debug$(if $(GIT_BRANCH),:$(subst /,-,$(GIT_BRANCH)))
CNFUZZ_IMAGE := $(APP_NAME)$(if $(GIT_COMMIT),:$(subst /,-,$(GIT_COMMIT)))
DEFAULT_HELM_DEV_ARGS := --set minio.persistence.size=1Gi,minio.resources.requests.memory=1Gi,minio.replicas=1,minio.mode=standalone --set redis.architecture=standalone,redis.replica.replicaCount=1 --set restler.timeBudget=0.001
KIND_EXAMPLE_IMAGE := $(APP_NAME)$(if $(GIT_COMMIT),-todo-api:$(subst /,-,$(GIT_COMMIT)))
IMAGE ?= "cnfuzz"

init:
	mkdir -p $(BIN_DIR)

helm-init:
	helm repo add bitnami https://charts.bitnami.com/bitnami
	helm repo add minio https://charts.min.io/
	helm dependency build chart/cnfuzz

run:
	go run src/main.go $(RUN_ARGS)

all: cnfuzz restlerwrapper

cnfuzz: init
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(BIN_DIR)/cnfuzz src/cmd/cnfuzz/main.go

cnfuzz-debug: init
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -gcflags "all=-N -l" -o dist/cnfuzz-debug src/main.go

restlerwrapper: init
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(BIN_DIR)/restlerwrapper src/cmd/restlerwrapper/main.go

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
	docker build -t $(IMAGE) -f hack/local.Dockerfile .

image-debug:
	docker build -t $(DEV_IMAGE) -f Dockerfile .

kind-init: build
	cd example && docker build -t $(KIND_EXAMPLE_IMAGE) -f Dockerfile . && cd ..
	docker build -t $(CNFUZZ_IMAGE) -f hack/local.Dockerfile .
	kind load docker-image $(CNFUZZ_IMAGE) && kind load docker-image $(KIND_EXAMPLE_IMAGE)
	helm install --wait --timeout 10m0s dev chart/cnfuzz $(DEFAULT_HELM_DEV_ARGS) $(if $(GIT_COMMIT),--set image.tag=$(subst /,-,$(GIT_COMMIT)))
	kubectl apply -f example/deployment.yaml
	kubectl set image deployment/todo-api todoapi=$(KIND_EXAMPLE_IMAGE)
	kubectl scale deployment --replicas=1 todo-api

kind-build: build
	docker build -t $(CNFUZZ_IMAGE) -f hack/local.Dockerfile .
	kind load docker-image $(CNFUZZ_IMAGE)
	helm upgrade --install dev chart/cnfuzz $(DEFAULT_HELM_DEV_ARGS) $(if $(GIT_COMMIT),--set image.tag=$(subst /,-,$(GIT_COMMIT)))

k8s-clean:
	helm delete dev
	kubectl delete pvc redis-data-dev-redis-master-0
	kubectl delete deployment todo-api

rancher-init: build
	cd example && nerdctl -n k8s.io build -t $(KIND_EXAMPLE_IMAGE) -f Dockerfile . && cd ..
	nerdctl -n k8s.io build -t $(CNFUZZ_IMAGE) -f hack/local.Dockerfile .
	helm install --wait --timeout 10m0s dev chart/cnfuzz $(DEFAULT_HELM_DEV_ARGS) $(if $(GIT_COMMIT),--set image.tag=$(subst /,-,$(GIT_COMMIT)))
	kubectl apply -f example/deployment.yaml
	kubectl set image deployment/todo-api todoapi=$(KIND_EXAMPLE_IMAGE)
	kubectl scale deployment --replicas=1 todo-api

rancher-build: build
	nerdctl -n k8s.io build -t $(CNFUZZ_IMAGE) -f hack/local.Dockerfile .
	helm upgrade --install dev chart/cnfuzz $(DEFAULT_HELM_DEV_ARGS) $(if $(GIT_COMMIT),--set image.tag=$(subst /,-,$(GIT_COMMIT)))

kill-jobs:
	# Kill running jobs
	JOBS=$(shell kubectl get jobs.batch --all-namespaces --no-headers | awk '{if ($$2 ~ "cnfuzz-") print $$2}')
	@if [ $(JOBS) ]; then\
        kubectl delete jobs.batch $$($(JOBS));\
    fi

.PHONY : clean
