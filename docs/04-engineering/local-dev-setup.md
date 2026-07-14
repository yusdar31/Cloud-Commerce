# Local Development Setup

**Project:** CloudCommerce

**Last Updated:** July 2026

---

## 1. Prerequisites

Install tools berikut sebelum memulai:

| Tool | Version | Install |
|------|---------|---------|
| Go | 1.25+ | https://go.dev/dl/ |
| Node.js | 22+ | https://nodejs.org/ |
| pnpm | 9+ | `npm install -g pnpm` |
| Docker Desktop | Latest | https://www.docker.com/products/docker-desktop/ |
| Git | Latest | https://git-scm.com/ |

**Opsional tapi sangat disarankan:**

| Tool | Kegunaan | Install |
|------|----------|---------|
| `goose` | Database migration | `go install github.com/pressly/goose/v3/cmd/goose@latest` |
| `golangci-lint` | Go linter | https://golangci-lint.run/usage/install/ |
| `make` | Task runner (jika pakai Makefile) | Sudah ada di macOS/Linux; Windows: chocolatey `choco install make` |

---

## 2. Setup Pertama Kali

### Clone Repository

```bash
git clone <repo-url> cloudcommerce
cd cloudcommerce
```

### Install Frontend Dependencies

```bash
# Di root project (pnpm workspace)
pnpm install
```

### Setup Environment Variables

Setiap app/service membutuhkan file `.env`. Buat dari template berikut:

#### `apps/storefront/.env.local`

```env
NEXT_PUBLIC_API_URL=http://localhost:8080
NEXT_PUBLIC_APP_NAME=CloudCommerce
```

#### Service Go (buat file `.env` di masing-masing service)

Semua service Go membaca dari environment variables. Untuk dev lokal, set via terminal atau buat `.env` dan load manual.

Lihat: [env-variables-reference.md](env-variables-reference.md) untuk daftar lengkap semua ENV.

---

## 3. Jalankan Infrastruktur Lokal

Infrastruktur (database, cache, message broker) dijalankan via Docker Compose.

### Start Infrastruktur Saja

```bash
# Di root project
docker compose up postgres redis nats -d
```

### Verifikasi Infrastruktur Berjalan

```bash
docker compose ps
```

Output yang diharapkan:

```
NAME              STATUS          PORTS
cc-postgres       running         0.0.0.0:5432->5432/tcp
cc-redis          running         0.0.0.0:6379->6379/tcp
cc-nats           running         0.0.0.0:4222->4222/tcp, 0.0.0.0:8222->8222/tcp
```

### Cek Database

```bash
# Connect ke postgres
docker exec -it cc-postgres psql -U postgres -l
```

Database yang tersedia setelah init:

| Database | Owner Service |
|----------|---------------|
| `identity_db` | user-service |
| `tenant_db` | store-service |
| `catalog_db` | product-service |
| `inventory_db` | inventory-service |
| `order_db` | order-service |
| `payment_db` | payment-service |
| `notification_db` | notification-service |
| `review_db` | review-service |

---

## 4. Jalankan Backend Services

### Product Service (Port 8083)

```bash
cd services/product-service

# Set environment variables
export PORT=8083
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/catalog_db?sslmode=disable"
export REDIS_URL="redis://localhost:6379"
export JWT_SECRET="dev-secret-change-me"
export APP_ENV="development"

# Jalankan migration dulu
goose -dir migrations postgres "$DATABASE_URL" up

# Jalankan service
go run ./cmd/server
```

Atau set semuanya dalam satu baris:

```bash
PORT=8083 DATABASE_URL="postgres://postgres:postgres@localhost:5432/catalog_db?sslmode=disable" REDIS_URL="redis://localhost:6379" JWT_SECRET="dev-secret-change-me" go run ./cmd/server
```

### User Service (Port 8081)

```bash
cd services/user-service
PORT=8081 DATABASE_URL="postgres://postgres:postgres@localhost:5432/identity_db?sslmode=disable" REDIS_URL="redis://localhost:6379" JWT_SECRET="dev-secret-change-me" JWT_EXPIRY="24h" go run ./cmd/server
```

### Store Service (Port 8082)

```bash
cd services/store-service
PORT=8082 DATABASE_URL="postgres://postgres:postgres@localhost:5432/tenant_db?sslmode=disable" JWT_SECRET="dev-secret-change-me" go run ./cmd/server
```

### API Gateway (Port 8080)

```bash
cd services/api-gateway
PORT=8080 JWT_SECRET="dev-secret-change-me" NATS_URL="nats://localhost:4222" REDIS_URL="redis://localhost:6379" USER_SERVICE_URL="http://localhost:8081" PRODUCT_SERVICE_URL="http://localhost:8083" go run ./cmd/server
```

---

## 5. Jalankan Frontend

```bash
cd apps/storefront
pnpm dev
```

Storefront tersedia di: **http://localhost:3000**

---

## 6. Jalankan Semua dengan Docker Compose

Jika ingin menjalankan seluruh stack termasuk semua service:

```bash
# Build & jalankan semua
docker compose up --build -d

# Lihat log semua service
docker compose logs -f

# Lihat log service tertentu
docker compose logs -f product-service

# Stop semua
docker compose down

# Stop + hapus data (fresh start)
docker compose down -v
```

> **⚠️ Catatan:** Saat menggunakan Docker Compose, storefront Next.js **tidak** termasuk.
> Jalankan `pnpm dev` terpisah untuk frontend.

---

## 7. Database Migration

CloudCommerce menggunakan **Goose** untuk migration.

### Jalankan Migration

```bash
cd services/product-service

# Jalankan semua migration pending
goose -dir migrations postgres "postgres://postgres:postgres@localhost:5432/catalog_db?sslmode=disable" up

# Rollback migration terakhir
goose -dir migrations postgres "postgres://postgres:postgres@localhost:5432/catalog_db?sslmode=disable" down

# Lihat status migration
goose -dir migrations postgres "postgres://postgres:postgres@localhost:5432/catalog_db?sslmode=disable" status
```

### Format File Migration

```
migrations/
├── 00001_init.sql
├── 00002_add_indexes.sql
└── 00003_add_column.sql
```

---

## 8. Health Check

Setelah service berjalan, verifikasi dengan:

```bash
# Product Service
curl http://localhost:8083/healthz

# User Service
curl http://localhost:8081/healthz

# API Gateway
curl http://localhost:8080/healthz
```

Response yang diharapkan:

```json
{ "status": "ok", "service": "product-service" }
```

---

## 9. Testing

### Unit Test (Go)

```bash
cd services/product-service
go test ./internal/...
```

### Integration Test (Go — butuh database)

```bash
# Pastikan Docker Compose infrastruktur sudah running
cd services/product-service
go test ./tests/...
```

### Frontend Test

```bash
cd apps/storefront
pnpm test        # Unit test (Vitest)
pnpm test:e2e    # E2E (Playwright)
```

---

## 10. Troubleshooting

### PostgreSQL tidak bisa connect

```bash
# Cek apakah container running
docker compose ps

# Cek logs postgres
docker compose logs postgres

# Restart postgres
docker compose restart postgres
```

### Port sudah digunakan

```bash
# Windows — cek proses di port 5432
netstat -ano | findstr :5432

# Kill process (ganti PID)
taskkill /PID <PID> /F
```

### Go module tidak ditemukan

```bash
# Di root services/{service-name}
go mod tidy

# Jika menggunakan shared-go local
# Pastikan replace directive ada di go.mod:
# replace github.com/cloudcommerce/shared-go => ../../packages/shared-go
```

### pnpm install error

```bash
# Clear cache dan reinstall
pnpm store prune
pnpm install
```

### Next.js build error

```bash
cd apps/storefront
rm -rf .next
pnpm dev
```

---

## 11. VS Code Setup (Rekomendasi)

Extensions yang disarankan:

- **Go** (golang.go)
- **ESLint** (dbaeumer.vscode-eslint)
- **Prettier** (esbenp.prettier-vscode)
- **Tailwind CSS IntelliSense** (bradlc.vscode-tailwindcss)
- **Docker** (ms-azuretools.vscode-docker)
- **REST Client** (humao.rest-client)

---

## 12. Related Documents

- [Technology Stack](01-tech-stack.md)
- [Monorepo Structure](02-monorepo-structure.md)
- [Coding Standards](03-coding-standards.md)
- [Backend Guidelines](04-backend-guidelines.md)
- [Frontend Guidelines](05-frontend-guidelines.md)
- [Environment Variables Reference](env-variables-reference.md)
