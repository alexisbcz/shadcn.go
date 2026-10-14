import { Progress } from "@/registry/new-york-v4/ui/progress"

// Port fixture: the Radix Progress props the upstream example does not
// use: no value (indeterminate), a complete bar, a custom max, a
// fractional value, zero, and a complete bar with a custom max.
export default function ProgressStates() {
  return (
    <div className="flex w-[60%] flex-col gap-4">
      <Progress />
      <Progress value={100} />
      <Progress value={30} max={40} />
      <Progress value={33.3} />
      <Progress value={0} />
      <Progress value={40} max={40} className="h-4 bg-muted" />
    </div>
  )
}
