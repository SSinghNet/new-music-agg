interface PaginationProps {
  total: number;
  limit: number;
  offset: number;
  onChange: (offset: number) => void;
}

export function Pagination({ total, limit, offset, onChange }: PaginationProps) {
  const totalPages = Math.ceil(total / limit);
  const currentPage = Math.floor(offset / limit) + 1;

  if (totalPages <= 1) return null;

  const pages = buildPageList(currentPage, totalPages);

  return (
    <nav className="flex items-center justify-center gap-1 py-8" aria-label="Pagination">
      <button
        onClick={() => onChange((currentPage - 2) * limit)}
        disabled={currentPage === 1}
        className="flex h-9 w-9 items-center justify-center rounded-lg text-sm text-zinc-400 transition hover:bg-zinc-800 disabled:pointer-events-none disabled:opacity-30"
        aria-label="Previous page"
      >
        ‹
      </button>

      {pages.map((page, i) =>
        page === "..." ? (
          <span key={`ellipsis-${i}`} className="flex h-9 w-9 items-center justify-center text-sm text-zinc-500">
            …
          </span>
        ) : (
          <button
            key={page}
            onClick={() => onChange((page - 1) * limit)}
            className={`flex h-9 w-9 items-center justify-center rounded-lg text-sm transition ${
              page === currentPage
                ? "bg-zinc-700 font-semibold text-white"
                : "text-zinc-400 hover:bg-zinc-800"
            }`}
          >
            {page}
          </button>
        )
      )}

      <button
        onClick={() => onChange(currentPage * limit)}
        disabled={currentPage === totalPages}
        className="flex h-9 w-9 items-center justify-center rounded-lg text-sm text-zinc-400 transition hover:bg-zinc-800 disabled:pointer-events-none disabled:opacity-30"
        aria-label="Next page"
      >
        ›
      </button>
    </nav>
  );
}

function buildPageList(current: number, total: number): (number | "...")[] {
  if (total <= 7) return Array.from({ length: total }, (_, i) => i + 1);

  if (current <= 4) {
    return [1, 2, 3, 4, 5, "...", total];
  }
  if (current >= total - 3) {
    return [1, "...", total - 4, total - 3, total - 2, total - 1, total];
  }
  return [1, "...", current - 1, current, current + 1, "...", total];
}
