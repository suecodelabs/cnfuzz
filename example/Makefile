IMAGE ?= todo-api
BIN_DIR ?= dist
BIN_NAME ?= todo-api
SWAG_BIN := $(shell go env GOPATH)/bin/swag

swagger: swag
swag:
	$(SWAG_BIN) init

build: swag
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(BIN_DIR)/$(BIN_NAME) ./main.go

image:
	docker build -t $(IMAGE) .

docker.start:
	docker run --rm -p 8080:8080 --name temp-todo-api $(IMAGE)
