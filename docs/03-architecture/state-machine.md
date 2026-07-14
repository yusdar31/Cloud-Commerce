# State Machine

**Project:** CloudCommerce

**Version:** 1.0.0

**Status:** Approved

**Owner:** Engineering Team

**Last Updated:** July 2026

---

## 1. Purpose

Dokumen ini mendefinisikan state machine (mesin status) untuk semua entitas bisnis utama CloudCommerce.

AI dan developer wajib menggunakan state diagram ini sebagai referensi saat:

- Membuat atau mengubah status entitas
- Menulis validasi transisi status
- Membuat migration database
- Membuat business logic di layer application

---

## 2. Product Status

### Status Values

| Status | Nilai | Deskripsi |
|--------|-------|-----------|
| Draft | `draft` | Produk baru dibuat, belum terlihat oleh buyer |
| Published | `published` | Produk aktif, terlihat di storefront |
| Archived | `archived` | Produk dinonaktifkan, tidak terlihat buyer |

### Transisi yang Diizinkan

```
[draft] ──────────────────────────────────► [published]
   │                                              │
   │                                              │
   └──────────────────────────────────► [archived]│
                                                  │
[published] ──────────────────────────► [archived]│
     │                                            │
     └──────────────────────────────── (tidak bisa kembali ke draft)
```

### Tabel Transisi

| From | To | Command | Kondisi |
|------|----|---------|---------|
| `draft` | `published` | `PublishProduct` | Nama, harga, stok wajib ada |
| `draft` | `archived` | `ArchiveProduct` | - |
| `published` | `archived` | `ArchiveProduct` | - |
| `archived` | `published` | `PublishProduct` | - |
| `published` | `draft` | ❌ TIDAK BISA | - |
| `archived` | `draft` | ❌ TIDAK BISA | - |

### Implementasi Go

```go
// domain/entity.go
type ProductStatus string

const (
    StatusDraft     ProductStatus = "draft"
    StatusPublished ProductStatus = "published"
    StatusArchived  ProductStatus = "archived"
)

var ErrInvalidStatusTransition = errors.New("invalid product status transition")

func (p *Product) Publish() error {
    if p.Status != StatusDraft && p.Status != StatusArchived {
        return ErrInvalidStatusTransition
    }
    p.Status = StatusPublished
    return nil
}

func (p *Product) Archive() error {
    if p.Status == StatusArchived {
        return ErrInvalidStatusTransition
    }
    p.Status = StatusArchived
    return nil
}
```

---

## 3. Order Status

### Status Values

| Status | Nilai | Deskripsi |
|--------|-------|-----------|
| Pending | `pending` | Order baru dibuat, menunggu pembayaran |
| Awaiting Payment | `awaiting_payment` | Invoice dikirim ke payment gateway |
| Paid | `paid` | Pembayaran berhasil dikonfirmasi |
| Processing | `processing` | Seller sedang memproses order |
| Shipped | `shipped` | Barang sudah dikirim |
| Completed | `completed` | Barang sudah diterima buyer |
| Cancelled | `cancelled` | Order dibatalkan |
| Refunded | `refunded` | Order di-refund |

### Diagram Transisi

```
[pending]
    │
    ▼
[awaiting_payment] ──(payment timeout)──► [cancelled]
    │
    ├──(payment success)──► [paid]
    │                          │
    │                          ▼
    │                     [processing]
    │                          │
    │                          ▼
    │                      [shipped]
    │                          │
    │                          ▼
    │                     [completed]
    │
    └──(payment failed)──► [cancelled]

[paid] ──────────────────────────────────► [refunded]
[processing] ────────────────────────────► [refunded]
[shipped] ───────────────────────────────► [refunded]
```

### Tabel Transisi

| From | To | Trigger | Actor |
|------|----|---------|-------|
| `pending` | `awaiting_payment` | Order dibuat → Invoice generated | System |
| `awaiting_payment` | `paid` | Webhook PaymentSucceeded | Payment Service |
| `awaiting_payment` | `cancelled` | Payment timeout / PaymentFailed | System |
| `paid` | `processing` | Seller konfirmasi order | Seller |
| `processing` | `shipped` | Seller input resi pengiriman | Seller |
| `shipped` | `completed` | Buyer konfirmasi penerimaan / Auto setelah 7 hari | Buyer / System |
| `paid` | `refunded` | Seller approve refund | Seller |
| `processing` | `refunded` | Seller approve refund | Seller |
| `shipped` | `refunded` | Seller approve refund | Seller |
| ❌ | `pending` | Tidak ada transisi ke pending setelah dibuat | - |

### Implementasi Go

```go
// domain/entity.go
type OrderStatus string

const (
    OrderPending         OrderStatus = "pending"
    OrderAwaitingPayment OrderStatus = "awaiting_payment"
    OrderPaid            OrderStatus = "paid"
    OrderProcessing      OrderStatus = "processing"
    OrderShipped         OrderStatus = "shipped"
    OrderCompleted       OrderStatus = "completed"
    OrderCancelled       OrderStatus = "cancelled"
    OrderRefunded        OrderStatus = "refunded"
)

var ErrInvalidOrderTransition = errors.New("invalid order status transition")

// ValidTransitions mendefinisikan transisi yang diizinkan
var validOrderTransitions = map[OrderStatus][]OrderStatus{
    OrderPending:         {OrderAwaitingPayment},
    OrderAwaitingPayment: {OrderPaid, OrderCancelled},
    OrderPaid:            {OrderProcessing, OrderRefunded},
    OrderProcessing:      {OrderShipped, OrderRefunded},
    OrderShipped:         {OrderCompleted, OrderRefunded},
    OrderCompleted:       {},
    OrderCancelled:       {},
    OrderRefunded:        {},
}

func (o *Order) CanTransitionTo(next OrderStatus) bool {
    allowed := validOrderTransitions[o.Status]
    for _, s := range allowed {
        if s == next {
            return true
        }
    }
    return false
}

func (o *Order) TransitionTo(next OrderStatus) error {
    if !o.CanTransitionTo(next) {
        return fmt.Errorf("%w: %s → %s", ErrInvalidOrderTransition, o.Status, next)
    }
    o.Status = next
    return nil
}
```

---

## 4. Payment Status

### Status Values

| Status | Nilai | Deskripsi |
|--------|-------|-----------|
| Pending | `pending` | Invoice dibuat, menunggu pembayaran |
| Processing | `processing` | Pembayaran sedang diproses gateway |
| Success | `success` | Pembayaran berhasil |
| Failed | `failed` | Pembayaran gagal |
| Expired | `expired` | Waktu pembayaran habis |
| Refunded | `refunded` | Pembayaran dikembalikan |
| Partial Refund | `partial_refund` | Sebagian dikembalikan |

### Diagram Transisi

```
[pending]
    │
    ▼
[processing]
    │
    ├──(webhook success)──► [success]
    │                           │
    │                           └──► [refunded]
    │                           │
    │                           └──► [partial_refund]
    │
    ├──(webhook failed)──► [failed]
    │
    └──(timeout)──► [expired]
```

### Tabel Transisi

| From | To | Trigger |
|------|----|---------|
| `pending` | `processing` | Payment gateway dipanggil |
| `processing` | `success` | Webhook payment success |
| `processing` | `failed` | Webhook payment failed |
| `processing` | `expired` | Timeout (biasanya 24 jam) |
| `success` | `refunded` | Refund diproses & berhasil |
| `success` | `partial_refund` | Partial refund diproses |
| `failed` | `pending` | Buyer retry pembayaran |
| `expired` | `pending` | Buyer retry pembayaran (buat invoice baru) |

---

## 5. Inventory Reservation Status

### Status Values

| Status | Nilai | Deskripsi |
|--------|-------|-----------|
| Available | `available` | Stok tersedia untuk dipesan |
| Reserved | `reserved` | Stok dikunci sementara saat checkout |
| Committed | `committed` | Stok dikurangi permanen setelah payment success |
| Released | `released` | Reservasi dibatalkan, stok dikembalikan |

### Diagram Transisi

```
[available]
    │
    ▼ (checkout: StockReserved event)
[reserved]
    │
    ├──(PaymentSucceeded)──► [committed]  (stok dikurangi permanen)
    │
    └──(PaymentFailed / Timeout)──► [released]  (stok dikembalikan ke available)
```

### Rules

```
reserved_quantity = jumlah stok yang di-lock
available_quantity = total_stock - reserved_quantity

// Saat Reserve:
available_quantity -= quantity
reserved_quantity  += quantity

// Saat Commit (payment success):
reserved_quantity -= quantity
// (stok sudah berkurang dari total)

// Saat Release (payment failed):
reserved_quantity  -= quantity
available_quantity += quantity
```

### Timeout Reservasi

Reservasi stok **otomatis expired** setelah **30 menit** jika pembayaran tidak selesai.

---

## 6. Subscription Status (Tenant)

### Status Values

| Status | Nilai | Deskripsi |
|--------|-------|-----------|
| Trial | `trial` | Periode trial gratis (14 hari) |
| Active | `active` | Berlangganan aktif |
| Past Due | `past_due` | Billing gagal, grace period |
| Suspended | `suspended` | Akses ditangguhkan |
| Cancelled | `cancelled` | Berlangganan dibatalkan |
| Expired | `expired` | Trial/langganan habis masa berlaku |

### Diagram Transisi

```
[trial]
    │
    ├──(upgrade)──► [active]
    │
    └──(14 hari habis tanpa upgrade)──► [expired]

[active]
    │
    ├──(billing gagal)──► [past_due]
    │                          │
    │                          ├──(bayar dalam 7 hari)──► [active]
    │                          │
    │                          └──(7 hari tidak bayar)──► [suspended]
    │
    └──(cancel)──► [cancelled]

[suspended]
    │
    ├──(bayar)──► [active]
    │
    └──(30 hari)──► [expired]
```

---

## 7. Notification Status

### Status Values

| Status | Nilai | Deskripsi |
|--------|-------|-----------|
| Queued | `queued` | Notification menunggu untuk dikirim |
| Sending | `sending` | Sedang dalam proses pengiriman |
| Sent | `sent` | Berhasil dikirim ke provider |
| Failed | `failed` | Pengiriman gagal |

### Retry Policy

```
failed → retry setelah: 1 menit, 5 menit, 30 menit (max 3 retry)
failed setelah 3 retry → alert ke monitoring
```

---

## 8. Related Documents

- [Event Storming](event-storming.md)
- [Domain Model](domain-model.md)
- [Service Boundaries](service-boundaries.md)
- [API Guidelines](api-guidelines.md)
