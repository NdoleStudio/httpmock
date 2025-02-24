"use client";

import {
  Box,
  Text,
  PageHeader,
  Spinner,
  Dialog,
  BranchName,
  Heading,
  Button,
  Label,
  TextInput,
  Textarea,
  RelativeTime,
} from "@primer/react";
import { usePathname } from "next/navigation";
import React, { MouseEvent, useCallback, useEffect, useState } from "react";
import { useAppStore } from "@/store/provider";
import {
  EntitiesProject,
  EntitiesProjectEndpoint,
  EntitiesProjectEndpointRequest,
} from "@/api/model";
import { MirrorIcon, TrashIcon } from "@primer/octicons-react";
import { CopyButton } from "@/components/copy-button";
import { getEndpointURL, labelColor } from "@/utils/filters";
import { BackButton } from "@/components/back-button";
import { UnderlinePanels } from "@primer/react/experimental";
import { toast } from "sonner";
import { useInView } from "react-intersection-observer";
import pusherClient from "@/utils/pusher";
import { useUser } from "@clerk/nextjs";

export default function EndpointShow() {
  const pathName = usePathname();
  const { ref, inView } = useInView();
  const user = useUser();

  const {
    showProject,
    showProjectEndpoint,
    deleteProjectEndpointRequest,
    indexProjectEndpointRequests,
  } = useAppStore((state) => state);
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
  const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false);
  const [deleteRequestId, setDeleteRequestId] = useState<string | undefined>(
    undefined,
  );
  const [
    canLoadMoreProjectEndpointRequests,
    setCanLoadMoreProjectEndpointRequests,
  ] = useState<boolean>(false);
  const [streamingIsActive, setStreamingIsActive] = useState<boolean>(false);

  const requestLimit: number = 20;
  const onDeleteDialogClose = useCallback(
    () => setIsDeleteDialogOpen(false),
    [],
  );

  const projectId = pathName.split("/")[2];
  const projectEndpointId = pathName.split("/")[4];

  const requestTabs = ["Request Summary", "Request Headers", "Request Body"];

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
    if (user.user?.id) {
      const channel = pusherClient.subscribe(user.user?.id);
      channel.bind("project.endpoint.request", () => {
        loadProjectEndpointRequests(projectId, projectEndpointId);
        setStreamingIsActive(true);
      });
    }
    return () => {
      if (user.user?.id) {
        pusherClient.unsubscribe(user.user?.id);
      }
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [user, projectId, projectEndpointId]);

  const loadProjectEndpointRequests = useCallback(
    (projectId: string, projectEndpointId: string, prev?: string) => {
      setLoadingProjectEndpointRequests(false);
      indexProjectEndpointRequests(
        projectId,
        projectEndpointId,
        requestLimit,
        prev,
      )
        .then((requests: Array<EntitiesProjectEndpointRequest>) => {
          const values = new Map<string, EntitiesProjectEndpointRequest>(
            [...projectEndpointRequests, ...requests].map((request) => [
              request.id,
              request,
            ]),
          );
          setProjectEndpointRequests(Array.from(values.values()));
          setCanLoadMoreProjectEndpointRequests(
            requests.length === requestLimit,
          );
        })
        .finally(() => {
          setLoadingProjectEndpointRequests(false);
        });
    },
    [indexProjectEndpointRequests, projectEndpointRequests],
  );

  useEffect(() => {
    loadProjectEndpointRequests(projectId, projectEndpointId);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [projectId, projectEndpointId]);

  useEffect(() => {
    if (inView && projectEndpointRequests.length > 0) {
      loadProjectEndpointRequests(
        projectId,
        projectEndpointId,
        projectEndpointRequests[projectEndpointRequests.length - 1].id,
      );
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [inView]);

  useEffect(() => {
    const id = setInterval(() => {
      if (streamingIsActive) {
        loadProjectEndpointRequests(projectId, projectEndpointId);
        setStreamingIsActive(false);
      }
    }, 5000);
    return () => clearInterval(id);
  }, [
    projectId,
    projectEndpointId,
    streamingIsActive,
    loadProjectEndpointRequests,
  ]);

  const onDeleteProjectEndpointRequest = async (event: MouseEvent) => {
    event.preventDefault();
    setLoadingProjectEndpointRequests(true);
    deleteProjectEndpointRequest(projectId, projectEndpointId, deleteRequestId!)
      .then(() => {
        setIsDeleteDialogOpen(false);
        setDeleteRequestId(undefined);
        setProjectEndpointRequests(
          projectEndpointRequests.filter(
            (request) => request.id !== deleteRequestId,
          ),
        );
        toast.info("Request deleted successfully.");
      })
      .finally(() => {
        setLoadingProjectEndpointRequests(false);
      });
  };

  const getHeaders = (
    headers: string | null,
  ): Array<Record<string, string>> => {
    const result = new Array<Record<string, string>>();
    if (headers == null || headers == "") {
      return result;
    }
    return JSON.parse(headers);
  };

  const getBodyIndentedJSON = (body: string | null): string => {
    if (body == null || body == "") {
      return "";
    }
    return JSON.stringify(JSON.parse(body), null, 4);
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
        <Text sx={{ ml: 2 }}>HTTP Request Stream</Text>
        {streamingIsActive && (
          <Spinner sx={{ ml: 2, color: "accent.emphasis" }} size="small" />
        )}
      </Heading>
      {!loadingProjectEndpointRequests && (
        <Box
          sx={{
            mt: 2,
            mb: 4,
            "> *": {
              borderWidth: 1,
              borderColor: "border.default",
              borderStyle: "solid",
              borderBottomWidth: 0,
              padding: 3,
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
                <Button
                  leadingVisual={TrashIcon}
                  onClick={(event) => {
                    event.preventDefault();
                    setDeleteRequestId(request.id);
                    setIsDeleteDialogOpen(true);
                  }}
                  sx={{ ml: "auto" }}
                  variant={"danger"}
                >
                  Delete
                </Button>
              </Box>
              <Text sx={{ color: "fg.muted", fontSize: "small" }}>
                Received{" "}
                <RelativeTime
                  date={new Date(request.created_at)}
                  noTitle={true}
                />
              </Text>
              <Box sx={{ mt: 2, display: "flex" }}>
                <UnderlinePanels aria-label="HTTP request details">
                  {requestTabs.map((tab: string, index: number) => (
                    <UnderlinePanels.Tab
                      key={index}
                      counter={undefined}
                      aria-selected={index === 0 ? true : undefined}
                    >
                      {tab}
                    </UnderlinePanels.Tab>
                  ))}
                  <UnderlinePanels.Panel key={0}>
                    <Box sx={{ mt: 2 }}>
                      <table>
                        <tbody>
                          <tr>
                            <td style={{ paddingTop: 16, paddingRight: 16 }}>
                              <Text weight={"semibold"} size={"medium"}>
                                Request ID
                              </Text>
                            </td>
                            <td style={{ minWidth: 400, paddingTop: 16 }}>
                              <TextInput
                                value={request.id}
                                readOnly={true}
                                block={true}
                                size={"small"}
                              />
                            </td>
                            <td style={{ paddingTop: 16, paddingLeft: 8 }}>
                              <CopyButton data={request.id} />
                            </td>
                          </tr>
                          <tr>
                            <td style={{ paddingTop: 16, paddingRight: 16 }}>
                              <Text weight={"semibold"} size={"medium"}>
                                IP Address
                              </Text>
                            </td>
                            <td style={{ paddingTop: 16 }}>
                              <TextInput
                                value={request.request_ip_address}
                                readOnly={true}
                                block={true}
                                size={"small"}
                              />
                            </td>
                            <td style={{ paddingTop: 16, paddingLeft: 8 }}>
                              <CopyButton data={request.request_ip_address} />
                            </td>
                          </tr>
                          <tr>
                            <td style={{ paddingTop: 16, paddingRight: 16 }}>
                              <Text weight={"semibold"} size={"medium"}>
                                Received At
                              </Text>
                            </td>
                            <td style={{ paddingTop: 16 }}>
                              <TextInput
                                value={request.created_at}
                                readOnly={true}
                                block={true}
                                size={"small"}
                              />
                            </td>
                            <td style={{ paddingTop: 16, paddingLeft: 8 }}>
                              <CopyButton data={request.created_at} />
                            </td>
                          </tr>
                        </tbody>
                      </table>
                    </Box>
                  </UnderlinePanels.Panel>
                  <UnderlinePanels.Panel key={1}>
                    <Box sx={{ mt: 2 }}>
                      <table>
                        <tbody>
                          {getHeaders(request.request_headers).map(
                            (header: Record<string, string>) => (
                              <tr
                                key={
                                  Object.keys(header)[0] +
                                  header[Object.keys(header)[0]]
                                }
                              >
                                <td
                                  style={{ paddingTop: 16, paddingRight: 16 }}
                                >
                                  <Text weight={"semibold"} size={"medium"}>
                                    {Object.keys(header)[0]}
                                  </Text>
                                </td>
                                <td style={{ paddingTop: 16, minWidth: 400 }}>
                                  <TextInput
                                    value={header[Object.keys(header)[0]]}
                                    readOnly={true}
                                    block={true}
                                    size={"small"}
                                  />
                                </td>
                                <td style={{ paddingTop: 16, paddingLeft: 8 }}>
                                  <CopyButton
                                    data={header[Object.keys(header)[0]]}
                                  />
                                </td>
                              </tr>
                            ),
                          )}
                        </tbody>
                      </table>
                    </Box>
                  </UnderlinePanels.Panel>
                  <UnderlinePanels.Panel key={2}>
                    <Box sx={{ mt: 2 }}>
                      <table>
                        <tbody>
                          <tr>
                            <td style={{ paddingTop: 16, minWidth: 400 }}>
                              <Textarea
                                value={getBodyIndentedJSON(
                                  request.request_body,
                                )}
                                readOnly={true}
                                block={true}
                                rows={20}
                              />
                            </td>
                            <td
                              style={{
                                paddingTop: 16,
                                paddingLeft: 8,
                                verticalAlign: "top",
                              }}
                            >
                              <CopyButton
                                data={getBodyIndentedJSON(request.request_body)}
                              />
                            </td>
                          </tr>
                        </tbody>
                      </table>
                    </Box>
                  </UnderlinePanels.Panel>
                </UnderlinePanels>
              </Box>
            </div>
          ))}
        </Box>
      )}
      {canLoadMoreProjectEndpointRequests && (
        <Box ref={ref} sx={{ mb: 2, color: "fg.muted" }}>
          <Spinner size={"small"} /> Loading more requests...
        </Box>
      )}
      {isDeleteDialogOpen && (
        <Dialog
          width={"large"}
          onClose={onDeleteDialogClose}
          title={`Delete HTTP Request`}
          footerButtons={[
            {
              disabled: loadingProjectEndpointRequests,
              content: "Close",
              onClick: onDeleteDialogClose,
            },
            {
              buttonType: "danger",
              block: true,
              loading: loadingProjectEndpointRequests,
              disabled: loadingProjectEndpointRequests,
              content: "Delete this Request",
              onClick: onDeleteProjectEndpointRequest,
            },
          ]}
        >
          <div>
            <p>
              Confirm that you want to delete the HTTP request because this is a
              permanent action and it cannot be reversed.
            </p>
          </div>
          <Box sx={{ mt: 2 }}>
            <Text sx={{ color: "fg.muted" }}>{deleteRequestId}</Text>
          </Box>
        </Dialog>
      )}
    </Box>
  );
}
