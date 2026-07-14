# API Guidelines

**Project:** CloudCommerce

**Version:** 1.0.0

**Status:** Draft

**Architecture:** REST API + Event Driven

**Specification:** OpenAPI 3.1

**Owner:** Backend Architecture Team

**Last Updated:** July 2026

---

# 1. Purpose

Dokumen ini mendefinisikan standar API yang wajib digunakan oleh seluruh microservice CloudCommerce.

Tujuan:

- Konsisten
- Mudah dipahami
- Mudah diuji
- Mudah didokumentasikan
- Aman
- Mendukung evolusi API

Semua service harus mengikuti dokumen ini.

---

# 2. API Philosophy

CloudCommerce menggunakan prinsip:

- Resource-Oriented API
- Stateless
- Versioned
- Secure by Default
- Backward Compatible
- Developer Friendly

---

# 3. Base URL

Production

```

https://api.cloudcommerce.com

```

Development

```

http://localhost:8080

```

---

# 4. API Versioning

Semua endpoint menggunakan URI versioning.

```

/api/v1

```

Contoh

```

GET /api/v1/products

POST /api/v1/orders

```

Versi lama tetap dipertahankan selama masa deprecation.

---

# 5. Resource Naming

Gunakan noun.

✔

```

/products

/orders

/users

```

Bukan

✖

```

/createProduct

/getOrders

/deleteUser

```

---

# 6. HTTP Methods

GET

Membaca data

POST

Membuat resource

PUT

Replace seluruh resource

PATCH

Update sebagian

DELETE

Soft Delete (default)

---

# 7. Status Codes

| Code | Meaning |
|------|---------|
|200|Success|
|201|Created|
|202|Accepted|
|204|No Content|
|400|Bad Request|
|401|Unauthorized|
|403|Forbidden|
|404|Not Found|
|409|Conflict|
|422|Validation Error|
|429|Too Many Requests|
|500|Internal Server Error|

---

# 8. Request Headers

Authorization

```

Bearer <JWT>

```

Content-Type

```

application/json

```

Accept

```

application/json

```

Correlation ID

```

X-Correlation-ID

```

Idempotency

```

Idempotency-Key

```

Tenant

```

X-Tenant-ID

```

---

# 9. Response Format

Semua endpoint menggunakan struktur berikut.

```json
{
  "data": {},
  "meta": {},
  "links": {}
}
```

---

Contoh

```json
{
  "data": {
    "id": "prd_001",
    "name": "MacBook Pro"
  }
}
```

---

# 10. Error Format

Gunakan RFC 9457 Problem Details.

Contoh

```json
{
  "type": "https://api.cloudcommerce.com/errors/validation",
  "title": "Validation Failed",
  "status": 422,
  "detail": "SKU already exists",
  "instance": "/api/v1/products",
  "traceId": "abc123"
}
```

---

# 11. Validation

Semua validasi dilakukan di server.

Gunakan:

- Required
- Min Length
- Max Length
- Regex
- Enum
- Custom Rule

---

# 12. Pagination

Gunakan Cursor Pagination.

Contoh

```

GET /products?limit=20&cursor=abc123

```

Response

```json
{
  "data": [],
  "meta": {
    "nextCursor": "xyz789",
    "hasMore": true
  }
}
```

---

# 13. Filtering

Gunakan query parameter.

```

GET /products?status=published

```

```

GET /orders?paymentStatus=paid

```

---

# 14. Sorting

```

GET /products?sort=name

```

Descending

```

GET /products?sort=-createdAt

```

---

# 15. Searching

Gunakan parameter:

```

?q=macbook

```

Contoh

```

GET /products?q=macbook

```

---

# 16. Field Selection

```

GET /products?fields=id,name,price

```

---

# 17. Includes

```

GET /orders?include=items,payment

```

---

# 18. Authentication

Menggunakan JWT Access Token.

Refresh Token disimpan secara aman.

Role

- Seller
- Buyer
- Admin

---

# 19. Authorization

RBAC.

Contoh

Seller

↓

Hanya dapat mengakses tenant miliknya.

Admin

↓

Seluruh tenant.

---

# 20. Multi-Tenant Rules

Semua request harus mengetahui konteks tenant.

Melalui:

- JWT Claim
- X-Tenant-ID (internal)
- Store Domain

Service tidak boleh mengakses data tenant lain.

---

# 21. Idempotency

Endpoint berikut wajib mendukung Idempotency-Key.

- Create Order
- Create Payment
- Subscription Billing
- Refund

---

# 22. Rate Limiting

Default

100 request / menit

Endpoint sensitif

- Login
- Register
- Checkout

memiliki limit lebih ketat.

---

# 23. File Upload

Multipart Form Data.

Response

```json
{
  "url": "...",
  "key": "...",
  "size": 102400
}
```

File disimpan di MinIO/S3.

---

# 24. Time Format

Gunakan ISO-8601 UTC.

```

2026-07-14T10:15:30Z

```

---

# 25. Money Format

Tidak menggunakan floating point.

Gunakan integer.

Contoh

```json
{
  "amount": 150000,
  "currency": "IDR"
}
```

---

# 26. API Documentation

Semua service wajib memiliki:

```

/openapi.json

/docs

```

OpenAPI menjadi sumber kebenaran.

---

# 27. Health Endpoints

```

GET /healthz

GET /readyz

GET /livez

GET /metrics

```

---

# 28. Observability Headers

Semua request membawa:

- Trace ID
- Correlation ID
- Request ID

Untuk mendukung OpenTelemetry.

---

# 29. Security

- HTTPS Only
- JWT
- Rate Limiting
- Input Validation
- Output Encoding
- Audit Logging
- Security Headers

---

# 30. Deprecation Policy

Endpoint lama:

- diberi header `Deprecation`
- didokumentasikan
- dipertahankan selama periode migrasi

---

# 31. OpenAPI Standards

Setiap service wajib menyediakan:

- OpenAPI 3.1
- Example Request
- Example Response
- Error Response
- Security Scheme

---

# 32. API Review Checklist

Sebelum endpoint dirilis:

- Resource menggunakan noun
- HTTP Method sesuai
- Status code benar
- Validation lengkap
- Error mengikuti RFC 9457
- Mendukung observability
- Memiliki dokumentasi OpenAPI
- Memiliki integration test

---

# 33. Future Improvements

- GraphQL Gateway (opsional)
- gRPC untuk komunikasi internal tertentu
- API SDK Generator
- Async API Specification untuk event NATS