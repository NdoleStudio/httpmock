"use client";

import * as React from "react";
import { Button } from "@primer/styled-react";
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
      style={{ marginBottom: 24 }}
      leadingVisual={ArrowLeftIcon}
    >
      Go Back
    </Button>
  );
}
