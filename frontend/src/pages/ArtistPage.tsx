import { useParams, useSearchParams } from "react-router-dom";
import { useArtist } from "../hooks/useArtist";
import { useReleases, PAGE_SIZE } from "../hooks/useReleases";
import ReleaseCard from "../components/ReleaseCard";
import { Pagination } from "../components/ui/Pagination";
import {
  TypeAlbum,
  TypeTrack,
  SpecialBestNewAlbum,
  SpecialBestNewTrack,
  SpecialBestNewReissue,
} from "../types/api";
import type { ReleaseType, SpecialLabel } from "../types/api";
import MusicNoteIcon from "../assets/icons/MusicNoteIcon";

const TYPES: ReleaseType[] = [TypeAlbum, TypeTrack];
const SPECIALS: SpecialLabel[] = [SpecialBestNewAlbum, SpecialBestNewTrack, SpecialBestNewReissue];

function Chip({
  label,
  active,
  onClick,
}: {
  label: string;
  active: boolean;
  onClick: () => void;
}) {
  return (
    <button
      onClick={onClick}
      className={`rounded-full px-2.5 py-1 text-xs transition ${
        active ? "bg-zinc-700 text-white" : "text-zinc-400 hover:bg-zinc-800 hover:text-zinc-200"
      }`}
    >
      {label}
    </button>
  );
}

export default function ArtistPage() {
  const { id } = useParams<{ id: string }>();
  const [searchParams, setSearchParams] = useSearchParams();

  const { data: artist, isPending: artistPending, isError: artistError } = useArtist(Number(id));

  const type = (searchParams.get("type") as ReleaseType) || undefined;
  const special = (searchParams.get("special") as SpecialLabel) || undefined;
  const offset = Number(searchParams.get("offset") ?? 0);

  const { data: releasesData, isPending: releasesPending } = useReleases({
    artist: artist?.name,
    type,
    special,
    limit: PAGE_SIZE,
    offset,
  });

  function setFilter(key: string, value: string | undefined) {
    setParams((prev) => {
      const next = new URLSearchParams(prev);
      if (value) next.set(key, value);
      else next.delete(key);
      next.delete("offset");
      return next;
    });
  }

  function setParams(updater: (prev: URLSearchParams) => URLSearchParams) {
    setSearchParams(updater);
  }

  function handlePageChange(newOffset: number) {
    setSearchParams((prev) => {
      const next = new URLSearchParams(prev);
      if (newOffset === 0) next.delete("offset");
      else next.set("offset", String(newOffset));
      return next;
    });
    window.scrollTo({ top: 0, behavior: "smooth" });
  }

  if (artistError) {
    return (
      <div className="py-20 text-center text-sm text-red-400">
        Artist not found or failed to load.
      </div>
    );
  }

  return (
    <div>
      <div className="mb-6 flex items-center gap-4">
        <div className="flex h-14 w-14 shrink-0 items-center justify-center rounded-full bg-zinc-800">
          <MusicNoteIcon className="size-7 text-zinc-400" />
        </div>
        <div>
          {artistPending ? (
            <div className="h-8 w-48 animate-pulse rounded-lg bg-zinc-800" />
          ) : (
            <h1 className="text-2xl font-bold tracking-tight text-zinc-100">{artist?.name}</h1>
          )}
          {releasesData && (
            <p className="mt-0.5 text-sm text-zinc-500">
              {releasesData.meta.total.toLocaleString()} {releasesData.meta.total === 1 ? "release" : "releases"}
            </p>
          )}
        </div>
      </div>

      <div className="mb-6 flex flex-wrap items-center gap-2 rounded-xl border border-zinc-800 bg-zinc-900/50 p-3">
        <span className="text-xs text-zinc-500">Filter</span>
        <div className="h-3 w-px bg-zinc-700" />
        <Chip label="All types" active={!type} onClick={() => setFilter("type", undefined)} />
        {TYPES.map((t) => (
          <Chip key={t} label={t} active={type === t} onClick={() => setFilter("type", type === t ? undefined : t)} />
        ))}
        <div className="h-3 w-px bg-zinc-700" />
        <Chip label="All" active={!special} onClick={() => setFilter("special", undefined)} />
        {SPECIALS.map((s) => (
          <Chip
            key={s}
            label={s}
            active={special === s}
            onClick={() => setFilter("special", special === s ? undefined : s)}
          />
        ))}
      </div>

      {releasesPending || artistPending ? (
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
          {Array.from({ length: PAGE_SIZE }).map((_, i) => (
            <div key={i} className="h-44 animate-pulse rounded-xl bg-zinc-800/50" />
          ))}
        </div>
      ) : releasesData?.data.length === 0 ? (
        <div className="py-20 text-center text-sm text-zinc-500">No releases found.</div>
      ) : (
        <>
          <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
            {releasesData?.data.map((release) => (
              <ReleaseCard key={release.ID} release={release} />
            ))}
          </div>
          {releasesData && (
            <Pagination
              total={releasesData.meta.total}
              limit={releasesData.meta.limit}
              offset={offset}
              onChange={handlePageChange}
            />
          )}
        </>
      )}
    </div>
  );
}
