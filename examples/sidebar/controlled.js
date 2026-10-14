// The sidebar-controlled example: the header button toggles the provider's
// state (setOpen((open) => !open)) and shows it, like the React example
// whose state drives <SidebarProvider open onOpenChange>.
import { define } from "./core.js"
import { getSidebar } from "./sidebar.js"

define("example-sidebar-controlled", {
  selector: "[data-example-sidebar-controlled]",
  mount(wrapper) {
    wrapper.removeAttribute("data-example-sidebar-controlled")
    const header = wrapper.querySelector('[data-slot="sidebar-inset"] > header')
    const btn = header?.querySelector('[data-slot="button"]')
    if (!btn) return
    const icons = {
      open: header.querySelector('template[data-icon="open"]'),
      close: header.querySelector('template[data-icon="close"]'),
    }
    const render = (open) => {
      const icon = icons[open ? "close" : "open"].content.firstElementChild.cloneNode(true)
      // {open ? <PanelLeftCloseIcon /> : <PanelLeftOpenIcon />}
      // <span>{open ? "Close" : "Open"} Sidebar</span>: two text nodes.
      const span = document.createElement("span")
      span.append(open ? "Close" : "Open", " Sidebar")
      btn.replaceChildren(icon, span)
    }
    let open = getSidebar(wrapper)?.open ?? true
    render(open)
    const onChange = (event) => {
      if (event.target !== wrapper) return
      queueMicrotask(() => {
        if (event.defaultPrevented) return
        open = event.detail.open
        render(open)
      })
    }
    const onClick = () => getSidebar(wrapper)?.setOpen((o) => !o)
    wrapper.addEventListener("shadcn:open-change", onChange)
    btn.addEventListener("click", onClick)
    return () => {
      wrapper.removeEventListener("shadcn:open-change", onChange)
      btn.removeEventListener("click", onClick)
    }
  },
})
