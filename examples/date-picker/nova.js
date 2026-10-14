// The radix-nova date-picker docs previews keep the picked date in React
// state: the trigger (or an input) shows it formatted, the calendar the
// popover mounts next starts from it, and some close the popover on pick.
// Mounted on [data-nova-date-picker] (JSON options), which holds the
// popover trigger.
import { define } from "./core.js"
import { setCalendarProps } from "./calendar.js"
import { getPopover } from "./popover.js"

const MONTHS = ["January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"]
const parse = (iso) => iso.split("-").map(Number)
function ordinal(n) {
  const r100 = n % 100
  if (r100 > 20 || r100 < 10) {
    if (r100 % 10 === 1) return n + "st"
    if (r100 % 10 === 2) return n + "nd"
    if (r100 % 10 === 3) return n + "rd"
  }
  return n + "th"
}
const formats = {
  // date-fns format(d, "PPP")
  PPP: (iso) => { const [y, m, d] = parse(iso); return `${MONTHS[m - 1]} ${ordinal(d)}, ${y}` },
  // format(d, "LLL dd, y")
  short: (iso) => { const [y, m, d] = parse(iso); return `${MONTHS[m - 1].slice(0, 3)} ${String(d).padStart(2, "0")}, ${y}` },
  // d.toLocaleDateString() (en-US)
  locale: (iso) => { const [y, m, d] = parse(iso); return `${m}/${d}/${y}` },
  // toLocaleDateString("en-US", { day: "2-digit", month: "long", year: "numeric" })
  long: (iso) => { const [y, m, d] = parse(iso); return `${MONTHS[m - 1]} ${String(d).padStart(2, "0")}, ${y}` },
}

define("example-nova-date-picker", {
  selector: "[data-nova-date-picker]",
  mount(root) {
    const opts = JSON.parse(root.getAttribute("data-nova-date-picker"))
    root.removeAttribute("data-nova-date-picker")
    const trigger = root.matches("[aria-controls]") ? root : root.querySelector("[aria-controls]")
    const contentId = trigger.getAttribute("aria-controls")
    const fmt = formats[opts.format]
    const template = () => document.querySelector(`template[data-content-id="${CSS.escape(contentId)}"]`)
    const setLabel = (text, empty) => {
      const label = root.querySelector("[data-dp-label]")
      if (label) label.textContent = text
      // Without a label span (range), the text nodes follow the leading icon.
      else if (opts.range) trigger.replaceChildren(...[...trigger.children].filter((c) => c.tagName === "svg"), text)
      if (opts.dataEmpty) trigger.setAttribute("data-empty", String(empty))
    }
    const onSelect = (e) => {
      const content = e.target.closest?.('[data-slot="popover-content"]')
      if (!content || content.id !== contentId) return
      const selected = e.detail.selected
      if (opts.range) {
        if (selected?.from) setLabel(selected.to ? `${fmt(selected.from)} - ${fmt(selected.to)}` : fmt(selected.from), false)
        else setLabel(opts.empty, true)
      } else {
        setLabel(selected ? fmt(selected) : opts.empty, !selected)
        if (opts.input) {
          const input = root.querySelector("input")
          if (input) {
            input.value = selected ? fmt(selected) : ""
            input.setAttribute("value", input.value)
          }
        }
        if (opts.echo) {
          const echo = root.querySelector("[data-dp-echo]")
          if (echo) echo.textContent = selected ? fmt(selected) : ""
        }
      }
      const tpl = template()
      if (tpl) setCalendarProps(tpl.content, opts.defaultMonth ? { selected, defaultMonth: (opts.range ? selected?.from : selected) ?? null } : { selected })
      if (opts.close) getPopover(contentId)?.setOpen(false)
    }
    const onKey = (e) => {
      if (e.key === "ArrowDown" && e.target.matches?.("input")) {
        e.preventDefault()
        getPopover(contentId)?.setOpen(true)
      }
    }
    const onClick = (e) => {
      if (e.target.closest?.("[data-dp-done]")?.closest('[data-slot="popover-content"]')?.id === contentId) getPopover(contentId)?.setOpen(false)
    }
    document.addEventListener("shadcn:select", onSelect)
    document.addEventListener("click", onClick)
    root.addEventListener("keydown", onKey)
    return () => {
      document.removeEventListener("shadcn:select", onSelect)
      document.removeEventListener("click", onClick)
      root.removeEventListener("keydown", onKey)
    }
  },
})
