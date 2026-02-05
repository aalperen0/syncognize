include .env

.PHONY: all build clean proto lint test run-gateway run-ingestion run-query run-extraction migrate-up migrate-down docker-build docker-up docker-down

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Binary names
GATEWAY_BINARY=syncognize-gateway
INGESTION_BINARY=syncognize-ingestion
QUERY_BINARY=syncognize-query
EXTRACTION_BINARY=syncognize-extraction

# Directories
CMD_DIR=./cmd
BIN_DIR=./bin
PROTO_DIR=./api/proto

# Build all binaries
all: build

build: build-gateway build-ingestion build-query build-extraction

build-gateway:
	$(GOBUILD) -o $(BIN_DIR)/$(GATEWAY_BINARY) $(CMD_DIR)/gateway

build-ingestion:
	$(GOBUILD) -o $(BIN_DIR)/$(INGESTION_BINARY) $(CMD_DIR)/ingestion

build-query:
	$(GOBUILD) -o $(BIN_DIR)/$(QUERY_BINARY) $(CMD_DIR)/query

build-extraction:
	$(GOBUILD) -o $(BIN_DIR)/$(EXTRACTION_BINARY) $(CMD_DIR)/extraction

# Proto generation
proto:
	buf generate

proto-lint:
	buf lint

proto-breaking:
	buf breaking --against '.git#branch=main'

# Dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Run services
run-gateway:
	$(GOCMD) run $(CMD_DIR)/gateway/main.go

run-ingestion:
	$(GOCMD) run $(CMD_DIR)/ingestion/main.go

run-query:
	$(GOCMD) run $(CMD_DIR)/query/main.go

run-extraction:
	$(GOCMD) run $(CMD_DIR)/extraction/main.go

# Testing
test:
	$(GOTEST) -v -race -cover ./...

test-coverage:
	$(GOTEST) -v -race -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Linting
lint:
	golangci-lint run ./...

# Database migrations
MIGRATE_CMD=migrate

migrate-up:
	$(MIGRATE_CMD) -path internal/adapter/db/migrations -database "$(SYNCOGNIZE_DATABASE_URL)" up

migrate-down:
	$(MIGRATE_CMD) -path internal/adapter/db/migrations -database "$(SYNCOGNIZE_DATABASE_URL)" down

migrate-create:
	$(MIGRATE_CMD) create -ext sql -dir internal/adapter/db/migrations -seq $(name)

# Docker
docker-build:
	docker-compose -f deployments/docker/docker-compose.yml build

docker-up:
	docker-compose -f deployments/docker/docker-compose.yml up -d

docker-down:
	docker-compose -f deployments/docker/docker-compose.yml down

docker-logs:
	docker-compose -f deployments/docker/docker-compose.yml logs -f

# Clean
clean:
	$(GOCLEAN)
	rm -rf $(BIN_DIR)
	rm -f coverage.out coverage.html

# Install development tools
tools:
	go install github.com/bufbuild/buf/cmd/buf@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Help
help:
	@echo "Available targets:"
	@echo "  build           - Build all binaries"
	@echo "  build-gateway   - Build gateway binary"
	@echo "  build-ingestion - Build ingestion binary"
	@echo "  build-query     - Build query binary"
	@echo "  build-extraction- Build extraction binary"
	@echo "  proto           - Generate protobuf code"
	@echo "  proto-lint      - Lint protobuf files"
	@echo "  deps            - Download and tidy dependencies"
	@echo "  run-gateway     - Run gateway service"
	@echo "  run-ingestion   - Run ingestion service"
	@echo "  run-query       - Run query service"
	@echo "  run-extraction  - Run extraction service"
	@echo "  test            - Run tests"
	@echo "  test-coverage   - Run tests with coverage"
	@echo "  lint            - Run linter"
	@echo "  migrate-up      - Run database migrations"
	@echo "  migrate-down    - Rollback database migrations"
	@echo "  migrate-create  - Create new migration (name=<name>)"
	@echo "  docker-build    - Build Docker images"
	@echo "  docker-up       - Start Docker services"
	@echo "  docker-down     - Stop Docker services"
	@echo "  clean           - Clean build artifacts"
	@echo "  tools           - Install development tools"
