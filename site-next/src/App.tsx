import { Hero } from "@/sections/Hero"
import { SiteHeader } from "@/components/SiteHeader"
import { SiteFooter } from "@/components/SiteFooter"
import { WhySection } from "@/sections/Why"
import { WorkflowSection } from "@/sections/Workflow"
import { MaintenanceSection } from "@/sections/Maintenance"
import { ArchitectureSection } from "@/sections/Architecture"
import { TerminalShowcase } from "@/sections/TerminalShowcase"

export default function App() {
  return (
    <div className="relative min-h-screen bg-background text-foreground">
      <div
        aria-hidden
        className="pointer-events-none fixed inset-0 -z-10 bg-grid bg-grid-fade opacity-60"
      />
      <SiteHeader />
      <main>
        <Hero />
        <TerminalShowcase />
        <WhySection />
        <WorkflowSection />
        <MaintenanceSection />
        <ArchitectureSection />
      </main>
      <SiteFooter />
    </div>
  )
}
