# System Context

**Project:** CloudCommerce

**Architecture Model:** C4 Model - Level 1 (System Context)

**Version:** 1.0.0

**Status:** Draft

**Owner:** Solution Architecture Team

**Last Updated:** July 2026

---

# 1. Overview

## Purpose

Dokumen ini menjelaskan posisi CloudCommerce di dalam ekosistem sistem yang lebih besar.

System Context membantu seluruh tim memahami:

- siapa saja pengguna sistem
- sistem eksternal yang berinteraksi
- batas (boundary) aplikasi
- dependency utama
- integrasi yang diperlukan

Dokumen ini merupakan **Level 1** dari C4 Model.

---

# 2. Vision

CloudCommerce adalah platform SaaS multi-tenant yang memungkinkan merchant membuat toko online sendiri, mengelola produk, menerima pembayaran, dan memonitor bisnis melalui dashboard modern berbasis cloud.

Platform dirancang menggunakan arsitektur **Cloud Native Microservices** sehingga mampu berkembang secara horizontal, mudah di-maintain, dan siap dijalankan di Kubernetes.

---

# 3. Goals

CloudCommerce harus mampu:

- melayani banyak tenant dalam satu platform
- mendukung ribuan transaksi
- mudah di-scale
- memiliki deployment otomatis
- mudah dimonitor
- aman secara default
- mudah dikembangkan oleh banyak tim

---

# 4. System Boundary

```
                    Internet

        +---------------------------+
        |                           |
        |      CloudCommerce        |
        |                           |
        |  Frontend                 |
        |  API Gateway              |
        |  Microservices            |
        |  Database                 |
        |                           |
        +---------------------------+

Outside Boundary

Payment Gateway

SMTP

Object Storage

Monitoring

GitHub

Kubernetes Cluster
```

---

# 5. Primary Actors

## Buyer

Pelanggan yang membeli produk.

Responsibilities

- Browse produk
- Checkout
- Payment
- Review produk

---

## Seller

Pemilik toko.

Responsibilities

- CRUD Product
- Kelola Order
- Analytics
- Inventory
- Subscription

---

## Platform Admin

Administrator platform.

Responsibilities

- Tenant Management
- Monitoring
- User Management
- Subscription
- Incident Handling

---

## DevOps Engineer

Mengelola infrastruktur platform.

Responsibilities

- Deployment
- Scaling
- Monitoring
- Security
- Backup
- Disaster Recovery

---

# 6. External Systems

CloudCommerce berintegrasi dengan beberapa sistem eksternal.

---

## Payment Gateway

Contoh

- Midtrans Sandbox
- Xendit Sandbox

Purpose

- Payment Processing
- Webhook
- Refund

Communication

REST API

Webhook

---

## Email Provider

Contoh

- Resend
- SendGrid

Purpose

- Email Verification
- Password Reset
- Notification

---

## Object Storage

Contoh

- MinIO
- Amazon S3

Purpose

- Product Images
- Invoice PDF
- Attachment

---

## DNS

Contoh

Cloudflare

Purpose

- DNS
- SSL
- CDN
- WAF (future)

---

## Identity Provider (Future)

OAuth

Google

GitHub

Microsoft

---

## Monitoring Stack

Prometheus

Grafana

Loki

Tempo

Purpose

- Metrics
- Logs
- Traces

---

# 7. Internal Systems

CloudCommerce terdiri dari beberapa domain utama.

```
Frontend

↓

API Gateway

↓

Identity

Tenant

Catalog

Inventory

Order

Payment

Notification

Analytics
```

Masing-masing domain akan berkembang menjadi microservice independen.

---

# 8. Context Diagram

```
                          Buyer
                            │
                            │
                    Browse / Checkout
                            │
                            ▼

              +----------------------------+
              |                            |
              |      CloudCommerce         |
              |                            |
              +----------------------------+
                 ▲      ▲        ▲
                 │      │        │

            Seller   Admin   DevOps

                 │      │        │

────────────── External Systems ──────────────

Payment Gateway

SMTP

Object Storage

Cloudflare

GitHub

Kubernetes

Monitoring Stack
```

---

# 9. High-Level Interaction

## Seller Flow

Seller

↓

Frontend

↓

API Gateway

↓

Catalog Service

↓

Database

---

## Buyer Flow

Buyer

↓

Storefront

↓

Gateway

↓

Order Service

↓

Payment Gateway

↓

Webhook

↓

Order Updated

---

## Admin Flow

Admin

↓

Dashboard

↓

Monitoring

↓

Tenant Management

---

# 10. Communication Patterns

| Communication | Usage |
|--------------|-------|
| HTTPS | Browser ↔ Gateway |
| REST | Sync Service Communication |
| NATS | Async Event Communication |
| Webhook | Payment Callback |
| SMTP | Email Delivery |

---

# 11. Security Boundary

Traffic masuk harus melewati:

Cloudflare

↓

Ingress Controller

↓

API Gateway

↓

Authentication

↓

Microservices

Tidak ada service yang dapat diakses langsung dari internet.

---

# 12. Trust Boundaries

External Users

↓

Internet

↓

Cloudflare

↓

CloudCommerce Cluster

↓

Internal Network

↓

Database

↓

Object Storage

Semua komunikasi antar boundary menggunakan TLS.

---

# 13. Non-Functional Considerations

Availability

Target

99.9%

---

Latency

API

<300 ms

---

Scalability

Horizontal Scaling

Kubernetes HPA

---

Security

JWT

HTTPS

Rate Limiting

RBAC

Audit Logging

---

Observability

Metrics

Logs

Tracing

---

# 14. Assumptions

- Semua service berjalan dalam Kubernetes.
- Semua komunikasi antar service menggunakan internal cluster networking.
- Setiap microservice memiliki database sendiri.
- Tidak ada shared database antar service.
- Semua deployment dilakukan melalui GitOps.

---

# 15. Architecture Decisions

## ADR-001

Decision

Menggunakan Microservices.

Reason

Mendukung scalability dan independent deployment.

---

## ADR-002

Decision

Menggunakan API Gateway.

Reason

Single entry point untuk seluruh client.

---

## ADR-003

Decision

Menggunakan Event Bus.

Reason

Mengurangi coupling antar service.

---

## ADR-004

Decision

Menggunakan Kubernetes.

Reason

Cloud-native deployment.

---

# 16. Out of Scope

Tidak termasuk:

- AI Recommendation
- Marketplace
- Affiliate
- Multi Warehouse
- ERP Integration

Semua fitur tersebut akan menjadi roadmap fase berikutnya.

---

# 17. Future Evolution

CloudCommerce dirancang agar mudah berkembang menjadi:

- Marketplace
- SaaS ERP
- POS Integration
- AI Recommendation
- AI Customer Support
- Warehouse Management
- CRM
- Loyalty Program

Tanpa mengubah arsitektur utama.

---

# 18. Related Documents

- Product Vision
- PRD
- Information Architecture
- User Flow
- Deployment Architecture
- Service Boundaries
- Event Storming