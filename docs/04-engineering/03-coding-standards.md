# Coding Standards

**Project:** CloudCommerce

**Version:** 1.0.0

**Status:** Approved

**Owner:** Engineering Team

**Last Updated:** July 2026

---

# 1. Purpose

Dokumen ini mendefinisikan standar penulisan kode yang digunakan dalam CloudCommerce.

Tujuan:

- Menjaga konsistensi kode
- Mempermudah maintenance
- Mempermudah code review
- Mengurangi technical debt
- Mendukung AI-assisted development

---

# 2. Engineering Principles

Seluruh kode mengikuti prinsip berikut:

- Readability over Cleverness
- Simplicity First
- SOLID Principles
- DRY (Don't Repeat Yourself)
- KISS (Keep It Simple)
- YAGNI (You Aren't Gonna Need It)

Kode harus mudah dibaca oleh developer lain.

---

# 3. Clean Architecture

Setiap microservice menggunakan struktur berikut:

```
internal/

domain/

application/

infrastructure/

transport/
```

## Domain

Berisi:

- Entity
- Value Object
- Business Rules

Tidak boleh bergantung pada framework.

---

## Application

Berisi:

- Use Case
- Business Flow

Menggunakan Domain.

---

## Infrastructure

Berisi implementasi:

- PostgreSQL
- Redis
- NATS
- MinIO

---

## Transport

Berisi:

- REST Handler
- Middleware
- Request / Response

---

# 4. Naming Convention

## Folder

Gunakan:

```
snake_case
```

Contoh:

```
order_service
payment_service
```

---

## File

Gunakan:

```
snake_case.go
```

Contoh:

```
order_handler.go
product_repository.go
create_order.go
```

---

## Package

Gunakan:

```
package order
```

Hindari package seperti:

```
package utils
package common
```

Jika memungkinkan, beri nama sesuai domain.

---

## Variable

Gunakan nama yang jelas.

Baik:

```go
totalPrice
customerID
```

Hindari:

```go
x
tmp
data
```

---

## Function

Nama fungsi harus menjelaskan aksi.

Contoh:

```
CreateOrder()

ReserveStock()

SendInvoice()
```

---

# 5. Error Handling

Selalu kembalikan error.

Contoh:

```go
order, err := repository.FindByID(id)

if err != nil {
    return err
}
```

Jangan mengabaikan error.

---

# 6. Logging

Gunakan structured logging.

Contoh informasi yang dicatat:

- request_id
- tenant_id
- user_id
- order_id
- duration

Jangan pernah mencatat:

- Password
- Token
- Informasi sensitif

---

# 7. Configuration

Semua konfigurasi berasal dari Environment Variable.

Jangan hardcode:

- URL database
- Secret
- API Key

---

# 8. API Standards

Semua API menggunakan:

- JSON
- REST
- Versioning (`/api/v1`)

Gunakan HTTP Status Code yang sesuai.

---

# 9. Database Standards

- Gunakan migration
- Jangan mengubah database secara manual
- Gunakan transaksi untuk operasi penting
- Hindari query N+1

---

# 10. Comments

Komentar digunakan untuk menjelaskan *mengapa*, bukan *apa*.

Kurang baik:

```go
// Increment stock
stock++
```

Lebih baik:

```go
// Reservasi stok dilakukan sebelum pembayaran
// untuk mencegah overselling.
```

---

# 11. Testing

Setiap fitur baru sebaiknya memiliki:

- Unit Test
- Integration Test (jika perlu)

---

# 12. Git Standards

- Satu branch untuk satu fitur
- Pull Request sebelum merge
- Gunakan Conventional Commit

---

# 13. AI-Assisted Development

AI boleh digunakan untuk:

- Membuat boilerplate
- Menulis test
- Membuat dokumentasi
- Refactoring

Namun, seluruh kode tetap harus dipahami dan direview sebelum digunakan.

---

# 14. Definition of Done

Sebuah fitur dianggap selesai jika:

- Kode berjalan dengan baik
- Tidak ada error lint
- Dokumentasi diperbarui jika diperlukan
- Test yang relevan berhasil dijalankan