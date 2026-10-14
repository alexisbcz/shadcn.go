// Site behaviors of the Go docs site that are not component behaviors
// (those live in the component runtime): the theme switch (next-themes:
// localStorage "theme", class on <html>, color-scheme).
//
// Extension point: header.go renders elements with data-site-action; add
// handlers below.

const root = document.documentElement

function resolved(theme) {
  if (theme === "system") {
    return // ThemeShortcut (components/theme-provider.tsx): "d" toggles the theme
// outside editable fields.
document.addEventListener("keydown", (e) => {
  if ((e.key === "d" || e.key === "D") && !e.metaKey && !e.ctrlKey && !e.altKey) {
    const t = e.target
    if ((t instanceof HTMLElement && t.isContentEditable) || t instanceof HTMLInputElement || t instanceof HTMLTextAreaElement || t instanceof HTMLSelectElement) return
    e.preventDefault()
    toggleTheme()
  }
})

matchMedia("(prefers-color-scheme: dark)").matches ? "dark" : "light"
  }
  return theme
}

export function setTheme(theme) {
  try {
    localStorage.setItem("theme", theme)
  } catch (_) {}
  const r = resolved(theme)
  // next-themes' disableTransitionOnChange: no transition runs while the
  // theme classes change.
  const css = document.createElement("style")
  css.appendChild(document.createTextNode("*,*::before,*::after{-webkit-transition:none!important;-moz-transition:none!important;-o-transition:none!important;-ms-transition:none!important;transition:none!important}"))
  document.head.appendChild(css)
  root.classList.remove("light", "dark")
  root.classList.add(r)
  root.style.colorScheme = r
  const meta = document.querySelector('meta[name="theme-color"]')
  if (meta) meta.setAttribute("content", r === "dark" ? "#09090b" : "#ffffff")
  window.getComputedStyle(document.body)
  setTimeout(() => document.head.removeChild(css), 1)
}

function toggleTheme() {
  setTheme(resolved(currentTheme()) === "dark" ? "light" : "dark")
}

function currentTheme() {
  try {
    return localStorage.getItem("theme") || "system"
  } catch (_) {
    return "system"
  }
}

document.addEventListener("click", (event) => {
  const el = event.target.closest("[data-site-action]")
  if (!el) return
  switch (el.dataset.siteAction) {
    case "toggle-theme":
      // ModeSwitcher: toggles between the resolved theme and its opposite.
      toggleTheme()
      break
  }
})

// ThemeShortcut (components/theme-provider.tsx): "d" toggles the theme
// outside editable fields.
document.addEventListener("keydown", (e) => {
  if ((e.key === "d" || e.key === "D") && !e.metaKey && !e.ctrlKey && !e.altKey) {
    const t = e.target
    if ((t instanceof HTMLElement && t.isContentEditable) || t instanceof HTMLInputElement || t instanceof HTMLTextAreaElement || t instanceof HTMLSelectElement) return
    e.preventDefault()
    toggleTheme()
  }
})

matchMedia("(prefers-color-scheme: dark)").addEventListener("change", () => {
  if (currentTheme() === "system") setTheme("system")
})

// Block viewer (components/block-viewer.tsx): Preview/Code tabs switch the
// wrapper's data-view, the size toggle resizes the preview panel, refresh
// reloads the iframe, the file tree selects the code shown, copy buttons.
document.addEventListener("click", (e) => {
  const t = e.target instanceof Element ? e.target : null
  const wrapper = t?.closest("[data-block-viewer]")
  if (!wrapper) return
  const tab = t.closest("[data-block-view]")
  if (tab) wrapper.dataset.view = tab.dataset.blockView
  const size = t.closest("[data-block-size] [data-slot=toggle-group-item]")
  if (size) {
    wrapper.dataset.view = "preview"
    const value = size.getAttribute("value") ?? size.dataset.value
    const panel = wrapper.querySelector("[data-block-panel]")
    if (value && panel?.resizablePanel) panel.resizablePanel.resize(value)
  }
  if (t.closest("[data-block-refresh]")) {
    const iframe = wrapper.querySelector("iframe")
    if (iframe) iframe.src = iframe.src
  }
  const file = t.closest("[data-block-select-file]")
  if (file) {
    const target = file.dataset.blockSelectFile
    for (const f of wrapper.querySelectorAll("[data-block-file]")) f.hidden = f.dataset.blockFile !== target
    for (const b of wrapper.querySelectorAll("[data-block-select-file]")) b.dataset.active = String(b === file)
  }
  const copy = t.closest("[data-block-copy]")
  if (copy) navigator.clipboard?.writeText(copy.dataset.blockCopy)
  const copyFile = t.closest("[data-block-copy-file]")
  if (copyFile) {
    const fig = copyFile.closest("[data-block-file]")
    navigator.clipboard?.writeText(fig?.querySelector("code")?.textContent ?? "")
  }
})
