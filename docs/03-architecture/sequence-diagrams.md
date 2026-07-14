# Sequence Diagrams

**Project:** CloudCommerce

**Version:** 1.0.0

**Status:** Approved

**Owner:** Solution Architecture Team

**Last Updated:** July 2026

---

## 1. Purpose

Dokumen ini mendefinisikan sequence diagram untuk alur bisnis utama CloudCommerce.

Gunakan sebagai referensi saat:
- Memahami alur komunikasi antar service
- Debugging error lintas service
- Membangun fitur baru yang melibatkan multiple service

---

## 2. Flow 1: Seller Register & Login

```
Browser          API Gateway      User Service     Store Service    NATS
   │                  │                │                 │            │
   │ POST /auth/register              │                 │            │
   │─────────────────►│               │                 │            │
   │                  │ POST /api/v1/auth/register      │            │
   │                  │──────────────►│                 │            │
   │                  │               │ Hash password   │            │
   │                  │               │ Create user     │            │
   │                  │               │ Publish ──────────────────► UserRegistered
   │                  │               │                 │            │
   │                  │               │                 │◄─────────── UserRegistered
   │                  │               │                 │ Create tenant
   │                  │               │                 │ Set trial subscription
   │                  │               │                 │ Publish ──► TenantCreated
   │                  │               │                 │            │
   │                  │◄──────────────│                 │            │
   │◄─────────────────│               │                 │            │
   │ 201 {userId, tenantId}           │                 │            │

   │ POST /auth/login │               │                 │
   │─────────────────►│               │                 │
   │                  │──────────────►│                 │
   │                  │               │ Verify password │
   │                  │               │ Generate JWT    │
   │                  │◄──────────────│                 │
   │◄─────────────────│               │                 │
   │ 200 {accessToken, refreshToken}  │                 │
```

---

## 3. Flow 2: Seller Create & Publish Product

```
Browser          API Gateway      Product Service   NATS           Inventory Service
   │                  │                │              │                  │
   │ POST /api/v1/products             │              │                  │
   │ Authorization: Bearer {JWT}       │              │                  │
   │─────────────────►│               │              │                  │
   │                  │ Validate JWT   │              │                  │
   │                  │ Extract tenant_id             │                  │
   │                  │──────────────►│              │                  │
   │                  │               │ Validate req │                  │
   │                  │               │ Check SKU unique               │
   │                  │               │ Create product (status: draft)  │
   │                  │               │ Publish ─────►│ ProductCreated  │
   │                  │               │              │                  │
   │                  │               │              │◄─── ProductCreated
   │                  │               │              │    Initialize inventory (qty=0)
   │                  │◄──────────────│              │                  │
   │◄─────────────────│               │              │                  │
   │ 201 {product}    │               │              │                  │

   │ POST /api/v1/products/{id}/publish               │                  │
   │─────────────────►│               │              │                  │
   │                  │──────────────►│              │                  │
   │                  │               │ Validate status transition      │
   │                  │               │ draft → published               │
   │                  │               │ Publish ─────►│ ProductPublished │
   │                  │◄──────────────│              │                  │
   │◄─────────────────│               │              │                  │
   │ 200 {product: status=published}  │              │                  │
```

---

## 4. Flow 3: Buyer Checkout (Happy Path)

```
Browser     API Gateway   Order Service   Inventory Service   Payment Service   NATS
   │             │              │                │                  │              │
   │ POST /api/v1/orders        │                │                  │              │
   │(cart items) │              │                │                  │              │
   │────────────►│              │                │                  │              │
   │             │─────────────►│                │                  │              │
   │             │              │ Validate cart  │                  │              │
   │             │              │ Calculate total│                  │              │
   │             │              │                │                  │              │
   │             │              │ POST /inventory/reserve           │              │
   │             │              │───────────────►│                  │              │
   │             │              │                │ Check stock      │              │
   │             │              │                │ Reserve stock    │              │
   │             │              │◄───────────────│                  │              │
   │             │              │ reservationId  │                  │              │
   │             │              │                │                  │              │
   │             │              │ Create order (status: pending)    │              │
   │             │              │ Publish ──────────────────────────────────────► OrderCreated
   │             │              │                │                  │              │
   │             │◄─────────────│                │                  │              │
   │◄────────────│              │                │                  │              │
   │ 201 {orderId, status: pending}              │                  │              │
   │             │              │                │         OrderCreated ◄──────────│
   │             │              │                │                  │ Generate invoice
   │             │              │                │                  │ Call Midtrans API
   │             │              │                │                  │ Save payment URL
   │             │              │                │                  │ Publish ────► PaymentRequested
   │             │              │                │                  │              │
   │ GET /api/v1/payments/{orderId}              │                  │              │
   │────────────►│              │                │                  │              │
   │             │──────────────────────────────►│                  │              │
   │◄────────────│              │                │                  │              │
   │ 200 {paymentUrl: "https://midtrans..."}     │                  │              │
   │             │              │                │                  │              │
   │ [User redirects to Midtrans payment page]  │                  │              │
   │ [User completes payment]   │                │                  │              │
```

---

## 5. Flow 4: Payment Webhook (Payment Success)

```
Midtrans    Payment Service    NATS     Order Service   Inventory Service   Notification Service
   │              │              │            │                │                    │
   │ POST /payments/webhook/midtrans          │                │                    │
   │─────────────►│              │            │                │                    │
   │              │ Validate signature         │                │                    │
   │              │ Check idempotency          │                │                    │
   │              │ (transaction_id exists?)   │                │                    │
   │              │ Update payment → success  │                │                    │
   │              │ Publish ─────►│ PaymentSucceeded           │                    │
   │◄─────────────│              │            │                │                    │
   │ 200 OK       │              │            │                │                    │
   │              │              │            │                │                    │
   │              │              │ PaymentSucceeded ───────────►│                   │
   │              │              │            │                │ Commit reservation │
   │              │              │            │                │ (reduce stock perm)│
   │              │              │            │                │                    │
   │              │              │ PaymentSucceeded ──────────►│                   │
   │              │              │            │ Update order   │                    │
   │              │              │            │ pending → paid │                    │
   │              │              │            │ Publish ────── ──────────────────► OrderConfirmed (via NATS)
   │              │              │            │                │                    │
   │              │              │ PaymentSucceeded ──────────────────────────────►│
   │              │              │            │                │  Send receipt email│
   │              │              │            │                │  Send seller alert │
```

---

## 6. Flow 5: Payment Failure & Stock Release

```
Midtrans    Payment Service    NATS     Order Service   Inventory Service
   │              │              │            │                │
   │ POST /webhook (payment failed)           │                │
   │─────────────►│              │            │                │
   │              │ Validate     │            │                │
   │              │ Update payment → failed   │                │
   │              │ Publish ─────►│ PaymentFailed             │
   │◄─────────────│              │            │                │
   │ 200 OK       │              │            │                │
   │              │              │ PaymentFailed ─────────────►│
   │              │              │            │                │ Release reservation
   │              │              │            │                │ (return stock)
   │              │              │ PaymentFailed ─────────────►│
   │              │              │            │ Update order   │
   │              │              │            │ → cancelled    │
```

---

## 7. Flow 6: Token Refresh

```
Browser          API Gateway      User Service    Redis
   │                  │                │              │
   │ POST /auth/refresh               │              │
   │ {refreshToken}   │               │              │
   │─────────────────►│               │              │
   │                  │──────────────►│              │
   │                  │               │ Validate refresh token
   │                  │               │ Check blacklist ──────►│
   │                  │               │                        │ token_exists?
   │                  │               │◄───────────────────────│
   │                  │               │ Generate new access token
   │                  │               │ Rotate refresh token   │
   │                  │               │ Blacklist old token ──►│
   │                  │◄──────────────│              │
   │◄─────────────────│               │              │
   │ 200 {accessToken, refreshToken}  │              │
```

---

## 8. Flow 7: Storefront Load Products (dengan Cache)

```
Browser    API Gateway   Product Service    Redis (Cache)    PostgreSQL
   │            │              │                  │               │
   │ GET /api/v1/products      │                  │               │
   │ (tenant from JWT or slug) │                  │               │
   │───────────►│              │                  │               │
   │            │─────────────►│                  │               │
   │            │              │ Generate cache key               │
   │            │              │ GET products:{tenant_id} ───────►│
   │            │              │                  │               │
   │            │              │◄──── HIT ────────│               │
   │            │◄─────────────│                  │               │
   │◄───────────│              │                  │               │
   │ 200 {products} (from cache, <50ms)           │               │

   [Cache MISS scenario]
   │ GET /api/v1/products      │                  │               │
   │───────────►│              │                  │               │
   │            │─────────────►│                  │               │
   │            │              │ GET cache ───────►│               │
   │            │              │◄──── MISS ────────│               │
   │            │              │ SELECT * FROM products WHERE tenant_id = $1
   │            │              │────────────────────────────────►│
   │            │              │◄────────────────────────────────│
   │            │              │ SET cache (TTL: 60s) ──────────►│
   │            │◄─────────────│                  │               │
   │◄───────────│              │                  │               │
   │ 200 {products}            │                  │               │
```

---

## 9. Error Flow: Product Not Found

```
Browser    API Gateway   Product Service
   │            │              │
   │ GET /api/v1/products/prd_xxx_invalid
   │───────────►│              │
   │            │─────────────►│
   │            │              │ SELECT * WHERE id = 'prd_xxx_invalid' AND tenant_id = ...
   │            │              │ → No rows found
   │            │◄─────────────│ 404
   │◄───────────│              │
   │ 404 {
   │   "type": "https://api.cloudcommerce.com/errors/product-not-found",
   │   "title": "Product Not Found",
   │   "status": 404,
   │   "detail": "Product with ID 'prd_xxx_invalid' not found",
   │   "traceId": "req_abc123"
   │ }
```

---

## 10. Related Documents

- [Event Storming](event-storming.md)
- [State Machine](state-machine.md)
- [API Guidelines](api-guidelines.md)
- [Service Boundaries](service-boundaries.md)
