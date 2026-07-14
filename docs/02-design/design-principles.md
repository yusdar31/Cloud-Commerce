# Design Principles

**Project:** CloudCommerce

**Document Version:** 1.0.0

**Status:** Draft

**Author:** UX Team

**Last Updated:** July 2026

---

# 1. Overview

## Purpose

Dokumen ini mendefinisikan prinsip desain yang menjadi pedoman dalam pengembangan antarmuka CloudCommerce.

Tujuan utama dokumen ini adalah memastikan seluruh halaman memiliki pengalaman pengguna yang:

- Konsisten
- Mudah dipelajari
- Cepat digunakan
- Responsif
- Mudah diakses
- Mudah dikembangkan

Design Principles berlaku untuk:

- UI Design
- Wireframe
- Design System
- Frontend Development
- Component Library

---

# 2. Design Philosophy

CloudCommerce mengadopsi filosofi berikut.

> **Simple enough for first-time users, powerful enough for growing businesses.**

UI harus sederhana bagi pengguna baru namun cukup fleksibel untuk merchant yang memiliki ribuan produk.

---

# 3. UX Principles

## 3.1 Clarity First

Setiap halaman harus memiliki tujuan yang jelas.

Pengguna harus mengetahui:

- Sedang berada di mana
- Apa yang bisa dilakukan
- Langkah berikutnya

### Good

Dashboard

↓

Products

↓

Product Detail

↓

Edit Product

### Avoid

Dashboard

↓

Products

↓

Unknown Screen

---

## 3.2 Minimize Cognitive Load

Jangan menampilkan informasi yang tidak diperlukan.

Gunakan:

- whitespace
- grouping
- section
- card
- hierarchy

Daripada:

```
100 informasi dalam satu halaman
```

---

## 3.3 Progressive Disclosure

Fitur kompleks hanya muncul ketika diperlukan.

Contoh

Saat membuat produk.

Yang muncul terlebih dahulu:

- Nama
- Harga
- Stok

SEO

Shipping

Metadata

baru muncul pada Advanced Settings.

---

## 3.4 Recognition Over Recall

Pengguna seharusnya tidak perlu mengingat.

Gunakan:

- breadcrumb
- icon
- tooltip
- autocomplete
- placeholder

---

## 3.5 Consistency

Komponen yang sama harus memiliki perilaku yang sama.

Contoh

Button Save

Selalu berwarna Primary.

Button Delete

Selalu Danger.

---

## 3.6 Accessibility

CloudCommerce mengikuti WCAG 2.2 AA.

Semua fitur harus dapat digunakan menggunakan:

- Keyboard
- Screen Reader
- High Contrast

---

## 3.7 Mobile Friendly

Dashboard tetap nyaman digunakan pada tablet dan mobile.

---

# 4. Visual Design Principles

## Minimal

Tidak menggunakan dekorasi berlebihan.

---

## Clean

Banyak whitespace.

---

## Modern

Menggunakan card layout.

---

## Professional

Inspirasi:

- Stripe
- Shopify
- Linear
- Vercel Dashboard
- GitHub

---

## Neutral

Gunakan warna hanya untuk memberikan makna.

Hijau

Success

Merah

Error

Orange

Warning

Blue

Primary

---

# 5. Layout Principles

## Dashboard

```
Top Navigation

↓

Sidebar

↓

Content

↓

Footer
```

---

## Storefront

```
Navbar

↓

Hero

↓

Categories

↓

Products

↓

Footer
```

---

# 6. Navigation Principles

Semua menu utama maksimal berada pada sidebar level pertama.

Nested menu maksimal:

Level 2.

Tidak menggunakan menu hingga Level 4 atau lebih.

---

# 7. Interaction Principles

Hover

↓

Clickable

Focus

↓

Keyboard Friendly

Loading

↓

Skeleton

Success

↓

Toast

Delete

↓

Confirmation Dialog

---

# 8. Feedback Principles

Semua aksi harus memberikan feedback.

Contoh

Product Saved

✓ Product berhasil disimpan.

Order Deleted

✓ Order berhasil dihapus.

Payment Failed

❌ Pembayaran gagal.

---

# 9. Form Principles

Gunakan:

Label

↓

Input

↓

Helper Text

↓

Validation

↓

Error Message

Bukan:

Input

↓

Error

↓

Label

---

# 10. Table Principles

Semua tabel harus mendukung:

- Search
- Filter
- Sort
- Pagination
- Bulk Action
- Export (future)

---

# 11. Empty State

Semua halaman wajib memiliki Empty State.

Contoh

Products

```
No products yet.

Create your first product.

[ Add Product ]
```

---

# 12. Loading State

Menggunakan Skeleton Loader.

Spinner hanya digunakan untuk:

- Button
- Dialog
- API singkat

---

# 13. Error State

Setiap halaman harus memiliki:

Retry

↓

Back

↓

Help

---

# 14. Notification Principles

Notification dibagi menjadi empat.

Information

Success

Warning

Error

Toast tidak boleh lebih dari:

5 detik.

---

# 15. Color Principles

Primary

Brand

Secondary

Neutral

Success

Warning

Danger

Info

Seluruh warna akan didefinisikan pada Design System.

---

# 16. Typography Principles

Menggunakan maksimal:

- Heading
- Subheading
- Body
- Caption

Tidak lebih.

---

# 17. Icon Principles

Menggunakan Lucide Icons.

Icon harus selalu disertai label jika aksi penting.

Contoh

🗑 Delete

✓ Bukan hanya icon.

---

# 18. Component Principles

Semua komponen harus:

Reusable

Composable

Accessible

Responsive

Dark Mode Ready

---

# 19. Performance Principles

Target:

First Load

<2 detik

Dashboard

<1 detik

Search

<300ms

Interaction

<100ms

---

# 20. Security UX

Password

Show / Hide

Session Expired

Auto Redirect

CSRF Error

Friendly Message

Permission Denied

403 Screen

---

# 21. Responsive Breakpoints

Desktop

≥1440px

Laptop

1024px

Tablet

768px

Mobile

390px

---

# 22. Future Design Goals

Dark Mode

White Label

RTL Language

Localization

Offline Mode

PWA

---

# 23. References

- Material Design 3
- Apple Human Interface Guidelines
- WCAG 2.2
- Stripe Dashboard
- Shopify Polaris
- Vercel Design
- GitHub Primer