// Proves that parity/docs/run.mjs detects and locates differences.
//
//   node parity/docs/stub-proof.mjs [--ref http://localhost:4000] [--jobs 2] [--case name,..]
//
// Serves reference pages through cmd/docs-stub (a reverse proxy of the
// reference) and runs the harness against it, once unmodified (must pass)
// and once per injected change (must fail, at the right place). Each case
// checks the report, not just the exit status:
//
//   unmodified       3 pages, every viewport and theme, all scenarios: 0 differences
//   class-removed    "tracking-tight" removed from the <h1> of one page: a DOM
//                    difference on that <h1> line and a pixel difference
//                    located at the <h1>
//   padding-1px      1px of left padding on one page's <h1> (CSS only, so the
//                    DOM is identical): a pixel difference located at the <h1>
//   anim-duration    the mobile menu popover's exit animation lasts 250ms
//                    instead of 100ms (CSS only): an animation-definition
//                    difference naming the duration and popover-content, in
//                    the scenario that closes it, while "initial" is identical
//
// Exit status 0 when every case behaves as expected.
import fs from "node:fs"
import os from "node:os"
import path from "node:path"
import { spawn, execFileSync } from "node:child_process"

const repo = path.resolve(import.meta.dirname, "../..")
const args = process.argv.slice(2)
const opt = (n, d) => (args.includes(n) ? args[args.indexOf(n) + 1] : d)
const REF = opt("--ref", "http://localhost:4000")
const JOBS = opt("--jobs", "2")
const outBase = path.join(import.meta.dirname, "out-stub")

const bin = path.join(os.tmpdir(), `docs-stub-${process.pid}`)
execFileSync("go", ["build", "-o", bin, "./cmd/docs-stub"], { cwd: repo, stdio: "inherit" })

function startStub(flags) {
  const proc = spawn(bin, ["-upstream", REF, "-addr", "127.0.0.1:0", ...flags], { stdio: ["ignore", "pipe", "pipe"] })
  let err = ""
  proc.stderr.on("data", (d) => (err += d))
  return new Promise((resolve, reject) => {
    let buf = ""
    proc.stdout.on("data", (d) => {
      buf += d
      const m = buf.match(/listening on (http:\/\/\S+)/)
      if (m) resolve({ url: m[1], proc, log: () => err })
    })
    proc.on("exit", (code) => reject(new Error(`docs-stub exited ${code}: ${err}`)))
  })
}

function runHarness(target, out, extra) {
  const r = spawn(process.execPath, [path.join(import.meta.dirname, "run.mjs"), "--ref", REF, "--target", target, "--out", out, "--jobs", JOBS, ...extra], { stdio: ["ignore", "pipe", "inherit"] })
  let text = ""
  r.stdout.on("data", (d) => { text += d; process.stdout.write(d) })
  return new Promise((resolve) => r.on("exit", (code) => resolve({ code, text })))
}

const reports = (out) => {
  const all = []
  for (const page of fs.existsSync(out) ? fs.readdirSync(out) : []) {
    const dir = path.join(out, page)
    if (!fs.statSync(dir).isDirectory()) continue
    for (const job of fs.readdirSync(dir)) {
      const f = path.join(dir, job, "report.json")
      if (fs.existsSync(f)) all.push(JSON.parse(fs.readFileSync(f, "utf8")))
    }
  }
  return all
}

const H1 = 'class="scroll-m-24 text-3xl font-semibold tracking-tight sm:text-3xl"'
const cases = [
  {
    name: "unmodified",
    stub: [],
    run: ["/docs/components/radix/button", "/docs/installation", "/"],
    expect: (rs) => {
      const bad = rs.filter((r) => !r.pass)
      const caps = rs.reduce((n, r) => n + r.captures.length, 0)
      return { ok: rs.length === 18 && bad.length === 0, detail: `${rs.length} jobs, ${caps} captures, ${bad.length} failing${bad.length ? ": " + bad.map((r) => `${r.path} ${r.theme}/${r.viewport} ${r.captures.filter((c) => !c.pass).map((c) => c.label)}`).join("; ") : ""}` }
    },
  },
  {
    name: "class-removed",
    stub: ["-paths", "/docs/components/radix/button", "-replace", `${H1}=>${H1.replace(" tracking-tight", "")}`],
    run: ["--no-scenarios", "/docs/components/radix/button"],
    expect: (rs) => {
      const checks = rs.map((r) => {
        const c = r.captures.find((x) => x.label === "initial")
        const domLine = c.dom.find((l) => l.startsWith("- ref:") && l.includes("<h1 ") && l.includes("tracking-tight") && l.endsWith("> h1"))
        const domTarget = c.dom.find((l) => l.startsWith("+ target:") && l.includes("<h1 ") && !l.includes("tracking-tight"))
        const at = (c.at ?? []).some((a) => /\bh1\b/.test(a) && a.includes('"Button"'))
        return { r, ok: !r.pass && !!domLine && !!domTarget && c.pixels > 0 && at, domLine, at: c.at?.[0] }
      })
      return { ok: rs.length === 6 && checks.every((c) => c.ok), detail: checks.map((c) => `${c.r.theme}/${c.r.viewport}: dom ${c.domLine ? "h1 line" : "MISSING"}; pixels at ${c.at}`).join("\n       ") }
    },
  },
  {
    name: "padding-1px",
    stub: ["-paths", "/docs/installation", "-css", "[data-slot=docs] h1{padding-left:1px}"],
    run: ["--no-scenarios", "/docs/installation"],
    expect: (rs) => {
      const checks = rs.map((r) => {
        const c = r.captures.find((x) => x.label === "initial")
        const at = (c.at ?? []).some((a) => /\bh1\b/.test(a) && a.includes('"Installation"'))
        return { r, ok: !r.pass && c.pixels > 0 && c.dom.length === 0 && at, c }
      })
      return { ok: rs.length === 6 && checks.every((c) => c.ok), detail: checks.map(({ r, c }) => `${r.theme}/${r.viewport}: ${c.pixels}px in ${c.regions?.length} region(s), box ${JSON.stringify(c.regions?.[0]?.box)} at ${c.at?.[0]}; dom lines ${c.dom.length}`).join("\n       ") }
    },
  },
  {
    name: "anim-duration",
    stub: ["-paths", "/docs/installation", "-css", "[data-slot=popover-content]{animation-duration:250ms!important}"],
    run: ["--viewports", "mobile", "--only", "mobile-menu-open,mobile-menu-close", "/docs/installation"],
    expect: (rs) => {
      const checks = rs.map((r) => {
        const init = r.captures.find((x) => x.label === "initial")
        const c = r.captures.find((x) => x.label === "mobile-menu-close")
        const line = c?.animations?.diff?.find((l) => l.includes("popover-content") && /duration \d+ vs 250\b/.test(l))
        return { r, ok: init?.pass && !r.pass && !!line, line, frames: c?.frames?.filter((f) => f.pixels).map((f) => `${f.t}ms:${f.pixels}px`).join(",") }
      })
      return { ok: rs.length === 2 && checks.every((c) => c.ok), detail: checks.map((c) => `${c.r.theme}/${c.r.viewport}: initial ${c.ok ? "identical" : "?"}; ${c.line?.slice(0, 200) ?? "NO ANIMATION DIFF"}; differing frames ${c.frames || "none"}`).join("\n       ") }
    },
  },
]

let failures = 0
const results = []
const pick = opt("--case") ? new Set(opt("--case").split(",")) : null
for (const c of cases.filter((c) => !pick || pick.has(c.name))) {
  const stub = await startStub(c.stub)
  const out = path.join(outBase, c.name)
  fs.rmSync(out, { recursive: true, force: true })
  console.log(`\n=== ${c.name}: docs-stub ${c.stub.join(" ") || "(no change)"} at ${stub.url}`)
  const { code } = await runHarness(stub.url, out, c.run)
  stub.proc.kill()
  const res = c.expect(reports(out))
  results.push({ name: c.name, exit: code, ...res })
  if (!res.ok) failures++
}
fs.rmSync(bin, { force: true })
console.log("\n=== stub proof")
for (const r of results) console.log(`${r.ok ? "OK  " : "BAD "} ${r.name} (harness exit ${r.exit})\n       ${r.detail}`)
process.exit(failures ? 1 : 0)
