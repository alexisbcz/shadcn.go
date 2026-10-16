// The input-group-button example's state: the copy button copies the URL
// (useCopyToClipboard: navigator.clipboard, else a hidden textarea and
// execCommand) and shows a check icon for 2 s; the star button toggles
// data-favorite on its icon.
import { define } from "./core.js"

const SVG = "http://www.w3.org/2000/svg"
const ICONS = {
  check: ["M5 12l5 5l10 -10"],
  copy: [
    "M7 7m0 2.667a2.667 2.667 0 0 1 2.667 -2.667h8.666a2.667 2.667 0 0 1 2.667 2.667v8.666a2.667 2.667 0 0 1 -2.667 2.667h-8.666a2.667 2.667 0 0 1 -2.667 -2.667z",
    "M4.012 16.737a2.005 2.005 0 0 1 -1.012 -1.737v-10c0 -1.1 .9 -2 2 -2h10c.75 0 1.158 .385 1.5 1",
  ],
}

function tablerIcon(name) {
  const svg = document.createElementNS(SVG, "svg")
  const attrs = { xmlns: SVG, width: "24", height: "24", viewBox: "0 0 24 24", fill: "none", stroke: "currentColor", "stroke-width": "2", "stroke-linecap": "round", "stroke-linejoin": "round", class: `tabler-icon tabler-icon-${name}` }
  for (const [k, v] of Object.entries(attrs)) svg.setAttribute(k, v)
  for (const d of ICONS[name]) {
    const path = document.createElementNS(SVG, "path")
    path.setAttribute("d", d)
    svg.append(path)
  }
  return svg
}

function legacyCopyToClipboard(value) {
  const textArea = document.createElement("textarea")
  textArea.value = value
  textArea.setAttribute("readonly", "")
  textArea.style.position = "fixed"
  textArea.style.opacity = "0"
  textArea.style.pointerEvents = "none"
  document.body.appendChild(textArea)
  textArea.focus()
  textArea.select()
  textArea.setSelectionRange(0, value.length)
  let hasCopied = false
  try {
    hasCopied = document.execCommand("copy")
  } catch {
    hasCopied = false
  }
  document.body.removeChild(textArea)
  return hasCopied
}

define("example-input-group-copy", {
  selector: '[data-slot="input-group"] button[aria-label="Copy"][title="Copy"]',
  mount(button) {
    let copied = false
    let timer
    const render = () => button.replaceChildren(tablerIcon(copied ? "check" : "copy"))
    const onClick = async () => {
      const value = "https://x.com/shadcn"
      let hasCopied = false
      if (navigator.clipboard?.writeText) {
        try {
          await navigator.clipboard.writeText(value)
          hasCopied = true
        } catch {
          hasCopied = legacyCopyToClipboard(value)
        }
      } else {
        hasCopied = legacyCopyToClipboard(value)
      }
      if (!hasCopied) return
      if (!copied) {
        copied = true
        render()
      }
      timer = setTimeout(() => {
        if (copied) {
          copied = false
          render()
        }
      }, 2000)
    }
    button.addEventListener("click", onClick)
    return () => {
      clearTimeout(timer)
      button.removeEventListener("click", onClick)
    }
  },
})

define("example-input-group-favorite", {
  selector: '[data-slot="input-group"] button:has(> svg[data-favorite])',
  mount(button) {
    const onClick = () => {
      const icon = button.querySelector(":scope > svg[data-favorite]")
      icon?.setAttribute("data-favorite", String(icon.getAttribute("data-favorite") !== "true"))
    }
    button.addEventListener("click", onClick)
    return () => button.removeEventListener("click", onClick)
  },
})
