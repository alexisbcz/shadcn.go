// The upstream app's own radix-ui (a bare "radix-ui" would resolve from
// parity/node_modules, a second React tree).
import { HoverCard as HoverCardPrimitive } from "@/node_modules/radix-ui"

import { Button } from "@/registry/new-york-v4/ui/button"
import {
  HoverCard,
  HoverCardContent,
  HoverCardTrigger,
} from "@/registry/new-york-v4/ui/hover-card"

// Port fixture: hover card sides and alignments with short delays, an
// arrow, a plain <a> trigger, and focusable content (taken out of the tab
// order by Radix).
const variants = [
  { side: "top", align: "start" },
  { side: "right", align: "center" },
  { side: "bottom", align: "end" },
  { side: "left", align: "center" },
] as const

export default function HoverCardSides() {
  return (
    <div className="flex flex-col items-center gap-8">
      <div className="flex gap-24">
        {variants.map((v) => (
          <HoverCard key={v.side} openDelay={100} closeDelay={100}>
            <HoverCardTrigger asChild>
              <Button variant="outline">{v.side}</Button>
            </HoverCardTrigger>
            <HoverCardContent side={v.side} align={v.align} className="w-48">
              <p className="text-sm">{`Card on the ${v.side}.`}</p>
              <a href="#more" className="text-sm underline">More</a>
              <HoverCardPrimitive.Arrow className="fill-popover" />
            </HoverCardContent>
          </HoverCard>
        ))}
      </div>
      <HoverCard openDelay={100} closeDelay={100}>
        <HoverCardTrigger href="#profile" className="text-sm underline">
          Plain link
        </HoverCardTrigger>
        <HoverCardContent>Default placement.</HoverCardContent>
      </HoverCard>
    </div>
  )
}
