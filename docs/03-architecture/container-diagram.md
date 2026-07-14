# Container Diagram (C4 Level 2)

**Project:** CloudCommerce

**Architecture Model:** C4 Model – Level 2 (Container)

**Version:** 1.0.0

**Status:** Draft

**Owner:** Solution Architecture Team

**Last Updated:** July 2026

---

# 1. Purpose

Dokumen ini menjelaskan bagaimana CloudCommerce dibagi menjadi beberapa container (application, database, message broker, storage, dan observability).

Pada C4 Model, **container bukan berarti Docker Container**, melainkan unit aplikasi yang dapat dijalankan secara independen.

Dokumen ini menjadi acuan untuk:

- Kubernetes Deployment
- Helm Chart
- GitHub Actions
- ArgoCD
- Service Communication
- Monitoring

---

# 2. High Level Architecture

```
                      Internet
                           │
                    Cloudflare CDN
                           │
                     HTTPS / TLS
                           │
                 NGINX Ingress Controller
                           │
                     API Gateway (Go)
                           │
        ┌──────────────────┼──────────────────┐
        │                  │                  │
        ▼                  ▼                  ▼
 Identity Service     Tenant Service     Catalog Service
                                              │
                                              ▼
                                       Inventory Service
                                              │
                                              ▼
                                         Order Service
                                              │
                                              ▼
                                        Payment Service

 Notification Service ◄────────────────────────────┐
 Analytics Service   ◄─────────────────────────────┤
 Subscription Service ◄────────────────────────────┘

                ───────── Infrastructure ─────────

                 PostgreSQL (Database per Service)

                       Redis

                    NATS JetStream

                    MinIO Object Storage

          Prometheus + Grafana + Loki + Tempo
```

---

# 3. Actors

## Buyer

Mengakses Storefront.

---

## Seller

Mengelola toko.

---

## Platform Admin

Mengelola platform.

---

## DevOps Engineer

Mengelola deployment, monitoring, scaling, dan keamanan.

---

# 4. Containers

## Web Storefront

Technology

- Next.js
- React
- TypeScript
- TailwindCSS

Responsibilities

- Public Store
- Checkout
- Product Browsing

---

## Seller Dashboard

Technology

- Next.js

Responsibilities

- Dashboard
- CRUD Product
- Orders
- Analytics

---

## API Gateway

Technology

Go + Gin

Responsibilities

- Autentikasi JWT
- Routing internal ke services
- Rate Limiting per tenant
- API Versioning
- Request Logging

Gateway tidak memiliki business logic.

## NGINX Ingress Controller

Technology

NGINX

Responsibilities

- TLS termination
- External traffic routing
- Request buffering

---

## Identity Service

Technology

Go

Responsibilities

- Authentication
- Authorization
- JWT
- Session

Database

identity_db

---

## Tenant Service

Responsibilities

- Tenant
- Branding
- Domain

Database

tenant_db

---

## Catalog Service

Responsibilities

- Product
- Category
- Media

Database

catalog_db

---

## Inventory Service

Responsibilities

- Stock
- Reservation

Database

inventory_db

---

## Order Service

Responsibilities

- Cart
- Checkout
- Order Lifecycle

Database

order_db

---

## Payment Service

Responsibilities

- Payment
- Invoice
- Webhook

Database

payment_db

External

- Midtrans Sandbox
- Xendit Sandbox

---

## Subscription Service

Responsibilities

- Billing
- Plan
- Renewal

Database

subscription_db

---

## Notification Service

Responsibilities

- Email
- Notification

Database

notification_db

---

## Analytics Service

Responsibilities

- Dashboard
- Reporting
- KPI

Database

analytics_db

---

# 5. Infrastructure Containers

## PostgreSQL

Satu database untuk setiap service.

```
identity_db

tenant_db

catalog_db

inventory_db

order_db

payment_db

subscription_db

notification_db

analytics_db
```

---

## Redis

Digunakan untuk:

- Cache
- Session
- Rate Limiting
- Distributed Lock

---

## NATS JetStream

Digunakan untuk:

- Domain Events
- Async Communication
- Retry
- Event Replay (future)

---

## MinIO

Digunakan untuk:

- Product Images
- Invoice PDF
- Attachment

---

# 6. External Systems

Payment Gateway

- Midtrans
- Xendit

Email Provider

- Resend
- SendGrid

DNS/CDN

- Cloudflare

OAuth (Future)

- Google
- GitHub

---

# 7. Communication

## Browser

↓

HTTPS

↓

Ingress

↓

API Gateway

↓

REST

↓

Microservices

---

## Internal

REST

Untuk request sinkron.

---

NATS

Untuk event.

---

Object Storage

S3 API.

---

PostgreSQL

Internal Cluster Only.

---

# 8. Deployment

Semua container dijalankan sebagai Kubernetes Deployment.

Masing-masing memiliki:

- Deployment
- Service
- ConfigMap
- Secret
- HPA
- NetworkPolicy

---

# 9. Observability

Semua service mengirim:

Metrics

↓

Prometheus

↓

Grafana

---

Logs

↓

Loki

↓

Grafana

---

Trace

↓

Tempo

↓

Grafana

---

# 10. CI/CD Flow

Developer

↓

GitHub

↓

GitHub Actions

↓

Container Registry

↓

ArgoCD

↓

Kubernetes Cluster

↓

Production

---

# 11. Security

Semua trafik:

HTTPS

↓

Ingress

↓

API Gateway

↓

JWT

↓

RBAC

↓

Service

Database tidak dapat diakses dari internet.

---

# 12. Scaling Strategy

| Container | Scaling |
|------------|---------|
| API Gateway | Horizontal |
| Catalog | Horizontal |
| Inventory | Horizontal |
| Order | Horizontal |
| Payment | Horizontal |
| Notification | Worker Based |
| Analytics | Worker Based |
| PostgreSQL | Vertical + Read Replica |
| Redis | Sentinel (future) |
| NATS | Cluster |

---

# 13. Related Documents

- System Context
- Domain Model
- Event Storming
- Service Boundaries
- Database Design
- Deployment Architecture