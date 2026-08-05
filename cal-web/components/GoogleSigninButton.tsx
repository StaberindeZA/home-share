"use client";

import { signIn } from "next-auth/react";

export default function GoogleSignnButton() {
  return (
    <button
      onClick={() => signIn("google", { redirectTo: "/" })}
      className="w-full bg-primary-button-light dark:bg-primary-button-dark text-primary-text-light dark:text-primary-text-dark p-2 rounded hover:bg-primary-button-light-darker hover:dark:bg-primary-button-dark-darker disabled:bg-gray-400"
    >
      Google Sign In
    </button>
  );
}
