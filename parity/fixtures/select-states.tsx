import { AppleIcon, CherryIcon } from "lucide-react"

import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/registry/new-york-v4/ui/select"

// Port fixture: a disabled select, a disabled trigger, an invalid trigger,
// items with icons (the value shows the item's content), a value with its
// own children, a right-to-left select, and disabled first items.
export default function SelectStates() {
  return (
    <div className="grid grid-cols-2 gap-4">
      <Select disabled defaultValue="apple">
        <SelectTrigger className="w-[180px]">
          <SelectValue placeholder="Disabled" />
        </SelectTrigger>
        <SelectContent>
          <SelectItem value="apple">Apple</SelectItem>
          <SelectItem value="cherry">Cherry</SelectItem>
        </SelectContent>
      </Select>
      <Select>
        <SelectTrigger className="w-[180px]" disabled>
          <SelectValue placeholder="Disabled trigger" />
        </SelectTrigger>
        <SelectContent>
          <SelectItem value="apple">Apple</SelectItem>
        </SelectContent>
      </Select>
      <Select defaultValue="cherry">
        <SelectTrigger className="w-[180px]" aria-invalid>
          <SelectValue placeholder="With icons" />
        </SelectTrigger>
        <SelectContent>
          <SelectItem value="apple">
            <AppleIcon />
            Apple
          </SelectItem>
          <SelectItem value="cherry">
            <CherryIcon />
            Cherry
          </SelectItem>
        </SelectContent>
      </Select>
      <Select defaultValue="b">
        <SelectTrigger className="w-[180px]">
          <SelectValue placeholder="Custom">Custom value</SelectValue>
        </SelectTrigger>
        <SelectContent>
          <SelectItem value="a">Alpha</SelectItem>
          <SelectItem value="b">Beta</SelectItem>
        </SelectContent>
      </Select>
      <Select dir="rtl">
        <SelectTrigger className="w-[180px]">
          <SelectValue placeholder="Right to left" />
        </SelectTrigger>
        <SelectContent>
          <SelectItem value="one">One</SelectItem>
          <SelectItem value="two">Two</SelectItem>
        </SelectContent>
      </Select>
      <Select>
        <SelectTrigger className="w-[180px]">
          <SelectValue placeholder="Disabled first" />
        </SelectTrigger>
        <SelectContent>
          <SelectItem value="one" disabled>
            One
          </SelectItem>
          <SelectItem value="two">Two</SelectItem>
          <SelectItem value="three" disabled>
            Three
          </SelectItem>
          <SelectItem value="four">Four</SelectItem>
        </SelectContent>
      </Select>
    </div>
  )
}
