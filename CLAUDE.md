# CloudCommerce — AI Coding Guide

Dokumen ini adalah panduan khusus untuk AI coding assistant (Claude, Gemini, Copilot, dll).
Baca dokumen ini PERTAMA KALI sebelum menyentuh kode apapun.

---

## 🎯 Apa Itu Project Ini

CloudCommerce adalah **Multi-Tenant SaaS E-commerce Platform**.

- Merchant (Seller) bisa membuka toko sendiri di platform ini
- Setiap merchant punya data yang terisolasi (multi-tenant)
- Buyer bisa berbelanja di storefront milik merchant
- Backend microservices di Go, Frontend di Next.js 15

---

## 📍 Status Phase Saat Ini

**Phase 1 — Aktif dikerjakan:**

- [x] Monorepo setup (pnpm + Turborepo)
- [x] Docker Compose infrastruktur (postgres, redis, nats)
- [x] `product-service` — CRUD produk (sebagian)
- [x] `user-service` — auth & JWT (sebagian)
- [x] `store-service` — tenant management (sebagian)
- [x] `apps/storefront` — Next.js 15 landing page + auth UI
- [ ] `api-gateway` — routing & JWT validation (belum selesai)
- [ ] Integrasi storefront ↔ product-service via API

---

## 🗂️ Dokumen Wajib Dibaca Sebelum Coding

| Dokumen | Isi |
|---------|-----|
| [tech-stack.md](docs/04-engineering/01-tech-stack.md) | Stack teknologi & versi |
| [monorepo-structure.md](docs/04-engineering/02-monorepo-structure.md) | Struktur folder & aturan |
| [coding-standards.md](docs/04-engineering/03-coding-standards.md) | Naming, error handling, logging |
| [backend-guidelines.md](docs/04-engineering/04-backend-guidelines.md) | Pola Go: Clean Arch, handler, migration |
| [frontend-guidelines.md](docs/04-engineering/05-frontend-guidelines.md) | Next.js App Router, state management |
| [api-guidelines.md](docs/03-architecture/api-guidelines.md) | Response format, error format, versioning |
| [service-boundaries.md](docs/03-architecture/service-boundaries.md) | Tanggung jawab setiap service |
| [event-storming.md](docs/03-architecture/event-storming.md) | Event bisnis & flow antar service |
| [state-machine.md](docs/03-architecture/state-machine.md) | Status Order, Payment, Inventory |
| [local-dev-setup.md](docs/04-engineering/local-dev-setup.md) | Setup environment lokal |

---

## ⚡ Commands Penting

### Frontend (Storefront)

```bash
cd apps/storefront
pnpm dev          # Dev server → http://localhost:3000
pnpm build        # Production build
pnpm lint         # ESLint
pnpm type-check   # TypeScript check
```

### Backend (Go Service)

```bash
cd services/product-service
go run ./cmd/server        # Jalankan service
go test ./...              # Semua test
go test ./internal/...     # Unit + integration test

# Migration (gunakan goose)
goose -dir migrations postgres "$DATABASE_URL" up
goose -dir migrations postgres "$DATABASE_URL" status
```

### Infrastruktur (Docker Compose)

```bash
# Di root project
docker compose up postgres redis nats -d         # Infra saja
docker compose up -d                              # Semua service
docker compose logs -f product-service            # Log service tertentu
docker compose down                               # Stop semua
docker compose down -v                            # Stop + hapus volumes
```

### Monorepo (Turborepo)

```bash
# Di root project
pnpm dev          # Jalankan semua apps sekaligus
pnpm build        # Build semua
pnpm lint         # Lint semua
```

---

## 🏗️ Struktur Service Go (WAJIB DIIKUTI)

Setiap microservice Go mengikuti Clean Architecture:

```
services/product-service/
├── cmd/server/
│   └── main.go              # Entry point SAJA — wire dependencies
├── internal/
│   ├── domain/
│   │   ├── entity.go        # Struct entitas + domain errors
│   │   └── repository.go    # Interface repository (BUKAN implementasi)
│   ├── application/
│   │   └── service.go       # Business logic / use cases
│   ├── infrastructure/
│   │   └── postgres/
│   │       └── repository.go # Implementasi query database
│   └── transport/
│       ├── handler.go       # Gin handler (HTTP concern saja)
│       ├── routes.go        # Setup routing
│       └── dto.go           # Request/Response structs + validasi
├── migrations/              # File SQL migration (goose)
├── Dockerfile
└── go.mod
```

**Aturan lapisan:**
- `domain` → tidak boleh import apapun kecuali stdlib Go
- `application` → hanya boleh import `domain`
- `infrastructure` → implements interface dari `domain`
- `transport` → hanya boleh import `application` dan `domain`

---

## 🎨 Struktur Frontend Next.js (WAJIB DIIKUTI)

```
apps/storefront/src/
├── app/                     # Next.js App Router
│   ├── (public)/            # Landing, pricing — NO AUTH required
│   ├── (auth)/              # Login, register pages
│   ├── (dashboard)/         # Seller dashboard — REQUIRES auth
│   └── (store)/             # Buyer storefront — public + auth mixed
├── components/
│   ├── ui/                  # shadcn/ui base components
│   ├── layout/              # Navbar, Footer, Sidebar
│   └── shared/              # Reusable business components
├── features/                # Feature-based modules (auth, products, cart, dll)
├── hooks/                   # Custom React hooks
├── lib/                     # API client, utils
├── stores/                  # Zustand stores (cart, auth)
└── types/                   # TypeScript type definitions
```

---

## 🔑 Aturan Wajib (JANGAN DILANGGAR)

### Multi-Tenant

```
❌ DILARANG: Mengambil data tanpa filter tenant_id
✅ WAJIB: Setiap query database HARUS include tenant_id dari JWT claims
```

### API Response Format

```json
// Success
{ "data": { ... }, "meta": { ... } }

// Error (RFC 9457)
{ "type": "...", "title": "...", "status": 422, "detail": "...", "traceId": "..." }
```

### Error Handling Go

```go
// ✅ Benar
if err != nil {
    return nil, fmt.Errorf("find product: %w", err)
}

// ❌ Salah — jangan ignore error
product, _ := repo.Find(id)
```

### Money — Selalu Integer

```go
// ✅ Benar: dalam unit terkecil (sen / IDR)
Price int64 // 150000 = Rp 1.500,00

// ❌ Salah
Price float64
```

### Timestamp — Selalu UTC

```go
CreatedAt time.Time // selalu simpan & kirim UTC
```

---

## 🌐 Service URLs (Local Dev)

```
Storefront:         http://localhost:3000
API Gateway:        http://localhost:8080
User Service:       http://localhost:8081
Store Service:      http://localhost:8082
Product Service:    http://localhost:8083
Inventory Service:  http://localhost:8084
Order Service:      http://localhost:8085
Payment Service:    http://localhost:8086
Notification:       http://localhost:8087
Review Service:     http://localhost:8088

Health check:       GET http://localhost:{port}/healthz
```

---

## 📦 Shared Packages

### `packages/shared-go` (Go)

```go
import "github.com/cloudcommerce/shared-go/jwt"        // JWT manager
import "github.com/cloudcommerce/shared-go/middleware"  // CORS, RequestID, JWTAuth, TenantIsolation
```

### `packages/shared-types` (TypeScript)

```ts
// Type definitions yang di-share antar frontend apps
```

---

## 🚫 Yang TIDAK Boleh Dilakukan AI

1. **Jangan** cross-import antar service di `services/`
2. **Jangan** share database antar service
3. **Jangan** hardcode secret, URL, atau credentials
4. **Jangan** membuat file di luar struktur yang sudah ditentukan
5. **Jangan** mengubah interface di `domain/` tanpa mengupdate implementasi
6. **Jangan** commit `.env` files atau secrets
7. **Jangan** skip validasi `tenant_id` di handler manapun
8. **Jangan** membuat migration baru tanpa rollback (down migration)

---

## ✅ Checklist Sebelum PR

- [ ] Kode kompilasi tanpa error
- [ ] Tidak ada lint error (`golangci-lint run` / `pnpm lint`)
- [ ] Semua fungsi baru ada error handling
- [ ] Response format mengikuti API Guidelines
- [ ] `tenant_id` divalidasi di setiap endpoint yang butuh auth
- [ ] Migration (jika ada) sudah bisa `up` dan `down`
- [ ] Tidak ada hardcoded secret
