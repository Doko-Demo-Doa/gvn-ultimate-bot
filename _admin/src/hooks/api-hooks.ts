import { useMutation, useQuery } from "@tanstack/react-query";
import { customApiClient } from "~/datasource/rest/api-client";
import { IPinModuleInput } from "~/types/payload";
import { IBackendModuleType, BackendResponseType } from "~/types/types";

export const ModuleActivationStatus = {
  ENABLED: 1,
  DISABLED: 0,
};

type Params = {
  module_id: number;
  is_activated: number; // 0 = disabled, 1 = enabled
};

customApiClient.init({
  baseUrl: process.env.NEXT_PUBLIC_BASE_API_URL || "",
});

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
        await customApiClient.get(`/module/${id}`, {});
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

export function usePinThresholdMutation() {
  return useMutation({
    mutationFn: async (params: IPinModuleInput) => {
      const resp = await customApiClient.post("/module/update-config", params);
      return resp;
    },
  });
}
