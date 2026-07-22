"use client";

import {
  Button,
  Divider,
  Group,
  Modal,
  Paper,
  Stack,
  Table,
  Text,
  Title,
} from "@mantine/core";
import { useDisclosure } from "@mantine/hooks";
import { notifications } from "@mantine/notifications";
import { IconPlus } from "@tabler/icons-react";
import EmbedEditor from "~/components/embed-editor/embed-editor";
import { useDiscordRoles } from "~/hooks/api-hooks";
import {
  usePublishReactionRoleEmbed,
  useReactionRoleEmbeds,
} from "~/hooks/use-reaction-roles";
import MasterLayout from "~/layouts/master-layout";

export default function ReactionRolesPage() {
  const { data: embedsData, isLoading, refetch } = useReactionRoleEmbeds();
  const { data: rolesData } = useDiscordRoles();
  const { mutateAsync: publish, isPending: isPublishing } =
    usePublishReactionRoleEmbed();
  const [opened, { open, close }] = useDisclosure(false);

  const embeds = embedsData?.data ?? [];
  const roles = rolesData?.data ?? [];

  async function handlePublish(values: any) {
    try {
      await publish(values);
      notifications.show({
        color: "green",
        title: "Thành công",
        message: "Đã gửi reaction role message lên Discord.",
      });
      close();
      void refetch();
    } catch (err: any) {
      notifications.show({
        color: "red",
        title: "Lỗi",
        message: err?.message || "Không thể gửi message.",
      });
    }
  }

  return (
    <MasterLayout>
      <Stack>
        <Paper>
          <Group justify="space-between">
            <Title order={3}>Your reaction role messages</Title>
            <Button leftSection={<IconPlus size={14} />} onClick={open}>
              Tạo mới
            </Button>
          </Group>
          <Divider my="sm" />

          {isLoading ? (
            <Text>Đang tải...</Text>
          ) : embeds.length === 0 ? (
            <Text c="dimmed">Chưa có reaction role message nào.</Text>
          ) : (
            <Table>
              <Table.Thead>
                <Table.Tr>
                  <Table.Th>ID</Table.Th>
                  <Table.Th>Message ID</Table.Th>
                  <Table.Th>Tên</Table.Th>
                  <Table.Th>Chế độ</Table.Th>
                  <Table.Th>Phiên bản</Table.Th>
                </Table.Tr>
              </Table.Thead>
              <Table.Tbody>
                {embeds.map((e) => (
                  <Table.Tr key={e.ID}>
                    <Table.Td>{e.ID}</Table.Td>
                    <Table.Td>{e.NativeMessageId}</Table.Td>
                    <Table.Td>{e.Name || "—"}</Table.Td>
                    <Table.Td>{e.Mode || "default"}</Table.Td>
                    <Table.Td>{e.Version}</Table.Td>
                  </Table.Tr>
                ))}
              </Table.Tbody>
            </Table>
          )}
        </Paper>
      </Stack>

      <Modal opened={opened} onClose={close} title="Tạo reaction role message" size="xl">
        <EmbedEditor
          roles={roles}
          onPublish={handlePublish}
          isPublishing={isPublishing}
        />
      </Modal>
    </MasterLayout>
  );
}
