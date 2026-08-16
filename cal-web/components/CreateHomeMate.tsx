"use client";

import { addHomeMate, CreateHomeMateState } from "@/actions";
import { useActionState } from "react";

const initialState: CreateHomeMateState = {
  slug: "",
  email: "",
  name: "",
  error: undefined,
};

interface CreateHomeMateProps {
  slug: string;
}

export default function CreateHomeMate({ slug }: CreateHomeMateProps) {
  const [state, formAction, isPending] = useActionState(
    addHomeMate,
    initialState,
  );

  return (
    <>
      <form action={formAction} className="space-y-4">
        <div className="flex gap-4">
          <input
            type="text"
            name="slug"
            required
            className="hidden"
            defaultValue={slug}
          />
          <input
            type="text"
            name="name"
            required
            maxLength={32}
            placeholder="Name"
            className="w-full p-2 border rounded"
            disabled={isPending}
            defaultValue={state.name}
          />
          <input
            type="email"
            name="email"
            required
            placeholder="name@example.com"
            className="w-full p-2 border rounded"
            disabled={isPending}
            defaultValue={state.email}
          />
        </div>
        {state.error && (
          <p className="text-sm text-red-600 bg-red-50 p-2 rounded">
            {state.error}
          </p>
        )}
        <button
          type="submit"
          disabled={isPending}
          className="w-full bg-primary-button-light dark:bg-primary-button-dark text-text-light dark:text-text-dark p-2 rounded hover:bg-primary-button-light-darker hover:dark:bg-primary-button-dark-darker disabled:bg-gray-500"
        >
          {isPending ? "Processing..." : "Add Mate"}
        </button>

        {state.success && (
          <p className="text-sm text-green-600 bg-green-50 p-2 rounded">
            New mate added
          </p>
        )}
      </form>
    </>
  );
}
