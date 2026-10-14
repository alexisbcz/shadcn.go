import { Button } from "@/registry/new-york-v4/ui/button"
import {
  HoverCard,
  HoverCardContent,
  HoverCardTrigger,
} from "@/registry/new-york-v4/ui/hover-card"
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/registry/new-york-v4/ui/tooltip"

// Port fixture: a tooltip and a hover card open on load (defaultOpen).
export default function TooltipDefaultOpen() {
  return (
    <div className="flex gap-24">
      <TooltipProvider>
        <Tooltip defaultOpen>
          <TooltipTrigger asChild>
            <Button variant="outline">Tooltip</Button>
          </TooltipTrigger>
          <TooltipContent side="right">Open on load</TooltipContent>
        </Tooltip>
      </TooltipProvider>
      <HoverCard defaultOpen>
        <HoverCardTrigger asChild>
          <Button variant="link">Hover card</Button>
        </HoverCardTrigger>
        <HoverCardContent>Open on load too.</HoverCardContent>
      </HoverCard>
    </div>
  )
}
