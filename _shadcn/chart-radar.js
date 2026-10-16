// radar.js: recharts 3.8.0's RadarChart for the shadcn/ui Chart runtime:
// RadarChart (PolarChart with the centric layout), Radar, PolarGrid,
// PolarAngleAxis and PolarRadiusAxis, with the axis tooltip (pointer angle,
// keyboard), active dots, dots, legend and the radar animation.
//
// The chart reuses CartesianChart's machinery (DOM skeleton, z-index
// layers, tooltip and legend rendering, events, animations) and replaces
// the state (polar view box, angle and radius axes) and the items.
//
// The polar pieces are exported for other polar charts:
//
//   import { PolarState, renderPolarGrid, renderPolarAngleAxis, renderPolarRadiusAxis } from "./chart-radar.js"
import { CartesianChart, registerChart, registerFn, resolveSpecNode, withDefaults, svgProps, DefaultZIndexes, Layer, Text, Dot, h, rechartsScale } from "./chart.js"
import * as R from "./recharts.js"

const wellBehaved = R.isWellBehavedNumber

function clsx(...args) {
  const out = []
  for (const a of args) {
    if (!a) continue
    if (typeof a === "string") out.push(a)
    else if (typeof a === "object") for (const k in a) if (a[k]) out.push(k)
  }
  return out.join(" ")
}

function omit(obj, keys) {
  const out = { ...obj }
  for (const k of keys) delete out[k]
  return out
}

function svgPropsFromUnknown(input) {
  if (input == null) return null
  if (typeof input === "object" && !Array.isArray(input) && typeof input !== "function") return svgProps(input)
  return null
}

const getLegendItemColor = (stroke, fill) => (stroke && stroke !== "none" ? stroke : fill)

// --- defaults -----------------------------------------------------------------------

const POLAR_CHART_DEFAULTS = { accessibilityLayer: true, stackOffset: "none", barCategoryGap: "10%", barGap: 4, reverseStackOrder: false, syncMethod: "index", layout: "radial", responsive: false, cx: "50%", cy: "50%", innerRadius: 0, outerRadius: "80%" }
const RADAR_CHART_DEFAULTS = { ...POLAR_CHART_DEFAULTS, layout: "centric", startAngle: 90, endAngle: -270 }
const DEFAULT_MARGIN = { top: 5, right: 5, bottom: 5, left: 5 }

export const POLAR_ANGLE_AXIS_DEFAULTS = { allowDecimals: false, allowDuplicatedCategory: true, allowDataOverflow: false, angle: 0, angleAxisId: 0, axisLine: true, axisLineType: "polygon", cx: 0, cy: 0, hide: false, includeHidden: false, label: false, niceTicks: "auto", orientation: "outer", reversed: false, scale: "auto", tick: true, tickLine: true, tickSize: 8, type: "auto", zIndex: DefaultZIndexes.axis }
export const POLAR_RADIUS_AXIS_DEFAULTS = { allowDataOverflow: false, allowDecimals: false, allowDuplicatedCategory: true, angle: 0, axisLine: true, includeHidden: false, hide: false, niceTicks: "auto", label: false, orientation: "right", radiusAxisId: 0, reversed: false, scale: "auto", stroke: "#ccc", tick: true, tickCount: 5, tickLine: true, type: "auto", zIndex: DefaultZIndexes.axis }
export const POLAR_GRID_DEFAULTS = { angleAxisId: 0, radiusAxisId: 0, gridType: "polygon", radialLines: true, zIndex: DefaultZIndexes.grid }
const RADAR_DEFAULTS = { activeDot: true, angleAxisId: 0, animationBegin: 0, animationDuration: 1500, animationEasing: "ease", dot: false, hide: false, isAnimationActive: "auto", label: false, legendType: "rect", radiusAxisId: 0, zIndex: DefaultZIndexes.area }
const TOOLTIP_DEFAULTS = { allowEscapeViewBox: { x: false, y: false }, animationDuration: 400, animationEasing: "ease", axisId: 0, contentStyle: {}, cursor: true, filterNull: true, includeHidden: false, isAnimationActive: "auto", itemSorter: "name", itemStyle: {}, labelStyle: {}, offset: 10, reverseDirection: { x: false, y: false }, separator: " : ", trigger: "hover", useTranslate3d: false, wrapperStyle: {} }
const LEGEND_DEFAULTS = { align: "center", iconSize: 14, inactiveColor: "#ccc", itemSorter: "value", layout: "horizontal", verticalAlign: "bottom" }

const IMPLICIT_ANGLE = { allowDataOverflow: false, allowDecimals: false, allowDuplicatedCategory: false, dataKey: undefined, domain: undefined, id: 0, includeHidden: false, name: undefined, reversed: false, scale: "auto", tick: true, tickCount: undefined, ticks: undefined, type: "auto", unit: undefined, niceTicks: "auto" }
const IMPLICIT_RADIUS = { allowDataOverflow: false, allowDecimals: false, allowDuplicatedCategory: true, dataKey: undefined, domain: undefined, id: 0, includeHidden: false, name: undefined, reversed: false, scale: "auto", tick: true, tickCount: 5, ticks: undefined, type: "auto", unit: undefined, niceTicks: "auto" }

// --- polar state (polarAxisSelectors, polarScaleSelectors, polarSelectors) -------------

function axisTypeBasedOnLayout(layout, axisType, type) {
  if (type !== "auto") return type
  return R.isCategoricalAxis(layout, axisType) ? "category" : "number"
}

function makeDomain(val) {
  const num = (v) => {
    if (R.isNumOrStr(v) || v instanceof Date) {
      const n = Number(v)
      if (wellBehaved(n)) return n
    }
    return undefined
  }
  if (Array.isArray(val)) {
    const a = [num(val[0]), num(val[1])]
    return R.isWellFormedNumberDomain(a) ? a : undefined
  }
  const n = num(val)
  return n == null ? undefined : [n, n]
}

function getDomainDefinition(s) {
  if (s.domain != null) return s.domain
  if (s.ticks != null) {
    if (s.type === "number") {
      const all = s.ticks.map(Number).filter(wellBehaved)
      return [Math.min(...all), Math.max(...all)]
    }
    if (s.type === "category") return s.ticks.map(String)
  }
  return [0, "auto"]
}

/**
 * Everything recharts derives from the store for a polar chart: offset,
 * polar view box, angle and radius axes (domain, scale, ticks).
 */
export class PolarState {
  constructor(chart) {
    this.chart = chart
    this.registered = chart.registered !== false
    this.layout = chart.props.layout
    this.width = chart.width
    this.height = chart.height
    this.data = chart.data
    this.items = this.registered ? chart.items : []
    this.axes = new Map()
    const margin = chart.margin
    let offset = { top: margin.top || 0, bottom: margin.bottom || 0, left: margin.left || 0, right: margin.right || 0 }
    const brushBottom = offset.bottom
    if (this.registered) offset = R.appendOffsetOfLegend(offset, chart.legendSettings, chart.legendSize)
    this.offset = { brushBottom, ...offset, width: Math.max(this.width - offset.left - offset.right, 0), height: Math.max(this.height - offset.top - offset.bottom, 0) }
    const p = chart.props
    const maxRadius = R.getMaxRadius(this.width, this.height, this.offset)
    this.maxRadius = maxRadius
    const innerRadius = R.getPercentValue(p.innerRadius, maxRadius, 0)
    const outerRadius = R.getPercentValue(p.outerRadius, maxRadius, maxRadius * 0.8)
    this.viewBox = { cx: R.getPercentValue(p.cx, this.width, this.width / 2), cy: R.getPercentValue(p.cy, this.height, this.height / 2), innerRadius, outerRadius, startAngle: p.startAngle, endAngle: p.endAngle, clockWise: false }
  }

  axisSettings(axisType, id) {
    const list = this.registered ? (axisType === "angleAxis" ? this.chart.angleAxes : this.chart.radiusAxes) : []
    const found = list.find((a) => a.settings.id === id)
    if (found) return found.settings
    const implicit = axisType === "angleAxis" ? IMPLICIT_ANGLE : IMPLICIT_RADIUS
    return { ...implicit, type: axisTypeBasedOnLayout(this.layout, axisType, implicit.type) }
  }

  tooltipAxisType() {
    return this.layout === "centric" ? "angleAxis" : "radiusAxis"
  }

  /** The axis: settings, scale, ticks, tooltip ticks. */
  axis(axisType, id) {
    const key = `${axisType}:${id}`
    if (this.axes.has(key)) return this.axes.get(key)
    const a = this.computeAxis(axisType, id)
    this.axes.set(key, a)
    return a
  }

  computeAxis(axisType, id) {
    const settings = this.axisSettings(axisType, id)
    const layout = this.layout
    const isCategorical = R.isCategoricalAxis(layout, axisType)
    const items = this.items.filter((i) => i[axisType === "angleAxis" ? "angleAxisId" : "radiusAxisId"] === id).filter((i) => settings.includeHidden || !i.hide)
    const data = this.data
    // applied values
    let applied
    if (settings.dataKey != null) applied = data.map((e) => ({ value: R.getValueByDataKey(e, settings.dataKey) }))
    else if (items.length > 0) applied = items.flatMap((item) => data.map((e) => ({ value: R.getValueByDataKey(e, item.dataKey) })))
    else applied = data.map((e) => ({ value: e }))
    // domain
    let domain
    if (data.length === 0) domain = undefined
    else if (isCategorical && settings.dataKey == null) domain = Array.from({ length: data.length }, (_, i) => i)
    else if (settings.type === "category") {
      const cat = applied.map(({ value }) => (R.isNumOrStr(value) || value instanceof Date ? value : undefined)).filter((v) => v != null)
      domain = settings.allowDuplicatedCategory ? cat : Array.from(new Set(cat))
    } else {
      let lo, hi
      const take = (v) => {
        const d = makeDomain(v)
        if (d) {
          lo = lo == null ? d[0] : Math.min(lo, d[0])
          hi = hi == null ? d[1] : Math.max(hi, d[1])
        }
      }
      if (items.length > 0) data.forEach((e) => items.forEach((item) => take(R.getValueByDataKey(e, settings.dataKey ?? item.dataKey))))
      if (settings.dataKey != null) data.forEach((e) => take(R.getValueByDataKey(e, settings.dataKey)))
      const dataDomain = wellBehaved(lo) && wellBehaved(hi) ? [lo, hi] : undefined
      domain = R.parseNumericalUserDomain(getDomainDefinition(settings), dataDomain, settings.allowDataOverflow)
    }
    // scale type
    let realScaleType
    if (settings.scale === "auto") realScaleType = settings.type === "category" ? "band" : "linear"
    else realScaleType = ["band", "linear", "point"].includes(settings.scale) ? settings.scale : "point"
    // nice ticks
    let niceTicks
    if (settings.niceTicks !== "none" && realScaleType === "linear" && settings.tickCount && R.isWellFormedNumberDomain(domain)) {
      const def = getDomainDefinition(settings)
      const hasAuto = Array.isArray(def) && (def[0] === "auto" || def[1] === "auto")
      const mode = settings.niceTicks === "auto" ? "adaptive" : settings.niceTicks
      if (hasAuto) niceTicks = R.getNiceTickValues(domain, settings.tickCount, settings.allowDecimals, mode)
      else if (settings.type === "number") niceTicks = R.getTickValuesFixedDomain(domain, settings.tickCount, settings.allowDecimals, mode)
    }
    if (axisType !== "angleAxis" && settings.type === "number" && R.isWellFormedNumberDomain(domain) && niceTicks?.length) domain = [Math.min(domain[0], niceTicks[0] ?? 0), Math.max(domain[1], niceTicks[niceTicks.length - 1] ?? 0)]
    if (realScaleType === "linear" && domain && !R.isWellFormedNumberDomain(domain)) {
      const nums = domain.filter(wellBehaved)
      domain = nums.length ? [Math.min(...nums), Math.max(...nums)] : undefined
    }
    // range
    const vb = this.viewBox
    let range = axisType === "angleAxis" ? [vb.startAngle, vb.endAngle] : [vb.innerRadius, vb.outerRadius]
    if (settings.reversed) range = [range[1], range[0]]
    const factory = { band: R.scaleBand, linear: R.scaleLinear, point: R.scalePoint }[realScaleType]
    const scale = domain && factory ? rechartsScale(factory().domain(domain).range(range)) : undefined
    let ticks
    let tooltipTicks
    if (scale) {
      let offset = settings.type === "category" && scale.bandwidth ? scale.bandwidth() / 2 : 0
      if (axisType === "angleAxis") offset = Math.sign(range[0] - range[1]) * 2 * offset
      const entries = (list, valueOf) =>
        list
          .map((entry, index) => {
            const scaled = scale.map(valueOf ? valueOf(entry) : entry)
            return wellBehaved(scaled) ? { coordinate: scaled + offset, value: entry, index, offset } : null
          })
          .filter(Boolean)
      const t = settings.ticks || niceTicks
      if (t) ticks = entries(t)
      else if (scale.ticks) ticks = entries(scale.ticks(settings.tickCount))
      else ticks = entries(scale.domain())
      tooltipTicks = entries(scale.domain())
    }
    if (axisType === "angleAxis" && ticks) {
      const unique = new Map()
      for (const t of ticks) {
        const c = (t.coordinate + 360) % 360
        if (!unique.has(c)) unique.set(c, t)
      }
      ticks = [...unique.values()]
    }
    return { settings, scale, ticks, tooltipTicks, range, realScaleType, axisWithScale: scale ? { ...settings, scale } : undefined }
  }

  plotArea() {
    const o = this.offset
    return { x: o.left, y: o.top, width: o.width, height: o.height }
  }
}

// --- polar grid and axes ------------------------------------------------------------

function polygonPathOfGrid(radius, cx, cy, angles) {
  let path = ""
  angles.forEach((angle, i) => {
    const p = R.polarToCartesian(cx, cy, radius, angle)
    path += i ? `L ${p.x},${p.y}` : `M ${p.x},${p.y}`
  })
  return path + "Z"
}

/** PolarGrid: returns the <g class="recharts-polar-grid"> or null. */
export function renderPolarGrid(state, gridProps) {
  const { gridType, radialLines, angleAxisId, radiusAxisId, polarAngles: anglesIn, polarRadius: radiusIn, zIndex, ...inputs } = gridProps
  const vb = state.viewBox
  const polarAngles = Array.isArray(anglesIn) ? anglesIn : state.axis("angleAxis", angleAxisId).ticks?.map((t) => t.coordinate)
  const polarRadius = Array.isArray(radiusIn) ? radiusIn : state.axis("radiusAxis", radiusAxisId).ticks?.map((t) => t.coordinate)
  if (polarAngles == null || polarRadius == null) return null
  const props = { cx: vb?.cx ?? inputs.cx ?? 0, cy: vb?.cy ?? inputs.cy ?? 0, innerRadius: vb?.innerRadius ?? 0, outerRadius: vb?.outerRadius ?? 0, ...omit(inputs, ["cx", "cy", "innerRadius", "outerRadius"]) }
  if (props.outerRadius <= 0) return null
  const { cx, cy } = props
  const circle = (p, radius) => h("circle", { stroke: "#ccc", fill: "none", ...svgProps(p), className: clsx("recharts-polar-grid-concentric-circle", p.className), cx, cy, r: radius })
  const polygon = (p, radius) => h("path", { stroke: "#ccc", fill: "none", ...svgProps({ ...p, radius }), className: clsx("recharts-polar-grid-concentric-polygon", p.className), d: polygonPathOfGrid(radius, cx, cy, polarAngles) })
  const children = []
  if (polarRadius.length) {
    const max = Math.max(...polarRadius)
    const background = props.fill && props.fill !== "none"
    const rings = []
    if (background) rings.push(gridType === "circle" ? circle({ ...props, radius: max }, max) : polygon(props, max))
    for (const r of polarRadius) rings.push(gridType === "circle" ? circle({ ...props, fill: "none", radius: r }, r) : polygon({ ...props, fill: "none" }, r))
    children.push(h("g", { className: "recharts-polar-grid-concentric" }, rings))
  }
  if (polarAngles.length && radialLines) {
    const lineProps = { stroke: "#ccc", ...svgProps(props) }
    children.push(
      h(
        "g",
        { className: "recharts-polar-grid-angle" },
        polarAngles.map((a) => {
          const s = R.polarToCartesian(cx, cy, props.innerRadius, a)
          const e = R.polarToCartesian(cx, cy, props.outerRadius, a)
          return h("line", { ...lineProps, x1: s.x, y1: s.y, x2: e.x, y2: e.y })
        }),
      ),
    )
  }
  return h("g", { className: "recharts-polar-grid" }, children)
}

// recharts Polygon's path.
const validPoint = (p) => p != null && p.x === +p.x && p.y === +p.y
function parsedPoints(points = []) {
  let segs = [[]]
  points.forEach((entry) => {
    const last = segs[segs.length - 1]
    if (validPoint(entry)) last?.push(entry)
    else if (last && last.length > 0) segs.push([])
  })
  const first = points[0]
  const last = segs[segs.length - 1]
  if (validPoint(first) && last) last.push(first)
  if (segs[segs.length - 1]?.length <= 0) segs = segs.slice(0, -1)
  return segs
}
function singlePolygonPath(points, connectNulls) {
  let segs = parsedPoints(points)
  if (connectNulls) segs = [segs.flat()]
  const path = segs.map((seg) => seg.reduce((acc, p, i) => R.roundTemplateLiteral(["", "", "", ",", ""], acc, i === 0 ? "M" : "L", p.x, p.y), "")).join("")
  return segs.length === 1 ? `${path}Z` : path
}

/** recharts <Polygon> (single path; range polygons are not supported). */
export function Polygon(props) {
  const { points, className, baseLinePoints, connectNulls, ...others } = props
  if (!points || !points.length) return null
  const d = singlePolygonPath(points, connectNulls)
  return h("path", { ...svgProps(others), fill: d.slice(-1) === "Z" ? others.fill : "none", className: clsx("recharts-polygon", className), d })
}

const COS_45 = Math.cos(R.degreeToRadian(45))
const EPS = 1e-5

/** PolarAngleAxis: returns the axis layer or null. */
export function renderPolarAngleAxis(state, axisDef) {
  const inputs = axisDef.props
  const a = state.axis("angleAxis", inputs.angleAxisId)
  const vb = state.viewBox
  const ticks = a.ticks
  if (!vb || !ticks || !ticks.length || !a.scale || inputs.hide) return null
  const props = { ...inputs, ...vb, radius: vb.outerRadius, ticks }
  const { cx, cy, radius } = props
  let axisLine = null
  if (props.axisLine) {
    const lineProps = { ...svgProps(props), fill: "none", ...(svgPropsFromUnknown(props.axisLine) ?? {}) }
    if (props.axisLineType === "circle") axisLine = Dot({ className: "recharts-polar-angle-axis-line", ...lineProps, cx, cy, r: radius })
    else axisLine = Polygon({ className: "recharts-polar-angle-axis-line", ...lineProps, points: ticks.map((t) => R.polarToCartesian(cx, cy, radius, t.coordinate)) })
  }
  const axisProps = omit(svgProps(props), ["ref"])
  const customTickProps = svgPropsFromUnknown(props.tick) ?? {}
  const tickLineProps = { ...axisProps, fill: "none", ...(svgPropsFromUnknown(props.tickLine) ?? {}) }
  const items = ticks.map((entry, i) => {
    const size = props.tickSize || 8
    const p1 = R.polarToCartesian(cx, cy, radius, entry.coordinate)
    const p2 = R.polarToCartesian(cx, cy, radius + (props.orientation === "inner" ? -1 : 1) * size, entry.coordinate)
    const cos = Math.cos(R.degreeToRadian(-entry.coordinate))
    const sin = Math.sin(R.degreeToRadian(-entry.coordinate))
    const textAnchor = cos > EPS ? (props.orientation === "outer" ? "start" : "end") : cos < -EPS ? (props.orientation === "outer" ? "end" : "start") : "middle"
    const verticalAnchor = Math.abs(cos) <= COS_45 ? (sin > 0 ? "start" : "end") : "middle"
    const tickProps = { ...axisProps, textAnchor, verticalAnchor, stroke: "none", fill: props.stroke, ...customTickProps, index: i, payload: entry, x: p2.x, y: p2.y }
    const value = typeof props.tickFormatter === "function" ? props.tickFormatter(entry.value, i) : entry.value
    let text = null
    if (props.tick) text = typeof props.tick === "function" ? props.tick(tickProps) : Text({ ...tickProps, className: "recharts-polar-angle-axis-tick-value", children: value == null ? value : String(value) })
    return Layer({ className: clsx("recharts-polar-angle-axis-tick", typeof props.tick === "object" ? props.tick?.className : undefined) }, props.tickLine ? h("line", { className: "recharts-polar-angle-axis-tick-line", ...tickLineProps, x1: p1.x, y1: p1.y, x2: p2.x, y2: p2.y }) : null, text)
  })
  return Layer({ className: clsx("recharts-polar-angle-axis", "angleAxis", props.className) }, axisLine, Layer({ className: "recharts-polar-angle-axis-ticks" }, items))
}

/** PolarRadiusAxis: returns the axis layer or null. */
export function renderPolarRadiusAxis(state, axisDef) {
  const inputs = axisDef.props
  const a = state.axis("radiusAxis", inputs.radiusAxisId)
  const vb = state.viewBox
  const ticks = a.ticks
  if (!vb || !ticks || !ticks.length || !a.scale || inputs.hide) return null
  const props = { ...inputs, ...vb }
  const { cx, cy, angle } = props
  let line = null
  if (props.axisLine) {
    const extent = ticks.reduce((r, e) => [Math.min(r[0], e.coordinate), Math.max(r[1], e.coordinate)], [Infinity, -Infinity])
    const p0 = R.polarToCartesian(cx, cy, extent[0], angle)
    const p1 = R.polarToCartesian(cx, cy, extent[1], angle)
    line = h("line", { className: "recharts-polar-radius-axis-line", ...svgProps(omit(props, ["cx", "cy", "angle", "axisLine"])), fill: "none", ...(svgPropsFromUnknown(props.axisLine) ?? {}), x1: p0.x, y1: p0.y, x2: p1.x, y2: p1.y })
  }
  let tickLayer = null
  if (props.tick) {
    const textAnchor = props.orientation === "left" ? "end" : props.orientation === "right" ? "start" : "middle"
    const axisProps = svgProps(omit(props, ["angle", "tickFormatter", "stroke", "tick"]))
    const customTickProps = svgPropsFromUnknown(props.tick) ?? {}
    const items = ticks.map((entry, i) => {
      const coord = R.polarToCartesian(cx, cy, entry.coordinate, angle)
      const tickProps = { textAnchor, transform: `rotate(${90 - angle}, ${coord.x}, ${coord.y})`, ...axisProps, stroke: "none", fill: props.stroke, ...customTickProps, index: i, ...coord, payload: entry }
      const value = typeof props.tickFormatter === "function" ? props.tickFormatter(entry.value, i) : entry.value
      const node = typeof props.tick === "function" ? props.tick(tickProps) : Text({ ...tickProps, className: "recharts-polar-radius-axis-tick-value", children: value == null ? value : String(value) })
      return Layer({ className: clsx("recharts-polar-radius-axis-tick", typeof props.tick === "object" ? props.tick?.className : undefined) }, node)
    })
    tickLayer = Layer({ className: "recharts-polar-radius-axis-ticks" }, items)
  }
  return Layer({ className: clsx("recharts-polar-radius-axis", "radiusAxis", props.className) }, line, tickLayer)
}

// --- RadarChart ------------------------------------------------------------------------

function interpolatePoint(prev, factor, t) {
  return (entry, index) => {
    const p = prev && prev[Math.floor(index * factor)]
    if (p) return { ...entry, x: R.interpolate(p.x, entry.x, t), y: R.interpolate(p.y, entry.y, t) }
    return { ...entry, x: R.interpolate(entry.cx, entry.x, t), y: R.interpolate(entry.cy, entry.y, t) }
  }
}

export class RadarChartInstance extends CartesianChart {
  setSpec(spec, initial = false) {
    const root = resolveSpecNode(spec, this.ctx)
    this.spec = spec
    this.name = root.type
    this.props = withDefaults(root.props, RADAR_CHART_DEFAULTS)
    this.margin = root.props.margin ?? DEFAULT_MARGIN
    this.data = root.props.data ?? []
    this.children = root.children
    const prevItems = this.items ?? []
    this.items = []
    this.xAxes = []
    this.yAxes = []
    this.angleAxes = []
    this.radiusAxes = []
    this.polarGrids = []
    this.tooltipProps = undefined
    this.tooltipContent = undefined
    this.legendProps = undefined
    this.defs = []
    let n = 1
    for (const child of root.children) {
      const p = child.props
      switch (child.type) {
        case "PolarAngleAxis": {
          const props = withDefaults(p, POLAR_ANGLE_AXIS_DEFAULTS)
          const settings = { id: props.angleAxisId, scale: props.scale, type: axisTypeBasedOnLayout(this.props.layout, "angleAxis", props.type), dataKey: props.dataKey, name: props.name, allowDuplicatedCategory: false, allowDataOverflow: false, reversed: props.reversed, includeHidden: false, allowDecimals: props.allowDecimals, tickCount: props.tickCount, niceTicks: props.niceTicks ?? "auto", ticks: props.ticks, tick: props.tick, domain: props.domain }
          this.angleAxes.push({ props, settings, node: child })
          break
        }
        case "PolarRadiusAxis": {
          const props = withDefaults(p, POLAR_RADIUS_AXIS_DEFAULTS)
          const settings = { id: props.radiusAxisId, domain: props.domain, scale: props.scale, type: axisTypeBasedOnLayout(this.props.layout, "radiusAxis", props.type), dataKey: props.dataKey, name: props.name, allowDuplicatedCategory: props.allowDuplicatedCategory, allowDataOverflow: props.allowDataOverflow, reversed: props.reversed, includeHidden: props.includeHidden, allowDecimals: props.allowDecimals, niceTicks: props.niceTicks ?? "auto", ticks: props.ticks, tickCount: props.tickCount, tick: props.tick }
          this.radiusAxes.push({ props, settings, node: child })
          break
        }
        case "PolarGrid":
          this.polarGrids.push(withDefaults(p, POLAR_GRID_DEFAULTS))
          break
        case "Tooltip":
          this.tooltipProps = withDefaults(p, TOOLTIP_DEFAULTS)
          this.tooltipContent = child.children.find((c) => c.type === "ChartTooltipContent")
          break
        case "Legend":
          this.legendProps = withDefaults(p, LEGEND_DEFAULTS)
          this.legendContent = child.children.find((c) => c.type === "ChartLegendContent")
          break
        case "Radar": {
          const props = withDefaults(p, RADAR_DEFAULTS)
          const id = props.id ?? `recharts-radar-${this.ids.item(n++)}`
          this.items.push({ type: "radar", id, props: { ...props, id }, node: child, dataKey: props.dataKey, angleAxisId: props.angleAxisId, radiusAxisId: props.radiusAxisId, hide: props.hide, prev: prevItems.find((x) => x.id === id) })
          break
        }
        case "defs":
          this.defs.push(child)
          break
      }
    }
    if (!initial) this.commit()
  }

  ensureSvg() {
    const fresh = !this.svg
    super.ensureSvg()
    if (fresh) {
      // PolarChart passes cx and cy to the surface.
      this.svg.setAttribute("cx", String(this.props.cx))
      this.svg.setAttribute("cy", String(this.props.cy))
    }
  }

  render() {
    if (!this.wrapper) return
    const w = this.width
    const hgt = this.height
    if (!(w > 0) || !(hgt > 0)) {
      this.wrapper.remove()
      return
    }
    if (!this.wrapper.isConnected) this.inner.appendChild(this.wrapper)
    this.wrapper.setAttribute("width", String(w))
    this.wrapper.setAttribute("height", String(hgt))
    Object.assign(this.wrapper.style, { position: "relative", cursor: "default", width: `${w}px`, height: `${hgt}px` })
    this.ensureTooltipWrapper()
    this.ensureLegend()
    this.ensureSvg()
    for (const [k, v] of Object.entries({ width: w, height: hgt, viewBox: `0 0 ${w} ${hgt}` })) if (this.svg.getAttribute(k) !== String(v)) this.svg.setAttribute(k, String(v))
    const state = new PolarState(this)
    this.state = state
    const rendered = new Set()
    const pa = state.plotArea()
    const zero = [h("defs", null, h("clipPath", { id: this.clipPathId }, h("rect", { x: pa.x, y: pa.y, height: pa.height, width: pa.width })))]
    let gi = 0
    let ai = 0
    let ri = 0
    for (const child of this.children) {
      switch (child.type) {
        case "PolarGrid": {
          const props = this.polarGrids[gi]
          this.slot(rendered, props.zIndex, `grid-${gi++}`, state.registered ? renderPolarGrid(state, props) : null)
          break
        }
        case "PolarAngleAxis": {
          const def = this.angleAxes[ai]
          this.slot(rendered, def.props.zIndex, `angle-${ai++}`, state.registered ? renderPolarAngleAxis(state, def) : null)
          break
        }
        case "PolarRadiusAxis": {
          const def = this.radiusAxes[ri]
          this.slot(rendered, def.props.zIndex, `radius-${ri++}`, state.registered ? renderPolarRadiusAxis(state, def) : null)
          break
        }
        case "Radar": {
          const item = this.items.find((i) => i.node === child)
          if (item && state.registered) this.renderRadar(state, item, rendered)
          break
        }
      }
    }
    this.setZeroLayer(zero)
    for (const key of this.slots.keys()) if (!rendered.has(key)) this.slots.clear(key)
    this.renderTooltip(state)
    this.renderLegend(state)
  }

  computeRadarPoints(state, item) {
    const radius = state.axis("radiusAxis", item.radiusAxisId)
    const angle = state.axis("angleAxis", item.angleAxisId)
    const vb = state.viewBox
    if (!radius.scale || !angle.scale || !vb) return undefined
    const bandSize = R.isCategoricalAxis(state.layout, "radiusAxis") ? R.getBandSizeOfAxis(radius.axisWithScale, radius.ticks, false) : R.getBandSizeOfAxis(angle.axisWithScale, angle.ticks, false)
    const angleBandSize = angle.settings.type !== "number" ? bandSize ?? 0 : 0
    const { cx, cy } = vb
    return state.data.map((entry, i) => {
      const name = R.getValueByDataKey(entry, angle.settings.dataKey, i)
      const value = R.getValueByDataKey(entry, item.dataKey)
      const a = (angle.scale.map(name) ?? 0) + angleBandSize
      const pointValue = Array.isArray(value) ? value[value.length - 1] : value
      const r = pointValue == null ? 0 : radius.scale.map(pointValue) ?? 0
      return { ...R.polarToCartesian(cx, cy, r, a), name, value, cx, cy, radius: r, angle: a, payload: entry }
    })
  }

  renderRadar(state, item, rendered) {
    const props = item.props
    const key = item.id
    const points = props.hide ? undefined : this.computeRadarPoints(state, item)
    if (!points) {
      this.slot(rendered, props.zIndex, key, [])
      this.slot(rendered, DefaultZIndexes.scatter, `${key}-dots`, [])
      this.slot(rendered, DefaultZIndexes.activeDot, `${key}-active`, [])
      return
    }
    const rec = this.animationFor(item, JSON.stringify([points.map((p) => [p.x, p.y]), props.animationDuration, props.animationBegin, props.animationEasing]))
    const t = this.tOf(rec)
    const prev = rec.closurePrev
    const factor = prev ? prev.length / points.length : 1
    const step = t === 1 ? points : points.map(interpolatePoint(prev, factor, t))
    if (t > 0) rec.prev = step
    const polygonProps = { ...svgProps(omit(props, ["points", "className"])), points: step, connectNulls: props.connectNulls }
    let shape
    if (typeof props.shape === "function") shape = props.shape({ ...props, points: step })
    else shape = Polygon(polygonProps)
    this.slot(rendered, props.zIndex, key, Layer({ className: clsx("recharts-radar", props.className) }, Layer({ className: "recharts-radar-polygon" }, shape)))
    this.renderDots(rendered, `${key}-dots`, step, omit(props, ["id"]), "recharts-radar-dots", "recharts-radar-dot")
    this.renderActivePoints(state, rendered, `${key}-active`, points, getLegendItemColor(props.stroke, props.fill), props.dataKey, props.activeDot)
  }

  tooltipEntrySettings(item) {
    const p = item.props
    const name = p.name ? String(p.name) : typeof p.dataKey === "string" ? p.dataKey : undefined
    return { stroke: p.stroke, strokeWidth: p.strokeWidth, fill: p.fill, nameKey: undefined, dataKey: p.dataKey, name, hide: p.hide, type: p.tooltipType, color: getLegendItemColor(p.stroke, p.fill), unit: "", graphicalItemId: item.id }
  }

  legendPayload() {
    const payload = this.items.map((item) => {
      const p = item.props
      return { inactive: p.hide, dataKey: p.dataKey, type: p.legendType, color: getLegendItemColor(p.stroke, p.fill), value: p.name ? String(p.name) : typeof p.dataKey === "string" ? p.dataKey : undefined, payload: p }
    })
    const sorter = this.legendProps.itemSorter
    if (!sorter) return payload
    const keyOf = typeof sorter === "function" ? sorter : (e) => e[sorter]
    return payload
      .map((e, i) => [e, i])
      .sort(([a, i], [b, j]) => {
        const x = keyOf(a)
        const y = keyOf(b)
        return x < y ? -1 : x > y ? 1 : i - j
      })
      .map(([e]) => e)
  }

  // The polar cursor is not ported (the shadcn radar charts disable it).
  renderCursor() {}

  activePropsFromPointer(state, chartEvent) {
    const vb = state.viewBox
    if (!vb) return undefined
    const axisType = state.tooltipAxisType()
    const axis = state.axis(axisType, this.tooltipProps?.axisId ?? 0)
    const ticks = axis.tooltipTicks
    if (!ticks || !axis.range) return undefined
    const rangeObj = R.inRangeOfSector(chartEvent, vb)
    if (!rangeObj) return undefined
    const ordered = [...ticks].sort((a, b) => a.coordinate - b.coordinate)
    const pos = R.calculatePolarTooltipPos(rangeObj, state.layout)
    const activeIndex = R.calculateActiveTickIndex(pos, ordered, ticks, axisType, axis.range)
    const activeCoordinate = R.getActivePolarCoordinate(state.layout, ticks, activeIndex, rangeObj)
    return { activeIndex: String(activeIndex), activeCoordinate }
  }

  updateLineLengths() {
    return false
  }
}

registerChart("RadarChart", RadarChartInstance)

// chart-radar-label-custom's PolarAngleAxis tick.
registerFn("radarStackedTick", ([rows, a, b, label]) => ({ x, y, textAnchor, index, ...props }) => {
  const d = rows[index] ?? {}
  const yValue = typeof y === "number" ? y : 0
  return h(
    "text",
    // React writes the spread's unknown props too (verticalAnchor, payload).
    { x, y: yValue + (index === 0 ? -10 : 0), textAnchor, fontSize: 13, fontWeight: 500, ...svgProps(props), verticalAnchor: props.verticalAnchor, payload: props.payload == null ? undefined : String(props.payload) },
    h("tspan", null, String(d[a])),
    h("tspan", { className: "fill-muted-foreground" }, "/"),
    h("tspan", null, String(d[b])),
    h("tspan", { x, dy: "1rem", fontSize: 12, className: "fill-muted-foreground" }, String(d[label])),
  )
})
