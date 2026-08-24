"use server";

import { auth, SessionWithToken } from "@/auth";
import { API_URL } from "@/constants";
import { convertTimeToUTCISO } from "@/utils";
import {
  ForbiddenError,
  Home,
  HomeMate,
  MateProfile,
  UnauthorizedError,
} from "./types";

async function getAccessToken() {
  const session = (await auth()) as SessionWithToken;
  return session.accessToken;
}

async function makeRequest(url: URL, init?: RequestInit, isUnauthed?: boolean) {
  let headers = {};
  if (!isUnauthed) {
    const accessToken = await getAccessToken();
    headers = {
      Authorization: `Bearer ${accessToken}`,
    };
  }

  try {
    const response = await fetch(url, {
      ...init,
      headers: {
        ...headers,
        ...init?.headers,
      },
    });

    if (!response.ok) {
      if (response.status === 401) {
        throw new UnauthorizedError(`url: ${url.pathname}`);
      } else if (response.status === 403) {
        throw new ForbiddenError(`url: ${url.pathname}`);
      } else {
        throw new Error(`HTTP error! Status: ${response.status}`);
      }
    }

    if (response.headers.get("Content-Type") === "application/json") {
      return response.json();
    }

    return;
  } catch (error) {
    if (!(error instanceof ForbiddenError || UnauthorizedError)) {
      console.error("Error:", error);
    }
    throw error;
  }
}

export async function createEntry(start: string, end: string) {
  const data = {
    start,
    end,
  };
  const url = new URL(`${API_URL}/v1/entry`);
  await makeRequest(url, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(data),
  });
}

export async function updateEntry(entryId: number, entryValue: number) {
  const data = {
    value: entryValue,
  };
  const url = new URL(`${API_URL}/v1/entry/${entryId}`);
  await makeRequest(url, {
    method: "PUT",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(data),
  });
}

export async function deleteEntry(entryId: number) {
  const url = new URL(`${API_URL}/v1/entry/${entryId}`);
  await makeRequest(url, {
    method: "DELETE",
  });
}

export async function fetchListEntries(userId: number) {
  const startDateString = convertTimeToUTCISO("00:00:00");
  const endDateString = convertTimeToUTCISO("23:59:59");
  const params = {
    userId: `${userId}`,
    start: startDateString,
    end: endDateString,
  };
  const url = new URL(`${API_URL}/v1/entry`);
  url.search = new URLSearchParams(params).toString();

  return makeRequest(url);
}

export async function requestOtp(email: string) {
  const data = {
    email,
  };
  const url = new URL(`${API_URL}/v1/otp`);
  await makeRequest(
    url,
    {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(data),
    },
    true,
  );
}

export async function verifyOtp(email: string, otp: string) {
  const data = {
    email,
    otp,
  };
  const url = new URL(`${API_URL}/v1/otp/verify`);
  return makeRequest(
    url,
    {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(data),
    },
    true,
  );
}

export async function fetchMateProfile() {
  const url = new URL(`${API_URL}/v1/user`);

  const responseData = await makeRequest(url);
  return responseData as MateProfile;
}

export async function updateMateProfile(name: string) {
  const data = {
    name,
  };
  const url = new URL(`${API_URL}/v1/user`);
  await makeRequest(url, {
    method: "PUT",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(data),
  });
}

export async function fetchListHomes() {
  const url = new URL(`${API_URL}/v1/homes`);
  const responseData = await makeRequest(url);
  if (responseData) {
    return responseData as Home[];
  } else {
    return [] as Home[];
  }
}

export async function fetchHome(slug: string) {
  const url = new URL(`${API_URL}/v1/home/${slug}`);
  const responseData = await makeRequest(url);
  return responseData as Home;
}

export async function createHome(name: string, description: string) {
  const url = new URL(`${API_URL}/v1/home`);
  const data = {
    name,
    description,
  };
  await makeRequest(url, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(data),
  });
}

export async function createHomeMate(
  homeSlug: string,
  name: string,
  email: string,
) {
  const url = new URL(`${API_URL}/v1/home/${homeSlug}/mate`);
  const data = {
    email,
    name,
  };
  await makeRequest(url, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(data),
  });
}

export async function fetchListHomeMates(homeSlug: string) {
  const url = new URL(`${API_URL}/v1/home/${homeSlug}/mates`);
  const responseData = await makeRequest(url);
  return responseData as HomeMate[];
}

export async function deleteHomeMate(slug: string, email: string) {
  const params = {
    email,
  };
  const url = new URL(`${API_URL}/v1/home/${slug}/mate`);
  url.search = new URLSearchParams(params).toString();
  await makeRequest(url, {
    method: "DELETE",
  });
}

export async function verifyMateRole(homeSlug: string, role: "Admin" | "Mate") {
  const data = {
    homeSlug,
    role,
  };
  const url = new URL(`${API_URL}/v1/auth/role/verify`);
  await makeRequest(url, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(data),
  });
}
