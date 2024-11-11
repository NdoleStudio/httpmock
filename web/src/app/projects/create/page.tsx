"use client";

import {
  Box,
  Button,
  FormControl,
  Heading,
  Text,
  Textarea,
  TextInput,
  useTheme,
} from "@primer/react";
import { ArrowLeftIcon } from "@primer/octicons-react";
import { useRouter } from "next/navigation";

export default function ProjectCreate() {
  const router = useRouter();

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
        <Button
          onClick={() => router.push("/")}
          sx={{ mb: 3 }}
          leadingVisual={ArrowLeftIcon}
        >
          Back to dashboard
        </Button>
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
            <TextInput block={true} size={"large"} />
          </FormControl>
          <FormControl sx={{ mt: 4 }} required={true}>
            <FormControl.Label>Project Description</FormControl.Label>
            <Textarea block={true} rows={2} />
          </FormControl>
          <Button sx={{ mt: 4 }} variant={"primary"}>
            Create Project
          </Button>
        </Box>
      </Box>
    </Box>
  );
}
