import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectLabel,
  SelectSeparator,
  SelectTrigger,
  SelectValue,
} from "@/registry/new-york-v4/ui/select"

// Port fixture: the popper position (side, align, offsets, collision
// handling, scroll buttons in a popper), the small trigger, a default
// value, disabled items, a separator and an item found by its textValue.
export default function SelectPopper() {
  return (
    <div className="flex items-start gap-6">
      <Select defaultValue="banana">
        <SelectTrigger size="sm" className="w-[160px]">
          <SelectValue placeholder="Fruit" />
        </SelectTrigger>
        <SelectContent position="popper">
          <SelectGroup>
            <SelectLabel>Fruits</SelectLabel>
            <SelectItem value="apple">Apple</SelectItem>
            <SelectItem value="banana">Banana</SelectItem>
            <SelectItem value="blueberry" disabled>
              Blueberry
            </SelectItem>
          </SelectGroup>
          <SelectSeparator />
          <SelectGroup>
            <SelectLabel>Vegetables</SelectLabel>
            <SelectItem value="carrot">Carrot</SelectItem>
            <SelectItem value="potato" textValue="Spud">
              Potato
            </SelectItem>
          </SelectGroup>
        </SelectContent>
      </Select>
      <Select>
        <SelectTrigger className="w-[160px]">
          <SelectValue placeholder="Top, end" />
        </SelectTrigger>
        <SelectContent position="popper" side="top" align="end" sideOffset={8}>
          <SelectItem value="one">One</SelectItem>
          <SelectItem value="two">Two</SelectItem>
          <SelectItem value="three">Three</SelectItem>
        </SelectContent>
      </Select>
      <Select>
        <SelectTrigger className="w-[160px]">
          <SelectValue placeholder="Many" />
        </SelectTrigger>
        <SelectContent position="popper" side="right" alignOffset={4}>
          {Array.from({ length: 30 }, (_, i) => (
            <SelectItem key={i} value={`item-${i + 1}`}>
              {`Item ${i + 1}`}
            </SelectItem>
          ))}
        </SelectContent>
      </Select>
    </div>
  )
}
