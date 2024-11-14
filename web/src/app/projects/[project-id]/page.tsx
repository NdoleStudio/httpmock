"use client";

import {
  Box,
  Text,
  Button,
  PageHeader,
  ActionMenu,
  ActionList,
  Spinner,
} from "@primer/react";
import { usePathname, useRouter } from "next/navigation";
import { MouseEvent, useEffect, useState } from "react";
import { ErrorMessages } from "@/utils/errors";
import { useAppStore } from "@/store/provider";
import { EntitiesProject } from "@/api/model";
import {
  GearIcon,
  PencilIcon,
  PlusIcon,
  TrashIcon,
} from "@primer/octicons-react";
import { useAuth } from "@clerk/nextjs";
import { setAuthHeader } from "@/api/axios";

export default function ProjectShow() {
  const router = useRouter();
  const auth = useAuth();
  const pathName = usePathname();
  const { fetchProject } = useAppStore((state) => state);

  const [errorMessages, setErrorMessages] = useState<ErrorMessages>(
    ErrorMessages.create(),
  );
  const [loading, setLoading] = useState<boolean>(false);
  const [project, setProject] = useState<EntitiesProject | undefined>(
    undefined,
  );

  const projectId = pathName.split("/")[2];

  const loadProject = () => {
    fetchProject(projectId)
      .then((project: EntitiesProject) => {
        setProject(project);
      })
      .catch((errorMessages: ErrorMessages) => {
        setErrorMessages(errorMessages);
      });
  };

  useEffect(() => {
    if (auth.isLoaded) {
      auth.getToken().then((token: string | null) => {
        setAuthHeader(token);
        loadProject();
      });
    }
  }, [projectId, auth.isLoaded]);

  return (
    <Box
      sx={{
        maxWidth: "xlarge",
        mx: "auto",
        mt: 6,
        minHeight: "calc(100vh - 200px)",
      }}
    >
      <PageHeader role="banner" aria-label="project details">
        <PageHeader.TitleArea variant={"large"}>
          {project && (
            <PageHeader.Title>{project && project.name}</PageHeader.Title>
          )}
          {!project && (
            <PageHeader.Title>
              <Spinner size="large" />
            </PageHeader.Title>
          )}
        </PageHeader.TitleArea>
        {project && (
          <PageHeader.Description>
            <Text sx={{ fontSize: 1, color: "fg.muted" }}>
              {project?.description}
            </Text>
          </PageHeader.Description>
        )}
        <PageHeader.Actions>
          <Button
            disabled={!project}
            variant="primary"
            leadingVisual={PlusIcon}
          >
            New Endpoint
          </Button>
          <ActionMenu>
            <ActionMenu.Anchor aria-label={"Manage project"}>
              <Button disabled={!project} leadingVisual={GearIcon}></Button>
            </ActionMenu.Anchor>
            <ActionMenu.Overlay>
              <ActionList>
                <ActionList.Item onSelect={() => alert("Edit comment clicked")}>
                  Edit Project
                  <ActionList.LeadingVisual>
                    <PencilIcon></PencilIcon>
                  </ActionList.LeadingVisual>
                </ActionList.Item>
                <ActionList.Divider />
                <ActionList.Item
                  variant="danger"
                  onSelect={() => alert("Delete file clicked")}
                >
                  <ActionList.LeadingVisual>
                    <TrashIcon />
                  </ActionList.LeadingVisual>
                  Delete Project
                </ActionList.Item>
              </ActionList>
            </ActionMenu.Overlay>
          </ActionMenu>
        </PageHeader.Actions>
      </PageHeader>
    </Box>
  );
}
