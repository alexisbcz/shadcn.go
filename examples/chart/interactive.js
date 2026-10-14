// Interactive chart examples: the series buttons in the card header set the
// Bar/Line dataKey and color, like setActiveChart upstream. The buttons are
// desktop then mobile, as rendered by the example.
import { define } from "./core.js"

const SERIES = ["desktop", "mobile"]

define("chart-example-interactive", {
  selector: '[data-slot="card"]:has(> [data-slot="card-header"] button[data-active]):has([data-slot="chart"])',
  mount(root) {
    const buttons = () => [...root.querySelectorAll(':scope > [data-slot="card-header"] button[data-active]')]
    const onClick = (e) => {
      const button = e.target.closest("button[data-active]")
      const list = buttons()
      const i = list.indexOf(button)
      if (i < 0) return
      const key = SERIES[i]
      for (const b of list) b.setAttribute("data-active", String(b === button))
      const chart = root.querySelector('[data-slot="chart"]')?.chart
      chart?.update((spec) => {
        for (const child of spec.children) {
          if (child.type !== "Bar" && child.type !== "Line") continue
          child.props.dataKey = key
          if (child.type === "Bar") child.props.fill = `var(--color-${key})`
          else child.props.stroke = `var(--color-${key})`
        }
      })
    }
    root.addEventListener("click", onClick)
    return () => root.removeEventListener("click", onClick)
  },
})
