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
  Heading,
  Label,
  Link,
} from "@primer/styled-react";
import { BranchName, RelativeTime } from "@primer/react";
import { usePathname, useRouter } from "next/navigation";
import { useCallback, useEffect, MouseEvent, useState } from "react";
import { useAppStore } from "@/store/provider";
import {
  EntitiesProject,
  EntitiesProjectEndpoint,
  RepositoriesTimeSeriesData,
} from "@/api/model";
import {
  GearIcon,
  LinkIcon,
  PencilIcon,
  PlusIcon,
  TrashIcon,
} from "@primer/octicons-react";
import { toast } from "sonner";
import { CopyButton } from "@/components/copy-button";
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

export default function ProjectShow() {
  const router = useRouter();
  const pathName = usePathname();
  const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false);
  const {
    showProject,
    deleteProject,
    indexProjectEndpoint,
    getProjectTraffic,
  } = useAppStore((state) => state);
  const [loading, setLoading] = useState<boolean>(false);
  const [endpoints, setEndpoints] = useState<EntitiesProjectEndpoint[]>([]);
  const [loadingEndpoints, setLoadingEndpoints] = useState<boolean>(false);
  const [loadingTraffic, setLoadingTraffic] = useState<boolean>(true);
  const [chartData, setChartData] = useState<ChartData | undefined>(undefined);
  const [project, setProject] = useState<EntitiesProject | undefined>(
    undefined,
  );

  const onDeleteDialogClose = useCallback(
    () => setIsDeleteDialogOpen(false),
    [],
  );

  const projectId = pathName.split("/")[2];

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
        time: {
          tooltipFormat: "LLL Lo, yyyy",
        },
      },
    },
  };

  useEffect(() => {
    setLoadingTraffic(true);
    getProjectTraffic(projectId)
      .then((timeSeries: RepositoriesTimeSeriesData[]) => {
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
  }, [getProjectTraffic, projectId]);

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
      style={{
        maxWidth: 1200, // xlarge
        marginLeft: "auto",
        marginRight: "auto",
        marginTop: 32, // mt: 4
        paddingLeft: 16, // px: 2
        paddingRight: 16,
        minHeight: "calc(100vh - 200px)",
      }}
    >
      <PageHeader role="banner" aria-label="Project details">
        <PageHeader.TitleArea variant={"large"}>
          {project && <PageHeader.Title>{project.name}</PageHeader.Title>}
          {!project && (
            <PageHeader.Title>
              <Spinner size="large" />
            </PageHeader.Title>
          )}
        </PageHeader.TitleArea>
        {project && (
          <PageHeader.Description>
            <div>
              {project.description && (
                <Text as={"p"} style={{ color: "#6e7781", marginBottom: 4 }}>
                  {project.description}
                </Text>
              )}
              <p>
                <BranchName>{project?.subdomain}.httpmock.dev</BranchName>
              </p>
            </div>
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
            <Text style={{ display: "inline-block" }}>New Endpoint</Text>
            <Text style={{ display: "none" }}>Endpoint</Text>
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
        style={{
          borderBottomWidth: 1,
          marginTop: 16,
          borderBottomStyle: "solid",
          borderColor: "#d0d7de",
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
          <Box style={{ marginTop: 16 }}>
            <Text style={{ color: "#6e7781" }}>
              {project?.subdomain}.httpmock.dev
            </Text>
          </Box>
        </Dialog>
      )}

      <Box display={["none", "none", "block"]}>
        <Heading style={{ marginTop: 64 }} as={"h2"} variant={"medium"}>
          Project Traffic Insights
        </Heading>
        {!loadingTraffic && (
          <Box
            style={{
              marginTop: 8,
              borderStyle: "solid",
              borderWidth: 1,
              padding: 16,
              borderRadius: 4,
              borderColor: "#d0d7de",
            }}
          >
            {/* @ts-expect-error chart types are not consistent */}
            <Line options={chartOptions} data={chartData} />
          </Box>
        )}
      </Box>
      <div>
        <Heading as="h2" style={{ marginTop: 256 }} variant="medium">
          <LinkIcon size={24} />
          <Text style={{ marginLeft: 8 }}>Endpoints</Text>
        </Heading>
      </div>
      {loadingEndpoints && (
        <div>
          <Spinner size="large" />
        </div>
      )}
      {!loadingEndpoints && endpoints.length === 0 && (
        <p>
          <Text style={{ color: "#6e7781" }}>
            This project does not have any have any mocked HTTP endpoints.{" "}
          </Text>
          <Link
            style={{ color: "#0969da", cursor: "pointer" }}
            onClick={() =>
              router.push(`/projects/${projectId}/endpoints/create`)
            }
          >
            Create a new endpoint
          </Link>
          <Text style={{ color: "#6e7781" }}> to start mocking your APIs</Text>
        </p>
      )}
      {!loadingEndpoints && (
        <Box
          style={{
            marginTop: 16,
            // Note: Complex child selectors and responsive padding are omitted for inline styles
          }}
        >
          {endpoints.map((endpoint) => (
            <div key={endpoint.id}>
              <Box style={{ display: "flex", alignItems: "baseline" }}>
                <Label
                  style={{
                    display: "flex",
                    color: getLabelColor(endpoint.request_method),
                  }}
                >
                  {endpoint.request_method}
                </Label>
                <Text
                  style={{ marginLeft: 4, marginRight: 4, fontWeight: "bold" }}
                >
                  {getEndpointURL(endpoint.request_path)}
                </Text>
                <CopyButton data={getEndpointURL(endpoint.request_path)} />
              </Box>
              <Box style={{ marginTop: 16, display: "flex" }}>
                <Button
                  onClick={() =>
                    router.push(
                      `/projects/${projectId}/endpoints/${endpoint.id}`,
                    )
                  }
                  size={"small"}
                  count={endpoint.request_count}
                  variant="default"
                  style={{ marginRight: 8 }}
                >
                  HTTP Requests
                </Button>
                <Button
                  onClick={() =>
                    router.push(
                      `/projects/${projectId}/endpoints/${endpoint.id}/edit`,
                    )
                  }
                  size={"small"}
                  variant="default"
                  style={{ marginRight: 8 }}
                  leadingVisual={GearIcon}
                >
                  Manage
                </Button>
                <div>
                  <Text
                    size="small"
                    style={{ marginRight: 4, color: "#6e7781" }}
                  >
                    Updated
                  </Text>
                  <RelativeTime
                    style={{ color: "#6e7781", fontSize: "small" }}
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
