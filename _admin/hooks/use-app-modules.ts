import { useMutation, useQuery } from "@tanstack/react-query";
import { BackendModuleType, BackendResponseType } from "~/types/types";

const BASE_URL = process.env.NEXT_PUBLIC_BACKEND_BASE_URL;

export function useAppModules() {
  return useQuery([], async () => {
    const resp = await fetch(BASE_URL + "/module/list");
    const data: BackendResponseType<BackendModuleType[]> = await resp.json();
    return data;
  });
}

// TODO: Refactor
export function useAppModuleEnabler() {
  return useMutation([], async (id: number) => {
    const resp = await fetch(BASE_URL + `/module/${id}`);
    const data: any = await resp.json();
    return data;
  })
}