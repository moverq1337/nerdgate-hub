import { Boxes, Edit3, HeartPulse, ShieldCheck } from "lucide-react"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { SectionHeader } from "@/components/SectionHeader"

const items = [
  {
    icon: Boxes,
    title: "Several domains",
    body: "Create multiple domains at once with commas, spaces, or new lines.",
  },
  {
    icon: Edit3,
    title: "Route editing",
    body: "Update domain, target URL, and HTTPS directly in the routes table.",
  },
  {
    icon: HeartPulse,
    title: "Health and attach",
    body: "See target health, read diagnostics, and attach containers to nerdgate-proxy.",
  },
  {
    icon: ShieldCheck,
    title: "Hardened admin",
    body: "CSRF tokens, login/setup rate limits, strict cookies, security headers, and an audit log are built in.",
  },
]

export function WorkflowSection() {
  return (
    <section
      id="workflow"
      className="mx-auto w-full max-w-6xl px-4 py-20 sm:px-6 sm:py-24"
    >
      <SectionHeader
        eyebrow="Admin workflow"
        title="Edit, check, and attach without leaving the panel."
        description="Everything the operator needs is one click away. No CLI round-trips for everyday changes."
      />
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        {items.map(({ icon: Icon, title, body }) => (
          <Card
            key={title}
            className="group relative overflow-hidden border-border bg-card/60 backdrop-blur-sm transition-colors hover:bg-card"
          >
            <CardHeader>
              <div className="mb-3 inline-flex h-9 w-9 items-center justify-center rounded-lg border border-border bg-secondary/60 text-foreground/85">
                <Icon className="h-4 w-4" />
              </div>
              <CardTitle className="text-base font-semibold tracking-tight">{title}</CardTitle>
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
