"use client";

import * as React from "react";
import Image from "next/image";
import {
  Header,
  Button,
  ActionMenu,
  ActionList,
  Avatar,
  Box,
} from "@primer/react";
import { Heading } from "@primer/react";
import Logo from "@/app/logo.svg";
import {
  SignInIcon,
  SignOutIcon,
  PersonFillIcon,
} from "@primer/octicons-react";
import { usePathname, useRouter } from "next/navigation";
import { useClerk, useUser } from "@clerk/nextjs";

export function AppBar() {
  const router = useRouter();
  const pathname = usePathname();
  const user = useUser();
  const { signOut } = useClerk();

  return (
    <Header>
      <Header.Item sx={{ cursor: "pointer" }} onClick={() => router.push("/")}>
        <Image src={Logo} height={40} alt={"httpmock.dev logo"} />
        <Heading sx={{ fontWeight: "lighter" }} as={"h1"}>
          httpmock
        </Heading>
      </Header.Item>
      <Box sx={{ flex: "auto" }}></Box>
      <Header.Item>
        {pathname === "/sign-up" && (
          <Button
            leadingVisual={SignInIcon}
            variant={"primary"}
            size={"large"}
            onClick={() => router.push("/sign-in")}
          >
            Sign In
          </Button>
        )}
        {pathname === "/sign-in" && (
          <Button
            leadingVisual={SignInIcon}
            variant={"primary"}
            size={"large"}
            onClick={() => router.push("/sign-up")}
          >
            Sign Up
          </Button>
        )}
        {user.isSignedIn && (
          <ActionMenu>
            <ActionMenu.Anchor>
              <Avatar
                size={30}
                sx={{ cursor: "pointer" }}
                alt={"account image"}
                src={user.user.imageUrl}
              />
            </ActionMenu.Anchor>
            <ActionMenu.Overlay>
              <ActionList>
                <ActionList.Item onSelect={() => router.push("/user-profile")}>
                  <ActionList.LeadingVisual>
                    <PersonFillIcon />
                  </ActionList.LeadingVisual>
                  Manage Account
                </ActionList.Item>
                <ActionList.Divider />
                <ActionList.Item
                  variant="danger"
                  onSelect={() => signOut({ redirectUrl: "/" })}
                >
                  <ActionList.LeadingVisual>
                    <SignOutIcon />
                  </ActionList.LeadingVisual>
                  Logout
                </ActionList.Item>
              </ActionList>
            </ActionMenu.Overlay>
          </ActionMenu>
        )}
      </Header.Item>
    </Header>
  );
}
