import { useMutation, useQuery } from "@tanstack/react-query";
import { customApiClient } from "~/datasource/rest/api-client";
import type {
  BackendResponseType,
  IDiscordRoleReactionEmbed,
  IReactionRoleMessagePayload,
} from "~/types/types";

export enum ReactionType {
  Emoji = 0,
  Button = 1,
  Dropdown = 2,
}

export enum ReactionMode {
  Default = "default",
  Reverse = "reverse",
}

export function useReactionRoleEmbeds() {
  return useQuery({
    queryKey: ["reaction-role-embeds"],
    queryFn: async () => {
      const resp: BackendResponseType<IDiscordRoleReactionEmbed[]> =
        await customApiClient.get("/discord/role-reaction/list", {});
      return resp;
    },
  });
}

export function useUpsertReactionRoleEmbed() {
  return useMutation({
    mutationFn: async (params: {
      native_message_id: string;
      name: string;
      payload: IReactionRoleMessagePayload;
      mode?: string;
      tags?: string;
      version?: number;
    }) => {
      const resp: BackendResponseType<IDiscordRoleReactionEmbed> =
        await customApiClient.post("/discord/role-reaction/upsert", params);
      return resp;
    },
  });
}

export function usePublishReactionRoleEmbed() {
  return useMutation({
    mutationFn: async (payload: IReactionRoleMessagePayload) => {
      const resp: BackendResponseType<IDiscordRoleReactionEmbed> =
        await customApiClient.post("/discord/role-reaction/publish", {
          payload,
        });
      return resp;
    },
  });
}

export function useDeleteReactionRoleEmbed() {
  return useMutation({
    mutationFn: async (id: number) => {
      const resp = await customApiClient.delete(`/discord/role-reaction/${id}`);
      return resp;
    },
  });
}
