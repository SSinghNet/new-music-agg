import { Link, useLocation } from "react-router-dom";
import ArrowLeftIcon from "../assets/icons/ArrowLeftIcon";
import MusicNoteIcon from "../assets/icons/MusicNoteIcon";

interface LayoutProps {
  children: React.ReactNode;
}

export default function Layout({ children }: LayoutProps) {
  const { pathname } = useLocation();

  const backTo = pathname.startsWith("/artists/") && pathname !== "/artists"
    ? { href: "/artists", label: "All artists" }
    : pathname !== "/"
    ? { href: "/", label: "All releases" }
    : null;

  return (
    <div className="min-h-screen bg-zinc-950 text-zinc-100">
      <header className="sticky top-0 z-10 border-b border-zinc-800 bg-zinc-950/90 backdrop-blur">
        <div className="mx-auto flex max-w-7xl items-center gap-4 px-4 py-3 sm:px-6">
          <Link to="/" className="flex items-center gap-2 text-white transition hover:opacity-80">
            <MusicNoteIcon className="size-5 text-zinc-400" />
            <span className="text-lg font-bold tracking-tight">new music</span>
            <span className="rounded bg-zinc-800 px-1.5 py-0.5 text-xs font-medium text-zinc-400">agg</span>
          </Link>

          <nav className="ml-4 hidden items-center gap-1 sm:flex">
            <Link
              to="/"
              className={`rounded-lg px-3 py-1.5 text-sm transition ${
                pathname === "/" ? "text-zinc-100" : "text-zinc-500 hover:text-zinc-300"
              }`}
            >
              Releases
            </Link>
            <Link
              to="/artists"
              className={`rounded-lg px-3 py-1.5 text-sm transition ${
                pathname.startsWith("/artists") ? "text-zinc-100" : "text-zinc-500 hover:text-zinc-300"
              }`}
            >
              Artists
            </Link>
          </nav>

          {backTo && (
            <Link
              to={backTo.href}
              className="ml-auto flex items-center gap-1 text-sm text-zinc-400 transition hover:text-zinc-100"
            >
              <ArrowLeftIcon />
              {backTo.label}
            </Link>
          )}
        </div>
      </header>
      <main className="mx-auto max-w-7xl px-4 py-6 sm:px-6">{children}</main>
    </div>
  );
}
