import { useSearchParams, Link } from "react-router-dom";
import { useArtists, ARTISTS_PAGE_SIZE } from "../hooks/useArtists";
import { Pagination } from "../components/ui/Pagination";
import MusicNoteIcon from "../assets/icons/MusicNoteIcon";

export default function ArtistsPage() {
  const [searchParams, setSearchParams] = useSearchParams();
  const offset = Number(searchParams.get("offset") ?? 0);

  const { data, isPending, isError } = useArtists({ limit: ARTISTS_PAGE_SIZE, offset });

  function handlePageChange(newOffset: number) {
    setSearchParams((prev) => {
      const next = new URLSearchParams(prev);
      if (newOffset === 0) {
        next.delete("offset");
      } else {
        next.set("offset", String(newOffset));
      }
      return next;
    });
    window.scrollTo({ top: 0, behavior: "smooth" });
  }

  return (
    <div>
      <div className="mb-6">
        <h1 className="text-2xl font-bold tracking-tight text-zinc-100">Artists</h1>
        {data && (
          <p className="mt-1 text-sm text-zinc-500">
            {data.meta.total.toLocaleString()} {data.meta.total === 1 ? "artist" : "artists"} in the directory
          </p>
        )}
      </div>

      {isError && (
        <div className="rounded-xl border border-red-900/50 bg-red-950/30 p-6 text-center text-sm text-red-400">
          Failed to load artists. Make sure the backend is running.
        </div>
      )}

      {isPending ? (
        <div className="grid grid-cols-2 gap-3 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6">
          {Array.from({ length: 24 }).map((_, i) => (
            <div key={i} className="h-20 animate-pulse rounded-xl bg-zinc-800/50" />
          ))}
        </div>
      ) : data?.data.length === 0 ? (
        <div className="py-20 text-center text-sm text-zinc-500">No artists found.</div>
      ) : (
        <>
          <div className="grid grid-cols-2 gap-3 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6">
            {data?.data.map((artist) => {
              const releaseCount = artist.releases?.length ?? 0;
              return (
                <Link
                  key={artist.ID}
                  to={`/artists/${artist.ID}`}
                  className="flex flex-col gap-1.5 rounded-xl border border-zinc-800 bg-zinc-900 p-4 transition hover:border-zinc-600 hover:bg-zinc-800/70"
                >
                  <div className="flex h-8 w-8 items-center justify-center rounded-full bg-zinc-800">
                    <MusicNoteIcon className="size-4 text-zinc-400" />
                  </div>
                  <span className="mt-1 text-sm font-medium leading-tight text-zinc-100">
                    {artist.name}
                  </span>
                  {releaseCount > 0 && (
                    <span className="text-xs text-zinc-500">
                      {releaseCount} {releaseCount === 1 ? "release" : "releases"}
                    </span>
                  )}
                </Link>
              );
            })}
          </div>

          {data && (
            <Pagination
              total={data.meta.total}
              limit={data.meta.limit}
              offset={offset}
              onChange={handlePageChange}
            />
          )}
        </>
      )}
    </div>
  );
}
