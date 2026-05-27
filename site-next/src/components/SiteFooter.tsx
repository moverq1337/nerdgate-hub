export function SiteFooter() {
  return (
    <footer className="border-t border-border mt-32">
      <div className="mx-auto flex w-full max-w-6xl flex-col gap-4 px-4 py-10 text-sm text-muted-foreground sm:flex-row sm:items-center sm:justify-between sm:px-6">
        <div className="flex items-center gap-2.5">
          <img
            src={`${import.meta.env.BASE_URL}logo.png`}
            alt=""
            className="h-5 w-5"
          />
          <span>NerdGate Hub · MIT licensed</span>
        </div>
        <div className="flex flex-wrap items-center gap-x-5 gap-y-2">
          <a
            className="hover:text-foreground transition-colors"
            href="https://github.com/moverq1337/nerdgate-hub"
            target="_blank"
            rel="noopener noreferrer"
          >
            GitHub
          </a>
          <a
            className="hover:text-foreground transition-colors"
            href="https://github.com/moverq1337/nerdgate-hub/pkgs/container/nerdgate-hub"
            target="_blank"
            rel="noopener noreferrer"
          >
            GHCR
          </a>
          <a
            className="hover:text-foreground transition-colors"
            href="https://github.com/moverq1337/nerdgate-hub/issues"
            target="_blank"
            rel="noopener noreferrer"
          >
            Issues
          </a>
        </div>
      </div>
    </footer>
  )
}
