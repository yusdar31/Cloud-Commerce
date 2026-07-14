# CloudCommerce

Multi-tenant SaaS e-commerce platform dibangun di atas microservices architecture.
Dirancang untuk berjalan di Kubernetes (local → on-prem → cloud).

---

## Tech Stack

| Layer | Technology |
|-------|------------|
| Frontend | Next.js 15 (App Router) + Tailwind CSS v4 |
| Backend | Go 1.25 + Gin |
| Database | PostgreSQL 17 |
| Cache | Redis 7 |
| Event Bus | NATS JetStream |
| Object Storage | MinIO |
| API Gateway | Go + Gin |
| Orchestration | Kubernetes (Kind → k3s → GKE) |

---

## Quick Start (Local Dev)

### Prerequisites

```bash
# Required
- Go 1.25+
- Node.js 22+ + pnpm
- Docker Desktop
```

### 1. Jalankan Infrastruktur

```bash
docker compose up postgres redis nats -d
```

### 2. Jalankan Backend Service

```bash
cd services/product-service
go run ./cmd/server
```

### 3. Jalankan Frontend

```bash
cd apps/storefront
pnpm dev
```

Storefront tersedia di: **http://localhost:3000**

---

## Port Map

| Service | Port |
|---------|------|
| Storefront (Next.js) | 3000 |
| API Gateway | 8080 |
| User Service | 8081 |
| Store Service | 8082 |
| Product Service | 8083 |
| Inventory Service | 8084 |
| Order Service | 8085 |
| Payment Service | 8086 |
| Notification Service | 8087 |
| Review Service | 8088 |
| PostgreSQL | 5432 |
| Redis | 6379 |
| NATS | 4222 |
| NATS Monitor | 8222 |

---

## Struktur Repository

```
cloudcommerce/
├── apps/
│   ├── storefront/         # Next.js — Buyer storefront
│   └── web/                # Next.js — Marketing & landing page
├── services/
│   ├── api-gateway/        # Go — Entry point semua request
│   ├── user-service/       # Go — Auth & identity
│   ├── store-service/      # Go — Tenant management
│   ├── product-service/    # Go — Catalog produk
│   ├── inventory-service/  # Go — Stok
│   ├── order-service/      # Go — Checkout & order
│   ├── payment-service/    # Go — Payment gateway integration
│   ├── notification-service/ # Go — Email notification
│   └── review-service/     # Go — Review & rating
├── packages/
│   ├── shared-go/          # Shared Go utilities (JWT, middleware, dll)
│   ├── shared-types/       # Shared TypeScript types
│   └── shared-config/      # Shared config
├── infra/                  # Docker, Helm, K8s, Terraform
├── docs/                   # Dokumentasi lengkap
└── scripts/                # Automation scripts
```

---

## Roadmap Phase

| Phase | Status | Scope |
|-------|--------|-------|
| Phase 0 | ✅ Done | Documentation |
| Phase 1 | 🚧 In Progress | Monorepo, Storefront, User Service, Product Service |
| Phase 2 | ⏳ Planned | Docker Compose full stack, API Gateway, Redis |
| Phase 3 | ⏳ Planned | Inventory, Order, Payment |
| Phase 4 | ⏳ Planned | NATS event-driven, Saga pattern |
| Phase 5 | ⏳ Planned | Helm, ArgoCD, GitOps |
| Phase 6 | ⏳ Planned | Monitoring (Prometheus, Grafana, Loki, Tempo) |
| Phase 7 | ⏳ Planned | Production (Terraform, AWS/GCP) |

---

## Dokumentasi

| Kategori | Link |
|----------|------|
| Product | [docs/01-product/](docs/01-product/) |
| Design | [docs/02-design/](docs/02-design/) |
| Architecture | [docs/03-architecture/](docs/03-architecture/) |
| Engineering | [docs/04-engineering/](docs/04-engineering/) |
| Roadmap | [docs/roadmaps.md](docs/roadmaps.md) |

---

## AI Coding Guide

Jika menggunakan AI assistant, baca [CLAUDE.md](CLAUDE.md) terlebih dahulu.