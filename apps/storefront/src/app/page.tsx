import Link from "next/link";
import {
  Store,
  ShoppingCart,
  CreditCard,
  Package,
  BarChart3,
  Shield,
  Cloud,
  Zap,
  ArrowRight,
  Check,
  Layers,
  Globe,
} from "lucide-react";
import { Navbar } from "@/components/layout/navbar";
import { Footer } from "@/components/layout/footer";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardHeader,
  CardTitle,
  CardDescription,
  CardContent,
} from "@/components/ui/card";
import { FadeIn } from "@/components/animations/fade-in";
import { Stagger, StaggerItem } from "@/components/animations/stagger";

const features = [
  {
    icon: Store,
    title: "Multi-Tenant Stores",
    description:
      "Spin up isolated storefronts in seconds. Each merchant gets their own data, branding, and domain.",
  },
  {
    icon: Package,
    title: "Product Catalog",
    description:
      "Full-featured catalog management with categories, variants, and bulk operations.",
  },
  {
    icon: ShoppingCart,
    title: "Shopping Cart",
    description:
      "Responsive cart with real-time inventory checks and seamless checkout flow.",
  },
  {
    icon: CreditCard,
    title: "Payment Integration",
    description:
      "Secure payment processing with support for multiple gateways and sandbox testing.",
  },
  {
    icon: BarChart3,
    title: "Analytics Dashboard",
    description:
      "Real-time insights into revenue, orders, and inventory with modern data visualization.",
  },
  {
    icon: Shield,
    title: "Enterprise Security",
    description:
      "JWT auth, tenant isolation, RBAC, and audit trails built-in from day one.",
  },
];

const architecture = [
  {
    icon: Layers,
    title: "Microservices",
    description: "Go services with Clean Architecture, database-per-service pattern.",
  },
  {
    icon: Cloud,
    title: "Cloud Native",
    description: "Docker, Kubernetes, and GitOps with ArgoCD for seamless deployments.",
  },
  {
    icon: Zap,
    title: "Event Driven",
    description: "NATS JetStream for async communication and saga orchestration.",
  },
  {
    icon: Globe,
    title: "Infrastructure as Code",
    description: "Terraform for provisioning, fully reproducible environments.",
  },
];

const pricingTiers = [
  {
    name: "Starter",
    price: "$0",
    period: "/month",
    description: "Perfect for trying out CloudCommerce",
    features: [
      "1 storefront",
      "Up to 50 products",
      "Basic analytics",
      "Community support",
    ],
    cta: "Get started",
    highlighted: false,
  },
  {
    name: "Growth",
    price: "$29",
    period: "/month",
    description: "For growing merchants who need more",
    features: [
      "Unlimited storefronts",
      "Unlimited products",
      "Advanced analytics",
      "Priority support",
      "Custom domain",
      "Payment gateway integration",
    ],
    cta: "Start free trial",
    highlighted: true,
  },
  {
    name: "Enterprise",
    price: "Custom",
    period: "",
    description: "For large-scale operations",
    features: [
      "Everything in Growth",
      "Dedicated infrastructure",
      "SLA guarantee",
      "White-label option",
      "Custom integrations",
      "24/7 support",
    ],
    cta: "Contact sales",
    highlighted: false,
  },
];

const faqs = [
  {
    question: "What is CloudCommerce?",
    answer:
      "CloudCommerce is a cloud-native, multi-tenant e-commerce platform built with microservices architecture. It provides merchants with isolated storefronts, scalable infrastructure, and developer-first tooling.",
  },
  {
    question: "How does multi-tenancy work?",
    answer:
      "Each merchant gets a logically isolated workspace with their own data partition (tenant ID), API rate limits, and storefront routing. Data isolation is enforced at the database and application layer.",
  },
  {
    question: "Can I self-host CloudCommerce?",
    answer:
      "Yes. CloudCommerce is designed to be cloud-agnostic. You can deploy it on GCP, AWS, or your own bare-metal infrastructure using our Terraform modules and Kubernetes manifests.",
  },
  {
    question: "What technologies are used?",
    answer:
      "The frontend uses Next.js 15 with TypeScript and Tailwind CSS. The backend is built with Go 1.25 microservices using Clean Architecture. PostgreSQL for storage, Redis for caching, and NATS for event-driven communication.",
  },
];

export default function LandingPage() {
  return (
    <div className="min-h-screen bg-background">
      <Navbar />

      {/* Hero Section */}
      <section className="relative overflow-hidden">
        <div className="absolute inset-0 bg-gradient-to-b from-primary/5 to-background" />
        <div className="relative mx-auto max-w-[1280px] px-4 py-20 sm:px-6 lg:px-8 lg:py-32">
          <div className="mx-auto max-w-3xl text-center">
            <FadeIn delay={0.1}>
              <div className="mb-4 inline-flex items-center gap-2 rounded-full border border-border bg-card px-4 py-1.5 text-sm text-muted-foreground">
                <span className="flex h-2 w-2 rounded-full bg-success animate-pulse" />
                Cloud Native Commerce Platform
              </div>
            </FadeIn>
            <FadeIn delay={0.2}>
              <h1 className="font-heading text-4xl font-bold tracking-tight text-foreground sm:text-5xl lg:text-6xl lg:leading-tight">
                Launch your store on a{" "}
                <span className="bg-gradient-to-r from-primary to-primary/80 bg-clip-text text-transparent">cloud-native</span>{" "}
                commerce platform
              </h1>
            </FadeIn>
            <FadeIn delay={0.3}>
              <p className="mt-8 text-lg leading-relaxed text-muted-foreground sm:text-xl">
                Multi-tenant e-commerce built with Go microservices, Next.js, and
                Kubernetes. Scalable, isolated, and developer-first.
              </p>
            </FadeIn>
            <FadeIn delay={0.4}>
              <div className="mt-10 flex flex-col items-center justify-center gap-4 sm:flex-row">
                <Link href="/register">
                  <Button size="lg" className="w-full sm:w-auto">
                    Start free
                    <ArrowRight className="ml-1 h-4 w-4" />
                  </Button>
                </Link>
                <Link href="/login">
                  <Button variant="outline" size="lg" className="w-full sm:w-auto">
                    Sign in
                  </Button>
                </Link>
              </div>
            </FadeIn>
          </div>
        </div>
      </section>

      {/* Features Section */}
      <section id="features" className="border-t border-border py-20">
        <div className="mx-auto max-w-[1280px] px-4 sm:px-6 lg:px-8">
          <FadeIn>
            <div className="mx-auto max-w-2xl text-center">
              <h2 className="font-heading text-3xl font-bold tracking-tight sm:text-4xl">
                Everything you need to sell online
              </h2>
              <p className="mt-4 text-lg leading-relaxed text-muted-foreground">
                From product catalog to payment processing, CloudCommerce has you
                covered with enterprise-grade features.
              </p>
            </div>
          </FadeIn>
          <Stagger className="mt-16 grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3" staggerDelay={0.1}>
            {features.map((feature) => {
              const Icon = feature.icon;
              return (
                <StaggerItem key={feature.title}>
                  <Card className="border-border transition-all duration-300 hover:shadow-lg hover:border-primary/50 h-full">
                    <CardHeader>
                      <div className="mb-3 flex h-12 w-12 items-center justify-center rounded-lg bg-primary/10 transition-colors group-hover:bg-primary/20">
                        <Icon className="h-6 w-6 text-primary" />
                      </div>
                      <CardTitle className="font-heading">{feature.title}</CardTitle>
                      <CardDescription className="leading-relaxed">{feature.description}</CardDescription>
                    </CardHeader>
                  </Card>
                </StaggerItem>
              );
            })}
          </Stagger>
        </div>
      </section>

      {/* Architecture Section */}
      <section id="architecture" className="border-t border-border bg-muted/30 py-20">
        <div className="mx-auto max-w-[1280px] px-4 sm:px-6 lg:px-8">
          <FadeIn>
            <div className="mx-auto max-w-2xl text-center">
              <h2 className="font-heading text-3xl font-bold tracking-tight sm:text-4xl">
                Built on modern architecture
              </h2>
              <p className="mt-4 text-lg leading-relaxed text-muted-foreground">
                Production-grade infrastructure designed for scale, reliability,
                and developer experience.
              </p>
            </div>
          </FadeIn>
          <Stagger className="mt-16 grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-4" staggerDelay={0.1}>
            {architecture.map((item) => {
              const Icon = item.icon;
              return (
                <StaggerItem key={item.title}>
                  <div className="group rounded-lg border border-border bg-card p-6 transition-all duration-300 hover:shadow-lg hover:border-primary/50 h-full">
                    <div className="mb-3 flex h-12 w-12 items-center justify-center rounded-lg bg-primary/10 transition-colors group-hover:bg-primary/20">
                      <Icon className="h-6 w-6 text-primary" />
                    </div>
                    <h3 className="font-heading font-semibold">{item.title}</h3>
                    <p className="mt-2 text-sm leading-relaxed text-muted-foreground">
                      {item.description}
                    </p>
                  </div>
                </StaggerItem>
              );
            })}
          </Stagger>
        </div>
      </section>

      {/* Pricing Section */}
      <section id="pricing" className="border-t border-border py-20">
        <div className="mx-auto max-w-[1280px] px-4 sm:px-6 lg:px-8">
          <FadeIn>
            <div className="mx-auto max-w-2xl text-center">
              <h2 className="font-heading text-3xl font-bold tracking-tight sm:text-4xl">
                Simple, transparent pricing
              </h2>
              <p className="mt-4 text-lg leading-relaxed text-muted-foreground">
                Start free, upgrade as you grow. No hidden fees.
              </p>
            </div>
          </FadeIn>
          <Stagger className="mt-16 grid grid-cols-1 gap-6 lg:grid-cols-3" staggerDelay={0.15}>
            {pricingTiers.map((tier) => (
              <StaggerItem key={tier.name}>
                <Card
                  className={
                    tier.highlighted
                      ? "border-primary shadow-lg scale-105 transition-all duration-300 hover:shadow-xl h-full"
                      : "border-border transition-all duration-300 hover:shadow-lg hover:border-primary/50 h-full"
                  }
                >
                  <CardHeader>
                  <div className="flex items-center justify-between">
                    <CardTitle className="font-heading">{tier.name}</CardTitle>
                    {tier.highlighted && (
                      <span className="rounded-full bg-primary px-3 py-1 text-xs font-semibold text-primary-foreground">
                        Popular
                      </span>
                    )}
                  </div>
                  <CardDescription className="leading-relaxed">{tier.description}</CardDescription>
                  <div className="mt-6">
                    <span className="font-heading text-4xl font-bold">{tier.price}</span>
                    <span className="text-muted-foreground">{tier.period}</span>
                  </div>
                </CardHeader>
                <CardContent>
                  <ul className="space-y-3">
                    {tier.features.map((feature) => (
                      <li
                        key={feature}
                        className="flex items-center gap-2 text-sm"
                      >
                        <Check className="h-4 w-4 text-success" />
                        <span>{feature}</span>
                      </li>
                    ))}
                  </ul>
                  <Link href="/register" className="mt-6 block">
                    <Button
                      variant={tier.highlighted ? "primary" : "outline"}
                      className="w-full"
                    >
                      {tier.cta}
                    </Button>
                  </Link>
                </CardContent>
              </Card>
            </StaggerItem>
            ))}
          </Stagger>
        </div>
      </section>

      {/* FAQ Section */}
      <section id="faq" className="border-t border-border bg-muted/30 py-20">
        <div className="mx-auto max-w-3xl px-4 sm:px-6 lg:px-8">
          <FadeIn>
            <div className="text-center">
              <h2 className="font-heading text-3xl font-bold tracking-tight sm:text-4xl">
                Frequently asked questions
              </h2>
            </div>
          </FadeIn>
          <Stagger className="mt-12 space-y-6" staggerDelay={0.1}>
            {faqs.map((faq) => (
              <StaggerItem key={faq.question}>
                <div className="group rounded-lg border border-border bg-card p-6 transition-all duration-300 hover:shadow-md hover:border-primary/50">
                  <h3 className="font-heading font-semibold">{faq.question}</h3>
                  <p className="mt-3 text-sm leading-relaxed text-muted-foreground">
                    {faq.answer}
                  </p>
                </div>
              </StaggerItem>
            ))}
          </Stagger>
        </div>
      </section>

      {/* CTA Section */}
      <section className="border-t border-border py-20">
        <div className="mx-auto max-w-[1280px] px-4 sm:px-6 lg:px-8">
          <FadeIn delay={0.2}>
            <div className="rounded-xl border border-border bg-gradient-to-br from-primary/10 via-primary/5 to-background p-12 text-center shadow-lg">
              <h2 className="font-heading text-3xl font-bold tracking-tight sm:text-4xl">
                Ready to launch your store?
              </h2>
              <p className="mt-4 text-lg leading-relaxed text-muted-foreground">
                Join CloudCommerce today and start selling in minutes.
              </p>
              <div className="mt-8 flex flex-col items-center justify-center gap-4 sm:flex-row">
                <Link href="/register">
                  <Button size="lg" className="shadow-md hover:shadow-lg transition-shadow">
                    Create your store
                    <ArrowRight className="ml-2 h-4 w-4" />
                  </Button>
                </Link>
                <Link href="/login">
                  <Button variant="outline" size="lg" className="hover:bg-accent transition-colors">
                    Sign in
                  </Button>
                </Link>
              </div>
            </div>
          </FadeIn>
        </div>
      </section>

      <Footer />
    </div>
  );
}
