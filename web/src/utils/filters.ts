import { EntitiesProject } from "@/api/model";

export const capitalize = function (value: string | null) {
  if (!value) {
    return "";
  }
  value = value.toString();
  return value.charAt(0).toUpperCase() + value.slice(1);
};

export const labelColor = (requestMethod: string): string => {
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

export const getEndpointURL = (
  project: EntitiesProject | undefined,
  path: string,
): string => {
  return `https://${project?.subdomain}.httpmock.dev${path === "/" ? "" : path}`;
};
