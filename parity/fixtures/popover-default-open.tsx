import { Button } from "@/registry/new-york-v4/ui/button"
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/registry/new-york-v4/ui/popover"

// Port fixture: a popover open on load (defaultOpen).
export default function PopoverDefaultOpen() {
  return (
    <Popover defaultOpen>
      <PopoverTrigger asChild>
        <Button variant="outline">Open on load</Button>
      </PopoverTrigger>
      <PopoverContent align="start">
        Opened by defaultOpen. <a href="#x" className="underline">Link</a>{" "}
        <button className="underline">Button</button>
      </PopoverContent>
    </Popover>
  )
}
