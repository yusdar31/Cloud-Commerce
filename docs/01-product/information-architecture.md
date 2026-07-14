# Information Architecture (IA)

**Project:** CloudCommerce  
**Document Version:** 1.0.0  
**Status:** Draft  
**Author:** Engineering Team  
**Last Updated:** July 2026

---

# 1. Overview

## Purpose

Dokumen ini mendefinisikan struktur informasi, navigasi, hierarki halaman, serta hak akses pengguna pada platform **CloudCommerce**.

Information Architecture menjadi fondasi bagi:

- UX Design
- Wireframes
- Frontend Routing
- Backend Authorization
- API Design
- Navigation Design

---

# 2. Design Principles

CloudCommerce mengikuti beberapa prinsip berikut:

### 1. Simple Navigation

Pengguna tidak boleh membutuhkan lebih dari **3 klik** untuk mencapai fitur utama.

---

### 2. Role Based Navigation

Menu hanya muncul sesuai role pengguna.

Contoh:

Buyer tidak akan melihat menu Analytics.

Seller tidak akan melihat menu Infrastructure.

---

### 3. Scalable

Menu dapat bertambah tanpa mengubah struktur utama.

Misalnya nanti ditambahkan:

- AI Recommendation
- Marketing Automation
- Affiliate
- CRM

---

### 4. Mobile First

Semua halaman harus memiliki versi mobile.

---

# 3. User Roles

CloudCommerce memiliki empat role utama.

| Role | Description |
|---------|-------------------------------|
| Guest | Pengunjung tanpa login |
| Buyer | Pelanggan |
| Seller | Pemilik toko |
| Platform Admin | Administrator platform |

---

# 4. Global Sitemap

```
CloudCommerce

├── Public Website
│
├── Authentication
│
├── Seller Dashboard
│
├── Buyer Storefront
│
└── Platform Administration
```

---

# 5. Public Website

```
/

├── Landing
├── Features
├── Pricing
├── Documentation
├── About
├── Contact
├── Login
├── Register
├── Forgot Password
└── Verify Email
```

---

## Navigation

```
Logo

Features

Pricing

Documentation

Contact

Login

Register
```

---

# 6. Authentication

```
Authentication

├── Login
├── Register Seller
├── Register Buyer
├── Verify Email
├── Forgot Password
├── Reset Password
└── Logout
```

---

# 7. Seller Dashboard

```
Dashboard

├── Overview
│
├── Products
│   ├── Product List
│   ├── Product Detail
│   ├── Create Product
│   ├── Edit Product
│   └── Archive Product
│
├── Categories
│
├── Orders
│   ├── Order List
│   └── Order Detail
│
├── Customers
│
├── Inventory
│
├── Reviews
│
├── Discounts
│
├── Analytics
│
├── Billing
│
├── User Management
│
├── Store Settings
│
└── Audit Logs
```

---

## Seller Navigation

```
Dashboard

Products

Orders

Customers

Inventory

Analytics

Billing

Settings
```

---

# 8. Buyer Storefront

```
Store

├── Home
├── Categories
├── Search
├── Product Detail
├── Cart
├── Checkout
├── Payment
├── Order Success
├── Order Failed
├── Order History
├── Wishlist
├── Notifications
└── Profile
```

---

## Buyer Navigation

```
Home

Categories

Search

Orders

Wishlist

Profile

Cart
```

---

# 9. Platform Administration

```
Platform

├── Dashboard
├── Tenant Management
├── User Management
├── Subscription
├── Infrastructure
├── Monitoring
├── Audit Logs
├── Feature Flags
├── Support
└── System Settings
```

---

## Admin Navigation

```
Dashboard

Tenants

Subscriptions

Infrastructure

Monitoring

Logs

Support
```

---

# 10. URL Convention

## Public

```
/

 /pricing

 /features

 /documentation

 /contact
```

---

## Authentication

```
/login

/register

/forgot-password

/reset-password
```

---

## Seller

```
/dashboard

/dashboard/products

/dashboard/products/new

/dashboard/products/{id}

/dashboard/orders

/dashboard/orders/{id}

/dashboard/customers

/dashboard/inventory

/dashboard/analytics

/dashboard/billing

/dashboard/settings
```

---

## Buyer

```
/store/{slug}

/store/{slug}/products

/store/{slug}/products/{id}

/cart

/checkout

/orders

/profile
```

---

## Admin

```
/admin

/admin/tenants

/admin/users

/admin/subscriptions

/admin/infrastructure

/admin/monitoring

/admin/logs
```

---

# 11. Breadcrumb Rules

Example

```
Dashboard

↓

Products

↓

Product Detail
```

Displayed as

```
Dashboard

>

Products

>

MacBook Pro M4
```

---

# 12. Navigation Hierarchy

```
Level 1

Dashboard

Products

Orders

Customers

Analytics

Billing

Settings

-------------------

Level 2

Product Detail

Order Detail

User Detail

-------------------

Level 3

Edit Product

Edit Customer

Invoice Detail
```

---

# 13. Permission Matrix

| Module | Guest | Buyer | Seller | Admin |
|---------|:----:|:-----:|:------:|:-----:|
| Landing | ✅ | ✅ | ✅ | ✅ |
| Login | ✅ | ✅ | ✅ | ✅ |
| Register | ✅ | ❌ | ❌ | ❌ |
| Store | ✅ | ✅ | ✅ | ✅ |
| Product Detail | ✅ | ✅ | ✅ | ✅ |
| Cart | ❌ | ✅ | ❌ | ❌ |
| Checkout | ❌ | ✅ | ❌ | ❌ |
| Orders | ❌ | ✅ | ✅ | ✅ |
| Dashboard | ❌ | ❌ | ✅ | ❌ |
| Analytics | ❌ | ❌ | ✅ | ✅ |
| Billing | ❌ | ❌ | ✅ | ❌ |
| Tenant Management | ❌ | ❌ | ❌ | ✅ |
| Infrastructure | ❌ | ❌ | ❌ | ✅ |

---

# 14. Search Strategy

CloudCommerce menyediakan search pada beberapa modul.

### Public

- Product Search

---

### Seller

- Product Search
- Order Search
- Customer Search

---

### Admin

- Tenant Search
- User Search
- Log Search

---

# 15. Naming Convention

Gunakan istilah yang konsisten.

| Wrong | Correct |
|----------|----------------|
| Shop | Store |
| Item | Product |
| Client | Buyer |
| Merchant | Seller |
| Admin Panel | Dashboard |
| Invoice List | Billing |

---

# 16. Empty State

Semua halaman harus memiliki Empty State.

Contoh:

Products

```
No products found.

Create your first product.

[ Add Product ]
```

---

# 17. Error State

Contoh

```
Unable to load products.

Retry
```

---

# 18. Loading State

Semua halaman menggunakan Skeleton Loader.

Tidak menggunakan Spinner sebagai loading utama.

---

# 19. Responsive Navigation

## Desktop

Sidebar Navigation

---

## Tablet

Collapsible Sidebar

---

## Mobile

Bottom Navigation

+

Hamburger Menu

---

# 20. Future Scalability

Struktur IA telah dipersiapkan untuk mendukung:

- AI Recommendation
- Affiliate Program
- Marketplace
- Multi Warehouse
- CRM
- Loyalty Program
- Marketing Automation
- Multi Currency
- International Shipping

Tanpa mengubah struktur navigasi utama.

---

# 21. Related Documents

- Product Vision
- Product Brief
- PRD
- User Persona
- User Stories
- User Flow
- Wireframes
- System Architecture