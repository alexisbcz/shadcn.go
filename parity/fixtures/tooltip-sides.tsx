import { Button } from "@/registry/new-york-v4/ui/button"
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/registry/new-york-v4/ui/tooltip"

// Port fixture: tooltip sides, alignments and offsets (the arrow follows
// the placed side and alignment), one that flips at the viewport edge and
// one labelled with aria-label.
const variants = [
  { label: "top", props: {} },
  { label: "right", props: { side: "right" } },
  { label: "bottom", props: { side: "bottom", sideOffset: 8 } },
  { label: "left", props: { side: "left", sideOffset: 4 } },
  { label: "top start", props: { align: "start" } },
  { label: "bottom end", props: { side: "bottom", align: "end", alignOffset: 6 } },
] as const

export default function TooltipSides() {
  return (
    <TooltipProvider>
      <div className="grid grid-cols-3 gap-x-16 gap-y-10">
        {variants.map((v) => (
          <Tooltip key={v.label}>
            <TooltipTrigger asChild>
              <Button variant="outline">{v.label}</Button>
            </TooltipTrigger>
            <TooltipContent {...v.props}>
              <p>{`Tooltip ${v.label}`}</p>
            </TooltipContent>
          </Tooltip>
        ))}
      </div>
      <div className="fixed top-2 left-2">
        <Tooltip>
          <TooltipTrigger asChild>
            <Button variant="outline" size="sm">edge</Button>
          </TooltipTrigger>
          <TooltipContent>Flipped below</TooltipContent>
        </Tooltip>
      </div>
      <div className="fixed right-2 bottom-2">
        <Tooltip>
          <TooltipTrigger asChild>
            <Button variant="outline" size="sm">labelled</Button>
          </TooltipTrigger>
          <TooltipContent aria-label="Accessible label">Visible text</TooltipContent>
        </Tooltip>
      </div>
    </TooltipProvider>
  )
}
