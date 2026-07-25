"use client";

import {
  Button,
  Group,
  Loader,
  NativeSelect,
  Select,
  Stack,
  TextInput,
  Title,
} from "@mantine/core";
import { useForm } from "@mantine/form";
import { notifications } from "@mantine/notifications";
import { useState } from "react";
import { renderMemberOption, toMemberSelectData } from "~/app/_components/member-select";
import { useAssignRoleMutation, useSearchDiscordMembers } from "~/hooks/api-hooks";
import type { IDiscordRole } from "~/types/types";

type Props = {
  roles: IDiscordRole[];
  onAssigned: () => void;
};

type FormValues = {
  memberNativeId: string | null;
  roleNativeId: string;
  durationValue: string;
  durationUnit: string;
};

export default function AssignRoleForm({ roles, onAssigned }: Props) {
  const [memberSearch, setMemberSearch] = useState("");

  const { data: membersData, isLoading: membersSearching } =
    useSearchDiscordMembers(memberSearch);
  const { mutateAsync: assignRole, isPending: isAssigning } =
    useAssignRoleMutation();

  const memberSelectData = toMemberSelectData(membersData?.data ?? []);

  const form = useForm<FormValues>({
    initialValues: {
      memberNativeId: null,
      roleNativeId: "",
      durationValue: "1",
      durationUnit: "d",
    },
  });

  function getRoleName(nativeId: string) {
    return roles.find((r) => r.NativeId === nativeId)?.Name || nativeId;
  }

  async function handleSubmit(values: FormValues) {
    if (!values.memberNativeId || !values.roleNativeId) {
      notifications.show({
        color: "red",
        title: "Lỗi",
        message: "Vui lòng chọn user và role",
      });
      return;
    }

    const duration = `${values.durationValue}${values.durationUnit}`;

    try {
      await assignRole({
        user_native_id: values.memberNativeId,
        role_native_id: values.roleNativeId,
        duration,
      });

      notifications.show({
        color: "green",
        title: "Thành công",
        message: `Đã gán role ${getRoleName(values.roleNativeId)} cho user (${duration})`,
      });

      form.reset();
      setMemberSearch("");
      onAssigned();
    } catch (err: any) {
      notifications.show({
        color: "red",
        title: "Lỗi",
        message: err?.message || "Không thể gán role",
      });
    }
  }

  return (
    <form onSubmit={form.onSubmit(handleSubmit)}>
      <Stack>
        <Title order={4}>Gán role mới</Title>
        <Select
          label="Discord User"
          placeholder="Tìm kiếm user..."
          searchable
          clearable
          onSearchChange={setMemberSearch}
          searchValue={memberSearch}
          data={memberSelectData}
          renderOption={renderMemberOption}
          rightSection={membersSearching ? <Loader size={16} /> : null}
          allowDeselect={false}
          {...form.getInputProps("memberNativeId")}
        />
        <NativeSelect
          label="Role"
          data={[
            { value: "", label: "-- Chọn role --", disabled: true },
            ...roles.map((r) => ({
              value: r.NativeId,
              label: r.Name,
            })),
          ]}
          {...form.getInputProps("roleNativeId")}
        />
        <Group>
          <TextInput
            label="Thời hạn"
            type="number"
            min={1}
            style={{ width: 120 }}
            {...form.getInputProps("durationValue")}
          />
          <NativeSelect
            label="Đơn vị"
            data={[
              { value: "m", label: "Phút" },
              { value: "h", label: "Giờ" },
              { value: "d", label: "Ngày" },
              { value: "w", label: "Tuần" },
            ]}
            style={{ width: 120 }}
            {...form.getInputProps("durationUnit")}
          />
        </Group>
        <Button type="submit" loading={isAssigning}>
          Gán role
        </Button>
      </Stack>
    </form>
  );
}
