"use client";

import * as React from "react";
import {
  ThemeProvider as PrimerThemeProvider,
  BaseStyles,
} from "@primer/react";
import { Toaster } from "sonner";
import { AppStoreProvider } from "@/store/provider";

export function ThemeProvider({
  children,
  ...props
}: React.ComponentProps<typeof PrimerThemeProvider>) {
  return (
    <AppStoreProvider>
      <PrimerThemeProvider {...props} colorMode={"night"}>
        <BaseStyles>
          <Toaster closeButton={true} richColors={true} />
          {children}
        </BaseStyles>
      </PrimerThemeProvider>
    </AppStoreProvider>
  );
}
