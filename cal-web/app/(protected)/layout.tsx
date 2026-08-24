import Navbar from "@/components/Navbar";

export default function OpenLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <>
      <Navbar isSignedIn={true} />
      <main className="min-h-full flex flex-col items-center">{children}</main>
    </>
  );
}
