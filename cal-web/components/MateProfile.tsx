"use client";

import { MateProfileState, saveMateProfile } from "@/actions";
import { useActionState } from "react";

const initialState: MateProfileState = {
  name: "",
  error: undefined,
  success: undefined,
};

interface MateProfileProps {
  email: string;
  name: string;
}

export default function MateProfile({ email, name }: MateProfileProps) {
  initialState.name = name;

  const [state, formAction, isPending] = useActionState(
    saveMateProfile,
    initialState,
  );

  return (
    <section className="mx-2">
      <h2 className="text-2xl mb-4">User</h2>
      <form action={formAction} className="flex flex-col gap-4">
        <div>
          <label className="block text-sm font-medium mb-1">
            Email Address
          </label>
          <input
            type="text"
            name="email"
            required
            defaultValue={email}
            className="w-full p-2 border rounded"
            disabled={true}
          />
        </div>
        <div>
          <label className="block text-sm font-medium mb-1">Name</label>
          <input
            type="text"
            name="name"
            required
            defaultValue={state.name}
            className="w-full p-2 border rounded"
            disabled={isPending}
          />
        </div>
        {state.error && (
          <p className="text-sm text-red-600 bg-red-50 p-2 rounded">
            {state.error}
          </p>
        )}

        <button type="submit" className="border border-white max-w-24">
          Save
        </button>

        {state.success && (
          <p className="text-sm text-green-600 bg-green-50 p-2 rounded">
            Profile saved
          </p>
        )}
      </form>
    </section>
  );
}
