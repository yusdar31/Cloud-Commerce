# Service Boundaries

**Project:** CloudCommerce

**Architecture:** Domain-Driven Design (DDD) + Microservices

**Version:** 1.0.0

**Status:** Draft

**Owner:** Solution Architecture Team

**Last Updated:** July 2026

---

# 1. Purpose

Dokumen ini mendefinisikan batas tanggung jawab setiap microservice di CloudCommerce.

Tujuan utama:

- Menghindari distributed monolith
- Menentukan ownership service
- Mendefinisikan API publik
- Mendefinisikan event yang dipublish dan dikonsumsi
- Menjadi dasar deployment Kubernetes
- Menjadi dasar CI/CD per service

---

# 2. Architecture Principles

Setiap service harus:

- Memiliki satu business capability
- Memiliki satu database sendiri
- Tidak mengakses database service lain
- Dapat di-deploy secara independen
- Dapat di-scale secara independen
- Berkomunikasi melalui REST (sinkron) atau NATS (asinkron)

---

# 3. Service Landscape

```
                    API Gateway
                         │
 ┌───────────────────────┼────────────────────────┐
 │                       │                        │
 ▼                       ▼                        ▼
Identity            Tenant                 Catalog
                                              │
                                              ▼
                                        Inventory
                                              │
                                              ▼
                                           Order
                                              │
                                              ▼
                                           Payment

Notification ◀───────────────────────────────┐
Analytics   ◀────────────────────────────────┘
Subscription ────────────────────────────────┘
```

---

# 4. Identity Service

## Purpose

Mengelola autentikasi dan identitas pengguna.

---

## Responsibilities

- Register
- Login
- JWT
- Refresh Token
- Session
- Role
- Permission

---

## Owns

- User
- Role
- Permission
- Session

---

## Database

identity_db

---

## REST APIs

```
POST /auth/register

POST /auth/login

POST /auth/logout

POST /auth/refresh

GET /users/me
```

---

## Publish Events

- UserRegistered
- UserLoggedIn
- PasswordChanged

---

## Consume Events

Tidak ada.

---

# 5. Tenant Service

## Purpose

Mengelola toko (tenant).

---

## Owns

- Tenant
- Branding
- Store Domain
- Store Settings

---

## Database

tenant_db

---

## Publish Events

- TenantCreated
- TenantActivated
- TenantSuspended

---

## Consume Events

- UserRegistered

---

# 6. Catalog Service

## Purpose

Mengelola katalog produk.

---

## Owns

- Product
- Category
- ProductImage

---

## Database

catalog_db

---

## REST APIs

```
GET /products

POST /products

PUT /products/{id}

DELETE /products/{id}
```

---

## Publish Events

- ProductCreated
- ProductUpdated
- ProductPublished
- ProductArchived

---

## Consume Events

- TenantCreated

---

# 7. Inventory Service

## Purpose

Mengelola stok produk.

---

## Owns

- Inventory
- Reservation
- StockMovement

---

## Database

inventory_db

---

## Publish Events

- StockReserved
- StockReleased
- StockAdjusted
- InventoryLow

---

## Consume Events

- ProductCreated
- OrderCreated

---

# 8. Order Service

## Purpose

Mengelola siklus hidup pesanan.

---

## Owns

- Cart
- CartItem
- Order
- OrderItem

---

## Database

order_db

---

## Publish Events

- OrderCreated
- OrderConfirmed
- OrderCancelled
- OrderCompleted

---

## Consume Events

- StockReserved
- PaymentSucceeded
- PaymentFailed

---

# 9. Payment Service

## Purpose

Mengelola pembayaran.

---

## Owns

- Payment
- Invoice
- Refund

---

## Database

payment_db

---

## External Systems

- Midtrans
- Xendit

---

## Publish Events

- PaymentRequested
- PaymentSucceeded
- PaymentFailed
- RefundCompleted

---

## Consume Events

- OrderCreated

---

# 10. Notification Service

## Purpose

Mengirim email dan notifikasi.

---

## Owns

- Notification
- EmailTemplate

---

## Database

notification_db

---

## Publish Events

- NotificationQueued
- NotificationSent
- NotificationFailed

---

## Consume Events

- UserRegistered
- OrderCreated
- PaymentSucceeded
- SubscriptionStarted

---

# 11. Subscription Service

## Purpose

Mengelola paket langganan tenant.

---

## Owns

- Subscription
- Plan
- Billing
- Quota

---

## Database

subscription_db

---

## Publish Events

- SubscriptionStarted
- SubscriptionRenewed
- SubscriptionExpired

---

## Consume Events

- TenantCreated
- PaymentSucceeded

---

# 12. Analytics Service

## Purpose

Membangun dashboard dan laporan.

---

## Owns

- DashboardSnapshot
- SalesMetrics
- Reports

---

## Database

analytics_db

---

## Publish Events

- DashboardUpdated
- DailyReportGenerated

---

## Consume Events

- ProductPublished
- OrderCreated
- PaymentSucceeded
- InventoryLow

---

# 13. API Gateway

## Responsibilities

- Authentication
- Rate Limiting
- Request Routing
- CORS
- API Versioning
- Request Logging

Gateway tidak memiliki business logic.

---

# 14. Service Dependency Matrix

| Service | REST Calls | Event Consumers |
|----------|------------|-----------------|
| Identity | - | Tenant, Notification |
| Tenant | Identity | Subscription, Catalog |
| Catalog | Tenant | Inventory, Analytics |
| Inventory | Catalog | Order |
| Order | Inventory, Payment | Notification, Analytics |
| Payment | Order | Notification |
| Notification | - | - |
| Subscription | Payment | Notification |
| Analytics | - | - |

---

# 15. Database Ownership

| Service | Database |
|----------|----------|
| Identity | identity_db |
| Tenant | tenant_db |
| Catalog | catalog_db |
| Inventory | inventory_db |
| Order | order_db |
| Payment | payment_db |
| Subscription | subscription_db |
| Notification | notification_db |
| Analytics | analytics_db |

Tidak ada database yang dibagi antar service.

---

# 16. Communication Rules

Gunakan REST jika:

- Membutuhkan respons langsung
- Operasi query
- Validasi sinkron

Gunakan NATS jika:

- Event bisnis
- Notifikasi
- Analytics
- Sinkronisasi antar service

---

# 17. Scaling Strategy

| Service | Scaling Priority |
|----------|------------------|
| API Gateway | High |
| Catalog | High |
| Inventory | High |
| Order | Critical |
| Payment | Critical |
| Notification | Medium |
| Analytics | Medium |
| Identity | Medium |
| Subscription | Low |

---

# 18. Deployment Strategy

Setiap service memiliki:

- Repository (atau folder monorepo)
- Dockerfile
- Helm Chart
- Kubernetes Deployment
- Horizontal Pod Autoscaler
- GitHub Actions Workflow

Semua service dapat dirilis secara independen.

---

# 19. Cross-Cutting Concerns

Semua service wajib memiliki:

- Structured Logging
- OpenTelemetry Tracing
- Prometheus Metrics
- Health Check
- Readiness Probe
- Liveness Probe
- Configuration via Environment Variables
- Secret Management

---

# 20. Anti-Patterns

CloudCommerce menghindari:

❌ Shared Database

❌ Shared Business Logic

❌ Circular Dependency

❌ Distributed Transaction

❌ Synchronous Chain > 3 Services

❌ Service saling membaca tabel

❌ API Gateway berisi business logic

---

# 21. Future Services

Roadmap berikutnya:

- Search Service
- Recommendation Service
- Review Service
- Shipping Service
- Promotion Service
- Fraud Detection Service
- AI Assistant Service
- Audit Service

Semua akan mengikuti prinsip yang sama.

---

# 22. Related Documents

- System Context
- Domain Model
- Event Storming
- API Guidelines
- Database Design
- Deployment Architecture