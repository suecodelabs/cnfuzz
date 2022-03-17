IMAGE ?= todo-api
BIN_DIR ?= dist
BIN_NAME ?= todo-api
SWAG_BIN := $(shell go env GOPATH)/bin/swag
KIND_IMAGE := $(APP_NAME)$(if $(GIT_COMMIT),:$(subst /,-,$(GIT_COMMIT)))
DEPLOYMENT_NAME := todo-api

swagger: swag
swag:
	$(SWAG_BIN) init

build: swag
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(BIN_DIR)/$(BIN_NAME) ./main.go

image:
	docker build -t $(IMAGE) .

docker.start:
	docker run --rm -p 8080:8080 --name temp-todo-api $(IMAGE)

kind: build
	docker build -t $(KIND_IMAGE) -f Dockerfile .
	kind load docker-image $(KIND_IMAGE)
	helm upgrade --install $(DEPLOYMENT_NAME) ./chart --set image.tag=$(subst /,-,$(GIT_COMMIT))

kind-clean:
	helm uninstall $(DEPLOYMENT_NAME)
