"use client";

import { Button, Group, Loader, Paper, Select, Stack } from "@mantine/core";
import { DatePickerInput } from "@mantine/dates";
import { IconCalendar, IconClearAll, IconSearch } from "@tabler/icons-react";
import { useEffect, useMemo, useState } from "react";
import { renderMemberOption, toMemberSelectData } from "~/app/_components/member-select";
import { useSearchDiscordMembers } from "~/hooks/api-hooks";
import type { IDiscordChannel } from "~/types/types";

type Props = {
  channels: IDiscordChannel[];
  datePreset: string;
  onDatePresetChange: (val: string | null) => void;
  fromDate: string;
  onFromDateChange: (val: string | null) => void;
  toDate: string;
  onToDateChange: (val: string | null) => void;
  channelId: string;
  onChannelIdChange: (val: string | null) => void;
  authorId: string;
  onAuthorIdChange: (val: string | null) => void;
  onSearch: () => void;
  onReset: () => void;
};

export default function AuditLogFilters({
  channels,
  datePreset,
  onDatePresetChange,
  fromDate,
  onFromDateChange,
  toDate,
  onToDateChange,
  channelId,
  onChannelIdChange,
  authorId,
  onAuthorIdChange,
  onSearch,
  onReset,
}: Props) {
  const [authorSearch, setAuthorSearch] = useState("");

  // Keep the local search text in sync when the selection is cleared
  // elsewhere (e.g. the page-level "Đặt lại" reset button).
  useEffect(() => {
    if (!authorId) setAuthorSearch("");
  }, [authorId]);

  const { data: authorSearchData, isLoading: authorSearching } =
    useSearchDiscordMembers(authorSearch);
  // Resolves the currently selected author by id so its name still renders
  // after a page refresh, when authorSearch is empty and the live search
  // above hasn't returned anything yet.
  const { data: selectedAuthorData } = useSearchDiscordMembers(authorId);

  const authorSelectData = useMemo(() => {
    const merged = [
      ...(selectedAuthorData?.data ?? []),
      ...(authorSearchData?.data ?? []),
    ];
    const seen = new Set<string>();
    const deduped = merged.filter((m) => {
      if (seen.has(m.native_id)) return false;
      seen.add(m.native_id);
      return true;
    });
    return toMemberSelectData(deduped);
  }, [selectedAuthorData, authorSearchData]);

  return (
    <Paper p="md" withBorder>
      <Stack gap="md">
        <Group align="flex-end">
          <Select
            label="Ngày"
            placeholder="Chọn khoảng thời gian"
            value={datePreset || ""}
            onChange={(val) => {
              onDatePresetChange(val);
              if (val !== "custom") {
                onFromDateChange(null);
                onToDateChange(null);
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
                onChange={(val) => onFromDateChange(val ?? null)}
                leftSection={<IconCalendar size={16} />}
                style={{ width: 160 }}
              />
              <DatePickerInput
                label="Đến ngày"
                placeholder="Chọn ngày"
                value={toDate || null}
                onChange={(val) => onToDateChange(val ?? null)}
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
            onChange={onChannelIdChange}
            data={channels.map((ch) => ({
              value: ch.id,
              label: `#${ch.name}`,
            }))}
            style={{ width: 220 }}
          />
          <Select
            label="Người gửi"
            placeholder="Tìm kiếm user..."
            searchable
            clearable
            value={authorId || null}
            onChange={onAuthorIdChange}
            onSearchChange={setAuthorSearch}
            searchValue={authorSearch}
            data={authorSelectData}
            renderOption={renderMemberOption}
            rightSection={authorSearching ? <Loader size={16} /> : null}
            style={{ width: 220 }}
          />
          <Button leftSection={<IconSearch size={14} />} onClick={onSearch}>
            Tìm kiếm
          </Button>
          <Button
            variant="light"
            color="gray"
            leftSection={<IconClearAll size={14} />}
            onClick={onReset}
          >
            Đặt lại
          </Button>
        </Group>
      </Stack>
    </Paper>
  );
}
