import { users } from "../constants";
import { UserCheckbox } from "../components/UserCheckbox";
import { saveOptions } from "../actions";
import { SearchParamsPromise } from "../types";
import { determineUser } from "../utils";
import { redirect } from "next/navigation";

interface OptionsPageProps {
  searchParams: SearchParamsPromise
}

export default async function Options({ searchParams }: OptionsPageProps) {
  const resolvedParams = await searchParams;

  const user = determineUser(resolvedParams);
  if (!user) {
    redirect('/user-select')
  }

  return (
    <main className="m-4">
      <h1 className="text-4xl mb-12">Options Page</h1>
      <section className="mx-2">
        <h2 className="text-2xl mb-4">Change User</h2>
        <form action={saveOptions} className="flex flex-col gap-4">
          <UserCheckbox
            userOne={users[0]}
            userTwo={users[1]}
            startChecked={users[0].id === user.id ? false : true}
          />
          <button type="submit" className="border border-white max-w-24">Save</button>
        </form>
      </section>
    </main>
  )
}
