// Produces components/ui/runtime/js/recharts.js: the framework-free parts of
// recharts 3.8.0 (the exact version ui/apps/v4 locks) that the chart runtime
// reuses unmodified: d3-scale and d3-shape through victory-vendor 37.3.6,
// decimal.js-light (nice ticks), and the pure recharts modules that compute
// scales, ticks, stacks, bar positions, tooltip positions and easings. The
// React components themselves are ported by hand in chart.js.
//
//   node components/ui/chart/vendor/build.mjs
import fs from "node:fs"
import path from "node:path"
import { createRequire } from "node:module"

const repo = path.resolve(import.meta.dirname, "../../../..")
const esbuild = createRequire(path.join(repo, "parity/package.json"))("esbuild")
const app = path.join(repo, "ui/apps/v4")
const recharts = fs.realpathSync(path.join(app, "node_modules/recharts"))
const version = JSON.parse(fs.readFileSync(path.join(recharts, "package.json"), "utf8")).version
const r = (p) => JSON.stringify(path.join(recharts, "es6", p))

const entry = `
export { scaleBand, scaleLinear, scalePoint } from "victory-vendor/d3-scale"
export * from "victory-vendor/d3-shape"
export { getNiceTickValues, getTickValuesFixedDomain } from ${r("util/scale/getNiceTickValues.js")}
export { getTicks } from ${r("cartesian/getTicks.js")}
export { getStringSize } from ${r("util/DOMUtils.js")}
export * from ${r("util/ChartUtils.js")}
export * from ${r("util/DataUtils.js")}
export * from ${r("util/round.js")}
export * from ${r("util/isWellBehavedNumber.js")}
export * from ${r("util/isDomainSpecifiedByUser.js")}
export * from ${r("util/getActiveCoordinate.js")}
export * from ${r("util/cursor/getCursorRectangle.js")}
export * from ${r("util/tooltip/translate.js")}
export * from ${r("util/ReduceCSSCalc.js")}
export * from ${r("util/PolarUtils.js")}
export * from ${r("util/CartesianUtils.js")}
export { configEasing } from ${r("animation/easing.js")}
export { combineAllBarPositions } from ${r("state/selectors/combiners/combineAllBarPositions.js")}
export { combineBarSizeList } from ${r("state/selectors/combiners/combineBarSizeList.js")}
`

const result = await esbuild.build({
  stdin: { contents: entry, resolveDir: path.join(recharts, "es6"), sourcefile: "recharts-entry.js" },
  bundle: true,
  format: "esm",
  platform: "browser",
  target: "es2022",
  write: false,
  legalComments: "none",
  define: { "process.env.NODE_ENV": '"production"' },
  logLevel: "warning",
})
const header = `// recharts ${version} (MIT, Copyright (c) 2015-present recharts), d3-scale 4.0.2, d3-shape 3.2.0, d3-path 3.1.0,
// d3-array 3.2.4, d3-interpolate 3.0.1, d3-format 3.1.0, d3-color 3.1.0, internmap 2.0.3
// (ISC, Copyright Mike Bostock) through victory-vendor 37.3.6, decimal.js-light 2.5.1 (MIT, Copyright (c) 2017 Michael Mclaughlin),
// es-toolkit 1.45.1 (MIT), clsx 2.1.1 (MIT). The framework-free modules of recharts, bundled unmodified with
// esbuild by components/ui/chart/vendor/build.mjs; the React components are ported in components/ui/chart/chart.js.
`
const out = path.join(repo, "components/ui/runtime/js/recharts.js")
fs.writeFileSync(out, header + "\n" + result.outputFiles[0].text)
console.log(`${path.relative(repo, out)}: ${(fs.statSync(out).size / 1024).toFixed(0)} KiB`)
