"use client";

import Link from "next/link";
import { useRouter, usePathname } from "next/navigation";
import { useCartStore } from "@/stores/cart-store";
import { useAuthStore } from "@/stores/auth-store";
import { Button } from "@/components/ui/button";
import { ShoppingCart, User, LogOut, Package } from "lucide-react";

export function Header() {
  const router = useRouter();
  const pathname = usePathname();
  const getItemCount = useCartStore((s) => s.getItemCount);
  const { user, isAuthenticated, clearAuth } = useAuthStore();
  const itemCount = getItemCount();

  // Hide header on landing page and auth pages
  if (pathname === "/" || pathname === "/login" || pathname === "/register") {
    return null;
  }

  const handleLogout = () => {
    clearAuth();
    router.push("/");
  };

  const isActive = (path: string) => pathname === path;

  return (
    <header className="sticky top-0 z-50 w-full border-b border-neutral-300 bg-white">
      <div className="container mx-auto px-4">
        <div className="flex h-16 items-center justify-between">
          {/* Logo */}
          <Link href="/" className="flex items-center gap-2">
            <div className="h-8 w-8 rounded bg-brand-700 flex items-center justify-center">
              <span className="text-white font-bold text-sm">CC</span>
            </div>
            <span className="font-bold text-xl text-neutral-900">CloudCommerce</span>
          </Link>

          {/* Navigation Links */}
          <nav className="hidden md:flex items-center gap-6">
            <Link
              href="/products"
              className={`text-sm font-medium transition-colors hover:text-brand-700 ${
                isActive("/products") ? "text-brand-700" : "text-neutral-700"
              }`}
            >
              Products
            </Link>
            {isAuthenticated && (
              <Link
                href="/dashboard"
                className={`text-sm font-medium transition-colors hover:text-brand-700 ${
                  isActive("/dashboard") ? "text-brand-700" : "text-neutral-700"
                }`}
              >
                Dashboard
              </Link>
            )}
          </nav>

          {/* Right Side Actions */}
          <div className="flex items-center gap-3">
            {/* Cart Button */}
            <Button
              variant="ghost"
              size="sm"
              onClick={() => router.push("/cart")}
              className="relative gap-2"
            >
              <ShoppingCart className="h-5 w-5" />
              {itemCount > 0 && (
                <span className="absolute -top-1 -right-1 h-5 w-5 rounded-full bg-brand-700 text-white text-xs flex items-center justify-center">
                  {itemCount > 9 ? "9+" : itemCount}
                </span>
              )}
              <span className="hidden sm:inline">Cart</span>
            </Button>

            {/* User Menu */}
            {isAuthenticated ? (
              <div className="flex items-center gap-2">
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={() => router.push("/dashboard")}
                  className="gap-2"
                >
                  <User className="h-5 w-5" />
                  <span className="hidden sm:inline">{user?.full_name || "Account"}</span>
                </Button>
                <Button variant="ghost" size="sm" onClick={handleLogout} className="gap-2">
                  <LogOut className="h-4 w-4" />
                  <span className="hidden sm:inline">Logout</span>
                </Button>
              </div>
            ) : (
              <div className="flex items-center gap-2">
                <Button variant="ghost" size="sm" onClick={() => router.push("/login")}>
                  Login
                </Button>
                <Button size="sm" onClick={() => router.push("/register")}>
                  Sign Up
                </Button>
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Mobile Navigation */}
      <div className="md:hidden border-t border-neutral-300 bg-white">
        <div className="container mx-auto px-4 py-2 flex gap-4">
          <Link
            href="/products"
            className={`text-sm font-medium transition-colors hover:text-brand-700 ${
              isActive("/products") ? "text-brand-700" : "text-neutral-700"
            }`}
          >
            Products
          </Link>
          {isAuthenticated && (
            <Link
              href="/dashboard"
              className={`text-sm font-medium transition-colors hover:text-brand-700 ${
                isActive("/dashboard") ? "text-brand-700" : "text-neutral-700"
              }`}
            >
              Dashboard
            </Link>
          )}
        </div>
      </div>
    </header>
  );
}
