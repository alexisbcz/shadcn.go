import { Button } from "@/registry/new-york-v4/ui/button"
import {
  Popover,
  PopoverContent,
  PopoverDescription,
  PopoverHeader,
  PopoverTitle,
  PopoverTrigger,
} from "@/registry/new-york-v4/ui/popover"

// Port fixture: every side and alignment of PopoverContent, plus side and
// align offsets, with the header parts.
const variants = [
  { side: "top", align: "start" },
  { side: "top", align: "center" },
  { side: "top", align: "end" },
  { side: "right", align: "start" },
  { side: "right", align: "center", sideOffset: 16 },
  { side: "right", align: "end" },
  { side: "bottom", align: "start", alignOffset: 12 },
  { side: "bottom", align: "center", sideOffset: 0 },
  { side: "bottom", align: "end", alignOffset: -8 },
  { side: "left", align: "start", alignOffset: 6 },
  { side: "left", align: "center" },
  { side: "left", align: "end", sideOffset: 12 },
] as const

export default function PopoverSides() {
  return (
    <div className="grid grid-cols-3 gap-x-24 gap-y-6">
      {variants.map((v, i) => (
        <Popover key={i}>
          <PopoverTrigger asChild>
            <Button variant="outline" size="sm">
              {`${v.side} ${v.align}`}
            </Button>
          </PopoverTrigger>
          <PopoverContent
            className="w-48"
            side={v.side}
            align={v.align}
            {...("sideOffset" in v ? { sideOffset: v.sideOffset } : {})}
            {...("alignOffset" in v ? { alignOffset: v.alignOffset } : {})}
          >
            <PopoverHeader>
              <PopoverTitle>{`${v.side} / ${v.align}`}</PopoverTitle>
              <PopoverDescription>Placed by the popper.</PopoverDescription>
            </PopoverHeader>
          </PopoverContent>
        </Popover>
      ))}
    </div>
  )
}
