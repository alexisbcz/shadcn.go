import { Bold, Italic } from "lucide-react"

import { Toggle } from "@/registry/new-york-v4/ui/toggle"

// Port fixture: Toggle states the upstream examples do not use:
// defaultPressed (default and outline variants) and a pressed disabled
// toggle.
export default function ToggleStates() {
  return (
    <div className="flex items-center gap-2">
      <Toggle defaultPressed aria-label="Toggle bold">
        <Bold />
      </Toggle>
      <Toggle variant="outline" defaultPressed aria-label="Toggle italic">
        <Italic />
        Italic
      </Toggle>
      <Toggle defaultPressed disabled aria-label="Toggle locked">
        <Bold />
      </Toggle>
    </div>
  )
}
