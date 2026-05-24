import { fetchJSON } from "./client";
import type { ListReleasesResponse, Release, ListParams } from "../types/api";

export function listReleases(params: ListParams): Promise<ListReleasesResponse> {
  return fetchJSON<ListReleasesResponse>("/releases", params as Record<string, string | number | undefined>);
}

export function getRelease(id: number): Promise<Release> {
  return fetchJSON<Release>(`/releases/${id}`);
}
