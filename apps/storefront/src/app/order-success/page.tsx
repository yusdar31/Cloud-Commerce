"use client";

import { useRouter } from "next/navigation";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Header } from "@/components/layout/header";
import { CheckCircle } from "lucide-react";

export default function OrderSuccessPage() {
  const router = useRouter();

  return (
    <div className="min-h-screen bg-neutral-100 flex items-center justify-center p-4">
      <Header />
      <Card className="max-w-lg w-full">
        <CardContent className="p-12 text-center">
          <div className="inline-flex items-center justify-center w-20 h-20 bg-success/10 rounded-full mb-6">
            <CheckCircle className="h-10 w-10 text-success" />
          </div>
          
          <h1 className="text-3xl font-bold text-neutral-900 mb-3">
            Order Placed Successfully!
          </h1>
          
          <p className="text-neutral-700 mb-2">
            Thank you for your order. We've received your payment and will process your order shortly.
          </p>
          
          <p className="text-sm text-neutral-500 mb-8">
            You will receive an email confirmation with your order details.
          </p>

          <div className="flex flex-col sm:flex-row gap-3 justify-center">
            <Button onClick={() => router.push("/products")} variant="outline" size="lg">
              Continue Shopping
            </Button>
            <Button onClick={() => router.push("/dashboard")} size="lg">
              View Orders
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
