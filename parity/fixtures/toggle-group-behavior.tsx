import { Bold, Italic, Underline } from "lucide-react"

import {
  ToggleGroup,
  ToggleGroupItem,
} from "@/registry/new-york-v4/ui/toggle-group"

// Port fixture: the Radix ToggleGroup props the upstream examples do not
// use: defaultValue, orientation="vertical", loop={false}, a disabled item,
// dir="rtl" and rovingFocus={false}.
export default function ToggleGroupBehavior() {
  return (
    <div className="flex flex-col gap-6">
      <ToggleGroup
        type="multiple"
        variant="outline"
        defaultValue={["italic"]}
        orientation="vertical"
        loop={false}
        className="flex-col"
      >
        <ToggleGroupItem value="bold" aria-label="Toggle bold">
          <Bold className="h-4 w-4" />
        </ToggleGroupItem>
        <ToggleGroupItem value="italic" aria-label="Toggle italic">
          <Italic className="h-4 w-4" />
        </ToggleGroupItem>
        <ToggleGroupItem value="underline" aria-label="Toggle underline" disabled>
          <Underline className="h-4 w-4" />
        </ToggleGroupItem>
      </ToggleGroup>
      <ToggleGroup type="single" defaultValue="center" dir="rtl" spacing={1}>
        <ToggleGroupItem value="left">Left</ToggleGroupItem>
        <ToggleGroupItem value="center">Center</ToggleGroupItem>
        <ToggleGroupItem value="right">Right</ToggleGroupItem>
      </ToggleGroup>
      <ToggleGroup type="single" rovingFocus={false}>
        <ToggleGroupItem value="a">A</ToggleGroupItem>
        <ToggleGroupItem value="b">B</ToggleGroupItem>
      </ToggleGroup>
    </div>
  )
}
