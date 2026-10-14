import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
} from "@/registry/new-york-v4/ui/collapsible"

// Port fixture: a collapsible nested in another one's content, each with
// its own trigger and state, and a second trigger for the outer one.
export default function CollapsibleNested() {
  return (
    <div className="flex w-[350px] flex-col gap-2">
      <Collapsible className="flex flex-col gap-2">
        <div className="flex gap-2">
          <CollapsibleTrigger className="rounded-md border px-4 py-2 text-sm">Outer</CollapsibleTrigger>
          <CollapsibleTrigger className="rounded-md border px-4 py-2 text-sm">Outer too</CollapsibleTrigger>
        </div>
        <CollapsibleContent className="flex flex-col gap-2">
          <div className="rounded-md border px-4 py-2 text-sm">Outer content</div>
          <Collapsible defaultOpen className="flex flex-col gap-2">
            <CollapsibleTrigger className="rounded-md border px-4 py-2 text-sm">Inner</CollapsibleTrigger>
            <CollapsibleContent>
              <div className="rounded-md border px-4 py-2 text-sm">Inner content</div>
            </CollapsibleContent>
          </Collapsible>
        </CollapsibleContent>
      </Collapsible>
    </div>
  )
}
