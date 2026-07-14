# Inventory Service

**Service:** inventory-service

**Port:** `8084`

**Database:** `inventory_db` (PostgreSQL)

**Owner:** Backend Team

---

## Responsibility

Manajemen stok produk dan reservasi stok saat checkout.

- Inisialisasi inventori saat produk dibuat
- Reservasi stok saat checkout (soft reserve)
- Commit stok setelah payment berhasil
- Release stok jika payment gagal / timeout
- Adjustment stok manual oleh seller
- Alert stok menipis (low stock notification)

---

## API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/v1/inventory/:productId` | Cek stok produk |
| `PUT` | `/api/v1/inventory/:productId` | Adjust stok (seller) |
| `POST` | `/api/v1/inventory/reserve` | Reserve stok (dipanggil order-service) |
| `POST` | `/api/v1/inventory/release` | Release reservasi |
| `POST` | `/api/v1/inventory/commit` | Commit reservasi (kurangi permanen) |
| `GET` | `/healthz` | Health check |

---

## Environment Variables

| Variable | Wajib | Default (Dev) |
|----------|-------|---------------|
| `PORT` | ✅ | `8084` |
| `DATABASE_URL` | ✅ | `postgres://postgres:postgres@localhost:5432/inventory_db?sslmode=disable` |
| `JWT_SECRET` | ✅ | `dev-secret-change-me` |
| `NATS_URL` | ❌ | `nats://localhost:4222` |
| `RESERVATION_TIMEOUT_MINUTES` | ❌ | `30` |
| `APP_ENV` | ✅ | `development` |

---

## Running Locally

```bash
cd services/inventory-service
export PORT=8084
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/inventory_db?sslmode=disable"
export JWT_SECRET="dev-secret-change-me"
export APP_ENV="development"

goose -dir migrations postgres "$DATABASE_URL" up
go run ./cmd/server
```

---

## Stock Reservation Logic

```
RESERVE:   available -= qty, reserved += qty
COMMIT:    reserved  -= qty  (stok dikurangi permanen)
RELEASE:   reserved  -= qty, available += qty
```

Timeout reservasi: **30 menit** (konfigurasikan via `RESERVATION_TIMEOUT_MINUTES`).

---

## Related Documents

- [State Machine](../../docs/03-architecture/state-machine.md)
- [Event Storming](../../docs/03-architecture/event-storming.md)
- [Environment Variables](../../docs/04-engineering/env-variables-reference.md)
