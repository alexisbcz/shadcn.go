import { Button } from "@/registry/new-york-v4/ui/button"
import {
  Drawer,
  DrawerClose,
  DrawerContent,
  DrawerDescription,
  DrawerFooter,
  DrawerHeader,
  DrawerTitle,
  DrawerTrigger,
} from "@/registry/new-york-v4/ui/drawer"

// Port fixture: a drawer from each edge and a non-modal drawer.
const directions = ["top", "right", "bottom", "left"] as const

export default function DrawerDirections() {
  return (
    <div className="flex flex-wrap items-center gap-2">
      {directions.map((direction) => (
        <Drawer key={direction} direction={direction}>
          <DrawerTrigger asChild>
            <Button variant="outline" className="capitalize">
              {direction}
            </Button>
          </DrawerTrigger>
          <DrawerContent>
            <DrawerHeader>
              <DrawerTitle>{`Drawer from the ${direction}`}</DrawerTitle>
              <DrawerDescription>Drag it back to its edge to close it.</DrawerDescription>
            </DrawerHeader>
            <div className="px-4 text-sm">Content</div>
            <DrawerFooter>
              <Button>Submit</Button>
              <DrawerClose asChild>
                <Button variant="outline">Cancel</Button>
              </DrawerClose>
            </DrawerFooter>
          </DrawerContent>
        </Drawer>
      ))}
      <Drawer modal={false}>
        <DrawerTrigger asChild>
          <Button variant="outline">Non-modal</Button>
        </DrawerTrigger>
        <DrawerContent>
          <DrawerHeader>
            <DrawerTitle>Non-modal drawer</DrawerTitle>
            <DrawerDescription>The page stays interactive.</DrawerDescription>
          </DrawerHeader>
          <DrawerFooter>
            <DrawerClose asChild>
              <Button variant="outline">Close</Button>
            </DrawerClose>
          </DrawerFooter>
        </DrawerContent>
      </Drawer>
    </div>
  )
}
