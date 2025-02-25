"use client";

import {
  Box,
  Button,
  FormControl,
  Heading,
  Text,
  Textarea,
  TextInput,
} from "@primer/react";
import { useRouter } from "next/navigation";
import { MouseEvent, useState } from "react";
import { ErrorMessages } from "@/utils/errors";
import { BackButton } from "@/components/back-button";
import { useAppStore } from "@/store/provider";
import { EntitiesProject } from "@/api/model";

export default function ProjectCreate() {
  const router = useRouter();
  const { storeProject } = useAppStore((state) => state);

  const [errorMessages, setErrorMessages] = useState<ErrorMessages>(
    ErrorMessages.create(),
  );
  const [loading, setLoading] = useState<boolean>(false);
  const [projectName, setProjectName] = useState<string>("");
  const [projectSubdomain, setProjectSubdomain] = useState<string>("");
  const [projectDescription, setProjectDescription] = useState<string>("");

  const onCreateProject = (event: MouseEvent<HTMLButtonElement>) => {
    event.preventDefault();

    setLoading(true);
    setErrorMessages(ErrorMessages.create());

    storeProject({
      name: projectName,
      subdomain: projectSubdomain,
      description: projectDescription,
    })
      .then((project: EntitiesProject) => {
        router.push(`/projects/${project.id}`);
      })
      .catch((errorMessages: ErrorMessages) => {
        setErrorMessages(errorMessages);
      })
      .finally(() => {
        setLoading(false);
      });
  };

  return (
    <Box
      sx={{
        mt: 4,
        px: 2,
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
      }}
    >
      <Box sx={{ maxWidth: "500px" }}>
        <BackButton href={"/"}></BackButton>
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
          <Heading as={"h2"}>Create Project</Heading>
          <Text sx={{ color: "fg.muted" }}>
            Your mocked endpoints are grouped into projects for better
            organization.
          </Text>
          <FormControl sx={{ mt: 4 }} required={true} disabled={loading}>
            <FormControl.Label>Project Subdomain</FormControl.Label>
            <TextInput
              trailingVisual={
                <Text size="large" weight="semibold">
                  .httpmock.dev
                </Text>
              }
              validationStatus={
                errorMessages.has("subdomain") ? "error" : undefined
              }
              value={projectSubdomain}
              onChange={(event) => {
                setProjectSubdomain(event.target.value);
              }}
              block={true}
              size={"large"}
            />
            {errorMessages.has("subdomain") && (
              <FormControl.Validation variant="error">
                {errorMessages.first("subdomain")}
              </FormControl.Validation>
            )}
          </FormControl>
          <FormControl sx={{ mt: 4 }} required={true} disabled={loading}>
            <FormControl.Label>Project Name</FormControl.Label>
            <TextInput
              validationStatus={errorMessages.has("name") ? "error" : undefined}
              value={projectName}
              onChange={(event) => {
                setProjectName(event.target.value);
              }}
              block={true}
              size={"large"}
            />
            {errorMessages.has("name") && (
              <FormControl.Validation variant="error">
                {errorMessages.first("name")}
              </FormControl.Validation>
            )}
          </FormControl>
          <FormControl sx={{ mt: 4 }} disabled={loading}>
            <FormControl.Label>Project Description</FormControl.Label>
            <Textarea
              validationStatus={
                errorMessages.has("description") ? "error" : undefined
              }
              value={projectDescription}
              onChange={(event) => {
                setProjectDescription(event.target.value);
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
            onClick={onCreateProject}
            sx={{ mt: 4 }}
            variant={"primary"}
          >
            Create Project
          </Button>
        </Box>
      </Box>
    </Box>
  );
}
