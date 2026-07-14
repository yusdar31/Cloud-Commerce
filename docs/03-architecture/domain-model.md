# Domain Model

**Project:** CloudCommerce

**Architecture:** Domain Driven Design (DDD)

**Version:** 1.0.0

**Status:** Draft

**Owner:** Solution Architecture Team

**Last Updated:** July 2026

---

# 1. Purpose

Dokumen ini mendefinisikan domain bisnis CloudCommerce.

Tujuannya adalah:

- Memahami objek bisnis utama
- Menentukan batas setiap domain
- Menghindari coupling antar service
- Menjadi dasar Event Storming
- Menjadi dasar Database Design
- Menjadi dasar Service Boundaries

Dokumen ini **bukan** ERD (Entity Relationship Diagram).

Dokumen ini berfokus pada bahasa bisnis (Ubiquitous Language).

---

# 2. Ubiquitous Language

Istilah berikut harus digunakan secara konsisten di seluruh proyek.

| Term | Description |
|-------|-------------|
| Tenant | Toko yang terdaftar di platform |
| Seller | Pemilik tenant |
| Buyer | Pelanggan yang membeli produk |
| Product | Barang yang dijual |
| Inventory | Persediaan produk |
| Cart | Keranjang belanja |
| Order | Transaksi pembelian |
| Payment | Pembayaran pesanan |
| Subscription | Paket langganan tenant |
| Notification | Pesan sistem |
| Event | Peristiwa bisnis yang dipublikasikan |

---

# 3. Core Domain

Core Domain adalah pembeda utama CloudCommerce.

```
Commerce Platform

↓

Catalog

Inventory

Order

Payment
```

Domain ini harus menjadi prioritas pengembangan.

---

# 4. Supporting Domain

Mendukung Core Domain.

- Identity
- Tenant
- Notification
- Analytics
- Subscription

---

# 5. Generic Domain

Menggunakan solusi umum.

- Authentication
- Email
- Logging
- Metrics
- Storage

---

# 6. Bounded Context

CloudCommerce dibagi menjadi beberapa bounded context.

```
Identity

Tenant

Catalog

Inventory

Order

Payment

Subscription

Notification

Analytics
```

Setiap bounded context nantinya menjadi microservice independen.

---

# 7. Identity Context

## Responsibilities

- Login
- Register
- JWT
- Refresh Token
- Password Reset
- Role Management

---

## Aggregate

User

---

## Entities

- User
- Role
- Permission
- Session

---

## Value Objects

- Email
- Password Hash

---

## Domain Events

- UserRegistered
- UserLoggedIn
- PasswordChanged

---

# 8. Tenant Context

## Responsibilities

- Tenant Registration
- Tenant Settings
- Branding
- Domain Configuration
- Subscription Ownership

---

## Aggregate

Tenant

---

## Entities

- Tenant
- TenantSettings
- Domain
- Brand

---

## Value Objects

- TenantId
- StoreSlug

---

## Domain Events

- TenantCreated
- TenantActivated
- TenantSuspended

---

# 9. Catalog Context

## Responsibilities

- Product
- Category
- Product Images
- Product Search
- Product Status

---

## Aggregate

Product

---

## Entities

- Product
- Category
- ProductImage

---

## Value Objects

- SKU
- Price

---

## Domain Events

- ProductCreated
- ProductUpdated
- ProductPublished
- ProductArchived

---

# 10. Inventory Context

## Responsibilities

- Stock
- Reservation
- Stock Movement
- Low Stock Alert

---

## Aggregate

Inventory

---

## Entities

- Inventory
- StockMovement
- Reservation

---

## Value Objects

- Quantity

---

## Domain Events

- StockReserved
- StockReleased
- StockAdjusted
- InventoryLow

---

# 11. Order Context

## Responsibilities

- Cart
- Checkout
- Order
- Order Status

---

## Aggregate

Order

---

## Entities

- Cart
- CartItem
- Order
- OrderItem

---

## Value Objects

- OrderNumber
- ShippingAddress
- Money

---

## Domain Events

- OrderCreated
- CheckoutStarted
- CheckoutCompleted
- OrderCancelled

---

# 12. Payment Context

## Responsibilities

- Payment
- Invoice
- Webhook
- Refund

---

## Aggregate

Payment

---

## Entities

- Payment
- Invoice
- Refund

---

## Value Objects

- TransactionId
- PaymentMethod

---

## Domain Events

- PaymentRequested
- PaymentSucceeded
- PaymentFailed
- RefundCompleted

---

# 13. Subscription Context

## Responsibilities

- Billing
- Plan
- Quota
- Renewal

---

## Aggregate

Subscription

---

## Entities

- Subscription
- Plan
- Invoice

---

## Value Objects

- BillingCycle

---

## Domain Events

- SubscriptionStarted
- SubscriptionRenewed
- SubscriptionExpired

---

# 14. Notification Context

## Responsibilities

- Email
- Webhook
- In-App Notification

---

## Aggregate

Notification

---

## Entities

- Notification
- EmailTemplate

---

## Domain Events

- NotificationQueued
- NotificationSent
- NotificationFailed

---

# 15. Analytics Context

## Responsibilities

- Sales Analytics
- Dashboard
- KPI
- Reporting

---

## Aggregate

AnalyticsSnapshot

---

## Entities

- DashboardMetric
- SalesReport

---

## Domain Events

- DailyReportGenerated
- DashboardUpdated

---

# 16. Context Relationship

```
Identity

↓

Tenant

↓

Catalog

↓

Inventory

↓

Order

↓

Payment

↓

Notification
```

Analytics menerima event dari semua domain.

---

# 17. Aggregate Rules

Setiap Aggregate:

- Memiliki satu Root.
- Menjaga konsistensi internal.
- Tidak boleh diubah langsung oleh aggregate lain.

Contoh:

Inventory hanya dapat diubah melalui Inventory Aggregate.

---

# 18. Shared Kernel

Shared package hanya boleh berisi:

- Event Contracts
- Common Types
- Error Codes
- Utilities

Tidak boleh berisi business logic.

---

# 19. Anti-Corruption Layer

Integrasi dengan sistem eksternal menggunakan adapter.

Contoh:

Payment Service

↓

Payment Adapter

↓

Midtrans

atau

↓

Xendit

Sehingga vendor dapat diganti tanpa mengubah domain.

---

# 20. Design Rules

- Setiap microservice memiliki satu bounded context.
- Tidak ada akses langsung ke database service lain.
- Komunikasi sinkron melalui REST/gRPC.
- Komunikasi asinkron melalui NATS.
- Domain event bersifat immutable.

---

# 21. Future Domains

Roadmap berikutnya dapat menambahkan:

- Review
- Wishlist
- Coupon
- Shipment
- Loyalty
- CRM
- AI Recommendation
- Marketplace
- POS

Tanpa mengubah bounded context utama.

---

# 22. Related Documents

- System Context
- Event Storming
- Service Boundaries
- Database Design
- API Guidelines

| Domain       | Owner                | Database       | Events Published       | Events Consumed              |
| ------------ | -------------------- | -------------- | ---------------------- | ---------------------------- |
| Identity     | Identity Service     | identity_db    | UserRegistered         | -                            |
| Tenant       | Tenant Service       | tenant_db      | TenantCreated          | UserRegistered               |
| Catalog      | Catalog Service      | catalog_db     | ProductPublished       | TenantCreated                |
| Inventory    | Inventory Service    | inventory_db   | StockReserved          | ProductCreated, OrderCreated |
| Order        | Order Service        | order_db       | OrderCreated           | StockReserved, PaymentSucceeded |
| Payment      | Payment Service      | payment_db     | PaymentSucceeded       | OrderCreated                 |
| Subscription | Subscription Service | subscription_db | SubscriptionStarted    | TenantCreated, PaymentSucceeded |
| Notification | Notification Service | notification_db | NotificationSent       | UserRegistered, OrderCreated, PaymentSucceeded |
| Analytics    | Analytics Service    | analytics_db   | DashboardUpdated       | ProductPublished, OrderCreated, PaymentSucceeded |
