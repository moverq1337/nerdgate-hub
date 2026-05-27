import { PageShell, DocSection } from "@/components/PageShell"
import { CodeBlock } from "@/components/CodeBlock"

export function OperationsRU() {
  return (
    <PageShell
      lang="ru"
      eyebrow="Операции"
      title="Обслуживание должно быть понятным."
      lead="В MVP основные операции явные: диагностика, backup, restore, сброс пароля, удаление, проверка DNS и обновление контейнеров."
    >
      <DocSection heading="Диагностика">
        <CodeBlock code="curl -fsSL https://raw.githubusercontent.com/moverq1337/nerdgate-hub/main/scripts/diagnose.sh | sh" />
        <p>
          Без helper script можно открыть Diagnostics в панели и использовать встроенную
          DNS-команду при необходимости:
        </p>
        <CodeBlock code="docker compose run --rm --no-deps nerdgate-hub check-domain nerdgate.example.com" />
        <CodeBlock code="docker compose logs --tail=160 traefik" />
        <p>
          Панель тоже показывает target health, статус Traefik config, Docker API и свежие
          Traefik logs. Shell-скрипт просто собирает host-level проверки в один интерактивный
          flow.
        </p>
      </DocSection>

      <DocSection heading="Backup и restore">
        <p>
          В панели используй <strong>Diagnostics → Download backup</strong> и{" "}
          <strong>Maintenance → Restore backup</strong>. Restore upload проверяет zip, ставит его
          в pending, перезапускает NerdGate Hub и применяет restore до открытия SQLite.
        </p>
        <CodeBlock code="curl -fsSL https://raw.githubusercontent.com/moverq1337/nerdgate-hub/main/scripts/backup.sh | sh" />
        <CodeBlock code="curl -fsSL https://raw.githubusercontent.com/moverq1337/nerdgate-hub/main/scripts/restore.sh | sh" />
        <p>Те же действия можно запускать напрямую через бинарь NerdGate Hub внутри Compose:</p>
        <CodeBlock code="mkdir -p backups" />
        <CodeBlock code='docker compose run --rm --no-deps -v "$PWD/backups:/backup" nerdgate-hub backup /data /acme/acme.json /backup/nerdgate.zip' />
        <CodeBlock code="docker compose down --remove-orphans" />
        <CodeBlock code='docker compose run --rm --no-deps -v "$PWD/backups/nerdgate.zip:/backup.zip:ro" nerdgate-hub restore /data /acme/acme.json /backup.zip' />
        <CodeBlock code="docker compose up -d" />
        <p>
          Backup содержит SQLite-данные, optional legacy routes и Let's Encrypt{" "}
          <code className="rounded bg-muted/60 px-1 py-0.5 font-mono text-xs">acme.json</code>.
          Restore останавливает Compose stack, заменяет данные и запускает stack обратно.
        </p>
      </DocSection>

      <DocSection heading="Security hardening">
        <p>
          Панель использует CSRF-токены на POST-формах, rate limit для login и first setup, более
          строгие session cookies, security headers браузера и audit log важных действий.
        </p>
      </DocSection>

      <DocSection heading="Tagged releases">
        <CodeBlock code="git tag v0.1.0" />
        <CodeBlock code="git push origin v0.1.0" />
        <p>Push тега создаёт GitHub Release и публикует такой же GHCR image tag.</p>
      </DocSection>

      <DocSection heading="Сброс пароля">
        <p>
          Если доступ к панели ещё есть, используй <strong>Maintenance → Account password</strong>.
        </p>
        <CodeBlock code="curl -fsSL https://raw.githubusercontent.com/moverq1337/nerdgate-hub/main/scripts/reset-password.sh | sh" />
        <p>Прямая команда без helper script:</p>
        <CodeBlock code="docker compose run --rm --no-deps nerdgate-hub reset-password /data admin 'new-password'" />
        <p>
          Скрипт обновляет credentials в SQLite и ротирует session secret, поэтому старые сессии
          перестают работать.
        </p>
      </DocSection>

      <DocSection heading="Удаление">
        <CodeBlock code="curl -fsSL https://raw.githubusercontent.com/moverq1337/nerdgate-hub/main/scripts/uninstall.sh | sh" />
        <p>Прямая команда без helper script:</p>
        <CodeBlock code="docker compose down --remove-orphans" />
        <p>
          Скрипт останавливает Compose stack и может удалить папку установки вместе с данными.
          Прямая команда только останавливает контейнеры и оставляет файлы на диске.
        </p>
      </DocSection>

      <DocSection heading="Обновление">
        <CodeBlock code="cd /opt/nerdgate-hub" />
        <CodeBlock code="docker compose pull" />
        <CodeBlock code="docker compose up -d" />
      </DocSection>

      <DocSection heading="DNS-диагностика">
        <CodeBlock code="docker compose run --rm --no-deps nerdgate-hub check-domain nerdgate.example.com" />
        <CodeBlock code="docker compose run --rm --no-deps nerdgate-hub check-domain nerdgate.example.com 203.0.113.10" />
      </DocSection>
    </PageShell>
  )
}
