import { memo } from "react"
import { LiquidMetal } from "@paper-design/shaders-react"
import { cn } from "@/lib/utils"

const MemoLiquid = memo(LiquidMetal)

interface LiquidOrbProps {
  image: string
  size: number
  colorTint?: string
  contour?: number
  distortion?: number
  softness?: number
  speed?: number
  scale?: number
  shape?: "circle" | "soft" | "none"
  className?: string
}

const shapeClass: Record<NonNullable<LiquidOrbProps["shape"]>, string> = {
  circle: "rounded-full",
  soft: "rounded-[36%]",
  none: "",
}

export function LiquidOrb({
  image,
  size,
  colorTint = "#cfd4dc",
  contour = 0.4,
  distortion = 0.55,
  softness = 0.9,
  speed = 0.5,
  scale = 0.78,
  shape = "circle",
  className,
}: LiquidOrbProps) {
  return (
    <div
      className={cn(
        "pointer-events-none overflow-hidden ring-1 ring-border/50 shadow-[0_30px_60px_-20px_rgba(0,0,0,0.55)]",
        shapeClass[shape],
        className,
      )}
      style={{ width: size, height: size }}
      aria-hidden
    >
      <MemoLiquid
        width={size}
        height={size}
        image={image}
        colorBack="#00000000"
        colorTint={colorTint}
        contour={contour}
        distortion={distortion}
        softness={softness}
        speed={speed}
        scale={scale}
        fit="contain"
        minPixelRatio={1}
        maxPixelCount={size * size}
      />
    </div>
  )
}
