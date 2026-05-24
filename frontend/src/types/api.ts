// Supplements tygo-generated types (models.gen.ts, httputil.gen.ts).
// gorm.Model has no json tags, so its fields serialize as PascalCase:
// ID, CreatedAt, UpdatedAt, DeletedAt.

export type { Meta } from "./httputil.gen";
export {
  TypeAlbum,
  TypeTrack,
  SpecialBestNewAlbum,
  SpecialBestNewTrack,
  SpecialBestNewReissue,
  SourceBandcamp,
  SourceNPR,
  SourcePitchfork,
  SourceStereogum,
  SourceXXL,
} from "./models.gen";
export type { ReleaseType, SpecialLabel, SourceName } from "./models.gen";

import type { Meta } from "./httputil.gen";
import type { ReleaseType, SpecialLabel, SourceName } from "./models.gen";

// gorm.Model embedded fields — no json tags → PascalCase keys in JSON
interface GormModel {
  ID: number;
  CreatedAt: string;
  UpdatedAt: string;
  DeletedAt: string | null;
}

export interface Artist extends GormModel {
  name: string;
  releases: Release[];
}

export interface Release extends GormModel {
  name: string;
  artists: Artist[];
  publish_date: string;
  link: string;
  source: SourceName;
  release_type: ReleaseType;
  special?: SpecialLabel;
}

export interface ListReleasesResponse {
  data: Release[];
  meta: Meta;
}

export interface ListArtistsResponse {
  data: Artist[];
  meta: Meta;
}

export interface ListParams {
  source?: SourceName;
  type?: ReleaseType;
  special?: SpecialLabel;
  artist?: string;
  from?: string;
  to?: string;
  limit?: number;
  offset?: number;
  order_dir?: "asc" | "desc";
}

export interface ListArtistsParams {
  limit?: number;
  offset?: number;
}
