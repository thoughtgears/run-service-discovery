ifneq (,$(wildcard .env))
include .env
export $(shell sed 's/=.*//' .env)
endif

GIT_SHA := $(shell git rev-parse --short HEAD)
DOCKER_REPO := $(GCP_REGION)-docker.pkg.dev/$(GCP_PROJECT_ID)/services
SERVICE_NAME := service-discovery

.PHONY: clean dev build docker-build

clean:
	rm -rf builds/*

dev:
	@echo "Running service discovery development mode"
	@GCP_PROJECT_ID=$(GCP_PROJECT_ID) GIN_MODE=debug go run main.go

build: clean
	@echo "Building service discovery"
	GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags "-s -w" -o builds/$(SERVICE_NAME)-linux-amd64 .
	GOOS=linux GOARCH=arm64 go build -a -installsuffix cgo -ldflags "-s -w" -o builds/$(SERVICE_NAME)-linux-arm64 .
	GOOS=darwin GOARCH=amd64 go build -a -installsuffix cgo -ldflags "-s -w" -o builds/$(SERVICE_NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build -a -installsuffix cgo -ldflags "-s -w" -o builds/$(SERVICE_NAME)-darwin-arm64 .

docker-build:
	@echo "Building service discovery docker image"
	docker build -t $(DOCKER_REPO)/$(SERVICE_NAME):$(GIT_SHA) .
	docker tag $(DOCKER_REPO)/$(SERVICE_NAME):$(GIT_SHA) $(DOCKER_REPO)/$(SERVICE_NAME):latest
