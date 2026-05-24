import { useQuery } from "@tanstack/react-query";
import { listArtists } from "../api/artists";
import type { ListArtistsParams } from "../types/api";

export const ARTISTS_PAGE_SIZE = 50;

export function useArtists(params: ListArtistsParams) {
  return useQuery({
    queryKey: ["artists", params],
    queryFn: () => listArtists(params),
    placeholderData: (prev) => prev,
  });
}
