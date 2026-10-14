// focus-guards.js: port of Radix FocusGuards (@radix-ui/react-focus-guards
// 1.1.3).
//
// While at least one guarded layer is open (modal Dialog content, Popover,
// Menu, Select, ... content), two invisible focusable spans sit at the very
// start and end of <body>, so focus can always be caught when it would
// otherwise leave the document (tabbing past the last element), keeping
// FocusScope's trap effective. They are shared: nested layers reuse the
// same pair, and the last release removes them.
//
// API
//
//   import { acquireFocusGuards } from "./focus-guards.js"
//
//   const release = acquireFocusGuards()  // Radix useFocusGuards() effect:
//                                         // (re)inserts the guards at the
//                                         // body edges, counts the user
//   release()                             // the effect cleanup
//
// Call it where Radix runs the effect: for a Dialog, after the content was
// mounted and focused, before hideOthers (see dialog.js), so the guards end
// up aria-hidden like upstream.

let count = 0

export function acquireFocusGuards() {
  const edgeGuards = document.querySelectorAll("[data-radix-focus-guard]")
  document.body.insertAdjacentElement("afterbegin", edgeGuards[0] ?? createFocusGuard())
  document.body.insertAdjacentElement("beforeend", edgeGuards[1] ?? createFocusGuard())
  count++
  let released = false
  return () => {
    if (released) return
    released = true
    if (count === 1) {
      document.querySelectorAll("[data-radix-focus-guard]").forEach((node) => node.remove())
    }
    count--
  }
}

function createFocusGuard() {
  const element = document.createElement("span")
  element.setAttribute("data-radix-focus-guard", "")
  element.tabIndex = 0
  element.style.outline = "none"
  element.style.opacity = "0"
  element.style.position = "fixed"
  element.style.pointerEvents = "none"
  return element
}
