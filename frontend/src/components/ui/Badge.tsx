import type { SourceName, ReleaseType, SpecialLabel } from "../../types/api";
import StarIcon from "../../assets/icons/StarIcon";

const SOURCE_STYLES: Record<string, string> = {
  Pitchfork: "bg-emerald-500/15 text-emerald-400 ring-1 ring-emerald-500/30",
  NPR: "bg-red-500/15 text-red-400 ring-1 ring-red-500/30",
  Bandcamp: "bg-cyan-500/15 text-cyan-400 ring-1 ring-cyan-500/30",
  Stereogum: "bg-violet-500/15 text-violet-400 ring-1 ring-violet-500/30",
  XXL: "bg-amber-500/15 text-amber-400 ring-1 ring-amber-500/30",
};

interface SourceBadgeProps {
  source: SourceName;
}

export function SourceBadge({ source }: SourceBadgeProps) {
  const cls = SOURCE_STYLES[source] ?? "bg-zinc-700/50 text-zinc-300 ring-1 ring-zinc-600/50";
  return (
    <span className={`inline-flex items-center rounded-full px-2 py-0.5 text-xs font-medium ${cls}`}>
      {source}
    </span>
  );
}

interface TypeBadgeProps {
  type: ReleaseType;
}

export function TypeBadge({ type }: TypeBadgeProps) {
  return (
    <span className="inline-flex items-center rounded-full bg-zinc-700/50 px-2 py-0.5 text-xs font-medium text-zinc-300 ring-1 ring-zinc-600/50">
      {type}
    </span>
  );
}

interface SpecialBadgeProps {
  label: SpecialLabel;
}

export function SpecialBadge({ label }: SpecialBadgeProps) {
  return (
    <span className="inline-flex items-center gap-1 rounded-full bg-amber-400/15 px-2.5 py-0.5 text-xs font-semibold text-amber-300 ring-1 ring-amber-400/30">
      <StarIcon className="size-3" />
      {label}
    </span>
  );
}
