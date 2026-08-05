import { signOut } from "@/auth";

export default async function HomeErrorPage() {
  return (
    <main>
      <h1>Workspace error</h1>
      <p>Your user is not registered for this workspace.</p>
      <button
        type="button"
        onClick={async () => {
          "use server";
          console.log("ok time to signout now ");
          await signOut({ redirectTo: "/login" });
        }}
      >
        Sign Out
      </button>
    </main>
  );
}
