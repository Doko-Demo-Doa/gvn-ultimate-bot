"use client";

import {
  Avatar,
  Badge,
  Button,
  Group,
  Loader,
  NativeSelect,
  Select,
  Stack,
  Table,
  Text,
  TextInput,
  Title,
} from "@mantine/core";
import { notifications } from "@mantine/notifications";
import { useState } from "react";
import {
  useAssignRoleMutation,
  useDiscordRoles,
  useRevokeRoleMutation,
  useRoleAssignments,
  useSearchDiscordMembers,
} from "~/hooks/api-hooks";
import MasterLayout from "~/layouts/master-layout";

function renderMemberOption({ option }: { option: any }) {
  return (
    <Group gap="xs">
      {option.image_url && (
        <Avatar src={option.image_url} size={24} radius="xl" />
      )}
      <span>{option.label}</span>
    </Group>
  );
}

export default function RoleManagerPage() {
  const { data: rolesData, isLoading: rolesLoading } = useDiscordRoles();
  const {
    data: assignmentsData,
    isLoading: assignmentsLoading,
    refetch,
  } = useRoleAssignments();
  const { mutateAsync: assignRole, isPending: isAssigning } =
    useAssignRoleMutation();
  const { mutateAsync: revokeRole, isPending: isRevoking } =
    useRevokeRoleMutation();

  const [memberSearch, setMemberSearch] = useState("");
  const [selectedMember, setSelectedMember] = useState<string | null>(null);
  const [roleNativeId, setRoleNativeId] = useState("");
  const [durationValue, setDurationValue] = useState("1");
  const [durationUnit, setDurationUnit] = useState("d");

  const { data: membersData, isLoading: membersSearching } =
    useSearchDiscordMembers(memberSearch);

  const roles = rolesData?.data ?? [];
  const assignments = assignmentsData?.data ?? [];
  const members = membersData?.data ?? [];

  const isLoading = rolesLoading || assignmentsLoading;

  const memberSelectData = members.map((m) => ({
    value: m.native_id,
    label: m.nickname
      ? `${m.nickname} (${m.username})`
      : m.username,
    image_url: m.avatar,
  }));

  async function handleAssign() {
    if (!selectedMember || !roleNativeId) {
      notifications.show({
        color: "red",
        title: "Lỗi",
        message: "Vui lòng chọn user và role",
      });
      return;
    }

    const duration = `${durationValue}${durationUnit}`;

    try {
      await assignRole({
        user_native_id: selectedMember,
        role_native_id: roleNativeId,
        duration,
      });

      notifications.show({
        color: "green",
        title: "Thành công",
        message: `Đã gán role ${getRoleName(roleNativeId)} cho user (${duration})`,
      });

      setSelectedMember(null);
      setMemberSearch("");
      void refetch();
    } catch (err: any) {
      notifications.show({
        color: "red",
        title: "Lỗi",
        message: err?.message || "Không thể gán role",
      });
    }
  }

  async function handleRevoke(id: number) {
    try {
      await revokeRole(id);
      notifications.show({
        color: "green",
        title: "Thành công",
        message: "Đã thu hồi role",
      });
      void refetch();
    } catch (err: any) {
      notifications.show({
        color: "red",
        title: "Lỗi",
        message: err?.message || "Không thể thu hồi role",
      });
    }
  }

  function getRoleName(nativeId: string) {
    return roles.find((r) => r.NativeId === nativeId)?.Name || nativeId;
  }

  return (
    <MasterLayout>
      <Stack>
        <Title order={3}>Role Manager</Title>
        <Text>
          Gán role cho user với thời hạn. Hệ thống sẽ tự động thu hồi role khi
          hết hạn.
        </Text>

        {isLoading ? (
          <Loader />
        ) : (
          <>
            <Stack>
              <Title order={4}>Gán role mới</Title>
              <Select
                label="Discord User"
                placeholder="Tìm kiếm user..."
                searchable
                clearable
                value={selectedMember}
                onChange={setSelectedMember}
                onSearchChange={setMemberSearch}
                searchValue={memberSearch}
                data={memberSelectData}
                renderOption={renderMemberOption}
                rightSection={membersSearching ? <Loader size={16} /> : null}
                allowDeselect={false}
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
                value={roleNativeId}
                onChange={(e) => setRoleNativeId(e.currentTarget.value)}
              />
              <Group>
                <TextInput
                  label="Thời hạn"
                  type="number"
                  min={1}
                  value={durationValue}
                  onChange={(e) => setDurationValue(e.currentTarget.value)}
                  style={{ width: 120 }}
                />
                <NativeSelect
                  label="Đơn vị"
                  data={[
                    { value: "m", label: "Phút" },
                    { value: "h", label: "Giờ" },
                    { value: "d", label: "Ngày" },
                    { value: "w", label: "Tuần" },
                  ]}
                  value={durationUnit}
                  onChange={(e) => setDurationUnit(e.currentTarget.value)}
                  style={{ width: 120 }}
                />
              </Group>
              <Button loading={isAssigning} onClick={handleAssign}>
                Gán role
              </Button>
            </Stack>

            <Stack mt="xl">
              <Title order={4}>Danh sách gán role</Title>
              {assignments.length === 0 ? (
                <Text c="dimmed">Không có role nào được gán.</Text>
              ) : (
                <Table>
                  <Table.Thead>
                    <Table.Tr>
                      <Table.Th>User ID</Table.Th>
                      <Table.Th>Role</Table.Th>
                      <Table.Th>Gán lúc</Table.Th>
                      <Table.Th>Hết hạn</Table.Th>
                      <Table.Th>Trạng thái</Table.Th>
                      <Table.Th>Thời gian còn lại</Table.Th>
                      <Table.Th>Hành động</Table.Th>
                    </Table.Tr>
                  </Table.Thead>
                  <Table.Tbody>
                    {assignments.map((a) => (
                      <Table.Tr key={a.ID}>
                        <Table.Td>{a.UserNativeID}</Table.Td>
                        <Table.Td>{getRoleName(a.RoleNativeID)}</Table.Td>
                        <Table.Td>
                          {new Date(a.GrantedDate).toLocaleString()}
                        </Table.Td>
                        <Table.Td>
                          {new Date(a.ExpirationDate).toLocaleString()}
                        </Table.Td>
                        <Table.Td>
                          {a.Status === "active" ? (
                            <Badge color="green">Đang hoạt động</Badge>
                          ) : (
                            <Badge color="red">Đã hết hạn</Badge>
                          )}
                        </Table.Td>
                        <Table.Td>
                          {a.Status === "active" ? a.TimeRemaining : "—"}
                        </Table.Td>
                        <Table.Td>
                          {a.Status === "active" && (
                            <Button
                              size="xs"
                              color="red"
                              variant="outline"
                              loading={isRevoking}
                              onClick={() => handleRevoke(a.ID)}
                            >
                              Thu hồi
                            </Button>
                          )}
                        </Table.Td>
                      </Table.Tr>
                    ))}
                  </Table.Tbody>
                </Table>
              )}
            </Stack>
          </>
        )}
      </Stack>
    </MasterLayout>
  );
}
