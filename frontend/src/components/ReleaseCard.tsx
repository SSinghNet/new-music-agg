import { Link, useNavigate } from "react-router-dom";
import type { Release } from "../types/api";
import { SourceBadge, TypeBadge, SpecialBadge } from "./ui/Badge";
import ExternalLinkIcon from "../assets/icons/ExternalLinkIcon";

interface ReleaseCardProps {
  release: Release;
}

function formatDate(iso: string) {
  return new Date(iso).toLocaleDateString("en-US", {
    month: "short",
    day: "numeric",
    year: "numeric",
  });
}

export default function ReleaseCard({ release }: ReleaseCardProps) {
  const navigate = useNavigate();

  const handleCardClick = () => navigate(`/releases/${release.ID}`);

  const handleLinkClick = (e: React.MouseEvent) => {
    e.stopPropagation();
    if (release.link) window.open(release.link, "_blank", "noopener,noreferrer");
  };

  return (
    <article
      onClick={handleCardClick}
      className="group relative flex cursor-pointer flex-col gap-3 rounded-xl border border-zinc-800 bg-zinc-900 p-5 transition hover:border-zinc-600 hover:bg-zinc-800/70"
    >
      {release.special && (
        <div>
          <SpecialBadge label={release.special} />
        </div>
      )}

      <div className="flex-1">
        {release.artists?.length > 0 && (
          <p className="mb-1 flex flex-wrap gap-x-1 text-xs text-zinc-500">
            {release.artists.map((artist, i) => (
              <span key={artist.ID}>
                <Link
                  to={`/artists/${artist.ID}`}
                  onClick={(e) => e.stopPropagation()}
                  className="hover:text-zinc-300 hover:underline underline-offset-2"
                >
                  {artist.name}
                </Link>
                {i < release.artists.length - 1 && <span className="select-none">,</span>}
              </span>
            ))}
          </p>
        )}
        <h2 className="text-base font-semibold leading-snug text-zinc-100 group-hover:text-white">
          {release.name}
        </h2>
      </div>

      <div className="flex items-end justify-between gap-2">
        <div className="flex flex-col gap-2">
          <div className="flex flex-wrap gap-1.5">
            <SourceBadge source={release.source} />
            <TypeBadge type={release.release_type} />
          </div>
          <span className="text-xs text-zinc-500">{formatDate(release.publish_date)}</span>
        </div>

        {release.link && (
          <button
            onClick={handleLinkClick}
            title="Open original article"
            className="shrink-0 rounded-lg p-1.5 text-zinc-500 transition hover:bg-zinc-700 hover:text-zinc-200"
          >
            <ExternalLinkIcon />
          </button>
        )}
      </div>
    </article>
  );
}
