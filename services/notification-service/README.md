# Notification Service

**Service:** notification-service

**Port:** `8087`

**Database:** `notification_db` (PostgreSQL)

**Owner:** Backend Team

---

## Responsibility

Event-driven notification dispatcher. Mengirim email ke buyer dan seller.

- Subscribe ke business events dari NATS
- Render email template
- Kirim email via SMTP (Mailhog untuk dev, SendGrid/SES untuk prod)
- Tracking status pengiriman notifikasi

---

## Events yang Di-subscribe

| Event | NATS Subject | Email Dikirim Ke |
|-------|-------------|-----------------|
| `UserRegistered` | `identity.user.registered` | Buyer/Seller: Welcome email |
| `OrderCreated` | `order.created` | Buyer: Order confirmation |
| `PaymentSucceeded` | `payment.succeeded` | Buyer: Payment receipt; Seller: New order alert |
| `PaymentFailed` | `payment.failed` | Buyer: Payment failed notification |
| `OrderShipped` | `order.shipped` | Buyer: Shipping notification + tracking |
| `InventoryLow` | `inventory.low` | Seller: Low stock alert |

---

## API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/v1/notifications` | List notifikasi user |
| `PUT` | `/api/v1/notifications/:id/read` | Tandai sudah dibaca |
| `GET` | `/healthz` | Health check |

---

## Environment Variables

| Variable | Wajib | Default (Dev) |
|----------|-------|---------------|
| `PORT` | ✅ | `8087` |
| `DATABASE_URL` | ✅ | `postgres://postgres:postgres@localhost:5432/notification_db?sslmode=disable` |
| `NATS_URL` | ✅ | `nats://localhost:4222` |
| `JWT_SECRET` | ✅ | `dev-secret-change-me` |
| `SMTP_HOST` | ✅ | `localhost` |
| `SMTP_PORT` | ✅ | `1025` (Mailhog) |
| `SMTP_USERNAME` | ❌ | `` |
| `SMTP_PASSWORD` | ❌ | `` |
| `EMAIL_FROM` | ✅ | `noreply@cloudcommerce.com` |
| `APP_ENV` | ✅ | `development` |

---

## Running Locally

```bash
# Jalankan Mailhog (SMTP mock)
docker run -d -p 1025:1025 -p 8025:8025 mailhog/mailhog

cd services/notification-service
export PORT=8087
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/notification_db?sslmode=disable"
export NATS_URL="nats://localhost:4222"
export JWT_SECRET="dev-secret-change-me"
export SMTP_HOST="localhost"
export SMTP_PORT="1025"
export EMAIL_FROM="noreply@cloudcommerce.com"
export APP_ENV="development"

go run ./cmd/server
```

Mailhog UI (lihat email masuk): **http://localhost:8025**

---

## Related Documents

- [Event Storming](../../docs/03-architecture/event-storming.md)
- [Environment Variables](../../docs/04-engineering/env-variables-reference.md)
