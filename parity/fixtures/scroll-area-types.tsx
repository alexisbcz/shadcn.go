import * as React from "react"

import { ScrollArea } from "@/registry/new-york-v4/ui/scroll-area"
import { Separator } from "@/registry/new-york-v4/ui/separator"

function Tags({ n }: { n: number }) {
  const tags = Array.from({ length: n }).map((_, i, a) => `v1.2.0-beta.${a.length - i}`)
  return (
    <div className="p-4">
      <h4 className="mb-4 text-sm leading-none font-medium">Tags</h4>
      {tags.map((tag) => (
        <React.Fragment key={tag}>
          <div className="text-sm">{tag}</div>
          <Separator className="my-2" />
        </React.Fragment>
      ))}
    </div>
  )
}

// Port fixture: the Radix types the examples do not use (scroll with a
// shorter hide delay, auto with and without overflow, always).
export default function ScrollAreaTypes() {
  return (
    <div className="grid grid-cols-2 gap-4">
      <ScrollArea type="scroll" scrollHideDelay={300} className="h-40 w-40 rounded-md border">
        <Tags n={20} />
      </ScrollArea>
      <ScrollArea type="auto" className="h-40 w-40 rounded-md border">
        <Tags n={20} />
      </ScrollArea>
      <ScrollArea type="auto" className="h-40 w-40 rounded-md border">
        <Tags n={2} />
      </ScrollArea>
      <ScrollArea type="always" className="h-40 w-40 rounded-md border">
        <Tags n={20} />
      </ScrollArea>
    </div>
  )
}
