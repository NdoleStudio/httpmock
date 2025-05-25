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
  Link,
} from "@primer/react";
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
      sx={{
        maxWidth: "xlarge",
        mx: "auto",
        mt: 4,
        px: 2,
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
                <Text as={"p"} sx={{ color: "fg.muted", mb: 1 }}>
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
            <Text sx={{ display: ["none", "inline-block"] }}>New Endpoint</Text>
            <Text sx={{ display: ["inline-block", "none"] }}>Endpoint</Text>
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

      <Box display={["none", "none", "block"]}>
        <Heading sx={{ mt: 4 }} as={"h2"} variant={"medium"}>
          Project Traffic Insights
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
      {!loadingEndpoints && endpoints.length === 0 && (
        <p>
          <Text sx={{ color: "fg.muted" }}>
            This project does not have any have any mocked HTTP endpoints.{" "}
          </Text>
          <Link
            sx={{ color: "accent.emphasis", cursor: "pointer" }}
            onClick={() =>
              router.push(`/projects/${projectId}/endpoints/create`)
            }
          >
            Create a new endpoint
          </Link>
          <Text sx={{ color: "fg.muted" }}> to start mocking your APIs</Text>
        </p>
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
              padding: [2, 4],
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
                <Label
                  sx={{
                    display: ["none", "flex"],
                    color: getLabelColor(endpoint.request_method),
                  }}
                >
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
                  count={endpoint.request_count}
                  variant="default"
                  sx={{ mr: 2 }}
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
