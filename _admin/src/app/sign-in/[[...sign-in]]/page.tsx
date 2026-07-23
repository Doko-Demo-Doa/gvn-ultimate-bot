import { Stack, Text, Title } from "@mantine/core";
import DiscordAuthButton from "~/components/auth/discord-auth-button";

export default function SignInPage() {
  return (
    <Stack
      align="center"
      justify="center"
      gap="md"
      style={{ minHeight: "100vh" }}
    >
      <Title order={2}>DautoVN Bot CP</Title>
      <Text c="dimmed">
        Continue with Discord to access the admin dashboard.
      </Text>
      <DiscordAuthButton />
    </Stack>
  );
}
