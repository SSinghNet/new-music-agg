import { fetchJSON } from "./client";
import type { ListArtistsResponse, Artist, ListArtistsParams } from "../types/api";

export function listArtists(params: ListArtistsParams): Promise<ListArtistsResponse> {
  return fetchJSON<ListArtistsResponse>("/artists", params as Record<string, string | number | undefined>);
}

export function getArtist(id: number): Promise<Artist> {
  return fetchJSON<Artist>(`/artists/${id}`);
}
