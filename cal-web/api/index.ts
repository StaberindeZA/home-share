"use server";

import { auth, SessionWithToken } from "@/auth";
import { API_URL } from "@/constants";
import { convertTimeToUTCISO } from "@/utils";
import { HomeMate, UnauthorizedError } from "./types";

async function getAccessToken() {
  const session = (await auth()) as SessionWithToken;
  return session.accessToken;
}

export async function createEntry(start: string, end: string) {
  const accessToken = await getAccessToken();
  const data = {
    start,
    end,
  };
  const url = new URL(`${API_URL}/v1/entry`);
  try {
    const response = await fetch(url, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${accessToken}`,
      },
      body: JSON.stringify(data),
    });

    if (!response.ok) {
      throw new Error(`HTTP error! Status: ${response.status}`);
    }
  } catch (error) {
    console.error("Error:", error);
    throw error;
  }
}

export async function updateEntry(entryId: number, entryValue: number) {
  const accessToken = await getAccessToken();
  const data = {
    value: entryValue,
  };
  const url = new URL(`${API_URL}/v1/entry/${entryId}`);
  try {
    const response = await fetch(url, {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${accessToken}`,
      },
      body: JSON.stringify(data),
    });

    if (!response.ok) {
      throw new Error(`HTTP error! Status: ${response.status}`);
    }

    return;
  } catch (error) {
    console.error("Error:", error);
    throw error;
  }
}

export async function deleteEntry(entryId: number) {
  const accessToken = await getAccessToken();
  const url = new URL(`${API_URL}/v1/entry/${entryId}`);
  try {
    const response = await fetch(url, {
      method: "DELETE",
      headers: {
        Authorization: `Bearer ${accessToken}`,
      },
    });

    if (!response.ok) {
      throw new Error(`HTTP error! Status: ${response.status}`);
    }

    return await response.json();
  } catch (error) {
    console.error("Error:", error);
    throw error;
  }
}

export async function fetchListEntries(userId: number) {
  const accessToken = await getAccessToken();
  const startDateString = convertTimeToUTCISO("00:00:00");
  const endDateString = convertTimeToUTCISO("23:59:59");
  const params = {
    userId: `${userId}`,
    start: startDateString,
    end: endDateString,
  };
  const url = new URL(`${API_URL}/v1/entry`);
  url.search = new URLSearchParams(params).toString();
  try {
    const response = await fetch(url, {
      headers: {
        Authorization: `Bearer ${accessToken}`,
      },
    });
    if (!response.ok) {
      throw new Error(`HTTP error! Status: ${response.status}`);
    }
    return response.json();
  } catch (err) {
    console.error("Something bad happend", err);
    throw err;
  }
}

export async function requestOtp(email: string) {
  const data = {
    email,
  };
  const url = new URL(`${API_URL}/v1/otp`);
  try {
    const response = await fetch(url, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(data),
    });

    if (!response.ok) {
      throw new Error(`HTTP error! Status: ${response.status}`);
    }

    return;
  } catch (error) {
    console.error("Error:", error);
    throw error;
  }
}

export async function verifyOtp(email: string, otp: string) {
  const data = {
    email,
    otp,
  };
  const url = new URL(`${API_URL}/v1/otp/verify`);
  try {
    const response = await fetch(url, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(data),
    });

    if (!response.ok) {
      throw new Error(`HTTP error! Status: ${response.status}`);
    }

    return await response.json();
  } catch (error) {
    console.error("Error:", error);
    throw error;
  }
}

export async function fetchListHomeMates(homeSlug: string) {
  const accessToken = await getAccessToken();
  const url = new URL(`${API_URL}/v1/home/${homeSlug}/mates`);
  try {
    const response = await fetch(url, {
      headers: {
        Authorization: `Bearer ${accessToken}`,
      },
    });
    if (!response.ok) {
      if (response.status === 401) {
        throw new UnauthorizedError("fetchListHomeMates");
      } else {
        throw new Error(`HTTP error! Status: ${response.status}`);
      }
    }
    const responseData = await response.json();
    return responseData as HomeMate[];
  } catch (err) {
    console.error("Error in fetchListHomeMates", err);
    throw err;
  }
}
