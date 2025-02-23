"use client";

import { UserProfile } from "@clerk/nextjs";
import { Box } from "@primer/react";
import { BackButton } from "@/components/back-button";

export default function Page() {
  return (
    <Box
      sx={{
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
        minHeight: "calc(100vh - 200px)",
      }}
    >
      <div>
        <BackButton href={"/"}></BackButton>
        <UserProfile path="/user-profile" />
      </div>
    </Box>
  );
}
