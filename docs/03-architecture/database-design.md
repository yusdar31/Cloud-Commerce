# Database Design

**Project:** CloudCommerce

**Architecture:** Database per Service

**Database:** PostgreSQL 17

**Version:** 1.0.0

**Status:** Draft

**Owner:** Database Architecture Team

**Last Updated:** July 2026

---

# 1. Purpose

Dokumen ini mendefinisikan strategi desain database untuk seluruh microservice CloudCommerce.

Tujuan utama:

- Memisahkan database setiap service
- Menjaga independensi deployment
- Mengurangi coupling
- Mendukung horizontal scaling
- Menjadi dasar migration dan backup strategy

Dokumen ini **bukan** ERD implementasi, melainkan panduan arsitektur data.

---

# 2. Database Architecture Principles

CloudCommerce mengikuti prinsip:

- Database per Service
- No Shared Database
- Schema Evolution
- Immutable Events
- Auditability
- Soft Delete
- UUID/ULID sebagai Primary Key

---

# 3. Database Landscape

```
Identity Service
    │
identity_db

Tenant Service
    │
tenant_db

Catalog Service
    │
catalog_db

Inventory Service
    │
inventory_db

Order Service
    │
order_db

Payment Service
    │
payment_db

Subscription Service
    │
subscription_db

Notification Service
    │
notification_db

Analytics Service
    │
analytics_db
```

Setiap database berjalan secara independen.

---

# 4. Shared Rules

Semua tabel wajib memiliki field berikut:

```
id

created_at

updated_at

deleted_at

created_by

updated_by

version
```

---

# 5. Primary Key Strategy

Gunakan UUID v7 atau ULID.

Contoh

```
usr_01J6...

prd_01J6...

ord_01J6...
```

Tidak menggunakan integer auto increment.

---

# 6. Audit Fields

Semua entity memiliki:

```
created_at

updated_at

deleted_at

created_by

updated_by
```

Soft delete menggunakan `deleted_at`.

---

# 7. Identity Database

Database

```
identity_db
```

Tabel

```
users

roles

permissions

sessions

refresh_tokens
```

Owner

Identity Service

---

# 8. Tenant Database

Database

```
tenant_db
```

Tabel

```
tenants

tenant_settings

tenant_domains

branding
```

Owner

Tenant Service

---

# 9. Catalog Database

Database

```
catalog_db
```

Tabel

```
products

categories

product_images

product_tags
```

Indexes

```
tenant_id

sku

status

slug

created_at
```

---

# 10. Inventory Database

Database

```
inventory_db
```

Tabel

```
inventories

stock_movements

stock_reservations
```

Indexes

```
product_id

tenant_id

updated_at
```

---

# 11. Order Database

Database

```
order_db
```

Tabel

```
carts

cart_items

orders

order_items
```

Indexes

```
order_number

buyer_id

tenant_id

status

created_at
```

---

# 12. Payment Database

Database

```
payment_db
```

Tabel

```
payments

invoices

refunds

payment_webhooks
```

Indexes

```
transaction_id

invoice_number

payment_status

tenant_id
```

---

# 13. Subscription Database

Database

```
subscription_db
```

Tabel

```
plans

subscriptions

billing_cycles

subscription_invoices
```

---

# 14. Notification Database

Database

```
notification_db
```

Tabel

```
notifications

email_templates

email_logs
```

---

# 15. Analytics Database

Database

```
analytics_db
```

Tabel

```
dashboard_metrics

daily_reports

sales_snapshots
```

Analytics menggunakan data hasil event, bukan query langsung ke database service lain.

---

# 16. Naming Convention

Table

snake_case

Column

snake_case

Primary Key

```
id
```

Foreign Key

```
user_id

product_id

tenant_id
```

---

# 17. Relationships

Relationship antar service tidak menggunakan Foreign Key database.

Contoh:

Order Service hanya menyimpan

```
product_id
tenant_id
buyer_id
```

Validasi dilakukan melalui API atau event.

---

# 18. Index Strategy

Semua tabel minimal memiliki index pada:

- id
- tenant_id
- created_at
- updated_at

Tambahkan composite index untuk query yang sering digunakan.

Contoh:

```
(tenant_id, status)

(tenant_id, created_at)

(order_id, payment_status)
```

---

# 19. Constraints

Gunakan:

- NOT NULL
- UNIQUE
- CHECK
- ENUM (secukupnya)

Contoh:

SKU harus unik dalam satu tenant.

---

# 20. Soft Delete

Gunakan

```
deleted_at TIMESTAMP NULL
```

Data tidak dihapus secara fisik kecuali melalui proses archival.

---

# 21. Migration Strategy

Gunakan migration versioning.

Contoh

```
0001_initial.sql

0002_add_products.sql

0003_add_inventory.sql
```

Migration harus:

- idempotent
- reversible
- terdokumentasi

---

# 22. Seed Data

Pisahkan dari migration.

Contoh:

```
roles

permissions

subscription_plans

admin_user
```

---

# 23. Backup Strategy

- Daily Full Backup
- WAL Archiving
- Point-in-Time Recovery (PITR)
- Enkripsi backup
- Uji restore secara berkala

---

# 24. Performance Strategy

Gunakan:

- Pagination
- Index
- Read Replica (future)
- Connection Pooling
- Query Timeout

Hindari query N+1.

---

# 25. Multi-Tenant Strategy

Semua tabel bisnis memiliki:

```
tenant_id
```

Tenant menjadi bagian dari seluruh query.

Tidak ada query lintas tenant.

---

# 26. Security

- TLS untuk koneksi database
- Secret disimpan di Kubernetes Secret / Vault
- Least Privilege
- Audit Logging
- Row Level Security (opsional untuk fase lanjutan)

---

# 27. Data Retention

Contoh kebijakan:

- Session: 30 hari
- Notification Log: 90 hari
- Audit Log: 1 tahun
- Payment Record: sesuai regulasi

---

# 28. Future Improvements

- Read Replica
- Partitioning untuk tabel besar (orders, payments)
- CQRS Read Database
- Event Store
- Data Warehouse
- CDC (Change Data Capture)

---

# 29. Related Documents

- Domain Model
- Service Boundaries
- API Guidelines
- Deployment Architecture