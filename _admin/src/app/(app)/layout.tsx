import { auth, currentUser } from "@clerk/nextjs/server";
import type React from "react";
import ForbiddenScreen from "~/app/_components/forbidden-screen";
import UnauthenticatedScreen from "~/app/_components/unauthenticated-screen";

async function checkDiscordAccess(discordUserId: string) {
  const res = await fetch(
    `${process.env.NEXT_PUBLIC_BASE_API_URL}/api/admin/access-check?discord_user_id=${discordUserId}`,
    { cache: "no-store" },
  );

  if (!res.ok) {
    return { allowed: false };
  }

  const body = await res.json();
  return body.data as { allowed: boolean; access_level?: string };
}

export default async function ProtectedLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const { userId } = await auth();

  if (!userId) {
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

  const access = await checkDiscordAccess(discordAccount.providerUserId);

  if (!access.allowed) {
    return (
      <ForbiddenScreen reason="Your Discord roles don't grant access to the admin dashboard." />
    );
  }

  return children;
}
