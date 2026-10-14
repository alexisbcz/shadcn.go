import { Checkbox } from "@/registry/new-york-v4/ui/checkbox"
import { Label } from "@/registry/new-york-v4/ui/label"

// Port fixture: the Radix Checkbox props the upstream examples do not use:
// an indeterminate checkbox, checkboxes in a form (hidden bubble inputs)
// with name, value, required and a reset button, a checked disabled
// checkbox and an aria-invalid one.
export default function CheckboxStates() {
  return (
    <div className="flex flex-col gap-6">
      <div className="flex items-center gap-3">
        <Checkbox defaultChecked="indeterminate" id="mixed" />
        <Label htmlFor="mixed">Select all</Label>
      </div>
      <form className="flex flex-col gap-3">
        <div className="flex items-center gap-3">
          <Checkbox name="news" value="yes" required defaultChecked id="news" />
          <Label htmlFor="news">Newsletter</Label>
        </div>
        <div className="flex items-center gap-3">
          <Checkbox name="partial" defaultChecked="indeterminate" id="partial" />
          <Label htmlFor="partial">Partial</Label>
        </div>
        <div className="flex items-center gap-3">
          <Checkbox name="offers" id="offers" />
          <Label htmlFor="offers">Offers</Label>
        </div>
        <button type="reset" className="w-fit rounded-md border px-2 text-sm">
          Reset
        </button>
      </form>
      <div className="flex items-center gap-3">
        <Checkbox defaultChecked disabled id="locked" />
        <Label htmlFor="locked">Locked</Label>
      </div>
      <div className="flex items-center gap-3">
        <Checkbox aria-invalid="true" id="invalid" />
        <Label htmlFor="invalid">Invalid</Label>
      </div>
    </div>
  )
}
