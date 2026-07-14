# API Gateway

**Service:** api-gateway

**Port:** `8080`

**Owner:** Backend Team

---

## Responsibility

Single entry point untuk semua request eksternal ke CloudCommerce.

- Routing request ke microservice yang tepat
- Validasi JWT (autentikasi)
- Rate limiting per tenant
- Request ID injection
- CORS handling
- API versioning enforcement

**Catatan:** API Gateway tidak memiliki business logic. Hanya routing & cross-cutting concerns.

---

## Upstream Services

| Service | Internal URL |
|---------|-------------|
| user-service | `http://user-service:8081` |
| store-service | `http://store-service:8082` |
| product-service | `http://product-service:8083` |
| inventory-service | `http://inventory-service:8084` |
| order-service | `http://order-service:8085` |
| payment-service | `http://payment-service:8086` |
| notification-service | `http://notification-service:8087` |
| review-service | `http://review-service:8088` |

---

## Routing Rules

| Path Prefix | Upstream Service | Auth Required |
|-------------|-----------------|---------------|
| `/api/v1/auth/*` | user-service | ❌ No |
| `/api/v1/users/*` | user-service | ✅ Yes |
| `/api/v1/stores/*` | store-service | ✅ Yes (seller only) |
| `/api/v1/products/*` | product-service | ✅ Yes |
| `/api/v1/categories/*` | product-service | ✅ Yes |
| `/api/v1/inventory/*` | inventory-service | ✅ Yes |
| `/api/v1/orders/*` | order-service | ✅ Yes |
| `/api/v1/payments/*` | payment-service | ✅ Yes |
| `/api/v1/payments/webhook` | payment-service | ❌ No (signature validation) |
| `/api/v1/reviews/*` | review-service | ✅ Yes |

---

## Environment Variables

| Variable | Wajib | Default (Dev) |
|----------|-------|---------------|
| `PORT` | ✅ | `8080` |
| `JWT_SECRET` | ✅ | `dev-secret-change-me` |
| `NATS_URL` | ✅ | `nats://localhost:4222` |
| `REDIS_URL` | ✅ | `redis://localhost:6379` |
| `USER_SERVICE_URL` | ✅ | `http://localhost:8081` |
| `PRODUCT_SERVICE_URL` | ✅ | `http://localhost:8083` |
| `STORE_SERVICE_URL` | ❌ | `http://localhost:8082` |
| `ORDER_SERVICE_URL` | ❌ | `http://localhost:8085` |
| `PAYMENT_SERVICE_URL` | ❌ | `http://localhost:8086` |
| `RATE_LIMIT_RPS` | ❌ | `100` |
| `APP_ENV` | ✅ | `development` |

---

## Running Locally

```bash
cd services/api-gateway

export PORT=8080
export JWT_SECRET="dev-secret-change-me"
export NATS_URL="nats://localhost:4222"
export REDIS_URL="redis://localhost:6379"
export USER_SERVICE_URL="http://localhost:8081"
export PRODUCT_SERVICE_URL="http://localhost:8083"
export APP_ENV="development"

go run ./cmd/server
```

---

## Related Documents

- [API Guidelines](../../docs/03-architecture/api-guidelines.md)
- [Service Boundaries](../../docs/03-architecture/service-boundaries.md)
- [Environment Variables](../../docs/04-engineering/env-variables-reference.md)
