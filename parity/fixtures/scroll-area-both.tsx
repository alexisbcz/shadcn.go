import { ScrollArea, ScrollBar } from "@/registry/new-york-v4/ui/scroll-area"

function Grid() {
  return (
    <div className="grid w-[480px] grid-cols-6 gap-2 p-4">
      {Array.from({ length: 48 }).map((_, i) => (
        <div key={i} className="flex h-12 items-center justify-center rounded-md border text-sm">
          {i + 1}
        </div>
      ))}
    </div>
  )
}

// Port fixture: both scrollbars (and the corner) with type="always", left
// to right and right to left.
export default function ScrollAreaBoth() {
  return (
    <div className="flex gap-6">
      <ScrollArea type="always" className="h-56 w-64 rounded-md border">
        <Grid />
        <ScrollBar orientation="horizontal" />
      </ScrollArea>
      <ScrollArea type="always" dir="rtl" className="h-56 w-64 rounded-md border">
        <Grid />
        <ScrollBar orientation="horizontal" />
      </ScrollArea>
    </div>
  )
}
