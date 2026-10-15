// The create designer's preview page (/preview/{base}/{name}): the iframe
// side of components/preview.tsx (use-iframe-sync.tsx, DesignSystemProvider)
// and the scripts the preview page loads (RandomizeScript, HistoryScript,
// DarkModeScript, ActionMenuScript, OpenPresetScript: shortcuts forwarded
// to the designer).
//
// The designer posts {type: "design-system-params", params, css, fonts}.
// Colors, radius, fonts and menu colors apply in place; a style, icon
// library, base or item change reloads the page with the new params (the
// Go components render a style's markup on the server).
import { define } from "./core.js"
import { encodePreset } from "./create.js"

const MENU_SLOTS = ["dropdown-menu-content", "dropdown-menu-sub-content", "select-content", "context-menu-content", "context-menu-sub-content", "menubar-content", "menubar-sub-content", "combobox-content"]

define("create-preview", {
  selector: "[data-create-preview]",
  mount(root) {
    let rendered = JSON.parse(root.getAttribute("data-create-preview"))
    let menuColor = rendered.menuColor

    const applyMenuColor = () => {
      const inverted = menuColor === "inverted" || menuColor === "inverted-translucent"
      for (const el of document.querySelectorAll(MENU_SLOTS.map((s) => `[data-slot="${s}"]`).join(","))) {
        el.classList.toggle("dark", inverted)
      }
    }
    const observer = new MutationObserver(() => menuColor !== "default" && applyMenuColor())
    observer.observe(document.body, { childList: true, subtree: true })

    const onMessage = (e) => {
      if (e.origin !== location.origin || e.source !== window.parent || e.data?.type !== "design-system-params") return
      const { params, css, fonts, fontTitles, styleTitle } = e.data
      const q = new URLSearchParams()
      const code = encodePreset(params)
      if (code !== "b0") q.set("preset", code)
      if (params.pointer) q.set("pointer", "true")
      const url = location.pathname + (q.toString() ? "?" + q : "")
      if (params.style !== rendered.style || params.iconLibrary !== rendered.iconLibrary) {
        location.replace(url)
        return
      }
      history.replaceState(history.state, "", url)
      document.getElementById("design-system-theme-vars").textContent = css
      for (const c of [...document.body.classList]) if (c.startsWith("base-color-")) document.body.classList.remove(c)
      document.body.classList.add(`base-color-${params.baseColor}`)
      document.documentElement.style.setProperty("--font-sans", fonts.sans)
      document.documentElement.style.setProperty("--font-heading", fonts.heading)
      menuColor = params.menuColor
      applyMenuColor()
      // style-overview.tsx and typography-specimen.tsx name the fonts.
      const heading = fontTitles.heading && fontTitles.heading !== fontTitles.body ? fontTitles.heading : null
      for (const el of document.querySelectorAll("[data-create-style-overview]")) el.textContent = `${styleTitle} - ${heading ?? fontTitles.body}`
      for (const el of document.querySelectorAll("[data-create-typography-label]")) el.textContent = `${heading ?? "Inherit"} - ${fontTitles.body ?? "Default"}`
      rendered = { ...rendered, ...params }
    }
    window.addEventListener("message", onMessage)

    // Shortcuts, forwarded to the designer.
    const onKey = (e) => {
      if (window.parent === window) return
      const t = e.target
      if (t instanceof HTMLElement && (t.isContentEditable || t instanceof HTMLInputElement || t instanceof HTMLTextAreaElement || t instanceof HTMLSelectElement)) return
      const key = e.key.toLowerCase()
      let type = null
      if ((e.metaKey || e.ctrlKey) && (key === "k" || key === "p")) type = "cmd-k-forward"
      else if ((e.metaKey || e.ctrlKey) && key === "z") type = e.shiftKey ? "redo-forward" : "undo-forward"
      else if (!e.metaKey && !e.ctrlKey && !e.altKey) {
        if (key === "r") type = e.shiftKey ? "reset-forward" : "randomize-forward"
        else if (key === "o") type = "open-preset-forward"
        else if (key === "d") type = "dark-mode-forward"
      }
      if (!type) return
      e.preventDefault()
      window.parent.postMessage({ type, key: e.key }, "*")
    }
    document.addEventListener("keydown", onKey)

    // The designer's light/dark mode: next-themes in the iframe follows the
    // stored theme the designer writes.
    const onStorage = (e) => {
      if (e.key === "theme" && e.newValue) import("/_site/site.js").then((m) => m.setTheme(e.newValue))
    }
    window.addEventListener("storage", onStorage)
    return () => {
      window.removeEventListener("storage", onStorage)
      window.removeEventListener("message", onMessage)
      document.removeEventListener("keydown", onKey)
      observer.disconnect()
    }
  },
})
