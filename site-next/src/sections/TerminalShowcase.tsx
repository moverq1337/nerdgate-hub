import {
  TerminalAnimationRoot,
  TerminalAnimationContainer,
  TerminalAnimationWindow,
  TerminalAnimationContent,
  TerminalAnimationCommandBar,
  TerminalAnimationOutput,
  TerminalAnimationTabList,
  TerminalAnimationTabTrigger,
  type TabContent,
} from "@/components/ui/terminal-animation"
import { SectionHeader } from "@/components/SectionHeader"
import { cn } from "@/lib/utils"

const tabs: TabContent[] = [
  {
    label: "install",
    command: "curl -fsSL .../install.sh | sh",
    lines: [
      { text: "", delay: 60 },
      { text: "==> NerdGate Hub installer", color: "text-[#b39aff]", delay: 220 },
      { text: "", delay: 60 },
      { text: "  Panel domain:  nerdgate.example.com", color: "text-neutral-400", delay: 140 },
      { text: "  Let's Encrypt: ops@example.com", color: "text-neutral-400", delay: 140 },
      { text: "  Install dir:   /opt/nerdgate-hub", color: "text-neutral-400", delay: 140 },
      { text: "", delay: 80 },
      { text: "==> DNS preflight", color: "text-[#b39aff]", delay: 200 },
      { text: "  nerdgate.example.com -> 203.0.113.10", color: "text-neutral-400", delay: 140 },
      { text: "  ✓ A record resolves to this server", color: "text-[#22ff73]", delay: 180 },
      { text: "  ✓ Port 80 reachable", color: "text-[#22ff73]", delay: 120 },
      { text: "  ✓ Port 443 reachable", color: "text-[#22ff73]", delay: 120 },
      { text: "", delay: 80 },
      { text: "==> Pulling image", color: "text-[#b39aff]", delay: 200 },
      { text: "  ghcr.io/moverq1337/nerdgate-hub:latest", color: "text-neutral-500", delay: 200 },
      { text: "  Status: Downloaded newer image", color: "text-neutral-400", delay: 180 },
      { text: "", delay: 80 },
      { text: "==> Starting Compose stack", color: "text-[#b39aff]", delay: 220 },
      { text: "  ✓ traefik        running", color: "text-[#22ff73]", delay: 140 },
      { text: "  ✓ nerdgate-hub   running", color: "text-[#22ff73]", delay: 140 },
      { text: "", delay: 80 },
      { text: "==> Setup token (one-time)", color: "text-[#b39aff]", delay: 200 },
      { text: "  6f4d3a-9c1e7b-22fa55", color: "text-[#32f3e9]", delay: 200 },
      { text: "", delay: 80 },
      { text: "Done. Open https://nerdgate.example.com", color: "text-neutral-300", delay: 200 },
    ],
  },
  {
    label: "setup",
    command: "open https://nerdgate.example.com/setup",
    lines: [
      { text: "", delay: 60 },
      { text: "==> First-run setup", color: "text-[#b39aff]", delay: 200 },
      { text: "", delay: 60 },
      { text: "  Setup token:       ************", color: "text-neutral-400", delay: 140 },
      { text: "  Admin username:    admin", color: "text-neutral-400", delay: 140 },
      { text: "  Admin password:    ************", color: "text-neutral-400", delay: 140 },
      { text: "", delay: 80 },
      { text: "  ✓ Session secret generated", color: "text-[#22ff73]", delay: 180 },
      { text: "  ✓ Admin account created", color: "text-[#22ff73]", delay: 140 },
      { text: "  ✓ Session cookie set (HttpOnly · SameSite=Lax)", color: "text-[#22ff73]", delay: 140 },
      { text: "  ✓ Audit log: setup.complete", color: "text-neutral-400", delay: 140 },
      { text: "", delay: 80 },
      { text: "Welcome. You're inside the panel.", color: "text-neutral-300", delay: 200 },
    ],
  },
  {
    label: "route",
    command: "POST /routes  ·  app.example.com",
    lines: [
      { text: "", delay: 60 },
      { text: "==> Creating route", color: "text-[#b39aff]", delay: 200 },
      { text: "", delay: 60 },
      { text: "  Domain:   app.example.com", color: "text-neutral-400", delay: 140 },
      { text: "  Target:   http://host.docker.internal:3000", color: "text-neutral-400", delay: 140 },
      { text: "  TLS:      HTTPS (Let's Encrypt)", color: "text-neutral-400", delay: 140 },
      { text: "", delay: 80 },
      { text: "  ✓ CSRF token validated", color: "text-[#22ff73]", delay: 140 },
      { text: "  ✓ Stored in SQLite (data/app/nerdgate.db)", color: "text-[#22ff73]", delay: 140 },
      { text: "  ✓ Rendered data/traefik/routes.yml", color: "text-[#22ff73]", delay: 140 },
      { text: "  ✓ Traefik picked up dynamic config", color: "text-[#22ff73]", delay: 180 },
      { text: "  ✓ ACME challenge: app.example.com (HTTP-01)", color: "text-[#22ff73]", delay: 200 },
      { text: "  ✓ Certificate issued", color: "text-[#22ff73]", delay: 200 },
      { text: "", delay: 80 },
      { text: "302 → https://app.example.com  (live)", color: "text-neutral-300", delay: 220 },
    ],
  },
  {
    label: "backup",
    command: "Diagnostics → Download backup",
    lines: [
      { text: "", delay: 60 },
      { text: "==> Creating backup", color: "text-[#b39aff]", delay: 200 },
      { text: "", delay: 60 },
      { text: "  + nerdgate.db                 (SQLite, 312 KB)", color: "text-neutral-400", delay: 140 },
      { text: "  + acme.json                   (TLS, 18 KB)", color: "text-neutral-400", delay: 140 },
      { text: "  + routes.json                 (legacy, optional)", color: "text-neutral-500", delay: 140 },
      { text: "  + meta.json                   (backup metadata)", color: "text-neutral-400", delay: 140 },
      { text: "", delay: 80 },
      { text: "  Packed nerdgate-backup-20260527-153002.zip", color: "text-[#22ff73]", delay: 220 },
      { text: "  → 287 KB  ·  sha256: 4f3a…b921", color: "text-neutral-500", delay: 180 },
      { text: "", delay: 80 },
      { text: "Streaming to browser. Done.", color: "text-neutral-300", delay: 200 },
    ],
  },
]

interface TerminalShowcaseProps {
  eyebrow?: string
  title?: string
  description?: string
}

export function TerminalShowcase({
  eyebrow = "What you'll see",
  title = "One terminal session, end to end.",
  description = "Install. Setup. First route. Backup. The whole flow stays explicit — no hidden orchestration.",
}: TerminalShowcaseProps = {}) {
  return (
    <section
      id="terminal"
      className="mx-auto w-full max-w-6xl px-4 py-20 sm:px-6 sm:py-24"
    >
      <SectionHeader eyebrow={eyebrow} title={title} description={description} />

      <TerminalAnimationRoot tabs={tabs} alwaysDark>
        <TerminalAnimationContainer className="px-0 pt-0 md:pt-0">
          <div className="relative rounded-xl border border-border bg-card/70 shadow-[0_30px_80px_-20px_rgba(0,0,0,0.7)] backdrop-blur-sm overflow-hidden">
            <div className="absolute inset-x-0 top-0 h-px bg-gradient-to-r from-transparent via-border-strong to-transparent opacity-60" />

            <TerminalAnimationTabList className="flex flex-wrap gap-1 border-b border-border bg-black/40 px-3 py-2.5">
              {tabs.map((t, i) => (
                <TerminalAnimationTabTrigger
                  key={t.label}
                  index={i}
                  className={cn(
                    "rounded-md border border-transparent px-3 py-1 text-[11px] font-mono uppercase tracking-[0.18em] text-muted-foreground transition-colors",
                    "hover:text-foreground",
                    "data-[state=active]:border-border data-[state=active]:bg-secondary/60 data-[state=active]:text-foreground",
                  )}
                >
                  {t.label}
                </TerminalAnimationTabTrigger>
              ))}
            </TerminalAnimationTabList>

            <TerminalAnimationWindow
              className="rounded-none"
              backgroundColor="#070809"
              minHeight="24rem"
              animateOnVisible
            >
              <TerminalAnimationContent className="font-mono text-[12.5px] leading-relaxed text-neutral-300">
                <div className="flex items-baseline gap-2 pb-3">
                  <span className="text-emerald-400/80">$</span>
                  <TerminalAnimationCommandBar className="text-neutral-100" />
                </div>
                <TerminalAnimationOutput
                  renderLine={(line, _i, visible) =>
                    visible ? (
                      <div
                        className={cn(
                          "whitespace-pre",
                          line.color || "text-neutral-300",
                        )}
                      >
                        {line.text || " "}
                      </div>
                    ) : null
                  }
                />
              </TerminalAnimationContent>
            </TerminalAnimationWindow>
          </div>
        </TerminalAnimationContainer>
      </TerminalAnimationRoot>
    </section>
  )
}
