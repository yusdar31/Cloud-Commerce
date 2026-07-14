# Environment Variables Reference

**Project:** CloudCommerce

**Version:** 1.0.0

**Status:** Approved

**Last Updated:** July 2026

---

## 1. Purpose

Dokumen ini mendefinisikan **seluruh environment variable** yang digunakan oleh setiap service di CloudCommerce.

Gunakan sebagai referensi saat:
- Setup environment lokal
- Membuat file `.env`
- Mengkonfigurasi Kubernetes Secrets
- Membuat CI/CD pipeline

---

## 2. Infrastruktur (Docker Compose)

Credential default untuk development lokal:

| Variable | Nilai Default (Dev) | Deskripsi |
|----------|--------------------|-----------| 
| `POSTGRES_DB` | `cloudcommerce` | Database utama (init) |
| `POSTGRES_USER` | `postgres` | PostgreSQL user |
| `POSTGRES_PASSWORD` | `postgres` | PostgreSQL password |

**⚠️ JANGAN gunakan nilai ini di production.**

---

## 3. Shared Variables (Semua Service)

Semua service Go menggunakan variable berikut:

| Variable | Wajib | Contoh | Deskripsi |
|----------|-------|--------|-----------|
| `PORT` | ✅ | `8083` | Port HTTP service |
| `APP_ENV` | ✅ | `development` | Environment (`development` / `production`) |
| `JWT_SECRET` | ✅ | `dev-secret-change-me` | Secret untuk JWT signing (min 32 char di production) |
| `LOG_LEVEL` | ❌ | `info` | Log level (`debug`, `info`, `warn`, `error`) |

---

## 4. API Gateway

**Port:** `8080`

| Variable | Wajib | Contoh | Deskripsi |
|----------|-------|--------|-----------|
| `PORT` | ✅ | `8080` | Port HTTP |
| `JWT_SECRET` | ✅ | `dev-secret-change-me` | JWT secret (sama dengan semua service) |
| `NATS_URL` | ✅ | `nats://localhost:4222` | NATS connection string |
| `REDIS_URL` | ✅ | `redis://localhost:6379` | Redis connection string |
| `USER_SERVICE_URL` | ✅ | `http://localhost:8081` | URL user-service |
| `PRODUCT_SERVICE_URL` | ✅ | `http://localhost:8083` | URL product-service |
| `STORE_SERVICE_URL` | ❌ | `http://localhost:8082` | URL store-service |
| `ORDER_SERVICE_URL` | ❌ | `http://localhost:8085` | URL order-service |
| `PAYMENT_SERVICE_URL` | ❌ | `http://localhost:8086` | URL payment-service |
| `APP_ENV` | ✅ | `development` | Environment |
| `RATE_LIMIT_RPS` | ❌ | `100` | Rate limit per tenant (request per menit) |

---

## 5. User Service (Identity)

**Port:** `8081` | **Database:** `identity_db`

| Variable | Wajib | Contoh | Deskripsi |
|----------|-------|--------|-----------|
| `PORT` | ✅ | `8081` | Port HTTP |
| `DATABASE_URL` | ✅ | `postgres://postgres:postgres@localhost:5432/identity_db?sslmode=disable` | PostgreSQL connection string |
| `REDIS_URL` | ✅ | `redis://localhost:6379` | Redis untuk session/token blacklist |
| `JWT_SECRET` | ✅ | `dev-secret-change-me` | JWT signing secret |
| `JWT_EXPIRY` | ✅ | `24h` | Access token expiry |
| `REFRESH_TOKEN_EXPIRY` | ❌ | `720h` | Refresh token expiry (30 hari) |
| `APP_ENV` | ✅ | `development` | Environment |

---

## 6. Store Service (Tenant)

**Port:** `8082` | **Database:** `tenant_db`

| Variable | Wajib | Contoh | Deskripsi |
|----------|-------|--------|-----------|
| `PORT` | ✅ | `8082` | Port HTTP |
| `DATABASE_URL` | ✅ | `postgres://postgres:postgres@localhost:5432/tenant_db?sslmode=disable` | PostgreSQL connection string |
| `JWT_SECRET` | ✅ | `dev-secret-change-me` | JWT secret |
| `APP_ENV` | ✅ | `development` | Environment |

---

## 7. Product Service (Catalog)

**Port:** `8083` | **Database:** `catalog_db`

| Variable | Wajib | Contoh | Deskripsi |
|----------|-------|--------|-----------|
| `PORT` | ✅ | `8083` | Port HTTP |
| `DATABASE_URL` | ✅ | `postgres://postgres:postgres@localhost:5432/catalog_db?sslmode=disable` | PostgreSQL connection string |
| `REDIS_URL` | ✅ | `redis://localhost:6379` | Redis untuk product cache |
| `JWT_SECRET` | ✅ | `dev-secret-change-me` | JWT secret |
| `JWT_EXPIRY` | ❌ | `24h` | Token expiry |
| `APP_ENV` | ✅ | `development` | Environment |
| `CACHE_TTL_SECONDS` | ❌ | `60` | TTL cache produk di Redis (detik) |

---

## 8. Inventory Service

**Port:** `8084` | **Database:** `inventory_db`

| Variable | Wajib | Contoh | Deskripsi |
|----------|-------|--------|-----------|
| `PORT` | ✅ | `8084` | Port HTTP |
| `DATABASE_URL` | ✅ | `postgres://postgres:postgres@localhost:5432/inventory_db?sslmode=disable` | PostgreSQL connection string |
| `JWT_SECRET` | ✅ | `dev-secret-change-me` | JWT secret |
| `NATS_URL` | ❌ | `nats://localhost:4222` | NATS (untuk subscribe event ProductCreated) |
| `RESERVATION_TIMEOUT_MINUTES` | ❌ | `30` | Timeout reservasi stok (menit) |
| `APP_ENV` | ✅ | `development` | Environment |

---

## 9. Order Service

**Port:** `8085` | **Database:** `order_db`

| Variable | Wajib | Contoh | Deskripsi |
|----------|-------|--------|-----------|
| `PORT` | ✅ | `8085` | Port HTTP |
| `DATABASE_URL` | ✅ | `postgres://postgres:postgres@localhost:5432/order_db?sslmode=disable` | PostgreSQL connection string |
| `NATS_URL` | ✅ | `nats://localhost:4222` | NATS untuk publish OrderCreated event |
| `JWT_SECRET` | ✅ | `dev-secret-change-me` | JWT secret |
| `INVENTORY_SERVICE_URL` | ❌ | `http://localhost:8084` | URL inventory-service (untuk reserve stock) |
| `APP_ENV` | ✅ | `development` | Environment |

---

## 10. Payment Service

**Port:** `8086` | **Database:** `payment_db`

| Variable | Wajib | Contoh | Deskripsi |
|----------|-------|--------|-----------|
| `PORT` | ✅ | `8086` | Port HTTP |
| `DATABASE_URL` | ✅ | `postgres://postgres:postgres@localhost:5432/payment_db?sslmode=disable` | PostgreSQL connection string |
| `NATS_URL` | ✅ | `nats://localhost:4222` | NATS untuk subscribe OrderCreated, publish PaymentSucceeded |
| `JWT_SECRET` | ✅ | `dev-secret-change-me` | JWT secret |
| `MIDTRANS_SERVER_KEY` | ⚠️ | `SB-Mid-server-xxx` | Midtrans server key (Sandbox) |
| `MIDTRANS_CLIENT_KEY` | ⚠️ | `SB-Mid-client-xxx` | Midtrans client key (Sandbox) |
| `MIDTRANS_IS_PRODUCTION` | ❌ | `false` | Mode production Midtrans |
| `XENDIT_API_KEY` | ⚠️ | `xnd_development_xxx` | Xendit API key (alternatif Midtrans) |
| `WEBHOOK_SECRET` | ✅ | `webhook-secret-change-me` | Secret untuk validasi webhook signature |
| `APP_ENV` | ✅ | `development` | Environment |

> **⚠️ MIDTRANS_SERVER_KEY dan XENDIT_API_KEY tidak boleh ada di repository.**
> Gunakan Kubernetes Secrets atau environment variable langsung.

---

## 11. Notification Service

**Port:** `8087` | **Database:** `notification_db`

| Variable | Wajib | Contoh | Deskripsi |
|----------|-------|--------|-----------|
| `PORT` | ✅ | `8087` | Port HTTP |
| `DATABASE_URL` | ✅ | `postgres://postgres:postgres@localhost:5432/notification_db?sslmode=disable` | PostgreSQL connection string |
| `NATS_URL` | ✅ | `nats://localhost:4222` | NATS untuk subscribe events |
| `JWT_SECRET` | ✅ | `dev-secret-change-me` | JWT secret |
| `SMTP_HOST` | ✅ | `localhost` | SMTP server host (gunakan Mailhog untuk dev) |
| `SMTP_PORT` | ✅ | `1025` | SMTP port (Mailhog: 1025) |
| `SMTP_USERNAME` | ❌ | `` | SMTP username (kosong untuk Mailhog) |
| `SMTP_PASSWORD` | ❌ | `` | SMTP password (kosong untuk Mailhog) |
| `EMAIL_FROM` | ✅ | `noreply@cloudcommerce.com` | Alamat pengirim email |
| `APP_ENV` | ✅ | `development` | Environment |

---

## 12. Review Service

**Port:** `8088` | **Database:** `review_db`

| Variable | Wajib | Contoh | Deskripsi |
|----------|-------|--------|-----------|
| `PORT` | ✅ | `8088` | Port HTTP |
| `DATABASE_URL` | ✅ | `postgres://postgres:postgres@localhost:5432/review_db?sslmode=disable` | PostgreSQL connection string |
| `JWT_SECRET` | ✅ | `dev-secret-change-me` | JWT secret |
| `APP_ENV` | ✅ | `development` | Environment |

---

## 13. Frontend — Storefront

**Port:** `3000`

| Variable | Wajib | Contoh | Deskripsi |
|----------|-------|--------|-----------|
| `NEXT_PUBLIC_API_URL` | ✅ | `http://localhost:8080` | Base URL API Gateway |
| `NEXT_PUBLIC_APP_NAME` | ❌ | `CloudCommerce` | Nama aplikasi |
| `NEXT_PUBLIC_APP_URL` | ❌ | `http://localhost:3000` | URL aplikasi sendiri |

> **Prefix `NEXT_PUBLIC_`** membuat variable dapat diakses di client-side (browser).
> Variable tanpa prefix hanya tersedia di server-side Next.js.

---

## 14. Cara Setup `.env` Lokal

### Template `.env` untuk Product Service

```env
# Buat file: services/product-service/.env (JANGAN commit file ini)

PORT=8083
APP_ENV=development
DATABASE_URL=postgres://postgres:postgres@localhost:5432/catalog_db?sslmode=disable
REDIS_URL=redis://localhost:6379
JWT_SECRET=dev-secret-change-me
JWT_EXPIRY=24h
CACHE_TTL_SECONDS=60
LOG_LEVEL=debug
```

### Template `.env.local` untuk Storefront

```env
# Buat file: apps/storefront/.env.local (JANGAN commit file ini)

NEXT_PUBLIC_API_URL=http://localhost:8080
NEXT_PUBLIC_APP_NAME=CloudCommerce
NEXT_PUBLIC_APP_URL=http://localhost:3000
```

---

## 15. Rules

```
✅ Gunakan environment variables untuk semua config
✅ Buat .env.example di setiap service (tanpa nilai rahasia)
❌ JANGAN commit file .env dengan nilai nyata
❌ JANGAN hardcode URL, secret, atau API key di kode
❌ JANGAN push Midtrans/Xendit key ke repository
```

---

## 16. Kubernetes Secrets (Production)

Di production, semua variable sensitif disimpan sebagai Kubernetes Secret:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: product-service-secrets
type: Opaque
stringData:
  DATABASE_URL: "postgres://..."
  JWT_SECRET: "..."
  REDIS_URL: "redis://..."
```

---

## 17. Related Documents

- [Local Dev Setup](local-dev-setup.md)
- [Technology Stack](01-tech-stack.md)
- [Deployment Strategy](04-deployment-strategy.md)
- [Coding Standards](03-coding-standards.md)
