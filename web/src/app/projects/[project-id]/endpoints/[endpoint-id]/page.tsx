"use client";

import {
  Box,
  Text,
  PageHeader,
  Spinner,
  Dialog,
  BranchName,
  Heading,
  Label,
} from "@primer/react";
import { usePathname } from "next/navigation";
import { useCallback, useEffect, useState } from "react";
import { useAppStore } from "@/store/provider";
import {
  EntitiesProject,
  EntitiesProjectEndpoint,
  EntitiesProjectEndpointRequest,
} from "@/api/model";
import { MirrorIcon } from "@primer/octicons-react";
import { CopyButton } from "@/components/copy-button";
import { getEndpointURL, labelColor } from "@/utils/filters";
import { BackButton } from "@/components/back-button";

export default function EndpointShow() {
  const pathName = usePathname();
  const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false);
  const { showProject, showProjectEndpoint, indexProjectEndpointRequests } =
    useAppStore((state) => state);
  const [loadingEndpoint, setLoadingEndpoint] = useState<boolean>(false);
  const [loadingProjectEndpointRequests, setLoadingProjectEndpointRequests] =
    useState<boolean>(false);
  const [project, setProject] = useState<EntitiesProject | undefined>(
    undefined,
  );
  const [projectEndpoint, setProjectEndpoint] = useState<
    EntitiesProjectEndpoint | undefined
  >(undefined);
  const [projectEndpointRequests, setProjectEndpointRequests] = useState<
    Array<EntitiesProjectEndpointRequest>
  >([]);

  const onDeleteDialogClose = useCallback(
    () => setIsDeleteDialogOpen(false),
    [],
  );

  const projectId = pathName.split("/")[2];
  const projectEndpointId = pathName.split("/")[4];

  useEffect(() => {
    setLoadingEndpoint(true);
    showProjectEndpoint(projectId, projectEndpointId)
      .then((endpoint: EntitiesProjectEndpoint) => {
        setProjectEndpoint(endpoint);
      })
      .finally(() => {
        setLoadingEndpoint(false);
      });
  }, [projectId, projectEndpointId, showProjectEndpoint]);

  useEffect(() => {
    showProject(projectId).then((project: EntitiesProject) => {
      setProject(project);
    });
  }, [showProject, projectId]);

  useEffect(() => {
    setLoadingProjectEndpointRequests(false);
    indexProjectEndpointRequests(projectId, projectEndpointId)
      .then((requests: Array<EntitiesProjectEndpointRequest>) => {
        setProjectEndpointRequests(requests);
      })
      .finally(() => {
        setLoadingProjectEndpointRequests(false);
      });
  }, [indexProjectEndpointRequests, projectId, projectEndpointId]);

  return (
    <Box
      sx={{
        maxWidth: "xlarge",
        mx: "auto",
        mt: 6,
        minHeight: "calc(100vh - 200px)",
      }}
    >
      <PageHeader role="banner" aria-label="Project details">
        <PageHeader.TitleArea aria-label={"Project endpoint"} variant={"large"}>
          {project && projectEndpoint && (
            <PageHeader.Title>
              <Box sx={{ display: "flex", alignItems: "flexStart" }}>
                <Label
                  sx={{
                    mt: 3,
                    color: labelColor(projectEndpoint.request_method),
                  }}
                >
                  {projectEndpoint.request_method}
                </Label>
                <Text sx={{ ml: 1, mr: 1, fontWeight: "bold" }}>
                  {getEndpointURL(project, projectEndpoint.request_path)}
                </Text>
                <CopyButton
                  size={"medium"}
                  sx={{ mt: 2 }}
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
              disabled: loadingEndpoint,
              content: "Close",
              onClick: onDeleteDialogClose,
            },
            {
              buttonType: "danger",
              block: true,
              loading: loadingEndpoint,
              disabled: loadingEndpoint,
              content: "Delete this Request",
            },
          ]}
        >
          <div>
            <p>
              Are you sure you want to delete the{" "}
              <BranchName>{project?.name}</BranchName> project. This is a
              permanent action and it cannot be reversed.
            </p>
          </div>
          <Box sx={{ mt: 2 }}>
            <Text sx={{ color: "fg.muted" }}>
              {project?.subdomain}.httpmock.dev
            </Text>
          </Box>
        </Dialog>
      )}

      <Heading as="h2" sx={{ mt: 32 }} variant="medium">
        <MirrorIcon size={24} />
        <Text sx={{ ml: 2 }}>HTTP Requests</Text>
      </Heading>
      {!loadingProjectEndpointRequests && (
        <Box
          sx={{
            mt: 2,
            "> *": {
              borderWidth: 1,
              borderColor: "border.default",
              borderStyle: "solid",
              borderBottomWidth: 0,
              padding: 4,
              "&:first-child": {
                borderTopLeftRadius: 2,
                borderTopRightRadius: 2,
              },
              "&:last-child": {
                borderBottomLeftRadius: 2,
                borderBottomRightRadius: 2,
                borderBottomWidth: 1,
              },
              "&:hover": {
                bg: "canvas.inset",
              },
            },
          }}
        >
          {projectEndpointRequests.map((request) => (
            <div key={request.id}>
              <Box sx={{ display: "flex", alignItems: "baseline" }}>
                <Label sx={{ color: labelColor(request.request_method) }}>
                  {request.request_method}
                </Label>
                <Text sx={{ ml: 1, mr: 1, fontWeight: "bold" }}>
                  {request.request_url}
                </Text>
                <CopyButton data={request.request_url} />
              </Box>
              <Box sx={{ mt: 2, display: "flex" }}></Box>
            </div>
          ))}
        </Box>
      )}
    </Box>
  );
}
