import type { Metadata } from "next";
import "./globals.css";
import { ClerkProvider } from "@clerk/nextjs";
import { ThemeProvider } from "@/components/theme-provider";
import { AppBar } from "@/components/app-bar";
import { dark } from "@clerk/themes";

export const metadata: Metadata = {
  title: "Mock HTTP Server",
  description:
    "This is a mock http server which can be used to test HTTP requests and responses when building an HTTP client.",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <ClerkProvider appearance={{ baseTheme: dark }}>
      <html lang="en">
        <body>
          <ThemeProvider>
            <AppBar></AppBar>
            <main>{children}</main>
          </ThemeProvider>
        </body>
      </html>
    </ClerkProvider>
  );
}
