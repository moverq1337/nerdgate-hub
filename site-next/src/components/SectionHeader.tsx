interface SectionHeaderProps {
  eyebrow: string
  title: string
  description?: string
}

export function SectionHeader({ eyebrow, title, description }: SectionHeaderProps) {
  return (
    <div className="mb-10 flex flex-col gap-3 sm:mb-12">
      <p className="text-xs font-medium uppercase tracking-[0.22em] text-muted-foreground/80">
        {eyebrow}
      </p>
      <h2 className="text-balance text-2xl font-semibold tracking-[-0.02em] text-foreground sm:text-3xl md:text-4xl">
        {title}
      </h2>
      {description ? (
        <p className="max-w-2xl text-sm leading-relaxed text-muted-foreground sm:text-base">
          {description}
        </p>
      ) : null}
    </div>
  )
}
