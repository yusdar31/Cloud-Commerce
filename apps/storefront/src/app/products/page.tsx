"use client";

import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { productApi } from "@/lib/api";
import { useCartStore } from "@/stores/cart-store";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Header } from "@/components/layout/header";
import { Search, ShoppingCart } from "lucide-react";

interface Product {
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
}

export default function ProductsPage() {
  const router = useRouter();
  const addItem = useCartStore((s) => s.addItem);
  const [products, setProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState(true);
  const [search, setSearch] = useState("");

  useEffect(() => {
    loadProducts();
  }, []);

  const loadProducts = async () => {
    try {
      setLoading(true);
      const response = await productApi.list({ limit: 20 });
      setProducts(response.data);
    } catch (error) {
      console.error("Failed to load products:", error);
    } finally {
      setLoading(false);
    }
  };

  const handleAddToCart = (product: Product) => {
    addItem({
      id: product.id,
      product_id: product.id,
      name: product.name,
      price: product.price,
      slug: product.slug,
      image: product.images[0],
    });
  };

  const filteredProducts = products.filter((p) =>
    p.name.toLowerCase().includes(search.toLowerCase())
  );

  if (loading) {
    return (
      <div className="min-h-screen bg-neutral-100">
        <div className="container mx-auto px-4 py-8">
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
            {[...Array(8)].map((_, i) => (
              <Card key={i} className="overflow-hidden">
                <div className="aspect-square bg-neutral-200 animate-pulse" />
                <CardContent className="p-4 space-y-3">
                  <div className="h-4 bg-neutral-200 rounded animate-pulse" />
                  <div className="h-3 bg-neutral-200 rounded w-2/3 animate-pulse" />
                  <div className="h-8 bg-neutral-200 rounded animate-pulse" />
                </CardContent>
              </Card>
            ))}
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-neutral-100">
      <Header />
      {/* Header */}
      <div className="bg-white border-b border-neutral-300">
        <div className="container mx-auto px-4 py-6">
          <h1 className="text-3xl font-bold text-neutral-900 mb-4">Products</h1>
          <div className="relative max-w-md">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-5 w-5 text-neutral-500" />
            <Input
              type="text"
              placeholder="Search products..."
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              className="pl-10"
            />
          </div>
        </div>
      </div>

      {/* Products Grid */}
      <div className="container mx-auto px-4 py-8">
        {filteredProducts.length === 0 ? (
          <div className="text-center py-20">
            <div className="inline-flex items-center justify-center w-16 h-16 bg-neutral-200 rounded-full mb-4">
              <Search className="h-8 w-8 text-neutral-500" />
            </div>
            <h3 className="text-xl font-semibold text-neutral-900 mb-2">
              No products found
            </h3>
            <p className="text-neutral-700 mb-4">
              Try adjusting your search to find what you're looking for
            </p>
            <Button onClick={() => setSearch("")} variant="outline">
              Clear search
            </Button>
          </div>
        ) : (
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
            {filteredProducts.map((product) => (
              <Card
                key={product.id}
                className="overflow-hidden hover:shadow-lg transition-shadow cursor-pointer group"
              >
                <div
                  onClick={() => router.push(`/products/${product.slug}`)}
                  className="aspect-square bg-neutral-200 relative overflow-hidden"
                >
                  {product.images[0] ? (
                    <img
                      src={product.images[0]}
                      alt={product.name}
                      className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
                    />
                  ) : (
                    <div className="w-full h-full flex items-center justify-center text-neutral-400">
                      No image
                    </div>
                  )}
                  {product.stock <= 0 && (
                    <div className="absolute inset-0 bg-black/50 flex items-center justify-center">
                      <span className="text-white font-semibold">Out of Stock</span>
                    </div>
                  )}
                </div>

                <CardContent className="p-4 space-y-3">
                  <div
                    onClick={() => router.push(`/products/${product.slug}`)}
                    className="space-y-1"
                  >
                    <h3 className="font-semibold text-neutral-900 line-clamp-2 group-hover:text-brand-700 transition-colors">
                      {product.name}
                    </h3>
                    <p className="text-sm text-neutral-700 line-clamp-2">
                      {product.description}
                    </p>
                  </div>

                  <div className="flex items-center justify-between pt-2">
                    <div>
                      <p className="text-2xl font-bold text-brand-700">
                        {new Intl.NumberFormat("id-ID", {
                          style: "currency",
                          currency: product.currency || "IDR",
                        }).format(product.price)}
                      </p>
                      {product.stock > 0 && product.stock < 10 && (
                        <p className="text-xs text-warning">
                          Only {product.stock} left
                        </p>
                      )}
                    </div>

                    <Button
                      size="sm"
                      onClick={(e) => {
                        e.stopPropagation();
                        handleAddToCart(product);
                      }}
                      disabled={product.stock <= 0}
                      className="gap-2"
                    >
                      <ShoppingCart className="h-4 w-4" />
                      Add
                    </Button>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
