// collection.js: port of @radix-ui/react-collection 1.1.7 (createCollection,
// the "legacy" implementation every radix-ui 1.4.3 primitive uses).
//
// A collection is an element (the collection root: a tab list, a menu
// content, a toggle group, ...) whose items carry the attribute
// `data-radix-collection-item=""`, exactly the attribute Radix's
// Collection.ItemSlot renders (server-render it on every item). getItems()
// returns the items in document order, like Radix's useCollection():
// Radix queries `[data-radix-collection-item]` under the collection node and
// keeps only the items registered in its own scope; here an item belongs to
// the nearest enclosing collection root, so collections nest (a submenu
// inside a menu) without seeing each other's items.
//
//   import { createCollection } from "./collection.js"
//   const collection = createCollection(listEl)
//   collection.getItems()              // Element[] in document order
//   collection.setItemData(el, { focusable: true, active: false })
//   collection.getEntries()            // [{ element, ref: { current }, ...data }]
//   collection.destroy()               // when the root unmounts
//
// Several collections may share one root (Radix menus stack a Collection
// and a RovingFocusGroup collection on the same content): createCollection
// with the same root returns independent handles over the same items; pass
// `filter` to narrow one of them (e.g. to enabled items).
//
// Item data mirrors the props Radix passes to Collection.ItemSlot (id,
// focusable, active, disabled, textValue, ...). In the DOM most of it is
// derived from attributes, so data is optional: entries default to {}.

export const ITEM_DATA_ATTR = "data-radix-collection-item"
const ITEM_SELECTOR = `[${ITEM_DATA_ATTR}]`

const roots = new WeakMap() // root element -> number of live handles
const itemData = new WeakMap() // item element -> data

/** Returns the collection root that owns item: its nearest registered ancestor. */
export function collectionOf(item) {
  for (let el = item.parentElement; el; el = el.parentElement) {
    if (roots.has(el)) return el
  }
  return null
}

/**
 * Creates a collection over root.
 * @param {Element} root
 * @param {{ filter?: (el: Element) => boolean }} [options]
 */
export function createCollection(root, options = {}) {
  const { filter } = options
  roots.set(root, (roots.get(root) ?? 0) + 1)
  let destroyed = false

  const getItems = () => {
    if (destroyed || !root.isConnected) return []
    const out = []
    for (const el of root.querySelectorAll(ITEM_SELECTOR)) {
      if (collectionOf(el) !== root) continue
      if (filter && !filter(el)) continue
      out.push(el)
    }
    return out
  }

  return {
    root,
    /** The items, in document order. */
    getItems,
    /** The items with their data, like Radix's getItems() entries. */
    getEntries() {
      return getItems().map((element) => ({ ...(itemData.get(element) ?? {}), element, ref: { current: element } }))
    },
    /** Marks el as an item (adds the attribute) and records its data. */
    addItem(el, data) {
      if (!el.hasAttribute(ITEM_DATA_ATTR)) el.setAttribute(ITEM_DATA_ATTR, "")
      if (data) itemData.set(el, { ...data })
    },
    /** Records data for an item (merged with what is there). */
    setItemData(el, data) {
      itemData.set(el, { ...(itemData.get(el) ?? {}), ...data })
    },
    /** Returns the data recorded for an item. */
    getItemData(el) {
      return itemData.get(el) ?? {}
    },
    /** Whether el is one of this collection's items. */
    has(el) {
      return !!el && el.matches?.(ITEM_SELECTOR) && collectionOf(el) === root && (!filter || filter(el))
    },
    /**
     * Returns this collection's item containing node (or node itself), if
     * any. Items of nested collections are skipped, so an event from inside
     * a nested group still reaches the outer item, as React events bubble
     * through every ancestor item's handlers.
     */
    itemFor(node) {
      for (let el = node?.nodeType === 1 ? node : node?.parentElement; el && el !== root; el = el.parentElement) {
        if (el.matches(ITEM_SELECTOR) && collectionOf(el) === root) return !filter || filter(el) ? el : null
      }
      return null
    },
    destroy() {
      if (destroyed) return
      destroyed = true
      const n = (roots.get(root) ?? 1) - 1
      if (n <= 0) roots.delete(root)
      else roots.set(root, n)
    },
  }
}
