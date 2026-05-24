import { useQuery } from "@tanstack/react-query";
import { listReleases } from "../api/releases";
import type { ListParams } from "../types/api";

export const PAGE_SIZE = 24;

export function useReleases(params: ListParams) {
  return useQuery({
    queryKey: ["releases", params],
    queryFn: () => listReleases(params),
    placeholderData: (prev) => prev,
  });
}
