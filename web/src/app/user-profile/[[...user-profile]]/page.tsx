"use client";

import { UserProfile } from "@clerk/nextjs";
import { Box } from "@primer/styled-react";
import { BackButton } from "@/components/back-button";

export default function Page() {
  return (
    <Box
      style={{
        marginTop: 16,
        paddingLeft: 16,
        paddingRight: 16,
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
