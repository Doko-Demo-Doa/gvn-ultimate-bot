import { useMutation, useQuery } from "@tanstack/react-query";
import { customApiClient } from "~/datasource/rest/api-client";
import type { IModuleConfigInput } from "~/types/payload";
import type {
  BackendResponseType,
  IAuditLogListResponse,
  IBackendModuleType,
  IDiscordChannel,
  IDiscordEmoji,
  IDiscordMember,
  IDiscordRole,
  IDiscordUserRoleAssignment,
  ISystemEventLog,
  IUserSyncResult,
} from "~/types/types";

export const ModuleActivationStatus = {
  ENABLED: 1,
  DISABLED: 0,
};

type Params = {
  module_id: number;
  is_activated: number; // 0 = disabled, 1 = enabled
};

export function useAppModules() {
  return useQuery({
    queryKey: ["module-list"],
    queryFn: async () => {
      const resp: BackendResponseType<IBackendModuleType[]> =
        await customApiClient.get("/module/list", {});
      return resp;
    },
  });
}

export function useAppModule(id: number) {
  return useQuery({
    queryKey: ["module-dautovn", id],
    queryFn: async () => {
      const resp: BackendResponseType<IBackendModuleType> =
        await customApiClient.get(`/module/id/${id}`, {});
      return resp;
    },
  });
}

export function useAppModuleEnabler() {
  return useMutation({
    mutationFn: async (params: Params) => {
      const resp = await customApiClient.post("/module/on-off", params);
      return resp;
    },
  });
}

export function useModuleConfigMutation() {
  return useMutation({
    mutationFn: async (params: IModuleConfigInput) => {
      const payload = {
        module_id: params.module_id,
        new_config: JSON.stringify(params.new_config),
      };
      const resp = await customApiClient.post("/module/update-config", payload);
      return resp;
    },
  });
}

// ################# Discord Role APIs #################

export function useDiscordRoles() {
  return useQuery({
    queryKey: ["discord-roles"],
    queryFn: async () => {
      const resp: BackendResponseType<IDiscordRole[]> =
        await customApiClient.get("/discord/role/list", {});
      return resp;
    },
  });
}

export function useDiscordChannels() {
  return useQuery({
    queryKey: ["discord-channels"],
    queryFn: async () => {
      const resp: BackendResponseType<IDiscordChannel[]> =
        await customApiClient.get("/discord/channels", {});
      return resp;
    },
  });
}

export function useDiscordEmojis() {
  return useQuery({
    queryKey: ["discord-emojis"],
    queryFn: async () => {
      const resp: BackendResponseType<IDiscordEmoji[]> =
        await customApiClient.get("/discord/emojis", {});
      return resp;
    },
  });
}

export function useSearchDiscordMembers(query: string) {
  return useQuery({
    queryKey: ["discord-members-search", query],
    queryFn: async () => {
      const resp: BackendResponseType<IDiscordMember[]> =
        await customApiClient.get("/discord/members/search", { q: query });
      return resp;
    },
    enabled: query.length >= 2,
  });
}

export function useRoleAssignments() {
  return useQuery({
    queryKey: ["discord-role-assignments"],
    queryFn: async () => {
      const resp: BackendResponseType<IDiscordUserRoleAssignment[]> =
        await customApiClient.get("/discord/role/assignments", {});
      return resp;
    },
  });
}

export function useAssignRoleMutation() {
  return useMutation({
    mutationFn: async (params: {
      user_native_id: string;
      role_native_id: string;
      duration: string;
    }) => {
      const resp = await customApiClient.post("/discord/role/assign", params);
      return resp;
    },
  });
}

export function useRevokeRoleMutation() {
  return useMutation({
    mutationFn: async (id: number) => {
      const resp = await customApiClient.delete(`/discord/role/assign/${id}`);
      return resp;
    },
  });
}

export function useSyncDiscordUsers() {
  return useMutation({
    mutationFn: async () => {
      const resp: BackendResponseType<IUserSyncResult> =
        await customApiClient.post("/discord/users/sync", {});
      return resp;
    },
  });
}

export function useLastUserSync() {
  return useQuery({
    queryKey: ["discord-users-last-sync"],
    queryFn: async () => {
      const resp: BackendResponseType<ISystemEventLog> =
        await customApiClient.get("/discord/users/sync/last", {});
      return resp;
    },
  });
}

// ################# Message Audit Log APIs #################

export function useAuditLogs(params: {
  limit?: number;
  offset?: number;
  from_date?: string;
  to_date?: string;
  channel_id?: string;
  author_name?: string;
}) {
  const queryKey = ["audit-logs", params];
  return useQuery({
    queryKey,
    queryFn: async () => {
      const resp: BackendResponseType<IAuditLogListResponse> =
        await customApiClient.get("/discord/audit-logs", params);
      return resp;
    },
  });
}

export function useClearAuditLogs() {
  return useMutation({
    mutationFn: async () => {
      const resp = await customApiClient.delete("/discord/audit-logs");
      return resp;
    },
  });
}
