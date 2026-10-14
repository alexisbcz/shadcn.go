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
