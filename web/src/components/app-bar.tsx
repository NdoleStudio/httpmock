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
  Text,
} from "@primer/react";
import { Heading } from "@primer/react";
import Logo from "@/app/logo.svg";
import {
  SignInIcon,
  SignOutIcon,
  PersonFillIcon,
  PlusIcon,
} from "@primer/octicons-react";
import { usePathname, useRouter } from "next/navigation";
import { useClerk, useUser } from "@clerk/nextjs";
import { useAppStore } from "@/store/provider";
import { EntitiesProject } from "@/api/model";
import { useEffect, useState } from "react";

export function AppBar() {
  const router = useRouter();
  const pathname = usePathname();
  const user = useUser();
  const { signOut } = useClerk();
  const { indexProjects } = useAppStore((state) => state);
  const [projects, setProjects] = useState<Array<EntitiesProject>>([]);
  const [activeProject, setActiveProject] = useState<EntitiesProject | null>(
    null,
  );

  useEffect(() => {
    indexProjects().then((projects: Array<EntitiesProject>) => {
      setProjects(projects);
    });
  }, [indexProjects]);

  useEffect(() => {
    projects.forEach((project: EntitiesProject) => {
      if (pathname.startsWith(`/projects/${project.id}`)) {
        setActiveProject(project);
      }
    });
    if (
      !pathname.startsWith("/projects") ||
      pathname.startsWith("/projects/create")
    ) {
      setActiveProject(null);
    }
  }, [pathname, projects]);

  return (
    <Header>
      <Header.Item
        style={{ cursor: "pointer" }}
        onClick={() => router.push("/")}
      >
        <Image src={Logo} height={48} alt={"httpmock.dev logo"} />
        <Heading
          sx={{ fontWeight: "lighter", display: ["none", "block"] }}
          as={"h2"}
        >
          httpmock
        </Heading>
        {activeProject === null && (
          <Heading
            sx={{ fontWeight: "lighter", display: ["block", "none"] }}
            as={"h2"}
          >
            httpmock
          </Heading>
        )}
      </Header.Item>
      {activeProject && (
        <ActionMenu>
          <ActionMenu.Button size={"large"} sx={{ mt: 2 }}>
            <Text variant={"large"} weight={"semibold"}>
              {activeProject.name}
            </Text>
          </ActionMenu.Button>
          <ActionMenu.Overlay width="auto">
            <ActionList selectionVariant="single">
              {projects.map((project: EntitiesProject) => (
                <ActionList.Item
                  key={project.id}
                  onSelect={() => router.push(`/projects/${project.id}`)}
                  selected={project.id === activeProject?.id}
                >
                  <Text weight={"semibold"}>{project.name}</Text>
                </ActionList.Item>
              ))}
            </ActionList>
            <ActionList.Group>
              <ActionList.Divider />
              <ActionList.Item>
                <Button
                  onClick={() => router.push("/projects/create")}
                  variant={"primary"}
                  block={true}
                  leadingVisual={PlusIcon}
                >
                  Create Project
                </Button>
              </ActionList.Item>
            </ActionList.Group>
          </ActionMenu.Overlay>
        </ActionMenu>
      )}
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
