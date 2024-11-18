"use client";

import {
  Box,
  Text,
  Button,
  PageHeader,
  ActionMenu,
  ActionList,
  Spinner,
  Dialog,
  BranchName,
} from "@primer/react";
import { usePathname, useRouter } from "next/navigation";
import { useCallback, useEffect, MouseEvent, useState } from "react";
import { ErrorMessages } from "@/utils/errors";
import { useAppStore } from "@/store/provider";
import { EntitiesProject } from "@/api/model";
import {
  GearIcon,
  PencilIcon,
  PlusIcon,
  TrashIcon,
} from "@primer/octicons-react";
import { toast } from "sonner";

export default function ProjectShow() {
  const router = useRouter();
  const pathName = usePathname();
  const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false);
  const { showProject, deleteProject } = useAppStore((state) => state);

  const [errorMessages, setErrorMessages] = useState<ErrorMessages>(
    ErrorMessages.create(),
  );
  const [loading, setLoading] = useState<boolean>(false);
  const [project, setProject] = useState<EntitiesProject | undefined>(
    undefined,
  );

  const onDeleteDialogClose = useCallback(
    () => setIsDeleteDialogOpen(false),
    [],
  );

  const projectId = pathName.split("/")[2];

  const loadProject = () => {
    showProject(projectId)
      .then((project: EntitiesProject) => {
        setProject(project);
      })
      .catch((errorMessages: ErrorMessages) => {
        setErrorMessages(errorMessages);
      });
  };

  useEffect(() => {
    loadProject();
  }, [projectId]);

  const onDeleteProject = async (event: MouseEvent) => {
    event.preventDefault();
    setLoading(true);
    deleteProject(projectId)
      .then(() => {
        setIsDeleteDialogOpen(false);
        toast.info("Project deleted successfully.");
        router.push("/projects");
      })
      .finally(() => {
        setLoading(false);
      });
  };

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
            onClick={() =>
              router.push(`/projects/${projectId}/endpoints/create`)
            }
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
                <ActionList.Item
                  onClick={() => router.push(`/projects/${projectId}/edit`)}
                >
                  Edit Project
                  <ActionList.LeadingVisual>
                    <PencilIcon></PencilIcon>
                  </ActionList.LeadingVisual>
                </ActionList.Item>
                <ActionList.Divider />
                <ActionList.Item
                  variant="danger"
                  onSelect={() => setIsDeleteDialogOpen(true)}
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
      <Box
        sx={{
          borderBottomWidth: 1,
          mt: 2,
          borderBottomStyle: "solid",
          borderColor: "border.default",
        }}
      ></Box>
      {isDeleteDialogOpen && (
        <Dialog
          width={"large"}
          onClose={onDeleteDialogClose}
          title={`Delete ${project?.name}`}
          footerButtons={[
            {
              disabled: loading,
              content: "Close",
              onClick: onDeleteDialogClose,
            },
            {
              buttonType: "danger",
              block: true,
              loading: loading,
              disabled: loading,
              content: "Delete this project",
              onClick: onDeleteProject,
            },
          ]}
        >
          <Box>
            <Text>
              Are you sure you want to delete the{" "}
              <BranchName>{project?.name}</BranchName> project. This is a
              permanent action and it cannot be reversed.
            </Text>
          </Box>
          <Box sx={{ mt: 2 }}>
            <Text sx={{ color: "fg.muted" }}>
              {project?.subdomain}.httpmock.dev
            </Text>
          </Box>
        </Dialog>
      )}
    </Box>
  );
}
