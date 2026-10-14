// The message-scroller docs examples' React glue: the opening-position
// tabs, the "Jump to..." and transcript-outline commands, and the
// scrollable status footer.
import { define } from "./core.js"
import { getMessageScroller } from "./message-scroller.js"

const scrollerOf = (messageId) => {
  const item = document.querySelector(`[data-slot="message-scroller-item"][data-message-id="${CSS.escape(messageId)}"]`)
  return item ? getMessageScroller(item) : undefined
}

// Commands: menu items and outline buttons jump to a message.
document.addEventListener("click", (event) => {
  const el = event.target.closest?.("[data-example-message-scroller-jump]")
  if (!el) return
  const id = el.dataset.exampleMessageScrollerJump
  scrollerOf(id)?.scrollToMessage(id, { align: "start", behavior: "smooth" })
})

define("example-message-scroller-position", {
  selector: "[data-example-message-scroller-position]",
  mount(tabs) {
    const card = tabs.closest('[data-slot="card"]')
    const onChange = (e) => {
      if (e.target !== tabs) return
      const viewport = card?.querySelector('[data-slot="message-scroller-viewport"]')
      const s = viewport && getMessageScroller(viewport)
      if (!s) return
      const value = e.detail.value
      requestAnimationFrame(() => {
        if (value === "start") s.scrollToStart({ behavior: "auto" })
        else if (value === "end") s.scrollToEnd({ behavior: "auto" })
        else s.scrollToMessage("open-3", { align: "start", behavior: "auto", scrollMargin: 64 })
      })
    }
    tabs.addEventListener("shadcn:value-change", onChange)
    return () => tabs.removeEventListener("shadcn:value-change", onChange)
  },
})

define("example-message-scroller-status", {
  selector: "[data-example-message-scroller-status]",
  mount(footer) {
    const viewport = footer.closest('[data-slot="card"]')?.querySelector('[data-slot="message-scroller-viewport"]')
    if (!viewport) return
    const onChange = (e) => {
      const { start, end } = e.detail
      footer.textContent =
        start && end
          ? "You can scroll both ways."
          : end
            ? "You are at the top. You can only scroll down."
            : start
              ? "You are at the bottom. You can only scroll up."
              : "All messages fit in the viewport."
    }
    viewport.addEventListener("shadcn:scrollable-change", onChange)
    return () => viewport.removeEventListener("shadcn:scrollable-change", onChange)
  },
})
