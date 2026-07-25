import Link from "next/link";
import { users } from "../constants";

export default function UserSelectPage() {
  return (
    <main className="m-4">
      <h1 className="text-4xl text-center mb-12">Wait...who are you?!</h1>
      <div className="flex flex-col items-center justify-center gap-6">
        <Link href={`/?user=${users[0].id}`} className="flex items-center justify-center h-32 w-32 bg-green-600">{users[0].name}</Link>
        <Link href={`/?user=${users[1].id}`} className="flex items-center justify-center h-32 w-32 bg-blue-600">{users[1].name}</Link>
      </div>
    </main>
  )
}

