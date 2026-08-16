import { Home } from "@/api/types";
import Link from "next/link";

interface HomeSelectorProps {
  homes: Home[];
}

export default function HomeSelector({ homes }: HomeSelectorProps) {
  if (!homes.length) {
    return null;
  }

  return (
    <>
      {homes.map((home) => (
        <div key={home.slug} className="flex flex-col gap-2 border p-4 my-4">
          <h2>{home.name}</h2>
          <p>{home.description}</p>
          <Link
            href={`/home/${home.slug}`}
            className="w-full text-center bg-primary-button-light dark:bg-primary-button-dark text-text-light dark:text-text-dark p-2 rounded hover:bg-primary-button-light-darker hover:dark:bg-primary-button-dark-darker disabled:bg-gray-500"
          >
            Select
          </Link>
        </div>
      ))}
    </>
  );
}
