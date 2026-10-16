// checkbox.js: Radix Checkbox (@radix-ui/react-checkbox 1.3.3).
//
// - A click (or Space, natively) toggles the checkbox; an indeterminate
//   checkbox becomes checked. Enter does nothing.
// - The indicator is mounted while the checkbox is checked or indeterminate
//   and unmounted (after its exit animation, if any: Radix Presence)
//   otherwise. Its markup comes from the server-rendered template.
// - The hidden <input type="checkbox"> (Radix CheckboxBubbleInput) is kept
//   only inside a form; its checked state follows the checkbox and a
//   "click" is dispatched on it when it changes, so forms see the change.
//   A form reset restores the initial state.
// - "shadcn:checked-change" (cancelable, detail: { checked: true | false |
//   "indeterminate" }) is dispatched on the checkbox before every change.
import { define, emit } from "./core.js"
import { createPresence } from "./presence.js"
import { takeBubbleInput, syncSize, setChecked as bubble } from "./bubble-input.js"

const parse = (s) => (s === "indeterminate" ? "indeterminate" : s === "checked")
const stateName = (c) => (c === "indeterminate" ? "indeterminate" : c ? "checked" : "unchecked")

define("checkbox", {
  selector: '[data-slot="checkbox"]',
  mount(control) {
    const cleanups = []
    const on = (el, type, fn) => {
      el.addEventListener(type, fn)
      cleanups.push(() => el.removeEventListener(type, fn))
    }

    let checked = parse(control.getAttribute("data-state"))
    const initial = checked

    const next = control.nextElementSibling
    const candidate = next?.matches('input[type="checkbox"][aria-hidden="true"]') ? next : null
    // Radix's form prop is on the input only (not on the button).
    const input = takeBubbleInput(control, candidate, candidate ?? control)
    if (input) cleanups.push(syncSize(control, input))
    let stopped = false

    let template = null
    const tpl = control.querySelector(':scope > template[data-slot="checkbox-indicator-template"]')
    if (tpl) {
      template = tpl.content
      tpl.remove()
    }
    let indicator = control.querySelector(':scope > [data-slot="checkbox-indicator"]')
    let presence = null
    const track = (node) => {
      presence?.destroy()
      presence = createPresence(node, { present: true, onExitComplete: () => node.remove() })
    }
    if (indicator) {
      if (!template) template = { firstElementChild: indicator.cloneNode(true) }
      track(indicator)
    }

    const render = () => {
      const state = stateName(checked)
      control.setAttribute("aria-checked", checked === "indeterminate" ? "mixed" : String(checked))
      control.setAttribute("data-state", state)
      const present = checked === "indeterminate" || checked === true
      if (present) {
        if (!indicator || !indicator.isConnected) {
          const node = template?.firstElementChild?.cloneNode(true)
          if (!node) return
          node.setAttribute("data-state", state)
          control.appendChild(node)
          indicator = node
          track(node)
        } else {
          indicator.setAttribute("data-state", state)
          presence?.setPresent(true)
        }
      } else if (indicator?.isConnected) {
        indicator.setAttribute("data-state", state)
        if (presence) presence.setPresent(false)
        else indicator.remove()
      }
    }

    const set = (value) => {
      if (value === checked) return false
      if (!emit(control, "shadcn:checked-change", { checked: value })) return false
      checked = value
      render()
      return true
    }

    // Bubble input effect, run after the change is committed.
    const sync = () => {
      if (!input) return
      bubble(input, checked === "indeterminate" ? false : checked, {
        bubbles: !stopped,
        indeterminate: checked === "indeterminate",
      })
    }

    on(control, "keydown", (event) => {
      if (event.key === "Enter") event.preventDefault()
    })
    on(control, "click", (event) => {
      if (event.defaultPrevented) return
      const changed = set(checked === "indeterminate" ? true : !checked)
      if (input) {
        stopped = event.cancelBubble
        if (!stopped) event.stopPropagation()
      }
      if (changed) sync()
    })
    // Like Radix, the button's own form (an ancestor form: the form prop is
    // not on the button).
    const form = control.form
    if (form) on(form, "reset", () => set(initial) && sync())

    return () => {
      presence?.destroy()
      cleanups.forEach((fn) => fn())
    }
  },
})
