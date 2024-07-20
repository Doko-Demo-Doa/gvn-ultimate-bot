import { useMutation, useQuery } from "@tanstack/react-query";
import { IBackendModuleType, BackendResponseType } from "~/types/types";

const BASE_URL = process.env.NEXT_PUBLIC_BASE_API_URL;

type Params = {
  id: number;
  activated: boolean;
};

export function useAppModules() {
  return useQuery({
    queryKey: ["module-list"],
    queryFn: async () => {
      const resp = await fetch(BASE_URL + "/module/list");
      const data: BackendResponseType<IBackendModuleType[]> = await resp.json();
      return data;
    },
  });
}

export function useAppModuleEnabler() {
  return useMutation({
    mutationFn: async (params: Params) => {
      const resp = await fetch(BASE_URL + "/module/enable", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(params),
      });

      const data: BackendResponseType<IBackendModuleType[]> = await resp.json();

      return data;
    },
  });
}
