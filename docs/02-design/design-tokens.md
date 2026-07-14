# Design Tokens

**Project:** CloudCommerce

**Version:** 1.0.0

**Status:** Draft

**Owner:** Design System Team

**Last Updated:** July 2026

---

# 1. Overview

## Purpose

Design Tokens adalah sumber kebenaran (Single Source of Truth) untuk seluruh visual CloudCommerce.

Semua komponen UI harus menggunakan token ini.

Token akan digunakan pada:

- Design System
- Figma
- Tailwind CSS
- shadcn/ui
- React Components
- Mobile App (Future)

---

# 2. Design Philosophy

CloudCommerce menggunakan desain yang:

- Calm
- Modern
- Interactive
- Enterprise
- Developer-first

Token harus sederhana, mudah dipahami, dan konsisten.

---

# 3. 8px Grid System

Semua spacing mengikuti kelipatan 8.

| Token | Value |
|--------|------:|
| space-1 | 4px |
| space-2 | 8px |
| space-3 | 12px |
| space-4 | 16px |
| space-5 | 20px |
| space-6 | 24px |
| space-8 | 32px |
| space-10 | 40px |
| space-12 | 48px |
| space-16 | 64px |
| space-20 | 80px |
| space-24 | 96px |

---

# 4. Border Radius

CloudCommerce tidak menggunakan radius yang terlalu besar.

| Token | Value |
|--------|------:|
| radius-none | 0 |
| radius-xs | 4px |
| radius-sm | 6px |
| radius-md | 8px |
| radius-lg | 12px |
| radius-xl | 16px |
| radius-full | 9999px |

Default:

Cards → radius-md

Button → radius-md

Input → radius-md

Modal → radius-lg

---

# 5. Color System

CloudCommerce menggunakan semantic color.

Bukan:

Blue

Red

Green

Tetapi:

Primary

Success

Warning

Danger

Info

Neutral

---

## Primary

| Shade | Usage |
|--------|--------|
| 50 | Hover Background |
| 100 | Light Background |
| 200 | Selected |
| 300 | Border |
| 400 | Hover |
| 500 | Primary |
| 600 | Primary Hover |
| 700 | Active |
| 800 | Dark |
| 900 | Darkest |

---

## Neutral

Digunakan untuk:

Background

Border

Card

Text

Divider

Sidebar

---

## Success

Order Success

Payment Success

Published

Completed

---

## Warning

Pending

Stock Low

Subscription Expiring

---

## Danger

Delete

Payment Failed

Validation Error

Suspended

---

## Info

Notification

Tips

Information

---

# 6. Light Theme

Background

Surface

Card

Sidebar

Primary Text

Secondary Text

Border

Divider

---

# 7. Dark Theme

Dark Mode menggunakan token yang sama.

Hanya nilainya berbeda.

Tidak membuat token baru.

---

# 8. Typography

Font Family

Primary

Inter

Secondary

Geist

Monospace

JetBrains Mono

---

## Font Weight

Regular

400

Medium

500

Semibold

600

Bold

700

---

## Font Scale

| Token | Size |
|---------|------|
| Display | 48px |
| H1 | 36px |
| H2 | 30px |
| H3 | 24px |
| H4 | 20px |
| H5 | 18px |
| Body Large | 16px |
| Body | 14px |
| Small | 13px |
| Caption | 12px |

---

## Line Height

Display

120%

Heading

130%

Body

150%

Caption

140%

---

# 9. Shadow

CloudCommerce menggunakan shadow yang tipis.

| Token | Usage |
|---------|--------|
| shadow-xs | Input Focus |
| shadow-sm | Card Hover |
| shadow-md | Dropdown |
| shadow-lg | Modal |
| shadow-xl | Floating Panel |

Tidak menggunakan shadow berlebihan.

---

# 10. Elevation

Level 0

Background

Level 1

Card

Level 2

Dropdown

Level 3

Modal

Level 4

Command Palette

---

# 11. Border

Border Default

1px

Border Strong

2px

Divider

1px

Dashed

1px

---

# 12. Icon Size

XS

14px

SM

16px

MD

20px

LG

24px

XL

32px

Semua icon menggunakan:

Lucide Icons

---

# 13. Animation

Durasi standar.

Fast

150ms

Normal

200ms

Slow

300ms

Extra Slow

500ms

---

## Easing

Ease Out

Ease In Out

Spring

Digunakan hanya jika diperlukan.

---

# 14. Motion Rules

Hover

↓

150ms

Dropdown

↓

200ms

Sidebar

↓

250ms

Modal

↓

200ms

Page Transition

↓

300ms

---

# 15. Z Index

| Layer | Value |
|---------|------:|
| Base | 1 |
| Header | 10 |
| Sidebar | 20 |
| Dropdown | 30 |
| Modal | 40 |
| Toast | 50 |
| Command Palette | 60 |

---

# 16. Responsive Breakpoints

| Device | Width |
|---------|------:|
| Mobile | 390px |
| Tablet | 768px |
| Laptop | 1024px |
| Desktop | 1280px |
| Wide | 1440px |
| Ultra Wide | 1600px |

---

# 17. Container Width

Content

1200px

Dashboard

1440px

Landing

1280px

Documentation

960px

---

# 18. Component Tokens

## Button

Height

40px

Padding

16px

Radius

8px

---

## Input

Height

40px

Padding

12px

Radius

8px

---

## Card

Padding

24px

Radius

8px

Border

1px

---

## Modal

Radius

12px

Padding

24px

Width

640px

---

## Table

Row Height

48px

Header Height

52px

Cell Padding

16px

---

# 19. State Colors

Hover

Primary +10%

Active

Primary +20%

Disabled

Opacity 50%

Focus

Outline Ring

Loading

Skeleton

---

# 20. Accessibility Tokens

Minimum Contrast

WCAG AA

Focus Ring

2px

Minimum Touch Target

44px

Keyboard Navigation

Required

---

# 21. Future Tokens

CloudCommerce nantinya akan mendukung:

- White Label Theme
- Tenant Branding
- Theme Builder
- Seasonal Theme
- High Contrast Mode

Tanpa mengubah struktur token.

---

# 22. Naming Convention

Gunakan semantic naming.

Contoh:

✓ color-primary

✓ color-success

✓ radius-md

✓ shadow-sm

✓ space-4

✓ text-body

Bukan:

✖ blue-500

✖ red

✖ big-radius

✖ shadow-dark

---

# 23. Token Governance

Perubahan token harus:

1. Diusulkan melalui Pull Request.
2. Direview oleh Design System Owner.
3. Diuji pada Light & Dark Mode.
4. Diverifikasi tidak merusak komponen yang sudah ada.

---

# 24. References

- Material Design 3
- Atlassian Design System
- Shopify Polaris
- GitHub Primer
- Radix UI
- shadcn/ui
- Tailwind CSS