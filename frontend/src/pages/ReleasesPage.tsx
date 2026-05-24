import { useSearchParams } from "react-router-dom";
import { useReleases, PAGE_SIZE } from "../hooks/useReleases";
import ReleaseCard from "../components/ReleaseCard";
import FilterBar from "../components/FilterBar";
import { Pagination } from "../components/ui/Pagination";
import type { ListParams } from "../types/api";

export default function ReleasesPage() {
  const [searchParams, setSearchParams] = useSearchParams();

  const params: ListParams = {
    source: searchParams.get("source") || undefined,
    type: searchParams.get("type") || undefined,
    special: searchParams.get("special") || undefined,
    offset: Number(searchParams.get("offset") ?? 0),
    limit: PAGE_SIZE,
  } as ListParams;

  const { data, isPending, isError } = useReleases(params);

  const offset = Number(searchParams.get("offset") ?? 0);

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
        <h1 className="text-2xl font-bold tracking-tight text-zinc-100">New Music</h1>
        <p className="mt-1 text-sm text-zinc-500">
          Latest releases from Pitchfork, NPR, Bandcamp, Stereogum, and XXL
        </p>
      </div>

      <FilterBar total={data?.meta.total} />

      {isError && (
        <div className="rounded-xl border border-red-900/50 bg-red-950/30 p-6 text-center text-sm text-red-400">
          Failed to load releases. Make sure the backend is running.
        </div>
      )}

      {isPending ? (
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
          {Array.from({ length: PAGE_SIZE }).map((_, i) => (
            <div key={i} className="h-44 animate-pulse rounded-xl bg-zinc-800/50" />
          ))}
        </div>
      ) : data?.data.length === 0 ? (
        <div className="py-20 text-center text-sm text-zinc-500">No releases match your filters.</div>
      ) : (
        <>
          <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
            {data?.data.map((release) => (
              <ReleaseCard key={release.ID} release={release} />
            ))}
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
