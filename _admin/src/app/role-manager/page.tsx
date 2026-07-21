"use client";

import {
	Button,
	Group,
	Loader,
	NativeSelect,
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
	useRoleAssignments,
} from "~/hooks/api-hooks";
import MasterLayout from "~/layouts/master-layout";

export default function RoleManagerPage() {
	const { data: rolesData, isLoading: rolesLoading } = useDiscordRoles();
	const { data: assignmentsData, isLoading: assignmentsLoading } =
		useRoleAssignments();
	const { mutateAsync: assignRole, isPending: isAssigning } =
		useAssignRoleMutation();

	const [userId, setUserId] = useState("");
	const [roleNativeId, setRoleNativeId] = useState("");
	const [durationValue, setDurationValue] = useState("1");
	const [durationUnit, setDurationUnit] = useState("d");

	const roles = rolesData?.data ?? [];
	const assignments = assignmentsData?.data ?? [];

	const isLoading = rolesLoading || assignmentsLoading;

	async function handleAssign() {
		if (!userId || !roleNativeId) {
			notifications.show({
				color: "red",
				title: "Lỗi",
				message: "Vui lòng nhập User ID và chọn Role",
			});
			return;
		}

		// Convert duration to Go-parseable format
		const duration = `${durationValue}${durationUnit}`;

		try {
			await assignRole({
				user_native_id: userId,
				role_native_id: roleNativeId,
				duration,
			});

			notifications.show({
				color: "green",
				title: "Thành công",
				message: `Đã gán role ${getRoleName(roleNativeId)} cho user ${userId} (${duration})`,
			});
		} catch (err: any) {
			notifications.show({
				color: "red",
				title: "Lỗi",
				message: err?.message || "Không thể gán role",
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
					Gán role cho user với thời hạn. Hệ thống sẽ tự động thu hồi role
					khi hết hạn.
				</Text>

				{isLoading ? (
					<Loader />
				) : (
					<>
						<Stack>
							<Title order={4}>Gán role mới</Title>
							<TextInput
								label="Discord User ID"
								placeholder="123456789012345678"
								value={userId}
								onChange={(e) => setUserId(e.currentTarget.value)}
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
							<Title order={4}>Danh sách đã hết hạn / sắp hết hạn</Title>
							{assignments.length === 0 ? (
								<Text c="dimmed">Không có role nào đang được gán có thời hạn.</Text>
							) : (
								<Table>
									<Table.Thead>
										<Table.Tr>
											<Table.Th>User ID</Table.Th>
											<Table.Th>Role</Table.Th>
											<Table.Th>Gán lúc</Table.Th>
											<Table.Th>Hết hạn</Table.Th>
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
