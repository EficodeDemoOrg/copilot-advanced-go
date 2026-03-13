.PHONY: install dev build start test test-unit test-integration test-e2e lint fmt swagger help

## install: Download Go dependencies, install dev tools, and install e2e dependencies
install:
	go mod download
	go install github.com/air-verse/air@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest
	cd tests/e2e && npm install

## dev: Start the server with hot reload (requires air)
dev:
	air

## build: Compile the server binary
build:
	go build -o bin/server ./cmd/server

## start: Run the compiled binary
start:
	./bin/server

## test: Run all unit and integration tests
test:
	go test ./...

## test-unit: Run unit tests only (converters, models, repository, services)
test-unit:
	go test ./internal/utils/... ./internal/models/... ./internal/repository/... ./internal/services/...

## test-integration: Run integration tests only (handlers via httptest)
test-integration:
	go test ./internal/handlers/...

## test-e2e: Run Playwright browser tests (requires Node.js + app running)
test-e2e:
	cd tests/e2e && npx playwright test

## lint: Run golangci-lint
lint:
	golangci-lint run

## fmt: Format all Go files
fmt:
	gofmt -w .
	goimports -w .

## swagger: Regenerate OpenAPI spec from handler annotations
swagger:
	swag init -g cmd/server/main.go --output docs

## help: Show available make targets
help:
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'
