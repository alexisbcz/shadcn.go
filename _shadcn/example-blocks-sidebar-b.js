// Behavior of the sidebar-09 and sidebar-16 blocks.
import { define } from "./core.js"
import { getSidebar } from "./sidebar.js"

// sidebar-09: clicking a folder makes it active, shows its name, shuffles
// the mails (data.mails.sort(() => Math.random() - 0.5), in place, like
// upstream) and expands the sidebar.
define("example-sidebar-09", {
  selector: "[data-example-sidebar-09]",
  mount(root) {
    root.removeAttribute("data-example-sidebar-09")
    const list = root.querySelector("[data-mail-list]")
    const title = root.querySelector("[data-mail-title]")
    const buttons = [...root.querySelectorAll("[data-mail-nav]")]
    if (!list || !title) return
    const mails = [...list.querySelectorAll(":scope > [data-mail]")]
    const onClick = (event) => {
      const btn = event.target.closest?.("[data-mail-nav]")
      if (!btn || !root.contains(btn)) return
      const name = btn.getAttribute("data-mail-nav")
      for (const b of buttons) b.setAttribute("data-active", String(b === btn))
      title.textContent = name
      mails.sort(() => Math.random() - 0.5)
      list.replaceChildren(...mails.slice(0, Math.max(5, Math.floor(Math.random() * 10) + 1)))
      getSidebar(root)?.setOpen(true)
    }
    root.addEventListener("click", onClick)
    return () => root.removeEventListener("click", onClick)
  },
})

// sidebar-16: the site header's button toggles the sidebar.
define("example-sidebar-toggle", {
  selector: "[data-example-sidebar-toggle]",
  mount(btn) {
    btn.removeAttribute("data-example-sidebar-toggle")
    const onClick = () => getSidebar(btn)?.toggleSidebar()
    btn.addEventListener("click", onClick)
    return () => btn.removeEventListener("click", onClick)
  },
})
