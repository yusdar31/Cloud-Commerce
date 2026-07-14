---
name: CloudCommerce Global
colors:
  surface: '#faf8ff'
  surface-dim: '#d9d9e5'
  surface-bright: '#faf8ff'
  surface-container-lowest: '#ffffff'
  surface-container-low: '#f3f3fe'
  surface-container: '#ededf9'
  surface-container-high: '#e7e7f3'
  surface-container-highest: '#e1e2ed'
  on-surface: '#191b23'
  on-surface-variant: '#434655'
  inverse-surface: '#2e3039'
  inverse-on-surface: '#f0f0fb'
  outline: '#737686'
  outline-variant: '#c3c6d7'
  surface-tint: '#0053db'
  primary: '#004ac6'
  on-primary: '#ffffff'
  primary-container: '#2563eb'
  on-primary-container: '#eeefff'
  inverse-primary: '#b4c5ff'
  secondary: '#565e74'
  on-secondary: '#ffffff'
  secondary-container: '#dae2fd'
  on-secondary-container: '#5c647a'
  tertiary: '#943700'
  on-tertiary: '#ffffff'
  tertiary-container: '#bc4800'
  on-tertiary-container: '#ffede6'
  error: '#ba1a1a'
  on-error: '#ffffff'
  error-container: '#ffdad6'
  on-error-container: '#93000a'
  primary-fixed: '#dbe1ff'
  primary-fixed-dim: '#b4c5ff'
  on-primary-fixed: '#00174b'
  on-primary-fixed-variant: '#003ea8'
  secondary-fixed: '#dae2fd'
  secondary-fixed-dim: '#bec6e0'
  on-secondary-fixed: '#131b2e'
  on-secondary-fixed-variant: '#3f465c'
  tertiary-fixed: '#ffdbcd'
  tertiary-fixed-dim: '#ffb596'
  on-tertiary-fixed: '#360f00'
  on-tertiary-fixed-variant: '#7d2d00'
  background: '#faf8ff'
  on-background: '#191b23'
  surface-variant: '#e1e2ed'
  admin-accent: '#3B82F6'
  storefront-cta: '#84CC16'
  status-success: '#22C55E'
  status-warning: '#F59E0B'
  status-error: '#EF4444'
  surface-muted: '#F8FAFC'
typography:
  display-lg:
    fontFamily: Plus Jakarta Sans
    fontSize: 48px
    fontWeight: '800'
    lineHeight: 56px
    letterSpacing: -0.02em
  headline-lg:
    fontFamily: Plus Jakarta Sans
    fontSize: 32px
    fontWeight: '700'
    lineHeight: 40px
  headline-lg-mobile:
    fontFamily: Plus Jakarta Sans
    fontSize: 24px
    fontWeight: '700'
    lineHeight: 32px
  headline-md:
    fontFamily: Plus Jakarta Sans
    fontSize: 24px
    fontWeight: '600'
    lineHeight: 32px
  body-lg:
    fontFamily: Inter
    fontSize: 18px
    fontWeight: '400'
    lineHeight: 28px
  body-md:
    fontFamily: Inter
    fontSize: 16px
    fontWeight: '400'
    lineHeight: 24px
  body-sm:
    fontFamily: Inter
    fontSize: 14px
    fontWeight: '400'
    lineHeight: 20px
  label-md:
    fontFamily: Inter
    fontSize: 14px
    fontWeight: '600'
    lineHeight: 16px
  label-sm:
    fontFamily: Inter
    fontSize: 12px
    fontWeight: '500'
    lineHeight: 16px
rounded:
  sm: 0.25rem
  DEFAULT: 0.5rem
  md: 0.75rem
  lg: 1rem
  xl: 1.5rem
  full: 9999px
spacing:
  container-max-width: 1440px
  sidebar-width: 260px
  gutter: 24px
  margin-mobile: 16px
  margin-desktop: 32px
  stack-sm: 8px
  stack-md: 16px
  stack-lg: 32px
---

## Brand & Style

The design system is engineered for a dual-natured ecosystem: the precision-driven **SaaS/Admin** environment and the high-energy **Storefront** experience. The brand personality is professional, reliable, and invisible—serving as a robust stage for diverse merchants while maintaining its own modern, architectural integrity.

The visual direction follows a **Corporate Modern** style for administrative tasks, transitioning into a **Minimalist / High-Contrast** aesthetic for the storefront.

- **SaaS/Admin:** Focuses on clarity, data density, and utilitarian elegance. It utilizes structured layouts and a calm color palette to reduce cognitive load during long work sessions.
- **Storefront:** Focuses on lifestyle, urgency, and product showcase. It uses larger typography, aggressive white space, and higher contrast to drive conversion.

## Colors

The palette distinguishes between functional utility and consumer engagement. 

- **Primary Blue (#2563EB):** The "Trust Blue," used primarily in the SaaS and Admin dashboards for primary actions, sidebar highlights, and data visualizations.
- **Secondary Navy (#0F172A):** Provides the structural weight for navigation and high-level headings.
- **Neutral Scale:** A sophisticated range of cool grays (Slate) ensures that the UI remains clean and professional.
- **Storefront Accents:** While the system defaults to high-contrast black/white for storefronts, the `storefront-cta` (#84CC16) is a vibrant lime green used to drive "Add to Cart" actions, inspired by high-end streetwear aesthetics.

## Typography

This design system utilizes a two-font pairing strategy. **Plus Jakarta Sans** is used for headlines to provide a friendly yet geometric and modern feel. **Inter** is the workhorse for all body text, data tables, and labels, ensuring maximum legibility across all screen sizes.

- **Storefronts** should lean heavily on `display-lg` for hero sections.
- **Dashboards** should prioritize `body-sm` and `label-sm` for data-heavy views to maintain information density.

## Layout & Spacing

The system employs a **Fluid Grid** for dashboards and a **Fixed/Centered Grid** for storefront content.

- **SaaS/Admin:** A classic sidebar-plus-content model. The sidebar remains fixed at 260px while the main content area utilizes a 12-column fluid grid.
- **Storefront:** A 12-column grid with a max-width of 1440px. Padding is significantly more generous in storefront views (64px+ between sections) compared to admin views (24px-32px).
- **Mobile:** All layouts collapse to a single column. The Seller Dashboard uses a bottom navigation bar or a condensed hamburger menu, while the Storefront prioritizes a "sticky" mobile header with search and cart icons.

## Elevation & Depth

Visual hierarchy is achieved through **Tonal Layers** and **Ambient Shadows**.

- **Admin/Dashboard:** Uses a light gray background (`surface-muted`) with white cards. Elevation is subtle, using low-opacity shadows (Blur: 10px, Y: 4px, Color: rgba(0,0,0,0.05)) to separate content modules without creating clutter.
- **Storefront:** Uses "High-Contrast Outlines" and flat surfaces for a trendy, streetwear vibe. Shadows are only used on primary floating elements like "Add to Cart" buttons or specific product hover states to create a "pop" effect.
- **Glassmorphism:** Reserved for mobile navigation overlays and top-level headers to maintain context of the background content while providing a modern, premium feel.

## Shapes

The design system uses a varied approach to roundedness to distinguish the "work" environment from the "shopping" environment.

- **Admin Components:** Use `rounded` (0.5rem) for cards and inputs to feel professional and stable.
- **Storefront Components:** Use `rounded-xl` (1.5rem) for primary buttons, product images, and category chips to evoke a softer, more modern, and consumer-friendly aesthetic.
- **Icons:** Should always follow a consistent 2px stroke weight with rounded terminals.

## Components

### Buttons
- **Admin:** Medium height (40px), `rounded` (0.5rem), Primary Blue background with white text.
- **Storefront:** Large height (56px), `rounded-xl` (1.5rem), High-contrast Black background or vibrant Lime Green for primary CTA.

### Cards
- **Admin:** White background, subtle 1px border (#E2E8F0), no shadow or very soft shadow. Content padding: 24px.
- **Storefront:** Borderless, soft shadow or tonal background. Large images with `rounded-xl` corners.

### Input Fields
- **Admin:** Outlined style, 1px gray border, turns Primary Blue on focus.
- **Storefront:** Larger, often with a subtle background fill instead of an outline, prioritizing ease of use on mobile devices.

### Data Tables (Admin Only)
- High density. Row borders are light (`#F1F5F9`). Uses `body-sm` for text. Header labels are `label-sm` with 50% opacity.

### Status Chips
- Small, `rounded-lg`, using the `status-` color tokens with 10% background opacity and 100% foreground color for the label.