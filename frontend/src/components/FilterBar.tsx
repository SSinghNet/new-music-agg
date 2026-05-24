import { useSearchParams } from "react-router-dom";
import {
  SourceBandcamp,
  SourceNPR,
  SourcePitchfork,
  SourceStereogum,
  SourceXXL,
  TypeAlbum,
  TypeTrack,
  SpecialBestNewAlbum,
  SpecialBestNewTrack,
  SpecialBestNewReissue,
} from "../types/api";
import type { SourceName, ReleaseType, SpecialLabel } from "../types/api";

const SOURCES: SourceName[] = [SourcePitchfork, SourceNPR, SourceBandcamp, SourceStereogum, SourceXXL];
const TYPES: ReleaseType[] = [TypeAlbum, TypeTrack];
const SPECIALS: SpecialLabel[] = [SpecialBestNewAlbum, SpecialBestNewTrack, SpecialBestNewReissue];

interface ToggleGroupProps<T extends string> {
  label: string;
  options: T[];
  value: T | undefined;
  onChange: (v: T | undefined) => void;
}

function ToggleGroup<T extends string>({ label, options, value, onChange }: ToggleGroupProps<T>) {
  return (
    <div className="flex items-center gap-1.5">
      <span className="w-16 shrink-0 text-xs font-medium text-zinc-500">{label}</span>
      <div className="flex flex-wrap gap-1">
        <button
          onClick={() => onChange(undefined)}
          className={`rounded-full px-2.5 py-1 text-xs transition ${
            !value ? "bg-zinc-700 text-white" : "text-zinc-400 hover:bg-zinc-800 hover:text-zinc-200"
          }`}
        >
          All
        </button>
        {options.map((opt) => (
          <button
            key={opt}
            onClick={() => onChange(value === opt ? undefined : opt)}
            className={`rounded-full px-2.5 py-1 text-xs transition ${
              value === opt ? "bg-zinc-700 text-white" : "text-zinc-400 hover:bg-zinc-800 hover:text-zinc-200"
            }`}
          >
            {opt}
          </button>
        ))}
      </div>
    </div>
  );
}

interface FilterBarProps {
  total?: number;
}

export default function FilterBar({ total }: FilterBarProps) {
  const [params, setParams] = useSearchParams();

  const source = (params.get("source") as SourceName) || undefined;
  const type = (params.get("type") as ReleaseType) || undefined;
  const special = (params.get("special") as SpecialLabel) || undefined;

  function setFilter(key: string, value: string | undefined) {
    setParams((prev) => {
      const next = new URLSearchParams(prev);
      if (value) {
        next.set(key, value);
      } else {
        next.delete(key);
      }
      next.delete("offset");
      return next;
    });
  }

  const hasFilters = !!(source || type || special);

  return (
    <div className="mb-6 flex flex-col gap-3 rounded-xl border border-zinc-800 bg-zinc-900/50 p-4">
      <ToggleGroup<SourceName>
        label="Source"
        options={SOURCES}
        value={source}
        onChange={(v) => setFilter("source", v)}
      />
      <div className="border-t border-zinc-800/60" />
      <ToggleGroup<ReleaseType>
        label="Type"
        options={TYPES}
        value={type}
        onChange={(v) => setFilter("type", v)}
      />
      <div className="border-t border-zinc-800/60" />
      <ToggleGroup<SpecialLabel>
        label="Special"
        options={SPECIALS}
        value={special}
        onChange={(v) => setFilter("special", v)}
      />

      <div className="flex items-center justify-between border-t border-zinc-800/60 pt-2">
        {total !== undefined && (
          <span className="text-xs text-zinc-500">
            {total.toLocaleString()} {total === 1 ? "release" : "releases"}
          </span>
        )}
        {hasFilters && (
          <button
            onClick={() => setParams({})}
            className="ml-auto text-xs text-zinc-500 underline-offset-2 hover:text-zinc-300 hover:underline"
          >
            Clear filters
          </button>
        )}
      </div>
    </div>
  );
}
