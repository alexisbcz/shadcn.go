// switch.js: Radix Switch (@radix-ui/react-switch 1.2.6).
//
// - A click (Space and Enter activate the button natively) toggles the
//   switch: aria-checked and data-state on the button and its thumb.
// - The hidden <input type="checkbox"> (Radix SwitchBubbleInput) is kept
//   only inside a form; its checked attribute and property follow the
//   switch and a "click" is dispatched on it when it changes.
// - "shadcn:checked-change" (cancelable, detail: { checked }) is dispatched
//   on the switch before every change.
import { define, emit } from "./core.js"
import { takeBubbleInput, syncSize, setChecked as bubble } from "./bubble-input.js"

define("switch", {
  selector: '[data-slot="switch"]',
  mount(control) {
    const cleanups = []
    const next = control.nextElementSibling
    const candidate = next?.matches('input[type="checkbox"][aria-hidden="true"]') ? next : null
    // Radix's form prop is on the input only (not on the button).
    const input = takeBubbleInput(control, candidate, candidate ?? control)
    if (input) cleanups.push(syncSize(control, input))
    let stopped = false

    const onClick = (event) => {
      if (event.defaultPrevented) return
      const checked = control.getAttribute("aria-checked") !== "true"
      const changed = emit(control, "shadcn:checked-change", { checked })
      if (changed) {
        const state = checked ? "checked" : "unchecked"
        control.setAttribute("aria-checked", String(checked))
        control.setAttribute("data-state", state)
        for (const thumb of control.querySelectorAll(':scope > [data-slot="switch-thumb"]')) thumb.setAttribute("data-state", state)
        // React re-renders the bubble input with defaultChecked={checked}.
        if (input) input.defaultChecked = checked
      }
      if (input) {
        stopped = event.cancelBubble
        if (!stopped) event.stopPropagation()
        if (changed) bubble(input, checked, { bubbles: !stopped })
      }
    }
    control.addEventListener("click", onClick)
    cleanups.push(() => control.removeEventListener("click", onClick))
    return () => cleanups.forEach((fn) => fn())
  },
})
