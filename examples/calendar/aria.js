// react-aria-components' Calendar / RangeCalendar behavior for the aria-nova
// calendar previews (examples/calendar/aria.go renders the first state):
// selecting a day (a range takes two presses: anchor, then end), paging by
// the visible months with Previous / Next, and calendar-presets' buttons.
// Every change re-renders the grids from the state kept in data attributes.

const VARIANTS = {
  today: "rounded-(--cell-radius) bg-muted text-foreground data-[selected=true]:rounded-none",
  start: "relative isolate z-0 rounded-l-(--cell-radius) bg-muted after:absolute after:inset-y-0 after:right-0 after:w-4 after:bg-muted",
  end: "relative isolate z-0 rounded-r-(--cell-radius) bg-muted after:absolute after:inset-y-0 after:left-0 after:w-4 after:bg-muted",
  unavailable: "text-muted-foreground opacity-50 [&>div]:line-through",
  disabled: "text-muted-foreground opacity-50",
  outside: "text-muted-foreground aria-selected:text-muted-foreground",
}
const BASE = "group/day relative mt-2 aspect-square h-full w-full cursor-default rounded-(--cell-radius) p-0 text-center select-none [&:is(:last-child>[data-selected=true])>div]:rounded-r-(--cell-radius) [&:is(:first-child>[data-selected=true])>div]:rounded-l-(--cell-radius)"
const MONTHS = ["January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"]
const AR_MONTHS = ["يناير", "فبراير", "مارس", "أبريل", "مايو", "يونيو", "يوليو", "أغسطس", "سبتمبر", "أكتوبر", "نوفمبر", "ديسمبر"]
const DAYS = ["Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"]

const parse = (s) => (s ? new Date(s + "T00:00:00Z") : null)
const iso = (d) => d.toISOString().slice(0, 10)
const addDays = (d, n) => new Date(Date.UTC(d.getUTCFullYear(), d.getUTCMonth(), d.getUTCDate() + n))
const addMonths = (d, n) => new Date(Date.UTC(d.getUTCFullYear(), d.getUTCMonth() + n, 1))
const same = (a, b) => a && b && a.getTime() === b.getTime()
const monthName = (root, d) => (root.dataset.locale === "ar" ? AR_MONTHS : MONTHS)[d.getUTCMonth()]
const label = (root, d) => `${monthName(root, d)} ${d.getUTCFullYear()}`
// The ui-rtl calendar's logical classes.
const rtl = (root, c) =>
  root.dataset.locale === "ar"
    ? c.replaceAll("rounded-r-(--cell-radius)", "rounded-e-(--cell-radius)").replaceAll("rounded-l-(--cell-radius)", "rounded-s-(--cell-radius)").replaceAll("after:right-0", "after:end-0").replaceAll("after:left-0", "after:start-0")
    : c
function today() {
  const n = new Date()
  return new Date(Date.UTC(n.getFullYear(), n.getMonth(), n.getDate()))
}

function state(root) {
  const d = root.dataset
  const [rs, re] = (d.range ?? "").split(" ")
  return {
    month: parse(d.month),
    months: +d.months || 1,
    weeks: +d.weeks || 0,
    value: parse(d.value),
    range: d.range !== undefined ? [parse(rs), parse(re)] : null,
    anchor: parse(d.anchor),
    unavailable: new Set((d.unavailable ?? "").split(" ").filter(Boolean)),
  }
}

function render(root) {
  const s = state(root)
  const t = today()
  const visEnd = addDays(addMonths(s.month, s.months), -1)
  const focus = s.value ?? s.range?.[0] ?? t
  const grids = root.querySelectorAll("table[role=grid]")
  const prices = root.dataset.prices !== undefined
  grids.forEach((grid, i) => {
    const m = addMonths(s.month, i)
    grid.setAttribute("aria-label", label(root, m))
    const caption = grid.previousElementSibling
    const h2 = caption.querySelector("h2")
    if (h2) h2.textContent = label(root, m)
    const values = caption.querySelectorAll("[data-slot=select-value] > span:first-child")
    if (values.length === 2) {
      const name = monthName(root, m)
      values[0].textContent = root.dataset.longMonth !== undefined || root.dataset.locale === "ar" ? name : name.slice(0, 3)
      values[1].textContent = String(m.getUTCFullYear())
    }
    const tbody = grid.tBodies[0]
    const tmpl = tbody.querySelector("td > div")
    const dayClass = tmpl.firstElementChild.className
    const weekStart = root.dataset.locale === "ar" ? 6 : 0
    const first = addDays(m, -((m.getUTCDay() - weekStart + 7) % 7))
    const end = addDays(addMonths(m, 1), -1)
    const weeks = s.weeks || Math.floor((Math.round((end - first) / 864e5) + 7) / 7)
    tbody.replaceChildren()
    for (let w = 0; w < weeks; w++) {
      const tr = document.createElement("tr")
      for (let k = 0; k < 7; k++) {
        const day = addDays(first, w * 7 + k)
        const outside = day.getUTCMonth() !== m.getUTCMonth()
        const disabled = day < s.month || day > visEnd
        const unavailable = s.unavailable.has(iso(day))
        let selected = same(day, s.value)
        let start = false
        let stop = false
        if (s.range && s.range[0] && day >= s.range[0] && day <= s.range[1] && !unavailable) {
          selected = true
          start = same(day, s.range[0])
          stop = same(day, s.range[1])
        }
        if (disabled) selected = start = stop = false
        const isToday = same(day, t)
        let cls = BASE
        if (isToday) cls += " " + VARIANTS.today
        if (start) cls += " " + VARIANTS.start
        if (stop) cls += " " + VARIANTS.end
        if (unavailable) cls += " " + VARIANTS.unavailable
        if (disabled) cls += " " + VARIANTS.disabled
        if (outside) cls += " " + VARIANTS.outside
        const td = document.createElement("td")
        td.setAttribute("role", "gridcell")
        const cell = document.createElement("div")
        cell.dataset.reactAriaPressable = "true"
        if (!disabled) cell.tabIndex = same(day, focus) ? 0 : -1
        cell.setAttribute("role", "button")
        if (disabled || unavailable) {
          td.setAttribute("aria-disabled", "true")
          cell.setAttribute("aria-disabled", "true")
        }
        if (selected) td.setAttribute("aria-selected", "true")
        cell.setAttribute("aria-label", `${DAYS[day.getUTCDay()]}, ${MONTHS[day.getUTCMonth()]} ${day.getUTCDate()}, ${day.getUTCFullYear()}`)
        cell.dataset.date = iso(day)
        if (selected) cell.dataset.selected = "true"
        if (start) cell.dataset.selectionStart = "true"
        if (stop) cell.dataset.selectionEnd = "true"
        if (isToday) cell.dataset.today = "true"
        if (disabled) {
          cell.dataset.disabled = "true"
          cell.dataset.outsideVisibleRange = "true"
        }
        if (outside) cell.dataset.outsideMonth = "true"
        if (unavailable) cell.dataset.unavailable = "true"
        cell.className = rtl(root, cls)
        cell.dataset.rac = ""
        const inner = document.createElement("div")
        const isRange = !!s.range
        inner.dataset.selectedSingle = String(selected && !isRange)
        inner.dataset.rangeStart = String(start && isRange)
        inner.dataset.rangeEnd = String(stop && isRange)
        if (!selected || isRange) inner.dataset.rangeMiddle = String(selected && !start && !stop && isRange)
        inner.className = dayClass
        inner.textContent = String(day.getUTCDate())
        if (prices && !outside) {
          const span = document.createElement("span")
          const wd = day.getUTCDay()
          span.textContent = wd === 0 || wd === 6 ? "$120" : "$100"
          inner.append(span)
        }
        cell.append(inner)
        td.append(cell)
        tr.append(td)
      }
      tbody.append(tr)
    }
  })
}

document.addEventListener("click", (e) => {
  if (!(e.target instanceof Element)) return
  const preset = e.target.closest("[data-aria-calendar-preset]")
  if (preset) {
    const root = preset.closest("[data-aria-calendar-presets]")?.querySelector("[data-aria-calendar]")
    if (!root) return
    const d = addDays(today(), +preset.dataset.ariaCalendarPreset)
    root.dataset.value = iso(d)
    root.dataset.month = iso(addMonths(d, 0))
    render(root)
    return
  }
  const root = e.target.closest("[data-aria-calendar]")
  // The Persian (calendar-hijri) month is rendered by the server only.
  if (!root || root.dataset.locale === "fa") return
  const nav = e.target.closest("button[slot]")
  if (nav) {
    const s = state(root)
    const dir = nav.getAttribute("slot") === "next" ? 1 : -1
    root.dataset.month = iso(addMonths(s.month, dir * s.months))
    render(root)
    return
  }
  const cell = e.target.closest("[data-date]")
  if (!cell || cell.getAttribute("aria-disabled") === "true") return
  const d = parse(cell.dataset.date)
  if (root.dataset.range !== undefined) {
    const anchor = parse(root.dataset.anchor)
    if (anchor) {
      const [a, b] = anchor <= d ? [anchor, d] : [d, anchor]
      root.dataset.range = `${iso(a)} ${iso(b)}`
      delete root.dataset.anchor
    } else {
      root.dataset.anchor = iso(d)
      root.dataset.range = `${iso(d)} ${iso(d)}`
    }
  } else {
    root.dataset.value = iso(d)
  }
  render(root)
})
