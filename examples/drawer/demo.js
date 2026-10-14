// The drawer-demo example: the goal is React state of the example (it
// survives closing the drawer); the -/+ buttons change it by 10 within
// [200, 400] and are disabled at the bounds.
import { define } from "./core.js"

const goals = new Map() // drawer content id -> goal

define("example-drawer-demo", {
  selector: '[data-slot="drawer-content"]',
  mount(content) {
    const value = content.querySelector(".text-7xl.tracking-tighter")
    const buttons = [...content.querySelectorAll('[data-slot="button"][data-size="icon"]')]
    if (!value || buttons.length !== 2) return
    const [dec, inc] = buttons
    const render = () => {
      const goal = goals.get(content.id) ?? 350
      value.textContent = String(goal)
      dec.disabled = goal <= 200
      inc.disabled = goal >= 400
    }
    const change = (adjustment) => () => {
      const goal = goals.get(content.id) ?? 350
      goals.set(content.id, Math.max(200, Math.min(400, goal + adjustment)))
      render()
    }
    const onDec = change(-10)
    const onInc = change(10)
    dec.addEventListener("click", onDec)
    inc.addEventListener("click", onInc)
    render()
    return () => {
      dec.removeEventListener("click", onDec)
      inc.removeEventListener("click", onInc)
    }
  },
})

// --- the bar chart -------------------------------------------------------------------------------
//
// recharts 3.8 animates the bars of <Bar> when the chart mounts: 400ms,
// "ease", from height 0 at the baseline, driven by requestAnimationFrame
// through its RequestAnimationFrameTimeoutController (which compares rAF
// timestamps with performance.now()). The server renders the final bars
// (the page without JavaScript); here the chart starts empty and grows
// exactly like recharts does, frame by frame.

const DATA = [400, 300, 200, 300, 200, 278, 189, 239, 300, 200, 278, 189, 349]

const roundN = (num) => {
  const rounded = Math.round(num * 1e4) / 1e4
  return Object.is(rounded, -0) ? 0 : rounded
}
const interpolate = (start, end, t) => start + (end - start) * t

// recharts' configBezier("ease").
const ACCURACY = 1e-4
const cubicBezierFactor = (c1, c2) => [0, 3 * c1, 3 * c2 - 6 * c1, 3 * c1 - 3 * c2 + 1]
const evaluatePolynomial = (params, t) => params.map((param, i) => param * t ** i).reduce((pre, curr) => pre + curr)
const cubicBezier = (c1, c2) => (t) => evaluatePolynomial(cubicBezierFactor(c1, c2), t)
const derivativeCubicBezier = (c1, c2) => (t) => {
  const params = cubicBezierFactor(c1, c2)
  return evaluatePolynomial([...params.map((param, i) => param * i).slice(1), 0], t)
}
function createBezierEasing(x1, y1, x2, y2) {
  const curveX = cubicBezier(x1, x2)
  const curveY = cubicBezier(y1, y2)
  const derCurveX = derivativeCubicBezier(x1, x2)
  const rangeValue = (value) => (value > 1 ? 1 : value < 0 ? 0 : value)
  return (_t) => {
    const t = _t > 1 ? 1 : _t
    let x = t
    for (let i = 0; i < 8; ++i) {
      const evalT = curveX(x) - t
      const derVal = derCurveX(x)
      if (Math.abs(evalT - t) < ACCURACY || derVal < ACCURACY) return curveY(x)
      x = rangeValue(x - evalT / derVal)
    }
    return curveY(x)
  }
}
const ease = createBezierEasing(0.25, 0.1, 0.25, 1.0)

// RequestAnimationFrameTimeoutController.setTimeout
function rafTimeout(callback, delay = 0) {
  const startTime = performance.now()
  let requestId = null
  const executeCallback = (now) => {
    if (now - startTime >= delay) callback(now)
    else requestId = requestAnimationFrame(executeCallback)
  }
  requestId = requestAnimationFrame(executeCallback)
  return () => cancelAnimationFrame(requestId)
}

// The bars' geometry for a chart of width x height (BarChart defaults:
// 5px margins, band category axis with a 10% gap, nice [0, max] domain).
function layout(width, height) {
  const margin = 5
  const plotW = width - 2 * margin
  const plotH = height - 2 * margin
  const domainMax = 400 // nice ticks of [0, 400]
  const band = plotW / DATA.length
  const offset = band * 0.1
  const size = (band - 2 * offset) >> 0
  return DATA.map((v, i) => {
    const y = margin + plotH * (1 - v / domainMax)
    return { x: margin + i * band + offset, y, width: size, height: margin + plotH - y, base: margin + plotH }
  })
}

function renderBars(groups, rects, t, barStyle) {
  rects.forEach((r, i) => {
    const g = groups[i]
    if (!g) return
    const y = t === 1 ? r.y : interpolate(r.base, r.y, t)
    const height = t === 1 ? r.height : interpolate(0, r.height, t)
    let path = g.querySelector("path")
    if (height === 0) {
      path?.remove()
      return
    }
    if (!path) {
      path = document.createElementNS("http://www.w3.org/2000/svg", "path")
      path.setAttribute("name", "undefined")
      path.setAttribute("class", "recharts-rectangle")
      path.setAttribute("style", barStyle)
      g.appendChild(path)
    }
    path.setAttribute("x", String(roundN(r.x)))
    path.setAttribute("y", String(roundN(y)))
    path.setAttribute("width", String(roundN(r.width)))
    path.setAttribute("height", String(roundN(height)))
    path.setAttribute("radius", "0")
    path.setAttribute("d", `M ${roundN(r.x)},${roundN(y)} h ${roundN(r.width)} v ${roundN(height)} h ${roundN(-r.width)} Z`)
  })
}

define("example-drawer-demo-chart", {
  selector: '[data-slot="drawer-content"] .recharts-responsive-container',
  mount(container) {
    const wrapper = container.querySelector(".recharts-wrapper")
    const groups = [...container.querySelectorAll(".recharts-inactive-bar")]
    if (!wrapper || groups.length !== DATA.length) return
    const rects = layout(Number(wrapper.getAttribute("width")), Number(wrapper.getAttribute("height")))
    // The server-rendered bars carry the example's style prop.
    const barStyle = container.querySelector(".recharts-rectangle")?.getAttribute("style") ?? "fill: hsl(var(--foreground)); opacity: 0.9;"
    renderBars(groups, rects, 0, barStyle)
    const stops = []
    let started = false
    const start = () => {
      // animationManager.start([onAnimationStart, begin, onAnimationActive, duration, onAnimationEnd])
      stops.push(
        rafTimeout(() => {
          stops.push(
            rafTimeout(() => {
              let beginTime
              const timingUpdate = (now) => {
                if (!beginTime) beginTime = now
                const t = (now - beginTime) / 400
                renderBars(groups, rects, ease(t), barStyle)
                if (t < 1) stops.push(rafTimeout(timingUpdate))
                else renderBars(groups, rects, 1, barStyle)
              }
              stops.push(rafTimeout(timingUpdate))
            }, 0),
          )
        }),
      )
    }
    // ResponsiveContainer renders the chart once its ResizeObserver
    // measured the container; the animation starts then.
    const ro = new ResizeObserver(() => {
      if (started) return
      started = true
      start()
    })
    ro.observe(container)
    return () => {
      ro.disconnect()
      stops.forEach((s) => s())
    }
  },
})
