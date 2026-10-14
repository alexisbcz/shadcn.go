import { Separator } from "@/registry/new-york-v4/ui/separator"

// Port fixture: non-decorative separators (role="separator", and
// aria-orientation when vertical) in both orientations.
export default function SeparatorSemantic() {
  return (
    <div className="w-64">
      <div className="text-sm">Above</div>
      <Separator decorative={false} className="my-4" />
      <div className="flex h-5 items-center gap-4 text-sm">
        <div>Left</div>
        <Separator decorative={false} orientation="vertical" />
        <div>Right</div>
      </div>
    </div>
  )
}
