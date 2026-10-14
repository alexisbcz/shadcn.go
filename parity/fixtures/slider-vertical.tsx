import { Slider } from "@/registry/new-york-v4/ui/slider"

// Port fixture: vertical sliders, plain, range, inverted and disabled.
export default function SliderVertical() {
  return (
    <div className="flex h-48 gap-10">
      <Slider defaultValue={[50]} orientation="vertical" />
      <Slider defaultValue={[20, 70]} orientation="vertical" />
      <Slider defaultValue={[30]} orientation="vertical" inverted />
      <Slider defaultValue={[40]} orientation="vertical" disabled />
    </div>
  )
}
