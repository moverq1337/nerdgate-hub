import { Archive, KeyRound, Stethoscope, Trash2 } from "lucide-react"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { SectionHeader } from "@/components/SectionHeader"
import { CodeBlock } from "@/components/CodeBlock"

const items = [
  {
    icon: Archive,
    title: "Backup",
    body: "Saves SQLite data, legacy routes when present, and Let's Encrypt acme.json. Restore upload is staged from the panel and applied on restart.",
    code: "curl -fsSL https://raw.githubusercontent.com/moverq1337/nerdgate-hub/main/scripts/backup.sh | sh",
  },
  {
    icon: KeyRound,
    title: "Reset password",
    body: "Change it in Maintenance, or use the helper when you are locked out.",
    code: "curl -fsSL https://raw.githubusercontent.com/moverq1337/nerdgate-hub/main/scripts/reset-password.sh | sh",
  },
  {
    icon: Trash2,
    title: "Uninstall",
    body: "Stops the Compose stack and can remove the install directory.",
    code: "curl -fsSL https://raw.githubusercontent.com/moverq1337/nerdgate-hub/main/scripts/uninstall.sh | sh",
  },
  {
    icon: Stethoscope,
    title: "Diagnostics",
    body: "Prints Traefik logs and hints for DNS, ports, ACME, and certificate failures.",
    code: "curl -fsSL https://raw.githubusercontent.com/moverq1337/nerdgate-hub/main/scripts/diagnose.sh | sh",
  },
]

export function MaintenanceSection() {
  return (
    <section
      id="maintenance"
      className="mx-auto w-full max-w-6xl px-4 py-20 sm:px-6 sm:py-24"
    >
      <SectionHeader
        eyebrow="Maintenance"
        title="Backup, password reset, diagnostics, and uninstall."
        description="Day-two operations ship in the box, both from the panel and as plain shell scripts."
      />
      <div className="grid gap-4 sm:grid-cols-2">
        {items.map(({ icon: Icon, title, body, code }) => (
          <Card
            key={title}
            className="group relative overflow-hidden border-border bg-card/60 backdrop-blur-sm transition-colors hover:bg-card"
          >
            <CardHeader>
              <div className="mb-3 inline-flex h-9 w-9 items-center justify-center rounded-lg border border-border bg-secondary/60 text-foreground/85">
                <Icon className="h-4 w-4" />
              </div>
              <CardTitle className="text-lg font-semibold tracking-tight">{title}</CardTitle>
            </CardHeader>
            <CardContent className="flex flex-col gap-4">
              <CardDescription className="text-sm leading-relaxed text-muted-foreground">
                {body}
              </CardDescription>
              <CodeBlock code={code} />
            </CardContent>
          </Card>
        ))}
      </div>
    </section>
  )
}
