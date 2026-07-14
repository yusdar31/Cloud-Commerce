# User Flow

**Project:** CloudCommerce

**Document Version:** 1.0.0

**Status:** Draft

**Author:** Engineering Team

**Last Updated:** July 2026

---

# 1. Overview

## Purpose

Dokumen ini mendefinisikan alur pengguna (User Flow) untuk seluruh fitur utama CloudCommerce.

Tujuan User Flow adalah:

- Memastikan setiap user dapat menyelesaikan task dengan langkah sesedikit mungkin.
- Menjadi referensi Wireframe.
- Menjadi referensi API Design.
- Menjadi referensi Backend Workflow.
- Menjadi referensi QA Test Scenario.

---

# 2. User Roles

CloudCommerce memiliki empat jenis pengguna.

| Role | Description |
|------|-------------|
| Guest | Pengunjung tanpa akun |
| Buyer | Pelanggan yang membeli produk |
| Seller | Pemilik toko |
| Platform Admin | Administrator platform |

---

# 3. Global User Journey

```
Visitor

↓

Landing Page

↓

Register / Login

↓

Role Detection

↓

Buyer
Seller
Admin

↓

Dashboard
```

---

# 4. Seller Journey

## Goal

Seller ingin membuka toko dan mulai menjual produk.

---

```
Landing

↓

Register Seller
(Tenant ID & store dibuat otomatis)

↓

Email Verification

↓

Login

↓

Setup Store Profile
(nama toko, logo, branding)

↓

Dashboard

↓

Create Category

↓

Add Product

↓

Publish Product

↓

Store Online

↓

Receive Order

↓

Ship Order

↓

Order Completed
```

---

## Happy Path

```
Register

↓

Verify Email

↓

Login

↓

Setup Store Profile

↓

Create Product

↓

Product Published

↓

Buyer Checkout

↓

Seller Receives Order

↓

Update Shipping

↓

Order Completed
```

---

## Alternative Flow

Register

↓

Email belum diverifikasi

↓

Tidak dapat login

↓

Resend Verification

↓

Verify

↓

Login

---

# 5. Buyer Journey

## Goal

Buyer membeli produk dengan proses yang sederhana.

---

```
Landing

↓

Browse Store

↓

Search Product

↓

Product Detail

↓

Add To Cart

↓

Checkout

↓

Payment

↓

Payment Success

↓

Order Tracking

↓

Completed
```

---

## Happy Path

```
Browse Product

↓

Product Detail

↓

Cart

↓

Checkout

↓

Payment Gateway

↓

Webhook

↓

Order Paid

↓

Seller Ship

↓

Buyer Receive

↓

Done
```

---

## Alternative Flow

Payment Failed

↓

Retry Payment

↓

Success

---

Payment Expired

↓

Cancel Order

---

Stock Empty

↓

Notify Buyer

↓

Return Cart

---

# 6. Guest Journey

Guest belum login.

```
Landing

↓

Browse Store

↓

View Product

↓

Add To Cart

↓

Login Required

↓

Register

↓

Login

↓

Checkout
```

---

# 7. Platform Admin Journey

```
Login

↓

Dashboard

↓

View Tenants

↓

View Subscription

↓

View Monitoring

↓

Audit Logs

↓

Resolve Incident
```

---

# 8. Product Management Flow

Seller

↓

Dashboard

↓

Products

↓

Create Product

↓

Upload Images

↓

Save Draft

↓

Publish

↓

Visible on Store

---

Edit Product

↓

Save

↓

Product Updated

---

Archive Product

↓

Hidden from Store

---

# 9. Order Processing Flow

Buyer Checkout

↓

Create Order

↓

Reserve Inventory

↓

Waiting Payment

↓

Webhook Received

↓

Payment Success

↓

Order Confirmed

↓

Seller Notification

↓

Packing

↓

Shipping

↓

Completed

---

## Failed Payment

Buyer Checkout

↓

Payment Failed

↓

Retry

↓

Success

OR

↓

Cancel

---

# 10. Inventory Flow

Product Published

↓

Inventory Created

↓

Buyer Checkout

↓

Reserve Stock

↓

Payment Success

↓

Reduce Stock

↓

Stock Updated

---

Payment Failed

↓

Release Reserved Stock

---

# 11. Subscription Flow (Seller)

Register

↓

Free Trial

↓

Trial Expired

↓

Choose Plan

↓

Payment

↓

Subscription Active

↓

Renew

↓

Expired

↓

Suspend Store

---

# 12. Authentication Flow

Guest

↓

Register

↓

Email Verification

↓

Login

↓

JWT Generated

↓

Access Dashboard

↓

Refresh Token

↓

Logout

↓

Session Destroyed

---

# 13. Password Reset Flow

Forgot Password

↓

Enter Email

↓

Receive Reset Link

↓

Open Link

↓

Create New Password

↓

Login

---

# 14. Payment Flow

Checkout

↓

Payment Gateway

↓

Waiting Payment

↓

Webhook

↓

Verify Signature

↓

Update Payment Status

↓

Publish Event

↓

Order Service

↓

Notification Service

↓

Buyer

↓

Seller

---

# 15. Notification Flow (MVP)

Order Paid

↓

NATS Event

↓

Notification Service

↓

Email
(MVP hanya email via Mailhog)

↓

Push Notification (Future)

↓

Dashboard Notification (Future)

---

# 16. Review Flow (Future Phase)

Order Completed

↓

Buyer Opens Order

↓

Leave Rating

↓

Write Review

↓

Published

↓

Visible on Product

---

# 17. Error Flow

## Unauthorized

```
User

↓

Protected Page

↓

401

↓

Login
```

---

## Forbidden

```
Seller

↓

Admin Page

↓

403

↓

Dashboard
```

---

## Product Not Found

```
URL

↓

404

↓

Back to Store
```

---

# 18. Session Flow

Login

↓

JWT

↓

Access API

↓

Token Expired

↓

Refresh Token

↓

Continue

OR

↓

Login Again

---

# 19. Event Flow (Microservices)

Buyer Checkout

↓

Order Service

↓

NATS

↓

Inventory Service

↓

Payment Service

↓

Notification Service

↓

Analytics Service

↓

Audit Service

---

# 20. Cross-Service Sequence

```
Buyer

↓

API Gateway

↓

Order Service

↓

Inventory Service

↓

Reserve Stock

↓

Payment Service

↓

Gateway

↓

Webhook

↓

Payment Service

↓

NATS

↓

Order Service

↓

Notification Service

↓

Buyer
```

---

# 21. Flow Summary

| Flow | Status |
|--------|--------|
| Seller Registration | ✅ |
| Buyer Checkout | ✅ |
| Product Management | ✅ |
| Order Management | ✅ |
| Inventory | ✅ |
| Subscription | ✅ |
| Authentication | ✅ |
| Password Reset | ✅ |
| Payment | ✅ |
| Notification | ✅ |
| Review | ✅ |

---

# 22. Related Documents

- Product Vision
- Product Brief
- PRD
- User Persona
- User Stories
- Information Architecture
- Wireframes
- API Specification
- System Architecture