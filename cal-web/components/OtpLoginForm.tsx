"use client";
import { FormState, getOtp, submitOtp } from "@/actions";
import { useActionState } from "react";

const initialState: FormState = {
  step: "REQUEST_OTP",
  email: "",
  error: undefined,
};

export function OtpLoginForm() {
  async function formCoordinator(prevState: FormState, formData: FormData) {
    if (prevState.step === "REQUEST_OTP") {
      return getOtp(prevState, formData);
    } else {
      return submitOtp(prevState, formData);
    }
  }
  const [state, formAction, isPending] = useActionState(
    formCoordinator,
    initialState,
  );

  return (
    <>
      <h2 className="text-xl font-bold mb-4">
        {state.step === "REQUEST_OTP"
          ? "Sign In / Register"
          : "Verify Your Identity"}
      </h2>

      <form action={formAction} className="space-y-4">
        {state.step === "REQUEST_OTP" ? (
          /* STEP 1: Enter Email */
          <div>
            <label className="block text-sm font-medium mb-1">
              Email Address
            </label>
            <input
              type="email"
              name="email"
              required
              placeholder="name@example.com"
              className="w-full p-2 border rounded"
              disabled={isPending}
            />
          </div>
        ) : (
          /* STEP 2: Enter OTP */
          <div>
            <p className="text-sm text-gray-600 mb-2">
              We sent a 6-digit code to <strong>{state.email}</strong>
            </p>
            <label className="block text-sm font-medium mb-1">
              One-Time Password (OTP)
            </label>
            <input
              type="text"
              name="otp"
              maxLength={6}
              required
              placeholder="000000"
              className="w-full p-2 border tracking-widest text-center font-mono text-lg rounded"
              disabled={isPending}
            />
          </div>
        )}

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
          {isPending
            ? "Processing..."
            : state.step === "REQUEST_OTP"
              ? "Send OTP Code"
              : "Verify & Log In"}
        </button>
      </form>
    </>
  );
}
