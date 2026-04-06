import "./globals.css";
import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "FeedPulse_Cloud_Native_Platform",
  description: "AI-powered product feedback platform",
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en" data-theme="light">
      <body>{children}</body>
    </html>
  );
}

