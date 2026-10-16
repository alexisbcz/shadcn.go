// Tooltip formatters of chart-tooltip-formatter and chart-tooltip-advanced:
// the inline formatter closures of the upstream examples.
import { registerFn, h } from "./chart.js"

const label = (ctx, name) => ctx.config[name]?.label || name

const kcal = (value) =>
  h("div", { className: "ml-auto flex items-baseline gap-0.5 font-mono font-medium text-foreground tabular-nums" }, String(value), h("span", { className: "font-normal text-muted-foreground" }, "kcal"))

// (value, name) => <div ...>{label}<div ...>{value}<span>kcal</span></div></div>
registerFn("kcalFormatter", (args, ctx) => (value, name) =>
  h("div", { className: "flex min-w-[130px] items-center text-xs text-muted-foreground" }, label(ctx, name), kcal(value)),
)

// (value, name, item, index) => <>indicator, label, value, total after the last item</>
registerFn("kcalAdvanced", (args, ctx) => (value, name, item, index) => [
  h("div", { className: "h-2.5 w-2.5 shrink-0 rounded-[2px] bg-(--color-bg)", style: { "--color-bg": `var(--color-${name})` } }),
  label(ctx, name),
  kcal(value),
  index === 1
    ? h("div", { className: "mt-1.5 flex basis-full items-center border-t pt-1.5 text-xs font-medium text-foreground" }, "Total", kcal(item.payload.running + item.payload.swimming))
    : null,
])
