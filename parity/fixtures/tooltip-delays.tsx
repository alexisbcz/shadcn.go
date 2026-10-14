import { Button } from "@/registry/new-york-v4/ui/button"
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/registry/new-york-v4/ui/tooltip"

// Port fixture: TooltipProvider delays. Tooltips in one provider share the
// open delay (600ms) and skip it for 400ms after one closed; C overrides
// the delay (0), D disables hoverable content. The triggers sit at fixed
// positions so scenarios can move the pointer between them.
export default function TooltipDelays() {
  return (
    <TooltipProvider delayDuration={600} skipDelayDuration={400}>
      <div className="fixed top-[200px] left-[200px] flex gap-8">
        <Tooltip>
          <TooltipTrigger asChild>
            <Button variant="outline" className="w-32">A</Button>
          </TooltipTrigger>
          <TooltipContent>Tooltip A</TooltipContent>
        </Tooltip>
        <Tooltip>
          <TooltipTrigger asChild>
            <Button variant="outline" className="w-32">B</Button>
          </TooltipTrigger>
          <TooltipContent side="bottom">Tooltip B</TooltipContent>
        </Tooltip>
        <Tooltip delayDuration={0}>
          <TooltipTrigger asChild>
            <Button variant="outline" className="w-32">C</Button>
          </TooltipTrigger>
          <TooltipContent>Tooltip C</TooltipContent>
        </Tooltip>
        <Tooltip disableHoverableContent>
          <TooltipTrigger asChild>
            <Button variant="outline" className="w-32">D</Button>
          </TooltipTrigger>
          <TooltipContent>Tooltip D</TooltipContent>
        </Tooltip>
      </div>
    </TooltipProvider>
  )
}
