// chart.js: the shadcn/ui Chart runtime, a port of recharts 3.8.0's
// cartesian charts (AreaChart, BarChart, LineChart, ComposedChart) to
// vanilla JavaScript.
//
// The Go package renders the ChartContainer markup (data-slot="chart", its
// <style> and recharts' ResponsiveContainer) plus a JSON description of the
// recharts element tree in <script type="application/json"
// data-chart-spec>. This module measures the container, then renders the
// exact SVG recharts renders: same elements, attributes, z-index layers,
// text measurements, tick selection, scales, curves, animations (driven by
// requestAnimationFrame with recharts' own easing), tooltip (mouse, touch,
// click, bar item enter/leave, keyboard, Escape), cursor, active dots and
// legend. Like ResponsiveContainer, a container measured at a non-positive
// size unmounts the chart, which mounts afresh once it has a size. The
// spec <script> is removed once read (React renders none).
//
// The numeric parts of recharts that do not depend on React (d3 scales and
// shapes, nice ticks, getTicks, stacking, bar positions, tooltip
// translation, easing) are vendored unmodified in ./recharts.js; this file
// ports the React components and the redux selectors that glue them.
//
// Extension points for other chart types (pie, radar, radial bar):
//
//   import { registerChart, registerFn, h, ... } from "./chart.js"
//
//   registerChart(name, Instance)   // Instance: class with render(), see CartesianChart
//   registerFn(name, factory)       // factory(args, ctx) => function, for {"$fn": name}
//
// Each mounted container exposes el.chart (the chart instance) with
// setSpec(spec) and update(fn) for applications that change the chart
// (data, dataKey, ...) after load; changes animate like recharts.
//
// Events: "shadcn:chart-tooltip" (detail: {active, index, label}) on the
// container when the active tooltip index changes.
import { define, emit } from "./core.js"
import { getNonce } from "./style.js"
import * as R from "./recharts.js"

// --- element creation with React DOM semantics ---------------------------------

const SVG_NS = "http://www.w3.org/2000/svg"
const HTML_TAGS = new Set(["div", "span", "button", "p", "ul", "li"])

// Camel-case props React writes as hyphenated attributes (React DOM's
// attribute aliases).
const ALIASES = (() => {
  const m = {
    className: "class",
    htmlFor: "for",
    tabIndex: "tabindex",
    xlinkHref: "xlink:href",
    xmlnsXlink: "xmlns:xlink",
    xmlLang: "xml:lang",
    xmlSpace: "xml:space",
  }
  for (const k of "accentHeight alignmentBaseline arabicForm baselineShift capHeight clipPath clipRule colorInterpolation colorInterpolationFilters colorProfile colorRendering dominantBaseline enableBackground fillOpacity fillRule floodColor floodOpacity fontFamily fontSize fontSizeAdjust fontStretch fontStyle fontVariant fontWeight glyphName glyphOrientationHorizontal glyphOrientationVertical horizAdvX horizOriginX imageRendering letterSpacing lightingColor markerEnd markerMid markerStart overlinePosition overlineThickness paintOrder panose1 pointerEvents renderingIntent shapeRendering stopColor stopOpacity strikethroughPosition strikethroughThickness strokeDasharray strokeDashoffset strokeLinecap strokeLinejoin strokeMiterlimit strokeOpacity strokeWidth textAnchor textDecoration textRendering transformOrigin underlinePosition underlineThickness unicodeBidi unicodeRange unitsPerEm vAlphabetic vHanging vIdeographic vMathematical vectorEffect vertAdvY vertOriginX vertOriginY wordSpacing writingMode xHeight".split(" ")) {
    m[k] = k.replace(/[A-Z]/g, (c) => "-" + c.toLowerCase())
  }
  m.panose1 = "panose-1"
  return m
})()

// React's unitless CSS properties (numbers are not suffixed with px).
const UNITLESS = new Set("animationIterationCount aspectRatio borderImageOutset borderImageSlice borderImageWidth boxFlex boxFlexGroup boxOrdinalGroup columnCount columns flex flexGrow flexPositive flexShrink flexNegative flexOrder gridArea gridRow gridRowEnd gridRowSpan gridRowStart gridColumn gridColumnEnd gridColumnSpan gridColumnStart fontWeight lineClamp lineHeight opacity order orphans scale tabSize widows zIndex zoom fillOpacity floodOpacity stopOpacity strokeDasharray strokeDashoffset strokeMiterlimit strokeOpacity strokeWidth".split(" "))

// Sets an element's attributes and inline style to exactly the given ones
// (React semantics), touching the DOM only for actual changes.
export function syncElement(el, props, style) {
  const scratch = el.namespaceURI === SVG_NS ? document.createElementNS(SVG_NS, el.tagName) : document.createElement(el.tagName)
  setProps(scratch, props)
  for (const attr of [...el.attributes]) {
    if (attr.name !== "style" && !scratch.hasAttributeNS(attr.namespaceURI, attr.localName)) el.removeAttributeNS(attr.namespaceURI, attr.localName)
  }
  for (const attr of scratch.attributes) {
    if (el.getAttributeNS(attr.namespaceURI, attr.localName) !== attr.value) el.setAttributeNS(attr.namespaceURI, attr.name, attr.value)
  }
  // Inline style: like React, only properties whose value changed are
  // written, through the style object (attribute serialization would
  // round numbers).
  const prev = el.__chartStyle ?? {}
  const next = style ?? {}
  for (const k in prev) {
    if (!(k in next) || next[k] == null || next[k] === "" || typeof next[k] === "boolean") {
      if (k.startsWith("--")) el.style.removeProperty(k)
      else el.style[k] = ""
    }
  }
  for (const k in next) {
    if (prev[k] !== next[k]) applyStyle(el, { [k]: next[k] })
  }
  el.__chartStyle = { ...next }
  if (el.getAttribute("style") === "") el.removeAttribute("style")
}

function patchAttributes(a, b) {
  for (const attr of [...a.attributes]) {
    if (!b.hasAttributeNS(attr.namespaceURI, attr.localName)) a.removeAttributeNS(attr.namespaceURI, attr.localName)
  }
  for (const attr of b.attributes) {
    if (a.getAttributeNS(attr.namespaceURI, attr.localName) !== attr.value) a.setAttributeNS(attr.namespaceURI, attr.name, attr.value)
  }
}

function applyStyle(el, style) {
  if (!style) return
  for (const [k, v] of Object.entries(style)) {
    if (v == null || typeof v === "boolean" || v === "") continue
    if (k.startsWith("--")) {
      el.style.setProperty(k, String(v))
      continue
    }
    const value = typeof v === "number" && v !== 0 && !UNITLESS.has(k) ? v + "px" : String(v).trim()
    if (k === "float") el.style.cssFloat = value
    else el.style[k] = value
  }
}

/**
 * Creates an element like React DOM would for <tag {...props}>children</tag>.
 * SVG tags are created in the SVG namespace, the few HTML tags in HTML.
 */
export function h(tag, props, ...children) {
  const el = HTML_TAGS.has(tag) ? document.createElement(tag) : document.createElementNS(SVG_NS, tag)
  if (props) setProps(el, props)
  appendChildren(el, children)
  return el
}

function appendChildren(el, children) {
  for (const c of children) {
    if (c == null || c === false || c === true) continue
    if (Array.isArray(c)) appendChildren(el, c)
    else if (c instanceof Node) el.appendChild(c)
    else el.appendChild(document.createTextNode(String(c)))
  }
}

export function setProps(el, props) {
  for (const k in props) {
    if (!Object.prototype.hasOwnProperty.call(props, k)) continue
    const v = props[k]
    if (k === "children" || k === "key" || k === "ref" || k === "dangerouslySetInnerHTML") continue
    if (k === "style") {
      applyStyle(el, v)
      continue
    }
    if (v == null || typeof v === "function" || typeof v === "symbol") continue
    if (typeof v === "boolean" && !k.startsWith("data-") && !k.startsWith("aria-")) continue
    if (typeof v === "object" && !Array.isArray(v)) continue
    const name = ALIASES[k] ?? k
    if (name.startsWith("xlink:")) el.setAttributeNS("http://www.w3.org/1999/xlink", name, String(v))
    else el.setAttribute(name, String(v))
  }
}

// --- recharts helpers ported from React modules -------------------------------------

const SVG_KEYS = new Set(["aria-activedescendant", "aria-atomic", "aria-autocomplete", "aria-busy", "aria-checked", "aria-colcount", "aria-colindex", "aria-colspan", "aria-controls", "aria-current", "aria-describedby", "aria-details", "aria-disabled", "aria-errormessage", "aria-expanded", "aria-flowto", "aria-haspopup", "aria-hidden", "aria-invalid", "aria-keyshortcuts", "aria-label", "aria-labelledby", "aria-level", "aria-live", "aria-modal", "aria-multiline", "aria-multiselectable", "aria-orientation", "aria-owns", "aria-placeholder", "aria-posinset", "aria-pressed", "aria-readonly", "aria-relevant", "aria-required", "aria-roledescription", "aria-rowcount", "aria-rowindex", "aria-rowspan", "aria-selected", "aria-setsize", "aria-sort", "aria-valuemax", "aria-valuemin", "aria-valuenow", "aria-valuetext", "className", "color", "height", "id", "lang", "max", "media", "method", "min", "name", "style", "target", "width", "role", "tabIndex", "accentHeight", "accumulate", "additive", "alignmentBaseline", "allowReorder", "alphabetic", "amplitude", "arabicForm", "ascent", "attributeName", "attributeType", "autoReverse", "azimuth", "baseFrequency", "baselineShift", "baseProfile", "bbox", "begin", "bias", "by", "calcMode", "capHeight", "clip", "clipPath", "clipPathUnits", "clipRule", "colorInterpolation", "colorInterpolationFilters", "colorProfile", "colorRendering", "contentScriptType", "contentStyleType", "cursor", "cx", "cy", "d", "decelerate", "descent", "diffuseConstant", "direction", "display", "divisor", "dominantBaseline", "dur", "dx", "dy", "edgeMode", "elevation", "enableBackground", "end", "exponent", "externalResourcesRequired", "fill", "fillOpacity", "fillRule", "filter", "filterRes", "filterUnits", "floodColor", "floodOpacity", "focusable", "fontFamily", "fontSize", "fontSizeAdjust", "fontStretch", "fontStyle", "fontVariant", "fontWeight", "format", "from", "fx", "fy", "g1", "g2", "glyphName", "glyphOrientationHorizontal", "glyphOrientationVertical", "glyphRef", "gradientTransform", "gradientUnits", "hanging", "horizAdvX", "horizOriginX", "href", "ideographic", "imageRendering", "in2", "in", "intercept", "k1", "k2", "k3", "k4", "k", "kernelMatrix", "kernelUnitLength", "kerning", "keyPoints", "keySplines", "keyTimes", "lengthAdjust", "letterSpacing", "lightingColor", "limitingConeAngle", "local", "markerEnd", "markerHeight", "markerMid", "markerStart", "markerUnits", "markerWidth", "mask", "maskContentUnits", "maskUnits", "mathematical", "mode", "numOctaves", "offset", "opacity", "operator", "order", "orient", "orientation", "origin", "overflow", "overlinePosition", "overlineThickness", "paintOrder", "panose1", "pathLength", "patternContentUnits", "patternTransform", "patternUnits", "pointerEvents", "pointsAtX", "pointsAtY", "pointsAtZ", "preserveAlpha", "preserveAspectRatio", "primitiveUnits", "r", "radius", "refX", "refY", "renderingIntent", "repeatCount", "repeatDur", "requiredExtensions", "requiredFeatures", "restart", "result", "rotate", "rx", "ry", "seed", "shapeRendering", "slope", "spacing", "specularConstant", "specularExponent", "speed", "spreadMethod", "startOffset", "stdDeviation", "stemh", "stemv", "stitchTiles", "stopColor", "stopOpacity", "strikethroughPosition", "strikethroughThickness", "string", "stroke", "strokeDasharray", "strokeDashoffset", "strokeLinecap", "strokeLinejoin", "strokeMiterlimit", "strokeOpacity", "strokeWidth", "surfaceScale", "systemLanguage", "tableValues", "targetX", "targetY", "textAnchor", "textDecoration", "textLength", "textRendering", "to", "transform", "u1", "u2", "underlinePosition", "underlineThickness", "unicode", "unicodeBidi", "unicodeRange", "unitsPerEm", "vAlphabetic", "values", "vectorEffect", "version", "vertAdvY", "vertOriginX", "vertOriginY", "vHanging", "vIdeographic", "viewTarget", "visibility", "vMathematical", "widths", "wordSpacing", "writingMode", "x1", "x2", "x", "xChannelSelector", "xHeight", "xlinkActuate", "xlinkArcrole", "xlinkHref", "xlinkRole", "xlinkShow", "xlinkTitle", "xlinkType", "xmlBase", "xmlLang", "xmlns", "xmlnsXlink", "xmlSpace", "y1", "y2", "y", "yChannelSelector", "z", "zoomAndPan", "ref", "key", "angle"])

/** recharts' svgPropertiesNoEvents (events are never set here). */
export function svgProps(obj) {
  const out = {}
  if (typeof obj !== "object" || obj === null) return out
  for (const k in obj) {
    if (Object.prototype.hasOwnProperty.call(obj, k) && (SVG_KEYS.has(k) || k.startsWith("data-"))) out[k] = obj[k]
  }
  return out
}

/** svgPropertiesNoEventsFromUnknown: plain objects only (no React elements here). */
function svgPropsFromUnknown(input) {
  if (input == null) return null
  if (typeof input === "object" && !Array.isArray(input) && !input.$fn) return svgProps(input)
  return null
}

function omit(obj, keys) {
  const out = { ...obj }
  for (const k of keys) delete out[k]
  return out
}

/** resolveDefaultProps: defaults fill only undefined props. */
export function withDefaults(props, defaults) {
  const out = { ...props }
  for (const k in defaults) if (out[k] === undefined) out[k] = defaults[k]
  return out
}

export const DefaultZIndexes = { grid: -100, barBackground: -50, area: 100, cursorRectangle: 200, bar: 300, line: 400, axis: 500, scatter: 600, activeBar: 1000, cursorLine: 1100, activeDot: 1200, label: 2000 }

const isNum = R.isNumber
const isNullish = (v) => v == null
const wellBehaved = R.isWellBehavedNumber

// clsx for the few class lists recharts builds.
function clsx(...args) {
  const out = []
  for (const a of args) {
    if (!a) continue
    if (typeof a === "string") out.push(a)
    else if (typeof a === "object") for (const k in a) if (a[k]) out.push(k)
  }
  return out.join(" ")
}

// --- function props ({"$fn": name, ...}) -----------------------------------------------

const FNS = {
  // (value) => String(value).slice(a, b)
  slice: (args) => (v) => (v == null ? v : String(v).slice(...args)),
  // (value) => chartConfig[value]?.label
  configLabel: (args, ctx) => (v) => ctx.config[v]?.label,
  // (value) => new Date(value).toLocaleDateString(locale, options)
  date: (args) => (v) => new Date(v).toLocaleDateString(args[0], args[1]),
  // (value) => String(value) + suffix
  suffix: (args) => (v) => `${v}${args[0]}`,
}

/** Registers a named function usable as {"$fn": name, "args": [...]} in specs. */
export function registerFn(name, factory) {
  FNS[name] = factory
}

/** Turns {"$fn"} values into functions, recursively in plain objects. */
function resolveFn(v, ctx) {
  if (v && typeof v === "object" && !Array.isArray(v) && typeof v.$fn === "string") {
    const f = FNS[v.$fn]
    if (!f) throw new Error(`chart: unknown function ${v.$fn}`)
    const fn = f(v.args ?? [], ctx, v)
    fn.$spec = v
    return fn
  }
  return v
}

/** Reads {"$field": name} from an entry (payload), for shapes and dots. */
function fieldValue(v, entry) {
  if (v && typeof v === "object" && !Array.isArray(v) && typeof v.$field === "string") return entry?.[v.$field]
  return v
}

function resolveFields(props, entry) {
  const out = {}
  for (const k in props) out[k] = fieldValue(props[k], entry)
  return out
}

// --- in-place updates ---------------------------------------------------------------

/**
 * Updates a with b's attributes and children and returns the node now in
 * the document, like React reconciling an element of the same type: nodes
 * that did not change are left alone (so are their paint and transitions).
 */
function patch(a, b) {
  if (a.nodeType === 3 && b.nodeType === 3) {
    if (a.data !== b.data) a.data = b.data
    return a
  }
  if (a.nodeType !== 1 || b.nodeType !== 1 || a.namespaceURI !== b.namespaceURI || a.tagName !== b.tagName) {
    a.replaceWith(b)
    return b
  }
  patchAttributes(a, b)
  patchChildren(a, [...b.childNodes])
  return a
}

/** Makes parent's children match next, patching positionally. */
export function patchChildren(parent, next) {
  const old = [...parent.childNodes]
  const n = Math.min(old.length, next.length)
  for (let i = 0; i < n; i++) patch(old[i], next[i])
  for (let i = n; i < old.length; i++) old[i].remove()
  for (let i = n; i < next.length; i++) parent.appendChild(next[i])
}

// --- ZIndex layers and portal-like slots ------------------------------------------

/**
 * A slot is the content one component portals into a layer. Like React
 * portals, a slot that re-renders keeps its place; a slot that becomes
 * non-empty again is appended at the end of the layer.
 */
class Slots {
  constructor() {
    this.map = new Map()
  }
  set(container, key, nodes) {
    nodes = nodes.filter(Boolean)
    const prev = this.map.get(key)
    const old = prev && prev.container === container ? prev.nodes : []
    if (prev && prev.container !== container) for (const n of prev.nodes) n.remove()
    if (old.length && old[0].parentNode === container) {
      const k = Math.min(old.length, nodes.length)
      const result = []
      for (let i = 0; i < k; i++) result.push(patch(old[i], nodes[i]))
      const after = old[old.length - 1].nextSibling
      for (let i = k; i < old.length; i++) old[i].remove()
      for (let i = k; i < nodes.length; i++) {
        container.insertBefore(nodes[i], after)
        result.push(nodes[i])
      }
      nodes = result
    } else {
      for (const n of old) n.remove()
      for (const n of nodes) container.appendChild(n)
    }
    this.map.set(key, { container, nodes })
    return nodes
  }
  clear(key) {
    const prev = this.map.get(key)
    if (prev) for (const n of prev.nodes) n.remove()
    this.map.delete(key)
  }
  keys() {
    return [...this.map.keys()]
  }
}

// --- Text (recharts component/Text.js) ------------------------------------------------

const BREAKING_SPACES = /[ \f\n\r\t\v\u2028\u2029]+/

function calculateWordWidths({ children, breakAll, style }) {
  try {
    let words = []
    if (!isNullish(children)) words = breakAll ? children.toString().split("") : children.toString().split(BREAKING_SPACES)
    const wordsWithComputedWidth = words.map((word) => ({ word, width: R.getStringSize(word, style).width }))
    const spaceWidth = breakAll ? 0 : R.getStringSize(" ", style).width
    return { wordsWithComputedWidth, spaceWidth }
  } catch {
    return null
  }
}

function calculateLines(words, lineWidth, spaceWidth, scaleToFit) {
  return words.reduce((result, { word, width }) => {
    const currentLine = result[result.length - 1]
    if (currentLine && width != null && (lineWidth == null || scaleToFit || currentLine.width + width + spaceWidth < Number(lineWidth))) {
      currentLine.words.push(word)
      currentLine.width += width + spaceWidth
    } else {
      result.push({ words: [word], width })
    }
    return result
  }, [])
}

const findLongestLine = (words) => words.reduce((a, b) => (a.width > b.width ? a : b))

function checkOverflow(text, index, breakAll, style, maxLines, lineWidth, spaceWidth, scaleToFit) {
  const tempText = text.slice(0, index)
  const words = calculateWordWidths({ breakAll, style, children: tempText + "…" })
  if (!words) return [false, []]
  const result = calculateLines(words.wordsWithComputedWidth, lineWidth, spaceWidth, scaleToFit)
  const doesOverflow = result.length > maxLines || findLongestLine(result).width > Number(lineWidth)
  return [doesOverflow, result]
}

function calculateWordsByLines({ maxLines, children, style, breakAll }, initialWords, spaceWidth, lineWidth, scaleToFit) {
  const shouldLimitLines = isNum(maxLines)
  const text = String(children)
  const originalResult = calculateLines(initialWords, lineWidth, spaceWidth, scaleToFit)
  if (!shouldLimitLines || scaleToFit) return originalResult
  const overflows = originalResult.length > maxLines || findLongestLine(originalResult).width > Number(lineWidth)
  if (!overflows) return originalResult
  let start = 0
  let end = text.length - 1
  let iterations = 0
  let trimmedResult
  while (start <= end && iterations <= text.length - 1) {
    const middle = Math.floor((start + end) / 2)
    const prev = middle - 1
    const [doesPrevOverflow, result] = checkOverflow(text, prev, breakAll, style, maxLines, lineWidth, spaceWidth, scaleToFit)
    const [doesMiddleOverflow] = checkOverflow(text, middle, breakAll, style, maxLines, lineWidth, spaceWidth, scaleToFit)
    if (!doesPrevOverflow && !doesMiddleOverflow) start = middle + 1
    if (doesPrevOverflow && doesMiddleOverflow) end = middle - 1
    if (!doesPrevOverflow && doesMiddleOverflow) {
      trimmedResult = result
      break
    }
    iterations++
  }
  return trimmedResult || originalResult
}

function getWordsWithoutCalculate(children) {
  const words = !isNullish(children) ? children.toString().split(BREAKING_SPACES) : []
  return [{ words, width: undefined }]
}

function getWordsByLines({ width, scaleToFit, children, style, breakAll, maxLines }) {
  if (width || scaleToFit) {
    const ww = calculateWordWidths({ breakAll, children, style })
    if (!ww) return getWordsWithoutCalculate(children)
    return calculateWordsByLines({ breakAll, children, maxLines, style }, ww.wordsWithComputedWidth, ww.spaceWidth, width, Boolean(scaleToFit))
  }
  return getWordsWithoutCalculate(children)
}

const TEXT_DEFAULTS = { angle: 0, breakAll: false, capHeight: "0.71em", fill: "#808080", lineHeight: "1em", scaleToFit: false, textAnchor: "start", verticalAnchor: "end", x: 0, y: 0 }

/** recharts <Text>: returns a <text> element or null. */
export function Text(outsideProps) {
  const all = withDefaults(outsideProps, TEXT_DEFAULTS)
  const { x: propsX, y: propsY, lineHeight, capHeight, fill, scaleToFit, textAnchor, verticalAnchor } = all
  const props = omit(all, ["x", "y", "lineHeight", "capHeight", "fill", "scaleToFit", "textAnchor", "verticalAnchor"])
  const wordsByLines = getWordsByLines({ breakAll: props.breakAll, children: props.children, maxLines: props.maxLines, scaleToFit, style: props.style, width: props.width })
  const { dx, dy, angle, className, breakAll } = props
  const textProps = omit(props, ["dx", "dy", "angle", "className", "breakAll"])
  if (!R.isNumOrStr(propsX) || !R.isNumOrStr(propsY) || wordsByLines.length === 0) return null
  const x = Number(propsX) + (isNum(dx) ? dx : 0)
  const y = Number(propsY) + (isNum(dy) ? dy : 0)
  if (!wellBehaved(x) || !wellBehaved(y)) return null
  let startDy
  switch (verticalAnchor) {
    case "start":
      startDy = R.reduceCSSCalc(`calc(${capHeight})`)
      break
    case "middle":
      startDy = R.reduceCSSCalc(`calc(${(wordsByLines.length - 1) / 2} * -${lineHeight} + (${capHeight} / 2))`)
      break
    default:
      startDy = R.reduceCSSCalc(`calc(${wordsByLines.length - 1} * -${lineHeight})`)
      break
  }
  const transforms = []
  const firstLine = wordsByLines[0]
  if (scaleToFit && firstLine != null) {
    const lineWidth = firstLine.width
    transforms.push(`scale(${isNum(props.width) && isNum(lineWidth) ? props.width / lineWidth : 1})`)
  }
  if (angle) transforms.push(`rotate(${angle}, ${x}, ${y})`)
  if (transforms.length) textProps.transform = transforms.join(" ")
  const attrs = { ...svgProps(textProps), x, y, className: clsx("recharts-text", className), textAnchor, fill: String(fill).includes("url") ? "#808080" : fill }
  delete attrs.children
  return h(
    "text",
    attrs,
    wordsByLines.map((line, index) => h("tspan", { x, dy: index === 0 ? startDy : lineHeight }, line.words.join(breakAll ? "" : " "))),
  )
}

// --- shapes ------------------------------------------------------------------------

function roundT(strings, ...values) {
  return R.roundTemplateLiteral(strings, ...values)
}

/** recharts shape/Rectangle getRectanglePath. */
function getRectanglePath(x, y, width, height, radius) {
  const roundedWidth = R.round(width)
  const roundedHeight = R.round(height)
  const maxRadius = Math.min(Math.abs(roundedWidth) / 2, Math.abs(roundedHeight) / 2)
  const ySign = roundedHeight >= 0 ? 1 : -1
  const xSign = roundedWidth >= 0 ? 1 : -1
  const clockWise = (roundedHeight >= 0 && roundedWidth >= 0) || (roundedHeight < 0 && roundedWidth < 0) ? 1 : 0
  let path
  if (maxRadius > 0 && Array.isArray(radius)) {
    const newRadius = [0, 0, 0, 0]
    for (let i = 0; i < 4; i++) {
      const r = radius[i] ?? 0
      newRadius[i] = r > maxRadius ? maxRadius : r
    }
    path = roundT`M${x},${y + ySign * newRadius[0]}`
    if (newRadius[0] > 0) path += roundT`A ${newRadius[0]},${newRadius[0]},0,0,${clockWise},${x + xSign * newRadius[0]},${y}`
    path += roundT`L ${x + width - xSign * newRadius[1]},${y}`
    if (newRadius[1] > 0) path += roundT`A ${newRadius[1]},${newRadius[1]},0,0,${clockWise},
        ${x + width},${y + ySign * newRadius[1]}`
    path += roundT`L ${x + width},${y + height - ySign * newRadius[2]}`
    if (newRadius[2] > 0) path += roundT`A ${newRadius[2]},${newRadius[2]},0,0,${clockWise},
        ${x + width - xSign * newRadius[2]},${y + height}`
    path += roundT`L ${x + xSign * newRadius[3]},${y + height}`
    if (newRadius[3] > 0) path += roundT`A ${newRadius[3]},${newRadius[3]},0,0,${clockWise},
        ${x},${y + height - ySign * newRadius[3]}`
    path += "Z"
  } else if (maxRadius > 0 && radius === +radius && radius > 0) {
    const r = Math.min(maxRadius, radius)
    path = roundT`M ${x},${y + ySign * r}
            A ${r},${r},0,0,${clockWise},${x + xSign * r},${y}
            L ${x + width - xSign * r},${y}
            A ${r},${r},0,0,${clockWise},${x + width},${y + ySign * r}
            L ${x + width},${y + height - ySign * r}
            A ${r},${r},0,0,${clockWise},${x + width - xSign * r},${y + height}
            L ${x + xSign * r},${y + height}
            A ${r},${r},0,0,${clockWise},${x},${y + height - ySign * r} Z`
  } else {
    path = roundT`M ${x},${y} h ${width} v ${height} h ${-width} Z`
  }
  return path
}

/** recharts <Rectangle> (no update animation, the default). */
export function Rectangle(rectangleProps) {
  const props = withDefaults(rectangleProps, { x: 0, y: 0, width: 0, height: 0, radius: 0 })
  const { x, y, width, height, radius, className } = props
  if (x !== +x || y !== +y || width !== +width || height !== +height || width === 0 || height === 0) return null
  const other = omit(svgProps(props), ["radius"])
  return h("path", {
    ...other,
    x: R.round(x),
    y: R.round(y),
    width: R.round(width),
    height: R.round(height),
    radius: typeof radius === "number" ? radius : undefined,
    className: clsx("recharts-rectangle", className),
    d: getRectanglePath(x, y, width, height, radius),
  })
}

/** recharts <Dot>. */
export function Dot(props) {
  const { cx, cy, r, className } = props
  if (isNum(cx) && isNum(cy) && isNum(r)) return h("circle", { ...svgProps(props), className: clsx("recharts-dot", className), cx, cy, r })
  return null
}

const CURVES = {
  curveBasisClosed: R.curveBasisClosed,
  curveBasisOpen: R.curveBasisOpen,
  curveBasis: R.curveBasis,
  curveBumpX: R.curveBumpX,
  curveBumpY: R.curveBumpY,
  curveLinearClosed: R.curveLinearClosed,
  curveLinear: R.curveLinear,
  curveMonotoneX: R.curveMonotoneX,
  curveMonotoneY: R.curveMonotoneY,
  curveNatural: R.curveNatural,
  curveStep: R.curveStep,
  curveStepAfter: R.curveStepAfter,
  curveStepBefore: R.curveStepBefore,
}

function getCurveFactory(type, layout) {
  const name = `curve${R.upperFirst(type)}`
  if ((name === "curveMonotone" || name === "curveBump") && layout) {
    const f = CURVES[`${name}${layout === "vertical" ? "Y" : "X"}`]
    if (f) return f
  }
  return CURVES[name] || R.curveLinear
}

const curveDefined = (p) => wellBehaved(p.x) && wellBehaved(p.y)
const areaDefined = (d) => d.base != null && curveDefined(d.base) && curveDefined(d)

/** recharts shape/Curve getPath. */
export function getCurvePath({ type = "linear", points = [], baseLine, layout, connectNulls = false }) {
  const curveFactory = getCurveFactory(type, layout)
  const formatPoints = connectNulls ? points.filter(curveDefined) : points
  if (Array.isArray(baseLine)) {
    const areaPoints = points.map((entry, index) => ({ ...entry, base: baseLine[index] }))
    const fn = layout === "vertical" ? R.area().y((p) => p.y).x1((p) => p.x).x0((d) => d.base.x) : R.area().x((p) => p.x).y1((p) => p.y).y0((d) => d.base.y)
    const nullable = fn.defined(areaDefined).curve(curveFactory)
    return nullable(connectNulls ? areaPoints.filter(areaDefined) : areaPoints)
  }
  let fn
  if (layout === "vertical" && isNum(baseLine)) fn = R.area().y((p) => p.y).x1((p) => p.x).x0(baseLine)
  else if (isNum(baseLine)) fn = R.area().x((p) => p.x).y1((p) => p.y).y0(baseLine)
  else fn = R.line().x((p) => p.x).y((p) => p.y)
  return fn.defined(curveDefined).curve(curveFactory)(formatPoints)
}

/** recharts <Curve>. */
export function Curve(props) {
  const { className, points, path } = props
  if ((!points || !points.length) && !path) return null
  const realPath = points && points.length ? getCurvePath({ type: props.type, points: props.points, baseLine: props.baseLine, layout: props.layout, connectNulls: props.connectNulls }) : path
  return h("path", { ...svgProps(props), className: clsx("recharts-curve", className), d: realPath === null ? undefined : realPath })
}

/** recharts <Layer>. */
export function Layer(props, ...children) {
  const { className } = props
  return h("g", { ...svgProps(omit(props, ["className", "children"])), className: clsx("recharts-layer", className) }, ...children)
}

// --- animation (recharts JavascriptAnimate) ---------------------------------------------

const prefersReducedMotion = () => typeof matchMedia === "function" && matchMedia("(prefers-reduced-motion: reduce)").matches

/**
 * One running JavascriptAnimate: t goes from 0 to 1 over duration with
 * recharts' easing, on the frames recharts uses (the animation manager
 * spends three frames before the first tick).
 */
class Animation {
  constructor(chart, { duration, begin = 0, easing = "ease", onEnd }) {
    this.chart = chart
    this.duration = duration
    this.begin = begin
    this.easing = R.configEasing(easing)
    this.onEnd = onEnd
    this.t = 0
    this.done = false
    this.stopped = false
    this.frames = 0
    this.beginTime = null
    // timeline: [onAnimationStart, begin, onAnimationActive, duration, onAnimationEnd]
    this.queue = ["wait", this.begin > 0 ? { delay: this.begin } : "wait", "active"]
    // RequestAnimationFrameTimeoutController: every step (and every timing
    // frame) is a setTimeout that records performance.now() when it is
    // scheduled and fires on the first frame whose rAF timestamp is at
    // least that start plus the delay.
    this.gate(0)
    chart.scheduleFrame()
  }
  gate(delay) {
    this.gateStart = performance.now()
    this.gateDelay = delay
  }
  stop() {
    this.stopped = true
  }
  // Returns true while it still needs frames.
  tick(now) {
    if (this.stopped || this.done) return false
    if (this.easing == null) {
      this.t = 1
      this.finish()
      return false
    }
    if (now - this.gateStart < this.gateDelay) return true
    if (this.queue.length) {
      const step = this.queue[0]
      if (step === "wait" || step.delay != null) this.queue.shift()
      else if (step === "active") {
        this.queue.shift()
        this.activeAt = now
        this.queue.push("timing")
      } else if (step === "timing") {
        this.queue.shift()
        this.beginTime = now
        this.endTimer = now
        this.t = 0
      }
      this.gate(this.queue[0]?.delay ?? 0)
      return true
    }
    if (this.ending) {
      // The animation manager's end callback runs on a frame after the
      // last tick (its timer starts after the first tick's timestamp).
      this.finish()
      return false
    }
    const raw = (now - this.beginTime) / this.duration
    this.gate(0)
    if (raw < 1) {
      this.t = this.easing(raw)
      return true
    }
    this.t = this.easing(1)
    this.ending = true
    return true
  }
  finish() {
    this.done = true
    this.onEnd?.()
  }
}

// --- chart configuration ----------------------------------------------------------------

const CHART_DEFAULTS = { accessibilityLayer: true, barCategoryGap: "10%", barGap: 4, layout: "horizontal", responsive: false, reverseStackOrder: false, stackOffset: "none", syncMethod: "index" }
const DEFAULT_MARGIN = { top: 5, right: 5, bottom: 5, left: 5 }

const AXIS_DEFAULT = {
  axisLine: true,
  tickLine: true,
  tickSize: 6,
  tickMargin: 2,
  interval: "preserveEnd",
  minTickGap: 5,
  mirror: false,
  stroke: "#666",
  tick: true,
}

const X_AXIS_DEFAULTS = { allowDataOverflow: false, allowDecimals: true, allowDuplicatedCategory: true, angle: 0, axisLine: true, height: 30, hide: false, includeHidden: false, interval: "preserveEnd", label: false, minTickGap: 5, mirror: false, orientation: "bottom", padding: { left: 0, right: 0 }, reversed: false, scale: "auto", tick: true, tickCount: 5, tickLine: true, tickSize: 6, type: "category", niceTicks: "auto", xAxisId: 0 }
const Y_AXIS_DEFAULTS = { allowDataOverflow: false, allowDecimals: true, allowDuplicatedCategory: true, angle: 0, axisLine: true, hide: false, includeHidden: false, interval: "preserveEnd", label: false, minTickGap: 5, mirror: false, orientation: "left", padding: { top: 0, bottom: 0 }, reversed: false, scale: "auto", tick: true, tickCount: 5, tickLine: true, tickSize: 6, type: "number", niceTicks: "auto", width: 60, yAxisId: 0 }

const IMPLICIT_X = { allowDataOverflow: false, allowDecimals: true, allowDuplicatedCategory: true, angle: 0, dataKey: undefined, domain: undefined, height: 30, hide: true, id: 0, includeHidden: false, interval: "preserveEnd", minTickGap: 5, mirror: false, name: undefined, orientation: "bottom", padding: { left: 0, right: 0 }, reversed: false, scale: "auto", tick: true, tickCount: 5, tickFormatter: undefined, ticks: undefined, type: "category", unit: undefined, niceTicks: "auto", implicit: true }
const IMPLICIT_Y = { allowDataOverflow: false, allowDecimals: true, allowDuplicatedCategory: true, angle: 0, dataKey: undefined, domain: [0, "auto"], hide: true, id: 0, includeHidden: false, interval: "preserveEnd", minTickGap: 5, mirror: false, name: undefined, orientation: "left", padding: { top: 0, bottom: 0 }, reversed: false, scale: "auto", tick: true, tickCount: 5, tickFormatter: undefined, ticks: undefined, type: "number", unit: undefined, niceTicks: "auto", width: 60, implicit: true }

const AXIS_SETTING_KEYS = ["allowDataOverflow", "allowDecimals", "allowDuplicatedCategory", "angle", "dataKey", "domain", "height", "hide", "includeHidden", "interval", "minTickGap", "mirror", "name", "orientation", "padding", "reversed", "scale", "tick", "tickCount", "tickFormatter", "ticks", "type", "unit", "niceTicks", "width"]

const AREA_DEFAULTS = { activeDot: true, animationBegin: 0, animationDuration: 1500, animationEasing: "ease", connectNulls: false, dot: false, fill: "#3182bd", fillOpacity: 0.6, hide: false, isAnimationActive: "auto", legendType: "line", stroke: "#3182bd", strokeWidth: 1, type: "linear", label: false, xAxisId: 0, yAxisId: 0, zIndex: DefaultZIndexes.area }
const BAR_DEFAULTS = { activeBar: false, animationBegin: 0, animationDuration: 400, animationEasing: "ease", background: false, hide: false, isAnimationActive: "auto", label: false, legendType: "rect", minPointSize: 0, xAxisId: 0, yAxisId: 0, zIndex: DefaultZIndexes.bar }
const LINE_DEFAULTS = { activeDot: true, animateNewValues: true, animationBegin: 0, animationDuration: 1500, animationEasing: "ease", connectNulls: false, dot: true, fill: "#fff", hide: false, isAnimationActive: "auto", label: false, legendType: "line", stroke: "#3182bd", strokeWidth: 1, xAxisId: 0, yAxisId: 0, zIndex: DefaultZIndexes.line, type: "linear" }
const GRID_DEFAULTS = { horizontal: true, vertical: true, horizontalPoints: [], verticalPoints: [], stroke: "#ccc", fill: "none", verticalFill: [], horizontalFill: [], xAxisId: 0, yAxisId: 0, syncWithTicks: false, zIndex: DefaultZIndexes.grid }
const TOOLTIP_DEFAULTS = { allowEscapeViewBox: { x: false, y: false }, animationDuration: 400, animationEasing: "ease", axisId: 0, contentStyle: {}, cursor: true, filterNull: true, includeHidden: false, isAnimationActive: "auto", itemSorter: "name", itemStyle: {}, labelStyle: {}, offset: 10, reverseDirection: { x: false, y: false }, separator: " : ", trigger: "hover", useTranslate3d: false, wrapperStyle: {} }
const LEGEND_DEFAULTS = { align: "center", iconSize: 14, inactiveColor: "#ccc", itemSorter: "value", layout: "horizontal", verticalAlign: "bottom" }
const LABEL_DEFAULTS = { angle: 0, offset: 5, zIndex: DefaultZIndexes.label, position: "middle", textBreakAll: false }

const CHART_TYPES = {}

/** Registers a chart type (the name of the recharts root element). */
export function registerChart(name, Instance) {
  CHART_TYPES[name] = Instance
}

/** React useId's format: _r_<base 32 counter>_. */
function reactId(n) {
  return `_r_${n.toString(32)}_`
}

export function resolveSpecNode(node, ctx) {
  const props = {}
  for (const [k, v] of Object.entries(node.props ?? {})) props[k] = resolveFn(v, ctx)
  return { type: node.type, props, children: (node.children ?? []).map((c) => resolveSpecNode(c, ctx)) }
}

// --- cartesian selectors (state/selectors/axisSelectors.js and friends) -----------------------

const isCategoricalAxis = R.isCategoricalAxis

function axisTypeBasedOnLayout(layout, axisType, type) {
  if (type !== "auto") return type
  if (layout == null) return undefined
  return isCategoricalAxis(layout, axisType) ? "category" : "number"
}

function makeNumber(val) {
  if (R.isNumOrStr(val) || val instanceof Date) {
    const n = Number(val)
    if (wellBehaved(n)) return n
  }
  return undefined
}

function makeDomain(val) {
  if (Array.isArray(val)) {
    const attempt = [makeNumber(val[0]), makeNumber(val[1])]
    return R.isWellFormedNumberDomain(attempt) ? attempt : undefined
  }
  const n = makeNumber(val)
  return n == null ? undefined : [n, n]
}

const onlyAllowNumbers = (data) => data.map(makeNumber).filter((v) => v != null)

function getDomainDefinition(axisSettings) {
  if (axisSettings == null || !("domain" in axisSettings)) return [0, "auto"]
  if (axisSettings.domain != null) return axisSettings.domain
  if ("ticks" in axisSettings && axisSettings.ticks != null) {
    if (axisSettings.type === "number") {
      const all = onlyAllowNumbers(axisSettings.ticks)
      return [Math.min(...all), Math.max(...all)]
    }
    if (axisSettings.type === "category") return axisSettings.ticks.map(String)
  }
  return axisSettings?.domain ?? [0, "auto"]
}

function mergeDomains(...domains) {
  const all = domains.filter(Boolean)
  if (all.length === 0) return undefined
  const values = all.flat()
  return [Math.min(...values), Math.max(...values)]
}

function combineAppliedValues(data, axisSettings, items) {
  if (axisSettings?.dataKey != null) return data.map((item) => ({ value: R.getValueByDataKey(item, axisSettings.dataKey) }))
  if (items.length > 0) return items.map((item) => item.dataKey).flatMap((dataKey) => data.map((entry) => ({ value: R.getValueByDataKey(entry, dataKey) })))
  return data.map((entry) => ({ value: entry }))
}

function combineDomainOfAllAppliedNumericalValues(data, axisSettings, items) {
  let lowerEnd, upperEnd
  if (items.length > 0) {
    data.forEach((entry) => {
      items.forEach((item) => {
        const valueByDataKey = R.getValueByDataKey(entry, axisSettings.dataKey ?? item.dataKey)
        const d = makeDomain(valueByDataKey)
        if (d != null) {
          lowerEnd = lowerEnd == null ? d[0] : Math.min(lowerEnd, d[0])
          upperEnd = upperEnd == null ? d[1] : Math.max(upperEnd, d[1])
        }
      })
    })
  }
  if (axisSettings?.dataKey != null) {
    data.forEach((item) => {
      const d = makeDomain(R.getValueByDataKey(item, axisSettings.dataKey))
      if (d != null) {
        lowerEnd = lowerEnd == null ? d[0] : Math.min(lowerEnd, d[0])
        upperEnd = upperEnd == null ? d[1] : Math.max(upperEnd, d[1])
      }
    })
  }
  if (wellBehaved(lowerEnd) && wellBehaved(upperEnd)) return [lowerEnd, upperEnd]
  return undefined
}

function hasDuplicate(ary) {
  return R.hasDuplicate(ary)
}

const rangeOf = (n) => Array.from({ length: Math.max(0, n) }, (_, i) => i)

function computeDomainOfTypeCategory(allDataSquished, axisSettings, isCategorical) {
  const categoricalDomain = allDataSquished.map(({ value }) => (R.isNumOrStr(value) || value instanceof Date ? value : undefined)).filter((v) => v != null)
  if (isCategorical && (axisSettings.dataKey == null || (axisSettings.allowDuplicatedCategory && hasDuplicate(categoricalDomain)))) return rangeOf(allDataSquished.length)
  if (axisSettings.allowDuplicatedCategory) return categoricalDomain
  return Array.from(new Set(categoricalDomain))
}

function combineRealScaleType(axisConfig, hasBar, chartType) {
  if (axisConfig == null) return undefined
  const { scale, type } = axisConfig
  if (scale === "auto") {
    if (type === "category" && chartType && (chartType.indexOf("LineChart") >= 0 || chartType.indexOf("AreaChart") >= 0 || (chartType.indexOf("ComposedChart") >= 0 && !hasBar))) return "point"
    if (type === "category") return "band"
    return "linear"
  }
  if (typeof scale === "string") return ["band", "linear", "point"].includes(scale) ? scale : "point"
  return undefined
}

function combineNiceTicks(axisDomain, axisSettings, realScaleType) {
  const { niceTicks } = axisSettings
  if (niceTicks === "none") return undefined
  const domainDefinition = getDomainDefinition(axisSettings)
  const hasAuto = Array.isArray(domainDefinition) && (domainDefinition[0] === "auto" || domainDefinition[1] === "auto")
  if ((niceTicks === "snap125" || niceTicks === "adaptive") && axisSettings.tickCount && R.isWellFormedNumberDomain(axisDomain)) {
    if (hasAuto) return R.getNiceTickValues(axisDomain, axisSettings.tickCount, axisSettings.allowDecimals, niceTicks)
    if (axisSettings.type === "number") return R.getTickValuesFixedDomain(axisDomain, axisSettings.tickCount, axisSettings.allowDecimals, niceTicks)
  }
  if (niceTicks === "auto" && realScaleType === "linear" && axisSettings.tickCount) {
    if (hasAuto && R.isWellFormedNumberDomain(axisDomain)) return R.getNiceTickValues(axisDomain, axisSettings.tickCount, axisSettings.allowDecimals, "adaptive")
    if (axisSettings.type === "number" && R.isWellFormedNumberDomain(axisDomain)) return R.getTickValuesFixedDomain(axisDomain, axisSettings.tickCount, axisSettings.allowDecimals, "adaptive")
  }
  return undefined
}

function combineAxisDomainWithNiceTicks(axisSettings, domain, niceTicks, axisType) {
  if (axisType !== "angleAxis" && axisSettings?.type === "number" && R.isWellFormedNumberDomain(domain) && Array.isArray(niceTicks) && niceTicks.length > 0) {
    return [Math.min(domain[0], niceTicks[0] ?? 0), Math.max(domain[1], niceTicks[niceTicks.length - 1] ?? 0)]
  }
  return domain
}

function combineCheckedDomain(realScaleType, axisDomain) {
  if (axisDomain == null) return undefined
  if (realScaleType === "linear" && !R.isWellFormedNumberDomain(axisDomain)) {
    let min, max
    for (const value of axisDomain) {
      if (!wellBehaved(value)) continue
      if (min === undefined || value < min) min = value
      if (max === undefined || value > max) max = value
    }
    return min !== undefined && max !== undefined ? [min, max] : undefined
  }
  return axisDomain
}

const D3_SCALES = { band: R.scaleBand, linear: R.scaleLinear, point: R.scalePoint }

/** Registers extra d3 scale factories (name: "band", "linear", ...). */
export function registerScale(name, factory) {
  D3_SCALES[name] = factory
}

/** recharts rechartsScaleFactory. */
export function rechartsScale(d3Scale) {
  if (d3Scale == null) return undefined
  const ticksFn = d3Scale.ticks
  const bandwidthFn = d3Scale.bandwidth
  const d3Range = d3Scale.range()
  const range = [Math.min(...d3Range), Math.max(...d3Range)]
  return {
    domain: () => d3Scale.domain(),
    range: () => range,
    rangeMin: () => range[0],
    rangeMax: () => range[1],
    isInRange(value) {
      const [first, last] = range
      return first <= last ? value >= first && value <= last : value >= last && value <= first
    },
    bandwidth: bandwidthFn ? () => bandwidthFn.call(d3Scale) : undefined,
    ticks: ticksFn ? (count) => ticksFn.call(d3Scale, count) : undefined,
    map: (input, options) => {
      let base = d3Scale(input)
      if (base == null) return undefined
      if (d3Scale.bandwidth && options?.position) {
        const bw = d3Scale.bandwidth()
        if (options.position === "middle") base += bw / 2
        else if (options.position === "end") base += bw
      }
      return base
    },
  }
}

function makeScale(realScaleType, domain, range) {
  if (domain == null || range == null) return undefined
  const f = D3_SCALES[realScaleType]
  if (!f) return undefined
  return rechartsScale(f().domain(domain).range(range))
}

function combineCategoricalDomain(layout, appliedValues, axis, axisType) {
  if (axis == null || axis.dataKey == null) return undefined
  const { type, scale } = axis
  if (isCategoricalAxis(layout, axisType) && (type === "number" || scale !== "auto")) return appliedValues.map((d) => d.value)
  return undefined
}

function combineDuplicateDomain(layout, appliedValues, axis, axisType) {
  if (axis == null) return undefined
  const { allowDuplicatedCategory, type, dataKey } = axis
  const all = appliedValues.map((av) => av.value)
  if (dataKey && isCategoricalAxis(layout, axisType) && type === "category" && allowDuplicatedCategory && hasDuplicate(all)) return all
  return undefined
}

function tickEntries(list, scale, offset, valueOf) {
  return list
    .map((entry, index) => {
      const scaled = scale.map(valueOf ? valueOf(entry) : entry)
      if (!wellBehaved(scaled)) return null
      return { coordinate: scaled + offset, value: entry, index, offset }
    })
    .filter(Boolean)
}

function domainTicks(scale, offset, duplicateDomain) {
  return scale
    .domain()
    .map((entry, index) => {
      const scaled = scale.map(entry)
      if (!wellBehaved(scaled)) return null
      return { coordinate: scaled + offset, value: duplicateDomain ? duplicateDomain[entry] : entry, index, offset }
    })
    .filter(Boolean)
}

/** combineAxisTicks (rendered axis ticks). */
function combineAxisTicks(layout, axis, realScaleType, scale, niceTicks, axisRange, duplicateDomain, categoricalDomain, axisType) {
  if (axis == null || scale == null) return undefined
  const isCategorical = isCategoricalAxis(layout, axisType)
  const { type, ticks, tickCount } = axis
  const offsetForBand = realScaleType === "scaleBand" && typeof scale.bandwidth === "function" ? scale.bandwidth() / 2 : 2
  const offset = type === "category" && scale.bandwidth ? scale.bandwidth() / offsetForBand : 0
  const t = ticks || niceTicks
  if (t) return tickEntries(t, scale, offset, (e) => (duplicateDomain ? duplicateDomain.indexOf(e) : e))
  if (isCategorical && categoricalDomain) return tickEntries(categoricalDomain, scale, offset)
  if (scale.ticks) return tickEntries(scale.ticks(tickCount), scale, offset)
  return domainTicks(scale, offset, duplicateDomain)
}

/** combineGraphicalItemTicks. */
function combineGraphicalItemTicks(layout, axis, scale, axisRange, duplicateDomain, categoricalDomain, axisType) {
  if (axis == null || scale == null || axisRange == null || axisRange[0] === axisRange[1]) return undefined
  const isCategorical = isCategoricalAxis(layout, axisType)
  if (isCategorical && categoricalDomain) return tickEntries(categoricalDomain, scale, 0)
  if (scale.ticks) return tickEntries(scale.ticks(axis.tickCount), scale, 0)
  return domainTicks(scale, 0, duplicateDomain)
}

/** combineTicksOfTooltipAxis. */
function combineTooltipTicks(layout, axis, realScaleType, scale, range, duplicateDomain, categoricalDomain, axisType) {
  if (!axis || !scale) return undefined
  const isCategorical = isCategoricalAxis(layout, axisType)
  const offsetForBand = realScaleType === "scaleBand" && scale.bandwidth ? scale.bandwidth() / 2 : 2
  const offset = axis.type === "category" && scale.bandwidth ? scale.bandwidth() / offsetForBand : 0
  if (isCategorical && categoricalDomain) return tickEntries(categoricalDomain, scale, offset)
  return domainTicks(scale, offset, duplicateDomain)
}

function combineDisplayedStackedData(stackedItems, chartData, tooltipAxisSettings) {
  const { allowDuplicatedCategory, dataKey: tooltipDataKey } = tooltipAxisSettings
  const known = new Map()
  stackedItems.forEach((item) => {
    const resolvedData = item.data ?? chartData
    if (resolvedData == null || resolvedData.length === 0) return
    resolvedData.forEach((entry, index) => {
      const tooltipValue = tooltipDataKey == null || allowDuplicatedCategory ? index : String(R.getValueByDataKey(entry, tooltipDataKey, null))
      const numericValue = R.getValueByDataKey(entry, item.dataKey, 0)
      const curr = known.has(tooltipValue) ? known.get(tooltipValue) : {}
      Object.assign(curr, { [item.id]: numericValue })
      known.set(tooltipValue, curr)
    })
  })
  return Array.from(known.values())
}

function combineStackGroups(displayedData, items, stackOffsetType, reverseStackOrder) {
  const groups = {}
  for (const item of items) {
    if (item.stackId == null) continue
    ;(groups[item.stackId] ??= []).push(item)
  }
  return Object.fromEntries(
    Object.entries(groups).map(([stackId, graphicalItems]) => {
      const ordered = reverseStackOrder ? [...graphicalItems].reverse() : graphicalItems
      return [stackId, { stackedData: R.getStackedData(displayedData, ordered.map((i) => i.id), stackOffsetType), graphicalItems: ordered }]
    }),
  )
}

const isStacked = (item) => item.stackId != null && item.dataKey != null

/**
 * Everything recharts derives from the state for a cartesian chart, in one
 * pass: offsets, axes (domain, scale, ticks), stacks and bar positions.
 */
class CartesianState {
  constructor(chart) {
    this.chart = chart
    // Before the first commit, axes and items have not registered with the
    // store yet: everything uses the implicit axes and no items.
    this.registered = chart.registered !== false
    this.xAxes = this.registered ? chart.xAxes : []
    this.yAxes = this.registered ? chart.yAxes : []
    this.items = this.registered ? chart.items : []
    this.layout = chart.props.layout
    this.width = chart.width
    this.height = chart.height
    this.data = chart.data
    this.axes = { xAxis: new Map(), yAxis: new Map() }
    this.computeOffset()
  }

  axisSettings(axisType, id) {
    const list = axisType === "xAxis" ? this.xAxes : this.yAxes
    const found = list.find((a) => a.settings.id === id)
    if (found) return found.settings
    return axisType === "xAxis" ? IMPLICIT_X : IMPLICIT_Y
  }

  computeOffset() {
    const c = this.chart
    const margin = c.margin
    const yAxes = this.yAxes.map((a) => a.settings)
    const xAxes = this.xAxes.map((a) => a.settings)
    const yWidth = (e) => (typeof e.width === "number" ? e.width : 60)
    const left = yAxes.reduce((r, e) => (e.orientation === "left" && !e.mirror && !e.hide ? r + yWidth(e) : r), 0)
    const right = yAxes.reduce((r, e) => (e.orientation === "right" && !e.mirror && !e.hide ? r + yWidth(e) : r), 0)
    const top = xAxes.reduce((r, e) => (e.orientation === "top" && !e.mirror && !e.hide ? r + e.height : r), 0)
    const bottom = xAxes.reduce((r, e) => (e.orientation === "bottom" && !e.mirror && !e.hide ? r + e.height : r), 0)
    let offset = { top: (margin.top || 0) + top, bottom: (margin.bottom || 0) + bottom, left: (margin.left || 0) + left, right: (margin.right || 0) + right }
    const brushBottom = offset.bottom
    if (this.registered) offset = R.appendOffsetOfLegend(offset, c.legendSettings, c.legendSize)
    this.offset = { brushBottom, ...offset, width: Math.max(this.width - offset.left - offset.right, 0), height: Math.max(this.height - offset.top - offset.bottom, 0) }
  }

  // Graphical items attached to an axis.
  itemsOf(axisType, id, settings) {
    return this.items.filter((item) => item[axisType === "xAxis" ? "xAxisId" : "yAxisId"] === id).filter((item) => settings?.includeHidden === true || !item.hide)
  }

  stackGroups(axisType, id) {
    const key = `${axisType}:${id}`
    this._stacks ??= new Map()
    if (this._stacks.has(key)) return this._stacks.get(key)
    const settings = this.axisSettings(axisType, id)
    const items = this.itemsOf(axisType, id, settings).filter((i) => i.type === "area" || i.type === "bar").filter(isStacked)
    const displayed = combineDisplayedStackedData(items, this.data, this.tooltipAxisSettings())
    const groups = combineStackGroups(displayed, items, this.chart.props.stackOffset, this.chart.props.reverseStackOrder)
    this._stacks.set(key, groups)
    return groups
  }

  tooltipAxisType() {
    return this.layout === "horizontal" ? "xAxis" : this.layout === "vertical" ? "yAxis" : this.layout === "centric" ? "angleAxis" : "radiusAxis"
  }

  tooltipAxisSettings() {
    return this.axisSettings(this.tooltipAxisType(), this.chart.tooltipProps?.axisId ?? 0)
  }

  hasBar() {
    return this.items.some((i) => i.type === "bar")
  }

  axis(axisType, id) {
    const cache = this.axes[axisType]
    if (cache.has(id)) return cache.get(id)
    const a = this.computeAxis(axisType, id)
    cache.set(id, a)
    return a
  }

  computeAxis(axisType, id) {
    const settings = this.axisSettings(axisType, id)
    const layout = this.layout
    const items = this.itemsOf(axisType, id, settings)
    const displayedData = this.data
    const appliedValues = combineAppliedValues(displayedData, settings, items)
    const isCategorical = isCategoricalAxis(layout, axisType)
    // numerical domain
    const domainDefinition = getDomainDefinition(settings)
    const domainFromUserPreference = R.numericalDomainSpecifiedWithoutRequiringData(domainDefinition, settings.allowDataOverflow)
    let domainOfStackGroups
    if (domainFromUserPreference == null) {
      const d = R.getDomainOfStackGroups(this.stackGroups(axisType, id), 0, Math.max(0, this.data.length - 1))
      domainOfStackGroups = d != null && d[0] === 0 && d[1] === 0 ? undefined : d
    }
    const notStacked = items.filter((item) => item.stackId === undefined)
    const dataDomain = combineDomainOfAllAppliedNumericalValues(displayedData, settings, notStacked)
    let numericalDomain
    if (domainFromUserPreference != null) numericalDomain = domainFromUserPreference
    else {
      const includeStacks = (layout === "vertical" && axisType === "xAxis") || (layout === "horizontal" && axisType === "yAxis")
      const merged = includeStacks ? mergeDomains(domainOfStackGroups, undefined, dataDomain) : mergeDomains(undefined, dataDomain)
      numericalDomain = R.parseNumericalUserDomain(domainDefinition, merged, settings.allowDataOverflow)
    }
    let axisDomain
    if ((displayedData == null || displayedData.length === 0) && numericalDomain === undefined) axisDomain = undefined
    else if (isCategorical && settings.dataKey == null) axisDomain = rangeOf(displayedData?.length ?? 0)
    else if (settings.type === "category") axisDomain = computeDomainOfTypeCategory(appliedValues, settings, isCategorical)
    else if (this.chart.props.stackOffset === "expand") axisDomain = [0, 1]
    else axisDomain = numericalDomain
    const realScaleType = combineRealScaleType(settings, this.hasBar(), this.chart.name)
    const niceTicks = axisDomain == null ? undefined : combineNiceTicks(axisDomain, settings, realScaleType)
    const domainWithNice = combineAxisDomainWithNiceTicks(settings, axisDomain, niceTicks, axisType)
    const range = this.axisRange(axisType, settings)
    const rangeWithReverse = settings.reversed ? [range[1], range[0]] : range
    const checked = combineCheckedDomain(realScaleType, domainWithNice)
    const scale = makeScale(realScaleType, checked, rangeWithReverse)
    const categoricalDomain = combineCategoricalDomain(layout, appliedValues, settings, axisType)
    const duplicateDomain = combineDuplicateDomain(layout, appliedValues, settings, axisType)
    const ticks = combineAxisTicks(layout, settings, realScaleType, scale, niceTicks, range, duplicateDomain, categoricalDomain, axisType)
    const itemTicks = combineGraphicalItemTicks(layout, settings, scale, range, duplicateDomain, categoricalDomain, axisType)
    const tooltipTicks = combineTooltipTicks(layout, settings, realScaleType, scale, range, duplicateDomain, categoricalDomain, axisType)
    return { axisType, id, settings, realScaleType, niceTicks, domain: checked, range, scale, ticks, itemTicks, tooltipTicks, categoricalDomain, duplicateDomain, isCategorical, axisWithScale: scale ? { ...settings, scale } : undefined }
  }

  axisRange(axisType, settings) {
    const o = this.offset
    if (axisType === "xAxis") {
      const p = typeof settings.padding === "string" ? { left: 0, right: 0 } : { left: settings.padding?.left ?? 0, right: settings.padding?.right ?? 0 }
      return [o.left + p.left, o.left + o.width - p.right]
    }
    const p = typeof settings.padding === "string" ? { top: 0, bottom: 0 } : { top: settings.padding?.top ?? 0, bottom: settings.padding?.bottom ?? 0 }
    if (this.layout === "horizontal") return [o.top + o.height - p.bottom, o.top + p.top]
    return [o.top + p.top, o.top + o.height - p.bottom]
  }

  // Axis position (selectXAxisPosition / selectYAxisPosition).
  axisPosition(axisType, settings) {
    const o = this.offset
    const list = (axisType === "xAxis" ? this.xAxes : this.yAxes).map((a) => a.settings).filter((a) => a.orientation === settings.orientation && a.mirror === settings.mirror).sort((a, b) => (a.id < b.id ? -1 : a.id > b.id ? 1 : 0))
    let position
    let step
    for (const axis of list) {
      const size = axisType === "xAxis" ? axis.height : typeof axis.width === "number" ? axis.width : 60
      if (position == null) {
        if (axisType === "xAxis") position = settings.orientation === "top" ? o.top : settings.orientation === "bottom" ? this.height - o.bottom : 0
        else position = settings.orientation === "left" ? o.left : settings.orientation === "right" ? this.width - o.right : 0
      }
      const needSpace = axisType === "xAxis" ? (settings.orientation === "top" && !settings.mirror) || (settings.orientation === "bottom" && settings.mirror) : (settings.orientation === "left" && !settings.mirror) || (settings.orientation === "right" && settings.mirror)
      if (axis.id === settings.id) step = position - Number(needSpace) * size
      position += (needSpace ? -1 : 1) * size
    }
    if (axisType === "xAxis") return { x: o.left, y: step ?? 0 }
    return { x: step ?? 0, y: o.top }
  }

  axisSize(axisType, settings) {
    if (axisType === "xAxis") return { width: this.offset.width, height: settings.height }
    return { width: typeof settings.width === "number" ? settings.width : 60, height: this.offset.height }
  }

  // Axis props the grid ticks generator needs (selectAxisPropsNeededForCartesianGridTicksGenerator).
  gridAxisProps(axisType, id) {
    const a = this.axis(axisType, id)
    const s = a.settings
    return {
      angle: s.angle,
      interval: s.interval,
      minTickGap: s.minTickGap,
      orientation: s.orientation,
      tick: s.tick,
      tickCount: s.tickCount,
      tickFormatter: s.tickFormatter,
      ticks: s.ticks,
      type: s.type,
      unit: s.unit,
      axisType,
      categoricalDomain: a.categoricalDomain,
      duplicateDomain: a.duplicateDomain,
      isCategorical: a.isCategorical,
      niceTicks: a.niceTicks,
      range: a.range,
      realScaleType: a.realScaleType,
      scale: a.scale,
    }
  }

  categoricalAxisType() {
    return isCategoricalAxis(this.layout, "xAxis") ? "xAxis" : "yAxis"
  }

  // Line/Area band size (selectBandSize in lineSelectors/areaSelectors).
  itemBandSize(item) {
    const x = this.axis("xAxis", item.xAxisId)
    const y = this.axis("yAxis", item.yAxisId)
    if (isCategoricalAxis(this.layout, "xAxis")) return R.getBandSizeOfAxis(x.axisWithScale, x.itemTicks, false)
    return R.getBandSizeOfAxis(y.axisWithScale, y.itemTicks, false)
  }

  stackedDataOf(item) {
    const numericAxisType = isCategoricalAxis(this.layout, "xAxis") ? "yAxis" : "xAxis"
    const groups = this.stackGroups(numericAxisType, numericAxisType === "yAxis" ? item.yAxisId : item.xAxisId)
    if (item.stackId == null || groups == null) return undefined
    return groups[item.stackId]?.stackedData?.find((v) => v.key === item.id)
  }

  computeArea(item) {
    const xAxis = this.axis("xAxis", item.xAxisId)
    const yAxis = this.axis("yAxis", item.yAxisId)
    if (!xAxis.axisWithScale || !yAxis.axisWithScale || !xAxis.itemTicks?.length || !yAxis.itemTicks?.length) return undefined
    const found = this.stackedDataOf(item)
    const stackedData = found?.map((p) => [p[0], p[1]])
    const bandSize = this.itemBandSize(item)
    return computeArea({ areaSettings: item, stackedData, layout: this.layout, chartBaseValue: this.chart.props.baseValue, xAxis: xAxis.axisWithScale, yAxis: yAxis.axisWithScale, displayedData: this.data, dataStartIndex: 0, xAxisTicks: xAxis.itemTicks, yAxisTicks: yAxis.itemTicks, bandSize })
  }

  computeLine(item) {
    const xAxis = this.axis("xAxis", item.xAxisId)
    const yAxis = this.axis("yAxis", item.yAxisId)
    if (!xAxis.axisWithScale || !yAxis.axisWithScale || !xAxis.itemTicks?.length || !yAxis.itemTicks?.length) return undefined
    const bandSize = this.itemBandSize(item)
    return computeLinePoints({ layout: this.layout, xAxis: xAxis.axisWithScale, yAxis: yAxis.axisWithScale, xAxisTicks: xAxis.itemTicks, yAxisTicks: yAxis.itemTicks, dataKey: item.dataKey, bandSize, displayedData: this.data })
  }

  computeBar(item, cells) {
    const layout = this.layout
    const xAxis = this.axis("xAxis", item.xAxisId)
    const yAxis = this.axis("yAxis", item.yAxisId)
    if (!xAxis.axisWithScale || !yAxis.axisWithScale || xAxis.itemTicks == null || yAxis.itemTicks == null) return undefined
    const catAxis = layout === "horizontal" ? xAxis : yAxis
    const bandSize = R.getBandSizeOfAxis(catAxis.axisWithScale, catAxis.itemTicks)
    if (bandSize == null) return undefined
    const pos = this.barPosition(item)
    if (pos == null) return undefined
    const numericAxisType = layout === "horizontal" ? "yAxis" : "xAxis"
    const groups = this.stackGroups(numericAxisType, numericAxisType === "yAxis" ? item.yAxisId : item.xAxisId)
    const stackedData = item.stackId != null ? groups?.[item.stackId]?.stackedData?.find((sd) => sd.key === item.id) : undefined
    return computeBarRectangles({ layout, barSettings: item, pos, bandSize, xAxis: xAxis.axisWithScale, yAxis: yAxis.axisWithScale, xAxisTicks: xAxis.itemTicks, yAxisTicks: yAxis.itemTicks, stackedData, displayedData: this.data, offset: this.offset, cells, parentViewBox: { x: 0, y: 0, width: this.width, height: this.height }, dataStartIndex: 0 })
  }

  barPosition(item) {
    const layout = this.layout
    const allBars = this.items.filter((i) => (layout === "horizontal" ? i.xAxisId === item.xAxisId : i.yAxisId === item.yAxisId)).filter((i) => i.hide === false).filter((i) => i.type === "bar")
    const catAxis = layout === "horizontal" ? this.axis("xAxis", item.xAxisId) : this.axis("yAxis", item.yAxisId)
    const axisSize = layout === "horizontal" ? this.offset.width : this.offset.height
    const sizeList = R.combineBarSizeList(allBars, this.chart.props.barSize, axisSize)
    const maxBarSize = item.maxBarSize ?? this.chart.props.maxBarSize
    const barBandSize = R.getBandSizeOfAxis(catAxis.axisWithScale, catAxis.itemTicks, true) ?? maxBarSize ?? 0
    const bandSize = R.getBandSizeOfAxis(catAxis.axisWithScale, catAxis.itemTicks)
    const all = R.combineAllBarPositions(sizeList, this.chart.props.maxBarSize, this.chart.props.barGap, this.chart.props.barCategoryGap, barBandSize, bandSize, item.maxBarSize)
    const p = all?.find((x) => x.stackId === item.stackId && item.dataKey != null && x.dataKeys.includes(item.dataKey))
    return p?.position
  }

  plotArea() {
    const o = this.offset
    return { x: o.left, y: o.top, width: o.width, height: o.height }
  }
}

// Area.computeArea / getBaseValue.
function getBaseValue(layout, chartBaseValue, itemBaseValue, xAxis, yAxis) {
  const baseValue = itemBaseValue ?? chartBaseValue
  if (isNum(baseValue)) return baseValue
  const numericAxis = layout === "horizontal" ? yAxis : xAxis
  const domain = numericAxis.scale.domain()
  if (numericAxis.type === "number") {
    const domainMax = Math.max(domain[0], domain[1])
    const domainMin = Math.min(domain[0], domain[1])
    if (baseValue === "dataMin") return domainMin
    if (baseValue === "dataMax") return domainMax
    return domainMax < 0 ? domainMax : Math.max(Math.min(domain[0], domain[1]), 0)
  }
  if (baseValue === "dataMin") return domain[0]
  if (baseValue === "dataMax") return domain[1]
  return domain[0]
}

function computeArea({ areaSettings: { connectNulls, baseValue: itemBaseValue, dataKey }, stackedData, layout, chartBaseValue, xAxis, yAxis, displayedData, dataStartIndex, xAxisTicks, yAxisTicks, bandSize }) {
  const hasStack = stackedData && stackedData.length
  const baseValue = getBaseValue(layout, chartBaseValue, itemBaseValue, xAxis, yAxis)
  const horizontal = layout === "horizontal"
  let isRange = false
  const points = displayedData.map((entry, index) => {
    let valueAsArray
    if (hasStack) valueAsArray = stackedData[dataStartIndex + index]
    else {
      const rawValue = R.getValueByDataKey(entry, dataKey)
      if (!Array.isArray(rawValue)) valueAsArray = [baseValue, rawValue]
      else {
        valueAsArray = rawValue
        isRange = true
      }
    }
    const value1 = valueAsArray?.[1] ?? null
    const isBreakPoint = value1 == null || (hasStack && !connectNulls && R.getValueByDataKey(entry, dataKey) == null)
    if (horizontal) {
      return { x: R.getCateCoordinateOfLine({ axis: xAxis, ticks: xAxisTicks, bandSize, entry, index }), y: isBreakPoint ? null : (yAxis.scale.map(value1) ?? null), value: valueAsArray, payload: entry }
    }
    return { x: isBreakPoint ? null : (xAxis.scale.map(value1) ?? null), y: R.getCateCoordinateOfLine({ axis: yAxis, ticks: yAxisTicks, bandSize, entry, index }), value: valueAsArray, payload: entry }
  })
  let baseLine
  if (hasStack || isRange) {
    baseLine = points.map((entry) => {
      const x = Array.isArray(entry.value) ? entry.value[0] : null
      if (horizontal) return { x: entry.x, y: x != null && entry.y != null ? (yAxis.scale.map(x) ?? null) : null, payload: entry.payload }
      return { x: x != null ? (xAxis.scale.map(x) ?? null) : null, y: entry.y, payload: entry.payload }
    })
  } else {
    baseLine = horizontal ? yAxis.scale.map(baseValue) : xAxis.scale.map(baseValue)
  }
  return { points, baseLine: baseLine ?? 0, isRange }
}

function computeLinePoints({ layout, xAxis, yAxis, xAxisTicks, yAxisTicks, dataKey, bandSize, displayedData }) {
  return displayedData
    .map((entry, index) => {
      const value = R.getValueByDataKey(entry, dataKey)
      if (layout === "horizontal") {
        const x = R.getCateCoordinateOfLine({ axis: xAxis, ticks: xAxisTicks, bandSize, entry, index })
        const y = isNullish(value) ? null : yAxis.scale.map(value)
        return { x, y: y ?? null, value, payload: entry }
      }
      const x = isNullish(value) ? null : xAxis.scale.map(value)
      const y = R.getCateCoordinateOfLine({ axis: yAxis, ticks: yAxisTicks, bandSize, entry, index })
      if (x == null || y == null) return null
      return { x, y, value, payload: entry }
    })
    .filter(Boolean)
}

function computeBarRectangles({ layout, barSettings: { dataKey, minPointSize: minPointSizeProp, hasCustomShape }, pos, bandSize, xAxis, yAxis, xAxisTicks, yAxisTicks, stackedData, displayedData, offset, cells, parentViewBox, dataStartIndex }) {
  const numericAxis = layout === "horizontal" ? yAxis : xAxis
  const stackedDomain = stackedData ? numericAxis.scale.domain() : null
  const baseValue = R.getBaseValueOfBar({ numericAxis })
  const stackedBarStart = numericAxis.scale.map(baseValue)
  return displayedData
    .map((entry, index) => {
      let value, x, y, width, height, background
      if (stackedData) {
        const untruncated = stackedData[index + dataStartIndex]
        if (untruncated == null) return null
        value = R.truncateByDomain(untruncated, stackedDomain)
      } else {
        value = R.getValueByDataKey(entry, dataKey)
        if (!Array.isArray(value)) value = [baseValue, value]
      }
      const minPointSize = isNum(minPointSizeProp) ? minPointSizeProp : 0
      if (layout === "horizontal") {
        const baseValueScale = yAxis.scale.map(value[0])
        const currentValueScale = yAxis.scale.map(value[1])
        if (baseValueScale == null || currentValueScale == null) return null
        x = R.getCateCoordinateOfBar({ axis: xAxis, ticks: xAxisTicks, bandSize, offset: pos.offset, entry, index })
        y = currentValueScale ?? baseValueScale ?? undefined
        width = pos.size
        const computedHeight = baseValueScale - currentValueScale
        height = R.isNan(computedHeight) ? 0 : computedHeight
        background = { x, y: offset.top, width, height: offset.height }
        if (Math.abs(minPointSize) > 0 && Math.abs(height) < Math.abs(minPointSize)) {
          const delta = R.mathSign(height || minPointSize) * (Math.abs(minPointSize) - Math.abs(height))
          y -= delta
          height += delta
        }
      } else {
        const baseValueScale = xAxis.scale.map(value[0])
        const currentValueScale = xAxis.scale.map(value[1])
        if (baseValueScale == null || currentValueScale == null) return null
        x = baseValueScale
        y = R.getCateCoordinateOfBar({ axis: yAxis, ticks: yAxisTicks, bandSize, offset: pos.offset, entry, index })
        width = currentValueScale - baseValueScale
        height = pos.size
        background = { x: offset.left, y, width: offset.width, height }
        if (Math.abs(minPointSize) > 0 && Math.abs(width) < Math.abs(minPointSize)) {
          width += R.mathSign(width || minPointSize) * (Math.abs(minPointSize) - Math.abs(width))
        }
      }
      if (x == null || y == null || width == null || height == null || (!hasCustomShape && (width === 0 || height === 0))) return null
      return { ...entry, stackedBarStart, x, y, width, height, value: stackedData ? value : value[1], payload: entry, background, tooltipPosition: { x: x + width / 2, y: y + height / 2 }, parentViewBox, originalDataIndex: index, ...(cells && cells[index] ? cells[index].props : {}) }
    })
    .filter(Boolean)
}

// --- grid (cartesian/CartesianGrid.js) ---------------------------------------------------

function gridLine(option, lineItemProps) {
  const { x1, y1, x2, y2 } = lineItemProps
  const rest = omit(svgProps(omit(lineItemProps, ["x1", "y1", "x2", "y2", "key"])), ["offset"])
  return h("line", { ...rest, x1, y1, x2, y2, fill: "none" })
}

// --- the cartesian chart ---------------------------------------------------------------------

let chartCount = 0

/**
 * A mounted cartesian chart. render() rebuilds the SVG from the current
 * state; the persistent elements (wrapper, svg, layers, tooltip and legend
 * wrappers) are kept so focus and CSS transitions survive re-renders.
 */
export class CartesianChart {
  constructor(container, spec, ctx, ids) {
    this.container = container // data-slot=chart
    this.ctx = ctx
    this.config = ctx.config
    this.icons = ctx.icons
    this.ids = ids
    this.slots = new Slots()
    this.anims = new Map() // item id -> animation record
    this.tooltipState = { hover: { active: false, index: null, coordinate: undefined }, keyboard: { active: false, index: null, coordinate: undefined }, dismissed: false, dismissedAt: { x: 0, y: 0 } }
    this.legendBox = { height: 0, left: 0, top: 0, width: 0 }
    this.legendSize = undefined
    this.tooltipBox = { height: 0, left: 0, top: 0, width: 0 }
    this.fontState = new Map() // axis key -> {fontSize, letterSpacing}
    this.clipPathId = `${R.uniqueId("recharts")}-clip`
    this.setSpec(spec, true)
  }

  // Parses the recharts element tree.
  setSpec(spec, initial = false) {
    const root = resolveSpecNode(spec, this.ctx)
    this.spec = spec
    this.name = root.type
    this.props = withDefaults(root.props, CHART_DEFAULTS)
    this.margin = root.props.margin ?? DEFAULT_MARGIN
    this.data = root.props.data ?? []
    this.children = root.children
    const prevItems = this.items ?? []
    this.items = []
    this.xAxes = []
    this.yAxes = []
    this.grids = []
    this.tooltipProps = undefined
    this.tooltipContent = undefined
    this.legendProps = undefined
    this.defs = []
    let n = 1
    for (const child of root.children) {
      const p = child.props
      switch (child.type) {
        case "XAxis":
        case "YAxis": {
          const axisType = child.type === "XAxis" ? "xAxis" : "yAxis"
          const props = withDefaults(p, axisType === "xAxis" ? X_AXIS_DEFAULTS : Y_AXIS_DEFAULTS)
          const settings = {}
          for (const k of AXIS_SETTING_KEYS) if (k in props) settings[k] = props[k]
          if (axisType === "xAxis") delete settings.width
          settings.id = axisType === "xAxis" ? props.xAxisId : props.yAxisId
          settings.type = axisTypeBasedOnLayout(this.props.layout, axisType, props.type)
          ;(axisType === "xAxis" ? this.xAxes : this.yAxes).push({ props, settings, node: child })
          break
        }
        case "CartesianGrid":
          this.grids.push(withDefaults(p, GRID_DEFAULTS))
          break
        case "Tooltip": {
          this.tooltipProps = withDefaults(p, TOOLTIP_DEFAULTS)
          this.tooltipContent = child.children.find((c) => c.type === "ChartTooltipContent")
          break
        }
        case "Legend": {
          this.legendProps = withDefaults(p, LEGEND_DEFAULTS)
          this.legendContent = child.children.find((c) => c.type === "ChartLegendContent")
          break
        }
        case "Area":
        case "Bar":
        case "Line": {
          const type = child.type.toLowerCase()
          const defaults = type === "area" ? AREA_DEFAULTS : type === "bar" ? BAR_DEFAULTS : LINE_DEFAULTS
          const props = withDefaults(p, defaults)
          const id = props.id ?? `recharts-${type}-${this.ids.item(n++)}`
          const prev = prevItems.find((x) => x.id === id)
          this.items.push({
            type,
            id,
            props,
            node: child,
            dataKey: props.dataKey,
            xAxisId: props.xAxisId,
            yAxisId: props.yAxisId,
            stackId: props.stackId == null ? undefined : String(props.stackId),
            hide: props.hide,
            data: undefined,
            barSize: type === "bar" ? props.barSize : undefined,
            maxBarSize: type === "bar" ? props.maxBarSize : undefined,
            minPointSize: props.minPointSize,
            baseValue: props.baseValue,
            connectNulls: props.connectNulls,
            hasCustomShape: type === "bar" ? props.shape != null : undefined,
            cells: child.children.filter((c) => c.type === "Cell"),
            labelLists: child.children.filter((c) => c.type === "LabelList"),
            prev,
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

  get legendSettings() {
    if (!this.legendProps) return undefined
    const { layout, align, verticalAlign, itemSorter } = this.legendProps
    return { layout, align, verticalAlign, itemSorter }
  }

  // --- DOM skeleton -----------------------------------------------------------------------

  mountSkeleton(inner) {
    this.inner = inner
    const wrapper = h("div", { className: "recharts-wrapper" })
    this.wrapper = wrapper
    inner.appendChild(wrapper)
    this.tooltipWrapper = null
    this.legendWrapper = null
    this.svg = null
    this.layers = new Map()
    this.attachEvents()
  }

  ensureSvg() {
    if (this.svg) return
    const svg = h("svg", { role: this.props.accessibilityLayer ? "application" : undefined, tabIndex: this.props.accessibilityLayer ? 0 : undefined, className: "recharts-surface", style: { width: "100%", height: "100%", display: "block" } })
    svg.appendChild(h("title", null))
    svg.appendChild(h("desc", null))
    const zs = [...new Set(Object.values(DefaultZIndexes))].sort((a, b) => a - b)
    for (const z of zs.filter((z) => z < 0)) svg.appendChild(this.layer(z))
    this.zeroLayerStart = document.createComment("")
    this.zeroLayerEnd = document.createComment("")
    svg.appendChild(this.zeroLayerStart)
    svg.appendChild(this.zeroLayerEnd)
    for (const z of zs.filter((z) => z > 0)) svg.appendChild(this.layer(z))
    this.svg = svg
    this.wrapper.appendChild(svg)
  }

  layer(z) {
    const g = h("g", { tabIndex: -1, className: `recharts-zIndex-layer_${z}` })
    this.layers.set(z, g)
    return g
  }

  // Content rendered in place (zIndex 0): clip path defs and user <defs>.
  setZeroLayer(nodes) {
    const old = []
    for (let n = this.zeroLayerStart.nextSibling; n && n !== this.zeroLayerEnd; n = n.nextSibling) old.push(n)
    const k = Math.min(old.length, nodes.length)
    for (let i = 0; i < k; i++) patch(old[i], nodes[i])
    for (let i = k; i < old.length; i++) old[i].remove()
    for (let i = k; i < nodes.length; i++) this.svg.insertBefore(nodes[i], this.zeroLayerEnd)
  }

  put(z, key, nodes) {
    if (z === 0 || z === undefined) return nodes
    this.slots.set(this.layers.get(z), key, Array.isArray(nodes) ? nodes : [nodes])
    return []
  }

  // --- rendering ---------------------------------------------------------------------------

  setSize(width, height) {
    this.width = width
    this.height = height
  }

  /** Renders everything for the current state. */
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
    this.ensureLegend()
    this.ensureSvg()
    syncElement(this.svg, { role: this.props.accessibilityLayer ? "application" : undefined, tabIndex: this.props.accessibilityLayer ? 0 : undefined, className: "recharts-surface", width: w, height: hgt, viewBox: `0 0 ${w} ${hgt}` }, { width: "100%", height: "100%", display: "block" })
    const state = new CartesianState(this)
    this.state = state
    const rendered = new Set()
    const zero = []
    // ClipPathProvider
    const pa = state.plotArea()
    zero.push(h("defs", null, h("clipPath", { id: this.clipPathId }, h("rect", { x: pa.x, y: pa.y, height: pa.height, width: pa.width }))))
    let gridIndex = 0
    let xIndex = 0
    let yIndex = 0
    let refIndex = 0
    for (const child of this.children) {
      if (!state.registered && child.type !== "CartesianGrid") continue
      switch (child.type) {
        case "CartesianGrid":
          this.renderGrid(state, this.grids[gridIndex], `grid-${gridIndex++}`, rendered)
          break
        case "XAxis":
          this.renderAxis(state, "xAxis", this.xAxes[xIndex++], rendered)
          break
        case "YAxis":
          this.renderAxis(state, "yAxis", this.yAxes[yIndex++], rendered)
          break
        case "Area":
        case "Bar":
        case "Line":
        case "Pie": {
          const item = this.items.find((i) => i.node === child)
          if (item) this.renderItem(state, item, rendered)
          break
        }
        case "Tooltip":
          this.renderCursor(state, rendered)
          break
        case "ReferenceDot":
          this.renderReferenceDot(state, child, `reference-dot-${refIndex++}`, rendered)
          break
        case "defs":
          zero.push(renderSpecSvg(child))
          break
      }
    }
    this.setZeroLayer(zero)
    for (const key of this.slots.keys()) if (!rendered.has(key)) this.slots.clear(key)
    this.renderTooltip(state)
    this.renderLegend(state)
  }

  // Puts nodes in layer z under key; returns the nodes now in the document.
  slot(rendered, z, key, nodes) {
    rendered.add(key)
    const list = (Array.isArray(nodes) ? nodes : [nodes]).filter(Boolean)
    return this.slots.set(this.layers.get(z), key, list)
  }

  // ReferenceDot (ifOverflow "discard": nothing outside the plot area).
  renderReferenceDot(state, node, key, rendered) {
    const p = withDefaults(node.props, { r: 10, fill: "#fff", stroke: "#ccc", fillOpacity: 1, strokeWidth: 1, xAxisId: 0, yAxisId: 0, zIndex: DefaultZIndexes.scatter })
    const xs = state.axis("xAxis", p.xAxisId).scale
    const ys = state.axis("yAxis", p.yAxisId).scale
    const cx = xs?.map(p.x, { position: "middle" })
    const cy = ys?.map(p.y, { position: "middle" })
    if (!isNum(cx) || !isNum(cy) || !xs.isInRange(cx) || !ys.isInRange(cy)) {
      this.slot(rendered, p.zIndex, key, [])
      return
    }
    const dot = Dot({ ...svgProps(omit(p, ["x", "y", "xAxisId", "yAxisId", "zIndex", "ifOverflow"])), className: "recharts-reference-dot-dot", cx, cy, r: p.r })
    this.slot(rendered, p.zIndex, key, [Layer({ className: "recharts-reference-dot" }, dot)])
  }

  // CartesianGrid
  renderGrid(state, props, key, rendered) {
    const offset = state.offset
    const p = { ...props, x: isNum(props.x) ? props.x : offset.left, y: isNum(props.y) ? props.y : offset.top, width: isNum(props.width) ? props.width : offset.width, height: isNum(props.height) ? props.height : offset.height }
    const { x, y, width, height } = p
    if (!(width > 0) || !(height > 0) || !isNum(x) || !isNum(y)) {
      this.slot(rendered, p.zIndex, key, [])
      return
    }
    const xAxis = state.gridAxisProps("xAxis", p.xAxisId)
    const yAxis = state.gridAxisProps("yAxis", p.yAxisId)
    let horizontalPoints = p.horizontalPoints
    let verticalPoints = p.verticalPoints
    if (!horizontalPoints || !horizontalPoints.length) {
      horizontalPoints = R.getCoordinatesOfGrid(R.getTicks({ ...AXIS_DEFAULT, x: 0, y: 0, width: 0, height: 0, viewBox: { x: 0, y: 0, width: 0, height: 0 }, orientation: "bottom", ticks: [], ...yAxis, ticks: R.getTicksOfAxis(yAxis, true), viewBox: { x: 0, y: 0, width: state.width, height: state.height } }), offset.top, offset.top + offset.height, p.syncWithTicks)
    }
    if (!verticalPoints || !verticalPoints.length) {
      verticalPoints = R.getCoordinatesOfGrid(R.getTicks({ ...AXIS_DEFAULT, x: 0, y: 0, width: 0, height: 0, viewBox: { x: 0, y: 0, width: 0, height: 0 }, orientation: "bottom", ticks: [], ...xAxis, ticks: R.getTicksOfAxis(xAxis, true), viewBox: { x: 0, y: 0, width: state.width, height: state.height } }), offset.left, offset.left + offset.width, p.syncWithTicks)
    }
    const children = []
    if (p.fill && p.fill !== "none") children.push(h("rect", { x, y, ry: p.ry, width, height, stroke: "none", fill: p.fill, fillOpacity: p.fillOpacity, className: "recharts-cartesian-grid-bg" }))
    const lineProps = omit(p, ["xAxisId", "yAxisId"])
    if (p.horizontal && horizontalPoints?.length) {
      children.push(h("g", { className: "recharts-cartesian-grid-horizontal" }, horizontalPoints.map((entry, i) => gridLine(p.horizontal, { ...lineProps, offset, horizontalPoints, x1: x, y1: entry, x2: x + width, y2: entry, index: i }))))
    }
    if (p.vertical && verticalPoints?.length) {
      children.push(h("g", { className: "recharts-cartesian-grid-vertical" }, verticalPoints.map((entry, i) => gridLine(p.vertical, { ...lineProps, offset, verticalPoints, x1: entry, y1: y, x2: entry, y2: y + height, index: i }))))
    }
    this.slot(rendered, p.zIndex, key, h("g", { className: "recharts-cartesian-grid" }, children))
  }

  // XAxis / YAxis through CartesianAxis.
  renderAxis(state, axisType, axisDef, rendered) {
    const { settings } = axisDef
    const key = `${axisType}-${settings.id}`
    const a = state.axis(axisType, settings.id)
    const pos = state.axisPosition(axisType, settings)
    const size = state.axisSize(axisType, settings)
    const props = withDefaults(
      {
        ...omit(axisDef.props, ["ticks", "scale", "dangerouslySetInnerHTML"]),
        ...omit(settings, ["id", "scale"]),
        x: pos.x,
        y: pos.y,
        width: size.width,
        height: size.height,
        className: clsx(`recharts-${axisType} ${axisType}`, axisDef.props.className),
        viewBox: { x: 0, y: 0, width: state.width, height: state.height },
        ticks: a.ticks,
        axisType,
        axisId: settings.id,
        ...(axisType === "yAxis" ? { tickTextProps: settings.width === "auto" ? { width: undefined } : { width: settings.width } } : {}),
      },
      { ...AXIS_DEFAULT, x: 0, y: 0, width: 0, height: 0, orientation: "bottom", ticks: [], zIndex: DefaultZIndexes.axis },
    )
    const labelsKey = `${key}-labels`
    if (props.hide || (props.width != null && props.width <= 0) || (props.height != null && props.height <= 0)) {
      this.slot(rendered, props.zIndex, key, [])
      this.slot(rendered, DefaultZIndexes.label, labelsKey, [])
      return
    }
    const font = this.fontState.get(key) ?? { fontSize: "", letterSpacing: "" }
    const finalTicks = R.getTicks({ ...props, ticks: a.ticks ?? [] }, font.fontSize, font.letterSpacing)
    const axisProps = svgProps(props)
    const { orientation, mirror, tickLine, tick, tickFormatter, unit, x, y, width, height, tickSize, tickMargin } = props
    const textAnchor = ["start", "middle", "end", "inherit"].includes(axisProps.textAnchor) ? axisProps.textAnchor : orientation === "left" ? (mirror ? "start" : "end") : orientation === "right" ? (mirror ? "end" : "start") : "middle"
    const verticalAnchor = orientation === "left" || orientation === "right" ? "middle" : orientation === "top" ? (mirror ? "start" : "end") : mirror ? "end" : "start"
    const tickLineProps = { ...axisProps, fill: "none", ...(typeof tickLine === "object" ? tickLine : {}) }
    const coords = finalTicks.map((entry) => ({ entry, ...getTickLineCoord(entry, x, y, width, height, orientation, tickSize, mirror, tickMargin) }))
    const tickLines = coords.map(({ line }) => Layer({ className: "recharts-cartesian-axis-tick" }, tickLine ? h("line", { ...tickLineProps, ...line, className: clsx("recharts-cartesian-axis-tick-line", tickLine?.className) }) : null))
    const customTickProps = svgPropsFromUnknown(tick)
    const tickLabels = coords.map(({ entry, tick: tickCoord }, i) => {
      const tickProps = { verticalAnchor, ...axisProps, textAnchor, stroke: "none", fill: props.stroke, ...tickCoord, index: i, payload: entry, visibleTicksCount: finalTicks.length, tickFormatter, padding: props.padding, ...props.tickTextProps, angle: props.tickTextProps?.angle ?? axisProps.angle ?? 0 }
      const finalTickProps = { ...tickProps, ...customTickProps }
      const value = `${typeof tickFormatter === "function" ? tickFormatter(entry.value, i) : entry.value}${unit || ""}`
      let className = "recharts-cartesian-axis-tick-value"
      if (typeof tick === "object" && tick?.className) className = clsx(className, tick.className)
      return Layer({ className: "recharts-cartesian-axis-tick-label" }, tick ? Text({ ...finalTickProps, className, children: value }) : null)
    })
    this.renderedTicks ??= new Map()
    this.renderedTicks.set(key, finalTicks)
    const labelsG = tickLabels.length > 0 ? h("g", { className: `recharts-cartesian-axis-tick-labels recharts-${axisType}-tick-labels` }, tickLabels) : null
    const placedLabels = this.slot(rendered, DefaultZIndexes.label, labelsKey, labelsG)[0]
    const axisLine = renderAxisLine({ x: props.x, y: props.y, width, height, orientation, mirror, axisLine: props.axisLine, otherSvgProps: svgProps(props) })
    const ticksG = h("g", { className: `recharts-cartesian-axis-ticks recharts-${axisType}-ticks` }, tickLines.length > 0 ? h("g", { className: `recharts-cartesian-axis-tick-lines recharts-${axisType}-tick-lines` }, tickLines) : null)
    this.slot(rendered, props.zIndex, key, Layer({ className: clsx("recharts-cartesian-axis", props.className) }, axisLine, ticksG))
    // remember the first tick text for the font size measurement (CartesianAxis layerRef)
    this.firstTick ??= new Map()
    this.firstTick.set(key, placedLabels?.querySelector(".recharts-cartesian-axis-tick-value") ?? null)
  }

  // Reads the computed font of the first tick of each axis, like
  // CartesianAxis does after its first render. Returns true if it changed.
  measureFonts() {
    let changed = false
    for (const [key, el] of this.firstTick ?? []) {
      if (!el || !el.isConnected) continue
      const cs = getComputedStyle(el)
      const prev = this.fontState.get(key) ?? { fontSize: "", letterSpacing: "" }
      if (cs.fontSize !== prev.fontSize || cs.letterSpacing !== prev.letterSpacing) {
        this.fontState.set(key, { fontSize: cs.fontSize, letterSpacing: cs.letterSpacing })
        changed = true
      }
    }
    return changed
  }

  // --- graphical items ------------------------------------------------------------------

  isAnimationActive(item) {
    const v = item.props.isAnimationActive
    return v === "auto" ? !prefersReducedMotion() : v
  }

  // Animation bookkeeping (useAnimationId + JavascriptAnimate): a new
  // animation starts whenever the animation input changes.
  animationFor(item, inputKey) {
    let rec = this.anims.get(item.id)
    if (!rec) {
      rec = { item: item.id, inputKey: undefined, anim: null, isAnimating: false, prev: null, prevBase: undefined, longest: 0, startingPoint: 0, totalLength: 0 }
      this.anims.set(item.id, rec)
    }
    if (rec.inputKey !== inputKey) {
      rec.inputKey = inputKey
      R.uniqueId(`recharts-${item.type}-`)
      rec.anim?.stop()
      rec.startingPoint = rec.longest
      if (this.isAnimationActive(item)) {
        rec.pendingStart = true
        rec.anim = new Animation(this, {
          duration: item.props.animationDuration,
          begin: item.props.animationBegin,
          easing: item.props.animationEasing,
          onEnd: () => {
            rec.isAnimating = false
            rec.pendingStart = false
            rec.closureStale = true
          },
        })
      } else {
        rec.anim = null
        rec.isAnimating = false
        rec.pendingStart = false
      }
      rec.isNew = true
      rec.closureStale = true
    }
    // The animated render closure reads the previous shape when the
    // component itself renders (new input, animation start or end), not on
    // every animation frame.
    if (rec.closureStale) {
      rec.closureStale = false
      rec.closurePrev = rec.prev
      rec.closurePrevBase = rec.prevBase
      rec.closureLength = rec.totalLength
    }
    return rec
  }

  tOf(rec) {
    if (!rec.anim) return 1
    return rec.anim.done ? 1 : rec.anim.t
  }

  renderItem(state, item, rendered) {
    if (item.type === "area") return this.renderArea(state, item, rendered)
    if (item.type === "bar") return this.renderBar(state, item, rendered)
    if (item.type === "line") return this.renderLine(state, item, rendered)
  }

  renderArea(state, item, rendered) {
    const props = item.props
    const key = item.id
    const layout = state.layout
    if ((layout !== "horizontal" && layout !== "vertical") || (this.name !== "AreaChart" && this.name !== "ComposedChart")) return
    const computed = state.computeArea(item)
    if (!computed || !computed.points?.length || props.hide) {
      this.slot(rendered, props.zIndex, key, [])
      return
    }
    const { points, baseLine, isRange } = computed
    const plot = state.plotArea()
    const allProps = { ...omit(props, ["activeDot", "animationBegin", "animationDuration", "animationEasing", "connectNulls", "dot", "fill", "fillOpacity", "hide", "isAnimationActive", "legendType", "stroke", "xAxisId", "yAxisId"]), ...props, baseLine, height: plot.height, width: plot.width, left: plot.x, top: plot.y, layout, isRange, points, needClip: false, id: item.id }
    const inputKey = JSON.stringify([points.map((p) => [p.x, p.y]), Array.isArray(baseLine) ? baseLine.map((p) => [p.x, p.y]) : baseLine])
    const rec = this.animationFor(item, inputKey)
    const t = this.tOf(rec)
    const prevPoints = rec.closurePrev
    const prevBaseLine = rec.closurePrevBase
    rec.isNew = false
    let content
    const staticArea = (pts, bl) => this.staticArea(pts, bl, allProps)
    if (prevPoints) {
      const factor = prevPoints.length / points.length
      const stepPoints = t === 1 ? points : points.map((entry, index) => {
        const prev = prevPoints[Math.floor(index * factor)]
        return prev ? { ...entry, x: R.interpolate(prev.x, entry.x, t), y: R.interpolate(prev.y, entry.y, t) } : entry
      })
      let stepBaseLine
      if (isNum(baseLine)) stepBaseLine = R.interpolate(prevBaseLine, baseLine, t)
      else if (isNullish(baseLine) || R.isNan(baseLine)) stepBaseLine = R.interpolate(prevBaseLine, 0, t)
      else stepBaseLine = baseLine.map((entry, index) => {
        const prev = Array.isArray(prevBaseLine) ? prevBaseLine[Math.floor(index * factor)] : undefined
        return prev ? { ...entry, x: R.interpolate(prev.x, entry.x, t), y: R.interpolate(prev.y, entry.y, t) } : entry
      })
      if (t > 0) {
        rec.prev = stepPoints
        rec.prevBase = stepBaseLine
      }
      content = staticArea(stepPoints, stepBaseLine)
    } else {
      if (t > 0) {
        rec.prev = points
        rec.prevBase = baseLine
      }
      const clipId = `animationClipPath-${item.id}`
      content = [Layer({}, this.isAnimationActive(item) ? h("defs", null, h("clipPath", { id: clipId }, clipRect({ alpha: t, points, baseLine, layout, strokeWidth: props.strokeWidth }))) : null, Layer({ clipPath: `url(#${clipId})` }, ...staticArea(points, baseLine)))]
    }
    this.slot(rendered, props.zIndex, key, Layer({ className: clsx("recharts-area", props.className) }, ...content))
    this.renderDots(rendered, `${key}-dots`, points, props, "recharts-area-dots", "recharts-area-dot")
    this.renderActivePoints(state, rendered, `${key}-active`, points, getLegendItemColor(props.stroke, props.fill), props.dataKey, props.activeDot)
    this.renderLabelLists(state, rendered, item, rec, points.map((point) => labelEntryFromPoint(point)))
  }

  staticArea(points, baseLine, props) {
    const { layout, type, stroke, connectNulls, isRange } = props
    const propsWithoutId = omit(props, ["id"])
    const allOther = svgProps(propsWithoutId)
    const out = []
    if (points?.length > 1) {
      out.push(
        Layer(
          {},
          Curve({ ...svgProps(propsWithoutId), id: props.id, points, connectNulls, type, baseLine, layout, stroke: "none", className: "recharts-area-area" }),
          stroke !== "none" ? Curve({ ...allOther, className: "recharts-area-curve", layout, type, connectNulls, fill: "none", points }) : null,
          stroke !== "none" && isRange && Array.isArray(baseLine) ? Curve({ ...allOther, className: "recharts-area-curve", layout, type, connectNulls, fill: "none", points: baseLine }) : null,
        ),
      )
    }
    return out
  }

  renderBar(state, item, rendered) {
    const props = item.props
    const key = item.id
    const layout = state.layout
    if (layout !== "vertical" && layout !== "horizontal") return
    const cells = item.cells.map((c) => ({ props: c.props }))
    const rects = state.computeBar(item, cells.length ? cells : undefined)
    if (props.hide || rects == null) {
      this.slot(rendered, props.zIndex, key, [])
      return
    }
    const allProps = { ...props, layout, needClip: false, data: rects, id: item.id }
    const inputKey = JSON.stringify([rects.map((r) => [r.x, r.y, r.width, r.height]), props.dataKey, props.fill])
    const rec = this.animationFor(item, inputKey)
    const t = this.tOf(rec)
    const prevData = rec.closurePrev
    rec.isNew = false
    const stepData = t === 1 ? rects : rects.map((entry, index) => {
      const prev = prevData && prevData[index]
      if (prev) return { ...entry, x: R.interpolate(prev.x, entry.x, t), y: R.interpolate(prev.y, entry.y, t), width: R.interpolate(prev.width, entry.width, t), height: R.interpolate(prev.height, entry.height, t) }
      if (layout === "horizontal") return { ...entry, y: R.interpolate(entry.stackedBarStart, entry.y, t), height: R.interpolate(0, entry.height, t) }
      return { ...entry, width: R.interpolate(0, entry.width, t), x: R.interpolate(entry.stackedBarStart, entry.x, t) }
    })
    if (t > 0) rec.prev = stepData ?? null
    ;(this.barEntries ??= new Map()).set(item.id, stepData)
    const baseProps = omit(svgProps(allProps), ["id"])
    const rectangles = stepData.map((entry, i) =>
      Layer(
        { className: "recharts-bar-rectangle" },
        Layer({ className: "recharts-inactive-bar" }, barShape(props.shape, { ...baseProps, name: String(baseProps.name), ...entry, isActive: false, index: i, dataKey: props.dataKey })),
      ),
    )
    this.slot(rendered, props.zIndex, key, Layer({ className: clsx("recharts-bar", props.className), id: item.id }, Layer({ className: "recharts-bar-rectangles" }, Layer({}, rectangles))))
    this.renderLabelLists(state, rendered, item, rec, rects.map((entry) => ({ x: entry.x, y: entry.y, width: entry.width, lowerWidth: entry.width, upperWidth: entry.width, height: entry.height, value: entry.value, payload: entry.payload, parentViewBox: entry.parentViewBox, viewBox: { x: entry.x, y: entry.y, width: entry.width, lowerWidth: entry.width, upperWidth: entry.width, height: entry.height }, fill: entry.fill })))
  }

  renderLine(state, item, rendered) {
    const props = item.props
    const key = item.id
    const layout = state.layout
    const points = state.computeLine(item)
    if ((layout !== "horizontal" && layout !== "vertical") || points == null || props.hide) {
      this.slot(rendered, props.zIndex, key, [])
      return
    }
    const plot = state.plotArea()
    const all = { ...omit(props, ["activeDot", "animateNewValues", "animationBegin", "animationDuration", "animationEasing", "connectNulls", "dot", "hide", "isAnimationActive", "label", "legendType", "xAxisId", "yAxisId", "id"]), id: item.id, connectNulls: props.connectNulls, dot: props.dot, activeDot: props.activeDot, hide: props.hide, label: props.label, legendType: props.legendType, xAxisId: props.xAxisId, yAxisId: props.yAxisId, points, layout, height: plot.height, width: plot.width, left: plot.x, top: plot.y, needClip: false }
    const inputKey = JSON.stringify(points.map((p) => [p.x, p.y]))
    const rec = this.animationFor(item, inputKey)
    const t = this.tOf(rec)
    const prevPoints = rec.closurePrev
    rec.isNew = false
    const totalLength = rec.closureLength
    const startingPoint = rec.startingPoint
    const lengthInterpolated = R.interpolate(startingPoint, totalLength + startingPoint, t)
    const curLength = Math.min(lengthInterpolated, totalLength)
    let dash
    const strokeDasharray = props.strokeDasharray
    if (this.isAnimationActive(item)) {
      if (strokeDasharray) dash = getStrokeDasharray(curLength, totalLength, `${strokeDasharray}`.split(/[,\s]+/gim).map((n) => parseFloat(n)))
      else dash = `${curLength}px ${totalLength}px`
    } else dash = strokeDasharray == null ? undefined : String(strokeDasharray)
    if (t > 0 && totalLength > 0) {
      rec.prev = points
      rec.longest = Math.max(rec.longest, curLength)
    }
    let stepData = points
    if (prevPoints) {
      const factor = prevPoints.length / points.length
      stepData = t === 1 ? points : points.map((entry, index) => {
        const prev = prevPoints[Math.floor(index * factor)]
        if (prev) return { ...entry, x: R.interpolate(prev.x, entry.x, t), y: R.interpolate(prev.y, entry.y, t) }
        if (props.animateNewValues) return { ...entry, x: R.interpolate(plot.width * 2, entry.x, t), y: R.interpolate(plot.height / 2, entry.y, t) }
        return { ...entry, x: entry.x, y: entry.y }
      })
      rec.prev = stepData
    }
    const others = omit(all, ["type", "layout", "connectNulls", "needClip", "shape"])
    const curveProps = { ...svgProps(others), fill: "none", className: "recharts-line-curve", points: stepData, type: props.type, layout, connectNulls: props.connectNulls, strokeDasharray: dash ?? props.strokeDasharray }
    const curve = stepData?.length > 1 ? Layer({ className: "recharts-shape" }, Curve(curveProps)) : null
    const placed = this.slot(rendered, props.zIndex, key, Layer({ className: clsx("recharts-line", props.className) }, curve))[0]
    rec.pathEl = curve ? placed?.querySelector("path.recharts-line-curve") : null
    this.renderDots(rendered, `${key}-dots`, stepData, all, "recharts-line-dots", "recharts-line-dot")
    this.renderActivePoints(state, rendered, `${key}-active`, points, props.stroke, props.dataKey, props.activeDot)
    this.renderLabelLists(state, rendered, item, rec, points.map((point) => labelEntryFromPoint(point)))
  }

  // Line: getTotalLength of the path rendered last time (recharts reads it in render).
  updateLineLengths() {
    let changed = false
    for (const item of this.items) {
      if (item.type !== "line") continue
      const rec = this.anims.get(item.id)
      if (!rec?.pathEl?.isConnected) continue
      let len = 0
      try {
        len = rec.pathEl.getTotalLength() || 0
      } catch {
        len = 0
      }
      if (len !== rec.totalLength) {
        rec.totalLength = len
        changed = true
      }
    }
    return changed
  }

  renderDots(rendered, key, points, baseAllProps, className, dotClassName) {
    const dot = baseAllProps.dot
    if (points == null || !(dot || points.length === 1)) {
      this.slot(rendered, DefaultZIndexes.scatter, key, [])
      return
    }
    const baseProps = svgProps(omit(baseAllProps, ["id"]))
    const customDotProps = svgPropsFromUnknown(dot) ?? {}
    const dots = points.map((entry, i) => {
      const dotProps = { r: 3, ...baseProps, ...customDotProps, index: i, cx: entry.x ?? undefined, cy: entry.y ?? undefined, dataKey: baseAllProps.dataKey, value: entry.value, payload: entry.payload, points }
      if (typeof dot === "function") return dot(dotProps)
      const finalClassName = clsx(dotClassName, typeof dot !== "boolean" ? dot.className : "")
      return Dot({ ...omit(dotProps, ["points"]), className: finalClassName })
    })
    this.slot(rendered, DefaultZIndexes.scatter, key, Layer({ className }, dots))
  }

  renderActivePoints(state, rendered, key, points, mainColor, dataKey, activeDot) {
    const active = this.activeTooltip(state)
    let node = null
    if (points != null && active.dataPoints != null) {
      const point = points.find((p) => active.dataPoints.includes(p.payload))
      if (point && activeDot !== false && point.x != null && point.y != null) {
        const dotProps = { index: Number(active.index), dataKey, cx: point.x, cy: point.y, r: 4, fill: mainColor ?? "none", strokeWidth: 2, stroke: "#fff", payload: point.payload, value: point.value, ...(svgPropsFromUnknown(activeDot) ?? {}) }
        const dot = typeof activeDot === "function" ? activeDot(dotProps) : Dot(dotProps)
        node = Layer({ className: "recharts-active-dot" }, dot)
      }
    }
    this.slot(rendered, DefaultZIndexes.activeDot, key, node)
  }

  renderLabelLists(state, rendered, item, rec, entries) {
    const showLabels = !rec.isAnimating
    const lists = [...item.labelLists]
    const label = item.props.label
    if (label === true) lists.push({ props: {} })
    else if (label && typeof label === "object" && !label.$fn) lists.push({ props: { ...label } })
    lists.forEach((list, i) => {
      const key = `${item.id}-labels-${i}`
      if (!showLabels || !entries?.length) {
        this.slot(rendered, DefaultZIndexes.label, key, [])
        return
      }
      this.slot(rendered, list.props.zIndex ?? DefaultZIndexes.label, key, renderLabelList(list.props, entries, state))
    })
  }

  // --- tooltip ----------------------------------------------------------------------------

  tooltipEventType() {
    return "axis"
  }

  interactionState() {
    const ts = this.tooltipState
    if (ts.hover.active) return ts.hover
    if (ts.keyboard.active) return ts.keyboard
    // Tooltip defaultIndex: active until the pointer has activated an index.
    const di = this.tooltipProps?.defaultIndex
    if (ts.hover.index == null && di != null) return { active: true, index: String(di), dataKey: undefined, coordinate: undefined, graphicalItemId: undefined }
    return { active: false, index: null, dataKey: undefined, coordinate: ts.hover.coordinate, graphicalItemId: undefined }
  }

  // selectActiveTooltipIndex, payload, label, coordinate, data points.
  activeTooltip(state) {
    const inter = this.interactionState()
    const axisType = state.tooltipAxisType()
    const axis = state.axis(axisType, this.tooltipProps?.axisId ?? 0)
    const data = this.data
    let index = null
    if (inter.index != null) {
      const n = Number(inter.index)
      if (!wellBehaved(n)) index = inter.index
      else {
        const upper = data.length > 0 ? data.length - 1 : Infinity
        const clamped = Math.max(0, Math.min(n, upper))
        index = String(clamped)
      }
    }
    const ticks = axis.tooltipTicks
    const n = Number(index)
    const label = index == null || R.isNan(n) ? undefined : n >= 0 ? ticks?.[n]?.value : undefined
    let payload
    if (index != null) {
      payload = []
      for (const item of this.items) {
        const settings = this.tooltipEntrySettings(item)
        const dataKey = settings.dataKey ?? axis.settings.dataKey
        let entry
        if (axis.settings.dataKey && Array.isArray(data) && !Array.isArray(data[0])) entry = data.find((e) => R.getValueByDataKey(e, axis.settings.dataKey) === label)
        else entry = data[Number(index)]
        payload.push({ ...settings, dataKey, payload: entry, value: R.getValueByDataKey(entry, dataKey), name: R.getValueByDataKey(entry, settings.nameKey) ?? settings.name })
      }
    }
    const dataPoints = payload == null ? undefined : Array.from(new Set(payload.map((p) => p.payload).filter((p) => p != null)))
    let coordinate = inter.coordinate
    // selectActiveTooltipCoordinate: the defaultIndex coordinate otherwise.
    if (coordinate == null && this.tooltipProps?.defaultIndex != null) coordinate = this.coordinateForIndex(state, String(this.tooltipProps.defaultIndex))
    return { active: inter.active && index != null, index, label, payload, coordinate, dataPoints, axis }
  }

  tooltipEntrySettings(item) {
    const p = item.props
    const name = p.name ? String(p.name) : typeof p.dataKey === "string" ? p.dataKey : undefined
    const color = item.type === "area" ? getLegendItemColor(p.stroke, p.fill) : item.type === "bar" ? p.fill : p.stroke
    return { stroke: p.stroke, strokeWidth: p.strokeWidth, fill: p.fill, dataKey: p.dataKey, nameKey: undefined, name, hide: p.hide, type: p.tooltipType, color, unit: p.unit, graphicalItemId: item.id }
  }

  coordinateForIndex(state, index) {
    const axis = state.axis(state.tooltipAxisType(), this.tooltipProps?.axisId ?? 0)
    const tick = axis.tooltipTicks?.[Number(index)]
    if (!tick) return undefined
    if (state.layout === "horizontal") return { x: tick.coordinate, y: (state.offset.top + state.height) / 2 }
    return { x: (state.offset.left + state.width) / 2, y: tick.coordinate }
  }

  renderCursor(state, rendered) {
    const key = "cursor"
    const tp = this.tooltipProps
    const active = this.activeTooltip(state)
    const cursor = tp?.cursor
    let node = null
    let z = DefaultZIndexes.cursorLine
    if (tp && active.active && cursor && active.coordinate) {
      const offset = state.offset
      const axis = active.axis
      const bandSize = R.getBandSizeOfAxis(axis.axisWithScale, axis.tooltipTicks)
      let restProps
      let comp
      if (this.name === "BarChart") {
        restProps = R.getCursorRectangle(state.layout, active.coordinate, offset, bandSize)
        comp = Rectangle
        z = DefaultZIndexes.cursorRectangle
      } else {
        const c = active.coordinate
        restProps = { points: state.layout === "horizontal" ? [{ x: c.x, y: offset.top }, { x: c.x, y: offset.top + offset.height }] : [{ x: offset.left, y: c.y }, { x: offset.left + offset.width, y: c.y }] }
        comp = Curve
      }
      const extraClassName = typeof cursor === "object" && "className" in cursor ? cursor.className : undefined
      const cursorProps = { stroke: "#ccc", pointerEvents: "none", ...offset, ...restProps, ...(svgPropsFromUnknown(cursor) ?? {}), payload: active.payload, payloadIndex: active.index, className: clsx("recharts-tooltip-cursor", extraClassName) }
      if (typeof cursor === "object" && cursor.zIndex != null) z = cursor.zIndex
      node = comp(cursorProps)
    }
    // the cursor slot can change layers (line vs rectangle)
    for (const zz of [DefaultZIndexes.cursorLine, DefaultZIndexes.cursorRectangle]) if (zz !== z) this.slot(rendered, zz, `${key}-${zz}`, [])
    this.slot(rendered, z, `${key}-${z}`, node)
  }

  ensureTooltipWrapper() {
    if (!this.tooltipProps || this.tooltipWrapper) return
    const div = h("div", { xmlns: "http://www.w3.org/1999/xhtml", tabIndex: -1, className: "recharts-tooltip-wrapper" })
    this.tooltipWrapper = div
    this.wrapper.insertBefore(div, this.wrapper.firstChild)
  }

  renderTooltip(state) {
    const tp = this.tooltipProps
    const div = this.tooltipWrapper
    if (!tp || !div) return
    const active = this.activeTooltip(state)
    const finalIsActive = tp.active ?? active.active ?? false
    let payload = active.payload ?? []
    if (!finalIsActive) payload = []
    if (tp.filterNull && payload.length) payload = payload.filter((entry) => entry.value != null && (entry.hide !== true || tp.includeHidden))
    const hasPayload = payload.length > 0
    // content
    const contentEl = this.tooltipContent ? renderTooltipContent(this.tooltipContent.props, { active: finalIsActive, payload, label: active.label, config: this.config, icons: this.icons }) : null
    patchChildren(div, contentEl ? [contentEl] : [])
    this.emitTooltip(finalIsActive, active)
    // useElementOffset([payload, finalIsActive]): the bounding box is read
    // by a callback ref, which React calls again only when its deps change:
    // the payload (memoized on the item settings, data, active index and
    // label), finalIsActive or the last box.
    const deps = [JSON.stringify(this.items.map((item) => this.tooltipEntrySettings(item))), this.data, active.index, active.label, finalIsActive]
    const prev = this.tooltipMeasureDeps
    const measure = !prev || deps.some((d, i) => d !== prev[i])
    this.tooltipMeasureDeps = deps
    this.layoutTooltip(state, active.coordinate, finalIsActive, hasPayload, measure)
  }

  // TooltipBoundingBox: position from the last measured box; re-measures
  // and repositions (the change animates through the CSS transition).
  layoutTooltip(state, coordinate, active, hasPayload, measure = true) {
    const tp = this.tooltipProps
    const div = this.tooltipWrapper
    const ts = this.tooltipState
    if (ts.dismissed && ((coordinate?.x ?? 0) !== ts.dismissedAt.x || (coordinate?.y ?? 0) !== ts.dismissedAt.y)) ts.dismissed = false
    const apply = () => {
      const viewBox = state.plotArea()
      const { cssClasses, cssProperties } = R.getTooltipTranslate({ allowEscapeViewBox: tp.allowEscapeViewBox, coordinate, offsetLeft: typeof tp.offset === "number" ? tp.offset : tp.offset.x, offsetTop: typeof tp.offset === "number" ? tp.offset : tp.offset.y, position: tp.position, reverseDirection: tp.reverseDirection, tooltipBox: { height: this.tooltipBox.height, width: this.tooltipBox.width }, useTranslate3d: tp.useTranslate3d, viewBox })
      const animate = tp.isAnimationActive === "auto" && prefersReducedMotion() ? undefined : tp.isAnimationActive && active ? `transform ${tp.animationDuration}ms ${tp.animationEasing}` : undefined
      const style = { transition: animate, ...cssProperties, pointerEvents: "none", position: "absolute", top: 0, left: 0, visibility: !ts.dismissed && active && hasPayload ? "visible" : "hidden", ...tp.wrapperStyle }
      syncElement(div, { xmlns: "http://www.w3.org/1999/xhtml", tabIndex: -1, className: cssClasses }, style)
    }
    apply()
    if (!measure) return
    // useElementOffset: measure, update if changed by more than 1px.
    const rect = div.getBoundingClientRect()
    const box = { height: rect.height, left: rect.left, top: rect.top, width: rect.width }
    const prev = this.tooltipBox
    if (Math.abs(box.height - prev.height) > 1 || Math.abs(box.left - prev.left) > 1 || Math.abs(box.top - prev.top) > 1 || Math.abs(box.width - prev.width) > 1) {
      this.tooltipBox = box
      apply()
      const r2 = div.getBoundingClientRect()
      const b2 = { height: r2.height, left: r2.left, top: r2.top, width: r2.width }
      if (Math.abs(b2.height - box.height) > 1 || Math.abs(b2.left - box.left) > 1 || Math.abs(b2.top - box.top) > 1 || Math.abs(b2.width - box.width) > 1) {
        this.tooltipBox = b2
        apply()
      }
    }
  }

  emitTooltip(active, info) {
    const key = `${active}:${info.index}`
    if (this.lastTooltipKey === key) return
    this.lastTooltipKey = key
    if (this.lastTooltipKey !== undefined) emit(this.container, "shadcn:chart-tooltip", { active, index: info.index == null ? null : Number(info.index), label: info.label })
  }

  // --- legend --------------------------------------------------------------------------

  ensureLegend() {
    if (!this.legendProps || this.legendWrapper) return
    const div = h("div", { className: "recharts-legend-wrapper" })
    this.legendWrapper = div
    const after = this.tooltipWrapper ? this.tooltipWrapper.nextSibling : this.wrapper.firstChild
    this.wrapper.insertBefore(div, after)
  }

  legendPayload() {
    const payload = []
    for (const item of this.items) {
      const p = item.props
      const color = item.type === "area" ? getLegendItemColor(p.stroke, p.fill) : item.type === "bar" ? p.fill : p.stroke
      payload.push({ inactive: p.hide, dataKey: p.dataKey, type: p.legendType, color, value: p.name ? String(p.name) : typeof p.dataKey === "string" ? p.dataKey : undefined, payload: p })
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

  renderLegend(state) {
    const lp = this.legendProps
    const div = this.legendWrapper
    if (!lp || !div) return
    const margin = this.margin
    const maxWidth = this.width - (margin?.left || 0) - (margin?.right || 0)
    const widthOrHeight = lp.layout === "vertical" && lp.height != null ? { height: lp.height } : lp.layout === "horizontal" ? { width: lp.width || maxWidth } : null
    const box = this.legendBox
    const pos = {}
    const ws = lp.wrapperStyle
    if (!ws || (ws.left == null && ws.right == null)) {
      if (lp.align === "center" && lp.layout === "vertical") pos.left = ((this.width || 0) - box.width) / 2
      else if (lp.align === "right") pos.right = (margin && margin.right) || 0
      else pos.left = (margin && margin.left) || 0
    }
    if (!ws || (ws.top == null && ws.bottom == null)) {
      if (lp.verticalAlign === "middle") pos.top = ((this.height || 0) - box.height) / 2
      else if (lp.verticalAlign === "bottom") pos.bottom = (margin && margin.bottom) || 0
      else pos.top = (margin && margin.top) || 0
    }
    const style = { position: "absolute", width: widthOrHeight?.width || lp.width || "auto", height: widthOrHeight?.height || lp.height || "auto", ...pos, ...ws }
    syncElement(div, { className: "recharts-legend-wrapper" }, style)
    const payload = state.registered ? this.legendPayload() : []
    const content = this.legendContent ? renderLegendContent(this.legendContent.props, { payload, verticalAlign: lp.verticalAlign, config: this.config, icons: this.icons }) : null
    patchChildren(div, content ? [content] : [])
  }

  // Measures the legend (useElementOffset); true if the chart must re-render.
  measureLegend() {
    if (!this.legendWrapper || !this.legendWrapper.isConnected) return false
    const rect = this.legendWrapper.getBoundingClientRect()
    const box = { height: rect.height, left: rect.left, top: rect.top, width: rect.width }
    const prev = this.legendBox
    if (Math.abs(box.height - prev.height) > 1 || Math.abs(box.left - prev.left) > 1 || Math.abs(box.top - prev.top) > 1 || Math.abs(box.width - prev.width) > 1) {
      this.legendBox = box
    }
    const size = { width: this.legendBox.width, height: this.legendBox.height }
    if (!this.legendSize || this.legendSize.width !== size.width || this.legendSize.height !== size.height) {
      this.legendSize = size
      return true
    }
    return false
  }

  /**
   * Renders like a React commit followed by the effects that re-render:
   * axis fonts and legend size (layout effects), then animations starting
   * (passive effects), which hide label lists while animating.
   */
  commit() {
    this.render()
    for (let i = 0; i < 4; i++) {
      let again = false
      if (this.measureFonts()) again = true
      if (this.measureLegend()) again = true
      if (this.updateLineLengths()) again = true
      if (!again) break
      this.render()
    }
    let started = false
    for (const rec of this.anims.values()) {
      if (rec.pendingStart) {
        rec.pendingStart = false
        rec.isAnimating = true
        rec.closureStale = true
        started = true
      }
    }
    if (started) this.render()
  }

  // The renders of a fresh chart: before registration, then committed.
  mountSequence() {
    this.registered = false
    this.render()
    this.registered = true
    this.commit()
  }

  // --- events --------------------------------------------------------------------------

  attachEvents() {
    const w = this.wrapper
    let raf = null
    let latest = null
    const pointer = (e) => {
      const rect = w.getBoundingClientRect()
      const scaleX = w.offsetWidth > 0 ? rect.width / w.offsetWidth : 1
      const scaleY = w.offsetHeight > 0 ? rect.height / w.offsetHeight : 1
      return { relativeX: Math.round((e.clientX - rect.left) / scaleX), relativeY: Math.round((e.clientY - rect.top) / scaleY) }
    }
    const move = (e) => {
      latest = pointer(e)
      if (raf !== null) cancelAnimationFrame(raf)
      raf = requestAnimationFrame(() => {
        raf = null
        if (!latest || !this.state) return
        const props = this.activePropsFromPointer(this.state, latest)
        if (props?.activeIndex != null) {
          // setMouseOverAxisIndex
          this.tooltipState.hover = { active: true, index: props.activeIndex, dataKey: undefined, coordinate: props.activeCoordinate }
          this.tooltipState.keyboard = { ...this.tooltipState.keyboard, active: false }
        } else this.tooltipState.hover = { ...this.tooltipState.hover, active: false }
        this.render()
      })
    }
    w.addEventListener("mouseenter", move)
    w.addEventListener("mousemove", move)
    // mouseLeaveChart clears the active flags only; like recharts, a pending
    // (rAF-throttled) move still applies afterwards.
    w.addEventListener("mouseleave", () => {
      this.tooltipState.hover = { ...this.tooltipState.hover, active: false }
      if (this.tooltipState.item) this.tooltipState.item = { ...this.tooltipState.item, active: false }
      this.render()
    })
    // touchEventAction (axis tooltips): the first touch, rAF-throttled,
    // activates the axis hover like a mouse move; outside the plot area
    // nothing changes.
    let touchRaf = null
    let latestTouch = null
    w.addEventListener(
      "touchmove",
      (e) => {
        if (!e.touches || e.touches.length === 0) return
        latestTouch = pointer(e.touches[0])
        if (touchRaf !== null) cancelAnimationFrame(touchRaf)
        touchRaf = requestAnimationFrame(() => {
          touchRaf = null
          if (!latestTouch || !this.state) return
          const props = this.activePropsFromPointer(this.state, latestTouch)
          if (props?.activeIndex == null) return
          this.tooltipState.hover = { active: true, index: props.activeIndex, dataKey: undefined, coordinate: props.activeCoordinate }
          this.tooltipState.keyboard = { ...this.tooltipState.keyboard, active: false }
          this.render()
        })
      },
      { passive: true },
    )
    // Bar rectangles' item events (useMouseEnterItemDispatch,
    // useMouseLeaveItemDispatch). React emulates mouseenter/mouseleave from
    // mouseout (with relatedTarget), or from mouseover when the pointer comes
    // from outside. setActiveMouseOverItemIndex also deactivates the
    // keyboard interaction, which an axis tooltip shows.
    const barRect = (el) => (el instanceof Element ? el.closest(".recharts-bar-rectangle") : null)
    const enterBar = (g) => {
      const bar = g.closest(".recharts-bar")
      const item = this.items.find((i) => i.id === bar?.id)
      if (!item) return
      const index = [...g.parentNode.children].filter((c) => c.classList.contains("recharts-bar-rectangle")).indexOf(g)
      const entry = this.barEntries?.get(item.id)?.[index]
      this.tooltipState.item = { active: true, index: String(index), dataKey: item.props.dataKey, graphicalItemId: item.id, coordinate: entry?.tooltipPosition }
      this.tooltipState.keyboard = { ...this.tooltipState.keyboard, active: false }
      this.render()
    }
    const leaveBar = () => {
      if (!this.tooltipState.item?.active) return
      this.tooltipState.item = { ...this.tooltipState.item, active: false }
      this.render()
    }
    w.addEventListener("mouseout", (e) => {
      const from = barRect(e.target)
      const to = w.contains(e.relatedTarget) ? barRect(e.relatedTarget) : null
      if (from === to) return
      if (from) leaveBar()
      if (to) enterBar(to)
    })
    w.addEventListener("mouseover", (e) => {
      if (e.relatedTarget && w.contains(e.relatedTarget)) return
      const to = barRect(e.target)
      if (to) enterBar(to)
    })
    // mouseClickAction -> setMouseClickAxisIndex (not throttled): the click
    // interaction is recorded and the keyboard interaction deactivated.
    w.addEventListener("click", (e) => {
      if (!this.state) return
      const props = this.activePropsFromPointer(this.state, pointer(e))
      if (props?.activeIndex == null) return
      this.tooltipState.click = { active: true, index: props.activeIndex, dataKey: undefined, coordinate: props.activeCoordinate }
      this.tooltipState.keyboard = { ...this.tooltipState.keyboard, active: false }
      this.render()
    })
    w.addEventListener("focusin", () => {
      if (!this.props.accessibilityLayer || !this.state) return
      const kb = this.tooltipState.keyboard
      if (kb.active) return
      if (kb.index == null) {
        this.tooltipState.keyboard = { active: true, index: "0", coordinate: this.coordinateForIndex(this.state, "0") }
        this.render()
      }
    })
    w.addEventListener("focusout", () => {
      if (!this.props.accessibilityLayer) return
      const kb = this.tooltipState.keyboard
      if (kb.active) {
        this.tooltipState.keyboard = { ...kb, active: false }
        this.render()
      }
    })
    w.addEventListener("keydown", (e) => {
      if (!this.props.accessibilityLayer || !this.state) return
      const key = e.key
      if (key !== "ArrowRight" && key !== "ArrowLeft" && key !== "Enter") return
      const kb = this.tooltipState.keyboard
      if (kb.index == null) return
      const currentIndex = Number(kb.index)
      if (!Number.isFinite(currentIndex) || currentIndex < 0) return
      const state = this.state
      const ticks = state.axis(state.tooltipAxisType(), this.tooltipProps?.axisId ?? 0).tooltipTicks
      if (key === "Enter") {
        this.tooltipState.keyboard = { active: !kb.active, index: kb.index, coordinate: this.coordinateForIndex(state, String(kb.index)) }
        this.render()
        return
      }
      const leftToRight = state.layout === "horizontal" ? !this.xAxes.some((a) => a.settings.reversed) : state.layout === "centric" || state.layout === "radial"
      const movement = key === "ArrowRight" ? 1 : -1
      const nextIndex = currentIndex + movement * (leftToRight ? 1 : -1)
      if (ticks == null || nextIndex >= ticks.length || nextIndex < 0) return
      this.tooltipState.keyboard = { active: true, index: String(nextIndex), coordinate: this.coordinateForIndex(state, String(nextIndex)) }
      this.render()
    })
    this.onDocKeyDown = (e) => {
      if (e.key === "Escape" && this.tooltipProps) {
        const active = this.state ? this.activeTooltip(this.state) : null
        this.tooltipState.dismissed = true
        this.tooltipState.dismissedAt = { x: active?.coordinate?.x ?? 0, y: active?.coordinate?.y ?? 0 }
        this.render()
      }
    }
    document.addEventListener("keydown", this.onDocKeyDown)
  }

  activePropsFromPointer(state, chartEvent) {
    const offset = state.offset
    if (!(chartEvent.relativeX >= offset.left && chartEvent.relativeX <= offset.left + offset.width && chartEvent.relativeY >= offset.top && chartEvent.relativeY <= offset.top + offset.height)) return undefined
    const axisType = state.tooltipAxisType()
    const axis = state.axis(axisType, this.tooltipProps?.axisId ?? 0)
    const range = axis.settings.reversed ? [axis.range[1], axis.range[0]] : axis.range
    const ticks = axis.tooltipTicks
    if (!ticks || !range) return undefined
    const ordered = [...ticks].sort((a, b) => a.coordinate - b.coordinate)
    const pos = state.layout === "horizontal" ? chartEvent.relativeX : chartEvent.relativeY
    const activeIndex = R.calculateActiveTickIndex(pos, ordered, ticks, axisType, range)
    const activeCoordinate = R.getActiveCartesianCoordinate(state.layout, ticks, activeIndex, chartEvent)
    return { activeIndex: String(activeIndex), activeCoordinate }
  }

  // --- frames ------------------------------------------------------------------------------

  scheduleFrame() {
    if (this.frameRequested) return
    this.frameRequested = true
    requestAnimationFrame((now) => {
      this.frameRequested = false
      let more = false
      for (const rec of this.anims.values()) if (rec.anim && rec.anim.tick(now)) more = true
      if (this.updateLineLengths()) more = true
      this.render()
      if (more) this.scheduleFrame()
    })
  }

  destroy() {
    document.removeEventListener("keydown", this.onDocKeyDown)
    for (const rec of this.anims.values()) rec.anim?.stop()
  }
}

registerChart("AreaChart", CartesianChart)
registerChart("BarChart", CartesianChart)
registerChart("LineChart", CartesianChart)
registerChart("ComposedChart", CartesianChart)

function getLegendItemColor(stroke, fill) {
  return stroke && stroke !== "none" ? stroke : fill
}

function labelEntryFromPoint(point) {
  const viewBox = { x: point.x ?? 0, y: point.y ?? 0, width: 0, lowerWidth: 0, upperWidth: 0, height: 0 }
  return { ...viewBox, value: point.value, payload: point.payload, parentViewBox: undefined, viewBox, fill: undefined }
}

function getTickLineCoord(data, x, y, width, height, orientation, tickSize, mirror, tickMargin) {
  let x1, x2, y1, y2, tx, ty
  const sign = mirror ? -1 : 1
  const finalTickSize = data.tickSize || tickSize
  const tickCoord = isNum(data.tickCoord) ? data.tickCoord : data.coordinate
  switch (orientation) {
    case "top":
      x1 = x2 = data.coordinate
      y2 = y + +!mirror * height
      y1 = y2 - sign * finalTickSize
      ty = y1 - sign * tickMargin
      tx = tickCoord
      break
    case "left":
      y1 = y2 = data.coordinate
      x2 = x + +!mirror * width
      x1 = x2 - sign * finalTickSize
      tx = x1 - sign * tickMargin
      ty = tickCoord
      break
    case "right":
      y1 = y2 = data.coordinate
      x2 = x + +mirror * width
      x1 = x2 + sign * finalTickSize
      tx = x1 + sign * tickMargin
      ty = tickCoord
      break
    default:
      x1 = x2 = data.coordinate
      y2 = y + +mirror * height
      y1 = y2 + sign * finalTickSize
      ty = y1 + sign * tickMargin
      tx = tickCoord
      break
  }
  return { line: { x1, y1, x2, y2 }, tick: { x: tx, y: ty } }
}

function renderAxisLine({ x, y, width, height, orientation, mirror, axisLine, otherSvgProps }) {
  if (!axisLine) return null
  let props = { ...otherSvgProps, ...svgPropsFromUnknown(axisLine), fill: "none" }
  if (orientation === "top" || orientation === "bottom") {
    const needHeight = +((orientation === "top" && !mirror) || (orientation === "bottom" && mirror))
    props = { ...props, x1: x, y1: y + needHeight * height, x2: x + width, y2: y + needHeight * height }
  } else {
    const needWidth = +((orientation === "left" && !mirror) || (orientation === "right" && mirror))
    props = { ...props, x1: x + needWidth * width, y1: y, x2: x + needWidth * width, y2: y + height }
  }
  return h("line", { ...props, className: clsx("recharts-cartesian-axis-line", axisLine?.className) })
}

// Area's animation clip (ClipRect).
function clipRect({ alpha, layout, points, baseLine, strokeWidth }) {
  if (layout === "vertical") {
    const startY = points[0]?.y
    const endY = points[points.length - 1]?.y
    if (!wellBehaved(startY) || !wellBehaved(endY)) return null
    const height = alpha * Math.abs(startY - endY)
    let maxX = Math.max(...points.map((e) => e.x || 0))
    if (isNum(baseLine)) maxX = Math.max(baseLine, maxX)
    else if (baseLine && Array.isArray(baseLine) && baseLine.length) maxX = Math.max(...baseLine.map((e) => e.x || 0), maxX)
    if (isNum(maxX)) return h("rect", { x: 0, y: startY < endY ? startY : startY - height, width: maxX + (strokeWidth ? parseInt(`${strokeWidth}`, 10) : 1), height: Math.floor(height) })
    return null
  }
  const startX = points[0]?.x
  const endX = points[points.length - 1]?.x
  if (!wellBehaved(startX) || !wellBehaved(endX)) return null
  const width = alpha * Math.abs(startX - endX)
  let maxY = Math.max(...points.map((e) => e.y || 0))
  if (isNum(baseLine)) maxY = Math.max(baseLine, maxY)
  else if (baseLine && Array.isArray(baseLine) && baseLine.length) maxY = Math.max(...baseLine.map((e) => e.y || 0), maxY)
  if (isNum(maxY)) return h("rect", { x: startX < endX ? startX : startX - width, y: 0, width, height: Math.floor(maxY + (strokeWidth ? parseInt(`${strokeWidth}`, 10) : 1)) })
  return null
}

// Line getStrokeDasharray.
function getStrokeDasharray(length, totalLength, lines) {
  const lineLength = lines.reduce((pre, next) => pre + next, 0)
  if (!lineLength) return `${length}px ${totalLength}px`
  const count = Math.floor(length / lineLength)
  const remainLength = length % lineLength
  let remainLines = []
  for (let i = 0, sum = 0; i < lines.length; sum += lines[i] ?? 0, ++i) {
    const lineValue = lines[i]
    if (lineValue != null && sum + lineValue > remainLength) {
      remainLines = [...lines.slice(0, i), remainLength - sum]
      break
    }
  }
  const emptyLines = remainLines.length % 2 === 0 ? [0, totalLength] : [totalLength]
  const unit = lines.length % 2 !== 0 ? [...lines, 0] : lines
  const repeated = []
  for (let i = 0; i < count; ++i) repeated.push(...unit)
  return [...repeated, ...remainLines, ...emptyLines].map((line) => `${line}px`).join(", ")
}

// Bar shape option: undefined (Rectangle), a plain object (props merged
// over), or a function shape from a spec ($fn).
function barShape(option, props) {
  if (typeof option === "function") return option(props, props.index)
  if (option && typeof option === "object" && !Array.isArray(option)) return Rectangle({ ...props, ...option })
  return Rectangle(props)
}

// LabelList + Label (cartesian).
function renderLabelList(listProps, entries, state) {
  const { valueAccessor, dataKey, clockWise, id, textBreakAll, zIndex, ...others } = listProps
  const accessor = (entry) => {
    const val = Array.isArray(entry.value) ? entry.value[entry.value.length - 1] : entry.value
    return val == null || typeof val === "string" || typeof val === "number" || typeof val === "boolean" ? val : undefined
  }
  const labels = entries.map((entry, index) => {
    const value = isNullish(dataKey) ? accessor(entry, index) : R.getValueByDataKey(entry.payload, dataKey)
    const idProps = isNullish(id) ? {} : { id: `${id}-${index}` }
    return renderLabel({ ...svgProps(entry), ...others, ...idProps, fill: listProps.fill ?? entry.fill, parentViewBox: entry.parentViewBox, value, textBreakAll, viewBox: entry.viewBox, index, zIndex: 0 }, state)
  })
  return Layer({ className: "recharts-label-list" }, labels)
}

function renderLabel(outerProps, state) {
  const props = withDefaults(outerProps, LABEL_DEFAULTS)
  const { viewBox: vbProps, parentViewBox, position, value, children, content, className = "", textBreakAll, formatter } = props
  const viewBox = vbProps == null ? undefined : cartesianToTrapezoid(vbProps)
  if (!viewBox || (isNullish(value) && isNullish(children) && typeof content !== "function")) return null
  let label = isNullish(children) ? value : children
  if (typeof formatter === "function") label = formatter(label)
  const attrs = svgProps(props)
  const r = getCartesianPosition({ viewBox, position, offset: props.offset, parentViewBox, clamp: true })
  const positionAttrs = { x: r.x, y: r.y, textAnchor: r.horizontalAnchor, verticalAnchor: r.verticalAnchor, ...(r.width !== undefined ? { width: r.width } : {}), ...(r.height !== undefined ? { height: r.height } : {}) }
  return Text({ className: clsx("recharts-label", className), ...attrs, ...positionAttrs, textAnchor: ["start", "middle", "end", "inherit"].includes(attrs.textAnchor) ? attrs.textAnchor : positionAttrs.textAnchor, breakAll: textBreakAll, children: label })
}

function cartesianToTrapezoid(box) {
  if (!box) return undefined
  return { x: box.x, y: box.y, upperWidth: "upperWidth" in box ? box.upperWidth : box.width, lowerWidth: "lowerWidth" in box ? box.lowerWidth : box.width, width: box.width, height: box.height }
}

/** recharts cartesian/getCartesianPosition. */
export function getCartesianPosition({ viewBox, position, offset = 0, parentViewBox, clamp }) {
  const { x, y, height, upperWidth, lowerWidth } = cartesianToTrapezoid(viewBox)
  const upperX = x
  const lowerX = x + (upperWidth - lowerWidth) / 2
  const middleX = (upperX + lowerX) / 2
  const midHeightWidth = (upperWidth + lowerWidth) / 2
  const centerX = upperX + upperWidth / 2
  const verticalSign = height >= 0 ? 1 : -1
  const verticalOffset = verticalSign * offset
  const verticalEnd = verticalSign > 0 ? "end" : "start"
  const verticalStart = verticalSign > 0 ? "start" : "end"
  const horizontalSign = upperWidth >= 0 ? 1 : -1
  const horizontalOffset = horizontalSign * offset
  const horizontalEnd = horizontalSign > 0 ? "end" : "start"
  const horizontalStart = horizontalSign > 0 ? "start" : "end"
  if (position === "top") {
    const result = { x: upperX + upperWidth / 2, y: y - verticalOffset, horizontalAnchor: "middle", verticalAnchor: verticalEnd }
    if (clamp && parentViewBox) {
      result.height = Math.max(y - parentViewBox.y, 0)
      result.width = upperWidth
    }
    return result
  }
  if (position === "bottom") {
    const result = { x: lowerX + lowerWidth / 2, y: y + height + verticalOffset, horizontalAnchor: "middle", verticalAnchor: verticalStart }
    if (clamp && parentViewBox) {
      result.height = Math.max(parentViewBox.y + parentViewBox.height - (y + height), 0)
      result.width = lowerWidth
    }
    return result
  }
  if (position === "left") {
    const result = { x: middleX - horizontalOffset, y: y + height / 2, horizontalAnchor: horizontalEnd, verticalAnchor: "middle" }
    if (clamp && parentViewBox) {
      result.width = Math.max(result.x - parentViewBox.x, 0)
      result.height = height
    }
    return result
  }
  if (position === "right") {
    const result = { x: middleX + midHeightWidth + horizontalOffset, y: y + height / 2, horizontalAnchor: horizontalStart, verticalAnchor: "middle" }
    if (clamp && parentViewBox) {
      result.width = Math.max(parentViewBox.x + parentViewBox.width - result.x, 0)
      result.height = height
    }
    return result
  }
  const sizeAttrs = clamp && parentViewBox ? { width: midHeightWidth, height } : {}
  if (position === "insideLeft") return { x: middleX + horizontalOffset, y: y + height / 2, horizontalAnchor: horizontalStart, verticalAnchor: "middle", ...sizeAttrs }
  if (position === "insideRight") return { x: middleX + midHeightWidth - horizontalOffset, y: y + height / 2, horizontalAnchor: horizontalEnd, verticalAnchor: "middle", ...sizeAttrs }
  if (position === "insideTop") return { x: upperX + upperWidth / 2, y: y + verticalOffset, horizontalAnchor: "middle", verticalAnchor: verticalStart, ...sizeAttrs }
  if (position === "insideBottom") return { x: lowerX + lowerWidth / 2, y: y + height - verticalOffset, horizontalAnchor: "middle", verticalAnchor: verticalEnd, ...sizeAttrs }
  if (position === "insideTopLeft") return { x: upperX + horizontalOffset, y: y + verticalOffset, horizontalAnchor: horizontalStart, verticalAnchor: verticalStart, ...sizeAttrs }
  if (position === "insideTopRight") return { x: upperX + upperWidth - horizontalOffset, y: y + verticalOffset, horizontalAnchor: horizontalEnd, verticalAnchor: verticalStart, ...sizeAttrs }
  if (position === "insideBottomLeft") return { x: lowerX + horizontalOffset, y: y + height - verticalOffset, horizontalAnchor: horizontalStart, verticalAnchor: verticalEnd, ...sizeAttrs }
  if (position === "insideBottomRight") return { x: lowerX + lowerWidth - horizontalOffset, y: y + height - verticalOffset, horizontalAnchor: horizontalEnd, verticalAnchor: verticalEnd, ...sizeAttrs }
  if (!!position && typeof position === "object" && (isNum(position.x) || R.isPercent(position.x)) && (isNum(position.y) || R.isPercent(position.y))) {
    return { x: x + R.getPercentValue(position.x, midHeightWidth), y: y + R.getPercentValue(position.y, height), horizontalAnchor: "end", verticalAnchor: "end", ...sizeAttrs }
  }
  return { x: centerX, y: y + height / 2, horizontalAnchor: "middle", verticalAnchor: "middle", ...sizeAttrs }
}

// Renders a spec subtree of raw SVG (user <defs>): {type, props, children}.
// Elements a chart's <defs> may hold (gradients, patterns, clip paths,
// masks, filters and the shapes inside them). Anything else (script,
// foreignObject, a, ...) is dropped, as are event handler attributes and
// javascript: URLs, which React would never let run either.
const DEFS_TAGS = new Set("defs linearGradient radialGradient stop pattern clipPath mask filter symbol g rect circle ellipse line path polygon polyline text tspan title desc feBlend feColorMatrix feComponentTransfer feComposite feConvolveMatrix feDiffuseLighting feDisplacementMap feDistantLight feDropShadow feFlood feFuncA feFuncB feFuncG feFuncR feGaussianBlur feImage feMerge feMergeNode feMorphology feOffset fePointLight feSpecularLighting feSpotLight feTile feTurbulence".split(" "))

function renderSpecSvg(node) {
  if (typeof node === "string") return document.createTextNode(node)
  if (!node || !DEFS_TAGS.has(node.type)) return document.createTextNode("")
  const props = {}
  for (const [k, v] of Object.entries(node.props ?? {})) {
    if (/^on/i.test(k)) continue
    if (/href$/i.test(k) && /^javascript:/i.test(String(v).replace(/[\u0000-\u0020]/g, ""))) continue
    props[k] = v
  }
  const el = h(node.type, props)
  for (const c of node.children ?? []) el.appendChild(renderSpecSvg(c))
  return el
}

// --- shadcn ChartTooltipContent / ChartLegendContent ----------------------------------------

function getPayloadConfigFromPayload(config, payload, key) {
  if (typeof payload !== "object" || payload === null) return undefined
  const payloadPayload = "payload" in payload && typeof payload.payload === "object" && payload.payload !== null ? payload.payload : undefined
  let configLabelKey = key
  if (key in payload && typeof payload[key] === "string") configLabelKey = payload[key]
  else if (payloadPayload && key in payloadPayload && typeof payloadPayload[key] === "string") configLabelKey = payloadPayload[key]
  return configLabelKey in config ? config[configLabelKey] : config[key]
}

function iconFor(icons, key) {
  const tpl = icons?.get(key)
  if (!tpl) return null
  return tpl.content.firstElementChild.cloneNode(true)
}

// Classes are computed by the Go side (cn) and passed as props.classes.
function renderTooltipContent(props, { active, payload, label, config, icons }) {
  const { indicator = "dot", hideLabel = false, hideIndicator = false, labelFormatter, formatter, color, nameKey, labelKey, classes } = props
  if (!active || !payload?.length) return null
  let tooltipLabel = null
  if (!hideLabel && payload.length) {
    const [item] = payload
    const key = `${labelKey ?? item?.dataKey ?? item?.name ?? "value"}`
    const itemConfig = getPayloadConfigFromPayload(config, item, key)
    const value = !labelKey && typeof label === "string" ? (config[label]?.label ?? label) : itemConfig?.label
    if (labelFormatter) tooltipLabel = h("div", { className: classes.label }, labelFormatter(value, payload))
    else if (value) tooltipLabel = h("div", { className: classes.label }, value)
  }
  const nestLabel = payload.length === 1 && indicator !== "dot"
  const c = nestLabel ? classes.nest : classes.flat
  return h(
    "div",
    { className: classes.root },
    !nestLabel ? tooltipLabel : null,
    h(
      "div",
      { className: "grid gap-1.5" },
      payload
        .filter((item) => item.type !== "none")
        .map((item, index) => {
          const key = `${nameKey ?? item.name ?? item.dataKey ?? "value"}`
          const itemConfig = getPayloadConfigFromPayload(config, item, key)
          const indicatorColor = color ?? item.payload?.fill ?? item.color
          let inner
          if (formatter && item?.value !== undefined && item.name) inner = [formatter(item.value, item.name, item, index, item.payload)]
          else {
            const icon = itemConfig?.icon ? iconFor(icons, itemConfig.key) : null
            inner = [
              icon ?? (!hideIndicator ? h("div", { className: c.indicator, style: { "--color-bg": indicatorColor, "--color-border": indicatorColor } }) : null),
              h(
                "div",
                { className: c.value },
                h("div", { className: "grid gap-1.5" }, nestLabel ? tooltipLabel : null, h("span", { className: "text-muted-foreground" }, itemConfig?.label ?? item.name)),
                item.value != null ? h("span", { className: "font-mono font-medium text-foreground tabular-nums" }, typeof item.value === "number" ? item.value.toLocaleString() : String(item.value)) : null,
              ),
            ]
          }
          return h("div", { className: classes.item }, ...inner)
        }),
    ),
  )
}

function renderLegendContent(props, { payload, verticalAlign, config, icons }) {
  const { hideIcon = false, nameKey, classes } = props
  if (!payload?.length) return null
  return h(
    "div",
    { className: verticalAlign === "top" ? classes.rootTop : classes.rootBottom },
    payload
      .filter((item) => item.type !== "none")
      .map((item) => {
        const key = `${nameKey ?? item.dataKey ?? "value"}`
        const itemConfig = getPayloadConfigFromPayload(config, item, key)
        const icon = itemConfig?.icon && !hideIcon ? iconFor(icons, itemConfig.key) : null
        return h("div", { className: classes.item }, icon ?? h("div", { className: "h-2 w-2 shrink-0 rounded-[2px]", style: { backgroundColor: item.color } }), itemConfig?.label)
      }),
  )
}

// --- spec functions for shapes and dots ----------------------------------------------------

// Bar shape: index === args[0] ? <Rectangle {...props} {...args[1]}/> : <Rectangle {...props}/>
// ({index, ...props} drops the index before rendering the Rectangle).
registerFn("activeIndexShape", (args) => (props) => {
  const { index, ...rest } = props
  if (index === args[0]) return Rectangle({ ...rest, ...resolveFields(args[1] ?? {}, rest.payload) })
  return Rectangle(rest)
})

// Line dot: <Dot r cx cy {...args[0]} /> with $field values read from the payload.
registerFn("dot", (args) => (props) => Dot({ cx: props.cx, cy: props.cy, ...resolveFields(args[0] ?? {}, props.payload) }))

// Line dot: an icon (template args[0]) of size args[1] centered on the point,
// with extra attributes args[2]: <Icon x={cx - r/2} y={cy - r/2} width={r} height={r} .../>.
registerFn("iconDot", (args, ctx) => (props) => {
  const { cx, cy } = props
  if (cx == null || cy == null) return null
  const tpl = ctx.icons?.get(args[0])
  if (!tpl) return null
  const r = args[1]
  const svg = tpl.content.firstElementChild.cloneNode(true)
  setProps(svg, { x: cx - r / 2, y: cy - r / 2, width: r, height: r, ...resolveFields(args[2] ?? {}, props.payload) })
  return svg
})

// --- ChartContainer behavior -------------------------------------------------------------------

/** Allocates React-useId-like ids per page (the container uses _r_0_ by default). */
class IdScope {
  constructor(base) {
    this.base = base
  }
  item(n) {
    return reactId(this.base + n)
  }
}

let pageIdBase = 0
const payloads = new WeakMap()

define("chart", {
  selector: '[data-slot="chart"]',
  mount(root) {
    // The spec <script> is removed once read: React renders no script in
    // the container, and hide-others (aria-hidden) keeps every <script>
    // and its ancestors visible, so a leftover script would change which
    // elements a modal (dialog, select, ...) marks aria-hidden. The
    // payload is kept for remounts.
    let payload = payloads.get(root)
    const script = root.querySelector(":scope > script[data-chart-spec]")
    if (script) {
      try {
        payload = JSON.parse(script.textContent)
      } catch (e) {
        console.error("[shadcn] chart: invalid spec", e)
        return
      }
      payloads.set(root, payload)
      script.remove()
    }
    if (!payload) return
    // The <style> needs the page's CSP nonce; recreate it with the nonce
    // when the server could not set one.
    const style = root.querySelector(":scope > style")
    if (style && !style.nonce) {
      const nonce = getNonce()
      if (nonce) {
        const s = document.createElement("style")
        s.setAttribute("nonce", nonce)
        s.textContent = style.textContent
        style.replaceWith(s)
      }
    }
    // Several charts on one page: renumber duplicate generated ids.
    const base = pageIdBase
    if (payload.generatedId) {
      const id = root.getAttribute("data-chart")
      const dup = [...document.querySelectorAll(`[data-slot="chart"][data-chart="${CSS.escape(id)}"]`)].filter((el) => el !== root && el.chart)
      if (dup.length || base > 0) {
        const fresh = `chart-${reactId(base)}`
        const st = root.querySelector(":scope > style")
        if (st) st.textContent = st.textContent.split(`[data-chart=${id}]`).join(`[data-chart=${fresh}]`)
        root.setAttribute("data-chart", fresh)
      }
    }
    const config = {}
    for (const item of payload.config ?? []) config[item.key] = item
    const icons = new Map()
    for (const tpl of root.querySelectorAll(":scope > template[data-chart-icon]")) icons.set(tpl.getAttribute("data-chart-icon"), tpl)
    const ids = new IdScope(base)
    const Instance = CHART_TYPES[payload.spec?.type]
    if (!Instance) {
      console.error(`[shadcn] chart: unknown chart type ${payload.spec?.type}`)
      return
    }
    const responsive = root.querySelector(":scope > .recharts-responsive-container")
    const inner = responsive?.firstElementChild
    if (!responsive || !inner) return
    pageIdBase += 1 + countItems(payload.spec)
    const ctx = { config, icons }
    let chart = new Instance(root, payload.spec, ctx, ids)
    chart.mountSkeleton(inner)
    const initial = payload.initialDimension ?? { width: 320, height: 200 }
    let size = { width: initial.width, height: initial.height }
    const setContainerSize = (w, hh) => {
      const width = Math.round(w)
      const height = Math.round(hh)
      if (size.width === width && size.height === height) return false
      size = { width, height }
      return true
    }
    // ResponsiveContainerContextProvider renders nothing for a non-positive
    // size: the chart unmounts (store, tooltip state, animations) and mounts
    // afresh, from the registration passes on, once the size is positive.
    const acceptable = () => size.width > 0 && size.height > 0
    const r = responsive.getBoundingClientRect()
    setContainerSize(r.width, r.height)
    chart.setSize(size.width, size.height)
    chart.mountSequence()
    let unmounted = !acceptable()
    const ro = new ResizeObserver((entries) => {
      const e = entries[0]
      if (!e) return
      if (!setContainerSize(e.contentRect.width, e.contentRect.height)) return
      if (!acceptable()) {
        unmounted = true
        chart.setSize(size.width, size.height)
        chart.render()
        return
      }
      if (unmounted) {
        unmounted = false
        const spec = chart.spec
        chart.destroy()
        chart.wrapper?.remove()
        chart = new Instance(root, spec, ctx, ids)
        chart.mountSkeleton(inner)
        chart.setSize(size.width, size.height)
        chart.mountSequence()
        return
      }
      chart.setSize(size.width, size.height)
      chart.commit()
    })
    ro.observe(responsive)
    root.chart = {
      get instance() {
        return chart
      },
      spec: () => structuredClone(chart.spec),
      setSpec(spec) {
        chart.setSpec(spec)
      },
      update(fn) {
        const spec = structuredClone(chart.spec)
        fn(spec)
        chart.setSpec(spec)
      },
    }
    return () => {
      ro.disconnect()
      chart.destroy()
      delete root.chart
    }
  },
})

function countItems(spec) {
  return (spec?.children ?? []).filter((c) => c.type === "Area" || c.type === "Bar" || c.type === "Line" || c.type === "Pie" || c.type === "Radar" || c.type === "RadialBar").length
}

