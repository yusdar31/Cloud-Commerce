"use client";

import React, { useEffect } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useAuthStore } from "@/stores/auth-store";
import { LayoutDashboard, ShoppingBag, Settings, LogOut, User as UserIcon, HelpCircle } from "lucide-react";
import { Button } from "@/components/ui/button";

export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const router = useRouter();
  const { user, isAuthenticated, clearAuth } = useAuthStore();

  useEffect(() => {
    // Redirect if not authenticated (as fallback to middleware)
    if (!isAuthenticated) {
      router.push("/login");
    }
  }, [isAuthenticated, router]);

  const handleLogout = () => {
    clearAuth();
    router.push("/login");
  };

  if (!isAuthenticated || !user) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-background">
        <div className="h-8 w-8 animate-spin rounded-full border-4 border-primary border-t-transparent" />
      </div>
    );
  }

  return (
    <div className="flex min-h-screen bg-muted/20">
      {/* Sidebar */}
      <aside className="fixed inset-y-0 left-0 z-20 flex w-64 flex-col border-r border-border bg-card px-4 py-6">
        <div className="flex items-center gap-2 px-2">
          <div className="flex h-9 w-9 items-center justify-center rounded-lg bg-primary text-primary-foreground font-bold">
            CC
          </div>
          <span className="font-heading text-lg font-bold tracking-tight text-foreground">
            CloudCommerce
          </span>
        </div>

        {/* Navigation */}
        <nav className="mt-8 flex-1 space-y-1">
          <Link
            href="/dashboard"
            className="flex items-center gap-3 rounded-lg bg-primary/10 px-3 py-2 text-sm font-medium text-primary"
          >
            <LayoutDashboard className="h-4 w-4" />
            Dashboard
          </Link>
          <Link
            href="#"
            className="flex items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium text-muted-foreground transition-colors hover:bg-muted hover:text-foreground"
          >
            <ShoppingBag className="h-4 w-4" />
            Products
          </Link>
          <Link
            href="#"
            className="flex items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium text-muted-foreground transition-colors hover:bg-muted hover:text-foreground"
          >
            <Settings className="h-4 w-4" />
            Settings
          </Link>
        </nav>

        {/* Footer info & Logout */}
        <div className="border-t border-border pt-4">
          <div className="flex items-center gap-3 px-2 py-1.5 mb-4">
            <div className="flex h-9 w-9 items-center justify-center rounded-full bg-muted">
              <UserIcon className="h-4 w-4 text-muted-foreground" />
            </div>
            <div className="flex flex-col min-w-0">
              <span className="truncate text-sm font-medium text-foreground">
                {user.full_name}
              </span>
              <span className="truncate text-xs text-muted-foreground">
                {user.email}
              </span>
            </div>
          </div>
          <Button
            variant="outline"
            className="w-full justify-start gap-3 border-destructive/20 text-destructive hover:bg-destructive/5 hover:text-destructive"
            onClick={handleLogout}
          >
            <LogOut className="h-4 w-4" />
            Sign Out
          </Button>
        </div>
      </aside>

      {/* Main Content */}
      <div className="flex flex-1 flex-col pl-64">
        {/* Header */}
        <header className="sticky top-0 z-10 flex h-16 items-center justify-between border-b border-border bg-card/85 px-8 backdrop-blur-md">
          <div className="flex items-center gap-1.5 text-sm text-muted-foreground">
            <span>Console</span>
            <span>/</span>
            <span className="font-medium text-foreground">Dashboard</span>
          </div>

          <div className="flex items-center gap-4">
            <span className="rounded-full bg-success/15 px-2.5 py-1 text-xs font-semibold text-success capitalize">
              Role: {user.role}
            </span>
          </div>
        </header>

        {/* Main Workspace */}
        <main className="flex-1 p-8">{children}</main>
      </div>
    </div>
  );
}
