// calendar.js: shadcn/ui Calendar behavior, a port of react-day-picker
// 9.8.1 (DayPicker, useCalendar, useFocus, useSelection and helpers) to
// vanilla JS.
//
// The server renders the calendar (calendar.go) and puts its configuration
// in a data-calendar attribute on the root (removed on mount): the props
// that matter (mode, selected, month, numberOfMonths, captionLayout, ...),
// the class names shadcn computes and the calendar system. Dates are
// carried as Julian Day Numbers (whole days, no time zone), which the
// Gregorian and Persian (date-fns-jalali) systems below convert.
//
// On every state change the calendar is rendered to a light virtual tree
// and reconciled with the DOM the way React does (same keys, same element
// reuse), so focus, hover and CSS transitions behave like upstream.
//
// Events dispatched on the root ([data-slot=calendar]):
//   "shadcn:select"        detail { selected, date } after a day is picked;
//                          selected is an ISO date "YYYY-MM-DD", an array of
//                          them, or { from, to } (null when cleared).
//   "shadcn:month-change"  detail { month } (ISO date of the first day).
import { define, emit } from "./core.js"

// --- calendar systems ----------------------------------------------------------

const div = (a, b) => ~~(a / b)
const mod = (a, b) => a - ~~(a / b) * b
const fdiv = (a, b) => Math.floor(a / b)
const fmod = (a, b) => ((a % b) + b) % b

// Julian Day Number of a proleptic Gregorian date (month 1-12).
function g2d(gy, gm, gd) {
  return (
    div(1461 * (gy + 4800 + div(gm - 14, 12)), 4) +
    div(367 * (gm - 2 - 12 * div(gm - 14, 12)), 12) -
    div(3 * div(gy + 4900 + div(gm - 14, 12), 100), 4) +
    gd -
    32075
  )
}

function d2g(jdn) {
  let L = jdn + 68569
  const n = div(4 * L, 146097)
  L = L - div(146097 * n + 3, 4)
  const i = div(4000 * (L + 1), 1461001)
  L = L - div(1461 * i, 4) + 31
  const j = div(80 * L, 2447)
  const gd = L - div(2447 * j, 80)
  L = div(j, 11)
  const gm = j + 2 - 12 * L
  const gy = 100 * (n - 49) + i + L
  return [gy, gm, gd]
}

const PERSIAN_EPOCH = 1948320
const PERSIAN_NUM_DAYS = [0, 31, 62, 93, 124, 155, 186, 216, 246, 276, 306, 336]

function j2d(jy, jm, jd) {
  const month = jm - 1
  let julianDay = PERSIAN_EPOCH - 1 + 365 * (jy - 1) + div(8 * jy + 21, 33)
  if (month !== 0) julianDay += PERSIAN_NUM_DAYS[month]
  return julianDay + jd
}

function d2j(julianDay) {
  const daysSinceEpoch = julianDay - PERSIAN_EPOCH
  let year = 1 + div(33 * daysSinceEpoch + 3, 12053)
  let dayOfYear = daysSinceEpoch - (365 * (year - 1) + div(8 * year + 21, 33))
  if (dayOfYear < 0) {
    year--
    dayOfYear = daysSinceEpoch - (365 * (year - 1) + div(8 * year + 21, 33))
  }
  const month = dayOfYear < 216 ? div(dayOfYear, 31) : div(dayOfYear - 6, 30)
  return [year, month + 1, dayOfYear - PERSIAN_NUM_DAYS[month] + 1]
}

function isLeapJalaliYear(jy) {
  const m = mod(25 * jy + 11, 33)
  return (m < 8 && m >= -1) || m <= -27
}

// A system maps a JDN to [year, month 0-11, day] and back.
const gregory = {
  ymd(j) {
    const [y, m, d] = d2g(j)
    return [y, m - 1, d]
  },
  jdn(y, m, d) {
    return g2d(y + fdiv(m, 12), fmod(m, 12) + 1, d)
  },
  dim(y, m) {
    const ny = y + fdiv(m, 12)
    const nm = fmod(m, 12)
    return g2d(ny + fdiv(nm + 1, 12), fmod(nm + 1, 12) + 1, 1) - g2d(ny, nm + 1, 1)
  },
}

const persian = {
  ymd(j) {
    const [y, m, d] = d2j(j)
    return [y, m - 1, d]
  },
  jdn(y, m, d) {
    return j2d(y + fdiv(m, 12), fmod(m, 12) + 1, d)
  },
  dim(y, m) {
    const ny = y + fdiv(m, 12)
    const nm = fmod(m, 12)
    if (nm < 6) return 31
    if (nm < 11) return 30
    return isLeapJalaliYear(ny) ? 30 : 29
  },
}

// --- locales -------------------------------------------------------------------------

const EN_MONTHS = ["January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"]
const EN_MONTHS_SHORT = ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"]
const EN_DAYS = ["Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"]
const EN_DAYS_SHORT = ["Su", "Mo", "Tu", "We", "Th", "Fr", "Sa"]
const FA_MONTHS = ["فروردین", "اردیبهشت", "خرداد", "تیر", "مرداد", "شهریور", "مهر", "آبان", "آذر", "دی", "بهمن", "اسفند"]
const FA_DAYS = ["یک‌شنبه", "دوشنبه", "سه‌شنبه", "چهارشنبه", "پنج‌شنبه", "جمعه", "شنبه"]
const FA_DAYS_SHORT = ["1ش", "2ش", "3ش", "4ش", "5ش", "ج", "ش"]
// react-day-picker/locale arSA (Gregorian, Latin digits, weeks from Sunday).
const AR_MONTHS = ["يناير", "فبراير", "مارس", "أبريل", "مايو", "يونيو", "يوليو", "أغسطس", "سبتمبر", "أكتوبر", "نوفمبر", "ديسمبر"]
const AR_DAYS = ["الأحد", "الاثنين", "الثلاثاء", "الأربعاء", "الخميس", "الجمعة", "السبت"]
const AR_DAYS_SHORT = ["أحد", "اثنين", "ثلاثاء", "أربعاء", "خميس", "جمعة", "سبت"]

function enOrdinal(n) {
  const r100 = n % 100
  if (r100 > 20 || r100 < 10) {
    switch (r100 % 10) {
      case 1:
        return n + "st"
      case 2:
        return n + "nd"
      case 3:
        return n + "rd"
    }
  }
  return n + "th"
}

const pad = (n, len) => String(n).padStart(len, "0")

const NUMERALS = {
  latn: null,
  arab: "٠١٢٣٤٥٦٧٨٩",
  arabext: "۰۱۲۳۴۵۶۷۸۹",
  deva: "०१२३४५६७८९",
  beng: "০১২৩৪৫৬৭৮৯",
}

// --- the date library (DateLib over JDNs) -------------------------------------------

function makeLib(cfg) {
  const S = cfg.calendar === "persian" ? persian : gregory
  const fa = cfg.locale === "fa"
  const ar = cfg.locale === "ar"
  const digits = NUMERALS[cfg.numerals ?? "latn"] ?? null
  const num = (s) => (digits ? String(s).replace(/\d/g, (d) => digits[d]) : String(s))
  const weekStartsOn = cfg.weekStartsOn ?? (fa ? 6 : 0)
  const weekday = (j) => fmod(j + 1, 7)
  const lib = {
    S,
    ymd: (j) => S.ymd(j),
    weekday,
    startOfMonth(j) {
      const [y, m] = S.ymd(j)
      return S.jdn(y, m, 1)
    },
    endOfMonth(j) {
      const [y, m] = S.ymd(j)
      return S.jdn(y, m, S.dim(y, m))
    },
    startOfYear(j) {
      return S.jdn(S.ymd(j)[0], 0, 1)
    },
    endOfYear(j) {
      const y = S.ymd(j)[0]
      return S.jdn(y, 11, S.dim(y, 11))
    },
    addDays: (j, n) => j + n,
    addWeeks: (j, n) => j + 7 * n,
    addMonths(j, n) {
      const [y, m, d] = S.ymd(j)
      const ny = y + fdiv(m + n, 12)
      const nm = fmod(m + n, 12)
      return S.jdn(ny, nm, Math.min(d, S.dim(ny, nm)))
    },
    addYears(j, n) {
      return lib.addMonths(j, 12 * n)
    },
    setMonth(j, month) {
      const [y, , d] = S.ymd(j)
      return S.jdn(y, month, Math.min(d, S.dim(y, month)))
    },
    setYear(j, year) {
      const [, m, d] = S.ymd(j)
      return S.jdn(year, m, Math.min(d, S.dim(year, m)))
    },
    getMonth: (j) => S.ymd(j)[1],
    getYear: (j) => S.ymd(j)[0],
    diffMonths(a, b) {
      const [ya, ma] = S.ymd(a)
      const [yb, mb] = S.ymd(b)
      return ya * 12 + ma - (yb * 12 + mb)
    },
    isSameMonth(a, b) {
      return lib.diffMonths(a, b) === 0
    },
    isSameYear: (a, b) => S.ymd(a)[0] === S.ymd(b)[0],
    startOfWeek: (j) => j - fmod(weekday(j) - weekStartsOn, 7),
    endOfWeek: (j) => j - fmod(weekday(j) - weekStartsOn, 7) + 6,
    getWeek(j) {
      const y = S.ymd(j)[0]
      const sow = (yy) => lib.startOfWeek(S.jdn(yy, 0, 1))
      let start = sow(y)
      if (j >= sow(y + 1)) start = sow(y + 1)
      else if (j < start) start = sow(y - 1)
      return fdiv(lib.startOfWeek(j) - start, 7) + 1
    },
    // Formats.
    num,
    caption(j) {
      const [y, m] = S.ymd(j)
      return num(`${fa ? FA_MONTHS[m] : ar ? AR_MONTHS[m] : EN_MONTHS[m]} ${y}`)
    },
    fullDate(j) {
      const [y, m, d] = S.ymd(j)
      const w = weekday(j)
      if (fa) return num(`${FA_DAYS[w]} ${d}-ام ${FA_MONTHS[m]} ${y}`)
      if (ar) return num(`${AR_DAYS[w]}, ${AR_MONTHS[m]} ${d}, ${y}`)
      return num(`${EN_DAYS[w]}, ${EN_MONTHS[m]} ${enOrdinal(d)}, ${y}`)
    },
    weekdayShort: (j) => num((fa ? FA_DAYS_SHORT : ar ? AR_DAYS_SHORT : EN_DAYS_SHORT)[weekday(j)]),
    weekdayLong: (j) => num((fa ? FA_DAYS : ar ? AR_DAYS : EN_DAYS)[weekday(j)]),
    day: (j) => num(S.ymd(j)[2]),
    isoDay(j) {
      const [y, m, d] = S.ymd(j)
      return num(`${pad(y, 4)}-${pad(m + 1, 2)}-${pad(d, 2)}`)
    },
    isoMonth(j) {
      const [y, m] = S.ymd(j)
      return num(`${pad(y, 4)}-${pad(m + 1, 2)}`)
    },
    year: (j) => num(pad(S.ymd(j)[0], 4)),
    // shadcn's formatMonthDropdown: date.toLocaleString("default", { month: "short" }).
    monthDropdown(j) {
      const [, gm] = d2g(j)
      if (ar) return AR_MONTHS[gm - 1]
      return (cfg.monthLong && !fa ? EN_MONTHS : EN_MONTHS_SHORT)[gm - 1]
    },
    // day.date.toLocaleDateString() (en-US).
    localeDate(j) {
      const [gy, gm, gd] = d2g(j)
      if (ar) return `${gd}\u200f/${gm}\u200f/${gy}`.replace(/\d/g, (c) => NUMERALS.arab[c])
      return `${gm}/${gd}/${gy}`
    },
    weekNumber: (n) => (n < 10 ? num(`0${n}`) : num(`${n}`)),
  }
  return lib
}

function todayJdn() {
  const now = new Date()
  return g2d(now.getFullYear(), now.getMonth() + 1, now.getDate())
}

function isoOf(j) {
  const [y, m, d] = d2g(j)
  return `${pad(y, 4)}-${pad(m, 2)}-${pad(d, 2)}`
}

// --- matchers ---------------------------------------------------------------------------

function matches(lib, j, matchers) {
  if (!matchers) return false
  return matchers.some((m) => {
    if (m === true || m === false) return m
    if (typeof m === "number") return m === j
    if (m.from !== undefined || m.to !== undefined) return rangeIncludes(m, j, false)
    if (m.dayOfWeek !== undefined) return [].concat(m.dayOfWeek).includes(lib.weekday(j))
    if (m.before !== undefined && m.after !== undefined) {
      const isDayBefore = m.before - j > 0
      const isDayAfter = m.after - j < 0
      if (m.before > m.after) return isDayAfter && isDayBefore
      return isDayBefore || isDayAfter
    }
    if (m.after !== undefined) return j - m.after > 0
    if (m.before !== undefined) return m.before - j > 0
    return false
  })
}

function rangeIncludes(range, j, excludeEnds) {
  let { from, to } = range
  if (from != null && to != null) {
    if (to - from < 0) [from, to] = [to, from]
    return j - from >= (excludeEnds ? 1 : 0) && to - j >= (excludeEnds ? 1 : 0)
  }
  if (!excludeEnds && to != null) return to === j
  if (!excludeEnds && from != null) return from === j
  return false
}

function addToRange(j, range, min = 0, max = 0, required = false) {
  const from = range?.from ?? null
  const to = range?.to ?? null
  let out
  if (from == null && to == null) {
    out = { from: j, to: min > 0 ? null : j }
  } else if (from != null && to == null) {
    if (from === j) out = required ? { from, to: null } : undefined
    else if (j < from) out = { from: j, to: from }
    else out = { from, to: j }
  } else if (from != null && to != null) {
    if (from === j && to === j) out = required ? { from, to } : undefined
    else if (from === j) out = { from, to: min > 0 ? null : j }
    else if (to === j) out = { from: j, to: min > 0 ? null : j }
    else if (j < from) out = { from: j, to }
    else if (j > from) out = { from, to: j }
    else if (j > to) out = { from, to: j }
  }
  if (out && out.from != null && out.to != null) {
    const diff = out.to - out.from
    if (max > 0 && diff > max) out = { from: j, to: null }
    else if (min > 1 && diff < min) out = { from: j, to: null }
  }
  return out
}

// --- virtual tree and React-like reconciliation -------------------------------------

function v(tag, attrs, children, key) {
  return { tag, attrs: attrs ?? {}, children: children ?? [], key }
}

const SVG_NS = "http://www.w3.org/2000/svg"

function create(node, svg) {
  if (typeof node === "string") return document.createTextNode(node)
  const isSvg = svg || node.tag === "svg"
  const el = isSvg ? document.createElementNS(SVG_NS, node.tag) : document.createElement(node.tag)
  el.__vkey = node.key
  patchAttrs(el, {}, node.attrs)
  el.__vattrs = node.attrs
  for (const c of node.children) el.appendChild(create(c, isSvg))
  if (node.tag === "select" && node.attrs[".value"] !== undefined) el.value = node.attrs[".value"]
  return el
}

function patchAttrs(el, prev, next) {
  for (const name of Object.keys(prev)) {
    if (name[0] === "." || name in next) continue
    if (next[name] === undefined) el.removeAttribute(name)
  }
  for (const [name, value] of Object.entries(next)) {
    if (name[0] === ".") continue
    if (value === undefined || value === null || value === false) {
      if (el.hasAttribute(name)) el.removeAttribute(name)
    } else {
      const s = value === true ? "" : String(value)
      if (el.getAttribute(name) !== s) el.setAttribute(name, s)
    }
  }
}

function sameType(el, node) {
  if (typeof node === "string") return el.nodeType === 3
  return el.nodeType === 1 && el.localName === node.tag && el.__vkey === node.key
}

// Reconciles el's children with nodes. Keyed children are matched by key,
// unkeyed ones by position (React's rules); a matched element of the same
// tag is updated in place, anything else is replaced.
function reconcile(el, nodes, svg, created) {
  const old = [...el.childNodes]
  const byKey = new Map()
  old.forEach((c, i) => {
    if (c.nodeType === 1 && c.__vkey !== undefined) byKey.set(c.__vkey, c)
  })
  const used = new Set()
  const next = nodes.map((node, i) => {
    let match = null
    if (typeof node !== "string" && node.key !== undefined) match = byKey.get(node.key) ?? null
    else {
      const c = old[i]
      if (c && !used.has(c) && (c.nodeType === 3 ? typeof node === "string" : c.__vkey === undefined)) match = c
    }
    if (match && sameType(match, node)) {
      used.add(match)
      update(match, node, svg, created)
      return match
    }
    const fresh = create(node, svg)
    markCreated(fresh, created)
    return fresh
  })
  for (const c of old) if (!used.has(c)) c.remove()
  next.forEach((c, i) => {
    if (el.childNodes[i] !== c) el.insertBefore(c, el.childNodes[i] ?? null)
  })
}

function markCreated(el, created) {
  if (el.nodeType !== 1) return
  created.add(el)
  el.querySelectorAll("*").forEach((c) => created.add(c))
}

function update(el, node, svg, created) {
  if (typeof node === "string") {
    if (el.data !== node) el.data = node
    return
  }
  const isSvg = svg || node.tag === "svg"
  patchAttrs(el, el.__vattrs ?? {}, node.attrs)
  el.__vattrs = node.attrs
  reconcile(el, node.children, isSvg, created)
  if (node.tag === "select" && node.attrs[".value"] !== undefined && el.value !== node.attrs[".value"]) el.value = node.attrs[".value"]
}

// Adopts server-rendered DOM: records keys and attributes on existing
// elements so the first reconciliation updates them in place.
function adopt(el, node) {
  if (typeof node === "string" || el.nodeType !== 1) return
  el.__vkey = node.key
  el.__vattrs = node.attrs
  const kids = [...el.childNodes]
  node.children.forEach((c, i) => kids[i] && adopt(kids[i], c))
}

// --- lucide chevrons (shadcn's Chevron component) -------------------------------------

const CHEVRONS = {
  left: ["chevron-left", "m15 18-6-6 6-6"],
  right: ["chevron-right", "m9 18 6-6-6-6"],
  down: ["chevron-down", "m6 9 6 6 6-6"],
}

function chevron(cfg, orientation, size, disabled) {
  const [name, d] = CHEVRONS[orientation]
  return v(
    "svg",
    {
      xmlns: SVG_NS,
      width: String(size),
      height: String(size),
      viewBox: "0 0 24 24",
      fill: "none",
      stroke: "currentColor",
      "stroke-width": "2",
      "stroke-linecap": "round",
      "stroke-linejoin": "round",
      class: `lucide lucide-${name} ${cfg.chevronClass}`,
      disabled: disabled ? true : undefined,
    },
    [v("path", { d })],
  )
}

// --- the DayPicker -----------------------------------------------------------------------

const SR_ONLY_STYLE = "border: 0px; clip: rect(0px, 0px, 0px, 0px); height: 1px; margin: -1px; overflow: hidden; padding: 0px; position: absolute; width: 1px; white-space: nowrap; overflow-wrap: normal;"

class DayPicker {
  constructor(root, cfg) {
    this.root = root
    this.cfg = cfg
    this.lib = makeLib(cfg)
    this.today = cfg.today ?? todayJdn()
    this.selected = cfg.selected ?? undefined
    this.focused = undefined
    this.lastFocused = undefined
    this.focusedEffect = new WeakMap()
    const [navStart, navEnd] = this.navMonths()
    this.navStart = navStart
    this.navEnd = navEnd
    this.firstMonth = this.initialMonth()
  }

  navMonths() {
    const { cfg, lib } = this
    let start = cfg.startMonth ?? undefined
    let end = cfg.endMonth ?? undefined
    const hasYearDropdown = cfg.captionLayout === "dropdown" || cfg.captionLayout === "dropdown-years"
    if (start !== undefined) start = lib.startOfMonth(start)
    else if (hasYearDropdown) start = lib.startOfYear(lib.addYears(this.today, -100))
    if (end !== undefined) end = lib.endOfMonth(end)
    else if (hasYearDropdown) end = lib.endOfYear(this.today)
    return [start, end]
  }

  initialMonth() {
    const { cfg, lib, navStart, navEnd } = this
    const n = cfg.numberOfMonths ?? 1
    let initial = cfg.month ?? cfg.defaultMonth ?? this.today
    if (navEnd !== undefined && lib.diffMonths(navEnd, initial) < n - 1) initial = lib.addMonths(navEnd, -1 * (n - 1))
    if (navStart !== undefined && lib.diffMonths(initial, navStart) < 0) initial = navStart
    return lib.startOfMonth(initial)
  }

  // useCalendar
  calendar() {
    const { cfg, lib } = this
    const n = cfg.numberOfMonths ?? 1
    const displayMonths = []
    for (let i = 0; i < n; i++) {
      const m = lib.addMonths(this.firstMonth, i)
      if (this.navEnd !== undefined && m > this.navEnd) break
      displayMonths.push(m)
    }
    const maxDate = cfg.endMonth != null ? lib.endOfMonth(cfg.endMonth) : undefined
    const first = displayMonths[0]
    const last = displayMonths[displayMonths.length - 1]
    const startDate = lib.startOfWeek(first)
    const endDate = lib.endOfWeek(lib.endOfMonth(last))
    const dates = []
    for (let j = startDate; j <= endDate; j++) {
      if (maxDate !== undefined && j > maxDate) break
      dates.push(j)
    }
    const nMonths = lib.diffMonths(last, first) + 1
    if (cfg.fixedWeeks && dates.length < 42 * nMonths) {
      const add = 42 * nMonths - dates.length
      for (let i = 0; i < add; i++) dates.push(dates[dates.length - 1] + 1)
    }
    const months = displayMonths.map((month) => {
      const firstOfWeek = lib.startOfWeek(month)
      const lastOfWeek = lib.endOfWeek(lib.endOfMonth(month))
      const monthDates = dates.filter((d) => d >= firstOfWeek && d <= lastOfWeek)
      if (cfg.fixedWeeks && monthDates.length < 42) {
        const add = 42 - monthDates.length
        monthDates.push(...dates.filter((d) => d > lastOfWeek && d <= lastOfWeek + add))
      }
      const weeks = []
      for (const date of monthDates) {
        const weekNumber = lib.getWeek(date)
        let week = weeks.find((w) => w.weekNumber === weekNumber)
        const day = { date, displayMonth: month, outside: !lib.isSameMonth(date, month) }
        if (!week) weeks.push((week = { weekNumber, days: [] }))
        week.days.push(day)
      }
      return { date: month, weeks }
    })
    const days = months.flatMap((m) => m.weeks.flatMap((w) => w.days))
    return { months, days, previousMonth: this.previousMonth(), nextMonth: this.nextMonth() }
  }

  previousMonth() {
    const { cfg, lib } = this
    if (cfg.disableNavigation) return undefined
    const offset = cfg.pagedNavigation ? (cfg.numberOfMonths ?? 1) : 1
    const month = lib.startOfMonth(this.firstMonth)
    if (this.navStart === undefined) return lib.addMonths(month, -offset)
    if (lib.diffMonths(month, this.navStart) <= 0) return undefined
    return lib.addMonths(month, -offset)
  }

  nextMonth() {
    const { cfg, lib } = this
    if (cfg.disableNavigation) return undefined
    const n = cfg.numberOfMonths ?? 1
    const offset = cfg.pagedNavigation ? n : 1
    const month = lib.startOfMonth(this.firstMonth)
    if (this.navEnd === undefined) return lib.addMonths(month, offset)
    if (lib.diffMonths(this.navEnd, this.firstMonth) < n) return undefined
    return lib.addMonths(month, offset)
  }

  goToMonth(date) {
    const { cfg, lib } = this
    if (cfg.disableNavigation) return
    let m = lib.startOfMonth(date)
    if (this.navStart !== undefined && m < lib.startOfMonth(this.navStart)) m = lib.startOfMonth(this.navStart)
    if (this.navEnd !== undefined && m > lib.startOfMonth(this.navEnd)) m = lib.startOfMonth(this.navEnd)
    this.firstMonth = m
    emit(this.root, "shadcn:month-change", { month: isoOf(m) })
  }

  // createGetModifiers
  modifiersOf(day) {
    const { cfg, lib } = this
    const navStart = this.navStart !== undefined ? lib.startOfMonth(this.navStart) : undefined
    const navEnd = this.navEnd !== undefined ? lib.endOfMonth(this.navEnd) : undefined
    const outside = day.outside
    const hidden =
      matches(lib, day.date, cfg.hidden) ||
      (navStart !== undefined && day.date < navStart) ||
      (navEnd !== undefined && day.date > navEnd) ||
      (cfg.showOutsideDays === false && outside)
    return {
      focused: false,
      disabled: matches(lib, day.date, cfg.disabled),
      hidden,
      outside,
      today: day.date === this.today,
    }
  }

  isSelected(j) {
    const { cfg } = this
    const s = this.selected
    switch (cfg.mode) {
      case "single":
        return s != null ? s === j : false
      case "multiple":
        return s?.some((d) => d === j) ?? false
      case "range":
        return s && rangeIncludes(s, j, false)
      default:
        return false
    }
  }

  isSelectionRange() {
    return this.cfg.mode === "range" && this.selected != null && typeof this.selected === "object" && !Array.isArray(this.selected)
  }

  select(j) {
    const { cfg } = this
    const s = this.selected
    if (cfg.mode === "single") {
      let next = j
      if (!cfg.required && s != null && s === j) next = undefined
      this.selected = next
    } else if (cfg.mode === "multiple") {
      let next = [...(s ?? [])]
      if (this.isSelected(j)) {
        if (s?.length === cfg.min) return
        if (cfg.required && s?.length === 1) return
        next = s?.filter((d) => d !== j)
      } else if (s?.length === cfg.max) next = [j]
      else next = [...next, j]
      this.selected = next
    } else if (cfg.mode === "range") {
      let next = addToRange(j, s, cfg.min ?? 0, cfg.max ?? 0, Boolean(cfg.required))
      if (cfg.excludeDisabled && cfg.disabled && next?.from != null && next.to != null) {
        for (let d = next.from; d <= next.to; d++) {
          if (matches(this.lib, d, cfg.disabled)) {
            next = { from: j, to: null }
            break
          }
        }
      }
      this.selected = next
    }
    this.remount = true
    const sel = this.selected
    emit(this.root, "shadcn:select", {
      date: isoOf(j),
      selected:
        sel == null ? null : Array.isArray(sel) ? sel.map(isoOf) : typeof sel === "object" ? { from: sel.from != null ? isoOf(sel.from) : null, to: sel.to != null ? isoOf(sel.to) : null } : isoOf(sel),
    })
  }

  // calculateFocusTarget
  focusTarget(days) {
    let target
    let priority = -1
    for (const day of days) {
      const m = this.modifiersOf(day)
      if (m.disabled || m.hidden || m.outside) continue
      if (this.lastFocused && sameDay(this.lastFocused, day, this.lib) && priority < 2) {
        target = day
        priority = 2
      } else if (this.isSelected(day.date) && priority < 1) {
        target = day
        priority = 1
      } else if (m.today && priority < 0) {
        target = day
        priority = 0
      }
    }
    if (!target) target = days.find((d) => {
      const m = this.modifiersOf(d)
      return !m.disabled && !m.hidden && !m.outside
    })
    return target
  }

  // getNextFocus
  nextFocus(moveBy, moveDir, refDay, attempt = 0) {
    if (attempt > 365) return undefined
    const { lib, cfg } = this
    const amount = moveDir === "after" ? 1 : -1
    let date
    switch (moveBy) {
      case "day":
        date = lib.addDays(refDay.date, amount)
        break
      case "week":
        date = lib.addWeeks(refDay.date, amount)
        break
      case "month":
        date = lib.addMonths(refDay.date, amount)
        break
      case "year":
        date = lib.addYears(refDay.date, amount)
        break
      case "startOfWeek":
        date = lib.startOfWeek(refDay.date)
        break
      case "endOfWeek":
        date = lib.endOfWeek(refDay.date)
        break
    }
    if (moveDir === "before" && this.navStart !== undefined) date = Math.max(this.navStart, date)
    else if (moveDir === "after" && this.navEnd !== undefined) date = Math.min(this.navEnd, date)
    const day = { date, displayMonth: date, outside: false }
    if (!matches(lib, date, cfg.disabled) && !matches(lib, date, cfg.hidden)) return day
    return this.nextFocus(moveBy, moveDir, day, attempt + 1)
  }

  moveFocus(moveBy, moveDir) {
    if (!this.focused) return
    const next = this.nextFocus(moveBy, moveDir, this.focused)
    if (!next) return
    // goToDay
    const { days } = this.calendar()
    if (!days.some((d) => sameDay(d, next, this.lib))) this.goToMonth(next.date)
    this.focused = next
  }

  // --- render ---

  render() {
    const { cfg, lib } = this
    const cn = cfg.classNames
    const cal = this.calendar()
    const { months, days, previousMonth, nextMonth } = cal
    const target = this.focusTarget(days)
    const interactive = cfg.mode != null
    const captionLayout = cfg.captionLayout ?? "label"
    const n = cfg.numberOfMonths ?? 1
    const weekdays = []
    const ws = lib.startOfWeek(todayJdn())
    for (let i = 0; i < 7; i++) weekdays.push(ws + i)

    const nav = v("nav", { class: cn.nav, "aria-label": "" }, [
      v(
        "button",
        {
          type: "button",
          class: cn.button_previous,
          tabindex: previousMonth !== undefined ? undefined : "-1",
          "aria-disabled": previousMonth !== undefined ? undefined : "true",
          "aria-label": "Go to the Previous Month",
        },
        [chevron(cfg, "left", 24, previousMonth === undefined)],
      ),
      v(
        "button",
        {
          type: "button",
          class: cn.button_next,
          tabindex: nextMonth !== undefined ? undefined : "-1",
          "aria-disabled": nextMonth !== undefined ? undefined : "true",
          "aria-label": "Go to the Next Month",
        },
        [chevron(cfg, "right", 24, nextMonth === undefined)],
      ),
    ])

    const yearOptions = () => {
      if (this.navStart === undefined || this.navEnd === undefined) return undefined
      const out = []
      const lastYear = lib.endOfYear(this.navEnd)
      for (let y = lib.startOfYear(this.navStart); y < lastYear || lib.isSameYear(y, lastYear); y = lib.addYears(y, 1)) {
        out.push({ value: lib.getYear(y), label: lib.year(y), disabled: false })
      }
      return out
    }

    const dropdown = (kind, options, value, label) => {
      const selectedOption = options?.find((o) => o.value === value)
      return v("span", { "data-disabled": cfg.disableNavigation ? "true" : "false", class: cn.dropdown_root }, [
        v(
          "select",
          {
            class: `${cn.dropdown} ${kind === "months" ? cn.months_dropdown : cn.years_dropdown}`,
            "aria-label": label,
            disabled: cfg.disableNavigation ? true : undefined,
            ".value": String(value),
          },
          (options ?? []).map((o) => v("option", { value: String(o.value), disabled: o.disabled ? true : undefined }, [o.label], `o${o.value}`)),
        ),
        v("span", { class: cn.caption_label, "aria-hidden": "true" }, [...(selectedOption ? [selectedOption.label] : []), chevron(cfg, "down", 18, false)]),
      ])
    }

    const monthNodes = months.map((month, displayIndex) => {
      const caption = lib.caption(month.date)
      let captionChildren
      if (captionLayout.startsWith("dropdown")) {
        const monthOptions = []
        const sy = lib.startOfYear(month.date)
        for (let i = 0; i < 12; i++) {
          const m = lib.addMonths(sy, i)
          const disabled =
            (this.navStart !== undefined && m < lib.startOfMonth(this.navStart)) || (this.navEnd !== undefined && m > lib.startOfMonth(this.navEnd)) || false
          monthOptions.push({ value: lib.getMonth(m), label: lib.monthDropdown(m), disabled })
        }
        captionChildren = [
          v("div", { class: cn.dropdowns }, [
            captionLayout === "dropdown" || captionLayout === "dropdown-months"
              ? dropdown("months", monthOptions, lib.getMonth(month.date), "Choose the Month")
              : v("span", {}, [lib.monthDropdown(month.date)]),
            captionLayout === "dropdown" || captionLayout === "dropdown-years"
              ? dropdown("years", yearOptions(), lib.getYear(month.date), "Choose the Year")
              : v("span", {}, [lib.year(month.date)]),
            v("span", { role: "status", "aria-live": "polite", style: SR_ONLY_STYLE }, [caption]),
          ]),
        ]
      } else {
        captionChildren = [v("span", { class: cn.caption_label, role: "status", "aria-live": "polite" }, [caption])]
      }

      const head = v("thead", { "aria-hidden": "true" }, [
        v("tr", { class: cn.weekdays }, [
          ...(cfg.showWeekNumber ? [v("th", { "aria-label": "Week Number", class: cn.week_number_header, scope: "col" }, [], "wnh")] : []),
          ...weekdays.map((w, i) => v("th", { "aria-label": lib.weekdayLong(w), class: cn.weekday, scope: "col" }, [lib.weekdayShort(w)], `w${i}`)),
        ]),
      ])

      const weeks = month.weeks.map((week) =>
        v(
          "tr",
          { class: cn.week },
          [
            ...(cfg.showWeekNumber
              ? [
                  v("td", { week: "[object Object]", "aria-label": `Week ${week.weekNumber}`, class: cn.week_number, scope: "row", role: "rowheader" }, [
                    v("div", { class: "flex size-(--cell-size) items-center justify-center text-center" }, [lib.weekNumber(week.weekNumber)]),
                  ], "wn"),
                ]
              : []),
            ...week.days.map((day) => this.renderDay(day, target, interactive)),
          ],
          `k${week.weekNumber}`,
        ),
      )

      return v(
        "div",
        { class: cn.month },
        [
          v("div", { class: cn.month_caption }, captionChildren),
          v(
            "table",
            {
              role: "grid",
              "aria-multiselectable": cfg.mode === "multiple" || cfg.mode === "range" ? "true" : "false",
              "aria-label": caption || undefined,
              class: cn.month_grid,
            },
            [...(cfg.hideWeekdays ? [] : [head]), v("tbody", { class: cn.weeks }, weeks)],
          ),
        ],
        `m${displayIndex}`,
      )
    })

    return [v("div", { class: cn.months }, [nav, ...monthNodes])]
  }

  renderDay(day, target, interactive) {
    const { cfg, lib } = this
    const cn = cfg.classNames
    const m = this.modifiersOf(day)
    m.focused = !m.hidden && Boolean(this.focused && sameDay(this.focused, day, lib))
    m.selected = this.isSelected(day.date) || undefined
    if (this.isSelectionRange()) {
      const { from, to } = this.selected
      m.range_start = Boolean(from != null && to != null && day.date === from)
      m.range_end = Boolean(from != null && to != null && day.date === to)
      m.range_middle = rangeIncludes(this.selected, day.date, true)
    }
    const classes = [cn.day]
    for (const [k, active] of Object.entries(m)) if (active === true && cn[k]) classes.push(cn[k])
    for (const mod of cfg.modifiers ?? []) if (mod.class && matches(lib, day.date, mod.match)) classes.push(mod.class)
    const attr = (b) => (b ? "true" : undefined)
    const dataBool = (b) => (b === undefined ? undefined : b ? "true" : "false")
    const label = (() => {
      let l = lib.fullDate(day.date)
      if (m.today) l = `Today, ${l}`
      return l
    })()
    let children = []
    if (!m.hidden && interactive) {
      const selectedSingle = m.selected && !m.range_start && !m.range_end && !m.range_middle
      const btnLabel = m.selected ? `${label}, selected` : label
      children = [
        v(
          "button",
          {
            "data-slot": "button",
            "data-variant": cfg.dayButtonVariant ?? "ghost",
            "data-size": "icon",
            class: cfg.dayButtonClass,
            "data-day": lib.localeDate(day.date),
            "data-selected-single": dataBool(selectedSingle),
            "data-range-start": dataBool(m.range_start),
            "data-range-end": dataBool(m.range_end),
            "data-range-middle": dataBool(m.range_middle),
            type: "button",
            disabled: m.disabled ? true : undefined,
            tabindex: target && sameDay(target, day, lib) ? "0" : "-1",
            "aria-label": btnLabel,
          },
          [lib.day(day.date), ...(!day.outside && cfg.dayExtra?.[lib.weekday(day.date)] ? [v("span", {}, [cfg.dayExtra[lib.weekday(day.date)]])] : [])],
        ),
      ]
    } else if (!m.hidden) {
      children = [lib.day(day.date)]
    }
    return v(
      "td",
      {
        class: classes.join(" "),
        role: "gridcell",
        "aria-selected": m.selected ? "true" : undefined,
        "aria-label": !interactive && !m.hidden ? label : undefined,
        "data-day": lib.isoDay(day.date),
        "data-month": day.outside ? lib.isoMonth(day.date) : undefined,
        "data-selected": attr(m.selected),
        "data-disabled": attr(m.disabled),
        "data-hidden": attr(m.hidden),
        "data-outside": attr(day.outside),
        "data-focused": attr(m.focused),
        "data-today": attr(m.today),
      },
      children,
      `${day.date}_${day.displayMonth}`,
    )
  }
}

function sameDay(a, b, lib) {
  return a.date === b.date && lib.isSameMonth(a.displayMonth, b.displayMonth)
}

function jdnOf(iso) {
  const [y, m, d] = iso.split("-").map(Number)
  return g2d(y, m, d)
}

/**
 * Updates the props of a calendar that has not mounted yet (for instance
 * one inside a popover's <template>, which mounts a fresh copy every time
 * it opens, like React remounts it) and re-renders its markup: selected (ISO dates in the shape of
 * the "shadcn:select" detail, or null) and defaultMonth (an ISO date).
 * @param {Element | DocumentFragment} el the calendar root or an ancestor
 * @param {{ selected?: any, defaultMonth?: string | null }} props
 */
export function setCalendarProps(el, props) {
  // A mounted calendar re-renders in place.
  const liveRoot = live.has(el) ? el : [...(el.querySelectorAll?.('[data-slot="calendar"]') ?? [])].find((r) => live.has(r))
  if (liveRoot) {
    live.get(liveRoot)(props)
    return
  }
  const root = el instanceof Element && el.hasAttribute("data-calendar") ? el : el.querySelector("[data-calendar]")
  if (!root) return
  const cfg = JSON.parse(root.getAttribute("data-calendar"))
  if ("selected" in props) {
    const s = props.selected
    if (s == null) delete cfg.selected
    else if (typeof s === "string") cfg.selected = jdnOf(s)
    else if (Array.isArray(s)) cfg.selected = s.map(jdnOf)
    else cfg.selected = { from: s.from ? jdnOf(s.from) : null, to: s.to ? jdnOf(s.to) : null }
  }
  if ("defaultMonth" in props) {
    if (props.defaultMonth) cfg.defaultMonth = jdnOf(props.defaultMonth)
    else delete cfg.defaultMonth
  }
  root.setAttribute("data-calendar", JSON.stringify(cfg))
  // Re-render the markup the calendar will mount from.
  const dp = new DayPicker(root, cfg)
  root.replaceChildren(...dp.render().map((n) => create(n, false)))
}

// Mounted calendar roots -> their setCalendarProps handler.
const live = new WeakMap()

function selectedOf(s) {
  if (s == null) return undefined
  if (typeof s === "string") return jdnOf(s)
  if (Array.isArray(s)) return s.map(jdnOf)
  return { from: s.from ? jdnOf(s.from) : null, to: s.to ? jdnOf(s.to) : null }
}

function readConfig(root) {
  const raw = root.getAttribute("data-calendar")
  root.removeAttribute("data-calendar")
  if (!raw) return null
  try {
    return JSON.parse(raw)
  } catch {
    return null
  }
}

define("calendar", {
  selector: '[data-slot="calendar"]',
  mount(root) {
    const cfg = readConfig(root)
    if (!cfg) return
    const dp = new DayPicker(root, cfg)
    const prevFocused = new WeakMap() // button -> its focused modifier at the last render

    const render = () => {
      const tree = dp.render()
      const created = new Set()
      const focusedThisRender = dp.focused
      // Buttons removed by this render blur synchronously; React has already
      // detached their handlers, so those focusout events are ignored.
      dp.reconciling = true
      try {
        if (dp.remount) {
          // A selection re-renders the parent upstream (onSelect sets its
          // state), and shadcn's Calendar passes new component functions to
          // DayPicker on every render: everything under the root remounts.
          dp.remount = false
          root.replaceChildren()
        }
        reconcile(root, tree, false, created)
      } finally {
        dp.reconciling = false
      }
      dp.committedFocused = focusedThisRender
      // DayButton's effect: focus when modifiers.focused becomes true.
      for (const btn of root.querySelectorAll("td > button")) {
        const td = btn.parentElement
        const focused = td.getAttribute("data-focused") === "true"
        const before = created.has(btn) ? false : (prevFocused.get(btn) ?? false)
        prevFocused.set(btn, focused)
        if (focused && !before) btn.focus()
      }
    }

    // First render: adopt the server-rendered DOM.
    adopt(root, v("div", {}, dp.render()))
    {
      const tree = dp.render()
      reconcile(root, tree, false, new Set())
      for (const btn of root.querySelectorAll("td > button")) prevFocused.set(btn, false)
      dp.committedFocused = dp.focused
    }

    const dayOf = (btn) => {
      const td = btn.closest("td")
      const key = td?.__vkey
      if (!key) return null
      const [date, displayMonth] = key.split("_").map(Number)
      return { date, displayMonth, outside: !dp.lib.isSameMonth(date, displayMonth) }
    }
    const dayButton = (target) => {
      const btn = target instanceof Element ? target.closest("td > button") : null
      return btn && root.contains(btn) ? btn : null
    }

    // React flushes discrete updates synchronously at the end of the
    // event; updates made while rendering (focus moved by the DayButton
    // effect fires focusin/focusout) are batched into one more render.
    let rendering = false
    let dirty = false
    const schedule = () => {
      if (rendering) {
        dirty = true
        return
      }
      rendering = true
      try {
        do {
          dirty = false
          render()
        } while (dirty)
      } finally {
        rendering = false
      }
    }

    const onClick = (e) => {
      const nav = e.target instanceof Element ? e.target.closest(".rdp-button_previous, .rdp-button_next") : null
      if (nav && root.contains(nav)) {
        const { previousMonth, nextMonth } = dp.calendar()
        if (nav.classList.contains("rdp-button_previous")) {
          if (previousMonth === undefined) return
          dp.goToMonth(previousMonth)
        } else {
          if (nextMonth === undefined) return
          dp.goToMonth(nextMonth)
        }
        schedule()
        return
      }
      const btn = dayButton(e.target)
      if (!btn) return
      const day = dayOf(btn)
      if (!day) return
      e.preventDefault()
      e.stopPropagation()
      dp.focused = day
      dp.select(day.date)
      schedule()
    }
    const onFocusIn = (e) => {
      const btn = dayButton(e.target)
      if (!btn) return
      const day = dayOf(btn)
      if (!day) return
      dp.focused = day
      schedule()
    }
    const onFocusOut = (e) => {
      const btn = dayButton(e.target)
      if (!btn || dp.reconciling) return
      // React runs the handler of the last committed render.
      dp.lastFocused = dp.committedFocused
      dp.focused = undefined
      schedule()
    }
    const onKeyDown = (e) => {
      const btn = dayButton(e.target)
      if (!btn) return
      const rtl = cfg.dir === "rtl"
      const map = {
        ArrowLeft: [e.shiftKey ? "month" : "day", rtl ? "after" : "before"],
        ArrowRight: [e.shiftKey ? "month" : "day", rtl ? "before" : "after"],
        ArrowDown: [e.shiftKey ? "year" : "week", "after"],
        ArrowUp: [e.shiftKey ? "year" : "week", "before"],
        PageUp: [e.shiftKey ? "year" : "month", "before"],
        PageDown: [e.shiftKey ? "year" : "month", "after"],
        Home: ["startOfWeek", "before"],
        End: ["endOfWeek", "after"],
      }
      const move = map[e.key]
      if (!move) return
      e.preventDefault()
      e.stopPropagation()
      dp.moveFocus(move[0], move[1])
      schedule()
    }
    const onChange = (e) => {
      const select = e.target
      if (!(select instanceof HTMLSelectElement) || !root.contains(select)) return
      const monthEl = select.closest(".rdp-month")
      const monthIndex = [...root.querySelectorAll(".rdp-month")].indexOf(monthEl)
      const { months } = dp.calendar()
      const month = months[monthIndex]
      if (!month) return
      const value = Number(select.value)
      const start = dp.lib.startOfMonth(month.date)
      if (select.classList.contains("rdp-months_dropdown")) dp.goToMonth(dp.lib.setMonth(start, value))
      else dp.goToMonth(dp.lib.setYear(start, value))
      schedule()
    }

    // setCalendarProps on a mounted calendar: the parent re-renders with new
    // props, which remounts everything under the root (see render).
    live.set(root, (props) => {
      if ("selected" in props) dp.selected = selectedOf(props.selected)
      dp.remount = true
      schedule()
    })
    root.addEventListener("click", onClick)
    root.addEventListener("focusin", onFocusIn)
    root.addEventListener("focusout", onFocusOut)
    root.addEventListener("keydown", onKeyDown)
    root.addEventListener("change", onChange)
    return () => {
      live.delete(root)
      root.removeEventListener("click", onClick)
      root.removeEventListener("focusin", onFocusIn)
      root.removeEventListener("focusout", onFocusOut)
      root.removeEventListener("keydown", onKeyDown)
      root.removeEventListener("change", onChange)
    }
  },
})
