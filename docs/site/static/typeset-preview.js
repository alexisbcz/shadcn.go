// The typeset preview iframe (app/(view)/preview/typeset): applies the
// params the designer posts (preview-provider.tsx) as CSS variables on
// <html>, forwards its keyboard shortcuts to the designer as commands
// (forward-scripts.tsx), and follows theme changes made by the designer
// (next-themes' storage listener).
import { define } from "./core.js"

const PARAMS_MESSAGE = "typeset-params"
const COMMAND_MESSAGE = "typeset-command"

define("typeset-preview", {
  selector: "[data-typeset-preview-root]",
  mount(root) {
    const { fonts, measures } = JSON.parse(root.dataset.typesetPreviewData)
    const findFont = (id) => fonts.find((f) => f.id === id)?.value

    const apply = (params) => {
      const body = findFont(params.body)
      const values = {
        "--preview-size": `${params.scale}px`,
        "--preview-leading": params.leading,
        "--preview-flow": params.flow,
        "--preview-measure": measures[params.measure],
        "--preview-font": body,
        "--preview-font-heading": params.heading === "inherit" ? body : findFont(params.heading),
        "--preview-font-mono": findFont(params.mono),
      }
      const style = document.documentElement.style
      for (const [name, value] of Object.entries(values)) {
        if (value) style.setProperty(name, value)
        else style.removeProperty(name)
      }
    }
    const onMessage = (e) => {
      if (e.data?.type === PARAMS_MESSAGE && e.data.data) apply(e.data.data)
    }
    const send = (command) => {
      if (window.parent && window.parent !== window) {
        window.parent.postMessage({ type: COMMAND_MESSAGE, command }, location.origin)
      }
    }
    const onKeyDown = (e) => {
      const t = e.target
      if ((t instanceof HTMLElement && t.isContentEditable) || t instanceof HTMLInputElement || t instanceof HTMLTextAreaElement || t instanceof HTMLSelectElement) return
      if (e.metaKey || e.ctrlKey) {
        const key = e.key.toLowerCase()
        const command = (key === "z" && e.shiftKey) || (key === "y" && e.ctrlKey) ? "redo" : key === "z" ? "undo" : null
        if (command) {
          e.preventDefault()
          send(command)
        }
        return
      }
      if (e.key === "r" || e.key === "R") {
        e.preventDefault()
        send(e.shiftKey ? "reset" : "shuffle")
      } else if ((e.key === "d" || e.key === "D") && !e.altKey && window.parent !== window) {
        // The site's own "d" shortcut would toggle this document only.
        e.preventDefault()
        e.stopImmediatePropagation()
        send("toggle-theme")
      }
    }
    const onStorage = (e) => {
      if (e.key !== "theme") return
      const theme = e.newValue || "system"
      const resolved = theme === "system" ? (matchMedia("(prefers-color-scheme: dark)").matches ? "dark" : "light") : theme
      const html = document.documentElement
      html.classList.remove("light", "dark")
      html.classList.add(resolved)
      html.style.colorScheme = resolved
    }
    window.addEventListener("message", onMessage)
    window.addEventListener("keydown", onKeyDown, true)
    window.addEventListener("storage", onStorage)
    return () => {
      window.removeEventListener("message", onMessage)
      window.removeEventListener("keydown", onKeyDown, true)
      window.removeEventListener("storage", onStorage)
    }
  },
})
