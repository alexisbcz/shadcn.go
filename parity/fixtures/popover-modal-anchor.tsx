// The upstream app's own radix-ui (a bare "radix-ui" would resolve from
// parity/node_modules, a second React tree).
import { Popover as PopoverPrimitive } from "@/node_modules/radix-ui"

import { Button } from "@/registry/new-york-v4/ui/button"
import { Input } from "@/registry/new-york-v4/ui/input"
import {
  Popover,
  PopoverAnchor,
  PopoverContent,
  PopoverTrigger,
} from "@/registry/new-york-v4/ui/popover"

// Port fixture: a modal popover (scroll lock, outside pointer events
// blocked, focus trapped, rest of the page hidden) positioned against a
// custom PopoverAnchor, with an arrow and a close button (Radix parts the
// shadcn component does not wrap). Radix only honors a custom anchor
// rendered after the trigger: the first one (before its trigger) leaves the
// content positioned against a detached node, like upstream.
export default function PopoverModalAnchor() {
  return (
    <div className="flex flex-col items-center gap-6">
      <Popover modal>
        <PopoverAnchor className="w-72 rounded-md border border-dashed p-2 text-center text-sm">
          Anchor
        </PopoverAnchor>
        <PopoverTrigger asChild>
          <Button variant="outline">Open modal</Button>
        </PopoverTrigger>
        <PopoverContent className="grid gap-2" side="top">
          <Input defaultValue="First" />
          <Input defaultValue="Second" />
          <PopoverPrimitive.Close asChild>
            <Button size="sm">Done</Button>
          </PopoverPrimitive.Close>
          <PopoverPrimitive.Arrow className="fill-popover" />
        </PopoverContent>
      </Popover>
      <Popover>
        <PopoverTrigger asChild>
          <Button variant="outline">Anchored below</Button>
        </PopoverTrigger>
        <PopoverAnchor className="w-72 rounded-md border border-dashed p-2 text-center text-sm">
          Second anchor
        </PopoverAnchor>
        <PopoverContent align="start">Positioned against the anchor after the trigger.</PopoverContent>
      </Popover>
      <Button variant="secondary">Outside</Button>
    </div>
  )
}
