import { useState } from "react"
import { Check, Copy } from "lucide-react"
import { cn } from "@/lib/utils"

interface CodeBlockProps {
  code: string
  label?: string
  className?: string
}

export function CodeBlock({ code, label, className }: CodeBlockProps) {
  const [copied, setCopied] = useState(false)

  const copy = async () => {
    try {
      await navigator.clipboard.writeText(code)
      setCopied(true)
      setTimeout(() => setCopied(false), 1500)
    } catch {
      /* clipboard unavailable; silent */
    }
  }

  return (
    <div
      className={cn(
        "group relative flex w-full max-w-full items-center gap-3 overflow-hidden rounded-xl border border-border bg-card/70 px-4 py-3 ring-glow",
        className,
      )}
    >
      {label ? (
        <span className="shrink-0 text-xs font-medium uppercase tracking-wider text-muted-foreground/80">
          {label}
        </span>
      ) : (
        <span className="shrink-0 text-muted-foreground/60 select-none">$</span>
      )}
      <code className="block min-w-0 flex-1 overflow-hidden truncate font-mono text-xs text-foreground sm:text-sm">
        {code}
      </code>
      <button
        type="button"
        onClick={copy}
        aria-label="Copy to clipboard"
        className="shrink-0 rounded-md border border-transparent p-1.5 text-muted-foreground transition-colors hover:border-border hover:text-foreground"
      >
        {copied ? (
          <Check className="h-3.5 w-3.5 text-emerald-400" />
        ) : (
          <Copy className="h-3.5 w-3.5" />
        )}
      </button>
    </div>
  )
}
