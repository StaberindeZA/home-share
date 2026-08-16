"use client";

import { addHome, CreateHomeState } from "@/actions";
import { useActionState } from "react";

const initialState: CreateHomeState = {
  name: "",
  description: "",
  error: undefined,
};

export default function CreateHome() {
  const [state, formAction, isPending] = useActionState(addHome, initialState);

  return (
    <>
      <form action={formAction} className="space-y-4">
        <div>
          <label className="block text-sm font-medium mb-1">Name</label>
          <input
            type="text"
            name="name"
            required
            maxLength={32}
            placeholder="Name"
            className="w-full p-2 border rounded"
            disabled={isPending}
          />
          <label className="block text-sm font-medium mb-1">Description</label>
          <textarea
            name="description"
            rows={2}
            placeholder="Description"
            className="w-full p-2 border rounded"
            disabled={isPending}
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
          {isPending ? "Processing..." : "Create"}
        </button>
      </form>
    </>
  );
}
