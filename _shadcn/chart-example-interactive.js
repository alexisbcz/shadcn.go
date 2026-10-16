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

// chart-area-interactive: the time range select filters the data with
// upstream's date arithmetic (local-time setDate on a UTC date), so the
// kept entries are the same as React's in every time zone. The full data
// is the chart's initial data ("90d" keeps every entry of the example).
define("chart-example-time-range", {
  selector: '[data-slot="card"]:not([data-chart]):has(> [data-slot="card-header"] [data-slot="select-trigger"]):has([data-slot="chart"])',
  mount(root) {
    let all = null
    const onChange = (e) => {
      const chart = root.querySelector('[data-slot="chart"]')?.chart
      if (!chart) return
      const timeRange = e.detail?.value
      chart.update((spec) => {
        all ??= spec.props.data
        spec.props.data = all.filter((item) => {
          const date = new Date(item.date)
          const referenceDate = new Date("2024-06-30")
          let daysToSubtract = 90
          if (timeRange === "30d") {
            daysToSubtract = 30
          } else if (timeRange === "7d") {
            daysToSubtract = 7
          }
          const startDate = new Date(referenceDate)
          startDate.setDate(startDate.getDate() - daysToSubtract)
          return date >= startDate
        })
      })
    }
    root.addEventListener("shadcn:value-change", onChange)
    return () => root.removeEventListener("shadcn:value-change", onChange)
  },
})
