import { Link, useLocation } from "react-router-dom"
import { GithubIcon } from "@/components/icons"

const enLinks = [
  { href: "/", label: "Idea" },
  { href: "/en/get-started", label: "Get Started" },
  { href: "/en/operations", label: "Operations" },
]

const ruLinks = [
  { href: "/ru", label: "Идея" },
  { href: "/ru/get-started", label: "Старт" },
  { href: "/ru/operations", label: "Операции" },
]

export function SiteHeader() {
  const { pathname } = useLocation()
  const isRu = pathname.startsWith("/ru")
  const links = isRu ? ruLinks : enLinks
  const altLang = isRu
    ? { href: "/", label: "EN" }
    : { href: "/ru", label: "RU" }

  return (
    <header className="sticky top-0 z-40 w-full">
      <div className="absolute inset-x-0 top-0 h-full backdrop-blur-xl bg-background/60 border-b border-border" />
      <div className="relative mx-auto flex h-16 w-full max-w-6xl items-center justify-between px-4 sm:px-6">
        <Link
          to={isRu ? "/ru" : "/"}
          className="flex items-center gap-2.5 font-semibold tracking-tight"
        >
          <span className="grid h-8 w-8 place-items-center rounded-lg border border-border bg-secondary/60 shadow-inner">
            <img
              src={`${import.meta.env.BASE_URL}logo.png`}
              alt=""
              className="h-5 w-5"
            />
          </span>
          <span className="text-sm">NerdGate Hub</span>
        </Link>

        <nav className="hidden items-center gap-6 text-sm text-muted-foreground md:flex">
          {links.map((link) => (
            <Link
              key={link.href}
              to={link.href}
              className="transition-colors hover:text-foreground"
            >
              {link.label}
            </Link>
          ))}
        </nav>

        <div className="flex items-center gap-2 text-sm">
          <Link
            to={altLang.href}
            className="rounded-md border border-border bg-secondary/40 px-2.5 py-1.5 font-mono text-[11px] uppercase tracking-wider text-muted-foreground transition-colors hover:text-foreground hover:bg-secondary"
          >
            {altLang.label}
          </Link>
          <a
            href="https://github.com/moverq1337/nerdgate-hub"
            target="_blank"
            rel="noopener noreferrer"
            className="hidden sm:inline-flex items-center gap-2 rounded-lg border border-border bg-secondary/40 px-3 py-1.5 text-muted-foreground transition-colors hover:text-foreground hover:bg-secondary"
          >
            <GithubIcon className="h-4 w-4" />
            GitHub
          </a>
        </div>
      </div>
    </header>
  )
}
