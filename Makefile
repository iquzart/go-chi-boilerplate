.DEFAULT_GOAL := help
.PHONY: help clean build run docker-build docker-run docker-stop compose-up compose-down swagger

# Define variables
BINARY_NAME := go-chi-boilerplate 
DOCKER_IMAGE := go-chi-boilerplate
DOCKER_CONTAINER := go-chi-boilerplate
MAIN_FILE := cmd/api/main.go
LDFLAGS := -ldflags="-s -w"

help: ## Display this help message
	@echo "Usage: make <command>"
	@echo ""
	@echo "Commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

clean: ## Remove binary file and container resources
	rm -f $(BINARY_NAME)
	docker stop $(DOCKER_CONTAINER) || true
	docker rm $(DOCKER_CONTAINER) || true
	docker rmi $(DOCKER_IMAGE) || true

build: ## Build the binary file
	go build $(LDFLAGS) -o $(BINARY_NAME) $(MAIN_FILE)

run: swagger ## Build and run the binary file locally
	go run $(LDFLAGS) $(MAIN_FILE)

docker-build: ## Build Docker image
	docker build -t $(DOCKER_IMAGE) . -f ./Containerfile

docker-run: ## Run Docker container
	@if ! docker image inspect $(DOCKER_IMAGE) >/dev/null 2>&1; then \
		make docker-build; \
	fi
	docker run -d -p 8080:8080 --name $(DOCKER_CONTAINER) $(DOCKER_IMAGE)

docker-stop: ## Stop Docker container
	docker stop $(DOCKER_CONTAINER) || true
	docker rm $(DOCKER_CONTAINER) || true

compose-up: ## Start services using docker-compose
	docker-compose up -d --build

compose-down: ## Stop services using docker-compose
	docker-compose down

swagger: ## Generate swagger 2.0 and convert it to OpenAPI v3
	swag init -g cmd/api/main.go -o docs
	docker run --rm -v $(PWD)/docs:/docs shaelmaar/swagger2openapi:node-12.12.0 \
	  swagger2openapi -o /docs/v3/openapi.json /docs/swagger.json

