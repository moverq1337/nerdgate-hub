import { PageShell, DocSection } from "@/components/PageShell"
import { CodeBlock } from "@/components/CodeBlock"

export function OperationsEN() {
  return (
    <PageShell
      eyebrow="Operations"
      title="Keep the install understandable."
      lead="The current MVP keeps operational tasks explicit: diagnostics, backup, restore, reset password, uninstall, check DNS, and update containers."
    >
      <DocSection heading="Diagnostics">
        <CodeBlock code="curl -fsSL https://raw.githubusercontent.com/moverq1337/nerdgate-hub/main/scripts/diagnose.sh | sh" />
        <p>
          Without the helper script, open the in-panel Diagnostics block and use the built-in DNS
          command when needed:
        </p>
        <CodeBlock code="docker compose run --rm --no-deps nerdgate-hub check-domain nerdgate.example.com" />
        <CodeBlock code="docker compose logs --tail=160 traefik" />
        <p>
          The panel also shows target health, Traefik config status, Docker API status, and recent
          Traefik logs. The shell script simply collects those host-level checks in one
          interactive flow.
        </p>
      </DocSection>

      <DocSection heading="Backup and restore">
        <p>
          In the panel, use <strong>Diagnostics → Download backup</strong> and{" "}
          <strong>Maintenance → Restore backup</strong>. Restore upload validates the zip, stages
          it, restarts NerdGate Hub, and applies it before SQLite opens.
        </p>
        <CodeBlock code="curl -fsSL https://raw.githubusercontent.com/moverq1337/nerdgate-hub/main/scripts/backup.sh | sh" />
        <CodeBlock code="curl -fsSL https://raw.githubusercontent.com/moverq1337/nerdgate-hub/main/scripts/restore.sh | sh" />
        <p>The same actions can be run directly through the NerdGate Hub binary inside Compose:</p>
        <CodeBlock code="mkdir -p backups" />
        <CodeBlock code='docker compose run --rm --no-deps -v "$PWD/backups:/backup" nerdgate-hub backup /data /acme/acme.json /backup/nerdgate.zip' />
        <CodeBlock code="docker compose down --remove-orphans" />
        <CodeBlock code='docker compose run --rm --no-deps -v "$PWD/backups/nerdgate.zip:/backup.zip:ro" nerdgate-hub restore /data /acme/acme.json /backup.zip' />
        <CodeBlock code="docker compose up -d" />
        <p>
          Backups include SQLite data, optional legacy routes, and Let's Encrypt{" "}
          <code className="rounded bg-muted/60 px-1 py-0.5 font-mono text-xs">acme.json</code>.
          Restore stops the Compose stack before replacing data, then starts it again.
        </p>
      </DocSection>

      <DocSection heading="Security hardening">
        <p>
          The panel uses CSRF tokens on POST forms, rate limits login and first setup, sends
          stricter session cookies, adds browser security headers, and records important actions
          in the audit log.
        </p>
      </DocSection>

      <DocSection heading="Tagged releases">
        <CodeBlock code="git tag v0.1.0" />
        <CodeBlock code="git push origin v0.1.0" />
        <p>Tag pushes create a GitHub Release and publish a matching GHCR image tag.</p>
      </DocSection>

      <DocSection heading="Reset password">
        <p>
          Use <strong>Maintenance → Account password</strong> in the panel when you can still
          sign in.
        </p>
        <CodeBlock code="curl -fsSL https://raw.githubusercontent.com/moverq1337/nerdgate-hub/main/scripts/reset-password.sh | sh" />
        <p>Direct command without the helper script:</p>
        <CodeBlock code="docker compose run --rm --no-deps nerdgate-hub reset-password /data admin 'new-password'" />
        <p>
          This updates SQLite credentials and rotates the session secret, which invalidates
          existing sessions.
        </p>
      </DocSection>

      <DocSection heading="Uninstall">
        <CodeBlock code="curl -fsSL https://raw.githubusercontent.com/moverq1337/nerdgate-hub/main/scripts/uninstall.sh | sh" />
        <p>Direct command without the helper script:</p>
        <CodeBlock code="docker compose down --remove-orphans" />
        <p>
          The script stops the Compose stack and can optionally remove the install directory with
          data. The direct command only stops containers and keeps files on disk.
        </p>
      </DocSection>

      <DocSection heading="Update">
        <CodeBlock code="cd /opt/nerdgate-hub" />
        <CodeBlock code="docker compose pull" />
        <CodeBlock code="docker compose up -d" />
      </DocSection>

      <DocSection heading="DNS diagnostics">
        <CodeBlock code="docker compose run --rm --no-deps nerdgate-hub check-domain nerdgate.example.com" />
        <CodeBlock code="docker compose run --rm --no-deps nerdgate-hub check-domain nerdgate.example.com 203.0.113.10" />
      </DocSection>
    </PageShell>
  )
}
