// A port of react-textarea-autosize 8.5.9 for input-group-custom: the
// textarea's height is set (with !important) to fit its value, or its
// placeholder, measured on a hidden textarea appended to <body>, on mount,
// on input, on window resize, when fonts finish loading and after a reset
// of its form.
import { define } from "./core.js"

const HIDDEN_TEXTAREA_STYLE = {
  "min-height": "0",
  "max-height": "none",
  height: "0",
  visibility: "hidden",
  overflow: "hidden",
  position: "absolute",
  "z-index": "-1000",
  top: "0",
  right: "0",
  display: "block",
}
const forceHiddenStyles = (node) => {
  Object.keys(HIDDEN_TEXTAREA_STYLE).forEach((key) => node.style.setProperty(key, HIDDEN_TEXTAREA_STYLE[key], "important"))
}

let hiddenTextarea = null
const getHeight = (node, sizingData) => {
  const height = node.scrollHeight
  if (sizingData.sizingStyle.boxSizing === "border-box") return height + sizingData.borderSize
  return height - sizingData.paddingSize
}

function calculateNodeHeight(sizingData, value, minRows = 1, maxRows = Infinity) {
  if (!hiddenTextarea) {
    hiddenTextarea = document.createElement("textarea")
    hiddenTextarea.setAttribute("tabindex", "-1")
    hiddenTextarea.setAttribute("aria-hidden", "true")
    forceHiddenStyles(hiddenTextarea)
  }
  if (hiddenTextarea.parentNode === null) document.body.appendChild(hiddenTextarea)
  const { paddingSize, borderSize, sizingStyle } = sizingData
  const { boxSizing } = sizingStyle
  Object.keys(sizingStyle).forEach((key) => {
    hiddenTextarea.style[key] = sizingStyle[key]
  })
  forceHiddenStyles(hiddenTextarea)
  hiddenTextarea.value = value
  let height = getHeight(hiddenTextarea, sizingData)
  hiddenTextarea.value = value
  height = getHeight(hiddenTextarea, sizingData)
  hiddenTextarea.value = "x"
  const rowHeight = hiddenTextarea.scrollHeight - paddingSize
  let minHeight = rowHeight * minRows
  if (boxSizing === "border-box") minHeight = minHeight + paddingSize + borderSize
  height = Math.max(minHeight, height)
  let maxHeight = rowHeight * maxRows
  if (boxSizing === "border-box") maxHeight = maxHeight + paddingSize + borderSize
  height = Math.min(maxHeight, height)
  return [height, rowHeight]
}

const SIZING_STYLE = ["borderBottomWidth", "borderLeftWidth", "borderRightWidth", "borderTopWidth", "boxSizing", "fontFamily", "fontSize", "fontStyle", "fontWeight", "letterSpacing", "lineHeight", "paddingBottom", "paddingLeft", "paddingRight", "paddingTop", "tabSize", "textIndent", "textRendering", "textTransform", "width", "wordBreak", "wordSpacing", "scrollbarGutter"]

function getSizingData(node) {
  const style = window.getComputedStyle(node)
  if (style === null) return null
  const sizingStyle = {}
  for (const k of SIZING_STYLE) sizingStyle[k] = style[k]
  if (sizingStyle.boxSizing === "") return null
  const paddingSize = parseFloat(sizingStyle.paddingBottom) + parseFloat(sizingStyle.paddingTop)
  const borderSize = parseFloat(sizingStyle.borderBottomWidth) + parseFloat(sizingStyle.borderTopWidth)
  return { sizingStyle, paddingSize, borderSize }
}

define("example-textarea-autosize", {
  selector: '[data-slot="input-group"] > textarea[data-slot="input-group-control"].py-2\\.5',
  mount(node) {
    let lastHeight = 0
    const resize = () => {
      const sizing = getSizingData(node)
      if (!sizing) return
      const [height] = calculateNodeHeight(sizing, node.value || node.placeholder || "x")
      if (lastHeight !== height) {
        lastHeight = height
        node.style.setProperty("height", height + "px", "important")
      }
    }
    resize()
    const onReset = (ev) => {
      if (node.form !== ev.target) return
      const currentValue = node.value
      requestAnimationFrame(() => {
        if (currentValue !== node.value) resize()
      })
    }
    node.addEventListener("input", resize)
    document.body.addEventListener("reset", onReset)
    window.addEventListener("resize", resize)
    document.fonts?.addEventListener("loadingdone", resize)
    return () => {
      node.removeEventListener("input", resize)
      document.body.removeEventListener("reset", onReset)
      window.removeEventListener("resize", resize)
      document.fonts?.removeEventListener("loadingdone", resize)
    }
  },
})
