import { Label } from "@/registry/new-york-v4/ui/label"
import {
  RadioGroup,
  RadioGroupItem,
} from "@/registry/new-york-v4/ui/radio-group"

// Port fixture: the Radix RadioGroup props the upstream examples do not
// use: a group inside a form (hidden bubble inputs) with name, required,
// orientation="horizontal" and a disabled item; a group with no default
// value, dir="rtl" and loop={false}; a disabled group.
export default function RadioGroupBehavior() {
  return (
    <div className="flex flex-col gap-6">
      <form onSubmit={(event) => event.preventDefault()}>
        <RadioGroup
          name="plan"
          defaultValue="pro"
          orientation="horizontal"
          required
          className="flex"
        >
          <div className="flex items-center gap-2">
            <RadioGroupItem value="free" id="plan-free" />
            <Label htmlFor="plan-free">Free</Label>
          </div>
          <div className="flex items-center gap-2">
            <RadioGroupItem value="pro" id="plan-pro" />
            <Label htmlFor="plan-pro">Pro</Label>
          </div>
          <div className="flex items-center gap-2">
            <RadioGroupItem value="team" id="plan-team" disabled />
            <Label htmlFor="plan-team">Team</Label>
          </div>
        </RadioGroup>
      </form>
      <RadioGroup dir="rtl" loop={false}>
        <div className="flex items-center gap-2">
          <RadioGroupItem value="a" id="rtl-a" />
          <Label htmlFor="rtl-a">First</Label>
        </div>
        <div className="flex items-center gap-2">
          <RadioGroupItem value="b" id="rtl-b" />
          <Label htmlFor="rtl-b">Second</Label>
        </div>
      </RadioGroup>
      <RadioGroup defaultValue="x" disabled>
        <div className="flex items-center gap-2">
          <RadioGroupItem value="x" id="dis-x" />
          <Label htmlFor="dis-x">Locked</Label>
        </div>
      </RadioGroup>
    </div>
  )
}
