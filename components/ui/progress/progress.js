// progress.js: shadcn Progress on Radix Progress (@radix-ui/react-progress 1.1.7).
//
// Progress has no interaction; this module lets applications change the
// value like a React re-render with a new value prop would:
//
//   import { setProgress } from "./progress.js"
//   setProgress(root, 66)   // root: [data-slot="progress"]
//
// Like upstream, the value only moves the indicator (translateX): shadcn
// does not pass it to Radix, so the ARIA and data-state attributes stay
// indeterminate. "shadcn:value-change" (detail: { value }) is dispatched on
// the root.
import { emit } from "./core.js"

/**
 * @param {HTMLElement} root
 * @param {number | null} value
 */
export function setProgress(root, value) {
  const indicator = root.querySelector(':scope > [data-slot="progress-indicator"]')
  if (indicator) indicator.style.transform = `translateX(-${100 - (value || 0)}%)`
  emit(root, "shadcn:value-change", { value })
}
