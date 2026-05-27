import { SectionHeader } from "@/components/SectionHeader"

const steps = [
  "DNS points your panel domain to the server.",
  "Traefik listens on ports 80 and 443 and owns TLS via Let's Encrypt.",
  "NerdGate Hub stores routes, users, and settings in SQLite.",
  "Admin POST actions use CSRF tokens; important changes are written to the audit log.",
  "NerdGate Hub writes Traefik dynamic config to data/traefik/routes.yml.",
  "Traefik watches that file and applies changes automatically.",
]

export function ArchitectureSection() {
  return (
    <section
      id="architecture"
      className="mx-auto w-full max-w-6xl px-4 py-20 sm:px-6 sm:py-24"
    >
      <SectionHeader
        eyebrow="Architecture"
        title="How it works."
        description="Six moving parts. None of them try to do each other's job."
      />
      <ol className="relative space-y-3 border-l border-border pl-6 sm:pl-8">
        {steps.map((step, i) => (
          <li key={i} className="group relative">
            <span
              aria-hidden
              className="absolute -left-[33px] sm:-left-[41px] top-0.5 grid h-6 w-6 place-items-center rounded-full border border-border bg-secondary/80 text-[11px] font-mono text-muted-foreground transition-colors group-hover:text-foreground group-hover:bg-secondary"
            >
              {i + 1}
            </span>
            <p className="rounded-lg border border-border bg-card/50 px-4 py-3 text-sm leading-relaxed text-foreground/90 transition-colors group-hover:bg-card sm:text-base">
              {step}
            </p>
          </li>
        ))}
      </ol>
    </section>
  )
}
