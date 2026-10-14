// Prints porting status: for every upstream reference, whether it is ported
// to Go and whether its last parity run passed, grouped by component.
//
//   node parity/status.mjs            # summary per component
//   node parity/status.mjs --todo     # also list what is not ported or failing
import fs from "node:fs"
import path from "node:path"
import { execFileSync } from "node:child_process"

const repo = path.resolve(import.meta.dirname, "..")
const refs = JSON.parse(execFileSync("go", ["run", "./cmd/parity", "-list", "-root", repo], { cwd: repo, encoding: "utf8" }))
const manifest = JSON.parse(fs.readFileSync(path.join(repo, "parity/manifest.json"), "utf8"))
const owner = new Map()
for (const g of manifest.groups) for (const r of g.references ?? []) owner.set(r, g.id)

const rows = new Map()
let ported = 0
let passing = 0
for (const r of refs) {
  const report = path.join(repo, "parity/out", r.name, "report.json")
  const pass = r.ported && fs.existsSync(report) ? JSON.parse(fs.readFileSync(report, "utf8")).pass : null
  const key = owner.get(r.name) ?? "(unassigned)"
  if (!rows.has(key)) rows.set(key, { total: 0, ported: 0, passing: 0, todo: [] })
  const row = rows.get(key)
  row.total++
  if (r.ported) row.ported++, ported++
  if (pass) row.passing++, passing++
  else row.todo.push(`${r.name}${r.ported ? (pass === false ? " (failing)" : " (not run)") : ""}`)
}
for (const [k, v] of [...rows].sort()) {
  console.log(`${k.padEnd(28)} ${String(v.passing).padStart(3)}/${String(v.total).padEnd(3)} passing  (${v.ported} ported)`)
  if (process.argv.includes("--todo") && v.todo.length) console.log("    " + v.todo.join(", "))
}
console.log(`\n${passing}/${refs.length} references pass, ${ported} ported`)
