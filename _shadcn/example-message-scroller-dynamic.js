// The message-scroller-dynamic fixture keeps its rows in React state: the
// buttons append a reply, append a user turn or prepend an older row
// (cloned from the prototypes in template[data-example-ms-rows]), call the
// scroller's commands, and the status line mirrors the visibility and
// scrollable state (useMessageScrollerVisibility, useMessageScrollerScrollable).
import { define } from "./core.js"
import { getMessageScroller } from "./message-scroller.js"

const text = (row) =>
  row.long
    ? `Reply ${row.id}: a longer reply that wraps onto a few lines so new turns have to scroll into place.`
    : `${row.user ? "Question" : "Reply"} ${row.id}`

define("example-message-scroller-dynamic", {
  selector: "[data-example-ms-dynamic]",
  mount(box) {
    box.removeAttribute("data-example-ms-dynamic")
    const tpl = box.querySelector("template[data-example-ms-rows]")
    const viewport = box.querySelector('[data-slot="message-scroller-viewport"]')
    const content = box.querySelector('[data-slot="message-scroller-content"]')
    const status = box.querySelector("[data-example-ms-status]")
    status.removeAttribute("data-example-ms-status")
    let next = 7
    let older = 0
    const scroller = () => getMessageScroller(viewport)

    const makeRow = (row) => {
      const el = tpl.content.querySelector(row.user ? "[data-scroll-anchor=true]" : "[data-scroll-anchor=false]").cloneNode(true)
      el.setAttribute("data-message-id", `r${row.id}`)
      el.querySelector('[data-slot="bubble-content"]').textContent = text(row)
      return el
    }
    const spacer = () => content.querySelector(":scope > [data-message-scroller-spacer]")
    const rows = () => [...content.children].filter((c) => c !== spacer())

    const actions = {
      reply: () => content.insertBefore(makeRow({ id: next++, user: false, long: true }), spacer()),
      turn: () => content.insertBefore(makeRow({ id: next++, user: true, long: false }), spacer()),
      prepend: () => content.insertBefore(makeRow({ id: older--, user: false, long: false }), rows()[0] ?? spacer()),
      jump2: () => scroller()?.scrollToMessage("r2", { align: "center" }),
      jump4: () => scroller()?.scrollToMessage("r4", { align: "end", scrollMargin: 8 }),
      start: () => scroller()?.scrollToStart(),
      end: () => scroller()?.scrollToEnd(),
    }
    const buttons = [...box.querySelectorAll("[data-example-ms-action]")]
    const listeners = buttons.map((b) => {
      const fn = actions[b.getAttribute("data-example-ms-action")]
      b.removeAttribute("data-example-ms-action")
      b.addEventListener("click", fn)
      return [b, fn]
    })

    let visibility = { currentAnchorId: null, visibleMessageIds: [] }
    let scrollable = { start: false, end: false }
    const render = () =>
      status.replaceChildren(
        "Current: ",
        visibility.currentAnchorId ?? "none",
        " · Visible: ",
        visibility.visibleMessageIds.join(" "),
        " · Scrollable: ",
        scrollable.start ? "start" : "-",
        " ",
        scrollable.end ? "end" : "-"
      )
    let unsubscribe = null
    const onScrollable = (e) => {
      scrollable = e.detail
      render()
    }
    viewport.addEventListener("shadcn:scrollable-change", onScrollable)
    // Status subscribes in its passive effects, after the scroller mounted.
    queueMicrotask(() => {
      const s = scroller()
      if (!s) return
      scrollable = s.scrollable
      render()
      unsubscribe = s.subscribeVisibility((v) => {
        visibility = v
        render()
      })
    })
    return () => {
      listeners.forEach(([b, fn]) => b.removeEventListener("click", fn))
      viewport.removeEventListener("shadcn:scrollable-change", onScrollable)
      unsubscribe?.()
    }
  },
})
