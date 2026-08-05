"use server";

import { revalidatePath } from "next/cache";
import { redirect } from "next/navigation";
import {
  createEntry,
  deleteEntry,
  fetchListEntries,
  fetchListHomeMates,
  requestOtp,
  updateEntry,
} from "@/api";
import { rows, startIndexMap } from "@/constants";
import { convertTimeToUTCISO, convertToLocalTime } from "@/utils";
import { signIn } from "@/auth";
import { UnauthorizedError } from "@/api/types";

export type FormState = {
  step: "REQUEST_OTP" | "VERIFY_OTP";
  email: string;
  error?: string;
  success?: boolean;
};

export async function addEntry(startTime: string, endTime: string) {
  const startDateString = convertTimeToUTCISO(startTime);
  const endDateString = convertTimeToUTCISO(endTime);
  await createEntry(startDateString, endDateString);

  revalidatePath("/");
}

export async function changeEntry(entryId: number, entryValue: number) {
  await updateEntry(entryId, entryValue);

  revalidatePath("/");
}

export async function removeEntry(entryId: number) {
  await deleteEntry(entryId);

  revalidatePath("/");
}

export async function listEntries(
  workspaceUserIds: number[],
  timeZone: string,
) {
  const currentRows = structuredClone(rows);
  const entries = await Promise.all([
    fetchListEntries(workspaceUserIds[0]),
    fetchListEntries(workspaceUserIds[1]),
  ]);

  entries.forEach((entry, userIndex) => {
    entry.forEach((e: any) => {
      const startTime = convertToLocalTime(new Date(e.start), timeZone);
      const startIndex = startIndexMap.get(startTime);

      if (!!startIndex || startIndex === 0) {
        currentRows[startIndex].entryIds[userIndex] = {
          id: e.id,
          value: e.value,
        };
      } else {
        console.error("Could not find index for:", startTime);
      }
    });
  });

  return currentRows;
}

export async function saveOptions(formData: FormData) {}

export async function getOtp(
  prevState: FormState,
  formData: FormData,
): Promise<FormState> {
  const emailFormValue = formData.get("email");

  if (emailFormValue) {
    const email = emailFormValue.toString();
    if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) {
      return { ...prevState, error: "Please enter a valid email address." };
    }
    try {
      await requestOtp(email.toString());
      return {
        step: "VERIFY_OTP",
        email,
        error: undefined,
      };
    } catch (error) {
      return { ...prevState, error: "Failed to send OTP. Please try again." };
    }
  } else {
    return { ...prevState, error: "Email address is required." };
  }
}

export async function submitOtp(
  prevState: FormState,
  formData: FormData,
): Promise<FormState> {
  const otpFormValue = formData.get("otp");

  if (!otpFormValue) {
    return { ...prevState, error: "OTP is required." };
  }

  const email = prevState.email;
  const otp = otpFormValue.toString();
  if (!otp || otp.length !== 6) {
    return { ...prevState, error: "OTP must be exactly 6 digits." };
  }
  const credentials = { email, otp };
  let redirectUrl = "/";
  try {
    redirectUrl = await signIn("credentials", {
      ...credentials,
      redirect: false,
      redirectTo: "/",
    });
  } catch (error) {
    return { ...prevState, error: "Verification failed. Try again." };
  }

  redirect(redirectUrl);
}

export async function listHomeMates(homeSlug: string) {
  try {
    const matesRaw = await fetchListHomeMates(homeSlug);
    const matesIDs = matesRaw.map((mate) => parseInt(mate.id));
    const matesEmails = matesRaw.map((mate) => mate.email);

    return {
      raw: matesRaw,
      ids: matesIDs,
      emails: matesEmails,
    };
  } catch (error) {
    if (error instanceof UnauthorizedError) {
      redirect("/login");
    }

    throw error;
  }
}
