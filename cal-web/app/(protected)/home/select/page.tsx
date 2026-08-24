import { listHomes } from "@/actions";
import HomeSelector from "@/components/HomeSelector";
import Link from "next/link";

export default async function HomeSelectPage() {
  const homes = await listHomes();

  return (
    <section className="my-4">
      <h1>Choose your Home Space</h1>
      <HomeSelector homes={homes} />
      <Link
        href={`/home/create`}
        className="w-full text-center bg-primary-button-light dark:bg-primary-button-dark text-text-light dark:text-text-dark p-2 rounded hover:bg-primary-button-light-darker hover:dark:bg-primary-button-dark-darker disabled:bg-gray-500"
      >
        + New Home Space
      </Link>
    </section>
  );
}
