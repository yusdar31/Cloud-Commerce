# Order Service

**Service:** order-service

**Port:** `8085`

**Database:** `order_db` (PostgreSQL)

**Owner:** Backend Team

---

## Responsibility

Orkestrasi proses checkout dan manajemen state pesanan.

- Manajemen keranjang belanja (cart)
- Proses checkout: buat order, reserve stok
- Manajemen status order
- Publish event `OrderCreated`, `OrderConfirmed`, `OrderCancelled`
- Subscribe event `PaymentSucceeded` untuk update status order

---

## API Endpoints

### Cart

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/v1/cart` | Lihat isi cart |
| `POST` | `/api/v1/cart/items` | Tambah item ke cart |
| `PUT` | `/api/v1/cart/items/:id` | Update quantity |
| `DELETE` | `/api/v1/cart/items/:id` | Hapus item dari cart |
| `DELETE` | `/api/v1/cart` | Kosongkan cart |

### Orders

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/api/v1/orders` | Checkout — buat order dari cart |
| `GET` | `/api/v1/orders` | List order (buyer: milik sendiri; seller: per toko) |
| `GET` | `/api/v1/orders/:id` | Detail order |
| `POST` | `/api/v1/orders/:id/confirm` | Seller konfirmasi order (→ processing) |
| `POST` | `/api/v1/orders/:id/ship` | Seller tandai dikirim (→ shipped) |
| `POST` | `/api/v1/orders/:id/complete` | Buyer konfirmasi diterima (→ completed) |
| `POST` | `/api/v1/orders/:id/cancel` | Batalkan order |

---

## Environment Variables

| Variable | Wajib | Default (Dev) |
|----------|-------|---------------|
| `PORT` | ✅ | `8085` |
| `DATABASE_URL` | ✅ | `postgres://postgres:postgres@localhost:5432/order_db?sslmode=disable` |
| `NATS_URL` | ✅ | `nats://localhost:4222` |
| `JWT_SECRET` | ✅ | `dev-secret-change-me` |
| `INVENTORY_SERVICE_URL` | ❌ | `http://localhost:8084` |
| `APP_ENV` | ✅ | `development` |

---

## Running Locally

```bash
cd services/order-service
export PORT=8085
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/order_db?sslmode=disable"
export NATS_URL="nats://localhost:4222"
export JWT_SECRET="dev-secret-change-me"
export APP_ENV="development"

goose -dir migrations postgres "$DATABASE_URL" up
go run ./cmd/server
```

---

## NATS Events

**Published:**

| Event | Subject | Trigger |
|-------|---------|---------|
| `OrderCreated` | `order.created` | POST /orders (checkout) |
| `OrderConfirmed` | `order.confirmed` | Seller konfirmasi |
| `OrderCancelled` | `order.cancelled` | Cancel order |

**Subscribed:**

| Event | Subject | Action |
|-------|---------|--------|
| `PaymentSucceeded` | `payment.succeeded` | Update order → `paid` |
| `PaymentFailed` | `payment.failed` | Update order → `cancelled` |

---

## Order Status

Lihat: [state-machine.md](../../docs/03-architecture/state-machine.md)

```
pending → awaiting_payment → paid → processing → shipped → completed
                          ↘                                        
                          cancelled ← payment failed / timeout
```

---

## Related Documents

- [State Machine](../../docs/03-architecture/state-machine.md)
- [Event Storming](../../docs/03-architecture/event-storming.md)
- [Backend Guidelines](../../docs/04-engineering/04-backend-guidelines.md)
- [Environment Variables](../../docs/04-engineering/env-variables-reference.md)
