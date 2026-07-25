import { Avatar, Group } from "@mantine/core";
import type { IDiscordMember } from "~/types/types";

export function renderMemberOption({ option }: { option: any }) {
  return (
    <Group gap="xs">
      {option.image_url && (
        <Avatar src={option.image_url} size={24} radius="xl" />
      )}
      <span>{option.label}</span>
    </Group>
  );
}

export function toMemberSelectData(members: IDiscordMember[]) {
  return members.map((m) => ({
    value: m.native_id,
    label: m.nickname ? `${m.nickname} (${m.username})` : m.username,
    image_url: m.avatar,
  }));
}
