"use client";

import * as React from "react";
import {
  ThemeProvider as PrimerThemeProvider,
  BaseStyles,
} from "@primer/react";
import { StoreProvider } from "@/store/provider";

export function ThemeProvider({
  children,
  ...props
}: React.ComponentProps<typeof PrimerThemeProvider>) {
  return (
    <StoreProvider>
      <PrimerThemeProvider {...props} colorMode={"night"}>
        <BaseStyles>{children}</BaseStyles>
      </PrimerThemeProvider>
    </StoreProvider>
  );
}
