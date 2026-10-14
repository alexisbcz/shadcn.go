import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/registry/new-york-v4/ui/select"

// Port fixture: a select open on load (defaultOpen) with a default value.
export default function SelectDefaultOpen() {
  return (
    <Select defaultOpen defaultValue="two">
      <SelectTrigger className="w-[180px]">
        <SelectValue placeholder="Pick one" />
      </SelectTrigger>
      <SelectContent>
        <SelectItem value="one">One</SelectItem>
        <SelectItem value="two">Two</SelectItem>
        <SelectItem value="three">Three</SelectItem>
      </SelectContent>
    </Select>
  )
}
