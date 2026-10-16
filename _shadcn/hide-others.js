// hide-others.js: port of aria-hidden 1.2.6 (theKashey/aria-hidden), which
// Radix modal layers (Dialog, AlertDialog, modal Popover / Menu / Select)
// call on their content.
//
// It marks every element outside the given targets as aria-hidden="true"
// (plus a marker attribute, data-aria-hidden="true" by default), walking
// down from the parent (document.body) only into ancestors of the targets.
// Elements with aria-live and <script> elements are always kept. Calls nest
// with reference counts; the returned undo restores exactly what it changed.
//
// API
//
//   import { hideOthers, inertOthers, suppressOthers, supportsInert } from "./hide-others.js"
//
//   const undo = hideOthers(content)                 // or [el1, el2]
//   const undo = hideOthers(content, parentNode, "data-aria-hidden")
//   undo()
//
//   inertOthers(...)     // same, with the `inert` attribute (marker data-inert-ed)
//   suppressOthers(...)  // inert when supported, else aria-hidden (data-suppressed)

const getDefaultParent = (originalTarget) => {
  if (typeof document === "undefined") return null
  const sampleTarget = Array.isArray(originalTarget) ? originalTarget[0] : originalTarget
  return sampleTarget.ownerDocument.body
}

let counterMap = new WeakMap()
let uncontrolledNodes = new WeakMap()
let markerMap = {}
let lockCount = 0

const unwrapHost = (node) => node && (node.host || unwrapHost(node.parentNode))

const correctTargets = (parent, targets) =>
  targets
    .map((target) => {
      if (parent.contains(target)) return target
      const correctedTarget = unwrapHost(target)
      if (correctedTarget && parent.contains(correctedTarget)) return correctedTarget
      console.error("aria-hidden", target, "in not contained inside", parent, ". Doing nothing")
      return null
    })
    .filter((x) => Boolean(x))

const applyAttributeToOthers = (originalTarget, parentNode, markerName, controlAttribute) => {
  const targets = correctTargets(parentNode, Array.isArray(originalTarget) ? originalTarget : [originalTarget])
  if (!markerMap[markerName]) markerMap[markerName] = new WeakMap()
  const markerCounter = markerMap[markerName]
  const hiddenNodes = []
  const elementsToKeep = new Set()
  const elementsToStop = new Set(targets)
  const keep = (el) => {
    if (!el || elementsToKeep.has(el)) return
    elementsToKeep.add(el)
    keep(el.parentNode)
  }
  targets.forEach(keep)
  const deep = (parent) => {
    if (!parent || elementsToStop.has(parent)) return
    Array.prototype.forEach.call(parent.children, (node) => {
      if (elementsToKeep.has(node)) {
        deep(node)
      } else {
        try {
          const attr = node.getAttribute(controlAttribute)
          const alreadyHidden = attr !== null && attr !== "false"
          const counterValue = (counterMap.get(node) || 0) + 1
          const markerValue = (markerCounter.get(node) || 0) + 1
          counterMap.set(node, counterValue)
          markerCounter.set(node, markerValue)
          hiddenNodes.push(node)
          if (counterValue === 1 && alreadyHidden) uncontrolledNodes.set(node, true)
          if (markerValue === 1) node.setAttribute(markerName, "true")
          if (!alreadyHidden) node.setAttribute(controlAttribute, "true")
        } catch (e) {
          console.error("aria-hidden: cannot operate on ", node, e)
        }
      }
    })
  }
  deep(parentNode)
  elementsToKeep.clear()
  lockCount++
  return () => {
    hiddenNodes.forEach((node) => {
      const counterValue = counterMap.get(node) - 1
      const markerValue = markerCounter.get(node) - 1
      counterMap.set(node, counterValue)
      markerCounter.set(node, markerValue)
      if (!counterValue) {
        if (!uncontrolledNodes.has(node)) node.removeAttribute(controlAttribute)
        uncontrolledNodes.delete(node)
      }
      if (!markerValue) node.removeAttribute(markerName)
    })
    lockCount--
    if (!lockCount) {
      counterMap = new WeakMap()
      uncontrolledNodes = new WeakMap()
      markerMap = {}
    }
  }
}

/**
 * Marks everything except the targets as aria-hidden.
 * @param {Element | Element[]} originalTarget elements to keep
 * @param {Element} [parentNode] top element, defaults to document.body
 * @param {string} [markerName] marker attribute, default data-aria-hidden
 * @returns {() => void} undo
 */
export const hideOthers = (originalTarget, parentNode, markerName = "data-aria-hidden") => {
  const targets = Array.from(Array.isArray(originalTarget) ? originalTarget : [originalTarget])
  const activeParentNode = parentNode || getDefaultParent(originalTarget)
  if (!activeParentNode) return () => null
  // Do not hide aria-live elements (theKashey/aria-hidden#10) nor scripts.
  targets.push(...Array.from(activeParentNode.querySelectorAll("[aria-live], script")))
  return applyAttributeToOthers(targets, activeParentNode, markerName, "aria-hidden")
}

/** Like hideOthers, with the inert attribute. */
export const inertOthers = (originalTarget, parentNode, markerName = "data-inert-ed") => {
  const activeParentNode = parentNode || getDefaultParent(originalTarget)
  if (!activeParentNode) return () => null
  return applyAttributeToOthers(originalTarget, activeParentNode, markerName, "inert")
}

/** Reports whether the browser supports the inert attribute. */
export const supportsInert = () => typeof HTMLElement !== "undefined" && Object.prototype.hasOwnProperty.call(HTMLElement.prototype, "inert")

/** inertOthers when inert is supported, else hideOthers. */
export const suppressOthers = (originalTarget, parentNode, markerName = "data-suppressed") =>
  (supportsInert() ? inertOthers : hideOthers)(originalTarget, parentNode, markerName)
