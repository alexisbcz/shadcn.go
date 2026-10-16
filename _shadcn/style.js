// style.js: runtime <style> injection under the strict CSP.
//
// Go pages are served with a nonce-based Content-Security-Policy: a <style>
// element only applies when it carries the request's nonce. Libraries that
// inject styles at runtime upstream (react-remove-scroll, sonner, vaul, ...)
// use these helpers to do the same with the nonce of the runtime's own
// <script> tag.
//
// API
//
//   import { getNonce, injectStyle } from "./style.js"
//
//   getNonce()                       // the page's CSP nonce, or undefined
//   const remove = injectStyle(css)  // appends <style nonce=...>css</style>
//   remove()                         // to <head>; remove() takes it out

/** Returns the CSP nonce of the page's scripts, if any. */
export function getNonce() {
  for (const s of document.querySelectorAll("script")) {
    if (s.nonce) return s.nonce
  }
  return undefined
}

/**
 * Appends a <style> with css (and the page nonce) to <head>.
 * @param {string} css
 * @returns {() => void} removes it
 */
export function injectStyle(css) {
  const tag = document.createElement("style")
  tag.type = "text/css"
  const nonce = getNonce()
  if (nonce) tag.setAttribute("nonce", nonce)
  tag.appendChild(document.createTextNode(css))
  ;(document.head || document.getElementsByTagName("head")[0]).appendChild(tag)
  return () => tag.remove()
}
