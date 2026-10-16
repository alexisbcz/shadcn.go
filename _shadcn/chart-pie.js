// pie.js: recharts 3.8.0's PieChart (polar/Pie.js, shape/Sector.js, the
// polar parts of component/Label.js and LabelList.js, pieSelectors) on top
// of the chart runtime in chart.js. PieChart reuses CartesianChart's
// skeleton (wrapper, surface, z-index layers, clip path, tooltip, legend,
// animations, events) and replaces what is polar: the items, the item
// tooltip ("item" event type: hovering a sector), the legend payload (one
// entry per sector) and the keyboard tooltip coordinates.
import { CartesianChart, registerChart, registerFn, resolveSpecNode, h, svgProps, withDefaults, DefaultZIndexes, Text, Curve, Layer } from "./chart.js"
import * as R from "./recharts.js"

function clsx(...args) {
  return args.filter((a) => a && typeof a === "string").join(" ")
}

function omit(obj, keys) {
  const out = { ...obj }
  for (const k of keys) delete out[k]
  return out
}

const RADIAN = Math.PI / 180
export const polarToCartesian = (cx, cy, radius, angle) => ({ x: cx + Math.cos(-RADIAN * angle) * radius, y: cy + Math.sin(-RADIAN * angle) * radius })
const getMaxRadius = (width, height) => Math.min(Math.abs(width), Math.abs(height)) / 2
const roundT = (strings, ...values) => R.roundTemplateLiteral(strings, ...values)

// --- Sector (shape/Sector.js) ---------------------------------------------------------

function sectorDeltaAngle(startAngle, endAngle) {
  const sign = R.mathSign(endAngle - startAngle)
  return sign * Math.min(Math.abs(endAngle - startAngle), 359.999)
}

function getTangentCircle({ cx, cy, radius, angle, sign, isExternal, cornerRadius, cornerIsExternal }) {
  const centerRadius = cornerRadius * (isExternal ? 1 : -1) + radius
  const theta = Math.asin(cornerRadius / centerRadius) / RADIAN
  const centerAngle = cornerIsExternal ? angle : angle + sign * theta
  const center = polarToCartesian(cx, cy, centerRadius, centerAngle)
  const circleTangency = polarToCartesian(cx, cy, radius, centerAngle)
  const lineTangencyAngle = cornerIsExternal ? angle - sign * theta : angle
  const lineTangency = polarToCartesian(cx, cy, centerRadius * Math.cos(theta * RADIAN), lineTangencyAngle)
  return { center, circleTangency, lineTangency, theta }
}

function getSectorPath({ cx, cy, innerRadius, outerRadius, startAngle, endAngle }) {
  const angle = sectorDeltaAngle(startAngle, endAngle)
  const tempEndAngle = startAngle + angle
  const outerStartPoint = polarToCartesian(cx, cy, outerRadius, startAngle)
  const outerEndPoint = polarToCartesian(cx, cy, outerRadius, tempEndAngle)
  let path = roundT`M ${outerStartPoint.x},${outerStartPoint.y}
    A ${outerRadius},${outerRadius},0,
    ${+(Math.abs(angle) > 180)},${+(startAngle > tempEndAngle)},
    ${outerEndPoint.x},${outerEndPoint.y}
  `
  if (innerRadius > 0) {
    const innerStartPoint = polarToCartesian(cx, cy, innerRadius, startAngle)
    const innerEndPoint = polarToCartesian(cx, cy, innerRadius, tempEndAngle)
    path += roundT`L ${innerEndPoint.x},${innerEndPoint.y}
            A ${innerRadius},${innerRadius},0,
            ${+(Math.abs(angle) > 180)},${+(startAngle <= tempEndAngle)},
            ${innerStartPoint.x},${innerStartPoint.y} Z`
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

const SECTOR_DEFAULTS = { cx: 0, cy: 0, innerRadius: 0, outerRadius: 0, startAngle: 0, endAngle: 0, cornerRadius: 0, forceCornerRadius: false, cornerIsExternal: false }

/** recharts <Sector>: a <path> or null. */
export function Sector(sectorProps) {
  const props = withDefaults(sectorProps, SECTOR_DEFAULTS)
  const { cx, cy, innerRadius, outerRadius, cornerRadius, forceCornerRadius, cornerIsExternal, startAngle, endAngle, className } = props
  if (outerRadius < innerRadius || startAngle === endAngle) return null
  const deltaRadius = outerRadius - innerRadius
  const cr = R.getPercentValue(cornerRadius, deltaRadius, 0, true)
  const path =
    cr > 0 && Math.abs(startAngle - endAngle) < 360
      ? getSectorWithCorner({ cx, cy, innerRadius, outerRadius, cornerRadius: Math.min(cr, deltaRadius / 2), forceCornerRadius, cornerIsExternal, startAngle, endAngle })
      : getSectorPath({ cx, cy, innerRadius, outerRadius, startAngle, endAngle })
  return h("path", { ...svgProps(props), className: clsx("recharts-sector", className), d: path })
}

// util/ActiveShapeUtils Shape for shapeType "sector".
function sectorShape(option, props) {
  let shape
  if (typeof option === "function") shape = option(props, props.index)
  else if (option && typeof option === "object" && !Array.isArray(option)) shape = Sector({ ...props, ...option })
  else shape = Sector(props)
  return Layer({ className: props.isActive ? "recharts-active-shape" : "recharts-shape" }, shape)
}

// --- polar labels (component/Label.js) --------------------------------------------------

function getAttrsOfPolarLabel(viewBox, offset, position) {
  const { cx, cy, innerRadius, outerRadius, startAngle, endAngle } = viewBox
  const midAngle = (startAngle + endAngle) / 2
  if (position === "outside") {
    const { x, y } = polarToCartesian(cx, cy, outerRadius + offset, midAngle)
    return { x, y, textAnchor: x >= cx ? "start" : "end", verticalAnchor: "middle" }
  }
  if (position === "center") return { x: cx, y: cy, textAnchor: "middle", verticalAnchor: "middle" }
  if (position === "centerTop") return { x: cx, y: cy, textAnchor: "middle", verticalAnchor: "start" }
  if (position === "centerBottom") return { x: cx, y: cy, textAnchor: "middle", verticalAnchor: "end" }
  const r = (innerRadius + outerRadius) / 2
  const { x, y } = polarToCartesian(cx, cy, r, midAngle)
  return { x, y, textAnchor: "middle", verticalAnchor: "middle" }
}

const LABEL_DEFAULTS = { angle: 0, offset: 5, zIndex: DefaultZIndexes.label, position: "middle", textBreakAll: false }
const TEXT_ANCHORS = ["start", "middle", "end", "inherit"]

/**
 * recharts <Label> on a polar view box. Returns {z, node}: a content
 * function's element renders in place (z 0), a Text in its z-index layer.
 */
function renderPolarLabel(outerProps) {
  const props = withDefaults(outerProps, LABEL_DEFAULTS)
  const { viewBox, value, children, content, className = "", textBreakAll, formatter } = props
  if (!viewBox || (value == null && children == null && typeof content !== "function")) return null
  let label
  if (typeof content === "function") {
    const out = content(omit(props, ["content"]))
    if (out instanceof Node) return { z: 0, node: out }
    label = out
  } else {
    label = children == null ? value : children
    if (typeof formatter === "function") label = formatter(label)
  }
  const attrs = svgProps(props)
  const pos = getAttrsOfPolarLabel(viewBox, props.offset, props.position)
  const node = Text({ className: clsx("recharts-label", className), ...attrs, ...pos, textAnchor: TEXT_ANCHORS.includes(attrs.textAnchor) ? attrs.textAnchor : pos.textAnchor, breakAll: textBreakAll, children: label })
  return { z: props.zIndex, node }
}

// --- Pie ---------------------------------------------------------------------------------

const PIE_DEFAULTS = { animationBegin: 400, animationDuration: 1500, animationEasing: "ease", cx: "50%", cy: "50%", dataKey: "value", endAngle: 360, fill: "#808080", hide: false, innerRadius: 0, isAnimationActive: "auto", label: false, labelLine: true, legendType: "rect", minAngle: 0, nameKey: "name", outerRadius: "80%", paddingAngle: 0, rootTabIndex: 0, startAngle: 0, stroke: "#fff", zIndex: DefaultZIndexes.area }
const PIE_CHART_DEFAULTS = { accessibilityLayer: true, stackOffset: "none", barCategoryGap: "10%", barGap: 4, reverseStackOrder: false, syncMethod: "index", layout: "centric", responsive: false, cx: "50%", cy: "50%", innerRadius: 0, outerRadius: "80%", startAngle: 0, endAngle: 360 }
const TOOLTIP_DEFAULTS = { allowEscapeViewBox: { x: false, y: false }, animationDuration: 400, animationEasing: "ease", axisId: 0, contentStyle: {}, cursor: true, filterNull: true, includeHidden: false, isAnimationActive: "auto", itemSorter: "name", itemStyle: {}, labelStyle: {}, offset: 10, reverseDirection: { x: false, y: false }, separator: " : ", trigger: "hover", useTranslate3d: false, wrapperStyle: {} }
const LEGEND_DEFAULTS = { align: "center", iconSize: 14, inactiveColor: "#ccc", itemSorter: "value", layout: "horizontal", verticalAlign: "bottom" }
const DEFAULT_MARGIN = { top: 5, right: 5, bottom: 5, left: 5 }

const getTextAnchor = (x, cx) => (x > cx ? "start" : x < cx ? "end" : "middle")

function parseDeltaAngle(startAngle, endAngle) {
  const sign = R.mathSign(endAngle - startAngle)
  return sign * Math.min(Math.abs(endAngle - startAngle), 360)
}

function parseCoordinateOfPie(s, offset) {
  const { top, left, width, height } = offset
  const maxPieRadius = getMaxRadius(width, height)
  return {
    cx: left + R.getPercentValue(s.cx, width, width / 2),
    cy: top + R.getPercentValue(s.cy, height, height / 2),
    innerRadius: R.getPercentValue(s.innerRadius, maxPieRadius, 0),
    outerRadius: R.getPercentValue(s.outerRadius, maxPieRadius, maxPieRadius * 0.8),
    maxRadius: s.maxRadius || Math.sqrt(width * width + height * height) / 2,
  }
}

/** polar/Pie computePieSectors. */
function computePieSectors({ pieSettings, displayedData, cells, offset }) {
  const { cornerRadius, startAngle, endAngle, dataKey, nameKey, tooltipType } = pieSettings
  const minAngle = Math.abs(pieSettings.minAngle)
  const deltaAngle = parseDeltaAngle(startAngle, endAngle)
  const absDeltaAngle = Math.abs(deltaAngle)
  const paddingAngle = displayedData.length <= 1 ? 0 : (pieSettings.paddingAngle ?? 0)
  const notZeroItemCount = displayedData.filter((entry) => R.getValueByDataKey(entry, dataKey, 0) !== 0).length
  const totalPaddingAngle = (absDeltaAngle >= 360 ? notZeroItemCount : notZeroItemCount - 1) * paddingAngle
  const realTotalAngle = absDeltaAngle - notZeroItemCount * minAngle - totalPaddingAngle
  const sum = displayedData.reduce((result, entry) => {
    const val = R.getValueByDataKey(entry, dataKey, 0)
    return result + (R.isNumber(val) ? val : 0)
  }, 0)
  if (!(sum > 0)) return undefined
  let prev
  return displayedData.map((entry, i) => {
    const val = R.getValueByDataKey(entry, dataKey, 0)
    const name = R.getValueByDataKey(entry, nameKey, i)
    const coordinate = parseCoordinateOfPie(pieSettings, offset)
    const percent = (R.isNumber(val) ? val : 0) / sum
    const entryWithCellInfo = { ...entry, ...(cells && cells[i] && cells[i].props) }
    const sectorColor = entryWithCellInfo != null && "fill" in entryWithCellInfo && typeof entryWithCellInfo.fill === "string" ? entryWithCellInfo.fill : pieSettings.fill
    const tempStartAngle = i ? prev.endAngle + R.mathSign(deltaAngle) * paddingAngle * (val !== 0 ? 1 : 0) : startAngle
    const tempEndAngle = tempStartAngle + R.mathSign(deltaAngle) * ((val !== 0 ? minAngle : 0) + percent * realTotalAngle)
    const midAngle = (tempStartAngle + tempEndAngle) / 2
    const middleRadius = (coordinate.innerRadius + coordinate.outerRadius) / 2
    const tooltipPayload = [{ name, value: val, payload: entryWithCellInfo, dataKey, type: tooltipType, color: sectorColor, fill: sectorColor, graphicalItemId: pieSettings.id }]
    const tooltipPosition = polarToCartesian(coordinate.cx, coordinate.cy, middleRadius, midAngle)
    prev = {
      ...pieSettings.presentationProps,
      percent,
      cornerRadius: typeof cornerRadius === "string" ? parseFloat(cornerRadius) : cornerRadius,
      name,
      tooltipPayload,
      midAngle,
      middleRadius,
      tooltipPosition,
      ...entryWithCellInfo,
      ...coordinate,
      value: val,
      dataKey,
      startAngle: tempStartAngle,
      endAngle: tempEndAngle,
      payload: entryWithCellInfo,
      paddingAngle: R.mathSign(deltaAngle) * paddingAngle,
    }
    return prev
  })
}

/**
 * A mounted PieChart: CartesianChart's machinery with polar items. The
 * render loop calls renderItem for each Pie (in children order).
 */
export class PieChart extends CartesianChart {
  setSpec(spec, initial = false) {
    const root = resolveSpecNode(spec, this.ctx)
    this.spec = spec
    this.name = root.type
    this.props = withDefaults(root.props, PIE_CHART_DEFAULTS)
    this.margin = root.props.margin ?? DEFAULT_MARGIN
    this.data = root.props.data ?? []
    this.children = root.children
    this.specVersion = (this.specVersion ?? 0) + 1
    this.items = []
    this.xAxes = []
    this.yAxes = []
    this.grids = []
    this.defs = []
    this.tooltipProps = undefined
    this.tooltipContent = undefined
    this.legendProps = undefined
    this.sectors ??= new Map()
    let n = 1
    for (const child of root.children) {
      const p = child.props
      switch (child.type) {
        case "Tooltip":
          this.tooltipProps = withDefaults(p, TOOLTIP_DEFAULTS)
          this.tooltipContent = child.children.find((c) => c.type === "ChartTooltipContent")
          break
        case "Legend":
          this.legendProps = withDefaults(p, LEGEND_DEFAULTS)
          this.legendContent = child.children.find((c) => c.type === "ChartLegendContent")
          break
        case "Pie": {
          const props = withDefaults(p, PIE_DEFAULTS)
          const id = props.id ?? `recharts-pie-${this.ids.item(n++)}`
          this.items.push({
            type: "pie",
            id,
            props,
            node: child,
            dataKey: props.dataKey,
            hide: props.hide,
            cells: child.children.filter((c) => c.type === "Cell"),
          })
          break
        }
        case "defs":
          this.defs.push(child)
          break
      }
    }
    if (!initial) this.commit()
  }

  displayedData(item) {
    const p = item.props
    let data = p.data != null && p.data.length > 0 ? p.data : this.data
    if ((!data || !data.length) && item.cells.length) data = item.cells.map((cell) => ({ ...svgProps(omit(p, ["id"])), ...cell.props }))
    return data
  }

  pieSettings(item) {
    const p = item.props
    return {
      id: item.id,
      name: p.name,
      dataKey: p.dataKey,
      nameKey: p.nameKey,
      tooltipType: p.tooltipType,
      legendType: p.legendType,
      fill: p.fill,
      cx: p.cx,
      cy: p.cy,
      startAngle: p.startAngle,
      endAngle: p.endAngle,
      paddingAngle: p.paddingAngle,
      minAngle: p.minAngle,
      innerRadius: p.innerRadius,
      outerRadius: p.outerRadius,
      cornerRadius: p.cornerRadius,
      presentationProps: svgProps(omit(p, ["id"])),
      maxRadius: p.maxRadius,
    }
  }

  render() {
    super.render()
    // CategoricalChart spreads the chart's svg props (cx, cy) on the surface.
    if (this.svg) {
      for (const k of ["cx", "cy"]) {
        const v = this.props[k]
        if (v == null) this.svg.removeAttribute(k)
        else if (this.svg.getAttribute(k) !== String(v)) this.svg.setAttribute(k, String(v))
      }
    }
    // The keyboard handlers of CartesianChart read the tooltip axis: for a
    // PieChart, the angle axis, whose ticks are the pies' data entries.
    const state = this.state
    if (state) {
      const count = this.items.reduce((n, item) => n + (this.displayedData(item)?.length ?? 0), 0)
      state.tooltipAxisType = () => "angleAxis"
      state.axis = () => ({ tooltipTicks: Array.from({ length: count }, (_, i) => ({ index: i, coordinate: i, value: i })) })
    }
  }

  // Item mode: Tooltip draws no cursor.
  renderCursor(state, rendered) {
    for (const z of [DefaultZIndexes.cursorLine, DefaultZIndexes.cursorRectangle]) this.slot(rendered, z, `cursor-${z}`, [])
  }

  renderItem(state, item, rendered) {
    if (item.type !== "pie") return
    const props = item.props
    const key = item.id
    const displayed = this.displayedData(item)
    const sectors = displayed == null ? undefined : computePieSectors({ pieSettings: this.pieSettings(item), displayedData: displayed, cells: item.cells, offset: state.offset })
    this.sectors.set(item.id, sectors)
    const labelKeys = [`${key}-labels`, `${key}-labellists`, `${key}-label-children`]
    if (props.hide || sectors == null) {
      this.slot(rendered, props.zIndex, key, Layer({ tabIndex: props.rootTabIndex, className: clsx("recharts-pie", props.className) }))
      for (const k of labelKeys) this.slot(rendered, DefaultZIndexes.label, k, [])
      return
    }
    const inputKey = JSON.stringify([this.specVersion, sectors.map((s) => [s.cx, s.cy, s.innerRadius, s.outerRadius, s.startAngle, s.endAngle])])
    const rec = this.animationFor(item, inputKey)
    const t = this.tOf(rec)
    const prevSectors = rec.closurePrev
    rec.isNew = false
    // SectorsWithAnimation
    const stepData = []
    let curAngle = sectors[0]?.startAngle ?? 0
    sectors.forEach((entry, index) => {
      const prev = prevSectors && prevSectors[index]
      const paddingAngle = index > 0 ? (entry.paddingAngle ?? 0) : 0
      const angle = prev ? R.interpolate(prev.endAngle - prev.startAngle, entry.endAngle - entry.startAngle, t) : R.interpolate(0, entry.endAngle - entry.startAngle, t)
      const latest = { ...entry, startAngle: curAngle + paddingAngle, endAngle: curAngle + angle + paddingAngle }
      stepData.push(latest)
      curAngle = latest.endAngle
    })
    rec.prev = stepData
    // PieSectors
    const inter = this.chosenInteraction()
    const activeIndex = inter.active ? inter.index : null
    const sectorNodes = stepData.map((entry, i) => {
      if (entry?.startAngle === 0 && entry?.endAngle === 0 && stepData.length !== 1) return null
      const graphicalItemMatches = inter.graphicalItemId == null || inter.graphicalItemId === item.id
      const isActive = String(i) === activeIndex && (inter.dataKey == null || props.dataKey === inter.dataKey) && graphicalItemMatches
      const inactiveShape = activeIndex ? props.inactiveShape : null
      const sectorOptions = props.activeShape && isActive ? props.activeShape : inactiveShape
      const sectorProps = { ...entry, stroke: entry.stroke, tabIndex: -1, "data-recharts-item-index": i, "data-recharts-item-id": item.id }
      return Layer({ tabIndex: -1, className: "recharts-pie-sector" }, sectorShape(props.shape ?? sectorOptions, { index: i, isActive, ...sectorProps }))
    })
    const showLabels = !rec.isAnimating
    const inPlace = [Layer({}, sectorNodes)]
    // PieLabelList / PieLabels: each label in the label layer.
    const label = props.label
    const labelNodes = []
    if (showLabels && label) {
      if (label && typeof label === "object" && !Array.isArray(label) && "position" in label) {
        labelNodes.push(this.polarLabelList(label, sectors))
      } else {
        const pieProps = svgProps({ ...props, sectors })
        const customLabelProps = label && typeof label === "object" ? svgProps(label) : {}
        const labelLine = props.labelLine
        const customLabelLineProps = labelLine && typeof labelLine === "object" ? svgProps(labelLine) : {}
        const offsetRadius = (label && typeof label === "object" && typeof label.offsetRadius === "number" && label.offsetRadius) || 20
        sectors.forEach((entry, i) => {
          const midAngle = (entry.startAngle + entry.endAngle) / 2
          const endPoint = polarToCartesian(entry.cx, entry.cy, entry.outerRadius + offsetRadius, midAngle)
          const labelProps = { ...pieProps, ...entry, stroke: "none", ...customLabelProps, index: i, textAnchor: getTextAnchor(endPoint.x, entry.cx), ...endPoint }
          const lineProps = { ...pieProps, ...entry, fill: "none", stroke: entry.fill, ...customLabelLineProps, index: i, points: [polarToCartesian(entry.cx, entry.cy, entry.outerRadius, midAngle), endPoint] }
          let line = null
          if (labelLine) {
            if (typeof labelLine === "function") line = labelLine(lineProps)
            else line = Curve({ ...lineProps, type: "linear", className: clsx("recharts-pie-label-line", typeof labelLine !== "boolean" ? labelLine.className : "") })
          }
          let text
          let value = R.getValueByDataKey(entry, props.dataKey)
          if (typeof label === "function") {
            const out = label(labelProps)
            if (out instanceof Node) text = out
            else value = out
          }
          if (!text) text = Text({ ...labelProps, alignmentBaseline: "middle", className: clsx("recharts-pie-label-text", typeof label === "object" ? label.className : ""), children: value })
          labelNodes.push(Layer({}, line, text))
        })
      }
    }
    if (showLabels && label && !(typeof label === "object" && "position" in label)) inPlace.push(Layer({ className: "recharts-pie-labels" }))
    this.slot(rendered, DefaultZIndexes.label, `${key}-labels`, labelNodes)
    // Children: Label (polar view box of the chart) and LabelList (sectors).
    const childLabels = []
    const childLists = []
    for (const child of item.node.children) {
      if (child.type === "Label") {
        const r = renderPolarLabel({ ...child.props, viewBox: this.polarViewBox() })
        if (!r) continue
        if (r.z === 0) inPlace.push(r.node)
        else childLabels.push(r.node)
      } else if (child.type === "LabelList") {
        if (showLabels) childLists.push(this.polarLabelList(child.props, sectors))
      }
    }
    this.slot(rendered, DefaultZIndexes.label, `${key}-label-children`, childLabels)
    this.slot(rendered, DefaultZIndexes.label, `${key}-labellists`, childLists)
    const [pieLayer] = this.slot(rendered, props.zIndex, key, Layer({ tabIndex: props.rootTabIndex, className: clsx("recharts-pie", props.className) }, ...inPlace))
    // React keys each sector by its angles: a sector whose angles change
    // (every animation frame) is a new element, so a hovered sector that
    // is replaced never receives mouseout (and the item stays hovered).
    const keys = stepData.filter((entry) => !(entry?.startAngle === 0 && entry?.endAngle === 0 && stepData.length !== 1)).map((entry, i) => `sector-${entry.startAngle}-${entry.endAngle}-${entry.midAngle}`)
    const groups = pieLayer?.firstElementChild ? [...pieLayer.firstElementChild.children] : []
    groups.forEach((g, i) => {
      const k = keys[i]
      if (g.__sectorKey !== undefined && g.__sectorKey !== k) {
        const fresh = g.cloneNode(true)
        g.replaceWith(fresh)
        g = fresh
      }
      g.__sectorKey = k
    })
  }

  // LabelList with the polar entries of the sectors.
  polarLabelList(listProps, sectors) {
    const { valueAccessor, dataKey, clockWise, id, textBreakAll, zIndex, ...others } = listProps
    const labels = sectors.map((s, index) => {
      const entry = { value: s.value, payload: s.payload, clockWise: false, parentViewBox: undefined, viewBox: { cx: s.cx, cy: s.cy, innerRadius: s.innerRadius, outerRadius: s.outerRadius, startAngle: s.startAngle, endAngle: s.endAngle, clockWise: false }, fill: s.fill }
      let value
      if (dataKey == null) {
        const v = Array.isArray(entry.value) ? entry.value[entry.value.length - 1] : entry.value
        value = v == null || typeof v === "string" || typeof v === "number" || typeof v === "boolean" ? v : undefined
      } else value = R.getValueByDataKey(entry.payload, dataKey)
      const idProps = id == null ? {} : { id: `${id}-${index}` }
      const r = renderPolarLabel({ ...svgProps(entry), ...others, ...idProps, fill: listProps.fill ?? entry.fill, parentViewBox: entry.parentViewBox, value, textBreakAll, viewBox: entry.viewBox, index, zIndex: 0 })
      return r?.node ?? null
    })
    return Layer({ className: "recharts-label-list" }, labels)
  }

  // selectPolarViewBox
  polarViewBox() {
    const p = this.props
    return { cx: R.getPercentValue(p.cx, this.width, this.width / 2), cy: R.getPercentValue(p.cy, this.height, this.height / 2), innerRadius: p.innerRadius, outerRadius: p.outerRadius, startAngle: p.startAngle, endAngle: p.endAngle, clockWise: false }
  }

  // --- tooltip (item event type) ---------------------------------------------------

  tooltipEventType() {
    return "item"
  }

  itemHover() {
    return this.tooltipState.item ?? { active: false, index: null, dataKey: undefined, graphicalItemId: undefined, coordinate: undefined }
  }

  // combineTooltipInteractionState for "item" + hover.
  chosenInteraction() {
    const hover = this.itemHover()
    if (hover.active) return hover
    if (this.tooltipState.keyboard.active) return this.tooltipState.keyboard
    return { active: false, index: null, dataKey: undefined, graphicalItemId: undefined, coordinate: hover.coordinate }
  }

  interactionState() {
    return this.chosenInteraction()
  }

  // combineTooltipPayloadConfigurations: the hovered item only.
  tooltipConfigurations() {
    const id = this.itemHover().graphicalItemId
    return this.items.filter((item) => item.id === id)
  }

  activeTooltip() {
    const inter = this.chosenInteraction()
    const index = inter.index == null ? null : String(inter.index)
    let payload
    if (index != null) {
      payload = []
      for (const item of this.tooltipConfigurations()) {
        const settings = this.tooltipEntrySettings(item)
        const sectors = this.sectors.get(item.id) ?? []
        const n = Number.parseInt(index, 10)
        const found = Number.isNaN(n) ? undefined : sectors[n]?.tooltipPayload
        if (Array.isArray(found)) {
          for (const it of found) {
            const s = { ...settings, name: it.name, unit: undefined, color: it.color ?? settings.color, fill: it.fill ?? settings.fill }
            payload.push({ ...s, dataKey: it.dataKey, payload: it.payload, value: R.getValueByDataKey(it.payload, it.dataKey), name: it.name == null ? undefined : String(it.name) })
          }
        } else {
          payload.push({ ...settings, dataKey: settings.dataKey, payload: found, value: R.getValueByDataKey(found, settings.dataKey), name: R.getValueByDataKey(found, settings.nameKey) ?? settings.name })
        }
      }
    }
    return { active: inter.active && index != null, index, label: undefined, payload, coordinate: inter.coordinate }
  }

  tooltipEntrySettings(item) {
    const p = item.props
    return { stroke: p.stroke, strokeWidth: p.strokeWidth, fill: p.fill, dataKey: p.dataKey, nameKey: p.nameKey, name: R.getTooltipNameProp(p.name, p.dataKey), hide: p.hide, type: p.tooltipType, color: p.fill, unit: "", graphicalItemId: item.id }
  }

  // selectCoordinateForDefaultIndex: the first configuration's position.
  // (selected with the "axis" event type: every item, so the first pie).
  coordinateForIndex(state, index) {
    const item = this.items[0]
    return item ? this.sectors.get(item.id)?.[Number(index)]?.tooltipPosition : undefined
  }

  activePropsFromPointer() {
    return undefined
  }

  // --- legend ----------------------------------------------------------------------

  legendPayload() {
    const payload = []
    for (const item of this.items) {
      const p = item.props
      const data = this.displayedData(item)
      if (data == null) continue
      data.forEach((entry, i) => {
        const name = R.getValueByDataKey(entry, p.nameKey, p.name)
        const cell = item.cells[i]
        const color = cell?.props?.fill ? cell.props.fill : typeof entry === "object" && entry != null && "fill" in entry ? entry.fill : p.fill
        payload.push({ value: R.getTooltipNameProp(name, p.dataKey), color, payload: entry, type: p.legendType })
      })
    }
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

  // --- events: sectors' item enter/leave ---------------------------------------------

  attachEvents() {
    super.attachEvents()
    const w = this.wrapper
    // React emulates mouseleave from mouseout/mouseover pairs: when the
    // element last under the pointer was removed (a sector replaced during
    // the animation), no mouseout reaches React and the chart never sees
    // the pointer leave. Swallow the native mouseleave in that case.
    let lastOver = null
    w.addEventListener("mouseover", (e) => (lastOver = e.target), true)
    w.addEventListener("mouseout", (e) => (lastOver = e.relatedTarget), true)
    this.inner.addEventListener(
      "mouseleave",
      (e) => {
        if (e.target === w && lastOver && !lastOver.isConnected) e.stopPropagation()
      },
      true,
    )
    const sectorOf = (el) => (el instanceof Element ? el.closest(".recharts-pie-sector") : null)
    const enter = (g) => {
      const path = g.querySelector("[data-recharts-item-index]")
      if (!path) return
      const id = path.getAttribute("data-recharts-item-id")
      const index = Number(path.getAttribute("data-recharts-item-index"))
      const item = this.items.find((i) => i.id === id)
      if (!item) return
      const sector = this.sectors.get(id)?.[index]
      this.tooltipState.item = { active: true, index: String(index), dataKey: item.props.dataKey, graphicalItemId: id, coordinate: sector?.tooltipPosition }
      this.tooltipState.keyboard = { ...this.tooltipState.keyboard, active: false }
      this.render()
    }
    const leave = () => {
      if (!this.tooltipState.item?.active) return
      this.tooltipState.item = { ...this.tooltipState.item, active: false }
      this.render()
    }
    w.addEventListener("mouseout", (e) => {
      const from = sectorOf(e.target)
      const to = w.contains(e.relatedTarget) ? sectorOf(e.relatedTarget) : null
      if (from === to) return
      if (from) leave()
      if (to) enter(to)
    })
    w.addEventListener("mouseover", (e) => {
      if (e.relatedTarget && w.contains(e.relatedTarget)) return
      const to = sectorOf(e.target)
      if (to) enter(to)
    })
  }
}

registerChart("PieChart", PieChart)

// --- spec functions -------------------------------------------------------------------

// ({index, outerRadius = 0, ...props}) => index === i ? <Sector {...props} outerRadius={outerRadius + grow}/> : <Sector {...props} outerRadius={outerRadius}/>
registerFn("pieActiveSector", (args) => (props) => {
  const { index, outerRadius = 0, ...rest } = props
  return Sector({ ...rest, outerRadius: index === args[0] ? outerRadius + args[1] : outerRadius })
})

// chart-pie-interactive's renderPieShape.
registerFn("pieActiveSectorRing", (args) => (props) => {
  const { index, outerRadius = 0, ...rest } = props
  if (index === args[0]) return h("g", null, Sector({ ...rest, outerRadius: outerRadius + 10 }), Sector({ ...rest, outerRadius: outerRadius + 25, innerRadius: outerRadius + 12 }))
  return Sector({ ...rest, outerRadius })
})

// Label content: the value (large) and a caption at the view box center.
registerFn("pieCenterText", (args) => ({ viewBox }) => {
  if (!viewBox || !("cx" in viewBox) || !("cy" in viewBox)) return undefined
  const value = typeof args[0] === "number" ? args[0].toLocaleString() : String(args[0])
  return h(
    "text",
    { x: viewBox.cx, y: viewBox.cy, textAnchor: "middle", dominantBaseline: "middle" },
    h("tspan", { x: viewBox.cx, y: viewBox.cy, className: "fill-foreground text-3xl font-bold" }, value),
    h("tspan", { x: viewBox.cx, y: (viewBox.cy || 0) + 24, className: "fill-muted-foreground" }, args[1]),
  )
})

// Label content of lines at the view box center: args are [text, className,
// dy] (dy from cy), one <tspan> each, as the create designer's donut does.
registerFn("pieCenterLines", (args) => ({ viewBox }) => {
  if (!viewBox || !("cx" in viewBox) || !("cy" in viewBox)) return undefined
  return h(
    "text",
    { x: viewBox.cx, y: viewBox.cy, textAnchor: "middle", dominantBaseline: "middle" },
    ...args.map(([text, className, dy]) => h("tspan", { x: viewBox.cx, y: (viewBox.cy || 0) + dy, className }, text)),
  )
})

// Pie label: <text cx cy x y textAnchor dominantBaseline fill>{payload[field]}</text>.
registerFn("pieValueLabel", (args) => ({ payload, ...props }) =>
  h("text", { cx: props.cx, cy: props.cy, x: props.x, y: props.y, textAnchor: props.textAnchor, dominantBaseline: props.dominantBaseline, fill: args[1] }, String(payload?.[args[0]])),
)

// labelFormatter: (_, payload) => chartConfig[payload?.[0].dataKey].label
registerFn("piePayloadConfigLabel", (args, ctx) => (_, payload) => ctx.config[payload?.[0].dataKey]?.label)
