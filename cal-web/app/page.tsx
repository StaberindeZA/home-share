import { redirect } from "next/navigation";
import { listEntries } from "./actions";
import UserButton from "./components/UserButton";
import { users } from "./constants";
import { SearchParamsPromise } from "./types";
import { convert24to12, determineUser } from "./utils";

interface HomePageProps {
  searchParams: SearchParamsPromise
}

export default async function Home({ searchParams }: HomePageProps) {
  const resolvedParams = await searchParams;
  const queryString = new URLSearchParams(resolvedParams as Record<string, string>).toString();

  const user = determineUser(resolvedParams);
  if (!user) {
    redirect('/user-select')
  }

  const rows = await listEntries(Intl.DateTimeFormat().resolvedOptions().timeZone)

  return (
    <>
      <table className="m-4 w-full table-fixed max-w-2xl">
        <thead>
          <tr key="header" className="text-lg">
            <th>Time</th>
            <th>{users[0].name}</th>
            <th>{users[1].name}</th>
          </tr>
        </thead>
        <tbody>
          {rows.map((row) => (
            <tr key={row.start} className="h-8 border border-black dark:border-white hover:bg-gray-100 hover:dark:bg-gray-800">
              <th scope="row">{convert24to12(row.start)}</th>
              <td><UserButton startTime={row.start} endTime={row.end} user={users[0]} rowEntry={row.entryIds.at(0)} /></td>
              <td><UserButton startTime={row.start} endTime={row.end} user={users[1]} rowEntry={row.entryIds.at(1)} /></td>
            </tr>
          ))}
        </tbody>
      </table>
      <footer className="flex flex-col gap-2 mx-2">
        <p className="text-xl">Hi, I'm {user?.name || 'USER'}</p>
        <p className="text-3xl font-bold text-center">Do you have a meeting?</p>
        <p>Not you? <a href={`/options?${queryString}`} className="font-bold underline">Change user here</a></p>
      </footer>
    </>
  );
}
