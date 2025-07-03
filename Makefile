.PHONY: build run test clean docker-build docker-run docker-compose-up docker-compose-down k8s-deploy k8s-undeploy

# Go commands
build:
	go build -o bin/server cmd/server/main.go

run:
	go run cmd/server/main.go

test:
	go test -v ./...

clean:
	rm -rf bin/

# Generate swagger docs
swagger:
	swag init -g cmd/server/main.go -o docs/

# Docker commands
docker-build:
	docker build -f deployments/docker/Dockerfile -t golang-service:latest .

docker-run:
	docker run -p 8080:8080 --env-file .env golang-service:latest

# Docker Compose commands
docker-compose-up:
	cd deployments/docker && docker-compose up -d

docker-compose-down:
	cd deployments/docker && docker-compose down

docker-compose-logs:
	cd deployments/docker && docker-compose logs -f

# Kubernetes commands
k8s-deploy:
	kubectl apply -f deployments/kubernetes/

k8s-undeploy:
	kubectl delete -f deployments/kubernetes/

k8s-logs:
	kubectl logs -f deployment/golang-service -n golang-service

# Development helpers
dev-setup:
	cp .env.example .env
	go mod download
	@echo "Please update .env file with your actual configuration values"

# Linting
lint:
	golangci-lint run

# Format code
fmt:
	go fmt ./...

# Run all checks
check: fmt lint test

# Help
help:
	@echo "Available commands:"
	@echo "  build              - Build the application"
	@echo "  run                - Run the application"
	@echo "  test               - Run tests"
	@echo "  clean              - Clean build artifacts"
	@echo "  swagger            - Generate swagger documentation"
	@echo "  docker-build       - Build Docker image"
	@echo "  docker-run         - Run Docker container"
	@echo "  docker-compose-up  - Start services with Docker Compose"
	@echo "  docker-compose-down- Stop services with Docker Compose"
	@echo "  k8s-deploy         - Deploy to Kubernetes"
	@echo "  k8s-undeploy       - Remove from Kubernetes"
	@echo "  dev-setup          - Setup development environment"
	@echo "  lint               - Run linter"
	@echo "  fmt                - Format code"
	@echo "  check              - Run all checks"