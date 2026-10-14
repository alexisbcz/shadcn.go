import {
  ResizableHandle,
  ResizablePanel,
  ResizablePanelGroup,
} from "@/registry/new-york-v4/ui/resizable"

// Port fixture: panel constraints the upstream examples do not use: a
// collapsible panel with min, max and collapsed sizes, and a minimum size
// on its neighbour (keyboard collapse/expand with Enter, clamping).
export default function ResizableCollapsible() {
  return (
    <ResizablePanelGroup
      orientation="horizontal"
      className="min-h-[200px] max-w-md rounded-lg border md:min-w-[450px]"
    >
      <ResizablePanel
        defaultSize="25%"
        minSize="15%"
        maxSize="40%"
        collapsible
        collapsedSize="5%"
      >
        <div className="flex h-full items-center justify-center p-6">
          <span className="font-semibold">Sidebar</span>
        </div>
      </ResizablePanel>
      <ResizableHandle withHandle />
      <ResizablePanel defaultSize="75%" minSize="30%">
        <div className="flex h-full items-center justify-center p-6">
          <span className="font-semibold">Content</span>
        </div>
      </ResizablePanel>
    </ResizablePanelGroup>
  )
}
