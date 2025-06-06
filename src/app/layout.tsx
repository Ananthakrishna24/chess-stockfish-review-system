import type { Metadata } from "next";
import { Geist, Geist_Mono } from "next/font/google";
import "./globals.css";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title: "Chess Game Review | Professional Analysis with Stockfish",
  description: "Analyze your chess games with professional-grade Stockfish engine. Get detailed move analysis, accuracy ratings, and improve your chess skills with comprehensive game review.",
  keywords: ["chess", "game review", "stockfish", "chess analysis", "chess engine", "move analysis", "chess improvement"],
  authors: [{ name: "Chess Review System" }],
  creator: "Chess Review Professional",
  publisher: "Chess Analysis Platform",
  robots: "index, follow",
  viewport: "width=device-width, initial-scale=1",
  themeColor: "#2c3e50",
  openGraph: {
    title: "Chess Game Review | Professional Analysis",
    description: "Analyze your chess games with professional-grade Stockfish engine",
    type: "website",
    locale: "en_US"
  }
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <head>
        <link rel="icon" href="/favicon.ico" />
        <meta name="theme-color" content="#2c3e50" />
      </head>
      <body
        className={`${geistSans.variable} ${geistMono.variable} antialiased`}
      >
        {children}
      </body>
    </html>
  );
}
