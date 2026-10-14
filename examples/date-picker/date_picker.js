// The date-picker examples keep the picked date in React state: the
// trigger shows it, and the calendar the popover mounts on its next opening
// starts from it. Mounted on the trigger (data-example-date-picker).
import { define } from "./core.js"
import { setCalendarProps } from "./calendar.js"

const MONTHS = ["January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"]

function parse(iso) {
  const [y, m, d] = iso.split("-").map(Number)
  return { y, m, d }
}

function ordinal(n) {
  const r100 = n % 100
  if (r100 > 20 || r100 < 10) {
    if (r100 % 10 === 1) return n + "st"
    if (r100 % 10 === 2) return n + "nd"
    if (r100 % 10 === 3) return n + "rd"
  }
  return n + "th"
}

// date-fns format(date, "PPP") and format(date, "LLL dd, y") in en-US.
const ppp = (iso) => {
  const { y, m, d } = parse(iso)
  return `${MONTHS[m - 1]} ${ordinal(d)}, ${y}`
}
const short = (iso) => {
  const { y, m, d } = parse(iso)
  return `${MONTHS[m - 1].slice(0, 3)} ${String(d).padStart(2, "0")}, ${y}`
}

define("example-date-picker", {
  selector: "[data-example-date-picker]",
  mount(trigger) {
    const kind = trigger.getAttribute("data-example-date-picker")
    trigger.removeAttribute("data-example-date-picker")
    const contentId = trigger.getAttribute("aria-controls")
    const icon = trigger.querySelector("svg")
    const template = () => document.querySelector(`template[data-content-id="${CSS.escape(contentId)}"]`)

    const onSelect = (e) => {
      const content = e.target.closest?.('[data-slot="popover-content"]')
      if (!content || content.id !== contentId) return
      const selected = e.detail.selected
      const tpl = template()
      if (kind === "single") {
        trigger.classList.toggle("text-muted-foreground", !selected)
        if (selected) trigger.replaceChildren(icon, ppp(selected))
        else {
          const span = document.createElement("span")
          span.textContent = "Pick a date"
          trigger.replaceChildren(icon, span)
        }
        if (tpl) setCalendarProps(tpl.content, { selected })
      } else {
        trigger.classList.toggle("text-muted-foreground", !selected)
        if (selected?.from) {
          if (selected.to) trigger.replaceChildren(icon, short(selected.from), " -", " ", short(selected.to))
          else trigger.replaceChildren(icon, short(selected.from))
        } else {
          const span = document.createElement("span")
          span.textContent = "Pick a date"
          trigger.replaceChildren(icon, span)
        }
        // defaultMonth={date?.from}
        if (tpl) setCalendarProps(tpl.content, { selected, defaultMonth: selected?.from ?? null })
      }
    }
    document.addEventListener("shadcn:select", onSelect)
    return () => document.removeEventListener("shadcn:select", onSelect)
  },
})
