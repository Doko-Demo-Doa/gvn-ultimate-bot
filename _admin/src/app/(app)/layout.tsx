import { auth } from "@clerk/nextjs/server";
import type React from "react";
import UnauthenticatedScreen from "~/app/_components/unauthenticated-screen";

export default async function ProtectedLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const { userId } = await auth();

  if (!userId) {
    return <UnauthenticatedScreen />;
  }

  return children;
}
