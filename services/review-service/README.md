# Review Service

**Service:** review-service

**Port:** `8088`

**Database:** `review_db` (PostgreSQL)

**Owner:** Backend Team

---

## Responsibility

Manajemen review dan rating produk oleh buyer.

- Buyer dapat menulis review setelah order completed
- Rating produk (1-5 bintang)
- Tampilkan review di storefront
- Rata-rata rating produk

---

## API Endpoints

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `GET` | `/api/v1/reviews/products/:productId` | ❌ | List review untuk produk |
| `POST` | `/api/v1/reviews` | ✅ | Tulis review (buyer) |
| `PUT` | `/api/v1/reviews/:id` | ✅ | Edit review milik sendiri |
| `DELETE` | `/api/v1/reviews/:id` | ✅ | Hapus review milik sendiri |
| `GET` | `/healthz` | ❌ | Health check |

---

## Environment Variables

| Variable | Wajib | Default (Dev) |
|----------|-------|---------------|
| `PORT` | ✅ | `8088` |
| `DATABASE_URL` | ✅ | `postgres://postgres:postgres@localhost:5432/review_db?sslmode=disable` |
| `JWT_SECRET` | ✅ | `dev-secret-change-me` |
| `APP_ENV` | ✅ | `development` |

---

## Running Locally

```bash
cd services/review-service
export PORT=8088
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/review_db?sslmode=disable"
export JWT_SECRET="dev-secret-change-me"
export APP_ENV="development"

goose -dir migrations postgres "$DATABASE_URL" up
go run ./cmd/server
```

---

## Business Rules

- Buyer hanya bisa review produk yang sudah dibeli (order status: `completed`)
- Satu buyer hanya bisa satu review per produk per order
- Review bisa diedit dalam 7 hari setelah dibuat
- Rating: integer 1-5

---

## Related Documents

- [State Machine](../../docs/03-architecture/state-machine.md)
- [Backend Guidelines](../../docs/04-engineering/04-backend-guidelines.md)
- [Environment Variables](../../docs/04-engineering/env-variables-reference.md)
