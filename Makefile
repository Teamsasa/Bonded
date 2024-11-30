.PHONY: help start-all stop-all start-sam-api start-dynamodb dynamodb-init build fmt clean

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
	docker-compose up -d --force-recreate

dynamodb-init: ## Initialize DynamoDB Local using an external script
	@./init-dynamodb.sh

sam-api: ## Start SAM API
	sam local start-api --docker-network bonded_default

build: ## Build SAM application
	go mod tidy                        # 依存関係を整理
	sam build                          # SAMビルドを実行

# 変更点:
# - 不��なコマンドはありませんが、プロジェクトのクリーンアップを推奨します。
# - `cmd/bonded/main.go` を削除するか、別の用途に使用する場合は名前を変更してください。

start-all: compose-up dynamodb-init build sam-api ## Start and initialize DynamoDB, build SAM application, then start SAM API

compose-down: ## Stop and remove Docker containers
	docker-compose down

fmt: ## Format all Go code files
	@go fmt ./...
	@gofmt -s -w .

clean: ## Clean build artifacts
	rm -rf .aws-sam
