# Error Handling Catalog

**Project:** CloudCommerce

**Version:** 1.0.0

**Status:** Approved

**Owner:** Engineering Team

**Last Updated:** July 2026

---

## 1. Purpose

Dokumen ini mendefinisikan **katalog error code** yang digunakan di seluruh CloudCommerce.

Tujuan:

- Konsistensi error response antar service
- Memudahkan debugging oleh frontend
- Memudahkan monitoring & alerting
- Menjadi single source of truth untuk error handling

Semua error mengikuti format **RFC 9457 Problem Details**.

---

## 2. Format Error Response

```json
{
  "type": "https://api.cloudcommerce.com/errors/{error-code}",
  "title": "Human-readable error title",
  "status": 422,
  "detail": "Penjelasan spesifik kenapa error ini terjadi",
  "instance": "/api/v1/products",
  "traceId": "abc-123-def"
}
```

### Field Definitions

| Field | Tipe | Wajib | Deskripsi |
|-------|------|-------|-----------|
| `type` | string (URI) | ✅ | URI unik untuk jenis error |
| `title` | string | ✅ | Judul human-readable, tidak berubah |
| `status` | integer | ✅ | HTTP status code |
| `detail` | string | ✅ | Penjelasan spesifik untuk request ini |
| `instance` | string | ❌ | URI endpoint yang menyebabkan error |
| `traceId` | string | ❌ | ID untuk distributed tracing |

---

## 3. Error Catalog

### 3.1 Authentication & Authorization

| Error Code | HTTP Status | Title | Contoh Detail |
|------------|-------------|-------|---------------|
| `AUTH_TOKEN_MISSING` | 401 | Authentication Required | No authorization header provided |
| `AUTH_TOKEN_INVALID` | 401 | Invalid Token | JWT signature is invalid |
| `AUTH_TOKEN_EXPIRED` | 401 | Token Expired | Access token has expired, please refresh |
| `AUTH_TOKEN_MALFORMED` | 401 | Malformed Token | Token format is not valid JWT |
| `AUTH_INSUFFICIENT_PERMISSION` | 403 | Insufficient Permission | Role 'buyer' cannot access this resource |
| `AUTH_ACCOUNT_SUSPENDED` | 403 | Account Suspended | Your account has been suspended |
| `AUTH_INVALID_CREDENTIALS` | 401 | Invalid Credentials | Email or password is incorrect |
| `AUTH_EMAIL_NOT_VERIFIED` | 403 | Email Not Verified | Please verify your email before logging in |

### 3.2 Tenant

| Error Code | HTTP Status | Title | Contoh Detail |
|------------|-------------|-------|---------------|
| `TENANT_NOT_FOUND` | 404 | Tenant Not Found | Tenant with ID 'ten_xxx' not found |
| `TENANT_SUSPENDED` | 403 | Tenant Suspended | This store has been suspended |
| `TENANT_TRIAL_EXPIRED` | 403 | Trial Expired | Your trial period has ended, please upgrade |
| `TENANT_DOMAIN_TAKEN` | 409 | Domain Already Taken | The domain 'mystore' is already in use |
| `TENANT_LIMIT_REACHED` | 422 | Plan Limit Reached | Your plan does not allow more than 10 products |

### 3.3 Product (Catalog)

| Error Code | HTTP Status | Title | Contoh Detail |
|------------|-------------|-------|---------------|
| `PRODUCT_NOT_FOUND` | 404 | Product Not Found | Product with ID 'prd_xxx' not found |
| `PRODUCT_SKU_EXISTS` | 409 | SKU Already Exists | SKU 'PROD-001' is already used in your store |
| `PRODUCT_SLUG_EXISTS` | 409 | Slug Already Exists | Slug 'macbook-pro' is already taken |
| `PRODUCT_INVALID_STATUS_TRANSITION` | 422 | Invalid Status Transition | Cannot transition from 'published' to 'draft' |
| `PRODUCT_PRICE_INVALID` | 422 | Invalid Price | Price must be greater than 0 |
| `PRODUCT_NAME_REQUIRED` | 422 | Name Required | Product name cannot be empty |
| `CATEGORY_NOT_FOUND` | 404 | Category Not Found | Category with ID 'cat_xxx' not found |

### 3.4 Inventory

| Error Code | HTTP Status | Title | Contoh Detail |
|------------|-------------|-------|---------------|
| `INVENTORY_NOT_FOUND` | 404 | Inventory Not Found | Inventory record for product 'prd_xxx' not found |
| `INVENTORY_INSUFFICIENT_STOCK` | 422 | Insufficient Stock | Only 3 items available, requested 5 |
| `INVENTORY_RESERVATION_EXPIRED` | 422 | Reservation Expired | Stock reservation has expired, please retry checkout |
| `INVENTORY_RESERVATION_NOT_FOUND` | 404 | Reservation Not Found | Reservation 'res_xxx' not found |

### 3.5 Order

| Error Code | HTTP Status | Title | Contoh Detail |
|------------|-------------|-------|---------------|
| `ORDER_NOT_FOUND` | 404 | Order Not Found | Order with ID 'ord_xxx' not found |
| `ORDER_INVALID_STATUS_TRANSITION` | 422 | Invalid Status Transition | Cannot cancel an order that is already shipped |
| `ORDER_EMPTY_CART` | 422 | Cart is Empty | Cannot checkout with empty cart |
| `ORDER_PRODUCT_UNAVAILABLE` | 422 | Product Unavailable | Product 'MacBook Pro' is no longer available |
| `ORDER_ALREADY_PAID` | 409 | Order Already Paid | This order has already been paid |

### 3.6 Payment

| Error Code | HTTP Status | Title | Contoh Detail |
|------------|-------------|-------|---------------|
| `PAYMENT_NOT_FOUND` | 404 | Payment Not Found | Payment with ID 'pay_xxx' not found |
| `PAYMENT_GATEWAY_ERROR` | 502 | Payment Gateway Error | Midtrans is currently unavailable, please try again |
| `PAYMENT_DUPLICATE_WEBHOOK` | 200 | Duplicate Webhook | This webhook has already been processed (idempotent) |
| `PAYMENT_INVALID_SIGNATURE` | 401 | Invalid Webhook Signature | Webhook signature validation failed |
| `PAYMENT_AMOUNT_MISMATCH` | 422 | Amount Mismatch | Payment amount does not match order total |
| `PAYMENT_EXPIRED` | 422 | Payment Expired | Payment window has expired |
| `REFUND_NOT_ELIGIBLE` | 422 | Not Eligible for Refund | Orders in 'completed' status cannot be refunded after 7 days |

### 3.7 Validation

| Error Code | HTTP Status | Title | Contoh Detail |
|------------|-------------|-------|---------------|
| `VALIDATION_FAILED` | 422 | Validation Failed | One or more fields are invalid (include `errors` array) |
| `VALIDATION_REQUIRED_FIELD` | 422 | Required Field Missing | Field 'name' is required |
| `VALIDATION_INVALID_FORMAT` | 422 | Invalid Format | Field 'email' must be a valid email address |
| `VALIDATION_VALUE_TOO_SHORT` | 422 | Value Too Short | Field 'password' must be at least 8 characters |
| `VALIDATION_VALUE_TOO_LONG` | 422 | Value Too Long | Field 'name' cannot exceed 255 characters |
| `VALIDATION_INVALID_ENUM` | 422 | Invalid Enum Value | Status must be one of: draft, published, archived |

### 3.8 Rate Limiting

| Error Code | HTTP Status | Title | Contoh Detail |
|------------|-------------|-------|---------------|
| `RATE_LIMIT_EXCEEDED` | 429 | Rate Limit Exceeded | Too many requests, please wait 60 seconds |

### 3.9 Server Errors

| Error Code | HTTP Status | Title | Contoh Detail |
|------------|-------------|-------|---------------|
| `INTERNAL_SERVER_ERROR` | 500 | Internal Server Error | An unexpected error occurred |
| `SERVICE_UNAVAILABLE` | 503 | Service Unavailable | Service is temporarily unavailable |
| `DATABASE_ERROR` | 500 | Database Error | Failed to connect to database (JANGAN expose detail ke client) |

---

## 4. Validation Error Format (Multiple Fields)

Untuk error validasi dengan multiple field, tambahkan field `errors`:

```json
{
  "type": "https://api.cloudcommerce.com/errors/validation-failed",
  "title": "Validation Failed",
  "status": 422,
  "detail": "One or more fields are invalid",
  "instance": "/api/v1/products",
  "traceId": "abc-123",
  "errors": [
    {
      "field": "name",
      "code": "VALIDATION_REQUIRED_FIELD",
      "message": "Product name is required"
    },
    {
      "field": "price",
      "code": "VALIDATION_INVALID_FORMAT",
      "message": "Price must be a positive integer"
    }
  ]
}
```

---

## 5. Implementasi Go

### Domain Errors

```go
// internal/domain/errors.go

package domain

import "errors"

// Catalog
var (
    ErrProductNotFound          = errors.New("product not found")
    ErrCategoryNotFound         = errors.New("category not found")
    ErrSKUExists                = errors.New("sku already exists")
    ErrSlugExists               = errors.New("slug already exists")
    ErrInvalidStatusTransition  = errors.New("invalid product status transition")
)

// Inventory
var (
    ErrInventoryNotFound        = errors.New("inventory not found")
    ErrInsufficientStock        = errors.New("insufficient stock")
    ErrReservationExpired       = errors.New("reservation expired")
    ErrReservationNotFound      = errors.New("reservation not found")
)

// Order
var (
    ErrOrderNotFound            = errors.New("order not found")
    ErrInvalidOrderTransition   = errors.New("invalid order status transition")
    ErrEmptyCart                = errors.New("cart is empty")
    ErrOrderAlreadyPaid         = errors.New("order already paid")
)
```

### Handler Error Mapping

```go
// internal/transport/errors.go

package transport

import (
    "errors"
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/cloudcommerce/product-service/internal/domain"
)

type ProblemDetail struct {
    Type     string `json:"type"`
    Title    string `json:"title"`
    Status   int    `json:"status"`
    Detail   string `json:"detail"`
    Instance string `json:"instance,omitempty"`
    TraceID  string `json:"traceId,omitempty"`
}

const errorBaseURL = "https://api.cloudcommerce.com/errors"

func handleError(c *gin.Context, err error) {
    var status int
    var code, title string

    switch {
    case errors.Is(err, domain.ErrProductNotFound):
        status, code, title = 404, "product-not-found", "Product Not Found"
    case errors.Is(err, domain.ErrSKUExists):
        status, code, title = 409, "product-sku-exists", "SKU Already Exists"
    case errors.Is(err, domain.ErrSlugExists):
        status, code, title = 409, "product-slug-exists", "Slug Already Exists"
    case errors.Is(err, domain.ErrInvalidStatusTransition):
        status, code, title = 422, "product-invalid-status-transition", "Invalid Status Transition"
    case errors.Is(err, domain.ErrInsufficientStock):
        status, code, title = 422, "inventory-insufficient-stock", "Insufficient Stock"
    default:
        status, code, title = 500, "internal-server-error", "Internal Server Error"
    }

    c.JSON(status, ProblemDetail{
        Type:     errorBaseURL + "/" + code,
        Title:    title,
        Status:   status,
        Detail:   err.Error(),
        Instance: c.Request.URL.Path,
        TraceID:  c.GetString("request_id"),
    })
}
```

---

## 6. Implementasi TypeScript (Frontend)

```typescript
// lib/api-error.ts

export interface ProblemDetail {
    type: string
    title: string
    status: number
    detail: string
    instance?: string
    traceId?: string
    errors?: ValidationError[]
}

export interface ValidationError {
    field: string
    code: string
    message: string
}

export class ApiError extends Error {
    constructor(
        public readonly problem: ProblemDetail,
        public readonly httpStatus: number
    ) {
        super(problem.detail)
    }

    get isNotFound() { return this.httpStatus === 404 }
    get isValidationError() { return this.httpStatus === 422 }
    get isUnauthorized() { return this.httpStatus === 401 }
    get isForbidden() { return this.httpStatus === 403 }
    get isConflict() { return this.httpStatus === 409 }
}

// Penggunaan di TanStack Query
async function fetchProduct(id: string) {
    const res = await fetch(`/api/v1/products/${id}`)
    if (!res.ok) {
        const problem: ProblemDetail = await res.json()
        throw new ApiError(problem, res.status)
    }
    return res.json()
}
```

---

## 7. Rules

```
✅ Semua error HARUS menggunakan format RFC 9457
✅ Error code di type URI HARUS menggunakan kebab-case
✅ Detail message boleh dinamis (menjelaskan konteks spesifik)
✅ Title HARUS statis (tidak berubah untuk error yang sama)
✅ TraceID WAJIB ada untuk membantu debugging
❌ JANGAN expose stack trace atau detail database di response
❌ JANGAN gunakan pesan error Go internal (bungkus dengan domain error)
❌ JANGAN kembalikan 200 untuk error (kecuali duplicate webhook — idempotent)
```

---

## 8. HTTP Status Code Summary

| Status | Digunakan Untuk |
|--------|-----------------|
| 200 | Sukses GET, PUT, PATCH |
| 201 | Sukses POST (resource dibuat) |
| 204 | Sukses DELETE (no content) |
| 400 | Bad request (malformed JSON, dll) |
| 401 | Unauthenticated (tidak ada atau token invalid) |
| 403 | Forbidden (sudah auth tapi tidak punya permission) |
| 404 | Resource tidak ditemukan |
| 409 | Conflict (duplikat data) |
| 422 | Validation failed atau business rule violation |
| 429 | Rate limit exceeded |
| 500 | Internal server error |
| 502 | Upstream service error (payment gateway, dll) |
| 503 | Service unavailable |

---

## 9. Related Documents

- [API Guidelines](../03-architecture/api-guidelines.md)
- [Backend Guidelines](04-backend-guidelines.md)
- [State Machine](../03-architecture/state-machine.md)
