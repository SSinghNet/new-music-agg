import { useQuery } from "@tanstack/react-query";
import { getArtist } from "../api/artists";

export function useArtist(id: number) {
  return useQuery({
    queryKey: ["artist", id],
    queryFn: () => getArtist(id),
    enabled: !isNaN(id) && id > 0,
  });
}
