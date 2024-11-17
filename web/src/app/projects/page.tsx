"use client";

import { redirect, useRouter } from "next/navigation";
import {
  Box,
  BranchName,
  Button,
  Heading,
  PageHeader,
  Spinner,
  Text,
} from "@primer/react";
import { PlusIcon } from "@primer/octicons-react";
import { useEffect, useState } from "react";
import { EntitiesProject } from "@/api/model";
import { useAppStore } from "@/store/provider";
import { toast } from "sonner";

export default function ProjectIndex() {
  const router = useRouter();
  const [projects, setProjects] = useState<Array<EntitiesProject>>([]);
  const { indexProjects } = useAppStore((state) => state);

  const loadProjects = () => {
    indexProjects().then((projects: Array<EntitiesProject>) => {
      if (projects.length === 0) {
        toast.info("You don't have any projects yet. Let's create one.");
        redirect("/projects/create");
      } else {
        setProjects(projects);
      }
    });
  };

  useEffect(() => {
    loadProjects();
  }, []);

  return (
    <Box
      sx={{
        maxWidth: "xlarge",
        mx: "auto",
        mt: 6,
        minHeight: "calc(100vh - 200px)",
      }}
    >
      <PageHeader role="banner" aria-label="project list">
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
        sx={{
          borderBottomWidth: 1,
          mt: 2,
          borderBottomStyle: "solid",
          borderColor: "border.default",
        }}
      ></Box>

      {projects.length === 0 && (
        <Box sx={{ textAlign: "center", mt: 4 }}>
          <Spinner size="large" sx={{ color: "accent.emphasis" }} />
        </Box>
      )}
      <Box
        sx={{
          mt: 4,
          display: "grid",
          gridTemplateColumns: "1fr 1fr 1fr 1fr",
          gap: 3,
        }}
      >
        {projects.map((project: EntitiesProject) => (
          <Box
            key={project.id}
            sx={{
              cursor: "pointer",
              p: 3,
              backgroundColor: "canvas.inset",
              borderWidth: 1,
              borderRadius: 2,
              borderStyle: "solid",
              borderColor: "border.default",
              ":hover": {
                bg: "neutral.muted",
              },
            }}
            onClick={() => {
              router.push(`/projects/${project.id}`);
            }}
          >
            <Heading className="text--clamp-1" variant="medium" as="h2">
              {project.name}
            </Heading>
            <Text as="p" className="text--clamp-3" sx={{ color: "fg.muted" }}>
              {project.description}
            </Text>
            <BranchName>{project.subdomain}.httpmock.dev</BranchName>
          </Box>
        ))}
      </Box>
    </Box>
  );
}
