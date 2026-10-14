import { Button } from "@/registry/new-york-v4/ui/button"
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/registry/new-york-v4/ui/popover"

// Port fixture: collision handling. Triggers sit at the viewport edges so
// the content flips (top -> bottom, bottom -> top, left -> right) and
// shifts (sticky partial and always), or overflows (avoidCollisions
// false); collisionPadding grows the margin, and hideWhenDetached hides
// the content once its trigger scrolls out of its container.
export default function PopoverCollision() {
  return (
    <div>
      <div className="fixed top-2 left-2">
        <Popover>
          <PopoverTrigger asChild>
            <Button variant="outline" size="sm">flip top</Button>
          </PopoverTrigger>
          <PopoverContent side="top" align="end">Flipped to the bottom, shifted right.</PopoverContent>
        </Popover>
      </div>
      <div className="fixed right-2 bottom-2">
        <Popover>
          <PopoverTrigger asChild>
            <Button variant="outline" size="sm">flip bottom</Button>
          </PopoverTrigger>
          <PopoverContent collisionPadding={24}>Flipped to the top, padded.</PopoverContent>
        </Popover>
      </div>
      <div className="fixed top-40 left-2">
        <Popover>
          <PopoverTrigger asChild>
            <Button variant="outline" size="sm">flip left</Button>
          </PopoverTrigger>
          <PopoverContent side="left" collisionPadding={{ top: 40, left: 8 }}>Flipped to the right.</PopoverContent>
        </Popover>
      </div>
      <div className="fixed top-2 right-2">
        <Popover>
          <PopoverTrigger asChild>
            <Button variant="outline" size="sm">no avoid</Button>
          </PopoverTrigger>
          <PopoverContent avoidCollisions={false} side="top">Overflows the viewport.</PopoverContent>
        </Popover>
      </div>
      <div className="fixed bottom-2 left-2">
        <Popover>
          <PopoverTrigger asChild>
            <Button variant="outline" size="sm">sticky always</Button>
          </PopoverTrigger>
          <PopoverContent side="right" align="end" sticky="always">Always inside.</PopoverContent>
        </Popover>
      </div>
      <div className="h-32 w-64 overflow-auto rounded-md border p-4" data-testid="scroller">
        <div className="h-96 pt-8">
          <Popover>
            <PopoverTrigger asChild>
              <Button variant="outline" size="sm">detached</Button>
            </PopoverTrigger>
            <PopoverContent hideWhenDetached className="w-40" side="right">Hidden when detached.</PopoverContent>
          </Popover>
        </div>
      </div>
    </div>
  )
}
