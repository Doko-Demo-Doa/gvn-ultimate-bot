"use client";

import { Divider, Paper, Title } from "@mantine/core";
import EmbedEditor from "~/components/embed-editor/embed-editor";
import MasterLayout from "~/layouts/master-layout";

export default function ReactionRolesPage() {
  return (
    <MasterLayout>
      <Paper>
        <Title order={3}>Your reaction role messages</Title>
        <Divider my="sm" />

        <EmbedEditor messageId="" />
      </Paper>
    </MasterLayout>
  );
}
