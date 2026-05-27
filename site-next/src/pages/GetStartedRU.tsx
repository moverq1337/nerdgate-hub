import { PageShell, DocSection } from "@/components/PageShell"
import { CodeBlock } from "@/components/CodeBlock"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"

export function GetStartedRU() {
  return (
    <PageShell
      lang="ru"
      eyebrow="Старт"
      title="Установка за одну терминальную сессию."
      lead="Перед установкой создай DNS-запись для домена панели. Инсталлер остановится, если DNS не указывает на сервер."
    >
      <DocSection heading="1. Подготовь DNS">
        <p>Создай A-запись:</p>
        <CodeBlock code="nerdgate.example.com -> SERVER_PUBLIC_IP" label="dns" />
        <p>
          Порты <code className="rounded bg-muted/60 px-1 py-0.5 font-mono text-xs">80</code> и{" "}
          <code className="rounded bg-muted/60 px-1 py-0.5 font-mono text-xs">443</code> должны
          быть доступны из интернета для Let's Encrypt HTTP-01 validation.
        </p>
      </DocSection>

      <DocSection heading="2. Запусти инсталлер">
        <CodeBlock
          label="curl"
          code="curl -fsSL https://raw.githubusercontent.com/moverq1337/nerdgate-hub/main/scripts/install.sh | sh"
        />
        <CodeBlock
          label="wget"
          code="wget -qO- https://raw.githubusercontent.com/moverq1337/nerdgate-hub/main/scripts/install.sh | sh"
        />
        <p>
          Инсталлер спросит папку установки, домен панели и email для Let's Encrypt. После
          запуска stack он покажет одноразовый setup token.
        </p>
      </DocSection>

      <DocSection heading="3. Открой панель">
        <CodeBlock code="https://nerdgate.example.com" label="url" />
        <p>
          Вставь setup token из installer, затем создай первый admin username и password в
          браузере.
        </p>
      </DocSection>

      <DocSection heading="4. Создай первый route">
        <p>
          Поле Domains принимает один домен или несколько доменов через запятые, пробелы или
          новые строки.
        </p>
        <div className="grid gap-3 sm:grid-cols-3">
          {[
            { title: "Порт хоста", code: "app.example.com -> http://host.docker.internal:3000" },
            { title: "Другой сервер", code: "api.example.com -> http://203.0.113.10:8080" },
            { title: "Docker network", code: "Используй Attach в панели или `docker network connect nerdgate-proxy CONTAINER_NAME`" },
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

      <DocSection heading="5. Backup, диагностика, reset или uninstall">
        <CodeBlock code="curl -fsSL https://raw.githubusercontent.com/moverq1337/nerdgate-hub/main/scripts/backup.sh | sh" />
        <CodeBlock code="curl -fsSL https://raw.githubusercontent.com/moverq1337/nerdgate-hub/main/scripts/diagnose.sh | sh" />
        <CodeBlock code="curl -fsSL https://raw.githubusercontent.com/moverq1337/nerdgate-hub/main/scripts/reset-password.sh | sh" />
        <CodeBlock code="curl -fsSL https://raw.githubusercontent.com/moverq1337/nerdgate-hub/main/scripts/uninstall.sh | sh" />
        <p>
          Панель умеет скачивать backups, принимать restore uploads и менять admin password.
          Helper-скрипты остаются для host-level recovery, diagnostics и uninstall.
        </p>
      </DocSection>

      <DocSection heading="Локальная разработка">
        <CodeBlock code="cp .env.example .env" />
        <CodeBlock code="docker compose -f docker-compose.yml -f docker-compose.dev.yml up --build" />
        <p>
          Production-установка тянет{" "}
          <code className="rounded bg-muted/60 px-1 py-0.5 font-mono text-xs">
            ghcr.io/moverq1337/nerdgate-hub:latest
          </code>
          , а dev-режим билдит из исходников. Замени пример setup token перед тем, как открывать
          приложение в интернет.
        </p>
      </DocSection>
    </PageShell>
  )
}
