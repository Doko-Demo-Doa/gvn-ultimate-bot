import { useMutation, useQuery } from "@tanstack/react-query";
import { customApiClient } from "~/datasource/rest/api-client";
import type { IModuleConfigInput } from "~/types/payload";
import type {
  BackendResponseType,
  IBackendModuleType,
  IDiscordChannel,
  IDiscordEmoji,
  IDiscordRole,
  IDiscordUserRoleAssignment,
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
