// radial.js: recharts 3.8.0's RadialBarChart (RadialBar, PolarGrid,
// PolarRadiusAxis with a Label, LabelList on the bars, Tooltip) on top of
// the cartesian runtime of chart.js: the chart instance, tooltip, keyboard
// and animation machinery are CartesianChart's; this file ports the polar
// selectors (view box, radius band axis, numeric angle axis, stacks, bar
// positions) and the polar components (Sector, PolarGrid, radial labels).
import { CartesianChart, registerChart, registerFn, resolveSpecNode, syncElement, h, svgProps, withDefaults, DefaultZIndexes, Layer } from "./chart.js"
import * as R from "./recharts.js"

const isNum = R.isNumber

function clsx(...args) {
  return args.filter((a) => typeof a === "string" && a).join(" ")
}

const POLAR_CHART_DEFAULTS = { accessibilityLayer: true, barCategoryGap: "10%", barGap: 4, cx: "50%", cy: "50%", innerRadius: 0, outerRadius: "80%", layout: "radial", startAngle: 0, endAngle: 360, stackOffset: "none", reverseStackOrder: false }
const RADIAL_BAR_DEFAULTS = { angleAxisId: 0, animationBegin: 0, animationDuration: 1500, animationEasing: "ease", background: false, cornerIsExternal: false, cornerRadius: 0, forceCornerRadius: false, hide: false, isAnimationActive: "auto", label: false, legendType: "rect", minPointSize: 0, radiusAxisId: 0, zIndex: DefaultZIndexes.bar }
const POLAR_GRID_DEFAULTS = { angleAxisId: 0, radiusAxisId: 0, gridType: "polygon", radialLines: true, zIndex: DefaultZIndexes.grid }
const RADIUS_AXIS_DEFAULTS = { angle: 0, axisLine: true, tick: true, tickLine: true, zIndex: DefaultZIndexes.axis }
const LABEL_LIST_DEFAULTS = { angle: 0, offset: 5, zIndex: DefaultZIndexes.label, position: "middle", textBreakAll: false }
const TOOLTIP_DEFAULTS = { allowEscapeViewBox: { x: false, y: false }, animationDuration: 400, animationEasing: "ease", axisId: 0, contentStyle: {}, cursor: true, filterNull: true, includeHidden: false, isAnimationActive: "auto", itemSorter: "name", itemStyle: {}, labelStyle: {}, offset: 10, reverseDirection: { x: false, y: false }, separator: " : ", trigger: "hover", useTranslate3d: false, wrapperStyle: {} }

// --- Sector (shape/Sector.js) ---------------------------------------------------------

function roundT(strings, ...values) {
  return R.roundTemplateLiteral(strings, ...values)
}

function getDeltaAngle(startAngle, endAngle) {
  const sign = R.mathSign(endAngle - startAngle)
  return sign * Math.min(Math.abs(endAngle - startAngle), 359.999)
}

function getTangentCircle({ cx, cy, radius, angle, sign, isExternal, cornerRadius, cornerIsExternal }) {
  const centerRadius = cornerRadius * (isExternal ? 1 : -1) + radius
  const theta = Math.asin(cornerRadius / centerRadius) / R.RADIAN
  const centerAngle = cornerIsExternal ? angle : angle + sign * theta
  const center = R.polarToCartesian(cx, cy, centerRadius, centerAngle)
  const circleTangency = R.polarToCartesian(cx, cy, radius, centerAngle)
  const lineTangencyAngle = cornerIsExternal ? angle - sign * theta : angle
  const lineTangency = R.polarToCartesian(cx, cy, centerRadius * Math.cos(theta * R.RADIAN), lineTangencyAngle)
  return { center, circleTangency, lineTangency, theta }
}

function getSectorPath({ cx, cy, innerRadius, outerRadius, startAngle, endAngle }) {
  const angle = getDeltaAngle(startAngle, endAngle)
  const tempEndAngle = startAngle + angle
  const os = R.polarToCartesian(cx, cy, outerRadius, startAngle)
  const oe = R.polarToCartesian(cx, cy, outerRadius, tempEndAngle)
  let path = roundT`M ${os.x},${os.y}
    A ${outerRadius},${outerRadius},0,
    ${+(Math.abs(angle) > 180)},${+(startAngle > tempEndAngle)},
    ${oe.x},${oe.y}
  `
  if (innerRadius > 0) {
    const is = R.polarToCartesian(cx, cy, innerRadius, startAngle)
    const ie = R.polarToCartesian(cx, cy, innerRadius, tempEndAngle)
    path += roundT`L ${ie.x},${ie.y}
            A ${innerRadius},${innerRadius},0,
            ${+(Math.abs(angle) > 180)},${+(startAngle <= tempEndAngle)},
            ${is.x},${is.y} Z`
  } else {
    path += roundT`L ${cx},${cy} Z`
  }
  return path
}

function getSectorWithCorner({ cx, cy, innerRadius, outerRadius, cornerRadius, forceCornerRadius, cornerIsExternal, startAngle, endAngle }) {
  const sign = R.mathSign(endAngle - startAngle)
  const { circleTangency: soct, lineTangency: solt, theta: sot } = getTangentCircle({ cx, cy, radius: outerRadius, angle: startAngle, sign, cornerRadius, cornerIsExternal })
  const { circleTangency: eoct, lineTangency: eolt, theta: eot } = getTangentCircle({ cx, cy, radius: outerRadius, angle: endAngle, sign: -sign, cornerRadius, cornerIsExternal })
  const outerArcAngle = cornerIsExternal ? Math.abs(startAngle - endAngle) : Math.abs(startAngle - endAngle) - sot - eot
  if (outerArcAngle < 0) {
    if (forceCornerRadius) {
      return roundT`M ${solt.x},${solt.y}
        a${cornerRadius},${cornerRadius},0,0,1,${cornerRadius * 2},0
        a${cornerRadius},${cornerRadius},0,0,1,${-cornerRadius * 2},0
      `
    }
    return getSectorPath({ cx, cy, innerRadius, outerRadius, startAngle, endAngle })
  }
  let path = roundT`M ${solt.x},${solt.y}
    A${cornerRadius},${cornerRadius},0,0,${+(sign < 0)},${soct.x},${soct.y}
    A${outerRadius},${outerRadius},0,${+(outerArcAngle > 180)},${+(sign < 0)},${eoct.x},${eoct.y}
    A${cornerRadius},${cornerRadius},0,0,${+(sign < 0)},${eolt.x},${eolt.y}
  `
  if (innerRadius > 0) {
    const { circleTangency: sict, lineTangency: silt, theta: sit } = getTangentCircle({ cx, cy, radius: innerRadius, angle: startAngle, sign, isExternal: true, cornerRadius, cornerIsExternal })
    const { circleTangency: eict, lineTangency: eilt, theta: eit } = getTangentCircle({ cx, cy, radius: innerRadius, angle: endAngle, sign: -sign, isExternal: true, cornerRadius, cornerIsExternal })
    const innerArcAngle = cornerIsExternal ? Math.abs(startAngle - endAngle) : Math.abs(startAngle - endAngle) - sit - eit
    if (innerArcAngle < 0 && cornerRadius === 0) return `${path}L${cx},${cy}Z`
    path += roundT`L${eilt.x},${eilt.y}
      A${cornerRadius},${cornerRadius},0,0,${+(sign < 0)},${eict.x},${eict.y}
      A${innerRadius},${innerRadius},0,${+(innerArcAngle > 180)},${+(sign > 0)},${sict.x},${sict.y}
      A${cornerRadius},${cornerRadius},0,0,${+(sign < 0)},${silt.x},${silt.y}Z`
  } else {
    path += roundT`L${cx},${cy}Z`
  }
  return path
}

export function Sector(props) {
  const p = { cx: 0, cy: 0, innerRadius: 0, outerRadius: 0, startAngle: 0, endAngle: 0, cornerRadius: 0, forceCornerRadius: false, cornerIsExternal: false, ...props }
  const { cx, cy, innerRadius, outerRadius, cornerRadius, forceCornerRadius, cornerIsExternal, startAngle, endAngle, className } = p
  if (outerRadius < innerRadius || startAngle === endAngle) return null
  const deltaRadius = outerRadius - innerRadius
  const cr = R.getPercentValue(cornerRadius, deltaRadius, 0, true)
  let d
  if (cr > 0 && Math.abs(startAngle - endAngle) < 360) d = getSectorWithCorner({ cx, cy, innerRadius, outerRadius, cornerRadius: Math.min(cr, deltaRadius / 2), forceCornerRadius, cornerIsExternal, startAngle, endAngle })
  else d = getSectorPath({ cx, cy, innerRadius, outerRadius, startAngle, endAngle })
  return h("path", { ...svgProps(props), className: clsx("recharts-sector", className), d })
}

// RadialBarSector: Shape wraps the sector in a recharts-shape layer.
function radialBarSector(props) {
  const { option, isActive, index, ...rest } = props
  if (typeof option === "function") return Layer({ className: "recharts-shape" }, option(rest))
  return Layer({ className: "recharts-shape" }, Sector(rest))
}

// --- polar state -------------------------------------------------------------------------

class PolarState {
  constructor(chart) {
    this.chart = chart
    this.registered = chart.registered !== false
    this.layout = "radial"
    this.width = chart.width
    this.height = chart.height
    this.data = chart.data
    const m = chart.margin
    const offset = { top: m.top || 0, bottom: m.bottom || 0, left: m.left || 0, right: m.right || 0 }
    this.offset = { brushBottom: offset.bottom, ...offset, width: Math.max(this.width - offset.left - offset.right, 0), height: Math.max(this.height - offset.top - offset.bottom, 0) }
    const p = chart.props
    const maxRadius = R.getMaxRadius(this.width, this.height, this.offset)
    this.innerRadius = R.getPercentValue(p.innerRadius, maxRadius, 0)
    this.outerRadius = R.getPercentValue(p.outerRadius, maxRadius, maxRadius * 0.8)
    this.viewBox = { cx: R.getPercentValue(p.cx, this.width, this.width / 2), cy: R.getPercentValue(p.cy, this.height, this.height / 2), innerRadius: this.innerRadius, outerRadius: this.outerRadius, startAngle: p.startAngle, endAngle: p.endAngle, clockWise: false }
    const items = this.registered ? chart.items.filter((i) => !i.hide) : []
    this.items = items
    // Radius axis: categories (the data indexes), band scale.
    const n = this.data.length
    const radiusScale = R.scaleBand().domain(Array.from({ length: n }, (_, i) => i)).range([this.innerRadius, this.outerRadius])
    const bandSize = radiusScale.bandwidth()
    this.bandSize = bandSize
    this.radiusScale = radiusScale
    this.radiusTicks = Array.from({ length: n }, (_, i) => ({ coordinate: radiusScale(i) + bandSize / 2, value: i, index: i, offset: bandSize / 2 }))
    this.radiusItemTicks = Array.from({ length: n }, (_, i) => ({ coordinate: radiusScale(i), value: i, index: i, offset: 0 }))
    // Stacks.
    const groups = new Map()
    for (const item of items) {
      if (item.stackId == null) continue
      if (!groups.has(item.stackId)) groups.set(item.stackId, [])
      groups.get(item.stackId).push(item)
    }
    this.stacked = new Map()
    for (const [, list] of groups) {
      const series = R.getStackedData(this.data, list.map((i) => i.dataKey), p.stackOffset, p.reverseStackOrder)
      list.forEach((item, i) => this.stacked.set(item.id, series[i]))
    }
    // Angle axis: numbers, domain [min(0, values), max(values)]. Like
    // recharts (selectAllPolarAppliedNumericalValues), the raw values, not
    // the stacked ones: stacked sectors are truncated to this domain.
    let min = 0
    let max = 0
    for (const item of items) {
      for (let i = 0; i < n; i++) {
        const vals = [R.getValueByDataKey(this.data[i], item.dataKey)]
        for (const v of vals) {
          const x = Number(v)
          if (Number.isFinite(x)) {
            min = Math.min(min, x)
            max = Math.max(max, x)
          }
        }
      }
    }
    this.angleDomain = [min, max]
    this.angleScale = R.scaleLinear().domain(this.angleDomain).range([p.startAngle, p.endAngle])
    this.angleTicks = this.angleScale.ticks().map((v, i) => ({ coordinate: this.angleScale(v), value: v, index: i, offset: 0 }))
    // Bar positions: one bar per stack group or unstacked item.
    const sizeList = []
    const seen = new Set()
    for (const item of items) {
      const key = item.stackId ?? item.id
      if (seen.has(key)) continue
      seen.add(key)
      sizeList.push({ key, barSize: item.barSize })
    }
    const offset0 = R.getPercentValue(p.barCategoryGap, bandSize, 0, true)
    let gap = R.getPercentValue(p.barGap, bandSize, 0, true)
    const len = sizeList.length
    if ((len - 1) * gap >= bandSize - 2 * offset0) gap = 0
    let originalSize = (bandSize - 2 * offset0 - (len - 1) * gap) / len
    if (originalSize > 1) originalSize >>= 0
    this.positions = new Map(sizeList.map((s, i) => [s.key, { offset: offset0 + (originalSize + gap) * i, size: originalSize }]))
  }

  tooltipAxisType() {
    return "radiusAxis"
  }

  axis() {
    return { tooltipTicks: this.radiusTicks, settings: {}, range: [this.innerRadius, this.outerRadius] }
  }

  plotArea() {
    const o = this.offset
    return { x: o.left, y: o.top, width: o.width, height: o.height }
  }

  // computeRadialBarDataItems (layout radial).
  sectors(item) {
    const { cx, cy, startAngle: rootStart, endAngle: rootEnd } = this.viewBox
    const pos = this.positions.get(item.stackId ?? item.id)
    const st = this.stacked.get(item.id)
    return this.data.map((entry, index) => {
      let value = st ? R.truncateByDomain(st[index], this.angleDomain) : R.getValueByDataKey(entry, item.dataKey)
      if (!Array.isArray(value)) value = [0, value]
      let startAngle = this.angleScale(value[0]) ?? rootStart
      let endAngle = this.angleScale(value[1]) ?? rootEnd
      const innerRadius = this.radiusItemTicks[index].coordinate + pos.offset
      const outerRadius = innerRadius + pos.size
      const minPointSize = item.props.minPointSize
      const deltaAngle = endAngle - startAngle
      if (Math.abs(minPointSize) > 0 && Math.abs(deltaAngle) < Math.abs(minPointSize)) endAngle += R.mathSign(deltaAngle || minPointSize) * (Math.abs(minPointSize) - Math.abs(deltaAngle))
      return { ...entry, background: { cx, cy, innerRadius, outerRadius, startAngle: rootStart, endAngle: rootEnd }, payload: entry, value: st ? value : value[1], cx, cy, innerRadius, outerRadius, startAngle, endAngle }
    })
  }
}

// --- the chart ------------------------------------------------------------------------------

class RadialBarChart extends CartesianChart {
  setSpec(spec, initial = false) {
    const root = resolveSpecNode(spec, this.ctx)
    this.spec = spec
    this.name = root.type
    this.props = withDefaults(root.props, POLAR_CHART_DEFAULTS)
    this.margin = root.props.margin ?? { top: 5, right: 5, bottom: 5, left: 5 }
    this.data = root.props.data ?? []
    this.children = root.children
    const prevItems = this.items ?? []
    this.items = []
    this.xAxes = []
    this.yAxes = []
    this.grids = []
    this.radiusAxes = []
    this.tooltipProps = undefined
    this.tooltipContent = undefined
    this.legendProps = undefined
    let n = 1
    for (const child of root.children) {
      const p = child.props
      switch (child.type) {
        case "PolarGrid":
          child.settings = withDefaults(p, POLAR_GRID_DEFAULTS)
          break
        case "PolarRadiusAxis":
          child.settings = withDefaults(p, RADIUS_AXIS_DEFAULTS)
          break
        case "Tooltip":
          this.tooltipProps = withDefaults(p, TOOLTIP_DEFAULTS)
          this.tooltipContent = child.children.find((c) => c.type === "ChartTooltipContent")
          break
        case "RadialBar": {
          const props = withDefaults(p, RADIAL_BAR_DEFAULTS)
          const id = props.id ?? `recharts-radialBar-${this.ids.item(n++)}`
          this.items.push({
            type: "radialBar",
            id,
            props,
            node: child,
            dataKey: props.dataKey,
            stackId: props.stackId == null ? undefined : String(props.stackId),
            hide: props.hide,
            barSize: props.barSize,
            labelLists: child.children.filter((c) => c.type === "LabelList"),
            prev: prevItems.find((x) => x.id === id),
          })
          break
        }
      }
    }
    if (!initial) this.commit()
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
    syncElement(this.wrapper, { className: "recharts-wrapper", width: w, height: hgt }, { position: "relative", cursor: "default", width: w, height: hgt })
    this.ensureTooltipWrapper()
    this.ensureSvg()
    const acc = this.props.accessibilityLayer
    syncElement(this.svg, { cx: this.props.cx, cy: this.props.cy, role: acc ? "application" : undefined, tabIndex: acc ? 0 : undefined, className: "recharts-surface", width: w, height: hgt, viewBox: `0 0 ${w} ${hgt}` }, { width: "100%", height: "100%", display: "block" })
    const state = new PolarState(this)
    this.state = state
    const rendered = new Set()
    const pa = state.plotArea()
    this.setZeroLayer([h("defs", null, h("clipPath", { id: this.clipPathId }, h("rect", { x: pa.x, y: pa.y, height: pa.height, width: pa.width })))])
    let k = 0
    for (const child of this.children) {
      k++
      if (!state.registered) continue
      switch (child.type) {
        case "PolarGrid":
          this.renderPolarGrid(state, child.settings, `polar-grid-${k}`, rendered)
          break
        case "PolarRadiusAxis":
          this.renderRadiusAxis(state, child, `radius-axis-${k}`, rendered)
          break
        case "RadialBar": {
          const item = this.items.find((i) => i.node === child)
          if (item) this.renderRadialBar(state, item, rendered)
          break
        }
      }
    }
    for (const key of this.slots.keys()) if (!rendered.has(key)) this.slots.clear(key)
    this.renderTooltip(state)
  }

  renderPolarGrid(state, settings, key, rendered) {
    const { gridType, radialLines, angleAxisId, radiusAxisId, cx: _cx, cy: _cy, innerRadius: _ir, outerRadius: _or, polarAngles: anglesIn, polarRadius: radiusIn, zIndex, ...inputs } = settings
    const vb = state.viewBox
    const polarAngles = Array.isArray(anglesIn) ? anglesIn : state.angleTicks.map((t) => t.coordinate)
    const polarRadius = Array.isArray(radiusIn) ? radiusIn : state.radiusTicks.map((t) => t.coordinate)
    const props = { cx: vb.cx, cy: vb.cy, innerRadius: vb.innerRadius, outerRadius: vb.outerRadius, ...inputs }
    if (!(props.outerRadius > 0)) {
      this.slot(rendered, zIndex, key, [])
      return
    }
    const circle = (radius, fill) => h("circle", { stroke: "#ccc", fill: "none", ...svgProps({ ...props, fill: fill ?? props.fill }), radius, className: clsx("recharts-polar-grid-concentric-circle", props.className), cx: props.cx, cy: props.cy, r: radius })
    const polygon = (radius, fill) => {
      let d = ""
      polarAngles.forEach((angle, i) => {
        const pt = R.polarToCartesian(props.cx, props.cy, radius, angle)
        d += `${i ? "L" : "M"} ${pt.x},${pt.y}`
      })
      d += "Z"
      return h("path", { stroke: "#ccc", fill: "none", ...svgProps({ ...props, fill: fill ?? props.fill }), className: clsx("recharts-polar-grid-concentric-polygon", props.className), d })
    }
    const shape = gridType === "circle" ? circle : polygon
    const children = []
    if (polarRadius.length) {
      const renderBackground = props.fill && props.fill !== "none"
      children.push(h("g", { className: "recharts-polar-grid-concentric" }, renderBackground ? shape(Math.max(...polarRadius)) : null, polarRadius.map((r) => shape(r, "none"))))
    }
    if (polarAngles.length && radialLines) {
      const lineProps = { stroke: "#ccc", ...svgProps(props) }
      children.push(
        h(
          "g",
          { className: "recharts-polar-grid-angle" },
          polarAngles.map((a) => {
            const s = R.polarToCartesian(props.cx, props.cy, props.innerRadius, a)
            const e = R.polarToCartesian(props.cx, props.cy, props.outerRadius, a)
            return h("line", { ...lineProps, x1: s.x, y1: s.y, x2: e.x, y2: e.y })
          }),
        ),
      )
    }
    this.slot(rendered, zIndex, key, h("g", { className: "recharts-polar-grid" }, children))
  }

  renderRadiusAxis(state, child, key, rendered) {
    const s = child.settings
    const vb = state.viewBox
    // getViewBox(angle, cx, cy, ticks): the label's polar view box.
    const ticks = state.radiusTicks
    if (!ticks.length) {
      this.slot(rendered, s.zIndex, key, [])
      return
    }
    const coords = ticks.map((t) => t.coordinate)
    const viewBox = { cx: vb.cx, cy: vb.cy, startAngle: s.angle, endAngle: s.angle, innerRadius: Math.min(...coords), outerRadius: Math.max(...coords), clockWise: false }
    const labels = child.children
      .filter((c) => c.type === "Label")
      .map((c) => {
        const content = c.props.content
        if (typeof content === "function") return content({ ...c.props, viewBox })
        return null
      })
    this.slot(rendered, s.zIndex, key, Layer({ className: clsx("recharts-polar-radius-axis", "radiusAxis", s.className) }, labels))
  }

  renderRadialBar(state, item, rendered) {
    const props = item.props
    const key = item.id
    if (props.hide) {
      this.slot(rendered, props.zIndex, key, [])
      this.slot(rendered, DefaultZIndexes.barBackground, `${key}-bg`, [])
      this.slot(rendered, DefaultZIndexes.label, `${key}-labels`, [])
      return
    }
    const sectors = state.sectors(item)
    this.sectorsOf ??= new Map()
    this.sectorsOf.set(item.id, sectors)
    const rec = this.animationFor(item, JSON.stringify(sectors.map((s) => [s.innerRadius, s.outerRadius, s.startAngle, s.endAngle])))
    const t = this.tOf(rec)
    const prevData = rec.closurePrev
    rec.isNew = false
    const step =
      t === 1
        ? sectors
        : sectors.map((entry, index) => {
            const prev = prevData && prevData[index]
            if (prev) return { ...entry, startAngle: R.interpolate(prev.startAngle, entry.startAngle, t), endAngle: R.interpolate(prev.endAngle, entry.endAngle, t) }
            return { ...entry, endAngle: R.interpolate(entry.startAngle, entry.endAngle, t) }
          })
    if (t > 0) rec.prev = step
    const { shape, activeShape, cornerRadius, id, background, children, ...others } = props
    const baseProps = svgProps(others)
    // background (renderBackground), in its own layer
    if (background) {
      const bgProps = svgProps(typeof background === "object" ? background : {})
      const bgs = sectors.map((entry, i) => {
        const { value, background: bg, ...rest } = entry
        if (!bg) return null
        return radialBarSector({ cornerRadius, ...rest, fill: "#eee", ...bg, ...bgProps, index: i, className: clsx("recharts-radial-bar-background-sector", String(bgProps.className)), option: background, isActive: false })
      })
      this.slot(rendered, DefaultZIndexes.barBackground, `${key}-bg`, bgs)
    } else this.slot(rendered, DefaultZIndexes.barBackground, `${key}-bg`, [])
    const sectorNodes = step.map((entry, i) => radialBarSector({ ...baseProps, cornerRadius, ...entry, className: `recharts-radial-bar-sector ${entry.className}`, forceCornerRadius: others.forceCornerRadius, cornerIsExternal: others.cornerIsExternal, isActive: false, option: shape, index: i }))
    const layer = Layer({ className: clsx("recharts-area", props.className) }, background ? Layer({ className: "recharts-radial-bar-background" }) : null, Layer({ className: "recharts-radial-bar-sectors" }, sectorNodes))
    this.slot(rendered, props.zIndex, key, layer)
    // LabelList children, hidden while animating.
    const showLabels = !rec.isAnimating
    const lists = showLabels ? item.labelLists.map((ll) => renderRadialLabelList(withDefaults(ll.props, LABEL_LIST_DEFAULTS), step)) : []
    this.slot(rendered, DefaultZIndexes.label, `${key}-labels`, lists)
  }

  tooltipEntrySettings(item) {
    const p = item.props
    return { graphicalItemId: item.id, stroke: p.stroke, strokeWidth: p.strokeWidth, fill: p.fill, nameKey: undefined, dataKey: p.dataKey, name: p.name != null ? String(p.name) : typeof p.dataKey === "string" ? p.dataKey : undefined, hide: p.hide, type: p.tooltipType, color: p.fill, unit: "" }
  }

  interactionState() {
    const ts = this.tooltipState
    if (ts.hover.active) return ts.hover
    if (ts.keyboard.active) return ts.keyboard
    const di = this.tooltipProps?.defaultIndex
    if (ts.hover.index == null && di != null) return { active: true, index: String(di), coordinate: undefined }
    return { active: false, index: null, coordinate: ts.hover.coordinate }
  }

  activeTooltip(state) {
    const inter = this.interactionState()
    const data = this.data
    let index = null
    if (inter.index != null) {
      const n = Number(inter.index)
      if (!R.isWellBehavedNumber(n)) index = inter.index
      else index = String(Math.max(0, Math.min(n, data.length > 0 ? data.length - 1 : Infinity)))
    }
    const label = index == null ? undefined : state.radiusTicks[Number(index)]?.value
    let payload
    if (index != null) {
      payload = []
      for (const item of state.items) {
        const settings = this.tooltipEntrySettings(item)
        const sectors = this.sectorsOf?.get(item.id) ?? state.sectors(item)
        const entry = sectors[Number(index)]
        payload.push({ ...settings, payload: entry, value: R.getValueByDataKey(entry, settings.dataKey), name: settings.name })
      }
    }
    let coordinate = inter.coordinate
    if (coordinate == null && this.tooltipProps?.defaultIndex != null) coordinate = this.coordinateForIndex(state, String(this.tooltipProps.defaultIndex))
    return { active: inter.active && index != null, index, label, payload, coordinate, axis: state.axis() }
  }

  // combineCoordinateForDefaultIndex: the non-horizontal branch.
  coordinateForIndex(state, index) {
    const tick = state.radiusTicks[Number(index)]
    if (!tick) return undefined
    return { x: (state.offset.left + state.width) / 2, y: tick.coordinate }
  }

  renderCursor() {}

  activePropsFromPointer(state, chartEvent) {
    const rangeObj = R.inRangeOfSector(chartEvent, state.viewBox)
    if (!rangeObj) return undefined
    const ticks = state.radiusTicks
    if (!ticks.length) return undefined
    const pos = R.calculatePolarTooltipPos(rangeObj, "radial")
    const ordered = [...ticks].sort((a, b) => a.coordinate - b.coordinate)
    const activeIndex = R.calculateActiveTickIndex(pos, ordered, ticks, "radiusAxis", [state.innerRadius, state.outerRadius])
    const activeCoordinate = R.getActivePolarCoordinate("radial", ticks, activeIndex, rangeObj)
    return { activeIndex: String(activeIndex), activeCoordinate }
  }
}

// LabelList on a RadialBar: polar labels (insideStart/insideEnd/end on the
// arc, others at the sector's middle).
function renderRadialLabelList(listProps, sectors) {
  const { valueAccessor, dataKey, clockWise, id, textBreakAll, zIndex, position, offset, className, ...others } = listProps
  const labels = sectors.map((sector, index) => {
    const value = dataKey == null ? sector.value : R.getValueByDataKey(sector.payload, dataKey)
    const viewBox = { cx: sector.cx, cy: sector.cy, innerRadius: sector.innerRadius, outerRadius: sector.outerRadius, startAngle: sector.startAngle, endAngle: sector.endAngle, clockWise: false }
    const attrs = svgProps({ fill: listProps.fill ?? sector.fill, className, ...others, offset })
    return renderRadialLabel({ position, offset, className, label: value, attrs, viewBox })
  })
  return Layer({ className: "recharts-label-list" }, labels)
}

function renderRadialLabel({ position, offset, className, label, attrs, viewBox }) {
  const { cx, cy, innerRadius, outerRadius, startAngle, endAngle, clockWise } = viewBox
  const radius = (innerRadius + outerRadius) / 2
  const sign0 = R.mathSign(endAngle - startAngle)
  const deltaAngle = sign0 * Math.min(Math.abs(endAngle - startAngle), 360)
  const sign = deltaAngle >= 0 ? 1 : -1
  let labelAngle
  let direction
  if (position === "insideStart" || position === "insideEnd" || position === "end") {
    if (position === "insideStart") {
      labelAngle = startAngle + sign * offset
      direction = clockWise
    } else if (position === "insideEnd") {
      labelAngle = endAngle - sign * offset
      direction = !clockWise
    } else {
      labelAngle = endAngle + sign * offset
      direction = clockWise
    }
    direction = deltaAngle <= 0 ? direction : !direction
    const s = R.polarToCartesian(cx, cy, radius, labelAngle)
    const e = R.polarToCartesian(cx, cy, radius, labelAngle + (direction ? 1 : -1) * 359)
    const d = `M${s.x},${s.y}\n    A${radius},${radius},0,1,${direction ? 0 : 1},\n    ${e.x},${e.y}`
    const pid = R.uniqueId("recharts-radial-line-")
    const tp = h("textPath", { xlinkHref: `#${pid}` }, label)
    return h("text", { ...attrs, dominantBaseline: "central", className: clsx("recharts-radial-bar-label", className) }, h("defs", null, h("path", { id: pid, d })), tp)
  }
  const mid = (startAngle + endAngle) / 2
  const p = R.polarToCartesian(cx, cy, radius, mid)
  return h("text", { ...attrs, x: p.x, y: p.y, textAnchor: "middle", className: clsx("recharts-label", className) }, label == null ? "" : String(label))
}

// Label content of the shadcn radial charts: a <text> at the view box
// center with <tspan> lines: args = [textProps, [{dy, className, text}]].
registerFn("centerText", (args) => ({ viewBox }) => {
  if (!viewBox || !("cx" in viewBox) || !("cy" in viewBox)) return null
  const [textProps = {}, lines = []] = args
  return h(
    "text",
    { x: viewBox.cx, y: viewBox.cy, ...textProps },
    lines.map((l) => h("tspan", { x: viewBox.cx, y: (viewBox.cy || 0) + (l.dy ?? 0), className: l.className }, l.text)),
  )
})

registerChart("RadialBarChart", RadialBarChart)
