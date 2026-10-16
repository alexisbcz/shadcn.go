// direction.js: port of @radix-ui/react-direction 1.1.1 (useDirection).
//
// Radix resolves a component's reading direction as
//   localDir (the component's `dir` prop) || DirectionProvider's dir || "ltr".
// Server-rendered components carry the resolved value in their `dir`
// attribute (Tabs, ToggleGroup, RadioGroup, menus, ... render dir="ltr" or
// dir="rtl" exactly like Radix does), and a DirectionProvider is any
// ancestor with a `dir` attribute, so the DOM answers the same question:
//
//   import { getDirection } from "./direction.js"
//   const dir = getDirection(groupEl)              // "ltr" | "rtl"
//   const dir = getDirection(groupEl, "rtl")       // explicit local value wins
//
// Components that do not render `dir` themselves (e.g. a roving item) pass
// their root element; the nearest `[dir]` ancestor (including the element
// itself) decides, and "ltr" is the default.

/**
 * Returns the reading direction for el.
 * @param {Element | null | undefined} el
 * @param {"ltr" | "rtl" | undefined | null} [localDir]
 * @returns {"ltr" | "rtl"}
 */
export function getDirection(el, localDir) {
  if (localDir === "ltr" || localDir === "rtl") return localDir
  const host = el?.closest?.("[dir]")
  const dir = host?.getAttribute("dir")?.toLowerCase()
  return dir === "rtl" ? "rtl" : "ltr"
}

/**
 * Maps ArrowLeft/ArrowRight to their reading-direction-aware meaning, like
 * Radix's getDirectionAwareKey: in RTL, ArrowLeft means "next".
 * @param {string} key
 * @param {"ltr" | "rtl"} dir
 */
export function getDirectionAwareKey(key, dir) {
  if (dir !== "rtl") return key
  return key === "ArrowLeft" ? "ArrowRight" : key === "ArrowRight" ? "ArrowLeft" : key
}
