"use client"

import * as React from "../../ui/apps/v4/node_modules/react"

import {
  ContextMenu,
  ContextMenuCheckboxItem,
  ContextMenuContent,
  ContextMenuItem,
  ContextMenuLabel,
  ContextMenuPortal,
  ContextMenuSeparator,
  ContextMenuSub,
  ContextMenuSubContent,
  ContextMenuSubTrigger,
  ContextMenuTrigger,
} from "@/registry/new-york-v4/ui/context-menu"

// Port fixture: a non-modal context menu with loop, a portalled nested
// submenu and a checkbox with state; a disabled trigger area.
export default function ContextMenuOptions() {
  const [grid, setGrid] = React.useState(false)
  return (
    <div className="flex flex-col gap-6">
      <ContextMenu modal={false}>
        <ContextMenuTrigger className="flex h-[120px] w-[300px] items-center justify-center rounded-md border border-dashed text-sm">
          Non-modal area
        </ContextMenuTrigger>
        <ContextMenuContent className="w-48" loop>
          <ContextMenuLabel>View</ContextMenuLabel>
          <ContextMenuCheckboxItem checked={grid} onCheckedChange={setGrid}>
            Show grid
          </ContextMenuCheckboxItem>
          <ContextMenuSeparator />
          <ContextMenuSub>
            <ContextMenuSubTrigger>Arrange</ContextMenuSubTrigger>
            <ContextMenuPortal>
              <ContextMenuSubContent className="w-40">
                <ContextMenuItem>Bring forward</ContextMenuItem>
                <ContextMenuSub>
                  <ContextMenuSubTrigger>Align</ContextMenuSubTrigger>
                  <ContextMenuSubContent>
                    <ContextMenuItem>Left</ContextMenuItem>
                    <ContextMenuItem>Right</ContextMenuItem>
                  </ContextMenuSubContent>
                </ContextMenuSub>
              </ContextMenuSubContent>
            </ContextMenuPortal>
          </ContextMenuSub>
          <ContextMenuItem disabled>Disabled</ContextMenuItem>
        </ContextMenuContent>
      </ContextMenu>
      <ContextMenu>
        <ContextMenuTrigger
          disabled
          className="flex h-[80px] w-[300px] items-center justify-center rounded-md border border-dashed text-sm"
        >
          Disabled area
        </ContextMenuTrigger>
        <ContextMenuContent>
          <ContextMenuItem>Never shown</ContextMenuItem>
        </ContextMenuContent>
      </ContextMenu>
    </div>
  )
}
