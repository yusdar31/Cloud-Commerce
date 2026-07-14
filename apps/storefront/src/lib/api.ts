import ky, { type KyInstance } from "ky";

const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

export const api: KyInstance = ky.create({
  prefixUrl: API_URL,
  timeout: 30000,
  hooks: {
    beforeRequest: [
      (request) => {
        if (typeof window !== "undefined") {
          const token = localStorage.getItem("access_token");
          if (token) {
            request.headers.set("Authorization", `Bearer ${token}`);
          }
          const tenantId = localStorage.getItem("tenant_id");
          if (tenantId) {
            request.headers.set("X-Tenant-ID", tenantId);
          }
        }
      },
    ],
  },
});

export const authApi = {
  register: async (data: {
    email: string;
    password: string;
    full_name: string;
    phone?: string;
  }) => {
    return api.post("api/v1/auth/register", { json: data }).json<{
      data: {
        user: {
          id: string;
          email: string;
          full_name: string;
          role: string;
        };
        tokens: {
          access_token: string;
          refresh_token: string;
          expires_at: string;
        };
      };
    }>();
  },

  login: async (data: { email: string; password: string }) => {
    return api.post("api/v1/auth/login", { json: data }).json<{
      data: {
        user: {
          id: string;
          email: string;
          full_name: string;
          role: string;
        };
        tokens: {
          access_token: string;
          refresh_token: string;
          expires_at: string;
        };
      };
    }>();
  },

  refresh: async (data: { refresh_token: string }) => {
    return api.post("api/v1/auth/refresh", { json: data }).json<{
      data: {
        access_token: string;
        refresh_token: string;
        expires_at: string;
      };
    }>();
  },

  logout: async () => {
    return api.post("api/v1/auth/logout").json();
  },

  getProfile: async () => {
    return api.get("api/v1/users/me").json<{
      data: {
        id: string;
        email: string;
        full_name: string;
        phone: string;
        role: string;
        tenant_id: string;
        is_active: boolean;
      };
    }>();
  },
};

export const productApi = {
  list: async (params?: { page?: number; limit?: number; search?: string }) => {
    const searchParams = new URLSearchParams();
    if (params?.page) searchParams.set("page", params.page.toString());
    if (params?.limit) searchParams.set("limit", params.limit.toString());
    if (params?.search) searchParams.set("search", params.search);
    
    return api.get(`api/v1/products?${searchParams}`).json<{
      data: Array<{
        id: string;
        name: string;
        slug: string;
        description: string;
        price: number;
        currency: string;
        stock: number;
        status: string;
        images: string[];
        created_at: string;
      }>;
    }>();
  },

  getBySlug: async (slug: string) => {
    return api.get(`api/v1/products/slug/${slug}`).json<{
      data: {
        id: string;
        name: string;
        slug: string;
        description: string;
        price: number;
        currency: string;
        stock: number;
        status: string;
        images: string[];
        created_at: string;
      };
    }>();
  },
};
