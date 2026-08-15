import { getMateProfile } from "@/actions";
import MateProfile from "@/components/MateProfile";

export default async function ProfilePage() {
  const mate = await getMateProfile();

  return (
    <main className="m-4">
      <h1 className="text-4xl mb-12">Profile Page</h1>
      <MateProfile email={mate.email} name={mate.name} />
    </main>
  );
}
