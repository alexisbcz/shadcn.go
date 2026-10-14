import { Slider } from "@/registry/new-york-v4/ui/slider"

// Port fixture: sliders tied to a form through the form attribute (outside
// the form: Radix keeps the bubble inputs), right to left combined with
// inverted (slides from the left again), and a right-to-left range.
export default function SliderForm() {
  return (
    <div className="flex w-[60%] flex-col gap-6">
      <form id="slider-form" />
      <Slider defaultValue={[40]} name="volume" form="slider-form" />
      <Slider defaultValue={[20, 70]} name="price" form="slider-form" />
      <Slider defaultValue={[30]} dir="rtl" inverted />
      <Slider defaultValue={[10, 90]} dir="rtl" />
    </div>
  )
}
