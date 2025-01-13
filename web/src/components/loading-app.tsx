"use client";

import * as React from "react";
import { Box, Heading, Spinner, Text } from "@primer/react";
import Logo from "@/app/logo.svg";
import Image from "next/image";

export function LoadingApp() {
  return (
    <Box
      sx={{
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
        minHeight: "calc(100vh - 200px)",
      }}
    >
      <Box sx={{ mt: "-32px" }}>
        <Box sx={{ textAlign: "center" }}>
          <Image src={Logo} height={120} alt={"httpmock.dev logo"} />
        </Box>
        <Box sx={{ display: "flex" }}>
          <Spinner sx={{ mt: 2, color: "accent.emphasis" }} size={"medium"} />
          <Heading as={"h2"} sx={{ ml: 2, fontWeight: "light" }}>
            Loading httpmock<Text sx={{ color: "accent.emphasis" }}>.</Text>dev
          </Heading>
        </Box>
      </Box>
    </Box>
  );
}
