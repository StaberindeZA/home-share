import { NextResponse } from "next/server";

export const dynamic = "force-dynamic";

export async function GET() {
  try {
    return NextResponse.json(
      { status: "UP", timestamp: Date.now() },
      {
        status: 200,
        headers: {
          "Cache-Control":
            "no-store, no-cache, must-revalidate, proxy-revalidate",
        },
      },
    );
  } catch (error) {
    const message = error instanceof Error ? error.message : "Unknonw error";
    return NextResponse.json(
      { status: "DOWN", error: message },
      {
        status: 503,
        headers: { "Cache-Control": "no-store" },
      },
    );
  }
}
