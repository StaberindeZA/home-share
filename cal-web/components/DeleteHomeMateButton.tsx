"use client";

import { removeHomeMate } from "@/actions";
import { useTransition } from "react";

interface DeleteHomeMateButtonProps {
  homeSlug: string;
  mateEmail: string;
}

export default function DeleteHomeMateButton({
  homeSlug,
  mateEmail,
}: DeleteHomeMateButtonProps) {
  const [isPending, startTransition] = useTransition();
  return (
    <button
      onClick={() => {
        startTransition(async () => {
          await removeHomeMate(homeSlug, mateEmail);
        });
      }}
      disabled={isPending}
    >
      Delete
    </button>
  );
}
