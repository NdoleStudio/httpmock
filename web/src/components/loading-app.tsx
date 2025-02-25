"use client";

import * as React from "react";
import { Box, Heading, Spinner, Text } from "@primer/react";
import Logo from "@/app/logo.svg";
import Image from "next/image";

export function LoadingApp() {
  const [show, setShow] = React.useState(false);
  React.useEffect(() => {
    setShow(true);
  }, []);
  return (
    <React.Fragment>
      {show && (
        <Box
          sx={{
            display: "flex",
            alignItems: "center",
            justifyContent: "center",
            minHeight: "calc(100vh - 200px)",
          }}
        >
          <div style={{ marginTop: -32 }}>
            <div style={{ textAlign: "center" }}>
              <Image src={Logo} height={120} alt={"httpmock.dev logo"} />
            </div>
            <div style={{ display: "flex" }}>
              <Spinner
                sx={{ mt: 2, color: "accent.emphasis" }}
                size={"medium"}
              />
              <Heading as={"h2"} sx={{ ml: 2, fontWeight: "light" }}>
                Loading httpmock<Text sx={{ color: "accent.emphasis" }}>.</Text>
                dev
              </Heading>
            </div>
          </div>
        </Box>
      )}
    </React.Fragment>
  );
}
