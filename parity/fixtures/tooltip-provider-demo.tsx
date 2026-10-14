import TooltipDemo from "@/registry/new-york-v4/examples/tooltip-demo"
import { TooltipProvider } from "@/registry/new-york-v4/ui/tooltip"

// Port fixture: tooltip-demo as the shadcn docs render it, inside the
// TooltipProvider (delayDuration 0) their layout mounts around every
// example. Radix Tooltip throws without a provider, so the bare
// tooltip-demo reference cannot render in the harness.
export default function TooltipProviderDemo() {
  return (
    <TooltipProvider>
      <TooltipDemo />
    </TooltipProvider>
  )
}
