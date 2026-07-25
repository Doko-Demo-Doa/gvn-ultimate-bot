"use client";

import { Box, Text, UnstyledButton } from "@mantine/core";

const CONTENT_PREVIEW_LENGTH = 120;

export type ContentModalState = {
  title: string;
  content: string;
} | null;

type Props = {
  title: string;
  content: string;
  onExpand: (state: ContentModalState) => void;
};

export default function MessageContentCell({ title, content, onExpand }: Props) {
  if (!content) {
    return <Text size="sm">—</Text>;
  }

  const isLong = content.length > CONTENT_PREVIEW_LENGTH;
  const preview = isLong
    ? `${content.slice(0, CONTENT_PREVIEW_LENGTH)}...`
    : content;

  return (
    <Box
      style={{
        maxWidth: 300,
        whiteSpace: "pre-wrap",
        wordBreak: "break-word",
      }}
    >
      <Text size="sm">{preview}</Text>
      {isLong && (
        <UnstyledButton
          onClick={() => onExpand({ title, content })}
          c="blue"
          style={{ fontSize: 12 }}
        >
          Xem thêm
        </UnstyledButton>
      )}
    </Box>
  );
}
