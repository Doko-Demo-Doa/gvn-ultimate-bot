"use client";

import { Button, Group, Text } from "@mantine/core";
import { notifications } from "@mantine/notifications";
import { IconRefresh } from "@tabler/icons-react";
import { useLastRoleSync, useSyncDiscordRoles } from "~/hooks/api-hooks";

function formatDateTime(iso: string) {
  return new Date(iso).toLocaleString();
}

type Props = {
  onSynced: () => void;
};

export default function RoleSyncBar({ onSynced }: Props) {
  const { data: lastRoleSync, refetch: refetchLastRoleSync } =
    useLastRoleSync();
  const { mutateAsync: syncRoles, isPending: isSyncingRoles } =
    useSyncDiscordRoles();

  const lastRoleSyncLog = lastRoleSync?.data;

  async function handleSyncRoles() {
    try {
      const resp = await syncRoles();
      notifications.show({
        color: "green",
        title: "Thành công",
        message: `Đã đồng bộ ${resp.data.synced_count} role, xóa ${resp.data.removed_count} role không còn tồn tại.`,
      });
      onSynced();
      void refetchLastRoleSync();
    } catch (err: any) {
      notifications.show({
        color: "red",
        title: "Lỗi",
        message: err?.message || "Không thể đồng bộ role",
      });
    }
  }

  return (
    <Group justify="space-between">
      <Button
        variant="light"
        leftSection={<IconRefresh size={14} />}
        loading={isSyncingRoles}
        onClick={handleSyncRoles}
      >
        Đồng bộ role từ Discord
      </Button>
      <Text size="sm" c="dimmed">
        {lastRoleSyncLog
          ? `Đồng bộ lần cuối ${formatDateTime(lastRoleSyncLog.CreatedAt)} (${lastRoleSyncLog.Status})`
          : "Chưa đồng bộ lần nào"}
      </Text>
    </Group>
  );
}
