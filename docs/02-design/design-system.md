# Design System

**Project:** CloudCommerce

**Version:** 1.0.0

**Status:** Draft

**Owner:** Design System Team

**Last Updated:** July 2026

---

# 1. Overview

## Purpose

Design System mendefinisikan seluruh komponen UI yang digunakan pada CloudCommerce agar:

- Konsisten
- Reusable
- Accessible
- Responsive
- Mudah dikembangkan

Semua halaman wajib menggunakan komponen dari Design System.

---

# 2. Design Principles

Setiap komponen harus memenuhi lima prinsip utama:

- Reusable
- Accessible (WCAG 2.2 AA)
- Responsive
- Themeable (Light & Dark)
- Composable

---

# 3. Component Categories

CloudCommerce membagi komponen menjadi beberapa kategori:

```
Foundations
│
├── Colors
├── Typography
├── Spacing
├── Motion
└── Icons

Components
│
├── Buttons
├── Forms
├── Navigation
├── Feedback
├── Data Display
├── Overlays
└── Charts

Patterns
│
├── Dashboard
├── Authentication
├── Checkout
├── Settings
└── Tables
```

---

# 4. Button

## Purpose

Digunakan untuk memicu aksi.

---

## Variants

- Primary
- Secondary
- Outline
- Ghost
- Danger
- Success
- Link

---

## Sizes

- Small
- Medium
- Large
- Icon

---

## States

- Default
- Hover
- Focus
- Active
- Disabled
- Loading

---

## Rules

✔ Maksimal satu Primary Button per section.

✔ Gunakan Danger hanya untuk aksi destruktif.

✔ Loading hanya muncul pada button yang ditekan.

---

# 5. Input

## Types

- Text
- Email
- Password
- Number
- Search
- URL
- Phone

---

## States

- Empty
- Filled
- Focus
- Error
- Disabled
- Success

---

## Features

- Label
- Placeholder
- Helper Text
- Validation
- Character Counter (optional)

---

# 6. Select

Support:

- Single Select
- Multi Select
- Searchable
- Async Search

---

# 7. Textarea

Support:

- Auto Resize
- Character Counter
- Validation

---

# 8. Checkbox

Support:

- Checked
- Unchecked
- Indeterminate

---

# 9. Radio Group

Gunakan jika pilihan hanya satu.

---

# 10. Switch

Digunakan untuk:

ON / OFF

Feature Toggle

---

# 11. Date Picker

Support

- Single Date
- Date Range
- Time
- Timezone (future)

---

# 12. File Upload

Support:

- Drag & Drop
- Browse
- Progress
- Preview
- Remove

---

# 13. Avatar

Variants

- Image
- Initial
- Icon

Sizes

XS

SM

MD

LG

XL

---

# 14. Badge

Variants

- Primary
- Success
- Warning
- Danger
- Info
- Neutral

---

# 15. Tag

Digunakan untuk kategori.

---

# 16. Chip

Digunakan untuk filter.

---

# 17. Card

Jenis

- Default
- Metric
- Product
- Activity
- Analytics

---

Card harus:

- Memiliki header yang jelas.
- Tidak menggunakan shadow berlebihan.
- Mendukung skeleton loading.

---

# 18. Table

## Features

- Search
- Sort
- Filter
- Pagination
- Bulk Action
- Column Visibility
- Row Selection
- Sticky Header

---

## Row Actions

Hover →

Edit

Duplicate

Archive

Delete

---

# 19. Data List

Digunakan untuk:

Notifications

Activity

Orders

Timeline

---

# 20. Tabs

Support:

- Horizontal
- Scrollable
- Icon + Label

---

# 21. Accordion

Digunakan untuk:

FAQ

Advanced Settings

---

# 22. Breadcrumb

Contoh

Dashboard

>

Products

>

MacBook Pro

---

# 23. Sidebar

Support

- Collapse
- Nested Menu (maks. 2 level)
- Active Indicator
- Notification Badge

---

# 24. Top Navigation

Berisi:

- Search
- Notifications
- Theme Toggle
- User Menu

---

# 25. Command Palette

Shortcut

Ctrl + K

Features

- Search Products
- Search Orders
- Search Customers
- Quick Navigation
- Quick Actions

---

# 26. Modal

Variants

- Confirmation
- Form
- Information
- Fullscreen (mobile)

---

# 27. Drawer

Digunakan untuk:

Quick Edit

Order Detail

Customer Detail

---

# 28. Tooltip

Maksimal dua baris.

---

# 29. Popover

Digunakan untuk aksi ringan.

---

# 30. Toast

Variants

- Success
- Error
- Warning
- Info

Support:

Undo Action

---

# 31. Alert

Jenis

- Information
- Warning
- Error
- Success

---

# 32. Empty State

Harus memiliki:

- Illustration/Icon
- Title
- Description
- CTA

---

# 33. Loading

Gunakan:

- Skeleton
- Button Spinner
- Progress Bar

Hindari spinner fullscreen.

---

# 34. Error State

Harus menyediakan:

- Penjelasan
- Retry
- Contact Support (opsional)

---

# 35. Charts

Jenis

- Line
- Area
- Bar
- Donut
- Sparkline

Tidak menggunakan chart 3D.

---

# 36. Dashboard Widgets

Widgets standar:

- KPI Card
- Revenue Chart
- Recent Orders
- Inventory Alert
- Activity Feed
- AI Insight
- Quick Actions

---

# 37. Authentication Pattern

Komponen:

- Login Form
- Register Form
- OTP Input
- Password Strength
- Social Login (future)

---

# 38. Checkout Pattern

Komponen:

- Cart Summary
- Shipping Form
- Payment Method
- Order Summary
- Confirmation

---

# 39. Settings Pattern

Gunakan layout:

```
Sidebar

↓

Section

↓

Card

↓

Save Button
```

---

# 40. Feedback Pattern

Semua aksi harus menghasilkan feedback:

- Toast
- Inline Validation
- Loading State
- Success State

---

# 41. Accessibility

Semua komponen harus:

- Keyboard Accessible
- Screen Reader Friendly
- Focus Visible
- ARIA Label bila diperlukan

---

# 42. Responsive Rules

Desktop

Sidebar

Tablet

Collapsible Sidebar

Mobile

Bottom Navigation

Floating Action Button (opsional)

---

# 43. Naming Convention

Gunakan nama yang konsisten:

- Button
- Input
- DataTable
- EmptyState
- PageHeader
- CommandPalette
- MetricCard

---

# 44. Implementation Mapping

| Design System | Frontend |
|---------------|----------|
| Button | shadcn Button |
| Dialog | Radix Dialog |
| Tooltip | Radix Tooltip |
| Table | TanStack Table |
| Chart | Tremor / Recharts |
| Icons | Lucide |
| Motion | Framer Motion |
| Forms | React Hook Form |
| Validation | Zod |

---

# 45. Component Lifecycle

```
Requirement

↓

Wireframe

↓

Component Design

↓

Review

↓

Implementation

↓

Testing

↓

Documentation

↓

Release
```

---

# 46. Future Components

CloudCommerce dapat menambahkan:

- Kanban Board
- Calendar
- Rich Text Editor
- AI Chat Panel
- Notification Center
- Workflow Builder
- File Manager
- White-label Theme Picker

Tanpa mengubah fondasi Design System.

---

# 47. References

- Shopify Polaris
- GitHub Primer
- Radix UI
- shadcn/ui
- TanStack Table
- Tremor
- Framer Motion
- WCAG 2.2