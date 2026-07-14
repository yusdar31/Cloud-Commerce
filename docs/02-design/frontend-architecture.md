# Frontend Architecture - CloudCommerce

**Version:** 1.0.0
**Last Updated:** July 2026
**Status:** In Progress

---

## Overview

CloudCommerce membutuhkan **3 frontend applications** yang berbeda untuk 3 jenis user:

| Application | Target User | Tech Stack | Status |
|-------------|-------------|------------|--------|
| **Storefront** (`apps/storefront`) | Buyers/Customers | Next.js 15, React 19, Tailwind v4 | 60% Complete |
| **Seller Dashboard** (`apps/web`) | Merchants/Sellers | Next.js 15, React 19, Tailwind v4 | Not Started |
| **Admin Panel** (TBD) | Platform Admin | Next.js 15, React 19, Tailwind v4 | Not Started |

---

## 1. Storefront (Buyer Application)

**Path:** `apps/storefront/`
**Base URL:** `/`
**Target Users:** End customers, guest visitors

### 1.1 Pages - Completed

| Page | Route | Description | Status |
|------|-------|-------------|--------|
| Landing Page | `/` | Marketing homepage with hero, features, pricing, FAQ | ✅ Complete |
| Login | `/login` | User login form | ✅ Complete |
| Register | `/register` | User registration form | ✅ Complete |
| Products List | `/products` | Browse all products with search/filter | ✅ Complete |
| Product Detail | `/products/[slug]` | Single product view with images, add to cart | ✅ Complete |
| Cart | `/cart` | Shopping cart with quantity management | ✅ Complete |
| Checkout | `/checkout` | Shipping & payment forms | ✅ Complete |
| Order Success | `/order-success` | Order confirmation page | ✅ Complete |

### 1.2 Pages - Pending

| Page | Route | Description | Status |
|------|-------|-------------|--------|
| User Dashboard | `/dashboard` | User profile, order history | 🔶 Shell Only |
| Order History | `/dashboard/orders` | List of past orders | ❌ Not Started |
| Order Detail | `/dashboard/orders/[id]` | Single order detail | ❌ Not Started |
| User Profile | `/dashboard/profile` | Edit profile, change password | ❌ Not Started |
| Wishlist | `/wishlist` | Saved products | ❌ Not Started |
| Search Results | `/search?q=query` | Search results page | ❌ Not Started |
| Categories | `/categories` | Browse by category | ❌ Not Started |
| Category Products | `/categories/[slug]` | Products in category | ❌ Not Started |

### 1.3 Components - Completed

| Component | Path | Description |
|-----------|------|-------------|
| Button | `components/ui/button.tsx` | Primary, secondary, ghost variants |
| Card | `components/ui/card.tsx` | Card container |
| Input | `components/ui/input.tsx` | Form input |
| Label | `components/ui/label.tsx` | Form label |
| Badge | `components/ui/badge.tsx` | Status badges |
| Navbar | `components/layout/navbar.tsx` | Top navigation |
| Footer | `components/layout/footer.tsx` | Footer links |
| Header | `components/layout/header.tsx` | App navigation header |
| FadeIn | `components/animations/fade-in.tsx` | Fade-in animation |
| Stagger | `components/animations/stagger.tsx` | Stagger animation |

### 1.4 Components - Needed

| Component | Description | Priority |
|-----------|-------------|----------|
| Skeleton | Loading placeholders | High |
| Toast | Notification system | High |
| Dialog/Modal | Confirmation dialogs | High |
| Select | Dropdown select | High |
| Checkbox | Form checkbox | Medium |
| Tabs | Tab navigation | Medium |
| Accordion | Collapsible content | Medium |
| Dropdown Menu | User menu dropdown | Medium |
| Separator | Visual divider | Low |
| Avatar | User avatar | Low |
| Progress | Loading progress | Low |

### 1.5 State Management

| Store | File | Description | Status |
|-------|------|-------------|--------|
| Auth Store | `stores/auth-store.ts` | User authentication state | ✅ Complete |
| Cart Store | `stores/cart-store.ts` | Shopping cart state | ✅ Complete |

### 1.6 State - Needed

| Store | Description | Priority |
|-------|-------------|----------|
| Wishlist Store | Saved products | Medium |
| Search Store | Search state | Low |
| UI Store | Theme, sidebar state | Low |

---

## 2. Seller Dashboard (Merchant Application)

**Path:** `apps/web/`
**Base URL:** `/` (separate app)
**Target Users:** Merchants, store owners

### 2.1 Core Pages

| Page | Route | Description | Priority |
|------|-------|-------------|----------|
| Dashboard Home | `/` | Overview stats, recent orders | High |
| Products List | `/products` | All products with filters | High |
| Product Create | `/products/new` | Add new product | High |
| Product Edit | `/products/[id]/edit` | Edit product | High |
| Orders List | `/orders` | All orders with filters | High |
| Order Detail | `/orders/[id]` | Order detail, fulfillment | High |
| Customers | `/customers` | Customer list | Medium |
| Analytics | `/analytics` | Sales, revenue charts | Medium |
| Settings | `/settings` | Store settings | Medium |
| Profile | `/settings/profile` | Merchant profile | Medium |
| Payments | `/settings/payments` | Payment gateway config | Medium |
| Team | `/settings/team` | Team members | Low |

### 2.2 Dashboard Components

| Component | Description | Priority |
|-----------|-------------|----------|
| Sidebar | Navigation sidebar | High |
| Top Bar | Header with user menu | High |
| Stats Card | Metric cards (revenue, orders) | High |
| Data Table | Sortable, filterable table | High |
| Charts | Line, bar, pie charts | High |
| Date Picker | Date range filter | High |
| Product Form | Multi-step product form | High |
| Order Status Badge | Status indicators | High |
| Image Uploader | Drag & drop upload | Medium |
| Rich Text Editor | Product description | Medium |
| Notification Dropdown | Alerts menu | Medium |

### 2.3 Dashboard Layout

```
┌─────────────────────────────────────────────────────────┐
│ Top Bar (Logo, Search, Notifications, User Menu)        │
├──────────┬──────────────────────────────────────────────┤
│          │                                              │
│ Sidebar  │              Content Area                   │
│          │                                              │
│ - Home   │   ┌─────────────────────────────────────┐   │
│ - Orders │   │                                     │   │
│ - Products│   │                                     │   │
│ - Customers│   │                                     │   │
│ - Analytics│   │                                     │   │
│ - Settings │   │                                     │   │
│          │   │                                     │   │
│          │   └─────────────────────────────────────┘   │
└──────────┴──────────────────────────────────────────────┘
```

---

## 3. Admin Panel (Platform Administration)

**Path:** `apps/admin/` (TBD)
**Base URL:** `/admin`
**Target Users:** Platform admins, support team

### 3.1 Core Pages

| Page | Route | Description | Priority |
|------|-------|-------------|----------|
| Dashboard | `/` | Platform metrics | High |
| Tenants | `/tenants` | All merchants | High |
| Tenant Detail | `/tenants/[id]` | Single tenant view | High |
| Users | `/users` | All users | High |
| User Detail | `/users/[id]` | User detail | Medium |
| Orders | `/orders` | All orders | Medium |
| Audit Logs | `/audit` | System audit trail | Medium |
| System Health | `/health` | Service status | Medium |
| Settings | `/settings` | Platform config | Low |

### 3.2 Admin Components

| Component | Description | Priority |
|-----------|-------------|----------|
| Admin Sidebar | Navigation | High |
| Tenant Card | Tenant overview | High |
| Audit Log Table | Activity log | Medium |
| Health Dashboard | Service monitoring | Medium |
| Metric Charts | Platform analytics | Medium |

---

## 4. Design System

### 4.1 Color Palette (Forest & Stone)

| Token | Hex | Usage |
|-------|-----|-------|
| `brand-800` | `#173C2C` | Text on tint, active state |
| `brand-700` | `#1E4D3B` | Primary action, links |
| `brand-600` | `#2F6B4F` | Hover state |
| `brand-100` | `#D7ECE0` | Background highlight |
| `brand-50` | `#F0F7F3` | Surface highlight |
| `neutral-900` | `#1C1917` | Heading text |
| `neutral-700` | `#44403C` | Body text |
| `neutral-500` | `#78716C` | Secondary text |
| `neutral-300` | `#D6D3D1` | Border |
| `neutral-100` | `#F5F5F4` | Surface |
| `success` | `#0D9488` | Success states |
| `warning` | `#D97706` | Warning states |
| `danger` | `#DC2626` | Error states |
| `info` | `#0284C7` | Info states |

### 4.2 Typography

| Token | Size | Weight | Usage |
|-------|------|--------|-------|
| `display` | 36px | 700 | Landing hero |
| `h1` | 28px | 700 | Page title |
| `h2` | 22px | 600 | Section title |
| `h3` | 18px | 600 | Card title |
| `body` | 15px | 400 | Body text |
| `body-sm` | 13px | 400 | Secondary text |
| `caption` | 12px | 500 | Labels |

### 4.3 Spacing Scale

```
4px / 8px / 12px / 16px / 20px / 24px / 32px / 40px / 48px / 64px
```

### 4.4 Border Radius

| Token | Value | Usage |
|-------|-------|-------|
| `sm` | 6px | Input, badge |
| `md` | 10px | Button, card |
| `lg` | 16px | Modal |

---

## 5. API Integration

### 5.1 API Client Structure

```typescript
// lib/api.ts
├── authApi
│   ├── register()
│   ├── login()
│   ├── refresh()
│   ├── logout()
│   └── getProfile()
├── productApi
│   ├── list()
│   ├── getBySlug()
│   ├── create()      // Seller only
│   ├── update()      // Seller only
│   └── delete()      // Seller only
├── orderApi
│   ├── create()
│   ├── list()
│   ├── getById()
│   └── updateStatus() // Seller only
├── cartApi
│   ├── get()
│   ├── addItem()
│   ├── updateItem()
│   └── removeItem()
└── userApi
    ├── getProfile()
    └── updateProfile()
```

### 5.2 Required API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `POST /api/v1/auth/register` | POST | User registration |
| `POST /api/v1/auth/login` | POST | User login |
| `POST /api/v1/auth/refresh` | POST | Refresh token |
| `GET /api/v1/users/me` | GET | Current user |
| `GET /api/v1/products` | GET | List products |
| `GET /api/v1/products/:id` | GET | Get product |
| `POST /api/v1/products` | POST | Create product |
| `PUT /api/v1/products/:id` | PUT | Update product |
| `DELETE /api/v1/products/:id` | DELETE | Delete product |
| `POST /api/v1/orders` | POST | Create order |
| `GET /api/v1/orders` | GET | List orders |
| `GET /api/v1/orders/:id` | GET | Get order |
| `POST /api/v1/payments` | POST | Create payment |
| `POST /webhooks/payments/:provider` | POST | Payment webhook |

---

## 6. File Structure

```
apps/
├── storefront/              # Buyer application
│   ├── src/
│   │   ├── app/
│   │   │   ├── (auth)/
│   │   │   │   ├── login/
│   │   │   │   └── register/
│   │   │   ├── (dashboard)/
│   │   │   │   ├── dashboard/
│   │   │   │   ├── orders/
│   │   │   │   └── profile/
│   │   │   ├── products/
│   │   │   ├── cart/
│   │   │   ├── checkout/
│   │   │   ├── order-success/
│   │   │   ├── layout.tsx
│   │   │   ├── page.tsx
│   │   │   └── globals.css
│   │   ├── components/
│   │   │   ├── ui/
│   │   │   ├── layout/
│   │   │   ├── animations/
│   │   │   └── providers/
│   │   ├── features/
│   │   ├── stores/
│   │   ├── lib/
│   │   └── types/
│   └── package.json
│
├── web/                     # Seller Dashboard
│   ├── src/
│   │   ├── app/
│   │   │   ├── (dashboard)/
│   │   │   │   ├── page.tsx
│   │   │   │   ├── products/
│   │   │   │   ├── orders/
│   │   │   │   ├── customers/
│   │   │   │   ├── analytics/
│   │   │   │   └── settings/
│   │   │   ├── layout.tsx
│   │   │   └── globals.css
│   │   ├── components/
│   │   ├── stores/
│   │   ├── lib/
│   │   └── types/
│   └── package.json
│
└── admin/                   # Platform Admin (TBD)
    └── ...
```

---

## 7. Implementation Priority

### Phase 1: Storefront Completion (Current)
1. ✅ Products list & detail pages
2. ✅ Cart & checkout flow
3. 🔶 User dashboard (profile, orders)
4. ❌ Wishlist functionality
5. ❌ Search & categories

### Phase 2: Seller Dashboard
1. ❌ Dashboard layout (sidebar, topbar)
2. ❌ Products management (CRUD)
3. ❌ Orders management
4. ❌ Analytics
5. ❌ Settings

### Phase 3: Admin Panel
1. ❌ Admin layout
2. ❌ Tenant management
3. ❌ User management
4. ❌ Audit logs
5. ❌ System health

---

## 8. Component Library Status

### shadcn/ui Components - Installed

| Component | Status |
|-----------|--------|
| Button | ✅ |
| Card | ✅ |
| Input | ✅ |
| Label | ✅ |
| Badge | ✅ |

### shadcn/ui Components - Needed

| Component | Priority |
|-----------|----------|
| Skeleton | High |
| Toast | High |
| Dialog | High |
| Select | High |
| Checkbox | Medium |
| Tabs | Medium |
| Accordion | Medium |
| Dropdown Menu | Medium |
| Avatar | Medium |
| Separator | Low |
| Progress | Low |
| Table | Medium |
| Form | Medium |
| Command | Low |
| Popover | Low |

---

## 9. Performance Targets

| Metric | Target |
|--------|--------|
| First Contentful Paint | < 1.5s |
| Largest Contentful Paint | < 2.5s |
| Time to Interactive | < 3.0s |
| Cumulative Layout Shift | < 0.1 |
| Bundle Size (Initial) | < 150KB |

---

## 10. Accessibility (WCAG 2.1 AA)

- [ ] Color contrast ratio ≥ 4.5:1
- [ ] Focus visible on all interactive elements
- [ ] Keyboard navigation support
- [ ] Screen reader compatible
- [ ] Form error messages accessible
- [ ] Image alt text
- [ ] ARIA labels where needed
