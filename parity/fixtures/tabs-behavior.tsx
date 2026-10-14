import {
  Tabs,
  TabsContent,
  TabsList,
  TabsTrigger,
} from "@/registry/new-york-v4/ui/tabs"

// Port fixture: the Radix Tabs props the upstream examples do not use:
// vertical orientation, manual activation, a disabled trigger, the line
// list variant, loop={false} and dir="rtl".
export default function TabsBehavior() {
  return (
    <div className="flex flex-col gap-6">
      <Tabs defaultValue="alpha" orientation="vertical" activationMode="manual">
        <TabsList>
          <TabsTrigger value="alpha">Alpha</TabsTrigger>
          <TabsTrigger value="beta" disabled>
            Beta
          </TabsTrigger>
          <TabsTrigger value="gamma">Gamma</TabsTrigger>
        </TabsList>
        <TabsContent value="alpha">Alpha panel</TabsContent>
        <TabsContent value="beta">Beta panel</TabsContent>
        <TabsContent value="gamma">Gamma panel</TabsContent>
      </Tabs>
      <Tabs defaultValue="one" dir="rtl">
        <TabsList variant="line" loop={false}>
          <TabsTrigger value="one">One</TabsTrigger>
          <TabsTrigger value="two">Two</TabsTrigger>
          <TabsTrigger value="three">Three</TabsTrigger>
        </TabsList>
        <TabsContent value="one">First panel</TabsContent>
        <TabsContent value="two">Second panel</TabsContent>
        <TabsContent value="three">Third panel</TabsContent>
      </Tabs>
    </div>
  )
}
