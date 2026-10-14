// The upstream app's copy of vaul (parity/node_modules has another one,
// whose Handle would not share the shadcn Drawer's context).
import { Drawer as DrawerPrimitive } from "../../ui/apps/v4/node_modules/vaul/dist/index.mjs"

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

// Port fixture: vaul options the shadcn examples do not use: snap points,
// a drawer that cannot be dismissed, background scaling and handle-only
// dragging.
export default function DrawerOptions() {
  return (
    <div data-vaul-drawer-wrapper="" className="bg-background flex flex-wrap items-center gap-2 p-4">
      <Drawer snapPoints={["240px", 1]}>
        <DrawerTrigger asChild>
          <Button variant="outline">Snap</Button>
        </DrawerTrigger>
        <DrawerContent className="h-full">
          <DrawerHeader>
            <DrawerTitle>Snap points</DrawerTitle>
            <DrawerDescription>Opens at 240px, drag up to expand.</DrawerDescription>
          </DrawerHeader>
          <DrawerFooter>
            <DrawerClose asChild>
              <Button variant="outline">Close</Button>
            </DrawerClose>
          </DrawerFooter>
        </DrawerContent>
      </Drawer>
      <Drawer dismissible={false}>
        <DrawerTrigger asChild>
          <Button variant="outline">Locked</Button>
        </DrawerTrigger>
        <DrawerContent>
          <DrawerHeader>
            <DrawerTitle>Not dismissible</DrawerTitle>
            <DrawerDescription>Only the button closes it.</DrawerDescription>
          </DrawerHeader>
          <DrawerFooter>
            <DrawerClose asChild>
              <Button variant="outline">Close</Button>
            </DrawerClose>
          </DrawerFooter>
        </DrawerContent>
      </Drawer>
      <Drawer shouldScaleBackground>
        <DrawerTrigger asChild>
          <Button variant="outline">Scale</Button>
        </DrawerTrigger>
        <DrawerContent>
          <DrawerHeader>
            <DrawerTitle>Scaled background</DrawerTitle>
            <DrawerDescription>The page shrinks behind the drawer.</DrawerDescription>
          </DrawerHeader>
          <DrawerFooter>
            <DrawerClose asChild>
              <Button variant="outline">Close</Button>
            </DrawerClose>
          </DrawerFooter>
        </DrawerContent>
      </Drawer>
      <Drawer handleOnly>
        <DrawerTrigger asChild>
          <Button variant="outline">Handle</Button>
        </DrawerTrigger>
        <DrawerContent>
          <DrawerPrimitive.Handle className="mt-2" />
          <DrawerHeader>
            <DrawerTitle>Handle only</DrawerTitle>
            <DrawerDescription>Drag the handle to close.</DrawerDescription>
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
