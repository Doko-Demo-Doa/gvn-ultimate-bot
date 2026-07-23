import { auth, currentUser } from "@clerk/nextjs/server";
import type React from "react";
import ForbiddenScreen from "~/app/_components/forbidden-screen";
import UnauthenticatedScreen from "~/app/_components/unauthenticated-screen";

export default async function ProtectedLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const { isAuthenticated } = await auth();

  if (!isAuthenticated) {
    return <UnauthenticatedScreen />;
  }

  const user = await currentUser();
  const discordAccount = user?.externalAccounts.find(
    (account) => account.provider === "oauth_discord",
  );

  if (!discordAccount) {
    return (
      <ForbiddenScreen reason="Link your Discord account to access the admin dashboard." />
    );
  }

  return children;
}
