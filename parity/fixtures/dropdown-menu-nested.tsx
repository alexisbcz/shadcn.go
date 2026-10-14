"use client"

import * as React from "../../ui/apps/v4/node_modules/react"

import { Button } from "@/registry/new-york-v4/ui/button"
import {
  DropdownMenu,
  DropdownMenuCheckboxItem,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuPortal,
  DropdownMenuRadioGroup,
  DropdownMenuRadioItem,
  DropdownMenuSeparator,
  DropdownMenuShortcut,
  DropdownMenuSub,
  DropdownMenuSubContent,
  DropdownMenuSubTrigger,
  DropdownMenuTrigger,
} from "@/registry/new-york-v4/ui/dropdown-menu"

// Port fixture: three levels of submenus (in place and portalled), loop,
// disabled items and sub triggers, an indeterminate checkbox, a disabled
// radio item, textValue, an asChild link item, inset label, shortcuts and
// the destructive variant.
export default function DropdownMenuNested() {
  const [bold, setBold] = React.useState<boolean | "indeterminate">("indeterminate")
  const [italic, setItalic] = React.useState(false)
  const [size, setSize] = React.useState("m")
  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="outline">Nested</Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent className="w-56" loop>
        <DropdownMenuLabel inset>Format</DropdownMenuLabel>
        <DropdownMenuCheckboxItem checked={bold} onCheckedChange={setBold}>
          Bold
          <DropdownMenuShortcut>⌘B</DropdownMenuShortcut>
        </DropdownMenuCheckboxItem>
        <DropdownMenuCheckboxItem checked={italic} onCheckedChange={setItalic}>
          Italic
        </DropdownMenuCheckboxItem>
        <DropdownMenuSeparator />
        <DropdownMenuSub>
          <DropdownMenuSubTrigger inset>Level 1</DropdownMenuSubTrigger>
          <DropdownMenuSubContent loop>
            <DropdownMenuItem>Alpha</DropdownMenuItem>
            <DropdownMenuSub>
              <DropdownMenuSubTrigger>Level 2</DropdownMenuSubTrigger>
              <DropdownMenuPortal>
                <DropdownMenuSubContent>
                  <DropdownMenuRadioGroup value={size} onValueChange={setSize}>
                    <DropdownMenuRadioItem value="s">Small</DropdownMenuRadioItem>
                    <DropdownMenuRadioItem value="m">Medium</DropdownMenuRadioItem>
                    <DropdownMenuRadioItem value="l" disabled>
                      Large
                    </DropdownMenuRadioItem>
                  </DropdownMenuRadioGroup>
                  <DropdownMenuSeparator />
                  <DropdownMenuSub>
                    <DropdownMenuSubTrigger>Level 3</DropdownMenuSubTrigger>
                    <DropdownMenuSubContent>
                      <DropdownMenuItem>Deep item</DropdownMenuItem>
                      <DropdownMenuItem variant="destructive">
                        Delete
                        <DropdownMenuShortcut>⌘⌫</DropdownMenuShortcut>
                      </DropdownMenuItem>
                    </DropdownMenuSubContent>
                  </DropdownMenuSub>
                </DropdownMenuSubContent>
              </DropdownMenuPortal>
            </DropdownMenuSub>
            <DropdownMenuItem disabled>Disabled</DropdownMenuItem>
          </DropdownMenuSubContent>
        </DropdownMenuSub>
        <DropdownMenuSub>
          <DropdownMenuSubTrigger disabled>Disabled sub</DropdownMenuSubTrigger>
          <DropdownMenuSubContent>
            <DropdownMenuItem>Never shown</DropdownMenuItem>
          </DropdownMenuSubContent>
        </DropdownMenuSub>
        <DropdownMenuItem textValue="Zebra">Custom text value</DropdownMenuItem>
        <DropdownMenuItem asChild>
          <a href="#docs">Docs link</a>
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  )
}
