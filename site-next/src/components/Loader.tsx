import { useEffect, useState } from "react"
import { motion, AnimatePresence } from "motion/react"
import { LiquidOrb } from "@/components/LiquidOrb"

const HOLD_MS = 3000
const FADE_MS = 700

const heroImage = `${import.meta.env.BASE_URL}brand-silhouette.svg`

export function Loader() {
  const [visible, setVisible] = useState(true)
  const [progress, setProgress] = useState(0)

  useEffect(() => {
    // Lock body scroll while loader is up
    const prevOverflow = document.body.style.overflow
    document.body.style.overflow = "hidden"

    const start = performance.now()
    let frame = 0
    const tick = (now: number) => {
      const elapsed = now - start
      const ratio = Math.min(1, elapsed / HOLD_MS)
      // ease-out cubic
      const eased = 1 - Math.pow(1 - ratio, 3)
      setProgress(Math.round(eased * 100))
      if (ratio < 1) {
        frame = requestAnimationFrame(tick)
      }
    }
    frame = requestAnimationFrame(tick)

    const hide = setTimeout(() => {
      setVisible(false)
      document.body.style.overflow = prevOverflow
    }, HOLD_MS)

    return () => {
      cancelAnimationFrame(frame)
      clearTimeout(hide)
      document.body.style.overflow = prevOverflow
    }
  }, [])

  return (
    <AnimatePresence>
      {visible && (
        <motion.div
          initial={{ opacity: 1 }}
          exit={{ opacity: 0, scale: 1.03, filter: "blur(8px)" }}
          transition={{ duration: FADE_MS / 1000, ease: [0.16, 1, 0.3, 1] }}
          className="fixed inset-0 z-[100] grid place-items-center bg-background"
          aria-hidden="true"
        >
          {/* ambient radial glow behind orb */}
          <div
            className="pointer-events-none absolute inset-0"
            style={{
              background:
                "radial-gradient(ellipse 50% 40% at 50% 45%, rgba(255,255,255,0.08), transparent 70%)",
            }}
          />

          <motion.div
            initial={{ opacity: 0, scale: 0.82 }}
            animate={{ opacity: 1, scale: 1 }}
            transition={{ duration: 0.6, ease: [0.16, 1, 0.3, 1] }}
            className="relative flex flex-col items-center gap-9"
          >
            <LiquidOrb
              image={heroImage}
              size={180}
              colorTint="#cfd4dc"
              contour={0.4}
              distortion={0.6}
              softness={0.92}
              speed={0.85}
              scale={0.82}
            />

            <div className="flex flex-col items-center gap-4">
              <motion.span
                initial={{ opacity: 0, y: 6 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: 0.25, duration: 0.5, ease: [0.16, 1, 0.3, 1] }}
                className="text-[13px] font-medium tracking-[0.16em] uppercase text-foreground/80"
              >
                NerdGate Hub
              </motion.span>

              <div
                className="relative h-px w-48 overflow-hidden rounded-full"
                style={{ background: "rgba(255,255,255,0.08)" }}
              >
                <motion.div
                  className="absolute inset-y-0 left-0"
                  style={{
                    width: `${progress}%`,
                    background:
                      "linear-gradient(90deg, rgba(255,255,255,0.4), rgba(255,255,255,0.95))",
                    boxShadow: "0 0 10px rgba(255,255,255,0.4)",
                  }}
                  transition={{ duration: 0.12, ease: "linear" }}
                />
              </div>

              <span className="font-mono text-[10.5px] tabular-nums tracking-[0.14em] text-muted-foreground">
                {String(progress).padStart(3, "0")}%
              </span>
            </div>
          </motion.div>
        </motion.div>
      )}
    </AnimatePresence>
  )
}
