// resizable.js: behavior of shadcn/ui Resizable, a port of
// react-resizable-panels v4.5.8
// (ui/apps/v4/node_modules/react-resizable-panels/dist/react-resizable-panels.js).
//
// The library keeps one global store (pointer interaction state, cursor
// flags and every mounted group with its layout) and listens on the
// document for pointer events; groups, panels and separators subscribe to
// the store. That core is ported function by function below, with the
// library's logic unchanged. The React components (Group, Panel,
// Separator) become DOM mounts that render the same attributes and styles
// at the same moments: updates made during discrete events (pointerdown,
// pointerup, keydown, dblclick) are rendered in a microtask, those made
// during continuous events (pointermove, ...) and timers in a later task,
// like React's lanes.
//
// Markup (server-rendered by resizable.go):
//   [data-group]      group: id, style flex-direction gives the orientation;
//                     optional data-resizable JSON {disabled, disableCursor,
//                     defaultLayout, resizeTargetMinimumSize}
//   [data-panel]      panel (direct child): id; optional data-resizable JSON
//                     {defaultSize, minSize, maxSize, collapsedSize, collapsible}
//   [data-separator]  separator (direct child): id
// data-resizable is read and removed when mounting.
//
// Events on the group element: "shadcn:layout-change" {layout} on every
// layout change, "shadcn:layout-changed" {layout} once a change completes.
// APIs: groupEl.resizableGroup = {getLayout, setLayout},
// panelEl.resizablePanel = {collapse, expand, getSize, isCollapsed, resize}.
import { define, emit } from "./core.js"

// --- library core ---------------------------------------------------------------------

function assert(e, t = "Assertion error") {
  if (!e) throw Error(t)
}

function getGroupSize({ group: e }) {
  const { orientation: t, panels: n } = e
  return n.reduce((o, i) => ((o += t === "horizontal" ? i.element.offsetWidth : i.element.offsetHeight), o), 0)
}

function sortByElementOffset(e, t) {
  return Array.from(t).sort(e === "horizontal" ? sortByOffsetLeft : sortByOffsetTop)
}
function sortByOffsetLeft(e, t) {
  const n = e.element.offsetLeft - t.element.offsetLeft
  return n !== 0 ? n : e.element.offsetWidth - t.element.offsetWidth
}
function sortByOffsetTop(e, t) {
  const n = e.element.offsetTop - t.element.offsetTop
  return n !== 0 ? n : e.element.offsetHeight - t.element.offsetHeight
}

function isHTMLElement(e) {
  return e !== null && typeof e == "object" && "nodeType" in e && e.nodeType === Node.ELEMENT_NODE
}

function getDistanceBetweenPointAndRect(e, t) {
  return {
    x: e.x >= t.left && e.x <= t.right ? 0 : Math.min(Math.abs(e.x - t.left), Math.abs(e.x - t.right)),
    y: e.y >= t.top && e.y <= t.bottom ? 0 : Math.min(Math.abs(e.y - t.top), Math.abs(e.y - t.bottom)),
  }
}

function findClosestRect({ orientation: e, rects: t, targetRect: n }) {
  const o = { x: n.x + n.width / 2, y: n.y + n.height / 2 }
  let i,
    r = Number.MAX_VALUE
  for (const l of t) {
    const { x: s, y: a } = getDistanceBetweenPointAndRect(o, l),
      u = e === "horizontal" ? s : a
    u < r && ((r = u), (i = l))
  }
  return assert(i, "No rect found"), i
}

let coarsePointer
export function isCoarsePointer() {
  return coarsePointer === void 0 && (typeof matchMedia == "function" ? (coarsePointer = !!matchMedia("(pointer:coarse)").matches) : (coarsePointer = !1)), coarsePointer
}

function calculateHitRegions(e) {
  const { element: t, orientation: n, panels: o, separators: i } = e,
    r = sortByElementOffset(n, Array.from(t.children).filter(isHTMLElement).map((c) => ({ element: c }))).map(({ element: c }) => c),
    l = []
  let s = !1,
    a,
    u = []
  for (const c of r)
    if (c.hasAttribute("data-panel")) {
      const m = o.find((y) => y.element === c)
      if (m) {
        if (a) {
          const y = a.element.getBoundingClientRect(),
            h = c.getBoundingClientRect()
          let f
          if (s) {
            const d = n === "horizontal" ? new DOMRect(y.right, y.top, 0, y.height) : new DOMRect(y.left, y.bottom, y.width, 0),
              p = n === "horizontal" ? new DOMRect(h.left, h.top, 0, h.height) : new DOMRect(h.left, h.top, h.width, 0)
            switch (u.length) {
              case 0: {
                f = [d, p]
                break
              }
              case 1: {
                const g = u[0],
                  v = findClosestRect({ orientation: n, rects: [y, h], targetRect: g.element.getBoundingClientRect() })
                f = [g, v === y ? p : d]
                break
              }
              default: {
                f = u
                break
              }
            }
          } else
            u.length
              ? (f = u)
              : (f = [
                  n === "horizontal" ? new DOMRect(y.right, h.top, h.left - y.right, h.height) : new DOMRect(h.left, y.bottom, h.width, h.top - y.bottom),
                ])
          for (const d of f) {
            let p = "width" in d ? d : d.element.getBoundingClientRect()
            const g = isCoarsePointer() ? e.resizeTargetMinimumSize.coarse : e.resizeTargetMinimumSize.fine
            if (p.width < g) {
              const v = g - p.width
              p = new DOMRect(p.x - v / 2, p.y, p.width + v, p.height)
            }
            if (p.height < g) {
              const v = g - p.height
              p = new DOMRect(p.x, p.y - v / 2, p.width, p.height + v)
            }
            l.push({ group: e, groupSize: getGroupSize({ group: e }), panels: [a, m], separator: "width" in d ? void 0 : d, rect: p })
          }
        }
        ;(s = !1), (a = m), (u = [])
      }
    } else if (c.hasAttribute("data-separator")) {
      const m = i.find((y) => y.element === c)
      m ? u.push(m) : ((a = void 0), (u = []))
    } else s = !0
  return l
}

function emToPx(e, t) {
  const n = getComputedStyle(e),
    o = parseFloat(n.fontSize)
  return t * o
}
function remToPx(e, t) {
  const n = getComputedStyle(e.ownerDocument.body),
    o = parseFloat(n.fontSize)
  return t * o
}
function vhToPx(e) {
  return (e / 100) * window.innerHeight
}
function vwToPx(e) {
  return (e / 100) * window.innerWidth
}
function parseSizeAndUnit(e) {
  switch (typeof e) {
    case "number":
      return [e, "px"]
    case "string": {
      const t = parseFloat(e)
      return e.endsWith("%")
        ? [t, "%"]
        : e.endsWith("px")
          ? [t, "px"]
          : e.endsWith("rem")
            ? [t, "rem"]
            : e.endsWith("em")
              ? [t, "em"]
              : e.endsWith("vh")
                ? [t, "vh"]
                : e.endsWith("vw")
                  ? [t, "vw"]
                  : [t, "%"]
    }
  }
}
function sizeStyleToPixels({ groupSize: e, panelElement: t, styleProp: n }) {
  let o
  const [i, r] = parseSizeAndUnit(n)
  switch (r) {
    case "%": {
      o = (i / 100) * e
      break
    }
    case "px": {
      o = i
      break
    }
    case "rem": {
      o = remToPx(t, i)
      break
    }
    case "em": {
      o = emToPx(t, i)
      break
    }
    case "vh": {
      o = vhToPx(i)
      break
    }
    case "vw": {
      o = vwToPx(i)
      break
    }
  }
  return o
}

function formatLayoutNumber(e) {
  return parseFloat(e.toFixed(3))
}

function calculatePanelConstraints(e) {
  const { panels: t } = e,
    n = getGroupSize({ group: e })
  return n === 0
    ? t.map((o) => ({
        collapsedSize: 0,
        collapsible: o.panelConstraints.collapsible === !0,
        defaultSize: void 0,
        minSize: 0,
        maxSize: 100,
        panelId: o.id,
      }))
    : t.map((o) => {
        const { element: i, panelConstraints: r } = o
        let l = 0
        if (r.collapsedSize !== void 0) {
          const c = sizeStyleToPixels({ groupSize: n, panelElement: i, styleProp: r.collapsedSize })
          l = formatLayoutNumber((c / n) * 100)
        }
        let s
        if (r.defaultSize !== void 0) {
          const c = sizeStyleToPixels({ groupSize: n, panelElement: i, styleProp: r.defaultSize })
          s = formatLayoutNumber((c / n) * 100)
        }
        let a = 0
        if (r.minSize !== void 0) {
          const c = sizeStyleToPixels({ groupSize: n, panelElement: i, styleProp: r.minSize })
          a = formatLayoutNumber((c / n) * 100)
        }
        let u = 100
        if (r.maxSize !== void 0) {
          const c = sizeStyleToPixels({ groupSize: n, panelElement: i, styleProp: r.maxSize })
          u = formatLayoutNumber((c / n) * 100)
        }
        return { collapsedSize: l, collapsible: r.collapsible === !0, defaultSize: s, minSize: a, maxSize: u, panelId: o.id }
      })
}

class EventEmitter {
  #e = {}
  addListener(t, n) {
    const o = this.#e[t]
    return (
      o === void 0 ? (this.#e[t] = [n]) : o.includes(n) || o.push(n),
      () => {
        this.removeListener(t, n)
      }
    )
  }
  emit(t, n) {
    const o = this.#e[t]
    if (o !== void 0)
      if (o.length === 1) o[0].call(null, n)
      else {
        let i = !1,
          r = null
        const l = Array.from(o)
        for (let s = 0; s < l.length; s++) {
          const a = l[s]
          try {
            a.call(null, n)
          } catch (u) {
            r === null && ((i = !0), (r = u))
          }
        }
        if (i) throw r
      }
  }
  removeListener(t, n) {
    const o = this.#e[t]
    if (o !== void 0) {
      const i = o.indexOf(n)
      i >= 0 && o.splice(i, 1)
    }
  }
}

function fuzzyNumbersEqual(e, t, n = 0) {
  return Math.abs(formatLayoutNumber(e) - formatLayoutNumber(t)) <= n
}

let store = {
  cursorFlags: 0,
  interactionState: { state: "inactive" },
  mountedGroups: new Map(),
}
const eventEmitter = new EventEmitter()
function read() {
  return store
}
function update(e) {
  const t = typeof e == "function" ? e(store) : e
  if (store === t) return store
  const n = store
  return (
    (store = { ...store, ...t }),
    t.cursorFlags !== void 0 && eventEmitter.emit("cursorFlagsChange", store.cursorFlags),
    t.interactionState !== void 0 && eventEmitter.emit("interactionStateChange", store.interactionState),
    t.mountedGroups !== void 0 &&
      (store.mountedGroups.forEach((o, i) => {
        o.derivedPanelConstraints.forEach((r) => {
          if (r.collapsible) {
            const { layout: l } = n.mountedGroups.get(i) ?? {}
            if (l) {
              const s = fuzzyNumbersEqual(r.collapsedSize, o.layout[r.panelId]),
                a = fuzzyNumbersEqual(r.collapsedSize, l[r.panelId])
              s && !a && (i.inMemoryLastExpandedPanelSizes[r.panelId] = l[r.panelId])
            }
          }
        })
      }),
      eventEmitter.emit("mountedGroupsChange", store.mountedGroups)),
    store
  )
}

function findClosestHitRegion(e, t, n) {
  let o,
    i = { x: 1 / 0, y: 1 / 0 }
  for (const r of t) {
    const l = getDistanceBetweenPointAndRect(n, r.rect)
    switch (e) {
      case "horizontal": {
        l.x <= i.x && ((o = r), (i = l))
        break
      }
      case "vertical": {
        l.y <= i.y && ((o = r), (i = l))
        break
      }
    }
  }
  return o ? { distance: i, hitRegion: o } : void 0
}

function isShadowRoot(e) {
  return e !== null && typeof e == "object" && "nodeType" in e && e.nodeType === Node.DOCUMENT_FRAGMENT_NODE
}

function compareStackingOrder(e, t) {
  if (e === t) throw new Error("Cannot compare node with itself")
  const n = { a: getAncestors(e), b: getAncestors(t) }
  let o
  for (; n.a.at(-1) === n.b.at(-1); ) (o = n.a.pop()), n.b.pop()
  assert(o, "Stacking order can only be calculated for elements with a common ancestor")
  const i = { a: getZIndex(findStackingContext(n.a)), b: getZIndex(findStackingContext(n.b)) }
  if (i.a === i.b) {
    const r = o.childNodes,
      l = { a: n.a.at(-1), b: n.b.at(-1) }
    let s = r.length
    for (; s--; ) {
      const a = r[s]
      if (a === l.a) return 1
      if (a === l.b) return -1
    }
  }
  return Math.sign(i.a - i.b)
}
const willChangeStackingProps = /\b(?:position|zIndex|opacity|transform|webkitTransform|mixBlendMode|filter|webkitFilter|isolation)\b/
function isFlexItem(e) {
  const t = getComputedStyle(getParent(e) ?? e).display
  return t === "flex" || t === "inline-flex"
}
function createsStackingContext(e) {
  const t = getComputedStyle(e)
  return !!(
    t.position === "fixed" ||
    (t.zIndex !== "auto" && (t.position !== "static" || isFlexItem(e))) ||
    +t.opacity < 1 ||
    ("transform" in t && t.transform !== "none") ||
    ("webkitTransform" in t && t.webkitTransform !== "none") ||
    ("mixBlendMode" in t && t.mixBlendMode !== "normal") ||
    ("filter" in t && t.filter !== "none") ||
    ("webkitFilter" in t && t.webkitFilter !== "none") ||
    ("isolation" in t && t.isolation === "isolate") ||
    willChangeStackingProps.test(t.willChange) ||
    t.webkitOverflowScrolling === "touch"
  )
}
function findStackingContext(e) {
  let t = e.length
  for (; t--; ) {
    const n = e[t]
    if ((assert(n, "Missing node"), createsStackingContext(n))) return n
  }
  return null
}
function getZIndex(e) {
  return (e && Number(getComputedStyle(e).zIndex)) || 0
}
function getAncestors(e) {
  const t = []
  for (; e; ) t.push(e), (e = getParent(e))
  return t
}
function getParent(e) {
  const { parentNode: t } = e
  return isShadowRoot(t) ? t.host : t
}

function intersects(e, t) {
  return e.x < t.x + t.width && e.x + e.width > t.x && e.y < t.y + t.height && e.y + e.height > t.y
}

function isViableHitTarget({ groupElement: e, hitRegion: t, pointerEventTarget: n }) {
  if (!isHTMLElement(n) || n.contains(e) || e.contains(n)) return !0
  if (compareStackingOrder(n, e) > 0) {
    let o = n
    for (; o; ) {
      if (o.contains(e)) return !0
      if (intersects(o.getBoundingClientRect(), t)) return !1
      o = o.parentElement
    }
  }
  return !0
}

function findMatchingHitRegions(e, t) {
  const n = []
  return (
    t.forEach((o, i) => {
      if (i.disabled) return
      const r = calculateHitRegions(i),
        l = findClosestHitRegion(i.orientation, r, { x: e.clientX, y: e.clientY })
      l &&
        l.distance.x <= 0 &&
        l.distance.y <= 0 &&
        isViableHitTarget({ groupElement: i.element, hitRegion: l.hitRegion.rect, pointerEventTarget: e.target }) &&
        n.push(l.hitRegion)
    }),
    n
  )
}

function arraysEqual(e, t) {
  if (e.length !== t.length) return !1
  for (let n = 0; n < e.length; n++) if (e[n] != t[n]) return !1
  return !0
}

function compareLayoutNumbers(e, t) {
  return fuzzyNumbersEqual(e, t) ? 0 : e > t ? 1 : -1
}

function validatePanelSize({ panelConstraints: e, size: t }) {
  const { collapsedSize: n = 0, collapsible: o, maxSize: i = 100, minSize: r = 0 } = e
  if (compareLayoutNumbers(t, r) < 0)
    if (o) {
      const l = (n + r) / 2
      compareLayoutNumbers(t, l) < 0 ? (t = n) : (t = r)
    } else t = r
  return (t = Math.min(i, t)), (t = formatLayoutNumber(t)), t
}

function adjustLayoutByDelta({ delta: e, initialLayout: t, panelConstraints: n, pivotIndices: o, prevLayout: i, trigger: r }) {
  if (fuzzyNumbersEqual(e, 0)) return t
  const l = Object.values(t),
    s = Object.values(i),
    a = [...l],
    [u, c] = o
  assert(u != null, "Invalid first pivot index"), assert(c != null, "Invalid second pivot index")
  let m = 0
  switch (r) {
    case "keyboard": {
      {
        const f = e < 0 ? c : u,
          d = n[f]
        assert(d, `Panel constraints not found for index ${f}`)
        const { collapsedSize: p = 0, collapsible: g, minSize: v = 0 } = d
        if (g) {
          const x = l[f]
          if ((assert(x != null, `Previous layout not found for panel index ${f}`), fuzzyNumbersEqual(x, p))) {
            const S = v - x
            compareLayoutNumbers(S, Math.abs(e)) > 0 && (e = e < 0 ? 0 - S : S)
          }
        }
      }
      {
        const f = e < 0 ? u : c,
          d = n[f]
        assert(d, `No panel constraints found for index ${f}`)
        const { collapsedSize: p = 0, collapsible: g, minSize: v = 0 } = d
        if (g) {
          const x = l[f]
          if ((assert(x != null, `Previous layout not found for panel index ${f}`), fuzzyNumbersEqual(x, v))) {
            const S = x - p
            compareLayoutNumbers(S, Math.abs(e)) > 0 && (e = e < 0 ? 0 - S : S)
          }
        }
      }
      break
    }
  }
  {
    const f = e < 0 ? 1 : -1
    let d = e < 0 ? c : u,
      p = 0
    for (;;) {
      const v = l[d]
      assert(v != null, `Previous layout not found for panel index ${d}`)
      const S = validatePanelSize({ panelConstraints: n[d], size: 100 }) - v
      if (((p += S), (d += f), d < 0 || d >= n.length)) break
    }
    const g = Math.min(Math.abs(e), Math.abs(p))
    e = e < 0 ? 0 - g : g
  }
  {
    let d = e < 0 ? u : c
    for (; d >= 0 && d < n.length; ) {
      const p = Math.abs(e) - Math.abs(m),
        g = l[d]
      assert(g != null, `Previous layout not found for panel index ${d}`)
      const v = g - p,
        x = validatePanelSize({ panelConstraints: n[d], size: v })
      if (!fuzzyNumbersEqual(g, x) && ((m += g - x), (a[d] = x), m.toFixed(3).localeCompare(Math.abs(e).toFixed(3), void 0, { numeric: !0 }) >= 0)) break
      e < 0 ? d-- : d++
    }
  }
  if (arraysEqual(s, a)) return i
  {
    const f = e < 0 ? c : u,
      d = l[f]
    assert(d != null, `Previous layout not found for panel index ${f}`)
    const p = d + m,
      g = validatePanelSize({ panelConstraints: n[f], size: p })
    if (((a[f] = g), !fuzzyNumbersEqual(g, p))) {
      let v = p - g,
        S = e < 0 ? c : u
      for (; S >= 0 && S < n.length; ) {
        const w = a[S]
        assert(w != null, `Previous layout not found for panel index ${S}`)
        const C = w + v,
          P = validatePanelSize({ panelConstraints: n[S], size: C })
        if ((fuzzyNumbersEqual(w, P) || ((v -= P - w), (a[S] = P)), fuzzyNumbersEqual(v, 0))) break
        e > 0 ? S-- : S++
      }
    }
  }
  const y = Object.values(a).reduce((f, d) => d + f, 0)
  if (!fuzzyNumbersEqual(y, 100, 0.1)) return i
  const h = Object.keys(i)
  return a.reduce((f, d, p) => ((f[h[p]] = d), f), {})
}

function layoutsEqual(e, t) {
  if (Object.keys(e).length !== Object.keys(t).length) return !1
  for (const n in e) if (t[n] === void 0 || compareLayoutNumbers(e[n], t[n]) !== 0) return !1
  return !0
}

function validatePanelGroupLayout({ layout: e, panelConstraints: t }) {
  const o = [...Object.values(e)],
    i = o.reduce((s, a) => s + a, 0)
  if (o.length !== t.length) throw Error(`Invalid ${t.length} panel layout: ${o.map((s) => `${s}%`).join(", ")}`)
  if (!fuzzyNumbersEqual(i, 100) && o.length > 0)
    for (let s = 0; s < t.length; s++) {
      const a = o[s]
      assert(a != null, `No layout data found for index ${s}`)
      const u = (100 / i) * a
      o[s] = u
    }
  let r = 0
  for (let s = 0; s < t.length; s++) {
    const a = o[s]
    assert(a != null, `No layout data found for index ${s}`)
    const u = validatePanelSize({ panelConstraints: t[s], size: a })
    a != u && ((r += a - u), (o[s] = u))
  }
  if (!fuzzyNumbersEqual(r, 0))
    for (let s = 0; s < t.length; s++) {
      const a = o[s]
      assert(a != null, `No layout data found for index ${s}`)
      const u = a + r,
        c = validatePanelSize({ panelConstraints: t[s], size: u })
      if (a !== c && ((r -= c - a), (o[s] = c), fuzzyNumbersEqual(r, 0))) break
    }
  const l = Object.keys(e)
  return o.reduce((s, a, u) => ((s[l[u]] = a), s), {})
}

function getImperativePanelMethods({ groupId: e, panelId: t }) {
  const n = () => {
      const { mountedGroups: s } = read()
      for (const [a, { defaultLayoutDeferred: u, derivedPanelConstraints: c, layout: m, separatorToPanels: y }] of s)
        if (a.id === e) return { defaultLayoutDeferred: u, derivedPanelConstraints: c, group: a, layout: m, separatorToPanels: y }
      throw Error(`Group ${e} not found`)
    },
    o = () => {
      const s = n().derivedPanelConstraints.find((a) => a.panelId === t)
      if (s !== void 0) return s
      throw Error(`Panel constraints not found for Panel ${t}`)
    },
    i = () => {
      const s = n().group.panels.find((a) => a.id === t)
      if (s !== void 0) return s
      throw Error(`Layout not found for Panel ${t}`)
    },
    r = () => {
      const s = n().layout[t]
      if (s !== void 0) return s
      throw Error(`Layout not found for Panel ${t}`)
    },
    l = (s) => {
      const a = r()
      if (s === a) return
      const { defaultLayoutDeferred: u, derivedPanelConstraints: c, group: m, layout: y, separatorToPanels: h } = n(),
        f = m.panels.findIndex((v) => v.id === t),
        d = f === m.panels.length - 1,
        p = adjustLayoutByDelta({
          delta: d ? a - s : s - a,
          initialLayout: y,
          panelConstraints: c,
          pivotIndices: d ? [f - 1, f] : [f, f + 1],
          prevLayout: y,
          trigger: "imperative-api",
        }),
        g = validatePanelGroupLayout({ layout: p, panelConstraints: c })
      layoutsEqual(y, g) ||
        update((v) => ({
          mountedGroups: new Map(v.mountedGroups).set(m, { defaultLayoutDeferred: u, derivedPanelConstraints: c, layout: g, separatorToPanels: h }),
        }))
    }
  return {
    collapse: () => {
      const { collapsible: s, collapsedSize: a } = o(),
        { mutableValues: u } = i(),
        c = r()
      s && c !== a && ((u.expandToSize = c), l(a))
    },
    expand: () => {
      const { collapsible: s, collapsedSize: a, minSize: u } = o(),
        { mutableValues: c } = i(),
        m = r()
      if (s && m === a) {
        let y = c.expandToSize ?? u
        y === 0 && (y = 1), l(y)
      }
    },
    getSize: () => {
      const { group: s } = n(),
        a = r(),
        { element: u } = i(),
        c = s.orientation === "horizontal" ? u.offsetWidth : u.offsetHeight
      return { asPercentage: a, inPixels: c }
    },
    isCollapsed: () => {
      const { collapsible: s, collapsedSize: a } = o(),
        u = r()
      return s && fuzzyNumbersEqual(a, u)
    },
    resize: (s) => {
      if (r() !== s) {
        let u
        switch (typeof s) {
          case "number": {
            const { group: c } = n(),
              m = getGroupSize({ group: c })
            u = formatLayoutNumber((s / m) * 100)
            break
          }
          case "string": {
            u = parseFloat(s)
            break
          }
        }
        l(u)
      }
    },
  }
}

function onDoubleClick(e) {
  if (e.defaultPrevented) return
  const { mountedGroups: t } = read()
  findMatchingHitRegions(e, t).forEach((o) => {
    if (o.separator) {
      const i = o.panels.find((r) => r.panelConstraints.defaultSize !== void 0)
      if (i) {
        const r = i.panelConstraints.defaultSize,
          l = getImperativePanelMethods({ groupId: o.group.id, panelId: i.id })
        l && r !== void 0 && (l.resize(r), e.preventDefault())
      }
    }
  })
}

function findSeparatorGroup(e) {
  const { mountedGroups: t } = read()
  for (const [n] of t) if (n.separators.some((o) => o.element === e)) return n
  throw Error("Could not find parent Group for separator element")
}

function getImperativeGroupMethods({ groupId: e }) {
  const t = () => {
    const { mountedGroups: n } = read()
    for (const [o, i] of n) if (o.id === e) return { group: o, ...i }
    throw Error(`Could not find Group with id "${e}"`)
  }
  return {
    getLayout() {
      const { defaultLayoutDeferred: n, layout: o } = t()
      return n ? {} : o
    },
    setLayout(n) {
      const { defaultLayoutDeferred: o, derivedPanelConstraints: i, group: r, layout: l, separatorToPanels: s } = t(),
        a = validatePanelGroupLayout({ layout: n, panelConstraints: i })
      return o
        ? l
        : (layoutsEqual(l, a) ||
            update((u) => ({
              mountedGroups: new Map(u.mountedGroups).set(r, { defaultLayoutDeferred: o, derivedPanelConstraints: i, layout: a, separatorToPanels: s }),
            })),
          a)
    },
  }
}

function getMountedGroup(e) {
  const { mountedGroups: t } = read(),
    n = t.get(e)
  return assert(n, `Mounted Group ${e.id} not found`), n
}

function adjustLayoutForSeparator(e, t) {
  const n = findSeparatorGroup(e),
    o = getMountedGroup(n),
    i = n.separators.find((m) => m.element === e)
  assert(i, "Matching separator not found")
  const r = o.separatorToPanels.get(i)
  assert(r, "Matching panels not found")
  const l = r.map((m) => n.panels.indexOf(m)),
    a = getImperativeGroupMethods({ groupId: n.id }).getLayout(),
    u = adjustLayoutByDelta({ delta: t, initialLayout: a, panelConstraints: o.derivedPanelConstraints, pivotIndices: l, prevLayout: a, trigger: "keyboard" }),
    c = validatePanelGroupLayout({ layout: u, panelConstraints: o.derivedPanelConstraints })
  layoutsEqual(a, c) ||
    update((m) => ({
      mountedGroups: new Map(m.mountedGroups).set(n, {
        defaultLayoutDeferred: o.defaultLayoutDeferred,
        derivedPanelConstraints: o.derivedPanelConstraints,
        layout: c,
        separatorToPanels: o.separatorToPanels,
      }),
    }))
}

function onSeparatorKeyDown(e) {
  if (e.defaultPrevented) return
  const t = e.currentTarget,
    n = findSeparatorGroup(t)
  if (!n.disabled)
    switch (e.key) {
      case "ArrowDown": {
        e.preventDefault(), n.orientation === "vertical" && adjustLayoutForSeparator(t, 5)
        break
      }
      case "ArrowLeft": {
        e.preventDefault(), n.orientation === "horizontal" && adjustLayoutForSeparator(t, -5)
        break
      }
      case "ArrowRight": {
        e.preventDefault(), n.orientation === "horizontal" && adjustLayoutForSeparator(t, 5)
        break
      }
      case "ArrowUp": {
        e.preventDefault(), n.orientation === "vertical" && adjustLayoutForSeparator(t, -5)
        break
      }
      case "End": {
        e.preventDefault(), adjustLayoutForSeparator(t, 100)
        break
      }
      case "Enter": {
        e.preventDefault()
        const o = findSeparatorGroup(t),
          { derivedPanelConstraints: i, layout: r, separatorToPanels: l } = getMountedGroup(o),
          s = o.separators.find((m) => m.element === t)
        assert(s, "Matching separator not found")
        const a = l.get(s)
        assert(a, "Matching panels not found")
        const u = a[0],
          c = i.find((m) => m.panelId === u.id)
        if ((assert(c, "Panel metadata not found"), c.collapsible)) {
          const m = r[u.id],
            y = c.collapsedSize === m ? (o.inMemoryLastExpandedPanelSizes[u.id] ?? c.minSize) : c.collapsedSize
          adjustLayoutForSeparator(t, y - m)
        }
        break
      }
      case "F6": {
        e.preventDefault()
        const i = findSeparatorGroup(t).separators.map((a) => a.element),
          r = Array.from(i).findIndex((a) => a === e.currentTarget)
        assert(r !== null, "Index not found")
        const l = e.shiftKey ? (r > 0 ? r - 1 : i.length - 1) : r + 1 < i.length ? r + 1 : 0
        i[l].focus()
        break
      }
      case "Home": {
        e.preventDefault(), adjustLayoutForSeparator(t, -100)
        break
      }
    }
}

function onPointerDown(e) {
  if (e.defaultPrevented) return
  if (e.pointerType === "mouse" && e.button > 0) return
  const { mountedGroups: t } = read(),
    n = findMatchingHitRegions(e, t),
    o = new Map()
  let i = !1
  n.forEach((r) => {
    r.separator && (i || ((i = !0), r.separator.element.focus()))
    const l = t.get(r.group)
    l && o.set(r.group, l.layout)
  }),
    update({
      interactionState: { hitRegions: n, initialLayoutMap: o, pointerDownAtPoint: { x: e.clientX, y: e.clientY }, state: "active" },
    }),
    n.length && e.preventDefault()
}

const CURSOR_FLAG_HORIZONTAL_MIN = 1,
  CURSOR_FLAG_HORIZONTAL_MAX = 2,
  CURSOR_FLAG_VERTICAL_MIN = 4,
  CURSOR_FLAG_VERTICAL_MAX = 8,
  CURSOR_FLAGS_HORIZONTAL = 3,
  CURSOR_FLAGS_VERTICAL = 12
let supportsAdvancedCursors
function supportsAdvancedCursorStyles() {
  return (
    supportsAdvancedCursors === void 0 &&
      ((supportsAdvancedCursors = !1),
      typeof window < "u" && (window.navigator.userAgent.includes("Chrome") || window.navigator.userAgent.includes("Firefox")) && (supportsAdvancedCursors = !0)),
    supportsAdvancedCursors
  )
}
function getCursorStyle({ cursorFlags: e, groups: t, state: n }) {
  let o = 0,
    i = 0
  switch (n) {
    case "active":
    case "hover":
      t.forEach((r) => {
        if (!r.disableCursor)
          switch (r.orientation) {
            case "horizontal": {
              o++
              break
            }
            case "vertical": {
              i++
              break
            }
          }
      })
  }
  if (o === 0 && i === 0) return null
  switch (n) {
    case "active": {
      if (e && supportsAdvancedCursorStyles()) {
        const r = (e & CURSOR_FLAG_HORIZONTAL_MIN) !== 0,
          l = (e & CURSOR_FLAG_HORIZONTAL_MAX) !== 0,
          s = (e & CURSOR_FLAG_VERTICAL_MIN) !== 0,
          a = (e & CURSOR_FLAG_VERTICAL_MAX) !== 0
        if (r) return s ? "se-resize" : a ? "ne-resize" : "e-resize"
        if (l) return s ? "sw-resize" : a ? "nw-resize" : "w-resize"
        if (s) return "s-resize"
        if (a) return "n-resize"
      }
      break
    }
  }
  return supportsAdvancedCursorStyles()
    ? o > 0 && i > 0
      ? "move"
      : o > 0
        ? "ew-resize"
        : "ns-resize"
    : o > 0 && i > 0
      ? "grab"
      : o > 0
        ? "col-resize"
        : "row-resize"
}
const documentStyles = new WeakMap()
function updateCursorStyle(e) {
  if (e.defaultView === null || e.defaultView === void 0) return
  let { prevStyle: t, styleSheet: n } = documentStyles.get(e) ?? {}
  n === void 0 && ((n = new e.defaultView.CSSStyleSheet()), e.adoptedStyleSheets.push(n))
  const { cursorFlags: o, interactionState: i } = read()
  switch (i.state) {
    case "active":
    case "hover": {
      const r = getCursorStyle({ cursorFlags: o, groups: i.hitRegions.map((s) => s.group), state: i.state }),
        l = `*, *:hover {cursor: ${r} !important; ${i.state === "active" ? "touch-action: none;" : ""} }`
      if (t === l) return
      ;(t = l), r ? (n.cssRules.length === 0 ? n.insertRule(l) : n.replaceSync(l)) : n.cssRules.length === 1 && n.deleteRule(0)
      break
    }
    case "inactive": {
      ;(t = void 0), n.cssRules.length === 1 && n.deleteRule(0)
      break
    }
  }
  documentStyles.set(e, { prevStyle: t, styleSheet: n })
}

function updateLayoutFromPointer({ document: e, event: t, hitRegions: n, initialLayoutMap: o, mountedGroups: i, pointerDownAtPoint: r, prevCursorFlags: l }) {
  let s = 0
  const a = new Map(i)
  n.forEach((c) => {
    const { group: m, groupSize: y } = c,
      { disableCursor: h, orientation: f, panels: d } = m
    let p = 0
    r
      ? f === "horizontal"
        ? (p = ((t.clientX - r.x) / y) * 100)
        : (p = ((t.clientY - r.y) / y) * 100)
      : f === "horizontal"
        ? (p = t.clientX < 0 ? -100 : 100)
        : (p = t.clientY < 0 ? -100 : 100)
    const g = o.get(m),
      { defaultLayoutDeferred: v, derivedPanelConstraints: x, layout: S, separatorToPanels: w } = i.get(m) ?? { defaultLayoutDeferred: !1 }
    if (x && g && S && w) {
      const C = adjustLayoutByDelta({
        delta: p,
        initialLayout: g,
        panelConstraints: x,
        pivotIndices: c.panels.map((P) => d.indexOf(P)),
        prevLayout: S,
        trigger: "mouse-or-touch",
      })
      if (layoutsEqual(C, S)) {
        if (p !== 0 && !h)
          switch (f) {
            case "horizontal": {
              s |= p < 0 ? CURSOR_FLAG_HORIZONTAL_MIN : CURSOR_FLAG_HORIZONTAL_MAX
              break
            }
            case "vertical": {
              s |= p < 0 ? CURSOR_FLAG_VERTICAL_MIN : CURSOR_FLAG_VERTICAL_MAX
              break
            }
          }
      } else {
        a.set(c.group, { defaultLayoutDeferred: v, derivedPanelConstraints: x, layout: C, separatorToPanels: w })
        const P = c.group.panels.map(({ id: T }) => T).join(",")
        c.group.inMemoryLayouts[P] = C
      }
    }
  })
  let u = 0
  t.movementX === 0 ? (u |= l & CURSOR_FLAGS_HORIZONTAL) : (u |= s & CURSOR_FLAGS_HORIZONTAL),
    t.movementY === 0 ? (u |= l & CURSOR_FLAGS_VERTICAL) : (u |= s & CURSOR_FLAGS_VERTICAL),
    update({ cursorFlags: u, mountedGroups: a }),
    updateCursorStyle(e)
}

function onPointerLeave(e) {
  const { cursorFlags: t, interactionState: n, mountedGroups: o } = read()
  switch (n.state) {
    case "active":
      updateLayoutFromPointer({
        document: e.currentTarget,
        event: e,
        hitRegions: n.hitRegions,
        initialLayoutMap: n.initialLayoutMap,
        mountedGroups: o,
        prevCursorFlags: t,
      })
  }
}

function onPointerMove(e) {
  if (e.defaultPrevented) return
  const { cursorFlags: t, interactionState: n, mountedGroups: o } = read()
  switch (n.state) {
    case "active": {
      if (e.buttons === 0) {
        update((i) => (i.interactionState.state === "inactive" ? i : { cursorFlags: 0, interactionState: { state: "inactive" } })),
          update((i) => ({ mountedGroups: new Map(i.mountedGroups) }))
        return
      }
      updateLayoutFromPointer({
        document: e.currentTarget,
        event: e,
        hitRegions: n.hitRegions,
        initialLayoutMap: n.initialLayoutMap,
        mountedGroups: o,
        pointerDownAtPoint: n.pointerDownAtPoint,
        prevCursorFlags: t,
      })
      break
    }
    default: {
      const i = findMatchingHitRegions(e, o)
      i.length === 0 ? n.state !== "inactive" && update({ interactionState: { state: "inactive" } }) : update({ interactionState: { hitRegions: i, state: "hover" } }),
        updateCursorStyle(e.currentTarget)
      break
    }
  }
}

function onPointerOut(e) {
  if (e.relatedTarget instanceof HTMLIFrameElement) {
    const { interactionState: t } = read()
    switch (t.state) {
      case "hover":
        update({ interactionState: { state: "inactive" } })
    }
  }
}

function onPointerUp(e) {
  if (e.defaultPrevented) return
  if (e.pointerType === "mouse" && e.button > 0) return
  const { interactionState: t } = read()
  switch (t.state) {
    case "active":
      update({ cursorFlags: 0, interactionState: { state: "inactive" } }),
        t.hitRegions.length > 0 && (updateCursorStyle(e.currentTarget), update((n) => ({ mountedGroups: new Map(n.mountedGroups) })), e.preventDefault())
  }
}

function getDefaultLayout(e) {
  let t = 0,
    n = 0
  const o = {}
  for (const r of e)
    if (r.defaultSize !== void 0) {
      t++
      const l = formatLayoutNumber(r.defaultSize)
      ;(n += l), (o[r.panelId] = l)
    } else o[r.panelId] = void 0
  const i = e.length - t
  if (i !== 0) {
    const r = formatLayoutNumber((100 - n) / i)
    for (const l of e) l.defaultSize === void 0 && (o[l.panelId] = r)
  }
  return o
}

function notifyPanelOnResize(e, t, n) {
  if (!n[0]) return
  const i = e.panels.find((u) => u.element === t)
  if (!i || !i.onResize) return
  const r = getGroupSize({ group: e }),
    l = e.orientation === "horizontal" ? i.element.offsetWidth : i.element.offsetHeight,
    s = i.mutableValues.prevSize,
    a = { asPercentage: formatLayoutNumber((l / r) * 100), inPixels: l }
  ;(i.mutableValues.prevSize = a), i.onResize(a, i.id, s)
}

function objectsEqual(e, t) {
  if (Object.keys(e).length !== Object.keys(t).length) return !1
  for (const o in e) if (e[o] !== t[o]) return !1
  return !0
}

function layoutNamesPanels(e, t) {
  const n = e.map((i) => i.id),
    o = Object.keys(t)
  if (n.length !== o.length) return !1
  for (const i of n) if (!o.includes(i)) return !1
  return !0
}

const documentGroupCount = new Map()
function mountGroup(e) {
  let t = !0
  assert(e.element.ownerDocument.defaultView, "Cannot register an unmounted Group")
  const n = e.element.ownerDocument.defaultView.ResizeObserver,
    o = new Set(),
    i = new Set(),
    r = new n((f) => {
      lane("default", () => {
        for (const d of f) {
          const { borderBoxSize: p, target: g } = d
          if (g === e.element) {
            if (t) {
              if (getGroupSize({ group: e }) === 0) return
              update((x) => {
                const S = x.mountedGroups.get(e)
                if (S) {
                  const w = calculatePanelConstraints(e),
                    C = S.defaultLayoutDeferred ? getDefaultLayout(w) : S.layout,
                    P = validatePanelGroupLayout({ layout: C, panelConstraints: w })
                  return !S.defaultLayoutDeferred && layoutsEqual(C, P) && objectsEqual(S.derivedPanelConstraints, w)
                    ? x
                    : {
                        mountedGroups: new Map(x.mountedGroups).set(e, {
                          defaultLayoutDeferred: !1,
                          derivedPanelConstraints: w,
                          layout: P,
                          separatorToPanels: S.separatorToPanels,
                        }),
                      }
                }
                return x
              })
            }
          } else notifyPanelOnResize(e, g, p)
        }
      })
    })
  r.observe(e.element),
    e.panels.forEach((f) => {
      assert(!o.has(f.id), `Panel ids must be unique; id "${f.id}" was used more than once`), o.add(f.id), f.onResize && r.observe(f.element)
    })
  const l = getGroupSize({ group: e }),
    s = calculatePanelConstraints(e),
    a = e.panels.map(({ id: f }) => f).join(",")
  let u = e.defaultLayout
  u && (layoutNamesPanels(e.panels, u) || (u = void 0))
  const c = e.inMemoryLayouts[a] ?? u ?? getDefaultLayout(s),
    m = validatePanelGroupLayout({ layout: c, panelConstraints: s }),
    y = calculateHitRegions(e),
    h = e.element.ownerDocument
  return (
    update((f) => {
      const d = new Map()
      return (
        documentGroupCount.set(h, (documentGroupCount.get(h) ?? 0) + 1),
        y.forEach((p) => {
          p.separator && d.set(p.separator, p.panels)
        }),
        { mountedGroups: new Map(f.mountedGroups).set(e, { defaultLayoutDeferred: l === 0, derivedPanelConstraints: s, layout: m, separatorToPanels: d }) }
      )
    }),
    e.separators.forEach((f) => {
      assert(!i.has(f.id), `Separator ids must be unique; id "${f.id}" was used more than once`), i.add(f.id), f.element.addEventListener("keydown", onSeparatorKeyDownLane)
    }),
    documentGroupCount.get(h) === 1 &&
      (h.addEventListener("dblclick", onDoubleClickLane, !0),
      h.addEventListener("pointerdown", onPointerDownLane, !0),
      h.addEventListener("pointerleave", onPointerLeaveLane),
      h.addEventListener("pointermove", onPointerMoveLane),
      h.addEventListener("pointerout", onPointerOutLane),
      h.addEventListener("pointerup", onPointerUpLane, !0)),
    function () {
      ;(t = !1),
        documentGroupCount.set(h, Math.max(0, (documentGroupCount.get(h) ?? 0) - 1)),
        update((d) => {
          const p = new Map(d.mountedGroups)
          return p.delete(e), { mountedGroups: p }
        }),
        e.separators.forEach((d) => {
          d.element.removeEventListener("keydown", onSeparatorKeyDownLane)
        }),
        documentGroupCount.get(h) ||
          (h.removeEventListener("dblclick", onDoubleClickLane, !0),
          h.removeEventListener("pointerdown", onPointerDownLane, !0),
          h.removeEventListener("pointerleave", onPointerLeaveLane),
          h.removeEventListener("pointermove", onPointerMoveLane),
          h.removeEventListener("pointerout", onPointerOutLane),
          h.removeEventListener("pointerup", onPointerUpLane, !0)),
        r.disconnect()
    }
  )
}

// --- React rendering, emulated -------------------------------------------------------------
//
// Components re-render (apply their DOM) when their state changes. Like
// React, updates made during a discrete event are flushed in a microtask,
// updates made during a continuous event or outside events in a later task,
// and updates made while mounting synchronously.

const pending = { sync: new Set(), default: new Set() }
let microtaskScheduled = false
let taskScheduled = false
let currentLane = null
const channel = new MessageChannel()
channel.port1.onmessage = () => {
  taskScheduled = false
  flush("default")
}

function flush(kind) {
  const list = [...pending[kind]]
  pending[kind].clear()
  for (const render of list) render()
}

function lane(kind, fn) {
  const prev = currentLane
  currentLane = kind
  try {
    return fn()
  } finally {
    currentLane = prev
  }
}

const withLane = (kind, fn) => (e) => lane(kind, () => fn(e))
const onDoubleClickLane = withLane("sync", onDoubleClick)
const onPointerDownLane = withLane("sync", onPointerDown)
const onPointerUpLane = withLane("sync", onPointerUp)
const onSeparatorKeyDownLane = withLane("sync", onSeparatorKeyDown)
const onPointerMoveLane = withLane("default", onPointerMove)
const onPointerLeaveLane = withLane("default", onPointerLeave)
const onPointerOutLane = withLane("default", onPointerOut)

// scheduleRender queues a component render in the current lane.
function scheduleRender(render) {
  switch (currentLane) {
    case "now":
      render()
      return
    case "sync":
      pending.sync.add(render)
      if (!microtaskScheduled) {
        microtaskScheduled = true
        queueMicrotask(() => {
          microtaskScheduled = false
          flush("sync")
        })
      }
      return
    default:
      pending.default.add(render)
      if (!taskScheduled) {
        taskScheduled = true
        channel.port2.postMessage(null)
      }
  }
}

function setAttr(el, name, value) {
  if (value === undefined || value === null) el.removeAttribute(name)
  else if (el.getAttribute(name) !== String(value)) el.setAttribute(name, String(value))
}

const configs = new WeakMap()
function readConfig(el) {
  if (!configs.has(el)) {
    let cfg = {}
    const raw = el.getAttribute("data-resizable")
    if (raw) {
      try {
        cfg = JSON.parse(raw)
      } catch {
        cfg = {}
      }
      el.removeAttribute("data-resizable")
    }
    configs.set(el, cfg)
  }
  return configs.get(el)
}

// Ids: the server generates them like React's useId within each outermost
// group; renumber duplicates when a page renders several groups.
let idCounter = 1024
function ensureUniqueId(el) {
  const id = el.id
  if (id && document.querySelectorAll(`[id="${CSS.escape(id)}"]`).length <= 1) return id
  let next
  do next = `_r_${(idCounter++).toString(32)}_`
  while (document.getElementById(next))
  el.id = next
  el.setAttribute("data-testid", next)
  return next
}

// Panel: style flexGrow from the layout, pointer-events while dragging.
function createPanel(el, groupCtx) {
  const cfg = readConfig(el)
  const id = ensureUniqueId(el)
  const panel = {
    element: el,
    id,
    idIsStable: true,
    mutableValues: { expandToSize: void 0, prevSize: void 0 },
    onResize: (size, panelId, prevSize) => emit(el, "shadcn:panel-resize", { size, prevSize }),
    panelConstraints: {
      collapsedSize: cfg.collapsedSize ?? "0%",
      collapsible: cfg.collapsible ?? !1,
      defaultSize: cfg.defaultSize,
      maxSize: cfg.maxSize ?? "100%",
      minSize: cfg.minSize ?? "0%",
    },
    scheduleUpdate: () => scheduleRender(render),
  }
  function render() {
    if (!el.isConnected) return
    const styles = groupCtx.getPanelStyles(groupCtx.id, id)
    const grow = String(styles.flexGrow)
    if (el.style.flexGrow !== grow) el.style.flexGrow = grow
    const pe = styles.pointerEvents ?? ""
    if (el.style.pointerEvents !== pe) el.style.pointerEvents = pe
  }
  el.resizablePanel = getImperativePanelMethods({ groupId: groupCtx.id, panelId: id })
  return { panel, render }
}

// Separator: data-separator and aria-value* from the store.
function createSeparator(el, groupCtx) {
  const id = ensureUniqueId(el)
  const separator = { element: el, id }
  let aria = {}
  let interaction = "inactive"
  function render() {
    if (!el.isConnected) return
    setAttr(el, "aria-controls", aria.valueControls)
    setAttr(el, "aria-valuemax", aria.valueMax)
    setAttr(el, "aria-valuemin", aria.valueMin)
    setAttr(el, "aria-valuenow", aria.valueNow)
    setAttr(el, "data-separator", interaction)
  }
  const offInteraction = eventEmitter.addListener("interactionStateChange", (C) => {
    const next = C.state !== "inactive" && C.hitRegions.some((P) => P.separator === separator) ? C.state : "inactive"
    if (next !== interaction) {
      interaction = next
      scheduleRender(render)
    }
  })
  const offGroups = eventEmitter.addListener("mountedGroupsChange", (C) => {
    C.forEach(({ derivedPanelConstraints: P, layout: T, separatorToPanels: j }, O) => {
      if (O.id === groupCtx.id) {
        const b = j.get(separator)
        if (b) {
          const L = b[0],
            k = O.panels.indexOf(L)
          aria = getSeparatorAria({ layout: T, panelConstraints: P, panelId: L.id, panelIndex: k })
          scheduleRender(render)
        }
      }
    })
  })
  return { separator, render, cleanup: () => (offInteraction(), offGroups()) }
}

function getSeparatorAria({ layout: e, panelConstraints: t, panelId: n, panelIndex: o }) {
  let i, r
  const l = e[n],
    s = t.find((a) => a.panelId === n)
  if (s) {
    const a = s.maxSize,
      u = s.collapsible ? s.collapsedSize : s.minSize,
      c = [o, o + 1]
    ;(r = validatePanelGroupLayout({
      layout: adjustLayoutByDelta({ delta: u - l, initialLayout: e, panelConstraints: t, pivotIndices: c, prevLayout: e }),
      panelConstraints: t,
    })[n]),
      (i = validatePanelGroupLayout({
        layout: adjustLayoutByDelta({ delta: a - l, initialLayout: e, panelConstraints: t, pivotIndices: c, prevLayout: e }),
        panelConstraints: t,
      })[n])
  }
  return { valueControls: n, valueMax: i, valueMin: r, valueNow: l }
}

// Group: registers itself, its panels and its separators with the store.
function createGroup(el) {
  const cfg = readConfig(el)
  const id = ensureUniqueId(el)
  const orientation = el.style.flexDirection === "column" ? "vertical" : "horizontal"
  const defaultLayout = cfg.defaultLayout
  const ctx = {
    id,
    getPanelStyles(b, L) {
      const { interactionState: k, mountedGroups: ne } = read()
      for (const oe of ne.keys())
        if (oe.id === b) {
          const K = ne.get(oe)
          if (K) {
            let X = !1
            switch (k.state) {
              case "active": {
                X = k.hitRegions.some((le) => le.group === oe)
                break
              }
            }
            return { flexGrow: K.layout[L] ?? 1, pointerEvents: X ? "none" : void 0 }
          }
        }
      return { flexGrow: defaultLayout?.[L] ?? 1 }
    },
  }
  const children = [...el.children]
  const panels = children.filter((c) => c.hasAttribute("data-panel")).map((c) => createPanel(c, ctx))
  const separators = children.filter((c) => c.hasAttribute("data-separator")).map((c) => createSeparator(c, ctx))
  const internal = {
    lastExpandedPanelSizes: {},
    layouts: {},
    panels: sortByElementOffset(orientation, panels.map((p) => p.panel)),
    resizeTargetMinimumSize: cfg.resizeTargetMinimumSize ?? { coarse: 20, fine: 10 },
    separators: sortByElementOffset(orientation, separators.map((s) => s.separator)),
  }
  const last = { onLayoutChange: {}, onLayoutChanged: {} }
  const onLayoutChange = (b) => {
    layoutsEqual(last.onLayoutChange, b) || ((last.onLayoutChange = b), emit(el, "shadcn:layout-change", { layout: b }))
  }
  const onLayoutChanged = (b) => {
    layoutsEqual(last.onLayoutChanged, b) || ((last.onLayoutChanged = b), emit(el, "shadcn:layout-changed", { layout: b }))
  }
  const k = {
    defaultLayout,
    disableCursor: !!cfg.disableCursor,
    disabled: !!cfg.disabled,
    element: el,
    id,
    inMemoryLastExpandedPanelSizes: internal.lastExpandedPanelSizes,
    inMemoryLayouts: internal.layouts,
    orientation,
    panels: internal.panels,
    resizeTargetMinimumSize: internal.resizeTargetMinimumSize,
    separators: internal.separators,
  }
  el.resizableGroup = getImperativeGroupMethods({ groupId: id })
  return lane("now", () => {
    const unmount = mountGroup(k)
    const K = read().mountedGroups.get(k)
    if (K) {
      const { defaultLayoutDeferred: q, derivedPanelConstraints: A, layout: V } = K
      !q &&
        A.length > 0 &&
        (onLayoutChange(V),
        onLayoutChanged(V),
        internal.panels.forEach((ue) => {
          ue.scheduleUpdate()
        }))
    }
    let X = !1
    const le = eventEmitter.addListener("interactionStateChange", (q) => {
      const A = q.state === "active"
      X !== A &&
        ((X = A),
        internal.panels.forEach((V) => {
          V.scheduleUpdate()
        }))
    })
    const Qe = eventEmitter.addListener("mountedGroupsChange", (q) => {
      const A = q.get(k)
      if (A) {
        const { defaultLayoutDeferred: V, derivedPanelConstraints: ue, layout: Se } = A
        if (V || ue.length === 0) return
        const { interactionState: et } = read(),
          tt = et.state !== "active"
        onLayoutChange(Se),
          tt && onLayoutChanged(Se),
          internal.panels.forEach((nt) => {
            nt.scheduleUpdate()
          })
      }
    })
    panels.forEach((p) => p.render())
    separators.forEach((s) => s.render())
    return () => {
      unmount(), le(), Qe()
      separators.forEach((s) => s.cleanup())
    }
  })
}

// React runs layout effects children first: register nested groups before
// the groups containing them, once every group on the page is known.
const queue = []
let queued = false
function enqueue(el, holder) {
  queue.push({ el, holder })
  if (queued) return
  queued = true
  queueMicrotask(() => {
    queued = false
    const list = queue.splice(0)
    list.sort((a, b) => {
      if (a.el.contains(b.el)) return 1
      if (b.el.contains(a.el)) return -1
      return a.el.compareDocumentPosition(b.el) & Node.DOCUMENT_POSITION_FOLLOWING ? -1 : 1
    })
    for (const { el, holder } of list) {
      if (!el.isConnected || holder.cancelled) continue
      try {
        holder.cleanup = createGroup(el)
      } catch (err) {
        console.error("[shadcn] resizable:", err)
      }
    }
  })
}

define("resizable", {
  selector: "[data-group]",
  mount(el) {
    const holder = { cleanup: null, cancelled: false }
    enqueue(el, holder)
    return () => {
      holder.cancelled = true
      holder.cleanup?.()
    }
  },
})
