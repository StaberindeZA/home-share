"use client";

import { useParams } from "next/navigation";
import { useEffect } from "react";

export default function HomeSlugTracker() {
  const params = useParams();
  const slug = params?.slug;

  useEffect(() => {
    if (slug) {
      sessionStorage.setItem("last_home_slug", String(slug));
    }
  }, [slug]);

  return null;
}
