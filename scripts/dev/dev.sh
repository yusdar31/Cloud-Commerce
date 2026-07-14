#!/usr/bin/env bash
set -euo pipefail

# CloudCommerce Development Script
# Usage: ./scripts/dev/dev.sh [service-name]
# If no service name is provided, starts all infrastructure only.

SERVICE="${1:-}"

# Always ensure infrastructure is running
echo "Starting infrastructure..."
docker compose up -d postgres redis nats
echo "Infrastructure is running."
echo ""

if [ -z "$SERVICE" ]; then
    echo "No service specified. Starting all services with Docker Compose..."
    docker compose up -d
    echo ""
    echo "All services started!"
    echo ""
    docker compose ps
    exit 0
fi

case "$SERVICE" in
    user|identity|user-service)
        echo "Starting Identity Service (user-service) on port 8081..."
        cd services/user-service
        export DATABASE_URL="postgres://postgres:postgres@localhost:5432/identity_db?sslmode=disable"
        export JWT_SECRET="dev-secret-change-me"
        export PORT="8081"
        go run ./cmd/server
        ;;
    product|catalog|product-service)
        echo "Starting Catalog Service (product-service) on port 8083..."
        cd services/product-service
        export DATABASE_URL="postgres://postgres:postgres@localhost:5432/catalog_db?sslmode=disable"
        export JWT_SECRET="dev-secret-change-me"
        export PORT="8083"
        go run ./cmd/server
        ;;
    gateway|api-gateway)
        echo "Starting API Gateway on port 8080..."
        cd services/api-gateway
        export PORT="8080"
        export JWT_SECRET="dev-secret-change-me"
        export USER_SERVICE_URL="http://localhost:8081"
        export PRODUCT_SERVICE_URL="http://localhost:8083"
        go run ./cmd/server
        ;;
    web|storefront)
        echo "Starting Storefront (Next.js)..."
        pnpm --filter storefront dev
        ;;
    *)
        echo "Unknown service: $SERVICE"
        echo ""
        echo "Available services:"
        echo "  user|identity       - Identity Service (port 8081)"
        echo "  product|catalog     - Catalog Service (port 8083)"
        echo "  gateway             - API Gateway (port 8080)"
        echo "  web|storefront      - Next.js Storefront"
        echo ""
        echo "Or run without arguments to start everything with Docker."
        exit 1
        ;;
esac
