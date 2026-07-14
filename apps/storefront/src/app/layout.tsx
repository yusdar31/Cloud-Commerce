import type { Metadata } from "next";
import { Poppins, Open_Sans } from "next/font/google";
import { QueryProvider } from "@/components/providers/query-provider";
import "./globals.css";

const poppins = Poppins({
  subsets: ["latin"],
  weight: ["400", "500", "600", "700"],
  variable: "--font-poppins",
  display: "swap",
});

const openSans = Open_Sans({
  subsets: ["latin"],
  weight: ["300", "400", "500", "600", "700"],
  variable: "--font-open-sans",
  display: "swap",
});

export const metadata: Metadata = {
  title: {
    default: "CloudCommerce - Cloud Native Commerce Platform",
    template: "%s | CloudCommerce",
  },
  description:
    "CloudCommerce is a cloud-native commerce platform for modern merchants. Launch your store in minutes with multi-tenant isolation, scalable infrastructure, and developer-first tools.",
  keywords: ["ecommerce", "saas", "multi-tenant", "cloud", "commerce"],
  authors: [{ name: "CloudCommerce" }],
  openGraph: {
    title: "CloudCommerce - Cloud Native Commerce Platform",
    description:
      "Launch your store in minutes with multi-tenant isolation and scalable infrastructure.",
    type: "website",
  },
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" className={`${poppins.variable} ${openSans.variable}`} suppressHydrationWarning>
      <body className="font-body antialiased">
        <QueryProvider>{children}</QueryProvider>
      </body>
    </html>
  );
}
