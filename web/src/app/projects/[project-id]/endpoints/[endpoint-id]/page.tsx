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
  Heading,
  Label,
} from "@primer/react";
import { usePathname, useRouter } from "next/navigation";
import { useCallback, useEffect, MouseEvent, useState } from "react";
import { useAppStore } from "@/store/provider";
import { EntitiesProject, EntitiesProjectEndpoint } from "@/api/model";
import {
  GearIcon,
  MirrorIcon,
  PencilIcon,
  PlusIcon,
  TrashIcon,
} from "@primer/octicons-react";
import { CopyButton } from "@/components/copy-button";
import { getEndpointURL, labelColor } from "@/utils/filters";
import { BackButton } from "@/components/back-button";

export default function EndpointShow() {
  const router = useRouter();
  const pathName = usePathname();
  const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false);
  const { showProject, showProjectEndpoint } = useAppStore((state) => state);
  const [loading, setLoading] = useState<boolean>(false);
  const [endpoints, setEndpoints] = useState<EntitiesProjectEndpoint[]>([]);
  const [loadingEndpoint, setLoadingEndpoint] = useState<boolean>(false);
  const [project, setProject] = useState<EntitiesProject | undefined>(
    undefined,
  );
  const [projectEndpoint, setProjectEndpoint] = useState<
    EntitiesProjectEndpoint | undefined
  >(undefined);

  const onDeleteDialogClose = useCallback(
    () => setIsDeleteDialogOpen(false),
    [],
  );

  const projectId = pathName.split("/")[2];
  const projectEndpointId = pathName.split("/")[4];

  const loadProjectEndpoint = () => {
    setLoadingEndpoint(true);
    showProjectEndpoint(projectId, projectEndpointId)
      .then((endpoint: EntitiesProjectEndpoint) => {
        setProjectEndpoint(endpoint);
      })
      .finally(() => {
        setLoadingEndpoint(false);
      });
  };
  useEffect(() => {
    loadProjectEndpoint();
  }, [projectId, projectEndpointId]);

  const loadProject = () => {
    showProject(projectId).then((project: EntitiesProject) => {
      setProject(project);
    });
  };
  useEffect(() => {
    loadProject();
  }, [projectId]);

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
          {project && projectEndpoint && (
            <PageHeader.Title>
              <Box sx={{ display: "flex", alignItems: "baseline" }}>
                <Label
                  sx={{ color: labelColor(projectEndpoint.request_method) }}
                >
                  {projectEndpoint.request_method}
                </Label>
                <Text sx={{ ml: 1, mr: 1, fontWeight: "bold" }}>
                  {getEndpointURL(project, projectEndpoint.request_path)}
                </Text>
                <CopyButton
                  data={getEndpointURL(project, projectEndpoint.request_path)}
                />
              </Box>
            </PageHeader.Title>
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
          <BackButton href={`/projects/${projectId}/`} />
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
              content: "Delete this Reuest",
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

      <Box>
        <Heading as="h2" sx={{ mt: 32 }} variant="medium">
          <MirrorIcon size={24} />
          <Text sx={{ ml: 2 }}>HTTP Requests</Text>
        </Heading>
      </Box>
      <Box>
        <Spinner size="large" />
      </Box>
    </Box>
  );
}
