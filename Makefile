.PHONY: up down sam-api dynamodb build fmt

# Start DynamoDB Local
dynamodb:
	docker-compose up -d

# Initialize DynamoDB Local using an external script
dynamodb-init:
	@./init-dynamodb.sh

# Start SAM API (change network name)
sam-api:
	sam local start-api --docker-network bonded_default

# Start both DynamoDB and SAM API
up: dynamodb dynamodb-init sam-api

# Stop and remove Docker containers
down:
	docker-compose down

# Build SAM application
build:
	sam build

# Format Go code
fmt:
	@find . -name "*.go" -exec go fmt {} +
