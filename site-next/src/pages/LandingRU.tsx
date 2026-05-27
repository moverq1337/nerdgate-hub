import { ArrowRight } from "lucide-react"
import { SiteHeader } from "@/components/SiteHeader"
import { SiteFooter } from "@/components/SiteFooter"
import { SectionHeader } from "@/components/SectionHeader"
import { CodeBlock } from "@/components/CodeBlock"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { MetalButton } from "@/components/ui/metal-button"
import { GithubIcon } from "@/components/icons"
import { Archive, Boxes, Edit3, HeartPulse, KeyRound, Minimize2, Network, ShieldCheck, Stethoscope, Trash2, Wand2 } from "lucide-react"
import {
  HeroLiquidMetalRoot,
  HeroLiquidMetalContainer,
  HeroLiquidMetalContent,
  HeroLiquidMetalHeading,
  HeroLiquidMetalDescription,
  HeroLiquidMetalActions,
  HeroLiquidMetalVisual,
  HeroLiquidMetalMobileVisual,
} from "@/components/ui/hero-liquid-metal"
import { LiquidOrb } from "@/components/LiquidOrb"
import { TerminalShowcase } from "@/sections/TerminalShowcase"

const BASE = import.meta.env.BASE_URL

const heroImage = `${BASE}brand-silhouette.svg`

const why = [
  { icon: Minimize2, title: "Меньше лишнего", body: "NerdGate Hub не пытается стать большой платформой. Traefik принимает трафик, Docker запускает приложения, Hub управляет маршрутами." },
  { icon: Network, title: "Traefik-first", body: "Маршруты генерируются как dynamic config для Traefik. HTTPS и Let's Encrypt остаются внутри Traefik." },
  { icon: Wand2, title: "Понятная установка", body: "Инсталлер проверяет DNS, спрашивает email для Let's Encrypt и даёт одноразовый setup token для onboarding в браузере." },
]

const workflow = [
  { icon: Boxes, title: "Несколько доменов", body: "Создавай несколько доменов сразу через запятые, пробелы или новые строки." },
  { icon: Edit3, title: "Редактирование", body: "Меняй domain, target URL и HTTPS прямо в таблице маршрутов." },
  { icon: HeartPulse, title: "Health и attach", body: "Смотри health targets, диагностику и подключай контейнеры к nerdgate-proxy." },
  { icon: ShieldCheck, title: "Hardened admin", body: "CSRF-токены, rate limit для login/setup, строгие cookies, security headers и audit log уже встроены." },
]

const maintenance = [
  { icon: Archive, title: "Backup", body: "SQLite, legacy routes (если есть) и Let's Encrypt acme.json. Restore upload ставится из панели и применяется при следующем старте.", code: "curl -fsSL https://raw.githubusercontent.com/moverq1337/nerdgate-hub/main/scripts/backup.sh | sh" },
  { icon: KeyRound, title: "Сброс пароля", body: "Меняй пароль в Maintenance или используй helper, если потерял доступ.", code: "curl -fsSL https://raw.githubusercontent.com/moverq1337/nerdgate-hub/main/scripts/reset-password.sh | sh" },
  { icon: Trash2, title: "Удаление", body: "Останавливает Compose stack и при желании удаляет папку установки.", code: "curl -fsSL https://raw.githubusercontent.com/moverq1337/nerdgate-hub/main/scripts/uninstall.sh | sh" },
  { icon: Stethoscope, title: "Диагностика", body: "Показывает Traefik logs и подсказки по DNS, портам, ACME и ошибкам сертификатов.", code: "curl -fsSL https://raw.githubusercontent.com/moverq1337/nerdgate-hub/main/scripts/diagnose.sh | sh" },
]

const architecture = [
  "DNS домена панели указывает на сервер.",
  "Traefik слушает порты 80 и 443 и владеет TLS через Let's Encrypt.",
  "NerdGate Hub хранит маршруты, пользователей и настройки в SQLite.",
  "Admin POST-действия используют CSRF-токены, важные изменения пишутся в audit log.",
  "NerdGate Hub пишет Traefik dynamic config в data/traefik/routes.yml.",
  "Traefik следит за файлом и применяет изменения автоматически.",
]

export function LandingRU() {
  return (
    <div className="relative min-h-screen bg-background text-foreground" lang="ru">
      <div
        aria-hidden
        className="pointer-events-none fixed inset-0 -z-10 bg-grid bg-grid-fade opacity-60"
      />
      <SiteHeader />
      <main>
        <HeroLiquidMetalRoot
          className="relative isolate"
          image={heroImage}
          colorBack="#00000000"
          colorTint="#cfd4dc"
          contour={0.4}
          distortion={0.55}
          softness={0.9}
          speed={0.5}
          scale={0.78}
          desktopShaderProps={{ width: 720, height: 720 }}
          title={
            <span className="bg-gradient-to-b from-white via-white to-white/50 bg-clip-text text-transparent">
              Домены к контейнерам.
            </span>
          }
          subtitle={
            <span className="text-muted-foreground/85">Без лишней платформы.</span>
          }
          description={
            <>
              Маленькая open-source панель поверх{" "}
              <span className="text-foreground font-medium">Traefik</span>: направляй домены к Docker-контейнерам, портам хоста и внешним HTTP-сервисам. Написана на Go, поставляется как один бинарь.
            </>
          }
          showCta={false}
          showBadges={false}
        >
          <div className="mx-auto w-full max-w-6xl px-4 pt-16 pb-10 sm:px-6 sm:pt-20 lg:pt-28">
            <HeroLiquidMetalContainer className="grid gap-8 lg:grid-cols-[1fr_minmax(280px,440px)] lg:gap-16 lg:items-center pb-0">
              <HeroLiquidMetalContent className="gap-6 text-center lg:text-left lg:gap-7">
                <Badge
                  variant="outline"
                  className="mx-auto w-fit rounded-full border-border bg-card/70 px-3 py-1 text-[11px] uppercase tracking-[0.18em] text-muted-foreground lg:mx-0"
                >
                  <span className="mr-1.5 h-1.5 w-1.5 rounded-full bg-emerald-400/90 shadow-[0_0_8px] shadow-emerald-400/70" />
                  Open source · MIT
                </Badge>

                <HeroLiquidMetalHeading headingClassName="font-semibold tracking-[-0.04em] text-4xl sm:text-5xl md:text-6xl lg:text-[64px] lg:leading-[1.02] text-balance" />

                <HeroLiquidMetalDescription descriptionClassName="mx-auto max-w-xl text-base sm:text-lg text-muted-foreground lg:mx-0" />

                <HeroLiquidMetalActions>
                  <div className="flex flex-col items-center gap-3 sm:flex-row lg:justify-start">
                    <MetalButton variant="default" size="lg" className="rounded-xl px-6 font-medium" asChild>
                      <a href={`${import.meta.env.BASE_URL}ru/get-started`}>
                        Начать
                        <ArrowRight className="ml-1 h-4 w-4" />
                      </a>
                    </MetalButton>
                    <MetalButton variant="outline" size="lg" className="rounded-xl px-6 font-medium" asChild>
                      <a href="https://github.com/moverq1337/nerdgate-hub" target="_blank" rel="noopener noreferrer">
                        <GithubIcon className="mr-1 h-4 w-4" />
                        На GitHub
                      </a>
                    </MetalButton>
                  </div>
                </HeroLiquidMetalActions>

                <div className="flex w-full flex-col gap-2 pt-2 lg:max-w-xl">
                  <CodeBlock
                    label="curl"
                    code="curl -fsSL https://raw.githubusercontent.com/moverq1337/nerdgate-hub/main/scripts/install.sh | sh"
                  />
                  <CodeBlock
                    label="wget"
                    code="wget -qO- https://raw.githubusercontent.com/moverq1337/nerdgate-hub/main/scripts/install.sh | sh"
                  />
                </div>
              </HeroLiquidMetalContent>

              <div className="relative hidden lg:block lg:h-[460px] xl:h-[520px]">
                <HeroLiquidMetalVisual
                  className="absolute inset-0 h-full"
                  desktopClassName="rounded-full ring-1 ring-border/60 shadow-[0_40px_80px_-20px_rgba(0,0,0,0.6)]"
                />
                <LiquidOrb
                  image={`${BASE}orb-hex.svg`}
                  size={96}
                  colorTint="#9ca3af"
                  speed={0.35}
                  contour={0.5}
                  className="absolute right-[-10px] top-[-12px]"
                />
                <LiquidOrb
                  image={`${BASE}orb-disc.svg`}
                  size={64}
                  colorTint="#b8c3cf"
                  speed={0.6}
                  className="absolute bottom-[-8px] left-[8%]"
                />
                <LiquidOrb
                  image={`${BASE}orb-pill.svg`}
                  size={56}
                  shape="soft"
                  colorTint="#a8b3c0"
                  speed={0.45}
                  className="absolute top-[45%] left-[-22px]"
                />
              </div>
            </HeroLiquidMetalContainer>
          </div>
          <HeroLiquidMetalMobileVisual className="-bottom-32 -z-10" />
        </HeroLiquidMetalRoot>

        <TerminalShowcase
          eyebrow="Что увидишь"
          title="Одна терминальная сессия от и до."
          description="Установка. Setup. Первый route. Backup. Весь flow остаётся явным — без скрытой оркестровки."
        />

        <section className="mx-auto w-full max-w-6xl px-4 py-20 sm:px-6 sm:py-24">
          <SectionHeader
            eyebrow="Зачем"
            title="Зачем это нужно"
            description="Маршрутизация решена. Панели поверх неё обычно слишком тяжёлые. NerdGate Hub выбирает самую маленькую полезную поверхность."
          />
          <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
            {why.map(({ icon: Icon, title, body }) => (
              <Card key={title} className="border-border bg-card/60 backdrop-blur-sm transition-colors hover:bg-card">
                <CardHeader>
                  <div className="mb-3 inline-flex h-9 w-9 items-center justify-center rounded-lg border border-border bg-secondary/60">
                    <Icon className="h-4 w-4" />
                  </div>
                  <CardTitle className="text-lg font-semibold tracking-tight">{title}</CardTitle>
                </CardHeader>
                <CardContent>
                  <CardDescription className="text-sm leading-relaxed text-muted-foreground">{body}</CardDescription>
                </CardContent>
              </Card>
            ))}
          </div>
        </section>

        <section className="mx-auto w-full max-w-6xl px-4 py-20 sm:px-6 sm:py-24">
          <SectionHeader
            eyebrow="Admin workflow"
            title="Редактирование, проверки и attach прямо в панели."
            description="Всё, что нужно оператору, в одном клике."
          />
          <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
            {workflow.map(({ icon: Icon, title, body }) => (
              <Card key={title} className="border-border bg-card/60 backdrop-blur-sm transition-colors hover:bg-card">
                <CardHeader>
                  <div className="mb-3 inline-flex h-9 w-9 items-center justify-center rounded-lg border border-border bg-secondary/60">
                    <Icon className="h-4 w-4" />
                  </div>
                  <CardTitle className="text-base font-semibold tracking-tight">{title}</CardTitle>
                </CardHeader>
                <CardContent>
                  <CardDescription className="text-sm leading-relaxed text-muted-foreground">{body}</CardDescription>
                </CardContent>
              </Card>
            ))}
          </div>
        </section>

        <section className="mx-auto w-full max-w-6xl px-4 py-20 sm:px-6 sm:py-24">
          <SectionHeader
            eyebrow="Обслуживание"
            title="Backup, сброс пароля, диагностика и uninstall."
            description="Day-two операции работают из коробки — из панели и из shell-скриптов."
          />
          <div className="grid gap-4 sm:grid-cols-2">
            {maintenance.map(({ icon: Icon, title, body, code }) => (
              <Card key={title} className="border-border bg-card/60 backdrop-blur-sm transition-colors hover:bg-card">
                <CardHeader>
                  <div className="mb-3 inline-flex h-9 w-9 items-center justify-center rounded-lg border border-border bg-secondary/60">
                    <Icon className="h-4 w-4" />
                  </div>
                  <CardTitle className="text-lg font-semibold tracking-tight">{title}</CardTitle>
                </CardHeader>
                <CardContent className="flex flex-col gap-4">
                  <CardDescription className="text-sm leading-relaxed text-muted-foreground">{body}</CardDescription>
                  <CodeBlock code={code} />
                </CardContent>
              </Card>
            ))}
          </div>
        </section>

        <section className="mx-auto w-full max-w-6xl px-4 py-20 sm:px-6 sm:py-24">
          <SectionHeader eyebrow="Архитектура" title="Как работает." description="Шесть частей. Ни одна не делает работу другой." />
          <ol className="relative space-y-3 border-l border-border pl-6 sm:pl-8">
            {architecture.map((step, i) => (
              <li key={i} className="group relative">
                <span
                  aria-hidden
                  className="absolute -left-[33px] sm:-left-[41px] top-0.5 grid h-6 w-6 place-items-center rounded-full border border-border bg-secondary/80 text-[11px] font-mono text-muted-foreground"
                >
                  {i + 1}
                </span>
                <p className="rounded-lg border border-border bg-card/50 px-4 py-3 text-sm leading-relaxed text-foreground/90 sm:text-base">
                  {step}
                </p>
              </li>
            ))}
          </ol>
        </section>
      </main>
      <SiteFooter />
    </div>
  )
}
