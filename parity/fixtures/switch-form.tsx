import { Label } from "@/registry/new-york-v4/ui/label"
import { Switch } from "@/registry/new-york-v4/ui/switch"

// Port fixture: switches tied to a form through the form attribute (outside
// the form), and a disabled switch inside a form (its input is disabled, so
// it is not submitted).
export default function SwitchForm() {
  return (
    <div className="flex flex-col gap-4">
      <form id="switch-form" className="flex gap-3">
        <button type="reset" className="w-fit rounded-md border px-2 text-sm">
          Reset
        </button>
      </form>
      <div className="flex items-center gap-2">
        <Switch form="switch-form" name="remote" value="yes" defaultChecked id="remote" />
        <Label htmlFor="remote">Remote</Label>
      </div>
      <div className="flex items-center gap-2">
        <Switch form="switch-form" name="remote-off" required id="remote-off" />
        <Label htmlFor="remote-off">Remote off</Label>
      </div>
      <form className="flex items-center gap-2">
        <Switch name="locked" disabled defaultChecked id="locked" />
        <Label htmlFor="locked">Locked</Label>
      </form>
    </div>
  )
}
