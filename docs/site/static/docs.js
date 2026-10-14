// Client behaviors of the docs page layout (docslayout.go):
//
//   - TOC scroll-spy (components/docs-toc.tsx useActiveItem): an
//     IntersectionObserver with rootMargin "0% 0% -80% 0%" over the headings;
//     the last intersecting one is active.
//   - "Copy Page" (components/docs-copy-page.tsx): copies the page source,
//     shows a check icon for 2s (useCopyToClipboard).
//   - Docs sidebar scroll saving (components/docs-sidebar.tsx): the position
//     is restored by the inline DOCS_SIDEBAR_SCROLL_RESTORE_SCRIPT.

const STORAGE_KEY = "docs-sidebar-scroll"

function initTOC() {
  const toc = document.querySelector("[data-docs-toc]")
  if (!toc) return
  const links = [...toc.querySelectorAll("a[href^='#']")]
  const setActive = (id) => {
    for (const a of links) {
      a.setAttribute("data-active", String(a.getAttribute("href") === "#" + id))
    }
  }
  const observer = new IntersectionObserver(
    (entries) => {
      for (const entry of entries) {
        if (entry.isIntersecting) setActive(entry.target.id)
      }
    },
    { rootMargin: "0% 0% -80% 0%" }
  )
  for (const a of links) {
    const el = document.getElementById(decodeURIComponent(a.getAttribute("href").slice(1)))
    if (el) observer.observe(el)
  }
}

const CHECK_SVG =
  '<path d="M5 12l5 5l10 -10"></path>'

function initCopyPage() {
  for (const btn of document.querySelectorAll("[data-docs-copy-page]")) {
    const svg = btn.querySelector("svg")
    const copyShapes = svg ? svg.innerHTML : ""
    const copyClass = svg ? svg.getAttribute("class") : ""
    let timer
    btn.addEventListener("click", () => {
      const src = btn.closest(".group\\/buttons")?.querySelector("template[data-docs-copy-page-source]")
      const text = src ? src.content.textContent : ""
      if (!navigator.clipboard?.writeText) return
      navigator.clipboard.writeText(text).then(() => {
        if (!svg) return
        svg.innerHTML = CHECK_SVG
        svg.setAttribute("class", "tabler-icon tabler-icon-check ")
        clearTimeout(timer)
        timer = setTimeout(() => {
          svg.innerHTML = copyShapes
          svg.setAttribute("class", copyClass)
        }, 2000)
      })
    })
  }
}

function initSidebarScroll() {
  const container = document.querySelector("[data-docs-sidebar-content]")
  if (!container) return
  const save = () => {
    try {
      sessionStorage.setItem(
        STORAGE_KEY,
        JSON.stringify({ pathname: location.pathname, scrollTop: container.scrollTop })
      )
    } catch (_) {}
  }
  save()
  container.addEventListener("scroll", save, { passive: true })
}

initTOC()
initCopyPage()
initSidebarScroll()
