"use client";

import { Button, Group, Modal, Pagination, Stack, Text, Title } from "@mantine/core";
import { notifications } from "@mantine/notifications";
import { IconTrash } from "@tabler/icons-react";
import { useQueryState } from "nuqs";
import { useState } from "react";
import { useAuditLogs, useClearAuditLogs, useDiscordChannels } from "~/hooks/api-hooks";
import MasterLayout from "~/layouts/master-layout";
import AuditLogFilters from "./_components/audit-log-filters";
import AuditLogTable from "./_components/audit-log-table";
import type { ContentModalState } from "./_components/message-content-cell";

const PAGE_SIZE = 50;

export default function ServerLogPage() {
  const { data: channelsData } = useDiscordChannels();
  const channels = channelsData?.data ?? [];

  // URL-based filter state via nuqs
  const [datePreset, setDatePreset] = useQueryState("preset", {
    defaultValue: "",
  });
  const [fromDate, setFromDate] = useQueryState("from", {
    defaultValue: "",
  });
  const [toDate, setToDate] = useQueryState("to", {
    defaultValue: "",
  });
  const [channelId, setChannelId] = useQueryState("channel", {
    defaultValue: "",
  });
  const [authorId, setAuthorId] = useQueryState("author_id", {
    defaultValue: "",
  });
  const [offsetStr, setOffsetStr] = useQueryState("offset", {
    defaultValue: "0",
  });

  const offset = Number.parseInt(offsetStr, 10) || 0;

  function getQueryDates() {
    const today = new Date();
    today.setHours(0, 0, 0, 0);

    switch (datePreset) {
      case "today": {
        const tomorrow = new Date(today);
        tomorrow.setDate(tomorrow.getDate() + 1);
        return {
          from: today.toISOString(),
          to: tomorrow.toISOString(),
        };
      }
      case "yesterday": {
        const yesterday = new Date(today);
        yesterday.setDate(yesterday.getDate() - 1);
        return {
          from: yesterday.toISOString(),
          to: today.toISOString(),
        };
      }
      case "week": {
        const weekAgo = new Date(today);
        weekAgo.setDate(weekAgo.getDate() - 7);
        return {
          from: weekAgo.toISOString(),
          to: new Date().toISOString(),
        };
      }
      case "custom":
      default: {
        const from = fromDate ? new Date(fromDate) : null;
        const to = toDate ? new Date(toDate) : null;
        if (to) {
          to.setHours(23, 59, 59, 999);
        }
        return {
          from: from?.toISOString() ?? undefined,
          to: to?.toISOString() ?? undefined,
        };
      }
    }
  }

  const dates = getQueryDates();

  const { data, isLoading, refetch } = useAuditLogs({
    limit: PAGE_SIZE,
    offset,
    from_date: dates.from,
    to_date: dates.to,
    channel_id: channelId || undefined,
    author_id: authorId || undefined,
  });

  const { mutateAsync: clearLogs, isPending: isClearing } = useClearAuditLogs();

  const [contentModal, setContentModal] = useState<ContentModalState>(null);

  const logs = data?.data?.items ?? [];
  const total = data?.data?.total ?? 0;
  const totalPages = Math.max(1, Math.ceil(total / PAGE_SIZE));
  const currentPage = Math.floor(offset / PAGE_SIZE) + 1;

  async function handleClear() {
    if (!confirm("Bạn có chắc muốn xóa toàn bộ audit log?")) return;
    try {
      await clearLogs();
      notifications.show({
        color: "green",
        title: "Thành công",
        message: "Đã xóa toàn bộ audit log.",
      });
      refetch();
    } catch (err: any) {
      notifications.show({
        color: "red",
        title: "Lỗi",
        message: err?.message || "Không thể xóa audit log.",
      });
    }
  }

  function handleSearch() {
    setOffsetStr("0");
  }

  function handleReset() {
    setDatePreset(null);
    setFromDate(null);
    setToDate(null);
    setChannelId(null);
    setAuthorId(null);
    setOffsetStr("0");
  }

  function setPage(newOffset: number) {
    setOffsetStr(String(newOffset));
  }

  return (
    <MasterLayout>
      <Stack>
        <Group justify="space-between">
          <Title order={3}>Server Message Audit Log</Title>
          <Button
            color="red"
            variant="light"
            leftSection={<IconTrash size={14} />}
            loading={isClearing}
            onClick={handleClear}
          >
            Xóa toàn bộ log
          </Button>
        </Group>

        <AuditLogFilters
          channels={channels}
          datePreset={datePreset}
          onDatePresetChange={setDatePreset}
          fromDate={fromDate}
          onFromDateChange={setFromDate}
          toDate={toDate}
          onToDateChange={setToDate}
          channelId={channelId}
          onChannelIdChange={setChannelId}
          authorId={authorId}
          onAuthorIdChange={setAuthorId}
          onSearch={handleSearch}
          onReset={handleReset}
        />

        {isLoading ? (
          <Text>Đang tải...</Text>
        ) : logs.length === 0 ? (
          <Text c="dimmed">Không có log nào.</Text>
        ) : (
          <AuditLogTable
            logs={logs}
            channels={channels}
            onExpandContent={setContentModal}
          />
        )}

        {total > PAGE_SIZE && (
          <Group justify="center" mt="md">
            <Pagination
              value={currentPage}
              onChange={(page) => setPage((page - 1) * PAGE_SIZE)}
              total={totalPages}
            />
            <Text size="sm" c="dimmed">
              Hiển thị {offset + 1} – {Math.min(offset + PAGE_SIZE, total)} /{" "}
              {total}
            </Text>
          </Group>
        )}
      </Stack>

      <Modal
        opened={!!contentModal}
        onClose={() => setContentModal(null)}
        title={contentModal?.title}
        size="lg"
      >
        <Text size="sm" style={{ whiteSpace: "pre-wrap", wordBreak: "break-word" }}>
          {contentModal?.content}
        </Text>
      </Modal>
    </MasterLayout>
  );
}
