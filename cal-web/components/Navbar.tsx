"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { useParams, usePathname } from "next/navigation";
import { verifyMateRole } from "@/api";
import { signOut } from "next-auth/react";

const signedInLinks = [{ name: "Profile", href: "/profile" }];
const signedOutLinks = [
  { name: "About", href: "/about" },
  {
    name: "Sign In",
    href: "/login",
  },
];

interface NavbarProps {
  isSignedIn: boolean;
}

type NavbarParams = {
  slug?: string;
};

// Example Navbar component structure with links and state for mobile menu
export default function Navbar({ isSignedIn }: NavbarProps) {
  const [isOpen, setIsOpen] = useState(false);
  const [sessionSlug, setSessionSlug] = useState("");
  const [isAdmin, setIsAdmin] = useState(false);
  const pathname = usePathname();
  const params = useParams<NavbarParams>();

  const slug = params?.slug || sessionSlug;

  useEffect(() => {
    if (params?.slug) {
      return;
    }
    const savedSlug = sessionStorage.getItem("last_home_slug");
    if (savedSlug) {
      setSessionSlug(savedSlug);
    }
  }, [pathname]);

  useEffect(() => {
    if (!params?.slug) {
      return;
    }

    const checkAdmin = async (slug: string) => {
      try {
        await verifyMateRole(slug, "Admin");
        setIsAdmin(true);
      } catch (error) {
        setIsAdmin(false);
      }
    };

    checkAdmin(params.slug);
  }, [pathname]);

  const navLinks = [];

  if (isSignedIn) {
    if (slug) {
      navLinks.push({
        name: "Manage Home",
        href: `/home/${slug}/manage`,
      });
      navLinks.push({ name: "Calendar", href: `/home/${slug}` });
    } else {
      navLinks.push({ name: "Select Home", href: `/home/select` });
    }

    signedInLinks.forEach((link) => navLinks.push(link));
  } else {
    signedOutLinks.forEach((link) => navLinks.push(link));
  }

  return (
    <nav className="sticky top-0 z-50 w-full border-b border-gray-200 bg-white/80 backdrop-blur-md dark:border-gray-800 dark:bg-gray-950/80">
      <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        <div className="flex h-16 items-center justify-between">
          {/* Logo */}
          <Link
            href={slug ? `/home/${slug}/manage` : "/home/select"}
            className="text-xl font-bold text-text-light dark:text-text-dark"
          >
            HomeShare
          </Link>

          {/* Desktop Navigation Links */}
          <div className="hidden md:flex space-x-4">
            {navLinks.map((link) => (
              <Link
                key={link.name}
                href={link.href}
                className={`rounded-md px-3 py-2 text-sm font-medium transition-colors ${
                  pathname === link.href
                    ? "bg-blue-50 text-blue-600 dark:bg-blue-950/50 dark:text-blue-400"
                    : "text-gray-600 hover:bg-gray-50 hover:text-gray-900 dark:text-gray-300 dark:hover:bg-gray-900 dark:hover:text-white"
                }`}
              >
                {link.name}
              </Link>
            ))}
            {isSignedIn && (
              <span
                className="rounded-md px-3 py-2 text-sm font-medium transition-colors text-gray-600 hover:bg-gray-50 hover:text-gray-900 dark:text-gray-300 dark:hover:bg-gray-900 dark:hover:text-white cursor-pointer"
                onClick={() => {
                  sessionStorage.clear();
                  signOut({
                    redirectTo: "/login",
                  });
                }}
              >
                Sign Out
              </span>
            )}
          </div>

          {/* Mobile Menu Button */}
          <div className="flex md:hidden">
            <button
              onClick={() => setIsOpen(!isOpen)}
              type="button"
              className="inline-flex items-center justify-center rounded-md p-2 text-gray-500 hover:bg-gray-100 hover:text-gray-700 focus:outline-none dark:text-gray-400 dark:hover:bg-gray-900 dark:hover:text-white"
              aria-controls="mobile-menu"
              aria-expanded={isOpen}
            >
              <span className="sr-only">Open main menu</span>
              {/* Dynamic Hamburger / X Icon using SVG */}
              <svg
                className="h-6 w-6"
                fill="none"
                viewBox="0 0 24 24"
                strokeWidth="1.5"
                stroke="currentColor"
              >
                {isOpen ? (
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    d="M6 18L18 6M6 6l12 12"
                  />
                ) : (
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    d="M3.75 6.75h16.5M3.75 12h16.5m-16.5 5.25h16.5"
                  />
                )}
              </svg>
            </button>
          </div>
        </div>
      </div>

      {/* Mobile Menu Panel */}
      <div
        className={`${isOpen ? "block" : "hidden"} md:hidden border-t border-gray-200 dark:border-gray-800 bg-white dark:bg-gray-950`}
        id="mobile-menu"
      >
        <div className="space-y-1 px-2 pb-3 pt-2">
          {navLinks.map((link) => (
            <Link
              key={link.name}
              href={link.href}
              onClick={() => setIsOpen(false)} // Closes menu when a link is clicked
              className={`block rounded-md px-3 py-2 text-base font-medium transition-colors ${
                pathname === link.href
                  ? "bg-blue-50 text-blue-600 dark:bg-blue-950/50 dark:text-blue-400"
                  : "text-gray-600 hover:bg-gray-50 hover:text-gray-900 dark:text-gray-300 dark:hover:bg-gray-900 dark:hover:text-white"
              }`}
            >
              {link.name}
            </Link>
          ))}
          {isSignedIn && (
            <span
              className="block rounded-md px-3 py-2 text-base font-medium transition-colors text-gray-600 hover:bg-gray-50 hover:text-gray-900 dark:text-gray-300 dark:hover:bg-gray-900 dark:hover:text-white cursor-pointer"
              onClick={() => {
                sessionStorage.clear();
                signOut({
                  redirectTo: "/login",
                });
              }}
            >
              Sign Out
            </span>
          )}
        </div>
      </div>
    </nav>
  );
}
