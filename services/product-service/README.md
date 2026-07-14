# Product Service

**Service:** product-service (Catalog Service)

**Port:** `8083`

**Database:** `catalog_db` (PostgreSQL)

**Owner:** Backend Team

---

## Responsibility

Mengelola semua data produk dan kategori per merchant (multi-tenant).

- CRUD produk: buat, baca, update, hapus produk
- Manajemen kategori produk
- Publish / archive produk
- Cache produk di Redis
- Publish event `ProductCreated`, `ProductPublished`, `ProductArchived` ke NATS

---

## Technology

| Component | Technology |
|-----------|------------|
| Language | Go 1.25 |
| Framework | Gin |
| Database | PostgreSQL 17 via pgx/v5 |
| Cache | Redis |
| Migration | Goose |
| Shared | `packages/shared-go` |

---

## Project Structure

```
product-service/
├── cmd/
│   └── server/
│       └── main.go              # Entry point, dependency injection
├── internal/
│   ├── domain/
│   │   ├── entity.go            # Product, Category struct + domain errors
│   │   └── repository.go        # Repository interfaces
│   ├── application/
│   │   └── service.go           # Business logic (use cases)
│   ├── infrastructure/
│   │   └── postgres/
│   │       └── repository.go    # SQL query implementations
│   └── transport/
│       ├── handler.go           # Gin HTTP handlers
│       ├── routes.go            # Route setup
│       └── dto.go               # Request/Response DTOs
├── migrations/
│   └── 00001_init.sql           # Database schema
├── Dockerfile
└── go.mod
```

---

## API Endpoints

Base URL: `http://localhost:8083`

All endpoints require: `Authorization: Bearer <JWT>`

### Products

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/v1/products` | List semua produk (tenant dari JWT) |
| `POST` | `/api/v1/products` | Buat produk baru |
| `GET` | `/api/v1/products/:id` | Detail produk |
| `PUT` | `/api/v1/products/:id` | Update produk |
| `DELETE` | `/api/v1/products/:id` | Soft delete produk |
| `POST` | `/api/v1/products/:id/publish` | Publish produk (draft → published) |
| `POST` | `/api/v1/products/:id/archive` | Archive produk |

### Categories

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/v1/categories` | List kategori |
| `POST` | `/api/v1/categories` | Buat kategori |
| `GET` | `/api/v1/categories/:id` | Detail kategori |
| `DELETE` | `/api/v1/categories/:id` | Hapus kategori |

### Health

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/healthz` | Health check |
| `GET` | `/` | Service info |

---

## Environment Variables

| Variable | Wajib | Default (Dev) | Deskripsi |
|----------|-------|---------------|-----------|
| `PORT` | ✅ | `8083` | Port HTTP |
| `DATABASE_URL` | ✅ | - | PostgreSQL connection string |
| `REDIS_URL` | ✅ | - | Redis connection string |
| `JWT_SECRET` | ✅ | - | JWT signing secret |
| `JWT_EXPIRY` | ❌ | `24h` | Token expiry |
| `APP_ENV` | ✅ | `development` | Environment |
| `CACHE_TTL_SECONDS` | ❌ | `60` | TTL cache produk di Redis |

---

## Running Locally

### 1. Pastikan infrastruktur berjalan

```bash
# Di root project
docker compose up postgres redis -d
```

### 2. Set environment & jalankan

```bash
cd services/product-service

export PORT=8083
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/catalog_db?sslmode=disable"
export REDIS_URL="redis://localhost:6379"
export JWT_SECRET="dev-secret-change-me"
export APP_ENV="development"

# Jalankan migration
goose -dir migrations postgres "$DATABASE_URL" up

# Jalankan service
go run ./cmd/server
```

### 3. Verifikasi

```bash
curl http://localhost:8083/healthz
# → {"status":"ok","service":"product-service"}
```

---

## Database

**Database:** `catalog_db`

### Tables

| Table | Description |
|-------|-------------|
| `products` | Data produk utama |
| `categories` | Kategori produk |
| `product_images` | Gambar produk (multiple) |
| `product_tags` | Tag produk |

### Schema Highlights

```sql
-- products table
CREATE TABLE products (
    id          TEXT PRIMARY KEY,          -- ULID format: prd_xxx
    tenant_id   TEXT NOT NULL,             -- WAJIB ada untuk multi-tenant
    name        TEXT NOT NULL,
    slug        TEXT NOT NULL,
    sku         TEXT NOT NULL,
    price       BIGINT NOT NULL,           -- dalam unit terkecil (IDR: sen)
    currency    TEXT NOT NULL DEFAULT 'IDR',
    status      TEXT NOT NULL DEFAULT 'draft',  -- draft | published | archived
    category_id TEXT REFERENCES categories(id),
    image_url   TEXT,
    weight      INTEGER DEFAULT 0,         -- dalam gram
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ,               -- soft delete
    UNIQUE (tenant_id, sku),               -- SKU unik per tenant
    UNIQUE (tenant_id, slug)               -- Slug unik per tenant
);
```

---

## Migration Commands

```bash
cd services/product-service

# Lihat status
goose -dir migrations postgres "$DATABASE_URL" status

# Jalankan semua pending migration
goose -dir migrations postgres "$DATABASE_URL" up

# Rollback 1 migration
goose -dir migrations postgres "$DATABASE_URL" down

# Buat migration baru
goose -dir migrations create add_product_weight sql
```

---

## Domain Events Published

| Event | Trigger | NATS Subject |
|-------|---------|--------------|
| `ProductCreated` | POST /products | `catalog.product.created` |
| `ProductUpdated` | PUT /products/:id | `catalog.product.updated` |
| `ProductPublished` | POST /products/:id/publish | `catalog.product.published` |
| `ProductArchived` | POST /products/:id/archive | `catalog.product.archived` |

---

## Testing

```bash
cd services/product-service

# Unit test
go test ./internal/domain/...
go test ./internal/application/...

# Integration test (butuh postgres running)
go test ./internal/infrastructure/...
go test ./tests/...

# Semua test
go test ./...

# Dengan verbose
go test ./... -v

# Dengan coverage
go test ./... -cover
```

---

## Build & Docker

```bash
# Build binary
go build -o main ./cmd/server

# Build Docker image (dari root project)
docker build -f services/product-service/Dockerfile -t cc-product-service .

# Jalankan via Docker
docker run -p 8083:8083 \
  -e PORT=8083 \
  -e DATABASE_URL="..." \
  cc-product-service
```

---

## Product Status Flow

```
draft ──► published ──► archived
  │                        │
  └────────────────────────┘
       (bisa archive dari draft)
```

Lihat detail: [state-machine.md](../../docs/03-architecture/state-machine.md)

---

## Related Documents

- [Backend Guidelines](../../docs/04-engineering/04-backend-guidelines.md)
- [API Guidelines](../../docs/03-architecture/api-guidelines.md)
- [Database Design](../../docs/03-architecture/database-design.md)
- [State Machine](../../docs/03-architecture/state-machine.md)
- [Error Handling Catalog](../../docs/04-engineering/error-handling-catalog.md)
- [Environment Variables](../../docs/04-engineering/env-variables-reference.md)
