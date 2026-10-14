// docs-mdx.js: behavior of the docs MDX components (docs/site/mdx).
//
// - Copy buttons ([data-slot="copy-button"]): copy-button.tsx and
//   code-block-command.tsx. The copied text is the code the button belongs
//   to (the active tab's command in a command block); the icon becomes a
//   check and data-copied="true" for 2 s.
// - Config tabs ([data-docs-config]): useConfig (jotai atomWithStorage
//   "config"): CodeTabs' installationType and CodeBlockCommand's
//   packageManager are restored on load, saved on change, and every tabs
//   root bound to the same key follows.
// - ComponentPreview's "View Code" button: shows the full source.
// - CodeCollapsibleWrapper ([data-docs-collapsible]): a Radix Collapsible
//   whose content stays mounted; "Expand"/"Collapse" labels.
import { define } from "./core.js"

const CONFIG_KEY = "config"
const DEFAULT_CONFIG = { style: "new-york-v4", packageManager: "pnpm", installationType: "cli" }

function readConfig() {
  try {
    const raw = localStorage.getItem(CONFIG_KEY)
    if (raw) return { ...DEFAULT_CONFIG, ...JSON.parse(raw) }
  } catch {}
  return { ...DEFAULT_CONFIG }
}

function writeConfig(config) {
  try {
    localStorage.setItem(CONFIG_KEY, JSON.stringify(config))
  } catch {}
}

const triggerValue = (t) => {
  const i = t.id.indexOf("-trigger-")
  return i < 0 ? t.id : t.id.slice(i + 9)
}

// selectTab activates the trigger with value in a tabs root (the tabs
// behavior selects on a primary-button mousedown).
function selectTab(root, value) {
  const trigger = [...root.querySelectorAll('[role="tab"]')].find(
    (t) => t.closest('[data-slot="tabs"]') === root && triggerValue(t) === value
  )
  if (!trigger || trigger.getAttribute("aria-selected") === "true") return
  trigger.dispatchEvent(new MouseEvent("mousedown", { bubbles: true, cancelable: true, button: 0 }))
}

define("docs-config-tabs", {
  selector: "[data-docs-config]",
  mount(root) {
    const key = root.getAttribute("data-docs-config")
    const onChange = (event) => {
      if (event.target !== root || event.defaultPrevented) return
      const value = event.detail?.value
      if (!value) return
      const config = readConfig()
      if (config[key] === value) return
      config[key] = value
      writeConfig(config)
      // Let this root update first, then sync the others.
      queueMicrotask(() => {
        for (const other of document.querySelectorAll(`[data-docs-config="${key}"]`)) {
          if (other !== root) selectTab(other, value)
        }
      })
    }
    root.addEventListener("shadcn:value-change", onChange)
    // The tabs behavior mounts on the same element; restore after it.
    const timer = setTimeout(() => {
      const value = readConfig()[key]
      if (value) selectTab(root, value)
    })
    return () => {
      clearTimeout(timer)
      root.removeEventListener("shadcn:value-change", onChange)
    }
  },
})

async function copyText(value) {
  if (!value) return false
  if (navigator.clipboard?.writeText) {
    try {
      await navigator.clipboard.writeText(value)
      return true
    } catch {}
  }
  const ta = document.createElement("textarea")
  ta.value = value
  ta.setAttribute("readonly", "")
  ta.style.position = "fixed"
  ta.style.opacity = "0"
  ta.style.pointerEvents = "none"
  document.body.appendChild(ta)
  ta.focus()
  ta.select()
  ta.setSelectionRange(0, value.length)
  let ok = false
  try {
    ok = document.execCommand("copy")
  } catch {}
  ta.remove()
  return ok
}

const SVG_NS = "http://www.w3.org/2000/svg"
function checkIcon(like) {
  const svg = like.cloneNode(false)
  svg.setAttribute("class", "tabler-icon tabler-icon-check")
  const path = document.createElementNS(SVG_NS, "path")
  path.setAttribute("d", "M5 12l5 5l10 -10")
  svg.appendChild(path)
  return svg
}

function codeFor(button) {
  const value = button.getAttribute("data-docs-value")
  if (value !== null) return value
  const scope = button.parentElement
  if (!scope) return ""
  const active = scope.querySelector('[role="tabpanel"][data-state="active"] code')
  const code = active ?? scope.querySelector("code")
  return code ? code.textContent.replace(/\n$/, "") : ""
}

define("docs-copy-button", {
  selector: '[data-slot="copy-button"]',
  mount(button) {
    let timer = 0
    let original = null
    const reset = () => {
      button.setAttribute("data-copied", "false")
      const svg = button.querySelector("svg")
      if (original && svg) svg.replaceWith(original)
      original = null
    }
    const onClick = async () => {
      if (!(await copyText(codeFor(button)))) return
      button.setAttribute("data-copied", "true")
      const svg = button.querySelector("svg")
      if (svg && !original) {
        original = svg
        svg.replaceWith(checkIcon(svg))
      }
      clearTimeout(timer)
      timer = setTimeout(reset, 2000)
    }
    button.addEventListener("click", onClick)
    return () => {
      clearTimeout(timer)
      button.removeEventListener("click", onClick)
    }
  },
})

define("docs-view-code", {
  selector: '[data-slot="code"][data-mobile-code-visible]',
  mount(panel) {
    const onClick = (event) => {
      if (!event.target.closest("[data-docs-view-code]")) return
      const tpl = panel.querySelector(":scope > template[data-docs-full-code]")
      if (!tpl) return
      panel.setAttribute("data-mobile-code-visible", "true")
      panel.replaceChildren(tpl.content.cloneNode(true))
    }
    panel.addEventListener("click", onClick)
    return () => panel.removeEventListener("click", onClick)
  },
})

define("docs-collapsible", {
  selector: "[data-docs-collapsible]",
  mount(root) {
    const own = (el) => el.closest("[data-docs-collapsible]") === root
    const triggers = [...root.querySelectorAll('[data-slot="collapsible-trigger"]')].filter(own)
    const content = [...root.querySelectorAll('[data-slot="collapsible-content"]')].find(own)
    const set = (open) => {
      const state = open ? "open" : "closed"
      root.setAttribute("data-state", state)
      content?.setAttribute("data-state", state)
      for (const t of triggers) {
        t.setAttribute("aria-expanded", String(open))
        t.setAttribute("data-state", state)
      }
      for (const label of root.querySelectorAll("[data-docs-toggle-label]")) {
        if (own(label)) label.textContent = open ? "Collapse" : "Expand"
      }
    }
    const onClick = () => set(root.getAttribute("data-state") !== "open")
    for (const t of triggers) t.addEventListener("click", onClick)
    return () => {
      for (const t of triggers) t.removeEventListener("click", onClick)
    }
  },
})
