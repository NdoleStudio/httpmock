"use client";

import * as React from "react";
import {
  ThemeProvider as PrimerThemeProvider,
  BaseStyles,
} from "@primer/react";
import { Toaster } from "sonner";
import { AppStoreProvider } from "@/store/provider";
import { useAuth } from "@clerk/nextjs";
import { useEffect } from "react";
import { setAuthHeader } from "@/api/axios";
import { LoadingApp } from "@/components/loading-app";

export function ThemeProvider({
  children,
  ...props
}: React.ComponentProps<typeof PrimerThemeProvider>) {
  const auth = useAuth();
  const [loading, setLoading] = React.useState<boolean>(true);

  useEffect(() => {
    console.log("auth.isLoaded", auth.isLoaded);
    if (auth.isLoaded) {
      auth
        .getToken()
        .then((token: string | null) => {
          if (token) {
            setAuthHeader(token);
          }
        })
        .finally(() => {
          setLoading(false);
        });
    }
  }, [auth.isLoaded]);

  return (
    <AppStoreProvider>
      <PrimerThemeProvider {...props} colorMode={"night"}>
        <BaseStyles>
          <Toaster closeButton={true} richColors={true} />
          {!loading && children}
          {loading && <LoadingApp />}
        </BaseStyles>
      </PrimerThemeProvider>
    </AppStoreProvider>
  );
}
