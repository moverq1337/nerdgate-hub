import type { ReactNode } from "react"
import { SiteHeader } from "@/components/SiteHeader"
import { SiteFooter } from "@/components/SiteFooter"
import { SectionHeader } from "@/components/SectionHeader"

interface PageShellProps {
  eyebrow: string
  title: string
  lead: string
  lang?: "en" | "ru"
  children: ReactNode
}

export function PageShell({ eyebrow, title, lead, lang = "en", children }: PageShellProps) {
  return (
    <div className="relative min-h-screen bg-background text-foreground" lang={lang}>
      <div
        aria-hidden
        className="pointer-events-none fixed inset-0 -z-10 bg-grid bg-grid-fade opacity-60"
      />
      <SiteHeader />
      <main className="mx-auto w-full max-w-4xl px-4 pt-20 pb-16 sm:px-6 sm:pt-24">
        <SectionHeader eyebrow={eyebrow} title={title} description={lead} />
        <div className="flex flex-col gap-14">{children}</div>
      </main>
      <SiteFooter />
    </div>
  )
}

interface DocSectionProps {
  heading: string
  children: ReactNode
}

export function DocSection({ heading, children }: DocSectionProps) {
  return (
    <section className="flex flex-col gap-4">
      <h2 className="text-xl font-semibold tracking-tight text-foreground sm:text-2xl">
        {heading}
      </h2>
      <div className="flex flex-col gap-3 text-sm leading-relaxed text-muted-foreground sm:text-base">
        {children}
      </div>
    </section>
  )
}
