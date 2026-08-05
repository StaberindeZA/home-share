import NextAuth, { Session } from "next-auth";
import GoogleProvider from "next-auth/providers/google";
import CredentialsProvider from "next-auth/providers/credentials";
import { JWT } from "next-auth/jwt";

interface CredentialUser {
  token: string;
}

interface JWTWithToken extends JWT {
  accessToken?: string;
}

export interface SessionWithToken extends Session {
  accessToken: string;
}

const CAL_API_URL = process.env.API_URL || "http://localhost:8080";

export const { handlers, signIn, signOut, auth } = NextAuth({
  providers: [
    GoogleProvider({
      clientId: process.env.GOOGLE_CLIENT_ID,
      clientSecret: process.env.GOOGLE_CLIENT_SECRET,
    }),
    CredentialsProvider({
      name: "OTP Login",
      credentials: {
        email: { label: "Email", type: "email" },
        otp: { label: "One-Time Password", type: "text" },
      },
      async authorize(credentials) {
        const email = credentials.email;
        const otp = credentials.otp;

        if (!email || !otp) {
          throw new Error("Email and OTP are required");
        }

        const data = {
          email,
          otp,
        };
        // Send OTP and Email to Go backend to verify
        const res = await fetch(`${CAL_API_URL}/v1/otp/verify`, {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify(data),
        });
        const token = await res.json();

        if (res.ok && token.token) {
          const userinfoResponse = await fetch(
            `${CAL_API_URL}/v1/auth/userinfo`,
            {
              headers: {
                Authorization: `Bearer ${token.token}`,
              },
            },
          );

          if (!userinfoResponse.ok) {
            console.error(
              "User Info fetch failed with code:",
              userinfoResponse.status,
            );
            return null;
          }

          const userinfo = await userinfoResponse.json();

          return {
            ...userinfo,
            ...token,
          };
        }
        return null;
      },
    }),
  ],
  callbacks: {
    async jwt({ token, user, account }) {
      // 1. If Google Login: exchange Google token for your BACKEND JWT
      if (account && account.provider === "google") {
        const res = await fetch(`${CAL_API_URL}/v1/auth/google`, {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ idToken: account.id_token }),
        });
        const backendUser = await res.json();
        token.accessToken = backendUser.token; // Your custom backend JWT
      }
      // 2. If OTP Login: User data already contains the backend token from authorize()
      else if (user) {
        token.accessToken = (user as CredentialUser).token;
      }
      return token as JWTWithToken;
    },
    async session({ session, token }) {
      const tokenWithAccess: JWTWithToken = token;
      if (!tokenWithAccess.accessToken) {
        throw new Error("Could not authenticate user");
      }

      const sessionWithToken: SessionWithToken = {
        ...session,
        accessToken: tokenWithAccess.accessToken,
      };
      return sessionWithToken;
    },
  },
});
