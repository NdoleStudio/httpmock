"use client";

import * as React from "react";
import { Button } from "@primer/react";
import { ArrowLeftIcon } from "@primer/octicons-react";
import { useRouter } from "next/navigation";

type BackButtonProps = {
  href: string;
};

export function BackButton({ ...props }: BackButtonProps) {
  const router = useRouter();
  return (
    <Button
      onClick={() => router.push(props.href)}
      sx={{ mb: 3 }}
      leadingVisual={ArrowLeftIcon}
    >
      Back to dashboard
    </Button>
  );
}
