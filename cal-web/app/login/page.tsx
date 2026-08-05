import GoogleSignnButton from "@/components/GoogleSigninButton";
import { OtpLoginForm } from "@/components/OtpLoginForm";

export default function LoginPage() {
  return (
    <main className="flex flex-col">
      <div className="max-w-md mx-auto mt-10 p-6 border rounded-lg shadow-sm space-y-4">
        <OtpLoginForm />
        <GoogleSignnButton />
      </div>
    </main>
  );
}
