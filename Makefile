SRCS = $(shell git ls-files '*.go' | grep -v '^vendor/')

CNFUZZ_EXT_IMG := ghcr.io/suecodelabs/cnfuzz
WRAPPER_EXT_IMG := ghcr.io/suecodelabs/cnfuzz-restlerwrapper

TAG_NAME := $(shell git tag -l --contains HEAD)
SHA := $(shell git rev-parse HEAD)
VERSION_GIT := $(if $(TAG_NAME),$(TAG_NAME),$(SHA))
VERSION := $(if $(VERSION),$(VERSION),$(VERSION_GIT))

GO_ENV_VARS ?= CGO_ENABLED=0 GOOS=linux GOARCH=amd64

BIN_DIR ?= dist
CNFUZZ_DOCKERFILE ?= "src/cmd/cnfuzz/Dockerfile"
CNFUZZ_LOCAL_DOCKERFILE ?= "src/cmd/cnfuzz/local.Dockerfile"
CNFUZZ_IMAGE ?= "cnfuzz"
RESTLERWRAPPER_IMAGE ?= "restlerwrapper"
RESTLERWRAPPER_DOCKERFILE ?= "src/cmd/restlerwrapper/Dockerfile"
RESTLERWRAPPER_LOCAL_DOCKERFILE ?= "src/cmd/restlerwrapper/local.Dockerfile"
EXAMPLE_API_IMAGE := cnfuzz-todo-api

DEFAULT_HELM_DEV_ARGS := --set controllerImage.repository=$(CNFUZZ_IMAGE),controllerImage.tag=latest,restlerwrapper.image.image=$(RESTLERWRAPPER_IMAGE),restlerwrapper.image.tag=latest --set minio.persistence.size=1Gi,minio.resources.requests.memory=1Gi,minio.replicas=1,minio.mode=standalone --set redis.architecture=standalone,redis.replica.replicaCount=1 --set restler.timeBudget=0.001 --set debugMode=true

init:
	mkdir -p $(BIN_DIR)

helm-init:
	helm repo add bitnami https://charts.bitnami.com/bitnami
	helm repo add minio https://charts.min.io/
	helm dependency build chart/cnfuzz

run:
	go run src/main.go $(RUN_ARGS)

test:
	go test ./...

clean:
	go clean
	rm -rf $(BIN_DIR)

fmt: format
format:
	gofmt -s -l -w $(SRCS)

all: cnfuzz restlerwrapper

cnfuzz: init
	$(GO_ENV_VARS) go build -o $(BIN_DIR)/cnfuzz src/cmd/cnfuzz/main.go

cnfuzz-debug: init
	$(GO_ENV_VARS) go build -gcflags "all=-N -l" -o $(BIN_DIR)/cnfuzz-debug src/main.go

restlerwrapper: init
	$(GO_ENV_VARS) go build -o $(BIN_DIR)/restlerwrapper src/cmd/restlerwrapper/main.go

cnfuzz-image:
	docker build -t $(CNFUZZ_IMAGE) -f $(CNFUZZ_DOCKERFILE) --no-cache .

cnfuzz-image.local: cnfuzz
	docker build -t $(CNFUZZ_IMAGE) -f $(CNFUZZ_LOCAL_DOCKERFILE) .

restlerwrapper-image:
	docker build -t $(RESTLERWRAPPER_IMAGE) -f $(RESTLERWRAPPER_DOCKERFILE) .

restlerwrapper-image.local: restlerwrapper
	docker build -t $(RESTLERWRAPPER_IMAGE) -f $(RESTLERWRAPPER_LOCAL_DOCKERFILE) .

kind-init: kind-load-images kind-fuzz-test

kind-load-images: all
	cd example && docker build -t $(EXAMPLE_API_IMAGE) -f Dockerfile . && cd ..
	docker build -t $(CNFUZZ_IMAGE) -f $(CNFUZZ_LOCAL_DOCKERFILE) .
	docker build -t $(RESTLERWRAPPER_IMAGE) -f $(RESTLERWRAPPER_LOCAL_DOCKERFILE) .
	kind load docker-image $(CNFUZZ_IMAGE) && kind load docker-image $(RESTLERWRAPPER_IMAGE) && kind load docker-image $(EXAMPLE_API_IMAGE)

kind-fuzz-test:
	helm install --wait --timeout 10m0s dev chart/cnfuzz $(DEFAULT_HELM_DEV_ARGS) # $(if $(GIT_COMMIT),--set image.tag=$(subst /,-,$(GIT_COMMIT)))
	kubectl apply -f example/deployment.yaml
	kubectl set image deployment/todo-api todoapi=$(EXAMPLE_API_IMAGE)
	kubectl scale deployment --replicas=1 todo-api

kind-build: all
	docker build -t $(CNFUZZ_IMAGE) -f $(CNFUZZ_LOCAL_DOCKERFILE) .
	docker build -t $(RESTLERWRAPPER_IMAGE) -f $(RESTLERWRAPPER_LOCAL_DOCKERFILE) .
	kind load docker-image $(CNFUZZ_IMAGE) && kind load docker-image $(RESTLERWRAPPER_IMAGE)
	helm upgrade --install dev chart/cnfuzz $(DEFAULT_HELM_DEV_ARGS) # $(if $(GIT_COMMIT),--set image.tag=$(subst /,-,$(GIT_COMMIT)))

kind-setup-dev: kind-load-images
	echo build example api image
	echo deploy example api
	@kubectl apply -f example/deployment.yaml
	@kubectl set image deployment/todo-api todoapi=$(EXAMPLE_API_IMAGE)
	@kubectl scale deployment --replicas=1 todo-api

# kind-load-ext-images:
# 	docker pull

k8s-clean:
	helm delete dev
	kubectl delete pvc redis-data-dev-redis-master-0
	kubectl delete deployment todo-api

rancher-init: rancher-load-images
	helm install --wait --timeout 10m0s dev chart/cnfuzz $(DEFAULT_HELM_DEV_ARGS) # $(if $(GIT_COMMIT),--set image.tag=$(subst /,-,$(GIT_COMMIT)))
	kubectl apply -f example/deployment.yaml
	kubectl set image deployment/todo-api todoapi=$(EXAMPLE_API_IMAGE)
	kubectl scale deployment --replicas=1 todo-api

rancher-build: all
	nerdctl -n k8s.io build -t $(CNFUZZ_IMAGE) -f $(CNFUZZ_LOCAL_DOCKERFILE) .
	nerdctl -n k8s.io build -t $(RESTLERWRAPPER_IMAGE) -f $(RESTLERWRAPPER_LOCAL_DOCKERFILE) .
	helm upgrade --install dev chart/cnfuzz $(DEFAULT_HELM_DEV_ARGS) # $(if $(GIT_COMMIT),--set image.tag=$(subst /,-,$(GIT_COMMIT)))

rancher-load-images: all
	cd example && nerdctl -n k8s.io build -t $(EXAMPLE_API_IMAGE) -f Dockerfile . && cd ..
	nerdctl -n k8s.io build -t $(CNFUZZ_IMAGE) -f $(CNFUZZ_LOCAL_DOCKERFILE) .
	nerctl -n k8s.io build -t $(RESTLERWRAPPER_IMAGE) -f $(RESTLERWRAPPER_LOCAL_DOCKERFILE) .

.PHONY : clean cnfuzz restlerwrapper
