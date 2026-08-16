import { redirect } from "next/navigation";
import { getMateProfile, listEntries, listHomeMates } from "@/actions";
import UserButton from "@/components/UserButton";
import { convert24to12 } from "@/utils";
import { auth, signOut } from "@/auth";
import { revalidatePath } from "next/cache";

interface HomeSlugPageProps {
  params: Promise<{ slug: string }>;
}

export default async function HomeSlugPage({ params }: HomeSlugPageProps) {
  const { slug } = await params;
  const homeMates = await listHomeMates(slug);
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
              const redirectPath = await signOut({
                redirectTo: "/login",
                redirect: false,
              });

              revalidatePath("/", "layout");
              redirect(redirectPath);
            }}
          >
            Change user here
          </span>
        </p>
      </footer>
    </>
  );
}
