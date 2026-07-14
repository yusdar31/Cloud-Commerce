import Link from "next/link";
import { Store } from "lucide-react";

export default function AuthLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="flex min-h-screen flex-col">
      <div className="flex flex-1">
        {/* Left panel - branding */}
        <div className="hidden w-1/2 flex-col justify-between bg-gradient-to-br from-primary to-primary/80 p-12 lg:flex">
          <Link href="/" className="flex items-center gap-2">
            <Store className="h-7 w-7 text-white" />
            <span className="text-xl font-semibold text-white">
              CloudCommerce
            </span>
          </Link>
          <div className="space-y-6 text-white">
            <h2 className="text-3xl font-bold tracking-tight">
              Build your commerce empire on the cloud
            </h2>
            <p className="text-lg text-white/80">
              Multi-tenant e-commerce platform with enterprise-grade security,
              scalable microservices, and developer-first tooling.
            </p>
            <div className="space-y-3">
              {[
                "Launch isolated storefronts in seconds",
                "Real-time analytics and insights",
                "Cloud-native, Kubernetes-ready infrastructure",
                "Event-driven architecture with NATS",
              ].map((item) => (
                <div key={item} className="flex items-center gap-2 text-white/90">
                  <svg
                    className="h-5 w-5 shrink-0"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                    strokeWidth={2}
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      d="M5 13l4 4L19 7"
                    />
                  </svg>
                  <span>{item}</span>
                </div>
              ))}
            </div>
          </div>
          <p className="text-sm text-white/60">
            &copy; {new Date().getFullYear()} CloudCommerce. All rights reserved.
          </p>
        </div>

        {/* Right panel - form */}
        <div className="flex w-full flex-col justify-center px-4 py-12 sm:px-6 lg:w-1/2 lg:px-12">
          <div className="mx-auto w-full max-w-sm">
            <Link
              href="/"
              className="mb-8 flex items-center gap-2 lg:hidden"
            >
              <Store className="h-6 w-6 text-primary" />
              <span className="text-lg font-semibold">CloudCommerce</span>
            </Link>
            {children}
          </div>
        </div>
      </div>
    </div>
  );
}
