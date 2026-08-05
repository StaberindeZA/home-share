import { saveOptions } from "@/actions";

export default async function Options() {
  return (
    <main className="m-4">
      <h1 className="text-4xl mb-12">Options Page</h1>
      <section className="mx-2">
        <h2 className="text-2xl mb-4">Change User</h2>
        <form action={saveOptions} className="flex flex-col gap-4">
          <button type="submit" className="border border-white max-w-24">
            Save
          </button>
        </form>
      </section>
    </main>
  );
}
