import { useQuery } from "@tanstack/react-query";
import { getRelease } from "../api/releases";

export function useRelease(id: number) {
  return useQuery({
    queryKey: ["release", id],
    queryFn: () => getRelease(id),
    enabled: !isNaN(id) && id > 0,
  });
}
