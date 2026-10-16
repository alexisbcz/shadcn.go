// bubble-input.js: the hidden native inputs Radix renders next to its form
// controls (CheckboxBubbleInput, SwitchBubbleInput, RadioBubbleInput,
// SliderBubbleInput) so that the controls take part in native forms.
//
// API
//
//   import { takeBubbleInput, syncSize, setChecked, setValue } from "./bubble-input.js"
//
//   // The server renders the input right after the control (React's first
//   // render, where isFormControl defaults to true). Radix keeps it only
//   // when the control is in a form (or has a form attribute): returns the
//   // input, or null after removing it.
//   const input = takeBubbleInput(control, input)
//
//   const stop = syncSize(control, input) // Radix useSize: width/height
//                                         // follow the control's border box
//
//   // The bubble input's effect: set the property and dispatch the event
//   // forms and listeners observe.
//   setChecked(input, true, { bubbles, indeterminate })
//   setValue(input, "42")

/** Returns input if control is a form control, else removes input and returns null. */
export function takeBubbleInput(control, input, formControl = control) {
  if (!input) return null
  const isFormControl = !!formControl.getAttribute("form") || !!formControl.closest("form")
  if (isFormControl) return input
  input.remove()
  return null
}

/** Mirrors the control's border-box size on input (Radix useSize). Returns a cleanup. */
export function syncSize(control, input) {
  const apply = (width, height) => {
    input.style.width = width + "px"
    input.style.height = height + "px"
  }
  apply(control.offsetWidth, control.offsetHeight)
  if (typeof ResizeObserver === "undefined") return () => {}
  const ro = new ResizeObserver((entries) => {
    const entry = entries[entries.length - 1]
    const box = Array.isArray(entry.borderBoxSize) ? entry.borderBoxSize[0] : entry.borderBoxSize
    if (box) apply(box.inlineSize, box.blockSize)
    else apply(control.offsetWidth, control.offsetHeight)
  })
  ro.observe(control, { box: "border-box" })
  return () => ro.disconnect()
}

/** Sets input.checked and dispatches a "click" event, like Radix's bubble input effect. */
export function setChecked(input, checked, { bubbles = true, indeterminate } = {}) {
  if (indeterminate !== undefined) input.indeterminate = indeterminate
  input.checked = checked
  input.dispatchEvent(new Event("click", { bubbles }))
}

/** Sets input.value and dispatches a bubbling "input" event. */
export function setValue(input, value) {
  input.value = value
  input.dispatchEvent(new Event("input", { bubbles: true }))
}
