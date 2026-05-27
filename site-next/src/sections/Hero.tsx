import { ArrowRight } from "lucide-react"
import { GithubIcon } from "@/components/icons"
import {
  HeroLiquidMetalRoot,
  HeroLiquidMetalContainer,
  HeroLiquidMetalContent,
  HeroLiquidMetalHeading,
  HeroLiquidMetalDescription,
  HeroLiquidMetalActions,
  HeroLiquidMetalBadges,
  HeroLiquidMetalVisual,
  HeroLiquidMetalMobileVisual,
  type HeroLiquidMetalTechItem,
} from "@/components/ui/hero-liquid-metal"
import { MetalButton } from "@/components/ui/metal-button"
import { Badge } from "@/components/ui/badge"
import { CodeBlock } from "@/components/CodeBlock"
import { LiquidOrb } from "@/components/LiquidOrb"

const BASE = import.meta.env.BASE_URL

const techStack: HeroLiquidMetalTechItem[] = [
  { name: "Go", version: "1.25" },
  { name: "Traefik" },
  { name: "SQLite" },
  { name: "Docker" },
]

const heroImage = `${BASE}brand-silhouette.svg`

export function Hero() {
  return (
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
          Domains to containers.
        </span>
      }
      subtitle={
        <span className="text-muted-foreground/85">No extra platform.</span>
      }
      description={
        <>
          A tiny open-source admin panel over{" "}
          <span className="text-foreground font-medium">Traefik</span> for routing
          domains to Docker containers, host ports, and remote HTTP services.
          Built in Go, shipped as a single binary.
        </>
      }
      showCta={false}
      techStack={techStack}
      renderBadge={(tech, _idx, _defaultBadge) => (
        <Badge
          key={tech.name}
          variant="outline"
          className="rounded-full border-border/70 bg-card/60 px-3 py-1 font-medium text-foreground/85 backdrop-blur-sm transition-colors hover:bg-card hover:text-foreground"
        >
          <span className="tracking-tight">{tech.name}</span>
          {tech.version ? (
            <span className="ml-1.5 font-mono text-[10px] uppercase text-muted-foreground">
              {tech.version}
            </span>
          ) : null}
        </Badge>
      )}
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

            <HeroLiquidMetalHeading
              headingClassName="font-semibold tracking-[-0.04em] text-4xl sm:text-5xl md:text-6xl lg:text-[64px] lg:leading-[1.02] text-balance"
            />

            <HeroLiquidMetalDescription
              descriptionClassName="mx-auto max-w-xl text-base sm:text-lg text-muted-foreground lg:mx-0"
            />

            <HeroLiquidMetalActions>
              <div className="flex flex-col items-center gap-3 sm:flex-row lg:justify-start">
                <MetalButton
                  variant="default"
                  size="lg"
                  className="rounded-xl px-6 font-medium"
                  asChild
                >
                  <a href="https://moverq1337.github.io/nerdgate-hub/en/get-started.html">
                    Get started
                    <ArrowRight className="ml-1 h-4 w-4" />
                  </a>
                </MetalButton>
                <MetalButton
                  variant="outline"
                  size="lg"
                  className="rounded-xl px-6 font-medium"
                  asChild
                >
                  <a
                    href="https://github.com/moverq1337/nerdgate-hub"
                    target="_blank"
                    rel="noopener noreferrer"
                  >
                    <GithubIcon className="mr-1 h-4 w-4" />
                    View on GitHub
                  </a>
                </MetalButton>
              </div>
            </HeroLiquidMetalActions>

            <div className="hidden lg:block">
              <HeroLiquidMetalBadges />
            </div>

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
  )
}
