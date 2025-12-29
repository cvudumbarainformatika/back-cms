.PHONY: help install run build test clean fmt lint deps-update setup-git reset-starter docker-build docker-run docker-stop

help:
	@echo "Go Gin Starter Kit - Available Commands"
	@echo ""
	@echo "  make install        Install dependencies"
	@echo "  make run            Run the application"
	@echo "  make build          Build executable"
	@echo "  make test           Run tests"
	@echo "  make clean          Remove build artifacts"
	@echo "  make fmt            Format code"
	@echo "  make lint           Run linter"
	@echo "  make deps-update    Update dependencies"
	@echo "  make setup-git      Initialize git repository"
	@echo "  make reset-starter  Reset to clean starter kit state"
	@echo "  make docker-build   Build Docker image"
	@echo "  make docker-run     Start Docker containers"
	@echo "  make docker-stop    Stop Docker containers"
	@echo ""

install:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy
	@echo "✓ Dependencies installed"

run:
	@echo "Running application..."
	go run main.go

build:
	@echo "Building application..."
	@mkdir -p bin
	go build -o bin/main main.go
	@echo "✓ Built: bin/main"

test:
	@echo "Running tests..."
	go test -v -cover ./...

clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	go clean
	@echo "✓ Cleaned"

fmt:
	@echo "Formatting code..."
	gofmt -s -w .
	@echo "✓ Formatted"

lint:
	@echo "Running linter..."
	@command -v golangci-lint >/dev/null 2>&1 || (echo "golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest" && exit 1)
	golangci-lint run ./...

deps-update:
	@echo "Updating dependencies..."
	go get -u ./...
	go mod tidy
	@echo "✓ Dependencies updated"

setup-git:
	@echo "Setting up git repository..."
	@chmod +x init_starter_git.sh
	./init_starter_git.sh

reset-starter:
	@echo "Resetting to clean starter kit..."
	@chmod +x clean_starter.sh
	./clean_starter.sh

docker-build:
	@echo "Building Docker image..."
	docker-compose build
	@echo "✓ Docker image built"

docker-run:
	@echo "Starting Docker containers..."
	docker-compose up -d
	@echo "✓ Containers started"

docker-stop:
	@echo "Stopping Docker containers..."
	docker-compose down
	@echo "✓ Containers stopped"
