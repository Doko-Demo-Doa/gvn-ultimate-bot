"use client";

import { Button, Group, Loader, Space, Stack, Switch, Text, Title } from "@mantine/core";
import { notifications } from "@mantine/notifications";
import { IconRefresh } from "@tabler/icons-react";
import {
  ModuleActivationStatus,
  useAppModuleEnabler,
  useAppModules,
  useLastUserSync,
  useSyncDiscordUsers,
} from "~/hooks/api-hooks";

function formatDateTime(iso: string) {
  return new Date(iso).toLocaleString();
}

const EnabledModules = () => {
  const { data, refetch } = useAppModules();
  const { mutateAsync } = useAppModuleEnabler();

  const { data: lastSync, refetch: refetchLastSync } = useLastUserSync();
  const { mutateAsync: syncUsers, isPending: isSyncing } =
    useSyncDiscordUsers();

  const lastSyncLog = lastSync?.data;

  async function handleSyncUsers() {
    try {
      const resp = await syncUsers();
      notifications.show({
        color: "green",
        title: "Success",
        message: `Synced ${resp.data.synced_count} Discord users.`,
      });
      refetchLastSync();
    } catch (err: any) {
      notifications.show({
        color: "red",
        title: "Error",
        message: err?.message || "Failed to sync Discord users.",
      });
    }
  }

  if (!data) {
    return <Loader />;
  }

  return (
    <>
      <Title order={3}>Enabled Modules</Title>
      <Space h="lg" />
      <Stack>
        {data?.data.map((module) => (
          <Switch
            key={module.ID}
            defaultChecked={!!module.IsActivated}
            label={module.ModuleLabel || ""}
            onChange={async (event) => {
              await mutateAsync({
                module_id: module.ID,
                is_activated: event.target.checked
                  ? ModuleActivationStatus.ENABLED
                  : ModuleActivationStatus.DISABLED,
              });

              refetch();
            }}
          />
        ))}
      </Stack>

      <Space h="xl" />
      <Title order={3}>Discord User Sync</Title>
      <Space h="lg" />
      <Group>
        <Button
          leftSection={<IconRefresh size={14} />}
          loading={isSyncing}
          onClick={handleSyncUsers}
        >
          Sync Discord Users
        </Button>
        <Text size="sm" c="dimmed">
          {lastSyncLog
            ? `Last synced ${formatDateTime(lastSyncLog.CreatedAt)} (${lastSyncLog.Status})`
            : "Never synced"}
        </Text>
      </Group>
    </>
  );
};

export default EnabledModules;
