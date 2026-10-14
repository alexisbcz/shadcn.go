// Site behaviors of the Go docs site that are not component behaviors
// (those live in the component runtime): the theme switch (next-themes:
// localStorage "theme", class on <html>, color-scheme).
//
// Extension point: header.go renders elements with data-site-action; add
// handlers below.

const root = document.documentElement

function resolved(theme) {
  if (theme === "system") {
    return matchMedia("(prefers-color-scheme: dark)").matches ? "dark" : "light"
  }
  return theme
}

export function setTheme(theme) {
  try {
    localStorage.setItem("theme", theme)
  } catch (_) {}
  const r = resolved(theme)
  root.classList.remove("light", "dark")
  root.classList.add(r)
  root.style.colorScheme = r
  const meta = document.querySelector('meta[name="theme-color"]')
  if (meta) meta.setAttribute("content", r === "dark" ? "#09090b" : "#ffffff")
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
      setTheme(resolved(currentTheme()) === "dark" ? "light" : "dark")
      break
  }
})

matchMedia("(prefers-color-scheme: dark)").addEventListener("change", () => {
  if (currentTheme() === "system") setTheme("system")
})
