// The Sera page (app/(app)/(styles)/sera): the analytics header's date range
// label follows its radio group, and "Load more content" reveals the other
// editorial rows after two seconds, showing a spinner meanwhile.
const ranges = { "last-7-days": "Last 7 days", "last-30-days": "Last 30 days", "this-month": "This month", "last-month": "Last month" }

document.addEventListener("shadcn:value-change", (e) => {
  const label = ranges[e.detail?.value]
  const el = document.querySelector("[data-sera-date-label]")
  if (label && el) el.textContent = label
})

document.addEventListener("click", (e) => {
  const b = e.target instanceof Element ? e.target.closest("[data-sera-load-more]") : null
  if (!b || b.disabled) return
  b.disabled = true
  const tpl = b.parentElement.querySelector("template[data-sera-spinner]")
  const icon = b.querySelector("svg[data-icon=inline-end]")
  if (tpl && icon) icon.replaceWith(tpl.content.firstElementChild.cloneNode(true))
  setTimeout(() => {
    for (const row of document.querySelectorAll("[data-sera-more]")) row.hidden = false
    b.remove()
  }, 2000)
})

// Lazy previews (components/lazy-preview.tsx): a preview replaces its
// placeholder once it comes within 800px of the viewport.
{
  const load = (wrapper) => {
    const tpl = wrapper.querySelector(":scope > template")
    if (!tpl) return
    wrapper.replaceChildren(tpl.content.cloneNode(true))
  }
  const wrappers = document.querySelectorAll("[data-sera-preview]")
  if (!("IntersectionObserver" in window)) wrappers.forEach(load)
  else {
    for (const w of wrappers) {
      const observer = new IntersectionObserver((entries) => {
        if (!entries.some((e) => e.isIntersecting)) return
        observer.disconnect()
        load(w)
      }, { rootMargin: "800px 0px" })
      observer.observe(w)
    }
  }
}

// Media library: a grid item selects its asset and shows its details.
document.addEventListener("click", (e) => {
  const b = e.target instanceof Element ? e.target.closest("[data-sera-asset]") : null
  if (!b) return
  const root = b.closest("[data-sera-preview]")
  for (const x of root.querySelectorAll("[data-sera-asset]")) x.setAttribute("aria-pressed", String(x === b))
  for (const d of root.querySelectorAll("[data-sera-asset-details]")) d.hidden = d.dataset.seraAssetDetails !== b.dataset.seraAsset
})

// Media library table: a row click toggles its checkbox; the row's
// data-state follows the checkbox.
document.addEventListener("click", (e) => {
  if (!(e.target instanceof Element)) return
  const row = e.target.closest("[data-sera-select-row]")
  if (!row || e.target.closest("[data-slot=dropdown-menu-trigger]")) return
  const box = row.querySelector("[data-slot=checkbox]")
  if (!box) return
  if (!e.target.closest("[data-slot=checkbox]")) box.click()
  queueMicrotask(() => {
    if (box.getAttribute("data-state") === "checked") row.setAttribute("data-state", "selected")
    else row.removeAttribute("data-state")
  })
})
