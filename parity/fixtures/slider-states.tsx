import { Slider } from "@/registry/new-york-v4/ui/slider"

// Port fixture: the Radix Slider props the upstream example does not use:
// a range (two thumbs) with minStepsBetweenThumbs, three thumbs, a
// fractional step with custom min and max, a coarse step, inverted, right
// to left, disabled, no default value (shadcn's hidden second thumb), and
// a slider in a form (hidden bubble inputs).
export default function SliderStates() {
  return (
    <div className="flex w-[60%] flex-col gap-6">
      <Slider defaultValue={[25, 75]} minStepsBetweenThumbs={10} />
      <Slider defaultValue={[10, 40, 80]} />
      <Slider defaultValue={[0.5]} min={-1} max={2} step={0.1} />
      <Slider defaultValue={[20]} step={10} />
      <Slider defaultValue={[30]} inverted />
      <Slider defaultValue={[30]} dir="rtl" />
      <Slider defaultValue={[60]} disabled />
      <Slider />
      <form>
        <Slider defaultValue={[20, 60]} name="price" />
      </form>
    </div>
  )
}
