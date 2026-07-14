"use client";

import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { productApi } from "@/lib/api";
import { useCartStore } from "@/stores/cart-store";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { ArrowLeft, ShoppingCart, Minus, Plus } from "lucide-react";

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

export default function ProductDetailPage({ params }: { params: { slug: string } }) {
  const router = useRouter();
  const addItem = useCartStore((s) => s.addItem);
  const [product, setProduct] = useState<Product | null>(null);
  const [loading, setLoading] = useState(true);
  const [quantity, setQuantity] = useState(1);
  const [selectedImage, setSelectedImage] = useState(0);

  useEffect(() => {
    loadProduct();
  }, [params.slug]);

  const loadProduct = async () => {
    try {
      setLoading(true);
      const response = await productApi.getBySlug(params.slug);
      setProduct(response.data);
    } catch (error) {
      console.error("Failed to load product:", error);
    } finally {
      setLoading(false);
    }
  };

  const handleAddToCart = () => {
    if (product) {
      addItem({
        id: product.id,
        product_id: product.id,
        name: product.name,
        price: product.price,
        slug: product.slug,
        image: product.images[0],
        quantity,
      });
      router.push("/cart");
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-neutral-100">
        <div className="container mx-auto px-4 py-8">
          <div className="grid md:grid-cols-2 gap-8">
            <div className="space-y-4">
              <div className="aspect-square bg-neutral-200 rounded-lg animate-pulse" />
              <div className="grid grid-cols-4 gap-2">
                {[...Array(4)].map((_, i) => (
                  <div key={i} className="aspect-square bg-neutral-200 rounded animate-pulse" />
                ))}
              </div>
            </div>
            <div className="space-y-6">
              <div className="h-8 bg-neutral-200 rounded animate-pulse" />
              <div className="h-4 bg-neutral-200 rounded w-1/3 animate-pulse" />
              <div className="space-y-2">
                <div className="h-3 bg-neutral-200 rounded animate-pulse" />
                <div className="h-3 bg-neutral-200 rounded animate-pulse" />
                <div className="h-3 bg-neutral-200 rounded w-2/3 animate-pulse" />
              </div>
            </div>
          </div>
        </div>
      </div>
    );
  }

  if (!product) {
    return (
      <div className="min-h-screen bg-neutral-100 flex items-center justify-center">
        <Card className="max-w-md w-full">
          <CardContent className="p-8 text-center">
            <h3 className="text-xl font-semibold text-neutral-900 mb-2">Product not found</h3>
            <p className="text-neutral-700 mb-4">
              The product you're looking for doesn't exist or has been removed.
            </p>
            <Button onClick={() => router.push("/products")}>Back to Products</Button>
          </CardContent>
        </Card>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-neutral-100">
      <div className="bg-white border-b border-neutral-300">
        <div className="container mx-auto px-4 py-4">
          <Button variant="ghost" onClick={() => router.back()} className="gap-2">
            <ArrowLeft className="h-4 w-4" />
            Back
          </Button>
        </div>
      </div>

      <div className="container mx-auto px-4 py-8">
        <div className="grid md:grid-cols-2 gap-8">
          {/* Images Section */}
          <div className="space-y-4">
            <div className="aspect-square bg-neutral-200 rounded-lg overflow-hidden">
              {product.images[selectedImage] ? (
                <img
                  src={product.images[selectedImage]}
                  alt={product.name}
                  className="w-full h-full object-cover"
                />
              ) : (
                <div className="w-full h-full flex items-center justify-center text-neutral-400">
                  No image available
                </div>
              )}
            </div>
            {product.images.length > 1 && (
              <div className="grid grid-cols-4 gap-2">
                {product.images.map((image, i) => (
                  <button
                    key={i}
                    onClick={() => setSelectedImage(i)}
                    className={`aspect-square rounded overflow-hidden border-2 transition-colors ${
                      i === selectedImage ? "border-brand-700" : "border-neutral-300"
                    }`}
                  >
                    <img src={image} alt={`${product.name} ${i + 1}`} className="w-full h-full object-cover" />
                  </button>
                ))}
              </div>
            )}
          </div>

          {/* Product Info Section */}
          <div className="space-y-6">
            <div>
              <h1 className="text-3xl font-bold text-neutral-900 mb-2">{product.name}</h1>
              <p className="text-4xl font-bold text-brand-700">
                {new Intl.NumberFormat("id-ID", {
                  style: "currency",
                  currency: product.currency || "IDR",
                }).format(product.price)}
              </p>
            </div>

            <div className="border-t border-neutral-300 pt-6">
              <h3 className="font-semibold text-neutral-900 mb-2">Description</h3>
              <p className="text-neutral-700 leading-relaxed">{product.description}</p>
            </div>

            <div className="border-t border-neutral-300 pt-6">
              <div className="flex items-center justify-between mb-2">
                <span className="font-semibold text-neutral-900">Availability</span>
                {product.stock > 0 ? (
                  <span className="text-success font-medium">{product.stock} in stock</span>
                ) : (
                  <span className="text-danger font-medium">Out of stock</span>
                )}
              </div>
            </div>

            {product.stock > 0 && (
              <div className="border-t border-neutral-300 pt-6">
                <div className="flex items-center gap-4 mb-4">
                  <span className="font-semibold text-neutral-900">Quantity</span>
                  <div className="flex items-center gap-2">
                    <Button
                      size="sm"
                      variant="outline"
                      onClick={() => setQuantity(Math.max(1, quantity - 1))}
                      disabled={quantity <= 1}
                    >
                      <Minus className="h-4 w-4" />
                    </Button>
                    <Input
                      type="number"
                      min="1"
                      max={product.stock}
                      value={quantity}
                      onChange={(e) => setQuantity(Math.max(1, Math.min(product.stock, parseInt(e.target.value) || 1)))}
                      className="w-20 text-center"
                    />
                    <Button
                      size="sm"
                      variant="outline"
                      onClick={() => setQuantity(Math.min(product.stock, quantity + 1))}
                      disabled={quantity >= product.stock}
                    >
                      <Plus className="h-4 w-4" />
                    </Button>
                  </div>
                </div>

                <Button onClick={handleAddToCart} className="w-full gap-2" size="lg">
                  <ShoppingCart className="h-5 w-5" />
                  Add to Cart
                </Button>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
