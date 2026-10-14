import { BoldIcon, ItalicIcon, UnderlineIcon } from "lucide-react"

import { DirectionProvider } from "@/registry/new-york-v4/ui/direction"
import { Label } from "@/registry/new-york-v4/ui/label"
import {
  RadioGroup,
  RadioGroupItem,
} from "@/registry/new-york-v4/ui/radio-group"
import { Slider } from "@/registry/new-york-v4/ui/slider"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/registry/new-york-v4/ui/tabs"
import {
  ToggleGroup,
  ToggleGroupItem,
} from "@/registry/new-york-v4/ui/toggle-group"

function Items() {
  return (
    <>
      <ToggleGroupItem value="bold" aria-label="Toggle bold">
        <BoldIcon />
      </ToggleGroupItem>
      <ToggleGroupItem value="italic" aria-label="Toggle italic">
        <ItalicIcon />
      </ToggleGroupItem>
      <ToggleGroupItem value="underline" aria-label="Toggle underline">
        <UnderlineIcon />
      </ToggleGroupItem>
    </>
  )
}

// Port fixture: DirectionProvider (dir="rtl") around direction-aware
// components (Tabs, RadioGroup, ToggleGroup, Slider): they render
// dir="rtl", lay out and handle arrow keys right to left. A component's
// own dir prop wins over the provider, and a nested provider given the
// direction prop (alias of dir) switches its subtree back to ltr.
export default function DirectionProviderFixture() {
  return (
    <DirectionProvider dir="rtl">
      <div className="flex w-full max-w-md flex-col gap-6">
        <Tabs defaultValue="one">
          <TabsList>
            <TabsTrigger value="one">One</TabsTrigger>
            <TabsTrigger value="two">Two</TabsTrigger>
            <TabsTrigger value="three">Three</TabsTrigger>
          </TabsList>
          <TabsContent value="one">First panel.</TabsContent>
          <TabsContent value="two">Second panel.</TabsContent>
          <TabsContent value="three">Third panel.</TabsContent>
        </Tabs>
        <RadioGroup defaultValue="a" className="flex">
          <div className="flex items-center gap-2">
            <RadioGroupItem value="a" id="dir-a" />
            <Label htmlFor="dir-a">First</Label>
          </div>
          <div className="flex items-center gap-2">
            <RadioGroupItem value="b" id="dir-b" />
            <Label htmlFor="dir-b">Second</Label>
          </div>
        </RadioGroup>
        <ToggleGroup type="multiple" variant="outline">
          <Items />
        </ToggleGroup>
        <Slider defaultValue={[30]} max={100} step={1} />
        <ToggleGroup type="single" dir="ltr">
          <Items />
        </ToggleGroup>
        <DirectionProvider direction="ltr">
          <Slider defaultValue={[30]} max={100} step={1} />
        </DirectionProvider>
      </div>
    </DirectionProvider>
  )
}
