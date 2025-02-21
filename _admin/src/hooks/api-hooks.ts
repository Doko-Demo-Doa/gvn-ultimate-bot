import { useMutation, useQuery } from "@tanstack/react-query";
import { customApiClient } from "~/datasource/rest/api-client";
import type { IModuleConfigInput } from "~/types/payload";
import type { BackendResponseType, IBackendModuleType } from "~/types/types";

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
