"use client";

import { Box, Divider, Paper, Title } from "@mantine/core";
import EmbedEditor from "~/components/embed-editor/embed-editor";

export default function ReactionRolesPage() {
  return (
    <Box p="lg">
      <Paper>
        <Title order={3}>Your reaction role messages</Title>
        <Divider my="sm" />

        <EmbedEditor messageId="" />
      </Paper>
    </Box>
  );
}
