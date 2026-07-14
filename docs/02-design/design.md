---
name: Forest & Stone
colors:
  surface: '#f9faf6'
  surface-dim: '#d9dad7'
  surface-bright: '#f9faf6'
  surface-container-lowest: '#ffffff'
  surface-container-low: '#f3f4f1'
  surface-container: '#edeeeb'
  surface-container-high: '#e7e8e5'
  surface-container-highest: '#e2e3e0'
  on-surface: '#1a1c1a'
  on-surface-variant: '#414944'
  inverse-surface: '#2e312f'
  inverse-on-surface: '#f0f1ee'
  outline: '#717974'
  outline-variant: '#c0c9c2'
  surface-tint: '#396754'
  primary: '#013626'
  on-primary: '#ffffff'
  primary-container: '#1e4d3b'
  on-primary-container: '#8cbda6'
  inverse-primary: '#a0d1b9'
  secondary: '#2d694d'
  on-secondary: '#ffffff'
  secondary-container: '#aeedca'
  on-secondary-container: '#326e51'
  tertiary: '#4c2120'
  on-tertiary: '#ffffff'
  tertiary-container: '#673735'
  on-tertiary-container: '#e3a29e'
  error: '#ba1a1a'
  on-error: '#ffffff'
  error-container: '#ffdad6'
  on-error-container: '#93000a'
  primary-fixed: '#bbeed5'
  primary-fixed-dim: '#a0d1b9'
  on-primary-fixed: '#002115'
  on-primary-fixed-variant: '#204f3d'
  secondary-fixed: '#b1f0cd'
  secondary-fixed-dim: '#96d4b2'
  on-secondary-fixed: '#002113'
  on-secondary-fixed-variant: '#105137'
  tertiary-fixed: '#ffdad7'
  tertiary-fixed-dim: '#fab5b1'
  on-tertiary-fixed: '#350f0f'
  on-tertiary-fixed-variant: '#693937'
  background: '#f9faf6'
  on-background: '#1a1c1a'
  surface-variant: '#e2e3e0'
  brand-800: '#173C2C'
  brand-100: '#D7ECE0'
  brand-50: '#F0F7F3'
  stone-900: '#1C1917'
  stone-700: '#44403C'
  stone-500: '#78716C'
  stone-300: '#D6D3D1'
  stone-100: '#F5F5F4'
  success: '#0D9488'
  warning: '#D97706'
  danger: '#DC2626'
  info: '#0284C7'
  dark-bg: '#141311'
typography:
  text-display:
    fontFamily: Inter
    fontSize: 36px
    fontWeight: '700'
    lineHeight: '1.2'
    letterSpacing: -0.02em
  text-h1:
    fontFamily: Inter
    fontSize: 28px
    fontWeight: '700'
    lineHeight: '1.2'
    letterSpacing: -0.02em
  text-h2:
    fontFamily: Inter
    fontSize: 22px
    fontWeight: '600'
    lineHeight: '1.2'
  text-h3:
    fontFamily: Inter
    fontSize: 18px
    fontWeight: '600'
    lineHeight: '1.2'
  text-body:
    fontFamily: Inter
    fontSize: 15px
    fontWeight: '400'
    lineHeight: '1.5'
  text-body-sm:
    fontFamily: Inter
    fontSize: 13px
    fontWeight: '400'
    lineHeight: '1.5'
  text-caption:
    fontFamily: Inter
    fontSize: 12px
    fontWeight: '500'
    lineHeight: '1.5'
  mono-data:
    fontFamily: jetbrainsMono
    fontSize: 13px
    fontWeight: '400'
    lineHeight: '1.5'
rounded:
  sm: 0.125rem
  DEFAULT: 0.25rem
  md: 0.375rem
  lg: 0.5rem
  xl: 0.75rem
  full: 9999px
spacing:
  base: 4px
  xs: 4px
  sm: 8px
  md: 16px
  lg: 24px
  xl: 32px
  2xl: 48px
  3xl: 64px
---

# Design System: CloudCommerce

**Project:** CloudCommerce
**Document Version:** 1.0.0
**Status:** Draft
**Author:** Engineering Team
**Last Updated:** July 2026

---

# 1. Overview

## Purpose

Dokumen ini mendefinisikan bahasa visual (Design Language) dan sistem komponen UI untuk CloudCommerce, mencakup tiga permukaan utama: **Public Website**, **Seller Dashboard**, **Buyer Storefront**, dan **Platform Admin** (lihat `information-architecture.md`).

Design System ini menjadi rujukan untuk:

- Konsistensi visual lintas produk
- Kecepatan development frontend (Next.js)
- Aksesibilitas dan skalabilitas komponen
- Referensi handoff desain ke engineering

## Design Direction

**"Minimalist, but not empty."** Arahan visual CloudCommerce adalah *minimalis modern*: whitespace lega, tipografi tegas, warna netral dengan satu accent color yang jelas, dan dekorasi seminimal mungkin — tapi tetap terasa hidup lewat micro-interaction, hirarki yang kuat, dan sedikit personality di ilustrasi/empty state. Ini bukan gaya "kosong generik SaaS template", melainkan bersih dengan sentuhan karakter.

Referensi rasa (bukan untuk ditiru identik): Linear, Vercel Dashboard, Stripe Dashboard, Shopify Polaris — clean, grid-disiplin, accent color yang berani tapi tidak ramai.

---

# 2. Design Principles

### 1. Clarity Over Decoration
Setiap elemen visual harus punya fungsi. Tidak ada shadow, gradient, atau border berlebihan hanya untuk estetika.

### 2. Consistent Density per Context
Storefront (Buyer) = lega, product-first, emosional.
Dashboard (Seller/Admin) = padat-informasi, efisien, data-first.
Satu design token system, dua density scale.

### 3. One Accent, Many Neutrals
Sistem warna didominasi neutral (grayscale), dengan satu accent color sebagai penanda aksi utama dan brand identity. Warna semantik (success/warning/error) dipakai secara fungsional saja.

### 4. Predictable Interaction
Pola interaksi (hover, focus, loading, empty, error) harus konsisten di semua modul — sejalan dengan Empty/Error/Loading State di `information-architecture.md` §16–18.

### 5. Accessible by Default
Kontras warna, ukuran target sentuh, dan navigasi keyboard mengikuti WCAG 2.1 AA minimum.

### 6. Mobile First, Role Aware
Layout menyesuaikan role (Guest/Buyer/Seller/Admin) sesuai Role Based Navigation di IA §2.2, dan selalu dirancang dari breakpoint mobile ke atas.

---

# 3. Color System

## 3.1 Brand & Neutral Palette — "Forest & Stone"

Accent **Forest Green** dipasangkan dengan neutral **Stone** (warm gray, bukan cool gray) supaya kesan hijau tidak terasa "dingin/klinis" seperti hijau medis, melainkan hangat dan grounded — cocok untuk narasi "kemandirian toko, pertumbuhan, kepercayaan" yang jadi value proposition CloudCommerce (`product-vision.md` §3).

| Token | Hex | Penggunaan |
|---|---|---|
| `color-brand-800` | `#173C2C` | Text-on-tint, active/pressed state |
| `color-brand-700` | `#1E4D3B` | Primary action, link, active state |
| `color-brand-600` | `#2F6B4F` | Hover state primary button |
| `color-brand-100` | `#D7ECE0` | Background highlight/badge lembut |
| `color-brand-50` | `#F0F7F3` | Surface highlight sangat subtle (banner, selected row) |
| `color-neutral-900` | `#1C1917` | Heading text, primary text |
| `color-neutral-700` | `#44403C` | Body text |
| `color-neutral-500` | `#78716C` | Secondary text, placeholder |
| `color-neutral-300` | `#D6D3D1` | Border, divider |
| `color-neutral-100` | `#F5F5F4` | Surface subtle / page background |
| `color-neutral-0` | `#FFFFFF` | Surface / card background |

## 3.2 Semantic Colors

| Token | Hex | Penggunaan |
|---|---|---|
| `color-success` | `#0D9488` (Teal) | Payment success, published, order completed |
| `color-warning` | `#D97706` | Pending payment, low stock |
| `color-danger` | `#DC2626` | Error, out of stock, destructive action |
| `color-info` | `#0284C7` | Informational banner, notice |

## 3.3 Role Accent (opsional, subtle)

Untuk membedakan konteks tanpa mengubah brand utama, gunakan aksen tipis (border-left 2px atau badge kecil):

| Role | Accent |
|---|---|
| Buyer / Storefront | `color-brand-700` (forest green untuk CTA & harga) |
| Seller Dashboard | Neutral (stone), brand hijau hanya muncul di CTA & status penting |
| Platform Admin | `color-neutral-900` (ink/stone tegas, minim hijau — kesan teknikal/serius) |

## 3.4 Dark Mode (Dashboard/Admin — opsional MVP+)

Dashboard operator (Yusdar persona) cocok punya dark mode karena penggunaan lama untuk monitoring. Storefront tetap light-only untuk MVP.

| Token | Hex |
|---|---|
| `color-dark-bg` | `#141311` (Stone, hampir hitam) |
| `color-dark-surface` | `#1C1917` (Stone 900) |
| `color-dark-border` | `#292524` (Stone 800) |
| `color-dark-text` | `#FAFAF9` |
| `color-dark-brand` | `#4ADE80` (hijau dicerahkan agar tetap kontras di atas dark surface) |

---

# 4. Typography

## 4.1 Font Family

- **Primary (UI & Body):** `Inter` — netral, sangat legible di ukuran kecil, standar de facto SaaS modern.
- **Monospace (opsional, untuk Admin/Log/ID/Tenant ID):** `JetBrains Mono` atau `IBM Plex Mono`.

```css
--font-sans: "Inter", -apple-system, BlinkMacSystemFont, sans-serif;
--font-mono: "JetBrains Mono", ui-monospace, monospace;
```

## 4.2 Type Scale (mobile-first, rem-based)

| Token | Size | Weight | Penggunaan |
|---|---|---|---|
| `text-display` | 36px / 2.25rem | 700 | Landing hero |
| `text-h1` | 28px | 700 | Page title |
| `text-h2` | 22px | 600 | Section title |
| `text-h3` | 18px | 600 | Card title / subsection |
| `text-body` | 15px | 400 | Body text default |
| `text-body-sm` | 13px | 400 | Secondary text, meta info |
| `text-caption` | 12px | 500 | Label, badge, timestamp |

Line-height: 1.5 untuk body, 1.2 untuk heading. Letter-spacing default (0), heading besar boleh `-0.02em`.

---

# 5. Spacing & Layout

## 5.1 Spacing Scale (4px base unit — konsisten dengan Tailwind default)

`4 / 8 / 12 / 16 / 20 / 24 / 32 / 40 / 48 / 64px`

## 5.2 Grid & Container

| Konteks | Max-width | Kolom |
|---|---|---|
| Public Website / Storefront | 1200px | 12 kolom |
| Dashboard (Seller/Admin) | Fluid (sidebar + content) | 12 kolom di content area |

## 5.3 Breakpoints

| Nama | Min-width |
|---|---|
| `sm` | 640px |
| `md` | 768px |
| `lg` | 1024px |
| `xl` | 1280px |

Sesuai IA §19: Desktop = sidebar navigation, Tablet = collapsible sidebar, Mobile = bottom navigation + hamburger.

## 5.4 Radius & Elevation

| Token | Value | Penggunaan |
|---|---|---|
| `radius-sm` | 6px | Input, badge |
| `radius-md` | 10px | Button, card kecil |
| `radius-lg` | 16px | Card besar, modal |
| `shadow-xs` | `0 1px 2px rgba(0,0,0,0.04)` | Card default (subtle, hampir flat) |
| `shadow-md` | `0 4px 12px rgba(0,0,0,0.08)` | Dropdown, popover |
| `shadow-lg` | `0 12px 32px rgba(0,0,0,0.12)` | Modal, dialog |

Prinsip minimalis: hindari shadow tebal. Card cukup dibedakan lewat border tipis (`1px solid neutral-300`) + shadow-xs, bukan shadow besar.

---

# 6. Core Components

## 6.1 Button

| Variant | Style | Penggunaan |
|---|---|---|
| Primary | Solid `brand-600`, text putih | Checkout, Publish, Submit |
| Secondary | Outline neutral-300, text neutral-900 | Cancel, aksi sekunder |
| Ghost | Tanpa border, hover surface subtle | Aksi tersier (icon button) |
| Destructive | Solid `danger` | Delete, Archive, Cancel Order |

Ukuran: `sm` (32px height), `md` (40px, default), `lg` (48px, CTA hero).

## 6.2 Form Elements

- Input height 40px, radius-sm, border `neutral-300`, focus ring `brand-500` (2px, offset).
- Label selalu di atas input (bukan placeholder-only), 13px, `neutral-700`.
- Error state: border `danger` + helper text merah di bawah field.
- Semua form wajib mendukung keyboard navigation & `aria-describedby` untuk error.

## 6.3 Card

Digunakan untuk Product Card, Order Card, Metric Card (Analytics), Tenant Card (Admin).

```
Card
├── Header (opsional): title + action icon
├── Body: content utama
└── Footer (opsional): meta / CTA
```

Style: `surface-0`, border `neutral-300` 1px, radius-lg, padding 20px (desktop) / 16px (mobile).

## 6.4 Data Table (Dashboard/Admin)

- Header row: `neutral-100` background, `text-caption` uppercase, sticky on scroll.
- Row hover: `neutral-100`.
- Row height: 48px (dashboard), 40px (admin — lebih padat, banyak data: tenant list, audit log).
- Aksi row di kolom paling kanan (icon button ghost).
- Selalu sediakan **Empty State**, **Loading (Skeleton row)**, **Error State** sesuai IA §16–18.

## 6.5 Badge / Status Tag

Dipetakan langsung ke state di `user-flow.md` (Order, Payment, Product):

| Status | Warna |
|---|---|
| Active / Published / Paid / Completed | `success` |
| Pending / Waiting Payment / Draft | `warning` |
| Failed / Suspended / Out of Stock | `danger` |
| Archived / Expired | Neutral (`neutral-500` on `neutral-100`) |

Bentuk: pill, radius-full, padding 4px 10px, `text-caption` medium.

## 6.6 Navigation

**Sidebar (Seller/Admin, Desktop):** fixed 240px, background `neutral-0`, active item = `brand-100` background + `brand-600` text + left border 2px accent.

**Bottom Nav (Buyer, Mobile):** 5 item max sesuai IA §8 (Home, Categories, Orders, Wishlist, Profile), icon + label 11px.

**Top Bar:** logo/tenant branding kiri, search tengah (jika ada), profile/notification kanan.

## 6.7 Empty / Loading / Error State (wajib di semua modul list)

- **Empty:** ilustrasi line-art sederhana (bukan foto/stok), 1 kalimat penjelas, 1 CTA. Contoh sudah didefinisikan di IA §16.
- **Loading:** Skeleton loader (bukan spinner) — sesuai IA §18, bentuk skeleton mengikuti layout final konten.
- **Error:** ikon minor + pesan singkat + tombol Retry, nada tidak menyalahkan user.

## 6.8 Modal & Toast

- Modal: max-width 480px (konfirmasi) / 640px (form kompleks), overlay `rgba(0,0,0,0.4)`.
- Toast: muncul kanan-atas (desktop) / bottom (mobile), auto-dismiss 4s, dipakai untuk feedback aksi async (mis. "Produk berhasil dipublikasikan").

---

# 7. Iconography & Imagery

- Icon set: **Lucide Icons** (outline, stroke 1.5–2px, konsisten dengan estetika minimal-modern, open-source, mudah dipakai di Next.js/React).
- Ukuran standar: 16px (inline), 20px (button/nav), 24px (empty state/feature).
- Produk (Storefront): foto asli merchant, rasio 1:1, object-fit cover, radius-md.
- Ilustrasi (empty state, onboarding): gaya line-art monokrom + 1 accent color, hindari ilustrasi ramai/berwarna-warni agar tetap selaras dengan prinsip minimalis.

---

# 8. Motion

Motion dipakai fungsional, bukan dekoratif:

| Interaksi | Durasi | Easing |
|---|---|---|
| Hover/focus | 120ms | ease-out |
| Modal/Dropdown open | 160ms | ease-out + fade |
| Page transition | 200ms | ease-in-out |
| Skeleton shimmer | 1.2s loop | linear |

Hindari animasi besar/bouncy — selaras dengan nada "serius tapi ramah" platform B2B2C.

---

# 9. Accessibility Checklist

- Kontras teks minimum 4.5:1 (body), 3:1 (large text/heading).
- Semua interactive element punya focus ring yang terlihat (jangan `outline: none` tanpa pengganti).
- Target sentuh minimum 40x40px di mobile.
- Form error selalu punya teks, bukan hanya warna.
- Skeleton/loading state punya `aria-busy="true"`.
- Role-based navigation (IA §2.2) tidak boleh merender menu yang tidak boleh diakses, bukan hanya disembunyikan via CSS.

---

# 10. Design Tokens (CSS Variables — referensi implementasi)

```css
:root {
  /* Color — Forest & Stone */
  --color-brand-800: #173C2C;
  --color-brand-700: #1E4D3B;
  --color-brand-600: #2F6B4F;
  --color-brand-100: #D7ECE0;
  --color-brand-50: #F0F7F3;
  --color-neutral-900: #1C1917;
  --color-neutral-700: #44403C;
  --color-neutral-500: #78716C;
  --color-neutral-300: #D6D3D1;
  --color-neutral-100: #F5F5F4;
  --color-neutral-0: #FFFFFF;
  --color-success: #0D9488;
  --color-warning: #D97706;
  --color-danger: #DC2626;
  --color-info: #0284C7;

  /* Typography */
  --font-sans: "Inter", -apple-system, BlinkMacSystemFont, sans-serif;
  --font-mono: "JetBrains Mono", ui-monospace, monospace;

  /* Radius */
  --radius-sm: 6px;
  --radius-md: 10px;
  --radius-lg: 16px;

  /* Shadow */
  --shadow-xs: 0 1px 2px rgba(0,0,0,0.04);
  --shadow-md: 0 4px 12px rgba(0,0,0,0.08);
  --shadow-lg: 0 12px 32px rgba(0,0,0,0.12);

  /* Spacing base unit */
  --space-unit: 4px;
}
```

---

# 11. Surface-Specific Guidelines

## 11.1 Public Website & Storefront (Buyer)
- Density: lega (padding besar, whitespace generous).
- Fokus: foto produk, CTA jelas ("Tambah ke Keranjang", "Checkout").
- Performa dianggap bagian dari desain — skeleton loader wajib untuk katalog (selaras NFR P95 < 100ms di `prd.md` §5.1).

## 11.2 Seller Dashboard
- Density: sedang, data-first tapi tetap bernafas.
- Prioritas: Orders & Products harus ≤ 2 klik dari Dashboard Home (IA §2.1: max 3 klik ke fitur utama).
- Gunakan metric card di Overview untuk ringkasan cepat (Total Order, Revenue, Stok Menipis).

## 11.3 Platform Admin
- Density: padat, mono-font untuk ID/tenant/log, tabel jadi elemen utama.
- Warna lebih netral/serius (minim accent warna cerah), mencerminkan konteks operasional/observability (selaras persona Yusdar di `user-persona.md`).

---

# 12. Related Documents

- Product Vision
- Product Brief
- PRD
- User Persona
- User Stories
- User Flow
- Information Architecture
- Wireframes (belum dibuat)
- Component Library / Storybook (belum dibuat)
