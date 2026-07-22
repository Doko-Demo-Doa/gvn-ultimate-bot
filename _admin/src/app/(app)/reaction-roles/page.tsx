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
import { IconEdit, IconPlus, IconTrash } from "@tabler/icons-react";
import { useState } from "react";
import EmbedEditor from "~/components/embed-editor/embed-editor";
import { useDiscordChannels, useDiscordRoles } from "~/hooks/api-hooks";
import {
  useDeleteReactionRoleEmbed,
  usePublishReactionRoleEmbed,
  useUpsertReactionRoleEmbed,
  useReactionRoleEmbeds,
} from "~/hooks/use-reaction-roles";
import type {
  IDiscordRoleReactionEmbed,
  IReactionRoleMessagePayload,
} from "~/types/types";
import MasterLayout from "~/layouts/master-layout";

export default function ReactionRolesPage() {
  const { data: embedsData, isLoading, refetch } = useReactionRoleEmbeds();
  const { data: rolesData } = useDiscordRoles();
  const { data: channelsData } = useDiscordChannels();
  const { mutateAsync: publish, isPending: isPublishing } =
    usePublishReactionRoleEmbed();
  const { mutateAsync: upsert, isPending: isUpserting } =
    useUpsertReactionRoleEmbed();
  const { mutateAsync: deleteEmbed, isPending: isDeleting } =
    useDeleteReactionRoleEmbed();
  const [opened, { open, close }] = useDisclosure(false);
  const [editingEmbed, setEditingEmbed] =
    useState<IDiscordRoleReactionEmbed | null>(null);

  const embeds = embedsData?.data ?? [];
  const roles = rolesData?.data ?? [];
  const channels = channelsData?.data ?? [];

  async function handlePublish(values: IReactionRoleMessagePayload) {
    try {
      await publish(values);
      notifications.show({
        color: "green",
        title: "Thành công",
        message: "Đã gửi reaction role message lên Discord.",
      });
      close();
      setEditingEmbed(null);
      void refetch();
    } catch (err: any) {
      notifications.show({
        color: "red",
        title: "Lỗi",
        message: err?.message || "Không thể gửi message.",
      });
    }
  }

  async function handleEditSave(values: IReactionRoleMessagePayload) {
    if (!editingEmbed) return;
    try {
      await upsert({
        native_message_id: editingEmbed.NativeMessageId,
        name: values.message || "",
        payload: values,
        mode: values.mode,
        version: editingEmbed.Version + 1,
      });
      notifications.show({
        color: "green",
        title: "Thành công",
        message: "Đã cập nhật reaction role message trên Discord.",
      });
      close();
      setEditingEmbed(null);
      void refetch();
    } catch (err: any) {
      notifications.show({
        color: "red",
        title: "Lỗi",
        message: err?.message || "Không thể cập nhật message.",
      });
    }
  }

  async function handleDelete(id: number) {
    if (!confirm("Bạn có chắc muốn xóa reaction role message này?")) return;
    try {
      await deleteEmbed(id);
      notifications.show({
        color: "green",
        title: "Thành công",
        message: "Đã xóa reaction role message.",
      });
      void refetch();
    } catch (err: any) {
      notifications.show({
        color: "red",
        title: "Lỗi",
        message: err?.message || "Không thể xóa message.",
      });
    }
  }

  function handleEdit(embed: IDiscordRoleReactionEmbed) {
    setEditingEmbed(embed);
    open();
  }

  function handleCreate() {
    setEditingEmbed(null);
    open();
  }

  function handleClose() {
    close();
    setEditingEmbed(null);
  }

  // Parse stored payload for pre-filling the editor when editing
  let editPayload: IReactionRoleMessagePayload | undefined;
  if (editingEmbed?.Payload) {
    try {
      editPayload = JSON.parse(editingEmbed.Payload);
    } catch {
      editPayload = undefined;
    }
  }

  return (
    <MasterLayout>
      <Stack>
        <Paper>
          <Group justify="space-between">
            <Title order={3}>Your reaction role messages</Title>
            <Button leftSection={<IconPlus size={14} />} onClick={handleCreate}>
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
                  <Table.Th>Hành động</Table.Th>
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
                    <Table.Td>
                      <Group gap="xs">
                        <Button
                          size="xs"
                          variant="light"
                          leftSection={<IconEdit size={14} />}
                          onClick={() => handleEdit(e)}
                        >
                          Sửa
                        </Button>
                        <Button
                          size="xs"
                          color="red"
                          variant="light"
                          leftSection={<IconTrash size={14} />}
                          loading={isDeleting}
                          onClick={() => handleDelete(e.ID)}
                        >
                          Xóa
                        </Button>
                      </Group>
                    </Table.Td>
                  </Table.Tr>
                ))}
              </Table.Tbody>
            </Table>
          )}
        </Paper>
      </Stack>

      <Modal
        opened={opened}
        onClose={handleClose}
        title={
          editingEmbed
            ? "Sửa reaction role message"
            : "Tạo reaction role message"
        }
        size="xl"
        padding="xl"
      >
        <EmbedEditor
          roles={roles}
          channels={channels}
          onPublish={editingEmbed ? handleEditSave : handlePublish}
          isPublishing={editingEmbed ? isUpserting : isPublishing}
          initialPayload={editPayload}
        />
      </Modal>
    </MasterLayout>
  );
}
