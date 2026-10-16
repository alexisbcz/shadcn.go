// chart-pie-interactive: the month select sets the active sector (the
// Pie's shape) and the center label's value, like the activeMonth state
// upstream. Data and months are read from the chart's spec.
import { define } from "./core.js"

define("chart-pie-example-interactive", {
  selector: '[data-slot="card"][data-chart]:has([data-slot="select-trigger"]):has([data-slot="chart"])',
  mount(root) {
    const onChange = (e) => {
      const chart = root.querySelector('[data-slot="chart"]')?.chart
      const month = e.detail?.value
      if (!chart || month == null) return
      chart.update((spec) => {
        for (const pie of spec.children) {
          if (pie.type !== "Pie") continue
          const data = pie.props.data ?? []
          const index = data.findIndex((item) => item.month === month)
          if (pie.props.shape?.$fn === "pieActiveSectorRing") pie.props.shape.args[0] = index
          for (const label of pie.children ?? []) {
            if (label.type === "Label" && label.props.content?.$fn === "pieCenterText") label.props.content.args[0] = data[index]?.[pie.props.dataKey]
          }
        }
      })
    }
    root.addEventListener("shadcn:value-change", onChange)
    return () => root.removeEventListener("shadcn:value-change", onChange)
  },
})
