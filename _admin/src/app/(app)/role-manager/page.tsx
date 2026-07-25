"use client";

import { Loader, Stack, Text, Title } from "@mantine/core";
import { useDiscordRoles, useRoleAssignments } from "~/hooks/api-hooks";
import MasterLayout from "~/layouts/master-layout";
import AssignRoleForm from "./_components/assign-role-form";
import RoleAssignmentsTable from "./_components/role-assignments-table";
import RoleSyncBar from "./_components/role-sync-bar";

export default function RoleManagerPage() {
  const {
    data: rolesData,
    isLoading: rolesLoading,
    refetch: refetchRoles,
  } = useDiscordRoles();
  const {
    data: assignmentsData,
    isLoading: assignmentsLoading,
    refetch: refetchAssignments,
  } = useRoleAssignments();

  const roles = rolesData?.data ?? [];
  const assignments = assignmentsData?.data ?? [];
  const isLoading = rolesLoading || assignmentsLoading;

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
            <RoleSyncBar onSynced={() => void refetchRoles()} />
            <AssignRoleForm
              roles={roles}
              onAssigned={() => void refetchAssignments()}
            />
            <RoleAssignmentsTable
              assignments={assignments}
              roles={roles}
              onRevoked={() => void refetchAssignments()}
            />
          </>
        )}
      </Stack>
    </MasterLayout>
  );
}
