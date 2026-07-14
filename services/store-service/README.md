# Store Service (Tenant Service)

**Service:** store-service

**Port:** `8082`

**Database:** `tenant_db` (PostgreSQL)

**Owner:** Backend Team

---

## Responsibility

Mengelola data merchant/tenant di platform CloudCommerce.

- Buat dan kelola profil toko (tenant)
- Konfigurasi toko: nama, domain, branding
- Manajemen domain toko (slug-based)
- Informasi toko untuk storefront buyer

---

## API Endpoints

Base URL: `http://localhost:8082`

### Protected (Butuh JWT)

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/v1/stores/me` | Info toko milik seller |
| `PUT` | `/api/v1/stores/me` | Update info toko |
| `GET` | `/api/v1/stores/me/settings` | Pengaturan toko |
| `PUT` | `/api/v1/stores/me/settings` | Update pengaturan |

### Public

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/v1/stores/:slug` | Info toko by slug (untuk buyer storefront) |

### Health

| Method | Path |
|--------|------|
| `GET` | `/healthz` |

---

## Environment Variables

| Variable | Wajib | Default (Dev) |
|----------|-------|---------------|
| `PORT` | ✅ | `8082` |
| `DATABASE_URL` | ✅ | `postgres://postgres:postgres@localhost:5432/tenant_db?sslmode=disable` |
| `JWT_SECRET` | ✅ | `dev-secret-change-me` |
| `APP_ENV` | ✅ | `development` |

---

## Running Locally

```bash
cd services/store-service

export PORT=8082
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/tenant_db?sslmode=disable"
export JWT_SECRET="dev-secret-change-me"
export APP_ENV="development"

goose -dir migrations postgres "$DATABASE_URL" up
go run ./cmd/server
```

---

## Related Documents

- [Backend Guidelines](../../docs/04-engineering/04-backend-guidelines.md)
- [Service Boundaries](../../docs/03-architecture/service-boundaries.md)
- [Environment Variables](../../docs/04-engineering/env-variables-reference.md)
