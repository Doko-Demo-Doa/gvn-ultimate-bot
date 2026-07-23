import { UserButton } from "@clerk/nextjs";
import { Stack, Text, Title } from "@mantine/core";

export default function ForbiddenScreen({ reason }: { reason: string }) {
  return (
    <Stack
      align="center"
      justify="center"
      gap="md"
      style={{ minHeight: "100vh" }}
    >
      <Title order={2}>Access denied</Title>
      <Text c="dimmed">{reason}</Text>
      <UserButton />
    </Stack>
  );
}
