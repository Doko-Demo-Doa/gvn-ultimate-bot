"use client";

import {
  Badge,
  Box,
  Button,
  Group,
  Paper,
  Select,
  Stack,
  Table,
  Text,
  TextInput,
  Title,
} from "@mantine/core";
import { DatePickerInput } from "@mantine/dates";
import { notifications } from "@mantine/notifications";
import {
  IconCalendar,
  IconClearAll,
  IconSearch,
  IconTrash,
} from "@tabler/icons-react";
import { useQueryState } from "nuqs";
import {
  useAuditLogs,
  useClearAuditLogs,
  useDiscordChannels,
} from "~/hooks/api-hooks";
import MasterLayout from "~/layouts/master-layout";

function formatDateTime(iso: string) {
  const d = new Date(iso);
  return d.toLocaleString("vi-VN", {
    year: "numeric",
    month: "2-digit",
    day: "2-digit",
    hour: "2-digit",
    minute: "2-digit",
  });
}

function parseAttachments(jsonStr: string): string[] {
  try {
    const arr = JSON.parse(jsonStr);
    return Array.isArray(arr) ? arr : [];
  } catch {
    return [];
  }
}

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
  const [authorName, setAuthorName] = useQueryState("author", {
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
    author_name: authorName || undefined,
  });

  const { mutateAsync: clearLogs, isPending: isClearing } = useClearAuditLogs();

  const logs = data?.data?.items ?? [];
  const total = data?.data?.total ?? 0;

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
    setAuthorName(null);
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

        <Paper p="md" withBorder>
          <Stack gap="md">
            <Group align="flex-end">
              <Select
                label="Ngày"
                placeholder="Chọn khoảng thời gian"
                value={datePreset || ""}
                onChange={(val) => {
                  setDatePreset(val);
                  if (val !== "custom") {
                    setFromDate(null);
                    setToDate(null);
                  }
                }}
                data={[
                  { value: "today", label: "Hôm nay" },
                  { value: "yesterday", label: "Hôm qua" },
                  { value: "week", label: "7 ngày qua" },
                  { value: "custom", label: "Tùy chọn" },
                ]}
                style={{ width: 200 }}
                allowDeselect
              />
              {datePreset === "custom" && (
                <>
                  <DatePickerInput
                    label="Từ ngày"
                    placeholder="Chọn ngày"
                    value={fromDate || null}
                    onChange={(val) => setFromDate(val ?? null)}
                    leftSection={<IconCalendar size={16} />}
                    style={{ width: 160 }}
                  />
                  <DatePickerInput
                    label="Đến ngày"
                    placeholder="Chọn ngày"
                    value={toDate || null}
                    onChange={(val) => setToDate(val ?? null)}
                    leftSection={<IconCalendar size={16} />}
                    style={{ width: 160 }}
                  />
                </>
              )}
              <Select
                label="Channel"
                placeholder="Tất cả channel"
                searchable
                clearable
                value={channelId || ""}
                onChange={(val) => setChannelId(val)}
                data={channels.map((ch) => ({
                  value: ch.id,
                  label: `#${ch.name}`,
                }))}
                style={{ width: 220 }}
              />
              <TextInput
                label="Người gửi"
                placeholder="Tên người gửi..."
                value={authorName}
                onChange={(e) => setAuthorName(e.currentTarget.value)}
                style={{ width: 200 }}
              />
              <Button
                leftSection={<IconSearch size={14} />}
                onClick={handleSearch}
              >
                Tìm kiếm
              </Button>
              <Button
                variant="light"
                color="gray"
                leftSection={<IconClearAll size={14} />}
                onClick={handleReset}
              >
                Đặt lại
              </Button>
            </Group>
          </Stack>
        </Paper>

        {isLoading ? (
          <Text>Đang tải...</Text>
        ) : logs.length === 0 ? (
          <Text c="dimmed">Không có log nào.</Text>
        ) : (
          <Paper withBorder>
            <Table striped>
              <Table.Thead>
                <Table.Tr>
                  <Table.Th>Thời gian</Table.Th>
                  <Table.Th>Hành động</Table.Th>
                  <Table.Th>Channel</Table.Th>
                  <Table.Th>Người gửi</Table.Th>
                  <Table.Th>Nội dung trước</Table.Th>
                  <Table.Th>Nội dung sau</Table.Th>
                  <Table.Th>Attachments</Table.Th>
                </Table.Tr>
              </Table.Thead>
              <Table.Tbody>
                {logs.map((log) => {
                  const attachments = parseAttachments(log.Attachments);
                  return (
                    <Table.Tr key={log.ID}>
                      <Table.Td>{formatDateTime(log.CreatedAt)}</Table.Td>
                      <Table.Td>
                        <Badge
                          color={log.Action === "delete" ? "red" : "blue"}
                        >
                          {log.Action === "delete" ? "Xóa" : "Sửa"}
                        </Badge>
                      </Table.Td>
                      <Table.Td>
                        <Text size="sm" c="dimmed">
                          #{channels.find((c) => c.id === log.ChannelId)?.name || log.ChannelId}
                        </Text>
                      </Table.Td>
                      <Table.Td>
                        <Text size="sm" fw={500}>
                          {log.AuthorName || "—"}
                        </Text>
                      </Table.Td>
                      <Table.Td>
                        <Box
                          style={{
                            maxWidth: 300,
                            whiteSpace: "pre-wrap",
                            wordBreak: "break-word",
                          }}
                        >
                          <Text size="sm">{log.BeforeContent || "—"}</Text>
                        </Box>
                      </Table.Td>
                      <Table.Td>
                        <Box
                          style={{
                            maxWidth: 300,
                            whiteSpace: "pre-wrap",
                            wordBreak: "break-word",
                          }}
                        >
                          <Text size="sm">{log.AfterContent || "—"}</Text>
                        </Box>
                      </Table.Td>
                      <Table.Td>
                        {attachments.length > 0 && (
                          <Stack gap={4}>
                            {attachments.map((url, i) => (
                              <a
                                key={i}
                                href={url}
                                target="_blank"
                                rel="noreferrer"
                                style={{ fontSize: 12 }}
                              >
                                Attachment {i + 1}
                              </a>
                            ))}
                          </Stack>
                        )}
                      </Table.Td>
                    </Table.Tr>
                  );
                })}
              </Table.Tbody>
            </Table>
          </Paper>
        )}

        {total > PAGE_SIZE && (
          <Group justify="center" mt="md">
            <Button
              variant="light"
              disabled={offset === 0}
              onClick={() => setPage(Math.max(0, offset - PAGE_SIZE))}
            >
              Trang trước
            </Button>
            <Text size="sm">
              Hiển thị {offset + 1} – {Math.min(offset + PAGE_SIZE, total)} /{" "}
              {total}
            </Text>
            <Button
              variant="light"
              disabled={offset + PAGE_SIZE >= total}
              onClick={() => setPage(offset + PAGE_SIZE)}
            >
              Trang sau
            </Button>
          </Group>
        )}
      </Stack>
    </MasterLayout>
  );
}
