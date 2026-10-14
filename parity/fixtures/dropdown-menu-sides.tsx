import { Button } from "@/registry/new-york-v4/ui/button"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuTrigger,
} from "@/registry/new-york-v4/ui/dropdown-menu"

// Port fixture: sides, alignments and offsets of DropdownMenuContent; the
// menus of the middle column are not modal.
const variants = [
  { side: "top", align: "start" },
  { side: "top", align: "center", modal: false },
  { side: "top", align: "end", sideOffset: 12 },
  { side: "right", align: "start", alignOffset: 8 },
  { side: "right", align: "center", modal: false },
  { side: "right", align: "end" },
  { side: "bottom", align: "start", sideOffset: 0 },
  { side: "bottom", align: "center", modal: false },
  { side: "bottom", align: "end", alignOffset: -6 },
  { side: "left", align: "start" },
  { side: "left", align: "center", modal: false, sideOffset: 10 },
  { side: "left", align: "end" },
] as const

export default function DropdownMenuSides() {
  return (
    <div className="grid grid-cols-3 gap-x-24 gap-y-6">
      {variants.map((v, i) => (
        <DropdownMenu key={i} modal={"modal" in v ? v.modal : true}>
          <DropdownMenuTrigger asChild>
            <Button variant="outline" size="sm">
              {`${v.side} ${v.align}`}
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent
            className="w-40"
            side={v.side}
            align={v.align}
            {...("sideOffset" in v ? { sideOffset: v.sideOffset } : {})}
            {...("alignOffset" in v ? { alignOffset: v.alignOffset } : {})}
          >
            <DropdownMenuLabel>{`${v.side} / ${v.align}`}</DropdownMenuLabel>
            <DropdownMenuItem>First</DropdownMenuItem>
            <DropdownMenuItem>Second</DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      ))}
    </div>
  )
}
