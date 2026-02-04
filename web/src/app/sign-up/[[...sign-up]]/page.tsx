"use client";

import { SignUp } from "@clerk/nextjs";
import { Box } from "@primer/styled-react";

export default function Page() {
  return (
    <Box
      style={{
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
        minHeight: "calc(100vh - 200px)",
      }}
    >
      <SignUp />
    </Box>
  );
}
