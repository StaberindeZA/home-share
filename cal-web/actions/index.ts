"use server";

import { revalidatePath } from "next/cache";
import { redirect } from "next/navigation";
import {
  createEntry,
  createHome,
  createHomeMate,
  deleteEntry,
  deleteHomeMate,
  fetchHome,
  fetchListEntries,
  fetchListHomeMates,
  fetchListHomes,
  fetchMateProfile,
  requestOtp,
  updateEntry,
  updateMateProfile,
  verifyMateRole,
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

export type MateProfileState = {
  name: string;
  error?: string;
  success?: boolean;
};

export type CreateHomeState = {
  name: string;
  description: string;
  error?: string;
};

export type CreateHomeMateState = {
  slug: string;
  email: string;
  name: string;
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
  const entries = await Promise.all(
    workspaceUserIds.map((ws) => fetchListEntries(ws)),
  );

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

export async function getMateProfile() {
  return fetchMateProfile();
}

export async function saveMateProfile(
  prevState: MateProfileState,
  formData: FormData,
) {
  const nameFormValue = formData.get("name");

  if (!nameFormValue) {
    return { ...prevState, error: "Name is required." };
  }

  const name = nameFormValue.toString();
  try {
    await updateMateProfile(name);
    return {
      ...prevState,
      name,
      error: undefined,
      success: true,
    };
  } catch (error) {
    return {
      ...prevState,
      error: "Failed to update Profile. Please try again.",
    };
  }
}

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
      await requestOtp(email);
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
    const matesIDs = matesRaw.map((mate) => mate.id);
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

export async function listHomes() {
  return fetchListHomes();
}

export async function readHome(slug: string) {
  return fetchHome(slug);
}

export async function addHome(prevState: CreateHomeState, formData: FormData) {
  const nameFormValue = formData.get("name");
  if (!nameFormValue) {
    return { ...prevState, error: "Name is required." };
  }
  const name = nameFormValue.toString();

  const descriptionFormValue = formData.get("description");
  const description = descriptionFormValue
    ? descriptionFormValue.toString()
    : "";

  try {
    await createHome(name, description);
  } catch (error) {
    console.error(error);
    return {
      ...prevState,
      name,
      description,
      error: "An error occurred. Please try again.",
    };
  }

  redirect("/home/select");
}

export async function addHomeMate(
  prevState: CreateHomeMateState,
  formData: FormData,
) {
  const slugFormValue = formData.get("slug");
  if (!slugFormValue) {
    return { ...prevState, error: "Slug is required." };
  }
  const slug = slugFormValue.toString();
  const emailFormValue = formData.get("email");
  if (!emailFormValue) {
    return { ...prevState, error: "Email is required." };
  }
  const email = emailFormValue.toString();
  const nameFormValue = formData.get("name");
  if (!nameFormValue) {
    return { ...prevState, error: "Name is required." };
  }
  const name = nameFormValue.toString();

  try {
    await createHomeMate(slug, name, email);
  } catch (error) {
    console.error(error);
    return {
      ...prevState,
      email,
      name,
      error: "An error occurred. Please try again.",
    };
  }

  revalidatePath("/home/[slug]/manage");

  return {
    ...prevState,
    email: "",
    name: "",
    error: undefined,
    success: true,
  };
}

export async function removeHomeMate(homeSlug: string, mateEmail: string) {
  await deleteHomeMate(homeSlug, mateEmail);

  revalidatePath("/home/[slug]/manage");
}

export async function verifyMateAdmin(homeSlug: string) {
  try {
    await verifyMateRole(homeSlug, "Admin");
    return true;
  } catch (error) {
    return false;
  }
}
