import { auth } from "@/auth";

export default auth((req) => {
  if (!req.auth) {
    const url = req.url;
    return Response.redirect(new URL("/login", url));
  }
});

export const config = {
  matcher: ["/", "/profile"],
};
