# Frontend Guidelines

**Project:** CloudCommerce

**Version:** 1.0.0

**Status:** Approved

**Owner:** Frontend Team

**Last Updated:** July 2026

---

# 1. Purpose

Dokumen ini mendefinisikan standar pengembangan frontend untuk CloudCommerce.

Tujuan:

- Konsistensi antar halaman
- Performance optimal
- Mudah di-maintain
- SEO friendly
- Mobile responsive

---

# 2. Technology Stack

| Component | Technology |
|-----------|------------|
| Framework | Next.js 15 (App Router) |
| Language | TypeScript (strict) |
| Styling | Tailwind CSS v4 |
| UI Library | shadcn/ui + Radix UI |
| Icons | Lucide React |
| Animation | Motion (framer-motion) |
| Forms | React Hook Form + Zod |
| Server State | TanStack Query |
| Client State | Zustand |
| HTTP Client | fetch / ky |
| Testing | Vitest + Playwright |

---

# 3. Project Structure

```
apps/storefront/
│
├── public/
│   └── images/
│
├── src/
│   ├── app/                    # Next.js App Router pages
│   │   ├── (public)/           # Landing, Pricing, etc.
│   │   ├── (auth)/             # Login, Register
│   │   ├── (dashboard)/        # Seller dashboard routes
│   │   └── (store)/            # Buyer storefront routes
│   │
│   ├── components/             # Shared components
│   │   ├── ui/                 # shadcn/ui components
│   │   ├── layout/             # Navbar, Sidebar, Footer
│   │   ├── shared/             # Reusable business components
│   │   └── widgets/            # Dashboard widgets
│   │
│   ├── features/               # Feature-based modules
│   │   ├── auth/
│   │   ├── products/
│   │   ├── cart/
│   │   ├── checkout/
│   │   └── orders/
│   │
│   ├── hooks/                  # Custom hooks
│   ├── lib/                    # Utilities, API client
│   ├── stores/                 # Zustand stores
│   ├── types/                  # TypeScript types
│   └── middleware.ts           # Next.js middleware (auth)
│
├── tailwind.config.ts
├── tsconfig.json
├── next.config.ts
└── package.json
```

---

# 4. Component Architecture

## Server Components (default)

```tsx
// app/(store)/products/page.tsx
export default async function ProductsPage() {
    const products = await getProducts()

    return <ProductList products={products} />
}
```

Gunakan Server Component default. Client Component hanya jika perlu interaktivitas.

## Client Components (when needed)

```tsx
"use client"

import { useState } from "react"

export function AddToCartButton({ productId }: { productId: string }) {
    const [loading, setLoading] = useState(false)

    return (
        <Button onClick={handleAdd} loading={loading}>
            Add to Cart
        </Button>
    )
}
```

---

# 5. Data Fetching

## Server-side (RSC)

```tsx
async function getProducts() {
    const res = await fetch(`${API_URL}/api/v1/products`, {
        next: { revalidate: 60 },   // ISR every 60s
    })
    return res.json()
}
```

## Client-side (TanStack Query)

```tsx
"use client"

import { useQuery } from "@tanstack/react-query"

function useProducts() {
    return useQuery({
        queryKey: ["products", tenantId],
        queryFn: () => fetchProducts(tenantId),
    })
}
```

## Mutation

```tsx
const mutation = useMutation({
    mutationFn: (data: CreateProduct) => api.createProduct(data),
    onSuccess: () => {
        queryClient.invalidateQueries({ queryKey: ["products"] })
        toast.success("Product created")
    },
})
```

---

# 6. State Management

## Server State (TanStack Query)

All data from API.

## Client State (Zustand)

```tsx
import { create } from "zustand"

interface CartStore {
    items: CartItem[]
    addItem: (item: CartItem) => void
    removeItem: (id: string) => void
    clear: () => void
}

export const useCartStore = create<CartStore>((set) => ({
    items: [],
    addItem: (item) => set((state) => ({ items: [...state.items, item] })),
    removeItem: (id) => set((state) => ({
        items: state.items.filter((i) => i.id !== id),
    })),
    clear: () => set({ items: [] }),
}))
```

Only use Zustand for:

- Cart state
- UI state (sidebar, theme)
- Auth tokens (cache)

---

# 7. Forms

```tsx
"use client"

import { useForm } from "react-hook-form"
import { z } from "zod"
import { zodResolver } from "@hookform/resolvers/zod"

const productSchema = z.object({
    name: z.string().min(1, "Name is required"),
    price: z.number().min(1, "Price must be > 0"),
    stock: z.number().int().min(0),
})

type ProductForm = z.infer<typeof productSchema>

function ProductForm() {
    const form = useForm<ProductForm>({
        resolver: zodResolver(productSchema),
    })

    return (
        <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit)}>
                <FormField name="name" label="Product Name" />
                <FormField name="price" label="Price" type="number" />
                <Button type="submit">Save</Button>
            </form>
        </Form>
    )
}
```

---

# 8. Routing (App Router)

## Public Routes

```
/                           → Landing
/pricing                    → Pricing
/login                      → Login
/register                   → Register
```

## Seller Dashboard

```
/dashboard                  → Overview
/dashboard/products         → Product list
/dashboard/products/new     → Create product
/dashboard/orders           → Order list
/dashboard/orders/:id       → Order detail
/dashboard/analytics        → Analytics
/dashboard/settings          → Store settings
```

## Buyer Storefront

```
/store/:slug                → Store home
/store/:slug/products/:id   → Product detail
/cart                       → Shopping cart
/checkout                   → Checkout
/orders                     → Order history
/profile                    → User profile
```

## Admin

```
/admin                      → Admin dashboard
/admin/tenants              → Tenant management
/admin/monitoring           → System monitoring
```

---

# 9. Middleware (Auth)

```ts
// middleware.ts
import { NextResponse } from "next/server"
import type { NextRequest } from "next/server"

export function middleware(request: NextRequest) {
    const token = request.cookies.get("token")
    const { pathname } = request.nextUrl

    // Protect dashboard routes
    if (pathname.startsWith("/dashboard") && !token) {
        return NextResponse.redirect(new URL("/login", request.url))
    }

    // Redirect if already logged in
    if (pathname === "/login" && token) {
        return NextResponse.redirect(new URL("/dashboard", request.url))
    }

    return NextResponse.next()
}

export const config = {
    matcher: ["/dashboard/:path*", "/login", "/register"],
}
```

---

# 10. Styling (Tailwind + shadcn/ui)

Always use `cn()` utility for conditional classes:

```tsx
import { cn } from "@/lib/utils"

<div className={cn(
    "flex items-center gap-2 p-4 rounded-lg",
    active ? "bg-primary text-white" : "bg-card",
)} />
```

Component variants using `cva`:

```tsx
import { cva } from "class-variance-authority"

const buttonVariants = cva(
    "inline-flex items-center justify-center rounded-md text-sm font-medium",
    {
        variants: {
            variant: {
                primary: "bg-primary text-white hover:bg-primary/90",
                ghost: "hover:bg-accent hover:text-accent-foreground",
                danger: "bg-destructive text-white",
            },
            size: {
                sm: "h-8 px-3",
                md: "h-10 px-4",
                lg: "h-12 px-6",
            },
        },
    }
)
```

---

# 11. Performance Rules

- Server Component by default
- Client Component only when needed
- Lazy load images: `next/image`
- Dynamic import heavy components: `next/dynamic`
- Debounce search input: 300ms
- Pagination with cursor-based
- Minimize bundle: tree-shake unused imports

---

# 12. Accessibility

- All images must have `alt` text
- Forms must have labels
- Interactive elements must be keyboard accessible
- Use semantic HTML (`<nav>`, `<main>`, `<section>`, `<button>`)
- Color contrast must meet WCAG AA

---

# 13. Error Handling

```tsx
function ErrorBoundary({ error, reset }: { error: Error; reset: () => void }) {
    return (
        <div className="flex flex-col items-center gap-4 p-8">
            <h2>Something went wrong</h2>
            <p className="text-muted-foreground">{error.message}</p>
            <Button onClick={reset}>Try again</Button>
        </div>
    )
}
```

Always implement:

- `error.tsx` (per route segment)
- `loading.tsx` (skeleton)
- `not-found.tsx` (custom 404)

---

# 14. Empty State

Every list page must have empty state:

```tsx
function EmptyState({ title, description, action }: Props) {
    return (
        <div className="flex flex-col items-center gap-4 py-16">
            <PackageIcon className="h-16 w-16 text-muted-foreground" />
            <h3 className="text-lg font-semibold">{title}</h3>
            <p className="text-sm text-muted-foreground">{description}</p>
            {action}
        </div>
    )
}
```

---

# 15. Related Documents

- Technology Stack
- Monorepo Structure
- Design System
- UI Direction
- API Guidelines
