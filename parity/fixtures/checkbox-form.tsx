import { Checkbox } from "@/registry/new-york-v4/ui/checkbox"
import { Label } from "@/registry/new-york-v4/ui/label"

// Port fixture: checkboxes tied to a form through the form attribute
// (outside the form: Radix keeps the bubble input but registers no reset
// listener, since the button itself has no form), and a disabled checkbox
// inside a form (its input is disabled, so it is not submitted).
export default function CheckboxForm() {
  return (
    <div className="flex flex-col gap-4">
      <form id="checkbox-form" className="flex gap-3">
        <button type="reset" className="w-fit rounded-md border px-2 text-sm">
          Reset
        </button>
      </form>
      <div className="flex items-center gap-3">
        <Checkbox form="checkbox-form" name="remote" value="yes" defaultChecked id="remote" />
        <Label htmlFor="remote">Remote</Label>
      </div>
      <div className="flex items-center gap-3">
        <Checkbox form="checkbox-form" name="remote-mixed" defaultChecked="indeterminate" required id="remote-mixed" />
        <Label htmlFor="remote-mixed">Remote mixed</Label>
      </div>
      <form className="flex items-center gap-3">
        <Checkbox name="locked" disabled defaultChecked id="locked" />
        <Label htmlFor="locked">Locked</Label>
      </form>
    </div>
  )
}
