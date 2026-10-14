import { Label } from "@/registry/new-york-v4/ui/label"
import { Switch } from "@/registry/new-york-v4/ui/switch"

// Port fixture: the Switch props the upstream example does not use: the sm
// size, defaultChecked, disabled (checked and unchecked), aria-invalid, and
// switches in a form (hidden bubble inputs) with name, value, required and
// a reset button.
export default function SwitchStates() {
  return (
    <div className="flex flex-col gap-4">
      <div className="flex items-center gap-2">
        <Switch size="sm" id="small" />
        <Label htmlFor="small">Small</Label>
      </div>
      <div className="flex items-center gap-2">
        <Switch size="sm" defaultChecked id="small-on" />
        <Label htmlFor="small-on">Small on</Label>
      </div>
      <div className="flex items-center gap-2">
        <Switch defaultChecked id="on" />
        <Label htmlFor="on">On</Label>
      </div>
      <div className="flex items-center gap-2">
        <Switch disabled id="off-disabled" />
        <Label htmlFor="off-disabled">Disabled</Label>
      </div>
      <div className="flex items-center gap-2">
        <Switch defaultChecked disabled id="on-disabled" />
        <Label htmlFor="on-disabled">Disabled on</Label>
      </div>
      <div className="flex items-center gap-2">
        <Switch aria-invalid="true" id="invalid" />
        <Label htmlFor="invalid">Invalid</Label>
      </div>
      <form className="flex flex-col gap-3">
        <div className="flex items-center gap-2">
          <Switch name="wifi" defaultChecked id="wifi" />
          <Label htmlFor="wifi">Wi-Fi</Label>
        </div>
        <div className="flex items-center gap-2">
          <Switch name="bt" value="yes" required id="bt" />
          <Label htmlFor="bt">Bluetooth</Label>
        </div>
        <button type="reset" className="w-fit rounded-md border px-2 text-sm">
          Reset
        </button>
      </form>
    </div>
  )
}
