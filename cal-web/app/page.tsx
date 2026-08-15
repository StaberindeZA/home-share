import { redirect } from "next/navigation";
import { getMateProfile, listEntries, listHomeMates } from "@/actions";
import UserButton from "@/components/UserButton";
import { convert24to12 } from "@/utils";
import { auth, signOut } from "@/auth";

const STATIC_HOME = "temp_hardcoded_home";

export default async function Home() {
  const homeMates = await listHomeMates(STATIC_HOME);
  const mate = await getMateProfile();

  const session = await auth();

  if (session?.user?.email && !homeMates.emails.includes(session.user.email)) {
    redirect("/home/error");
  }

  const rows = await listEntries(
    homeMates.ids,
    Intl.DateTimeFormat().resolvedOptions().timeZone,
  );

  return (
    <>
      <table className="m-4 w-full table-fixed max-w-2xl">
        <thead>
          <tr key="header" className="text-lg">
            <th>Time</th>
            {homeMates.raw.map((homeMate) => (
              <th key={`th-${homeMate.email}`}>
                {homeMate.email === mate.email ? mate.name : homeMate.name}
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {rows.map((row) => (
            <tr
              key={row.start}
              className="h-8 border border-black dark:border-white hover:bg-gray-100 hover:dark:bg-gray-800"
            >
              <th scope="row">{convert24to12(row.start)}</th>
              {homeMates.raw.map((homeMate, index) => (
                <td key={row.start + homeMate.id}>
                  <UserButton
                    startTime={row.start}
                    endTime={row.end}
                    homeMate={homeMate}
                    sessionUser={session?.user}
                    rowEntry={row.entryIds.at(index)}
                  />
                </td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>
      <footer className="flex flex-col gap-2 mx-2">
        <p className="text-xl">Hi, I'm {mate.name}</p>
        <p className="text-3xl font-bold text-center">Do you have a meeting?</p>
        <p>
          Not you?{" "}
          <span
            className="font-bold underline cursor-pointer"
            onClick={async () => {
              "use server";
              await signOut({ redirectTo: "/login" });
            }}
          >
            Change user here
          </span>
        </p>
      </footer>
    </>
  );
}
