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

export async function fetchMateProfile() {
  const accessToken = await getAccessToken();
  const url = new URL(`${API_URL}/v1/user`);
  try {
    const response = await fetch(url, {
      headers: {
        Authorization: `Bearer ${accessToken}`,
      },
    });
    if (!response.ok) {
      if (response.status === 401) {
        throw new UnauthorizedError("fetchMateProfile");
      } else {
        throw new Error(`HTTP error! Status: ${response.status}`);
      }
    }
    const responseData = await response.json();
    return responseData as MateProfile;
  } catch (err) {
    console.error("Error in fetchListHomeMates", err);
    throw err;
  }
}

export async function updateMateProfile(name: string) {
  const accessToken = await getAccessToken();
  const data = {
    name,
  };
  const url = new URL(`${API_URL}/v1/user`);
  try {
    const response = await fetch(url, {
      method: "PUT",
      headers: {
        Authorization: `Bearer ${accessToken}`,
        "Content-Type": "application/json",
      },
      body: JSON.stringify(data),
    });
    if (!response.ok) {
      if (response.status === 401) {
        throw new UnauthorizedError("updateMateProfile");
      } else {
        throw new Error(`HTTP error! Status: ${response.status}`);
      }
    }
    return;
  } catch (err) {
    console.error("Error in fetchListHomeMates", err);
    throw err;
  }
}

export async function fetchListHomes() {
  const accessToken = await getAccessToken();
  const url = new URL(`${API_URL}/v1/homes`);
  try {
    const response = await fetch(url, {
      headers: {
        Authorization: `Bearer ${accessToken}`,
      },
    });
    if (!response.ok) {
      if (response.status === 401) {
        throw new UnauthorizedError("fetchListHomes");
      } else {
        throw new Error(`HTTP error! Status: ${response.status}`);
      }
    }
    const responseData = await response.json();
    if (responseData) {
      return responseData as Home[];
    } else {
      return [] as Home[];
    }
  } catch (err) {
    console.error("Error in fetchListHomes", err);
    throw err;
  }
}

export async function fetchHome(slug: string) {
  const accessToken = await getAccessToken();
  const url = new URL(`${API_URL}/v1/home/${slug}`);
  try {
    const response = await fetch(url, {
      headers: {
        Authorization: `Bearer ${accessToken}`,
      },
    });
    if (!response.ok) {
      if (response.status === 401) {
        throw new UnauthorizedError("readHome");
      } else {
        throw new Error(`HTTP error! Status: ${response.status}`);
      }
    }
    const responseData = await response.json();
    return responseData as Home;
  } catch (err) {
    console.error("Error in fetchListHomes", err);
    throw err;
  }
}

export async function createHome(name: string, description: string) {
  const accessToken = await getAccessToken();
  const url = new URL(`${API_URL}/v1/home`);
  const data = {
    name,
    description,
  };
  try {
    const response = await fetch(url, {
      method: "POST",
      headers: {
        Authorization: `Bearer ${accessToken}`,
        "Content-Type": "application/json",
      },
      body: JSON.stringify(data),
    });
    if (!response.ok) {
      if (response.status === 401) {
        throw new UnauthorizedError("createHome");
      } else {
        throw new Error(`HTTP error! Status: ${response.status}`);
      }
    }
    return;
  } catch (err) {
    console.error("Error in createHome", err);
    throw err;
  }
}

export async function createHomeMate(
  homeSlug: string,
  name: string,
  email: string,
) {
  const accessToken = await getAccessToken();
  const url = new URL(`${API_URL}/v1/home/${homeSlug}/mate`);
  const data = {
    email,
    name,
  };
  try {
    const response = await fetch(url, {
      method: "POST",
      headers: {
        Authorization: `Bearer ${accessToken}`,
        "Content-Type": "application/json",
      },
      body: JSON.stringify(data),
    });
    if (!response.ok) {
      if (response.status === 401) {
        throw new UnauthorizedError("fetchListHomeMates");
      } else {
        throw new Error(`HTTP error! Status: ${response.status}`);
      }
    }
    return;
  } catch (err) {
    console.error("Error in fetchListHomeMates", err);
    throw err;
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

export async function deleteHomeMate(slug: string, email: string) {
  const accessToken = await getAccessToken();
  const params = {
    email,
  };
  const url = new URL(`${API_URL}/v1/home/${slug}/mate`);
  url.search = new URLSearchParams(params).toString();
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

    return;
  } catch (error) {
    console.error("Error:", error);
    throw error;
  }
}

export async function verifyMateRole(homeSlug: string, role: "Admin" | "Mate") {
  const accessToken = await getAccessToken();
  const data = {
    homeSlug,
    role,
  };
  const url = new URL(`${API_URL}/v1/auth/role/verify`);
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
      if (response.status === 403) {
        throw new ForbiddenError("fetchMateProfile");
      } else {
        throw new Error(`HTTP error! Status: ${response.status}`);
      }
    }
  } catch (error) {
    if (!(error instanceof ForbiddenError)) {
      console.error("Error:", error);
    }
    throw error;
  }
}
