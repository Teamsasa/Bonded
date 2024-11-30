.PHONY: help start-all stop-all start-sam-api start-dynamodb dynamodb-init build fmt

# Default target
.DEFAULT_GOAL := help

# Colors for output
GREEN  := $(shell tput setaf 2)
YELLOW := $(shell tput setaf 3)
RESET  := $(shell tput sgr0)

help: ## Display this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-20s$(RESET) %s\n", $$1, $$2}'
	@echo ''

compose-up: ## Start Docker containers
	docker-compose up -d

dynamodb-init: ## Initialize DynamoDB Local using an external script
	@./init-dynamodb.sh

sam-api: ## Start SAM API
	sam local start-api --docker-network bonded_default

start-all: compose-up dynamodb-init sam-api ## Start and initialize DynamoDB, then start SAM API

compose-down: ## Stop and remove Docker containers
	docker-compose down

build: ## Build SAM application
	sam build

fmt: ## Format all Go code files
	@go fmt ./...
