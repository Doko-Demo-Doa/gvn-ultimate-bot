"use client";

import { Badge, Paper, Stack, Table, Text } from "@mantine/core";
import type { IAuditLogItem, IDiscordChannel } from "~/types/types";
import MessageContentCell, { type ContentModalState } from "./message-content-cell";

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

type Props = {
  logs: IAuditLogItem[];
  channels: IDiscordChannel[];
  onExpandContent: (state: ContentModalState) => void;
};

export default function AuditLogTable({ logs, channels, onExpandContent }: Props) {
  return (
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
                  <Badge color={log.Action === "delete" ? "red" : "blue"}>
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
                  <MessageContentCell
                    title="Nội dung trước"
                    content={log.BeforeContent}
                    onExpand={onExpandContent}
                  />
                </Table.Td>
                <Table.Td>
                  <MessageContentCell
                    title="Nội dung sau"
                    content={log.AfterContent}
                    onExpand={onExpandContent}
                  />
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
  );
}
