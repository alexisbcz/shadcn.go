// Bundle entry: shiki 3.23.0 (the version rehype-pretty-code 0.14.1 resolves
// to in ui/apps/v4), its JavaScript regex engine, and the grammars and themes
// the docs use. The packages are bundled unmodified.
import { createHighlighterCoreSync } from "@shikijs/core"
import { createJavaScriptRegexEngine } from "@shikijs/engine-javascript"

import astro from "@shikijs/langs/astro"
import css from "@shikijs/langs/css"
import diff from "@shikijs/langs/diff"
import html from "@shikijs/langs/html"
import javascript from "@shikijs/langs/javascript"
import json from "@shikijs/langs/json"
import jsx from "@shikijs/langs/jsx"
import markdown from "@shikijs/langs/markdown"
import postcss from "@shikijs/langs/postcss"
import regexp from "@shikijs/langs/regexp"
import shellscript from "@shikijs/langs/shellscript"
import toml from "@shikijs/langs/toml"
import tsx from "@shikijs/langs/tsx"
import typescript from "@shikijs/langs/typescript"

import githubLightDefault from "@shikijs/themes/github-light-default"
import vesper from "@shikijs/themes/vesper"

globalThis.Shiki = {
  createHighlighterCoreSync,
  createJavaScriptRegexEngine,
  langs: { astro, css, diff, html, javascript, json, jsx, markdown, postcss, regexp, shellscript, toml, tsx, typescript },
  themes: { "github-light-default": githubLightDefault, vesper },
}
