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
  RelativeTime,
} from "@primer/react";
import { usePathname, useRouter } from "next/navigation";
import { useCallback, useEffect, MouseEvent, useState } from "react";
import { useAppStore } from "@/store/provider";
import { EntitiesProject, EntitiesProjectEndpoint } from "@/api/model";
import {
  GearIcon,
  LinkIcon,
  PencilIcon,
  PlusIcon,
  TrashIcon,
} from "@primer/octicons-react";
import { toast } from "sonner";
import { CopyButton } from "@/components/copy-button";

export default function ProjectShow() {
  const router = useRouter();
  const pathName = usePathname();
  const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false);
  const { showProject, deleteProject, indexProjectEndpoint } = useAppStore(
    (state) => state,
  );
  const [loading, setLoading] = useState<boolean>(false);
  const [endpoints, setEndpoints] = useState<EntitiesProjectEndpoint[]>([]);
  const [loadingEndpoints, setLoadingEndpoints] = useState<boolean>(false);
  const [project, setProject] = useState<EntitiesProject | undefined>(
    undefined,
  );

  const onDeleteDialogClose = useCallback(
    () => setIsDeleteDialogOpen(false),
    [],
  );

  const projectId = pathName.split("/")[2];

  useEffect(() => {
    setLoading(true);
    showProject(projectId)
      .then((project: EntitiesProject) => {
        setProject(project);
      })
      .finally(() => {
        setLoading(false);
      });
  }, [showProject, projectId]);

  useEffect(() => {
    setLoadingEndpoints(true);
    indexProjectEndpoint(projectId)
      .then((endpoints: EntitiesProjectEndpoint[]) => {
        setEndpoints(endpoints);
      })
      .finally(() => {
        setLoadingEndpoints(false);
      });
  }, [indexProjectEndpoint, projectId]);

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

  const getLabelColor = (requestMethod: string): string => {
    switch (requestMethod) {
      case "GET":
        return "accent.fg";
      case "POST":
        return "success.fg";
      case "PUT":
        return "attention.fg";
      case "DELETE":
        return "danger.fg";
      case "ANY":
        return "done.fg";
      case "OPTIONS":
        return "sponsors.fg";
      default:
        return "accent.fg";
    }
  };

  const getEndpointURL = (path: string): string => {
    return `https://${project?.subdomain}.httpmock.dev${path === "/" ? "" : path}`;
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
      <PageHeader role="banner" aria-label="Project details">
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
                  onSelect={() => router.push(`/projects/${projectId}/edit`)}
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

      <div>
        <Heading as="h2" sx={{ mt: 32 }} variant="medium">
          <LinkIcon size={24} />
          <Text sx={{ ml: 2 }}>Endpoints</Text>
        </Heading>
      </div>
      {loadingEndpoints && (
        <div>
          <Spinner size="large" />
        </div>
      )}
      {!loadingEndpoints && (
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
          {endpoints.map((endpoint) => (
            <div key={endpoint.id}>
              <Box sx={{ display: "flex", alignItems: "baseline" }}>
                <Label sx={{ color: getLabelColor(endpoint.request_method) }}>
                  {endpoint.request_method}
                </Label>
                <Text sx={{ ml: 1, mr: 1, fontWeight: "bold" }}>
                  {getEndpointURL(endpoint.request_path)}
                </Text>
                <CopyButton data={getEndpointURL(endpoint.request_path)} />
              </Box>
              <Box sx={{ mt: 2, display: "flex" }}>
                <Button
                  onClick={() =>
                    router.push(
                      `/projects/${projectId}/endpoints/${endpoint.id}`,
                    )
                  }
                  size={"small"}
                  variant="default"
                  sx={{ mr: 2 }}
                >
                  View Requests
                </Button>
                <Button
                  onClick={() =>
                    router.push(
                      `/projects/${projectId}/endpoints/${endpoint.id}/edit`,
                    )
                  }
                  size={"small"}
                  variant="default"
                  sx={{ mr: 2 }}
                  leadingVisual={GearIcon}
                >
                  Manage
                </Button>
                <div>
                  <Text size="small" sx={{ mr: 1, color: "fg.muted" }}>
                    Updated
                  </Text>
                  <RelativeTime
                    sx={{ color: "fg.muted", fontSize: "small" }}
                    date={new Date(endpoint.updated_at)}
                    noTitle={true}
                  ></RelativeTime>
                </div>
              </Box>
            </div>
          ))}
        </Box>
      )}
    </Box>
  );
}
