"use client";

import { useSignIn } from "@clerk/nextjs";
import { Button, type ButtonProps } from "@mantine/core";

type DiscordAuthButtonProps = ButtonProps & {
  children?: React.ReactNode;
};

export default function DiscordAuthButton({
  children = "Continue with Discord",
  ...props
}: DiscordAuthButtonProps) {
  const { fetchStatus, signIn } = useSignIn();

  const handleClick = async () => {
    await signIn.sso({
      strategy: "oauth_discord",
      redirectUrl: "/",
      redirectCallbackUrl: "/sso-callback",
    });
  };

  return (
    <Button
      loading={fetchStatus === "fetching"}
      onClick={handleClick}
      {...props}
    >
      {children}
    </Button>
  );
}
