// toggle-group.js: Radix ToggleGroup (@radix-ui/react-toggle-group 1.1.11).
//
// Items are toggles; in a single group (items with role="radio") pressing
// one releases the others and pressing the pressed one releases it, in a
// multiple group each item toggles on its own. Keyboard navigation is the
// roving focus group (loop by default; data-loop="false" turns it off,
// data-roving-focus="false" disables roving focus entirely).
import { define, emit } from "./core.js"
import { createRovingFocusGroup } from "./roving-focus.js"

const ITEM = '[data-slot="toggle-group-item"]'

// Radix props with no DOM trace (loop, rovingFocus, activationMode) are
// server-rendered as data attributes only when they differ from the default;
// takeOption reads one and removes it, so the DOM matches React's, and
// remembers it in case the element is mounted again.
const savedOptions = new WeakMap()
function takeOption(el, name) {
  let saved = savedOptions.get(el)
  if (!saved) savedOptions.set(el, (saved = {}))
  if (el.hasAttribute(name)) {
    saved[name] = el.getAttribute(name)
    el.removeAttribute(name)
  }
  return saved[name] ?? null
}

define("toggle-group", {
  selector: '[data-slot="toggle-group"]',
  mount(root) {
    const items = () => [...root.querySelectorAll(ITEM)].filter((el) => el.closest('[data-slot="toggle-group"]') === root)
    const isSingle = (item) => item.getAttribute("role") === "radio"
    const isPressed = (item) => item.getAttribute("data-state") === "on"
    const setPressed = (item, pressed) => {
      item.setAttribute(isSingle(item) ? "aria-checked" : "aria-pressed", String(pressed))
      item.setAttribute("data-state", pressed ? "on" : "off")
    }

    const onClick = (event) => {
      const item = event.target.closest(ITEM)
      if (event.defaultPrevented || !item || !items().includes(item) || item.disabled) return
      const pressed = !isPressed(item)
      const single = isSingle(item)
      const next = items().filter((it) => (it === item ? pressed : !single && isPressed(it)))
      if (!emit(root, "shadcn:value-change", { item, pressed, pressedItems: next })) return
      for (const it of items()) setPressed(it, next.includes(it))
    }
    root.addEventListener("click", onClick)

    let roving = null
    const loop = takeOption(root, "data-loop") !== "false"
    if (takeOption(root, "data-roving-focus") !== "false") {
      roving = createRovingFocusGroup(root, {
        // Radix gives the items its orientation; the styled (bases) group
        // renders data-orientation on the root only, without one for Radix.
        orientation: root.querySelector("[data-radix-collection-item]")?.getAttribute("data-orientation") || undefined,
        loop,
        focusable: (item) => !item.disabled,
        active: isPressed,
      })
    }
    return () => {
      root.removeEventListener("click", onClick)
      roving?.destroy()
    }
  },
})
