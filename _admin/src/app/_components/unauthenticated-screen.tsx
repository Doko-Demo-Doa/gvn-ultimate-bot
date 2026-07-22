import { SignInButton, SignUpButton } from "@clerk/nextjs";
import { Button, Group, Stack, Text, Title } from "@mantine/core";

export default function UnauthenticatedScreen() {
  return (
    <Stack
      align="center"
      justify="center"
      gap="md"
      style={{ minHeight: "100vh" }}
    >
      <Title order={2}>DautoVN Bot CP</Title>
      <Text c="dimmed">Sign in to access the admin dashboard.</Text>
      <Group>
        <SignInButton>
          <Button variant="default">Sign in</Button>
        </SignInButton>
        <SignUpButton>
          <Button>Sign up</Button>
        </SignUpButton>
      </Group>
    </Stack>
  );
}
