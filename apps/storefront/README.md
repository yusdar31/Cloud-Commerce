# Storefront

**App:** storefront

**Port:** `3000`

**Framework:** Next.js 15 (App Router)

**Owner:** Frontend Team

---

## Purpose

Storefront adalah aplikasi frontend utama CloudCommerce yang menggabungkan:

1. **Public Landing Page** — Marketing, pricing, konversi visitor ke seller
2. **Auth Pages** — Login & Register untuk seller dan buyer
3. **Seller Dashboard** — Manajemen produk, order, analytics
4. **Buyer Storefront** — Jelajah produk, cart, checkout

---

## Technology

| Component | Technology |
|-----------|------------|
| Framework | Next.js 15 (App Router) |
| Language | TypeScript (strict) |
| Styling | Tailwind CSS v4 |
| UI Components | shadcn/ui + Radix UI |
| Icons | Lucide React |
| Animation | Framer Motion |
| Forms | React Hook Form + Zod |
| Server State | TanStack Query v5 |
| Client State | Zustand v5 |
| HTTP Client | ky |

---

## Project Structure

```
apps/storefront/
├── public/
├── src/
│   ├── app/                          # Next.js App Router
│   │   ├── (public)/                 # Landing, pricing pages — NO auth
│   │   │   ├── page.tsx              # Landing page (SUDAH ADA)
│   │   │   ├── pricing/
│   │   │   └── about/
│   │   ├── (auth)/                   # Login, Register
│   │   │   ├── login/
│   │   │   │   └── page.tsx          # SUDAH ADA
│   │   │   └── register/
│   │   │       └── page.tsx          # SUDAH ADA
│   │   ├── (dashboard)/              # Seller dashboard — REQUIRES auth
│   │   │   ├── dashboard/
│   │   │   │   ├── page.tsx          # Overview
│   │   │   │   ├── products/
│   │   │   │   │   ├── page.tsx      # Product list
│   │   │   │   │   ├── new/page.tsx  # Create product
│   │   │   │   │   └── [id]/page.tsx # Edit product
│   │   │   │   ├── orders/
│   │   │   │   │   ├── page.tsx      # Order list
│   │   │   │   │   └── [id]/page.tsx # Order detail
│   │   │   │   └── settings/
│   │   │   │       └── page.tsx      # Store settings
│   │   │   └── layout.tsx
│   │   └── (store)/                  # Buyer storefront
│   │       ├── store/
│   │       │   └── [slug]/
│   │       │       ├── page.tsx      # Store home
│   │       │       └── products/
│   │       │           └── [id]/page.tsx # Product detail
│   │       ├── cart/page.tsx
│   │       ├── checkout/page.tsx
│   │       └── orders/page.tsx
│   │
│   ├── components/
│   │   ├── ui/                       # shadcn/ui base components
│   │   ├── layout/                   # Navbar, Footer, Sidebar
│   │   └── shared/                   # Reusable business components
│   │
│   ├── features/                     # Feature-based modules
│   │   ├── auth/                     # Login, register logic (SUDAH ADA)
│   │   ├── products/                 # Product CRUD
│   │   ├── cart/                     # Cart management
│   │   ├── checkout/                 # Checkout flow
│   │   └── orders/                   # Order list & detail
│   │
│   ├── hooks/                        # Custom React hooks
│   ├── lib/
│   │   ├── api.ts                    # ky API client setup
│   │   └── utils.ts                  # cn(), formatCurrency(), dll
│   ├── stores/
│   │   ├── cart.store.ts             # Zustand cart state
│   │   └── auth.store.ts             # Zustand auth state
│   ├── types/                        # TypeScript types
│   └── middleware.ts                 # Next.js auth middleware (SUDAH ADA)
│
├── next.config.ts
├── tailwind.config.ts
├── tsconfig.json
└── package.json
```

---

## Page Inventory

### ✅ Sudah Ada

| Page | Route | Status |
|------|-------|--------|
| Landing Page | `/` | ✅ Done |
| Login | `/(auth)/login` | ✅ Done |
| Register | `/(auth)/register` | ✅ Done |

### ⏳ Belum Ada (Perlu Dibuat)

| Page | Route | Priority |
|------|-------|----------|
| Seller Dashboard Overview | `/dashboard` | 🔴 P1 |
| Product List | `/dashboard/products` | 🔴 P1 |
| Create Product | `/dashboard/products/new` | 🔴 P1 |
| Edit Product | `/dashboard/products/[id]` | 🔴 P1 |
| Order List (Seller) | `/dashboard/orders` | 🟠 P2 |
| Order Detail (Seller) | `/dashboard/orders/[id]` | 🟠 P2 |
| Store Settings | `/dashboard/settings` | 🟡 P3 |
| Buyer Storefront | `/store/[slug]` | 🟠 P2 |
| Product Detail | `/store/[slug]/products/[id]` | 🟠 P2 |
| Shopping Cart | `/cart` | 🟠 P2 |
| Checkout | `/checkout` | 🟠 P2 |
| Order History (Buyer) | `/orders` | 🟡 P3 |

---

## Environment Variables

| Variable | Wajib | Default (Dev) |
|----------|-------|---------------|
| `NEXT_PUBLIC_API_URL` | ✅ | `http://localhost:8080` |
| `NEXT_PUBLIC_APP_NAME` | ❌ | `CloudCommerce` |
| `NEXT_PUBLIC_APP_URL` | ❌ | `http://localhost:3000` |

File: `apps/storefront/.env.local`

---

## Running Locally

```bash
cd apps/storefront

# Install dependencies (dari root sudah cukup: pnpm install)
pnpm install

# Jalankan dev server
pnpm dev
# → http://localhost:3000

# Lainnya
pnpm build          # Production build
pnpm lint           # ESLint
pnpm type-check     # TypeScript check
```

---

## Route Groups Explanation

```
(public)    → Tidak ada layout auth, tidak perlu login
(auth)      → Layout auth, redirect ke /dashboard jika sudah login
(dashboard) → Layout dashboard, redirect ke /login jika belum login
(store)     → Layout storefront, mixed (ada yang perlu login, ada yang tidak)
```

---

## Auth Flow

```
middleware.ts mengintersep semua request:

/dashboard/* → Cek cookie 'token' → Jika tidak ada: redirect /login
/login       → Jika sudah ada token: redirect /dashboard
/register    → Jika sudah ada token: redirect /dashboard
```

---

## API Client Setup

```typescript
// lib/api.ts — setup ky dengan base URL dan interceptor

import ky from 'ky'

export const api = ky.create({
    prefixUrl: process.env.NEXT_PUBLIC_API_URL,
    hooks: {
        beforeRequest: [
            (request) => {
                const token = getCookie('token')
                if (token) {
                    request.headers.set('Authorization', `Bearer ${token}`)
                }
            }
        ]
    }
})
```

---

## Currency Formatting

Harga dari API selalu dalam IDR (integer, dalam sen).

```typescript
// lib/utils.ts
export function formatPrice(amount: number, currency = 'IDR'): string {
    return new Intl.NumberFormat('id-ID', {
        style: 'currency',
        currency,
        minimumFractionDigits: 0,
    }).format(amount / 100) // convert sen ke rupiah
}

// Contoh: formatPrice(150000) → "Rp150.000"
```

---

## Related Documents

- [Frontend Guidelines](../../docs/04-engineering/05-frontend-guidelines.md)
- [Design System](../../docs/02-design/design-system.md)
- [UI Direction](../../docs/02-design/ui-direction.md)
- [API Guidelines](../../docs/03-architecture/api-guidelines.md)
- [Error Handling Catalog](../../docs/04-engineering/error-handling-catalog.md)
