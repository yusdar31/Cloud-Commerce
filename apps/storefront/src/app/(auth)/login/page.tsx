import Link from "next/link";
import { LoginForm } from "@/features/auth/login-form";

export const metadata = {
  title: "Sign in",
  description: "Sign in to your CloudCommerce account",
};

export default function LoginPage() {
  return (
    <div>
      <div className="mb-8">
        <h1 className="text-2xl font-bold tracking-tight">Welcome back</h1>
        <p className="mt-2 text-sm text-muted-foreground">
          Enter your credentials to access your dashboard
        </p>
      </div>

      <LoginForm />

      <div className="mt-6 text-center text-sm text-muted-foreground">
        Don&apos;t have an account?{" "}
        <Link
          href="/register"
          className="font-medium text-primary hover:underline"
        >
          Sign up
        </Link>
      </div>
    </div>
  );
}
