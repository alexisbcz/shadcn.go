// trees.mjs reduces a TSX source to the class-bearing syntax trees the Go
// tests replay (format: see extract.mjs). Shared by extract.mjs and
// upstream.mjs.
import path from "node:path"
import { createRequire } from "node:module"
import { fileURLToPath } from "node:url"

const here = path.dirname(fileURLToPath(import.meta.url))
const require = createRequire(path.resolve(here, "../../../../ui/apps/v4/package.json"))
const ts = require("typescript")

function calleeIdent(expr) {
  return ts.isIdentifier(expr) ? expr.text : null
}

function calleeRoot(expr) {
  let e = expr
  for (;;) {
    if (ts.isCallExpression(e)) e = e.expression
    else if (ts.isPropertyAccessExpression(e) || ts.isElementAccessExpression(e)) e = e.expression
    else break
  }
  return ts.isIdentifier(e) ? e.text : null
}

const ROOT_CALLEES = new Set(["cn", "cva", "mergeProps"])

function isRoot(node) {
  if (ts.isJsxAttribute(node)) return node.name.getText() === "className"
  if (ts.isCallExpression(node)) {
    const f = calleeIdent(node.expression)
    const r = calleeRoot(node.expression)
    return ROOT_CALLEES.has(f) || r === "cn" || r === "cva"
  }
  if (ts.isPropertyAssignment(node)) return node.name.getText() === "className"
  return false
}

function children(node) {
  const out = []
  ts.forEachChild(node, (c) => {
    out.push(c)
  })
  return out
}

function ser(node) {
  while (ts.isParenthesizedExpression(node)) node = node.expression
  if (ts.isStringLiteral(node)) return { k: "str", v: node.text }
  if (ts.isNoSubstitutionTemplateLiteral(node)) return { k: "str", v: node.text, t: 1 }
  if (ts.isTemplateExpression(node)) {
    const q = [[node.head.rawText ?? node.head.text, node.head.text]]
    const e = []
    for (const span of node.templateSpans) {
      e.push(ser(span.expression))
      q.push([span.literal.rawText ?? span.literal.text, span.literal.text])
    }
    return { k: "tpl", q, e }
  }
  if (ts.isCallExpression(node)) {
    return {
      k: "call",
      f: calleeIdent(node.expression),
      r: calleeRoot(node.expression),
      c: ser(node.expression),
      a: node.arguments.map(ser),
    }
  }
  if (ts.isObjectLiteralExpression(node)) return { k: "obj", p: node.properties.map(ser) }
  if (ts.isPropertyAssignment(node)) {
    const out = { k: "pa", n: node.name.getText(), id: ts.isIdentifier(node.name) ? 1 : 0, v: ser(node.initializer) }
    // A string key is a literal prettier sorts too (clsx object syntax).
    if (ts.isStringLiteral(node.name) || ts.isNoSubstitutionTemplateLiteral(node.name)) out.nk = ser(node.name)
    return out
  }
  if (ts.isJsxAttribute(node)) {
    return { k: "attr", n: node.name.getText(), v: node.initializer ? ser(node.initializer) : null }
  }
  if (ts.isJsxExpression(node)) return { k: "jx", e: node.expression ? ser(node.expression) : null }
  if (ts.isJsxElement(node)) {
    return {
      k: "jsx",
      a: node.openingElement.attributes.properties.map(ser),
      c: node.children.map(ser).filter(Boolean),
    }
  }
  if (ts.isJsxSelfClosingElement(node)) {
    return { k: "jsx", a: node.attributes.properties.map(ser), c: [] }
  }
  if (ts.isJsxText(node)) {
    // Prettier rewraps JSX text; whitespace-only text is not content.
    const t = node.text.replace(/\s+/g, " ").trim()
    return t ? { k: "g", s: "JsxText" } : null
  }
  if (ts.isBinaryExpression(node)) {
    return { k: "bin", op: node.operatorToken.getText(), l: ser(node.left), r: ser(node.right) }
  }
  if (ts.isIdentifier(node)) return { k: "id", v: node.text }
  const c = children(node).map(ser).filter(Boolean)
  const out = { k: "g", s: ts.SyntaxKind[node.kind] }
  if (c.length) out.c = c
  return out
}

// stringsOutsideRoots lists, in document order, every string literal,
// template part and JSX text outside the class roots, except import/export
// specifiers and the IconPlaceholder library attributes (transformIcons
// rewrites those). extract.mjs requires these to be identical in a base
// source and its generated counterpart: then the roots hold everything the
// build changes.
export function stringsOutsideRoots(file, text) {
  const sf = ts.createSourceFile(file, text, ts.ScriptTarget.Latest, true, ts.ScriptKind.TSX)
  const out = []
  const visit = (n) => {
    if (isRoot(n) || ts.isImportDeclaration(n) || ts.isExportDeclaration(n)) return
    if (ts.isJsxAttribute(n) && /^(lucide|tabler|hugeicons|phosphor|remixicon)$/.test(n.name.getText())) return
    if (ts.isStringLiteral(n) || ts.isNoSubstitutionTemplateLiteral(n) || ts.isTemplateHead(n) || ts.isTemplateMiddle(n) || ts.isTemplateTail(n)) out.push(n.text)
    if (ts.isJsxText(n)) {
      const t = n.text.replace(/\s+/g, " ").trim()
      if (t) out.push(t)
    }
    ts.forEachChild(n, visit)
  }
  visit(sf)
  return out
}

export function rootsOf(file, text) {
  const sf = ts.createSourceFile(file, text, ts.ScriptTarget.Latest, true, ts.ScriptKind.TSX)
  const out = []
  const visit = (node) => {
    if (isRoot(node)) {
      const { line } = sf.getLineAndCharacterOfPosition(node.getStart(sf))
      out.push({ line: line + 1, n: ser(node) })
      return
    }
    ts.forEachChild(node, visit)
  }
  visit(sf)
  return out
}

