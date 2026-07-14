# Event Storming

**Project:** CloudCommerce

**Architecture:** Event-Driven Microservices

**Version:** 1.0.0

**Status:** Draft

**Owner:** Solution Architecture Team

**Last Updated:** July 2026

---

# 1. Purpose

Dokumen ini mendefinisikan seluruh event bisnis utama pada CloudCommerce.

Tujuan:

- Mengidentifikasi business events
- Menentukan komunikasi antar microservice
- Mengurangi coupling
- Menentukan event contract
- Menjadi dasar implementasi NATS

Dokumen ini tidak membahas implementasi teknis.

Fokusnya adalah proses bisnis.

---

# 2. Event Storming Legend

| Color (Workshop) | Meaning |
|------------------|---------|
| Orange | Domain Event |
| Blue | Command |
| Yellow | Actor |
| Purple | Policy |
| Green | Read Model |
| Pink | External System |

Pada dokumentasi Markdown kita menggunakan label teks sebagai pengganti warna.

---

# 3. Core Business Flows

CloudCommerce memiliki beberapa alur bisnis utama:

1. Seller Registration
2. Product Management
3. Inventory Management
4. Checkout
5. Payment
6. Order Fulfillment
7. Subscription Billing
8. Notifications

---

# 4. Seller Registration Flow

## Actor

Seller

↓

## Command

Register Seller

↓

## Aggregate

Identity

↓

## Business Rules

- Email harus unik
- Password memenuhi kebijakan keamanan
- Role default adalah Seller

↓

## Domain Event

UserRegistered

↓

Tenant Service

↓

Create Tenant

↓

Domain Event

TenantCreated

↓

Subscription Service

↓

Create Trial Subscription

↓

Domain Event

TrialSubscriptionStarted

↓

Notification Service

↓

Send Welcome Email

↓

Domain Event

WelcomeEmailSent

---

# 5. Product Management Flow

Seller

↓

Create Product

↓

Catalog Aggregate

↓

Validation

- SKU unik
- Harga valid
- Nama wajib

↓

ProductCreated

↓

Inventory Service

↓

Initialize Inventory

↓

InventoryInitialized

---

# 6. Product Publishing Flow

Seller

↓

Publish Product

↓

Catalog Aggregate

↓

Validation

↓

ProductPublished

↓

Analytics

↓

Track Product Published

---

# 7. Inventory Flow

Seller

↓

Adjust Stock

↓

Inventory Aggregate

↓

Validation

↓

StockAdjusted

↓

Analytics

↓

InventoryUpdated

↓

Notification

↓

Low Stock Alert (jika threshold tercapai)

---

# 8. Checkout Flow

Buyer

↓

Create Cart

↓

CartCreated

↓

Add Product

↓

CartUpdated

↓

Checkout

↓

Order Aggregate

↓

Reserve Inventory

↓

StockReserved

↓

Create Order

↓

OrderCreated

↓

Payment Requested

↓

PaymentRequested

---

# 9. Payment Flow

Payment Service

↓

Generate Invoice

↓

InvoiceCreated

↓

Call Payment Gateway

↓

Waiting Payment

↓

Webhook Received

↓

PaymentSucceeded

↓

Publish Event

↓

OrderConfirmed

↓

InventoryCommitted

↓

Notification Service

↓

Send Receipt

↓

Analytics Updated

---

# 10. Payment Failure Flow

PaymentRequested

↓

Gateway Response

↓

PaymentFailed

↓

Release Reserved Stock

↓

StockReleased

↓

OrderCancelled

↓

Notify Buyer

---

# 11. Refund Flow

Seller

↓

Refund Requested

↓

Refund Validation

↓

RefundRequested

↓

Gateway

↓

RefundCompleted

↓

OrderRefunded

↓

Notification

---

# 12. Subscription Flow

Seller

↓

Choose Plan

↓

Subscription Aggregate

↓

SubscriptionStarted

↓

InvoiceCreated

↓

PaymentSucceeded

↓

SubscriptionActivated

---

# 13. Notification Flow

Business Event

↓

Notification Policy

↓

Queue Notification

↓

NotificationQueued

↓

Email Provider

↓

NotificationSent

---

# 14. Analytics Flow

Business Event

↓

Analytics Consumer

↓

Update Dashboard

↓

DashboardUpdated

Semua analytics bersifat asynchronous.

---

# 15. Domain Events

## Identity

- UserRegistered
- UserLoggedIn
- PasswordChanged

---

## Tenant

- TenantCreated
- TenantActivated
- TenantSuspended

---

## Catalog

- ProductCreated
- ProductUpdated
- ProductPublished
- ProductArchived

---

## Inventory

- InventoryInitialized
- StockReserved
- StockReleased
- StockAdjusted
- InventoryLow

---

## Order

- CartCreated
- CartUpdated
- OrderCreated
- OrderConfirmed
- OrderCancelled
- OrderCompleted

---

## Payment

- PaymentRequested
- InvoiceCreated
- PaymentSucceeded
- PaymentFailed
- RefundCompleted

---

## Subscription

- SubscriptionStarted
- SubscriptionRenewed
- SubscriptionExpired

---

## Notification

- NotificationQueued
- NotificationSent
- NotificationFailed

---

## Analytics

- DashboardUpdated
- DailyReportGenerated

---

# 16. Event Naming Convention

Gunakan bentuk lampau (Past Tense).

✔ UserRegistered

✔ OrderCreated

✔ PaymentSucceeded

Bukan

✖ RegisterUser

✖ CreateOrder

✖ Pay

---

# 17. Event Payload Guidelines

Minimal setiap event memiliki:

- Event ID
- Event Type
- Aggregate ID
- Tenant ID
- Timestamp
- Version
- Correlation ID
- Payload

---

Contoh

```json
{
  "eventId": "...",
  "eventType": "OrderCreated",
  "tenantId": "...",
  "aggregateId": "...",
  "timestamp": "...",
  "correlationId": "...",
  "version": 1,
  "payload": {}
}
```

---

# 18. Event Ordering

CloudCommerce tidak menjamin urutan event antar service.

Consumer harus:

- idempotent
- retryable
- tolerant terhadap duplicate event

---

# 19. Error Handling

Jika consumer gagal:

- Retry
- Dead Letter Queue (future)
- Alert

Tidak melakukan rollback lintas service.

---

# 20. Long Running Transactions

CloudCommerce menggunakan pola Saga.

Contoh Checkout:

Reserve Stock

↓

Create Order

↓

Payment

↓

Confirm Order

Jika Payment gagal:

↓

Release Stock

↓

Cancel Order

---

# 21. Event Bus

Transport:

NATS JetStream

Semua domain berkomunikasi melalui event.

REST hanya digunakan untuk operasi sinkron yang membutuhkan respons langsung.

---

# 22. Outbox Pattern

Setiap service menggunakan Outbox Pattern untuk memastikan perubahan database dan publish event berlangsung konsisten.

---

# 23. Related Documents

- Domain Model
- Service Boundaries
- API Guidelines
- Deployment Architecture