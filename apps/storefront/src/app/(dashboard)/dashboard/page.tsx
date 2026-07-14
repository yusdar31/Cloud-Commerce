"use client";

import React from "react";
import { useAuthStore } from "@/stores/auth-store";
import { 
  TrendingUp, 
  ShoppingBag, 
  Users, 
  DollarSign, 
  ArrowUpRight,
  ArrowDownRight,
  PackageCheck
} from "lucide-react";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";

export default function DashboardPage() {
  const { user } = useAuthStore();

  const stats = [
    {
      title: "Total Revenue",
      value: "$14,235.50",
      change: "+12.5%",
      isPositive: true,
      description: "from last month",
      icon: DollarSign,
    },
    {
      title: "Orders",
      value: "350",
      change: "+8.2%",
      isPositive: true,
      description: "from last month",
      icon: ShoppingBag,
    },
    {
      title: "Active Customers",
      value: "1,205",
      change: "+15.3%",
      isPositive: true,
      description: "from last month",
      icon: Users,
    },
    {
      title: "Fulfillment Rate",
      value: "98.2%",
      change: "-0.5%",
      isPositive: false,
      description: "from last month",
      icon: PackageCheck,
    },
  ];

  const recentOrders = [
    { id: "ORD-001", customer: "John Doe", status: "completed", date: "2026-07-14", amount: "$150.00" },
    { id: "ORD-002", customer: "Sarah Smith", status: "processing", date: "2026-07-14", amount: "$89.90" },
    { id: "ORD-003", customer: "Michael Brown", status: "completed", date: "2026-07-13", amount: "$245.00" },
    { id: "ORD-004", customer: "Emily Davis", status: "failed", date: "2026-07-12", amount: "$12.50" },
    { id: "ORD-005", customer: "David Wilson", status: "completed", date: "2026-07-11", amount: "$105.00" },
  ];

  return (
    <div className="space-y-8 animate-fade-in">
      {/* Welcome Banner */}
      <div>
        <h1 className="font-heading text-3xl font-bold tracking-tight text-foreground">
          Welcome back, {user?.full_name}!
        </h1>
        <p className="text-muted-foreground mt-1">
          Here is what is happening with your commerce store today.
        </p>
      </div>

      {/* Stats Grid */}
      <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-4">
        {stats.map((stat) => {
          const Icon = stat.icon;
          return (
            <Card key={stat.title} className="border-border bg-card">
              <CardHeader className="flex flex-row items-center justify-between pb-2">
                <CardTitle className="text-sm font-medium text-muted-foreground">
                  {stat.title}
                </CardTitle>
                <div className="rounded-lg bg-primary/10 p-2">
                  <Icon className="h-4 w-4 text-primary" />
                </div>
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{stat.value}</div>
                <div className="flex items-center gap-1 mt-1 text-xs">
                  {stat.isPositive ? (
                    <span className="flex items-center gap-0.5 font-medium text-success">
                      <ArrowUpRight className="h-3 w-3" />
                      {stat.change}
                    </span>
                  ) : (
                    <span className="flex items-center gap-0.5 font-medium text-destructive">
                      <ArrowDownRight className="h-3 w-3" />
                      {stat.change}
                    </span>
                  )}
                  <span className="text-muted-foreground">{stat.description}</span>
                </div>
              </CardContent>
            </Card>
          );
        })}
      </div>

      {/* Details Section */}
      <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
        {/* Recent Orders Card */}
        <Card className="col-span-2 border-border bg-card">
          <CardHeader className="flex flex-row items-center justify-between">
            <div>
              <CardTitle className="font-heading">Recent Orders</CardTitle>
              <CardDescription>Latest order transactions from your customers.</CardDescription>
            </div>
          </CardHeader>
          <CardContent>
            <div className="relative w-full overflow-auto">
              <table className="w-full caption-bottom text-sm">
                <thead>
                  <tr className="border-b border-border transition-colors hover:bg-muted/50 data-[state=selected]:bg-muted">
                    <th className="h-12 px-4 text-left align-middle font-medium text-muted-foreground">Order ID</th>
                    <th className="h-12 px-4 text-left align-middle font-medium text-muted-foreground">Customer</th>
                    <th className="h-12 px-4 text-left align-middle font-medium text-muted-foreground">Status</th>
                    <th className="h-12 px-4 text-left align-middle font-medium text-muted-foreground">Date</th>
                    <th className="h-12 px-4 text-right align-middle font-medium text-muted-foreground">Amount</th>
                  </tr>
                </thead>
                <tbody className="[&_tr:last-child]:border-0">
                  {recentOrders.map((order) => (
                    <tr key={order.id} className="border-b border-border transition-colors hover:bg-muted/50">
                      <td className="p-4 align-middle font-medium">{order.id}</td>
                      <td className="p-4 align-middle">{order.customer}</td>
                      <td className="p-4 align-middle">
                        <span className={`inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-semibold capitalize ${
                          order.status === "completed" 
                            ? "bg-success/15 text-success" 
                            : order.status === "processing" 
                            ? "bg-warning/15 text-warning" 
                            : "bg-destructive/15 text-destructive"
                        }`}>
                          {order.status}
                        </span>
                      </td>
                      <td className="p-4 align-middle text-muted-foreground">{order.date}</td>
                      <td className="p-4 align-middle text-right font-medium">{order.amount}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </CardContent>
        </Card>

        {/* Store Isolation Card */}
        <Card className="border-border bg-card">
          <CardHeader>
            <CardTitle className="font-heading">Tenant Metadata</CardTitle>
            <CardDescription>Details of your multi-tenant workspace isolation.</CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="flex flex-col gap-1 rounded-lg bg-muted/40 p-4">
              <span className="text-xs text-muted-foreground uppercase font-bold tracking-wider">Tenant ID</span>
              <span className="font-mono text-xs break-all text-foreground">
                {user?.tenant_id || "d1be6a24-c15b-426b-87a3-e847cd1192fa"}
              </span>
            </div>
            <div className="flex flex-col gap-1 rounded-lg bg-muted/40 p-4">
              <span className="text-xs text-muted-foreground uppercase font-bold tracking-wider">User ID</span>
              <span className="font-mono text-xs break-all text-foreground">{user?.id}</span>
            </div>
            <div className="flex flex-col gap-1 rounded-lg bg-muted/40 p-4">
              <span className="text-xs text-muted-foreground uppercase font-bold tracking-wider">Session Token Status</span>
              <span className="text-xs text-success flex items-center gap-1.5 font-medium">
                <span className="h-2 w-2 rounded-full bg-success" />
                Active (Validated by API Gateway)
              </span>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
