"use client";

import * as React from "react";
import {
  ThemeProvider as PrimerThemeProvider,
  BaseStyles,
} from "@primer/react";

export function ThemeProvider({
  children,
  ...props
}: React.ComponentProps<typeof PrimerThemeProvider>) {
  return (
    <PrimerThemeProvider {...props} colorMode={"night"}>
      <BaseStyles>{children}</BaseStyles>
    </PrimerThemeProvider>
  );
}
