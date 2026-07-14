export interface User {
  id: string;
  email: string;
  full_name: string;
  phone: string;
  role: UserRole;
  tenant_id: string;
  is_active: boolean;
}

export type UserRole = "seller" | "buyer" | "admin";

export interface AuthTokens {
  access_token: string;
  refresh_token: string;
  expires_at: string;
}

export interface AuthResponse {
  user: User;
  tokens: AuthTokens;
}

export interface Product {
  id: string;
  tenant_id: string;
  name: string;
  slug: string;
  description: string;
  sku: string;
  price: number;
  currency: string;
  status: ProductStatus;
  category_id: string;
  image_url: string;
  weight: number;
}

export type ProductStatus = "draft" | "published" | "archived";

export interface Category {
  id: string;
  tenant_id: string;
  name: string;
  slug: string;
  description: string;
}

export interface CartItem {
  product_id: string;
  name: string;
  price: number;
  quantity: number;
  image_url: string;
}

export interface ApiResponse<T> {
  success: boolean;
  data: T;
  meta?: {
    page: number;
    per_page: number;
    total: number;
  };
}
