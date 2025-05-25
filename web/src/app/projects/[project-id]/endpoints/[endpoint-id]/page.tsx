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
  IconButton,
} from "@primer/react";
import { usePathname } from "next/navigation";
import React, { MouseEvent, useCallback, useEffect, useState } from "react";
import { useAppStore } from "@/store/provider";
import {
  EntitiesProject,
  EntitiesProjectEndpoint,
  EntitiesProjectEndpointRequest,
  RepositoriesTimeSeriesData,
} from "@/api/model";
import { MirrorIcon, TrashIcon, GraphIcon } from "@primer/octicons-react";
import { CopyButton } from "@/components/copy-button";
import { getEndpointURL, labelColor } from "@/utils/filters";
import { BackButton } from "@/components/back-button";
import { UnderlinePanels } from "@primer/react/experimental";
import { toast } from "sonner";
import { useInView } from "react-intersection-observer";
import pusherClient from "@/utils/pusher";
import { useUser } from "@clerk/nextjs";
import { Line } from "react-chartjs-2";
import "chartjs-adapter-date-fns";
import {
  Chart,
  CategoryScale,
  LinearScale,
  LineElement,
  Title,
  PointElement,
  Tooltip,
  TimeScale,
  Legend,
  ChartData,
  Point,
} from "chart.js";

Chart.register(
  CategoryScale,
  LinearScale,
  LineElement,
  PointElement,
  TimeScale,
  Title,
  Tooltip,
  Legend,
);

export default function EndpointShow() {
  const pathName = usePathname();
  const { ref, inView } = useInView();
  const user = useUser();

  const {
    showProject,
    showProjectEndpoint,
    deleteProjectEndpointRequest,
    indexProjectEndpointRequests,
    getProjectEndpointTraffic,
  } = useAppStore((state) => state);
  const [loadingEndpoint, setLoadingEndpoint] = useState<boolean>(false);
  const [loadingProjectEndpointRequests, setLoadingProjectEndpointRequests] =
    useState<boolean>(true);
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
  const [loadingTraffic, setLoadingTraffic] = useState<boolean>(true);
  const [chartData, setChartData] = useState<ChartData | undefined>(undefined);

  const requestLimit: number = 20;
  const onDeleteDialogClose = useCallback(
    () => setIsDeleteDialogOpen(false),
    [],
  );

  const projectId = pathName.split("/")[2];
  const projectEndpointId = pathName.split("/")[4];

  const requestTabs = ["Summary", "Headers", "Body"];

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
      setLoadingProjectEndpointRequests(true);
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
          setProjectEndpointRequests(
            Array.from(values.values()).sort((a, b) => {
              return (
                new Date(b.created_at).getTime() -
                new Date(a.created_at).getTime()
              );
            }),
          );
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

  useEffect(() => {
    setLoadingTraffic(true);
    getProjectEndpointTraffic(projectId, projectEndpointId)
      .then((timeSeries: RepositoriesTimeSeriesData[]) => {
        console.log(timeSeries);
        setChartData({
          labels: [],
          datasets: [
            {
              label: "HTTP Request Count",
              data: timeSeries.map((dataPoint) => {
                return {
                  x: new Date(dataPoint.timestamp),
                  y: dataPoint.count,
                } as unknown as Point;
              }),
              borderColor: "rgb(9, 105, 218)",
              backgroundColor: "rgba(9, 105, 218, 0.2)",
              fill: true,
            },
          ],
        });
      })
      .finally(() => {
        setLoadingTraffic(false);
      });
  }, [getProjectEndpointTraffic, projectId, projectEndpointId]);

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
    try {
      return JSON.stringify(JSON.parse(body), null, 4);
    } catch (error) {
      console.error(error);
      return body;
    }
  };

  const secondsSince = (date: string): number => {
    return (new Date().getTime() - new Date(date).getTime()) / 1000;
  };

  const getCurlCode = (endpoint: EntitiesProjectEndpoint) => {
    return `curl --header \"Content-Type: application/json\" \\\n\t--request ${endpoint.request_method === "ANY" ? "POST" : endpoint.request_method} ${endpoint.request_method != "GET" ? "\\\n\t--data '{\"success\":true}' \\\n\t" : ""} \"${getEndpointURL(endpoint)}\"`;
  };

  const chartOptions = {
    plugins: {
      legend: {
        display: false,
      },
    },
    maintainAspectRatio: false,
    responsive: true,
    scales: {
      x: {
        type: "time",
      },
    },
  };

  return (
    <Box
      sx={{
        maxWidth: "xlarge",
        mx: "auto",
        mt: 4,
        px: 2,
        minHeight: "calc(100vh - 200px)",
      }}
    >
      <PageHeader role="banner" aria-label="Project details">
        <PageHeader.TitleArea
          aria-label={"Project endpoint"}
          variant={{ narrow: "medium", regular: "large" }}
        >
          {project && projectEndpoint && (
            <PageHeader.Title>
              <Box sx={{ display: "flex", alignItems: "flexStart" }}>
                <Label
                  sx={{
                    mt: 3,
                    display: ["none", "flex"],
                    color: labelColor(projectEndpoint.request_method),
                  }}
                >
                  {projectEndpoint.request_method}
                </Label>
                <Text
                  sx={{
                    ml: 1,
                    mr: 1,
                    fontWeight: "bold",
                    wordBreak: "break-all",
                  }}
                >
                  {getEndpointURL(projectEndpoint)}
                </Text>
                <CopyButton
                  size={"medium"}
                  sx={{ display: ["none", "block"], mt: 2 }}
                  data={getEndpointURL(projectEndpoint)}
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
            <Text as={"p"} sx={{ fontSize: 1, color: "fg.muted" }}>
              {project?.description}
            </Text>
          </PageHeader.Description>
        )}
        <PageHeader.Actions sx={{ display: ["none", "flex"] }}>
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

      <Box display={["none", "none", "block"]}>
        <Heading sx={{ mt: 4 }} as={"h2"} variant={"medium"}>
          <GraphIcon size={24} />
          <Text sx={{ ml: 2 }}>Endpoint Traffic Insights</Text>
        </Heading>
        {!loadingTraffic && (
          <Box
            sx={{
              marginTop: 1,
              borderStyle: "solid",
              borderWidth: 1,
              padding: 2,
              borderRadius: 4,
              borderColor: "border.default",
            }}
          >
            {/* @ts-expect-error chart types are not consistent */}
            <Line options={chartOptions} data={chartData} />
          </Box>
        )}
      </Box>

      <Heading as="h2" sx={{ mt: 32 }} variant="medium">
        <MirrorIcon size={24} />
        <Text sx={{ ml: 2 }}>HTTP Request Stream</Text>
        {streamingIsActive && (
          <Spinner sx={{ ml: 2, color: "accent.emphasis" }} size="small" />
        )}
      </Heading>
      {loadingProjectEndpointRequests &&
        projectEndpointRequests.length === 0 && (
          <Spinner sx={{ ml: 2, color: "accent.emphasis" }} size="large" />
        )}
      {!loadingProjectEndpointRequests &&
        projectEndpoint &&
        projectEndpointRequests.length === 0 && (
          <React.Fragment>
            <Text as={"p"} sx={{ color: "fg.muted" }}>
              The endpoint{" "}
              <BranchName>{getEndpointURL(projectEndpoint)}</BranchName> does
              not have any have any HTTP requests. You can make an http request
              using the <BranchName>curl</BranchName> command below
            </Text>
            <Box sx={{ display: "flex", mt: 2 }}>
              <Textarea
                value={getCurlCode(projectEndpoint)}
                readOnly={true}
                block={true}
                rows={4}
              />
              <CopyButton sx={{ ml: 2 }} data={getCurlCode(projectEndpoint)} />
            </Box>
          </React.Fragment>
        )}
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
                <Text
                  sx={{
                    ml: 1,
                    mr: 1,
                    fontWeight: "bold",
                    wordBreak: "break-all",
                  }}
                >
                  {request.request_url}
                </Text>
                <CopyButton
                  sx={{ display: ["none", "inline-block"] }}
                  data={request.request_url}
                />
                <IconButton
                  onClick={(event) => {
                    event.preventDefault();
                    setDeleteRequestId(request.id);
                    setIsDeleteDialogOpen(true);
                  }}
                  sx={{ ml: "auto", px: 2, display: ["inline-block", "none"] }}
                  variant={"danger"}
                  icon={TrashIcon}
                  aria-label="Delete"
                />
                <Button
                  leadingVisual={TrashIcon}
                  onClick={(event) => {
                    event.preventDefault();
                    setDeleteRequestId(request.id);
                    setIsDeleteDialogOpen(true);
                  }}
                  sx={{ ml: "auto", display: ["none", "inline-block"] }}
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
                {secondsSince(request.created_at) < 120 && (
                  <Label
                    size={"small"}
                    sx={{ color: "#4c8341", ml: 1, fontSize: "small" }}
                  >
                    NEW
                  </Label>
                )}
              </Text>
              <Box sx={{ mt: 2, display: "flex" }}>
                <UnderlinePanels aria-label="HTTP request details">
                  {requestTabs.map((tab: string, index: number) => (
                    <UnderlinePanels.Tab
                      key={index}
                      counter={undefined}
                      aria-selected={index === 0 ? true : undefined}
                    >
                      <p>
                        <Text sx={{ display: ["none", "inline-block"], mr: 1 }}>
                          Request
                        </Text>{" "}
                        {tab}
                      </p>
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
                            <td style={{ width: 400, paddingTop: 16 }}>
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
                                <td style={{ paddingTop: 16, width: 400 }}>
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
                            <td style={{ paddingTop: 16, width: 400 }}>
                              <Textarea
                                value={getBodyIndentedJSON(
                                  request.request_body,
                                )}
                                readOnly={true}
                                block={true}
                                rows={10}
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
