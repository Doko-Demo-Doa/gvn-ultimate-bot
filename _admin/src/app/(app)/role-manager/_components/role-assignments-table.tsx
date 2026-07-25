"use client";

import { Badge, Button, Stack, Table, Text, Title } from "@mantine/core";
import { notifications } from "@mantine/notifications";
import { useRevokeRoleMutation } from "~/hooks/api-hooks";
import type { IDiscordRole, IDiscordUserRoleAssignment } from "~/types/types";

type Props = {
  assignments: IDiscordUserRoleAssignment[];
  roles: IDiscordRole[];
  onRevoked: () => void;
};

export default function RoleAssignmentsTable({
  assignments,
  roles,
  onRevoked,
}: Props) {
  const { mutateAsync: revokeRole, isPending: isRevoking } =
    useRevokeRoleMutation();

  function getRoleName(nativeId: string) {
    return roles.find((r) => r.NativeId === nativeId)?.Name || nativeId;
  }

  async function handleRevoke(id: number) {
    try {
      await revokeRole(id);
      notifications.show({
        color: "green",
        title: "Thành công",
        message: "Đã thu hồi role",
      });
      onRevoked();
    } catch (err: any) {
      notifications.show({
        color: "red",
        title: "Lỗi",
        message: err?.message || "Không thể thu hồi role",
      });
    }
  }

  return (
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
  );
}
