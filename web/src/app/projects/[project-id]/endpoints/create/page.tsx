"use client";

import {
  Box,
  Button,
  FormControl,
  Heading,
  Text,
  Textarea,
  TextInput,
} from "@primer/styled-react";
import { Select } from "@primer/react";
import { usePathname, useRouter } from "next/navigation";
import { MouseEvent, useEffect, useState } from "react";
import { ErrorMessages } from "@/utils/errors";
import { BackButton } from "@/components/back-button";
import { useAppStore } from "@/store/provider";
import { EntitiesProject, EntitiesProjectEndpoint } from "@/api/model";

export default function EndpointsCreate() {
  const router = useRouter();
  const pathName = usePathname();
  const { showProject, storeProjectEndpoint } = useAppStore((state) => state);

  const [errorMessages, setErrorMessages] = useState<ErrorMessages>(
    ErrorMessages.create(),
  );
  const [loading, setLoading] = useState<boolean>(false);
  const [project, setProject] = useState<EntitiesProject | undefined>(
    undefined,
  );
  const [requestMethod, setRequestMethod] = useState<string>("GET");
  const [responseCode, setResponseCode] = useState<number>(200);
  const [requestPath, setRequestPath] = useState<string>("");
  const [responseBody, setResponseBody] = useState<string>("");
  const [responseHeaders, setResponseHeaders] = useState<string>("");
  const [responseDelayInMilliseconds, setResponseDelayInMilliseconds] =
    useState<number>(0);
  const [description, setDescription] = useState<string>("");

  const projectId = pathName.split("/")[2];

  const onEndpointCreate = (event: MouseEvent<HTMLButtonElement>) => {
    event.preventDefault();

    setLoading(true);
    setErrorMessages(ErrorMessages.create());

    storeProjectEndpoint(projectId, {
      request_method: requestMethod,
      request_path: requestPath,
      response_code: responseCode,
      response_body: responseBody,
      response_headers: responseHeaders,
      response_delay_in_milliseconds: responseDelayInMilliseconds,
      description: description,
    })
      .then((projectEndpoint: EntitiesProjectEndpoint) => {
        router.push(`/projects/${projectEndpoint.project_id}`);
      })
      .catch((errorMessages: ErrorMessages) => {
        setErrorMessages(errorMessages);
      })
      .finally(() => {
        setLoading(false);
      });
  };

  useEffect(() => {
    showProject(projectId)
      .then((project: EntitiesProject) => {
        setProject(project);
      })
      .catch((errorMessages: ErrorMessages) => {
        setErrorMessages(errorMessages);
      });
  }, [showProject, projectId]);

  return (
    <Box
      sx={{
        mt: 6,
        display: "flex",
        justifyContent: "center",
        minHeight: "calc(100vh - 200px)",
      }}
    >
      <Box sx={{ maxWidth: "100%", width: "small", mb: 16 }}>
        <BackButton href={`/projects/${projectId}`}></BackButton>
        <Box
          sx={{
            backgroundColor: "canvas.inset",
            borderWidth: 1,
            borderStyle: "solid",
            borderColor: "border.default",
            borderRadius: 2,
            p: 3,
          }}
        >
          <Heading as={"h2"}>Create Mock Endpoint</Heading>
          <FormControl sx={{ mt: 4 }} required={true} disabled={loading}>
            <FormControl.Label>Request Method</FormControl.Label>
            <FormControl.Caption>
              Use ANY if you want to match all HTTP methods (GET, POST, DELETE
              etc)
            </FormControl.Caption>
            <Select
              validationStatus={
                errorMessages.has("request_method") ? "error" : undefined
              }
              value={requestMethod}
              onChange={(event) => {
                setRequestMethod(event.target.value);
              }}
              block={true}
              size={"large"}
            >
              <Select.Option value="GET">GET</Select.Option>
              <Select.Option value="POST">POST</Select.Option>
              <Select.Option value="PUT">PUT</Select.Option>
              <Select.Option value="PATCH">PATCH</Select.Option>
              <Select.Option value="DELETE">DELETE</Select.Option>
              <Select.Option value="OPTIONS">OPTIONS</Select.Option>
              <Select.Option value="ANY">ANY (*)</Select.Option>
            </Select>
            {errorMessages.has("request_method") && (
              <FormControl.Validation variant="error">
                {errorMessages.first("request_method")}
              </FormControl.Validation>
            )}
          </FormControl>
          <FormControl sx={{ mt: 4 }} required={true} disabled={loading}>
            <FormControl.Label>Request Path</FormControl.Label>
            <FormControl.Caption>
              Your full URL will look like{" "}
              <Text sx={{ color: "accent.emphasis", fontWeight: "bold" }}>
                https://{project?.subdomain}.httmock.dev{requestPath}
              </Text>
            </FormControl.Caption>
            <TextInput
              placeholder={"e.g. /api/v1/users"}
              validationStatus={
                errorMessages.has("request_path") ? "error" : undefined
              }
              value={requestPath}
              onChange={(event) => {
                setRequestPath(event.target.value);
              }}
              block={true}
              size={"large"}
            />
            {errorMessages.has("request_path") && (
              <FormControl.Validation variant="error">
                {errorMessages.first("request_path")}
              </FormControl.Validation>
            )}
          </FormControl>
          <FormControl sx={{ mt: 4 }} required={true} disabled={loading}>
            <FormControl.Label>Response Code</FormControl.Label>
            <FormControl.Caption>
              HTTP status code to return in the response
            </FormControl.Caption>
            <TextInput
              placeholder={"e.g. 200"}
              validationStatus={
                errorMessages.has("response_code") ? "error" : undefined
              }
              value={responseCode}
              type={"number"}
              onChange={(event) => {
                setResponseCode(Number.parseInt(event.target.value));
              }}
              block={true}
              size={"large"}
            />
            {errorMessages.has("response_code") && (
              <FormControl.Validation variant="error">
                {errorMessages.first("response_code")}
              </FormControl.Validation>
            )}
          </FormControl>
          <FormControl sx={{ mt: 4 }} disabled={loading}>
            <FormControl.Label>Response Body</FormControl.Label>
            <FormControl.Caption>
              The response body can be any valid JSON, XML, or plain text
            </FormControl.Caption>
            <Textarea
              placeholder={
                'e.g \n{\n\t"message": "Hello World",\n\t"status": 200\n}'
              }
              validationStatus={
                errorMessages.has("response_body") ? "error" : undefined
              }
              value={responseBody}
              onChange={(event) => {
                setResponseBody(event.target.value);
              }}
              block={true}
              rows={10}
            />
            {errorMessages.has("response_body") && (
              <FormControl.Validation variant="error">
                {errorMessages.first("response_body")}
              </FormControl.Validation>
            )}
          </FormControl>
          <FormControl sx={{ mt: 4 }} disabled={loading}>
            <FormControl.Label>Response Headers</FormControl.Label>
            <FormControl.Caption>
              This should be a JSON array of headers that will be returned with
              the HTTP response
            </FormControl.Caption>
            <Textarea
              placeholder={
                'e.g\n[\n\t{"Content-Type": "application/json"},\n\t{"x-request-id": "abc-7c865c14b444"}\n]'
              }
              validationStatus={
                errorMessages.has("response_headers") ? "error" : undefined
              }
              value={responseHeaders}
              onChange={(event) => {
                setResponseHeaders(event.target.value);
              }}
              block={true}
              rows={5}
            />
            {errorMessages.has("response_headers") && (
              <FormControl.Validation variant="error">
                {errorMessages.first("response_headers")}
              </FormControl.Validation>
            )}
          </FormControl>
          <FormControl sx={{ mt: 4 }} disabled={loading}>
            <FormControl.Label>
              Response Delay in Milliseconds
            </FormControl.Label>
            <FormControl.Caption>
              The time in milliseconds to wait before sending the HTTP response
            </FormControl.Caption>
            <TextInput
              placeholder={"e.g. 1000"}
              validationStatus={
                errorMessages.has("response_delay_in_milliseconds")
                  ? "error"
                  : undefined
              }
              value={responseDelayInMilliseconds}
              type={"number"}
              onChange={(event) => {
                setResponseDelayInMilliseconds(
                  Number.parseInt(event.target.value),
                );
              }}
              block={true}
              size={"large"}
            />
            {errorMessages.has("delay_in_milliseconds") && (
              <FormControl.Validation variant="error">
                {errorMessages.first("delay_in_milliseconds")}
              </FormControl.Validation>
            )}
          </FormControl>
          <FormControl sx={{ mt: 4 }} disabled={loading}>
            <FormControl.Label>Endpoint Description</FormControl.Label>
            <FormControl.Caption>
              Use the description field to add more context to your mock
              endpoint.
            </FormControl.Caption>
            <Textarea
              placeholder={"e.g This is a mock of the GitHub API"}
              validationStatus={
                errorMessages.has("description") ? "error" : undefined
              }
              value={description}
              onChange={(event) => {
                setDescription(event.target.value);
              }}
              block={true}
              rows={3}
            />
            {errorMessages.has("description") && (
              <FormControl.Validation variant="error">
                {errorMessages.first("description")}
              </FormControl.Validation>
            )}
          </FormControl>
          <Button
            loading={loading}
            disabled={loading}
            onClick={onEndpointCreate}
            sx={{ mt: 4 }}
            variant={"primary"}
          >
            Create Endpoint
          </Button>
        </Box>
      </Box>
    </Box>
  );
}
