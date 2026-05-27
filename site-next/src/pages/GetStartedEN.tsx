import { PageShell, DocSection } from "@/components/PageShell"
import { CodeBlock } from "@/components/CodeBlock"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"

export function GetStartedEN() {
  return (
    <PageShell
      eyebrow="Get Started"
      title="Install in one terminal session."
      lead="Before installation, create a DNS record for the panel domain. The installer stops if DNS does not point to the server."
    >
      <DocSection heading="1. Prepare DNS">
        <p>Create an A record:</p>
        <CodeBlock code="nerdgate.example.com -> SERVER_PUBLIC_IP" label="dns" />
        <p>
          Ports <code className="rounded bg-muted/60 px-1 py-0.5 font-mono text-xs">80</code> and{" "}
          <code className="rounded bg-muted/60 px-1 py-0.5 font-mono text-xs">443</code> must be
          reachable from the internet for Let's Encrypt HTTP-01 validation.
        </p>
      </DocSection>

      <DocSection heading="2. Run installer">
        <CodeBlock
          label="curl"
          code="curl -fsSL https://raw.githubusercontent.com/moverq1337/nerdgate-hub/main/scripts/install.sh | sh"
        />
        <CodeBlock
          label="wget"
          code="wget -qO- https://raw.githubusercontent.com/moverq1337/nerdgate-hub/main/scripts/install.sh | sh"
        />
        <p>
          The installer asks for install directory, panel domain, and Let's Encrypt email. It
          prints a one-time setup token after the stack starts.
        </p>
      </DocSection>

      <DocSection heading="3. Open the panel">
        <CodeBlock code="https://nerdgate.example.com" label="url" />
        <p>
          Paste the setup token from the installer, then create the first admin username and
          password in the browser.
        </p>
      </DocSection>

      <DocSection heading="4. Create your first route">
        <p>
          The Domains field accepts one domain or several domains separated by commas, spaces, or
          new lines.
        </p>
        <div className="grid gap-3 sm:grid-cols-3">
          {[
            { title: "Host port", code: "app.example.com -> http://host.docker.internal:3000" },
            { title: "Remote server", code: "api.example.com -> http://203.0.113.10:8080" },
            {
              title: "Docker network",
              code: "Use Attach in the panel or `docker network connect nerdgate-proxy CONTAINER_NAME`",
            },
          ].map((c) => (
            <Card key={c.title} className="bg-card/60 border-border">
              <CardHeader>
                <CardTitle className="text-sm font-semibold tracking-tight">{c.title}</CardTitle>
              </CardHeader>
              <CardContent>
                <code className="block break-words rounded-md border border-border bg-muted/40 px-3 py-2 font-mono text-[11px] leading-relaxed text-foreground/85">
                  {c.code}
                </code>
              </CardContent>
            </Card>
          ))}
        </div>
      </DocSection>

      <DocSection heading="5. Backup, diagnose, reset, or uninstall">
        <CodeBlock code="curl -fsSL https://raw.githubusercontent.com/moverq1337/nerdgate-hub/main/scripts/backup.sh | sh" />
        <CodeBlock code="curl -fsSL https://raw.githubusercontent.com/moverq1337/nerdgate-hub/main/scripts/diagnose.sh | sh" />
        <CodeBlock code="curl -fsSL https://raw.githubusercontent.com/moverq1337/nerdgate-hub/main/scripts/reset-password.sh | sh" />
        <CodeBlock code="curl -fsSL https://raw.githubusercontent.com/moverq1337/nerdgate-hub/main/scripts/uninstall.sh | sh" />
        <p>
          The panel can download backups, stage restore uploads, and change the admin password.
          Helper scripts remain available for host-level recovery, diagnostics, and uninstall.
        </p>
      </DocSection>

      <DocSection heading="Local development">
        <CodeBlock code="cp .env.example .env" />
        <CodeBlock code="docker compose -f docker-compose.yml -f docker-compose.dev.yml up --build" />
        <p>
          Production install pulls{" "}
          <code className="rounded bg-muted/60 px-1 py-0.5 font-mono text-xs">
            ghcr.io/moverq1337/nerdgate-hub:latest
          </code>
          ; local development can still build from source. Replace the example setup token before
          exposing the app to the internet.
        </p>
      </DocSection>
    </PageShell>
  )
}
