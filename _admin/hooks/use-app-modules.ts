import { useMutation, useQuery } from "@tanstack/react-query";
import { BackendModuleType, BackendResponseType } from "~/types/types";

const BASE_URL = process.env.NEXT_PUBLIC_BASE_API_URL;

type Params = {
  id: number;
  activated: boolean;
};

export function useAppModules() {
  return useQuery([], async () => {
    const resp = await fetch(BASE_URL + "/module/list");
    const data: BackendResponseType<BackendModuleType[]> = await resp.json();

    return data;
  });
}

// TODO: Refactor
export function useAppModuleEnabler() {
  return useMutation([], async (params: Params) => {
    const resp = await fetch(BASE_URL + `/module/on-off`, {
      method: "POST",
      body: JSON.stringify({
        module_id: params.id,
        is_activated: params.activated ? 1 : 0,
      }),
    });
    const data: any = await resp.json();
    return data;
  });
}
