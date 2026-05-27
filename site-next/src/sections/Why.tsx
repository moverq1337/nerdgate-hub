import { Minimize2, Network, Wand2 } from "lucide-react"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { SectionHeader } from "@/components/SectionHeader"

const items = [
  {
    icon: Minimize2,
    title: "Less surface area",
    body: "NerdGate Hub does not try to become a full platform. Traefik handles traffic. Docker runs apps. The Hub manages routes and admin workflows.",
  },
  {
    icon: Network,
    title: "Traefik-first",
    body: "Routes are generated as Traefik dynamic config. HTTPS and Let's Encrypt stay inside Traefik, where they belong.",
  },
  {
    icon: Wand2,
    title: "Clear install path",
    body: "The installer checks DNS, asks for the Let's Encrypt email, and gives you a one-time setup token for browser onboarding.",
  },
]

export function WhySection() {
  return (
    <section
      id="why"
      className="mx-auto w-full max-w-6xl px-4 py-20 sm:px-6 sm:py-24"
    >
      <SectionHeader
        eyebrow="Why"
        title="Why this exists"
        description="Routing is solved. The dashboards on top of it are usually too heavy. NerdGate Hub picks the smallest useful surface."
      />
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        {items.map(({ icon: Icon, title, body }) => (
          <Card
            key={title}
            className="group relative overflow-hidden border-border bg-card/60 backdrop-blur-sm transition-colors hover:bg-card"
          >
            <div className="absolute inset-x-0 top-0 h-px bg-gradient-to-r from-transparent via-border to-transparent opacity-60" />
            <CardHeader>
              <div className="mb-3 inline-flex h-9 w-9 items-center justify-center rounded-lg border border-border bg-secondary/60 text-foreground/85">
                <Icon className="h-4 w-4" />
              </div>
              <CardTitle className="text-lg font-semibold tracking-tight">{title}</CardTitle>
            </CardHeader>
            <CardContent>
              <CardDescription className="text-sm leading-relaxed text-muted-foreground">
                {body}
              </CardDescription>
            </CardContent>
          </Card>
        ))}
      </div>
    </section>
  )
}
