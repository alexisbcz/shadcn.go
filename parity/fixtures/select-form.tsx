import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/registry/new-york-v4/ui/select"

// Port fixture: native form integration. Inside a form, Radix renders a
// hidden native <select> (name, required, autocomplete, disabled) whose
// value follows the selection; outside a form it does too when the form
// prop names one.
export default function SelectForm() {
  return (
    <div className="flex flex-col gap-4">
      <form id="select-form" className="flex gap-3">
        <Select name="fruit" required autoComplete="off">
          <SelectTrigger className="w-[180px]">
            <SelectValue placeholder="Fruit" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="apple">Apple</SelectItem>
            <SelectItem value="banana">Banana</SelectItem>
            <SelectItem value="cherry" disabled>
              Cherry
            </SelectItem>
          </SelectContent>
        </Select>
        <Select name="size" defaultValue="m" disabled>
          <SelectTrigger className="w-[120px]">
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="s">Small</SelectItem>
            <SelectItem value="m">Medium</SelectItem>
          </SelectContent>
        </Select>
      </form>
      <Select name="color" form="select-form" defaultValue="red">
        <SelectTrigger className="w-[180px]">
          <SelectValue placeholder="Color" />
        </SelectTrigger>
        <SelectContent>
          <SelectItem value="red">Red</SelectItem>
          <SelectItem value="green">Green</SelectItem>
        </SelectContent>
      </Select>
    </div>
  )
}
