import { Link, useParams } from "react-router-dom";
import { useRelease } from "../hooks/useRelease";
import { SourceBadge, TypeBadge, SpecialBadge } from "../components/ui/Badge";
import ExternalLinkIcon from "../assets/icons/ExternalLinkIcon";

function formatDate(iso: string) {
  return new Date(iso).toLocaleDateString("en-US", {
    weekday: "long",
    month: "long",
    day: "numeric",
    year: "numeric",
  });
}

export default function ReleasePage() {
  const { id } = useParams<{ id: string }>();
  const { data: release, isPending, isError } = useRelease(Number(id));

  if (isPending) {
    return (
      <div className="space-y-4 py-8">
        <div className="h-8 w-48 animate-pulse rounded-lg bg-zinc-800" />
        <div className="h-12 w-96 animate-pulse rounded-lg bg-zinc-800" />
        <div className="h-4 w-64 animate-pulse rounded-lg bg-zinc-800" />
      </div>
    );
  }

  if (isError || !release) {
    return (
      <div className="py-20 text-center text-sm text-red-400">
        Release not found or failed to load.
      </div>
    );
  }

  return (
    <div className="mx-auto max-w-2xl py-8">
      {release.special && (
        <div className="mb-4">
          <SpecialBadge label={release.special} />
        </div>
      )}

      {release.artists?.length > 0 && (
        <p className="mb-1 flex flex-wrap gap-x-1.5 text-sm font-medium text-zinc-400">
          {release.artists.map((artist, i) => (
            <span key={artist.ID}>
              <Link
                to={`/artists/${artist.ID}`}
                className="hover:text-zinc-100 hover:underline underline-offset-2"
              >
                {artist.name}
              </Link>
              {i < release.artists.length - 1 && <span className="text-zinc-600">,</span>}
            </span>
          ))}
        </p>
      )}

      <h1 className="mb-4 text-3xl font-bold leading-tight tracking-tight text-zinc-100 sm:text-4xl">
        {release.name}
      </h1>

      <div className="mb-8 flex flex-wrap items-center gap-2">
        <SourceBadge source={release.source} />
        <TypeBadge type={release.release_type} />
        <span className="text-xs text-zinc-500">{formatDate(release.publish_date)}</span>
      </div>

      <div className="rounded-xl border border-zinc-800 bg-zinc-900 p-6">
        <dl className="space-y-4">
          {release.artists?.length > 0 && (
            <div className="flex gap-4">
              <dt className="w-24 shrink-0 text-sm text-zinc-500">Artists</dt>
              <dd className="flex flex-wrap gap-x-1.5 text-sm text-zinc-100">
                {release.artists.map((artist, i) => (
                  <span key={artist.ID}>
                    <Link
                      to={`/artists/${artist.ID}`}
                      className="hover:text-zinc-400 hover:underline underline-offset-2"
                    >
                      {artist.name}
                    </Link>
                    {i < release.artists.length - 1 && <span className="text-zinc-500">,</span>}
                  </span>
                ))}
              </dd>
            </div>
          )}
          <Row label="Release" value={release.name} />
          <Row label="Type" value={release.release_type} />
          <Row label="Source" value={release.source} />
          <Row label="Published" value={formatDate(release.publish_date)} />
          {release.special && <Row label="Award" value={release.special} />}
        </dl>

        {release.link && (
          <div className="mt-6 border-t border-zinc-800 pt-6">
            <a
              href={release.link}
              target="_blank"
              rel="noopener noreferrer"
              className="inline-flex items-center gap-2 rounded-lg bg-zinc-800 px-4 py-2.5 text-sm font-medium text-zinc-100 transition hover:bg-zinc-700"
            >
              Read on {release.source}
              <ExternalLinkIcon />
            </a>
          </div>
        )}
      </div>
    </div>
  );
}

function Row({ label, value }: { label: string; value: string }) {
  return (
    <div className="flex gap-4">
      <dt className="w-24 shrink-0 text-sm text-zinc-500">{label}</dt>
      <dd className="text-sm text-zinc-100">{value}</dd>
    </div>
  );
}
