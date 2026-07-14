import Link from "next/link";
import { RegisterForm } from "@/features/auth/register-form";

export const metadata = {
  title: "Create account",
  description: "Create your CloudCommerce account",
};

export default function RegisterPage() {
  return (
    <div>
      <div className="mb-8">
        <h1 className="text-2xl font-bold tracking-tight">Create your account</h1>
        <p className="mt-2 text-sm text-muted-foreground">
          Start your commerce journey with CloudCommerce
        </p>
      </div>

      <RegisterForm />

      <div className="mt-6 text-center text-sm text-muted-foreground">
        Already have an account?{" "}
        <Link
          href="/login"
          className="font-medium text-primary hover:underline"
        >
          Sign in
        </Link>
      </div>

      <p className="mt-6 text-center text-xs text-muted-foreground">
        By creating an account, you agree to our{" "}
        <a href="#" className="hover:underline">
          Terms of Service
        </a>{" "}
        and{" "}
        <a href="#" className="hover:underline">
          Privacy Policy
        </a>
        .
      </p>
    </div>
  );
}
