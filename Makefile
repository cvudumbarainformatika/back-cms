.PHONY: help install run build test clean fmt lint deps-update migrate migrate-fresh setup-git reset-starter docker-build docker-run docker-stop

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
	@echo "  make migrate        Run all database migrations"
	@echo "  make migrate-fresh  Drop all tables and re-run migrations"
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

migrate:
	@echo "Running database migrations..."
	@if [ ! -f .env ]; then \
		echo "Error: .env file not found"; \
		exit 1; \
	fi
	@echo "Loading environment variables..."
	@export $$(grep -v '^#' .env | grep -v '^$$' | xargs) && \
	echo "Connecting to Docker MySQL container..." && \
	for file in database/migrations/*.sql; do \
		echo "Migrating: $$file"; \
		docker exec -i $$(docker ps -qf "name=mysql") mysql -u$$DB_USERNAME -p$$DB_PASSWORD $$DB_DATABASE < $$file || exit 1; \
	done
	@echo "✓ All migrations completed successfully"

migrate-fresh:
	@echo "WARNING: This will drop all tables and re-run migrations!"
	@echo "Press Ctrl+C to cancel, or Enter to continue..."
	@read confirm
	@echo "Dropping all tables..."
	@export $$(grep -v '^#' .env | grep -v '^$$' | xargs) && \
	docker exec -i $$(docker ps -qf "name=mysql") mysql -u$$DB_USERNAME -p$$DB_PASSWORD $$DB_DATABASE -e "SET FOREIGN_KEY_CHECKS = 0; \
		SELECT CONCAT('DROP TABLE IF EXISTS \`', table_name, '\`;') \
		FROM information_schema.tables \
		WHERE table_schema = '$$DB_DATABASE'; \
		SET FOREIGN_KEY_CHECKS = 1;" | grep "DROP TABLE" | docker exec -i $$(docker ps -qf "name=mysql") mysql -u$$DB_USERNAME -p$$DB_PASSWORD $$DB_DATABASE || true
	@echo "Running fresh migrations..."
	@$(MAKE) migrate

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
	docker compose build
	@echo "✓ Docker image built"

docker-run:
	@echo "Starting Docker containers..."
	docker compose up -d
	@echo "✓ Containers started"

docker-stop:
	@echo "Stopping Docker containers..."
	docker compose down
	@echo "✓ Containers stopped"
