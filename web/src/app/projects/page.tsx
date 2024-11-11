"use client";

import { redirect } from "next/navigation";

export default function ProjectIndex() {
  return redirect("/projects/create");
}
