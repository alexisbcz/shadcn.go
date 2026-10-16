// The sidebar-01/02 VersionSwitcher and sidebar-07 TeamSwitcher: selecting
// an item of the switcher's dropdown updates the trigger (React state) and,
// for versions, moves the check mark (the menu re-renders its items from
// its template on every open).
import { define } from "./core.js"

const ITEM = '[data-slot="dropdown-menu-item"]'

define("example-blocks-sidebar-switcher", {
  selector: "[data-example-blocks-sidebar-switcher]",
  mount(trigger) {
    const kind = trigger.getAttribute("data-example-blocks-sidebar-switcher")
    trigger.removeAttribute("data-example-blocks-sidebar-switcher")
    const tpl = document.querySelector(`template[data-trigger-id="${CSS.escape(trigger.id)}"]`)
    if (!tpl) return
    const contentId = tpl.getAttribute("data-content-id")
    const li = trigger.closest('[data-slot="sidebar-menu-item"]')

    const selectVersion = (index, item) => {
      const version = item.textContent.trim().replace(/^v/, "")
      const label = trigger.querySelector(".flex-col > span:last-child")
      if (label) label.replaceChildren("v", version)
      for (const root of [tpl.content, document.getElementById(contentId)]) {
        if (!root) continue
        const items = [...root.querySelectorAll(ITEM)]
        const check = items.map((el) => el.querySelector(":scope > svg")).find(Boolean)
        if (!check) continue
        items[index]?.append(check)
      }
    }

    const selectTeam = (index, item) => {
      const logo = li?.querySelector(`template[data-team-logo="${index}"]`)?.content.firstElementChild
      const box = trigger.querySelector(":scope > div")
      if (logo && box) box.replaceChildren(logo.cloneNode(true))
      const spans = trigger.querySelectorAll(".grid > span")
      const name = [...item.childNodes].filter((n) => n.nodeType === Node.TEXT_NODE).map((n) => n.textContent).join("")
      const plan = li?.querySelector(`template[data-team-logo="${index}"]`)?.getAttribute("data-plan") ?? ""
      if (spans[0]) spans[0].textContent = name
      if (spans[1]) spans[1].textContent = plan
    }

    const onSelect = (event) => {
      const item = event.target.closest?.(ITEM)
      const content = item?.closest('[data-slot="dropdown-menu-content"]')
      if (!content || content.id !== contentId) return
      const index = [...content.querySelectorAll(ITEM)].indexOf(item)
      if (kind === "version") selectVersion(index, item)
      else if (kind === "team" && index < content.querySelectorAll(ITEM).length - 1) selectTeam(index, item)
    }
    document.addEventListener("shadcn:select", onSelect)
    return () => document.removeEventListener("shadcn:select", onSelect)
  },
})
