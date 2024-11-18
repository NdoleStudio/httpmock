"use client";

import { MouseEvent } from "react";
import { Button } from "@primer/react";
import { CheckIcon, CopyIcon } from "@primer/octicons-react";
import { useState } from "react";

type CopyButtonProps = {
  data: string;
};

export function CopyButton({ ...props }: CopyButtonProps) {
  const [copied, setCopied] = useState<boolean>(false);

  const onClick = (event: MouseEvent) => {
    navigator.clipboard.writeText(props.data).then(() => {
      setCopied(true);
      setTimeout(() => setCopied(false), 4000);
    });
  };
  return (
    <Button onClick={onClick} sx={{}} size={"small"} variant="invisible">
      {!copied && <CopyIcon size={"small"} />}
      {copied && <CheckIcon fill={"#3fb950"} size={"small"} />}
    </Button>
  );
}
