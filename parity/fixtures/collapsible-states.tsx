import { Button } from "@/registry/new-york-v4/ui/button"
import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
} from "@/registry/new-york-v4/ui/collapsible"

function Repo({ name }: { name: string }) {
  return (
    <div className="rounded-md border px-4 py-2 font-mono text-sm">{name}</div>
  )
}

// Port fixture: the Radix Collapsible props the upstream examples do not
// use: defaultOpen (with an animated content), disabled, and a trigger
// without asChild.
export default function CollapsibleStates() {
  return (
    <div className="flex w-[350px] flex-col gap-6">
      <Collapsible defaultOpen className="flex flex-col gap-2">
        <CollapsibleTrigger className="rounded-md border px-4 py-2 text-sm">
          Toggle open
        </CollapsibleTrigger>
        <CollapsibleContent className="overflow-hidden data-[state=closed]:animate-collapsible-up data-[state=open]:animate-collapsible-down">
          <Repo name="@radix-ui/colors" />
          <Repo name="@stitches/react" />
        </CollapsibleContent>
      </Collapsible>
      <Collapsible disabled className="flex flex-col gap-2">
        <CollapsibleTrigger asChild>
          <Button variant="outline">Disabled</Button>
        </CollapsibleTrigger>
        <CollapsibleContent>
          <Repo name="@radix-ui/themes" />
        </CollapsibleContent>
      </Collapsible>
    </div>
  )
}
