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

export default function ProjectCreate() {
  const router = useRouter();
  const [errorMessages, setErrorMessages] = useState<ErrorMessages>(
    ErrorMessages.create(),
  );
  const [loading, setLoading] = useState<boolean>(false);
  const [projectName, setProjectName] = useState<string>("");
  const [projectDescription, setProjectDescription] = useState<string>("");

  const createProject = (event: MouseEvent<HTMLButtonElement>) => {
    event.preventDefault();
    setErrorMessages(ErrorMessages.create());
    setLoading(true);
  };

  return (
    <Box
      sx={{
        mt: 6,
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
      }}
    >
      <Box>
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
          <Heading>Create Project</Heading>
          <Text sx={{ color: "fg.muted" }}>
            Your mocked endpoints are grouped into projects for better
            organization.
          </Text>
          <FormControl sx={{ mt: 4 }} required={true}>
            <FormControl.Label>Project Name</FormControl.Label>
            <TextInput
              disabled={loading}
              value={projectName}
              onChange={(event) => {
                setProjectName(event.target.value);
              }}
              block={true}
              size={"large"}
            />
          </FormControl>
          <FormControl sx={{ mt: 4 }} required={true}>
            <FormControl.Label>Project Description</FormControl.Label>
            <Textarea
              disabled={loading}
              value={projectDescription}
              onChange={(event) => {
                setProjectDescription(event.target.value);
              }}
              block={true}
              rows={2}
            />
          </FormControl>
          <Button
            loading={loading}
            disabled={loading}
            onClick={createProject}
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
