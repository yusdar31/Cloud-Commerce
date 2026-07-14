# Payment Service

**Service:** payment-service

**Port:** `8086`

**Database:** `payment_db` (PostgreSQL)

**Owner:** Backend Team

---

## Responsibility

Integrasi dengan payment gateway (Midtrans / Xendit) dan penanganan webhook.

- Subscribe event `OrderCreated`, generate invoice
- Panggil Midtrans/Xendit API untuk mendapatkan payment URL
- Terima webhook dari payment gateway
- Validasi signature webhook
- Publish event `PaymentSucceeded` / `PaymentFailed`
- Manajemen refund

---

## API Endpoints

### Payments

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `GET` | `/api/v1/payments/:orderId` | ✅ | Status pembayaran untuk order |
| `POST` | `/api/v1/payments/webhook/midtrans` | ❌ | Webhook dari Midtrans (validasi signature) |
| `POST` | `/api/v1/payments/webhook/xendit` | ❌ | Webhook dari Xendit |

### Refunds

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `POST` | `/api/v1/refunds` | ✅ | Request refund |
| `GET` | `/api/v1/refunds/:id` | ✅ | Status refund |

---

## Environment Variables

| Variable | Wajib | Default (Dev) |
|----------|-------|---------------|
| `PORT` | ✅ | `8086` |
| `DATABASE_URL` | ✅ | `postgres://postgres:postgres@localhost:5432/payment_db?sslmode=disable` |
| `NATS_URL` | ✅ | `nats://localhost:4222` |
| `JWT_SECRET` | ✅ | `dev-secret-change-me` |
| `MIDTRANS_SERVER_KEY` | ⚠️ | - | Dapatkan dari dashboard Midtrans Sandbox |
| `MIDTRANS_CLIENT_KEY` | ⚠️ | - | Dapatkan dari dashboard Midtrans Sandbox |
| `MIDTRANS_IS_PRODUCTION` | ❌ | `false` | |
| `WEBHOOK_SECRET` | ✅ | `webhook-secret-change-me` | |
| `APP_ENV` | ✅ | `development` | |

---

## Running Locally

```bash
cd services/payment-service
export PORT=8086
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/payment_db?sslmode=disable"
export NATS_URL="nats://localhost:4222"
export JWT_SECRET="dev-secret-change-me"
export MIDTRANS_SERVER_KEY="SB-Mid-server-xxxxx"
export MIDTRANS_IS_PRODUCTION=false
export WEBHOOK_SECRET="webhook-secret-change-me"
export APP_ENV="development"

goose -dir migrations postgres "$DATABASE_URL" up
go run ./cmd/server
```

---

## NATS Events

**Subscribed:**

| Event | Subject | Action |
|-------|---------|--------|
| `OrderCreated` | `order.created` | Buat invoice, panggil payment gateway |

**Published:**

| Event | Subject | Trigger |
|-------|---------|---------|
| `PaymentSucceeded` | `payment.succeeded` | Webhook success diterima & validated |
| `PaymentFailed` | `payment.failed` | Webhook failed / expired |
| `RefundCompleted` | `payment.refunded` | Refund berhasil diproses |

---

## Idempotency

Webhook dari payment gateway **harus idempotent**.

```
Jika webhook yang sama diterima 2x:
→ Cek apakah transaction_id sudah ada di database
→ Jika sudah ada: return 200 OK tanpa memproses ulang
→ Jika belum ada: proses dan simpan
```

---

## Related Documents

- [State Machine](../../docs/03-architecture/state-machine.md)
- [Event Storming](../../docs/03-architecture/event-storming.md)
- [Error Handling Catalog](../../docs/04-engineering/error-handling-catalog.md)
- [Environment Variables](../../docs/04-engineering/env-variables-reference.md)
