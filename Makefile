
.PHONY: help build test clean up down logs fmt vet security-scan coverage benchmark

# Default target
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Development
build: ## Build the Go application
	@echo "Building Go application..."
	go build -o main ./cmd/server

test: ## Run tests
	@echo "Running tests..."
	go test -v ./...

coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

benchmark: ## Run benchmarks
	@echo "Running benchmarks..."
	go test -bench=. -benchmem ./...

fmt: ## Format Go code
	@echo "Formatting Go code..."
	go fmt ./...

vet: ## Run go vet
	@echo "Running go vet..."
	go vet ./...

security-scan: ## Run security scan with gosec
	@echo "Running security scan..."
	@if ! command -v gosec > /dev/null; then \
		echo "Installing gosec..."; \
		go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest; \
	fi
	gosec ./...

# Docker operations
docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t web-server-go:latest .

docker-run: ## Run Docker container
	@echo "Running Docker container..."
	docker run -p 8080:8080 --name web-server-go web-server-go:latest

docker-clean: ## Clean Docker images and containers
	@echo "Cleaning Docker resources..."
	docker stop web-server-go 2>/dev/null || true
	docker rm web-server-go 2>/dev/null || true
	docker rmi web-server-go:latest 2>/dev/null || true

# Docker Compose operations
up: ## Start all services with docker-compose
	@echo "Starting services..."
	docker-compose up -d

down: ## Stop all services
	@echo "Stopping services..."
	docker-compose down

restart: down up ## Restart all services

logs: ## Show logs from all services
	docker-compose logs -f

logs-web: ## Show logs from web service only
	docker-compose logs -f web

logs-prometheus: ## Show logs from prometheus service only
	docker-compose logs -f prometheus

logs-grafana: ## Show logs from grafana service only
	docker-compose logs -f grafana

# Monitoring setup
setup-monitoring: ## Initialize monitoring infrastructure
	@echo "Setting up monitoring infrastructure..."
	@./monitoring/scripts/setup.sh

monitoring-up: up ## Start full monitoring stack (same as 'make up')
	@echo ""
	@echo "✅ Full monitoring stack started from unified docker-compose!"
	@echo ""
	@echo "🌐 Access URLs:"
	@echo "📊 Grafana: http://localhost:3000 (admin/admin123)"
	@echo "🔍 Prometheus: http://localhost:9090"
	@echo "🚨 Alertmanager: http://localhost:9093"
	@echo "📈 Node Exporter: http://localhost:9100"
	@echo "🔍 Blackbox Exporter: http://localhost:9115"
	@echo "🚀 Web Server: http://localhost:8080"

monitoring-down: down ## Stop monitoring stack (same as 'make down')
	@echo "✅ Full monitoring stack stopped!"

monitoring-logs: ## Show monitoring logs
	@docker-compose logs -f prometheus grafana alertmanager node-exporter blackbox-exporter

monitoring-status: ## Show monitoring stack status
	@echo "📊 Monitoring Stack Status:"
	@docker-compose ps
	@echo ""
	@echo "🎯 Prometheus Targets Health:"
	@curl -s http://localhost:9090/api/v1/targets 2>/dev/null | grep -o '"health":"[^"]*"' | sort | uniq -c || echo "⚠️  Prometheus not ready yet"
	@echo ""
	@echo "🚨 Active Alerts:"
	@curl -s http://localhost:9090/api/v1/alerts 2>/dev/null | grep -o '"state":"[^"]*"' | sort | uniq -c || echo "⚠️  No alerts data"

monitoring-check: ## Quick health check of all endpoints
	@echo "🏥 Health Check Results:"
	@echo "Web Server: $$(curl -s http://localhost:8080/health -w "%{http_code}" -o /dev/null 2>/dev/null)"
	@echo "Prometheus: $$(curl -s http://localhost:9090/-/healthy -w "%{http_code}" -o /dev/null 2>/dev/null)"
	@echo "Grafana: $$(curl -s http://localhost:3000/api/health -w "%{http_code}" -o /dev/null 2>/dev/null)"
	@echo "Alertmanager: $$(curl -s http://localhost:9093/api/v2/status -w "%{http_code}" -o /dev/null 2>/dev/null)"
	@echo "Node Exporter: $$(curl -s http://localhost:9100/metrics -w "%{http_code}" -o /dev/null 2>/dev/null)"

monitoring-reload: ## Reload Prometheus configuration
	@echo "🔄 Reloading Prometheus configuration..."
	@curl -X POST http://localhost:9090/-/reload 2>/dev/null && echo "✅ Prometheus config reloaded" || echo "❌ Failed to reload"

# Monitoring access
prometheus: ## Open Prometheus UI
	@echo "Opening Prometheus at http://localhost:9090"
	@which xdg-open > /dev/null && xdg-open http://localhost:9090 || echo "Open http://localhost:9090 manually"

grafana: ## Open Grafana UI
	@echo "Opening Grafana at http://localhost:3000 (admin/admin123)"
	@which xdg-open > /dev/null && xdg-open http://localhost:3000 || echo "Open http://localhost:3000 manually"

alertmanager: ## Open Alertmanager UI
	@echo "Opening Alertmanager at http://localhost:9093"
	@which xdg-open > /dev/null && xdg-open http://localhost:9093 || echo "Open http://localhost:9093 manually"

# Health checks
health: ## Check application health
	@echo "Checking application health..."
	@curl -s http://localhost:8080/health | jq '.' || echo "Health check failed"

metrics: ## Show application metrics
	@echo "Application metrics:"
	@curl -s http://localhost:8080/metrics | jq '.' || echo "Metrics fetch failed"

# Cleanup
clean: docker-clean ## Clean build artifacts and Docker resources
	@echo "Cleaning build artifacts..."
	rm -f main coverage.out coverage.html

# Development workflow
dev: fmt vet test ## Run development checks (format, vet, test)

ci: fmt vet test security-scan coverage ## Run all CI checks

# Production deployment
prod-build: ## Build production Docker image
	@echo "Building production image..."
	docker build -t web-server-go:prod --target builder .
	docker build -t web-server-go:prod .

prod-deploy: prod-build ## Deploy to production (placeholder)
	@echo "Production deployment would go here..."
	@echo "Consider using Kubernetes, AWS ECS, or similar orchestration"
