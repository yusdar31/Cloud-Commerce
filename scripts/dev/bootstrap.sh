#!/usr/bin/env bash
set -euo pipefail

# CloudCommerce Development Bootstrap Script
# Usage: ./scripts/dev/bootstrap.sh

echo "=========================================="
echo "  CloudCommerce Development Bootstrap"
echo "=========================================="
echo ""

# Check prerequisites
echo "[1/5] Checking prerequisites..."

if ! command -v go &> /dev/null; then
    echo "  ✗ Go is not installed. Please install Go 1.25+"
    exit 1
fi
echo "  ✓ Go: $(go version)"

if ! command -v node &> /dev/null; then
    echo "  ✗ Node.js is not installed. Please install Node.js 22+"
    exit 1
fi
echo "  ✓ Node.js: $(node --version)"

if ! command -v pnpm &> /dev/null; then
    echo "  ✗ pnpm is not installed. Installing..."
    npm install -g pnpm@9.15.4
fi
echo "  ✓ pnpm: $(pnpm --version)"

if ! command -v docker &> /dev/null; then
    echo "  ✗ Docker is not installed. Please install Docker Desktop."
    exit 1
fi
echo "  ✓ Docker: $(docker --version)"

echo ""

# Start infrastructure
echo "[2/5] Starting infrastructure (PostgreSQL, Redis, NATS)..."
docker compose up -d postgres redis nats
echo "  ✓ Infrastructure started"

echo ""

# Wait for PostgreSQL to be ready
echo "[3/5] Waiting for PostgreSQL..."
until docker exec cc-postgres pg_isready -U postgres &> /dev/null; do
    sleep 1
    echo -n "."
done
echo ""
echo "  ✓ PostgreSQL is ready"

echo ""

# Run migrations
echo "[4/5] Running database migrations..."

if command -v goose &> /dev/null; then
    DATABASE_URL="postgres://postgres:postgres@localhost:5432/identity_db?sslmode=disable" goose -dir services/user-service/migrations postgres "$DATABASE_URL" up
    echo "  ✓ identity_db migrations applied"

    DATABASE_URL="postgres://postgres:postgres@localhost:5432/catalog_db?sslmode=disable" goose -dir services/product-service/migrations postgres "$DATABASE_URL" up
    echo "  ✓ catalog_db migrations applied"
else
    echo "  ! goose not installed. Skipping migrations."
    echo "  ! Install with: go install github.com/pressly/goose/v3/cmd/goose@latest"
fi

echo ""

# Install frontend dependencies
echo "[5/5] Installing frontend dependencies..."
pnpm install
echo "  ✓ Dependencies installed"

echo ""
echo "=========================================="
echo "  Bootstrap complete!"
echo "=========================================="
echo ""
echo "Next steps:"
echo "  1. Start backend services:"
echo "     cd services/user-service && go run ./cmd/server"
echo "     cd services/product-service && go run ./cmd/server"
echo ""
echo "  2. Start frontend:"
echo "     pnpm dev"
echo ""
echo "  3. Or start everything with Docker:"
echo "     docker compose up -d"
echo ""
