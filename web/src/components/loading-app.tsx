"use client";

import * as React from "react";
import { Box, Spinner, Heading, Text } from "@primer/styled-react";
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
          style={{
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
                style={{ marginTop: 16, color: "#0969da" }}
                size={"medium"}
              />
              <Heading
                as={"h2"}
                style={{ marginLeft: 16, fontWeight: "light" }}
              >
                Loading httpmock<Text style={{ color: "#0969da" }}>.</Text>
                dev
              </Heading>
            </div>
          </div>
        </Box>
      )}
    </React.Fragment>
  );
}
