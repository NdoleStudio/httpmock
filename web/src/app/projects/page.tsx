"use client";

import { redirect, useRouter } from "next/navigation";
import {
  Box,
  Button,
  Heading,
  PageHeader,
  Spinner,
  Text,
} from "@primer/styled-react";
import { BranchName } from "@primer/react";
import { PlusIcon } from "@primer/octicons-react";
import { useEffect, useState } from "react";
import { EntitiesProject } from "@/api/model";
import { useAppStore } from "@/store/provider";
import { toast } from "sonner";

export default function ProjectIndex() {
  const router = useRouter();
  const [projects, setProjects] = useState<Array<EntitiesProject>>([]);
  const { indexProjects } = useAppStore((state) => state);

  useEffect(() => {
    indexProjects().then((projects: Array<EntitiesProject>) => {
      if (projects.length === 0) {
        toast.info("You don't have any projects yet. Let's create one.");
        redirect("/projects/create");
      } else {
        setProjects(projects);
      }
    });
  }, [indexProjects]);

  return (
    <Box
      style={{
        maxWidth: 1200,
        marginLeft: "auto",
        marginRight: "auto",
        paddingLeft: 16,
        paddingRight: 16,
        marginTop: 48,
        minHeight: "calc(100vh - 200px)",
      }}
    >
      <PageHeader role="banner" aria-label="Project list">
        <PageHeader.TitleArea variant={"large"}>
          All Projects
        </PageHeader.TitleArea>
        <PageHeader.Actions>
          <Button
            variant="primary"
            onClick={() => {
              router.push("/projects/create");
            }}
            leadingVisual={PlusIcon}
          >
            Create Project
          </Button>
        </PageHeader.Actions>
      </PageHeader>
      <Box
        style={{
          borderBottomWidth: 1,
          marginTop: 16,
          borderBottomStyle: "solid",
          borderColor: "#d0d7de",
        }}
      ></Box>

      {projects.length === 0 && (
        <Box style={{ textAlign: "center", marginTop: 32 }}>
          <Spinner size="large" style={{ color: "#0969da" }} />
        </Box>
      )}

      <Box
        style={{
          marginTop: 32,
          display: "grid",
          gridTemplateColumns: "repeat(auto-fill, minmax(250px, 1fr))",
          gap: 24,
        }}
      >
        {projects.map((project: EntitiesProject) => (
          <Box
            key={project.id}
            style={{
              cursor: "pointer",
              padding: 24,
              backgroundColor: "#f6f8fa",
              borderWidth: 1,
              borderRadius: 6,
              borderStyle: "solid",
              borderColor: "#d0d7de",
            }}
            onClick={() => {
              router.push(`/projects/${project.id}`);
            }}
          >
            <Heading className="text--clamp-1" variant="medium" as="h2">
              {project.name}
            </Heading>
            <Text as="p" className="text--clamp-3" style={{ color: "#6e7781" }}>
              {project.description}
            </Text>
            <BranchName>{project.subdomain}.httpmock.dev</BranchName>
          </Box>
        ))}
      </Box>
    </Box>
  );
}
