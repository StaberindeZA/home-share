import { listHomeMates, readHome, verifyMateAdmin } from "@/actions";
import CreateHomeMate from "@/components/CreateHomeMate";
import DeleteHomeMateButton from "@/components/DeleteHomeMateButton";
import Link from "next/link";
import { redirect } from "next/navigation";

interface HomeSlugManagePageProps {
  params: Promise<{ slug: string }>;
}

export default async function HomeSlugManagePage({
  params,
}: HomeSlugManagePageProps) {
  const { slug } = await params;
  const home = await readHome(slug);
  const homeMates = await listHomeMates(slug);
  const isAdmin = await verifyMateAdmin(slug);

  return (
    <section>
      <h1>Manage your Home Space</h1>
      <h2>Details</h2>
      <div className="mb-10">
        <p>Name: {home.name}</p>
        <p>Description: {home.description}</p>
        <h3 className="mt-6">Not the Home Space you are looking for?</h3>
        <Link
          href="/home/select"
          className="w-full bg-primary-button-light dark:bg-primary-button-dark text-text-light dark:text-text-dark p-2 rounded hover:bg-primary-button-light-darker hover:dark:bg-primary-button-dark-darker disabled:bg-gray-500"
        >
          Change Home Space
        </Link>
      </div>
      {isAdmin && (
        <div className="mb-10">
          <h2>Admin Zone</h2>
          <div>
            <h3>Mates</h3>
            <table className="m-4 w-full table-fixed max-w-lg">
              <thead>
                <tr>
                  <th>Name</th>
                  <th>Email</th>
                  <th>Role</th>
                  <th>Actions</th>
                </tr>
              </thead>
              <tbody>
                {homeMates.raw.map((mate) => (
                  <tr
                    key={mate.email}
                    className="h-8 border border-black dark:border-white hover:bg-gray-100 hover:dark:bg-gray-800"
                  >
                    <th>{mate.name}</th>
                    <th>{mate.email}</th>
                    <th>{mate.role}</th>
                    <th>
                      {isAdmin && mate.role !== "Admin" && (
                        <DeleteHomeMateButton
                          homeSlug={slug}
                          mateEmail={mate.email}
                        />
                      )}
                    </th>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
          <div>
            <h3>Add more mates</h3>
            <CreateHomeMate slug={slug} />
          </div>
        </div>
      )}
    </section>
  );
}
