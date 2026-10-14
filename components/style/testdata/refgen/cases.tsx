// Edge cases of transformStyle and of prettier's class sorting, beyond
// upstream's unit tests. Each "// ---" section is one file, transformed with
// the nova map (upstream.mjs). Not compiled: only parsed.
// --- cva claims a token, later uses only lose it; duplicate tokens across variants
import { cva } from "class-variance-authority"
import { cn } from "@/lib/utils"

const a = cva("cn-button inline-flex", {
  variants: {
    variant: {
      default: "cn-button-variant-default cn-button",
      outline: "cn-button-variant-outline   extra   ",
      "with-quote": "cn-button-variant-default",
    },
    size: { sm: "cn-button-size-sm", none: "", allow: "cn-menu-target cn-unknown-token" },
  },
  compoundVariants: [{ variant: "outline", size: "sm", className: "  z-10 flex  p-2 " }],
  defaultVariants: { variant: "default" },
})

function A({ className }: { className?: string }) {
  return (
    <>
      <button className={cn(a({ className }))} />
      <div className="cn-button" />
      <div className={cn("cn-button-variant-default p-1", className)} />
      <div className={cn("cn-card", "cn-card-header", className)} />
    </>
  )
}
// --- first argument not a string: classes prepended, empty arguments dropped
import { cn } from "@/lib/utils"

function B({ className, active }: { className?: string; active?: boolean }) {
  return (
    <div>
      <span className={cn(className, "cn-card px-2", "", active && "cn-card-header z-1")} />
      <span className={cn(className, "   ", "cn-unknown-x")} />
      <span className={cn(`cn-card-title text-sm`, className)} />
      <span className={cn("cn-card-footer  ", `  m-1`, "cn-unknown")} />
      <span className={"cn-card-description mt-2"} />
      <span className={`cn-card-action ml-1`} />
      <span className="cn-menu-target cn-card-content" />
      <span className="cn-nothing-here  flex" />
      <span className={active ? "cn-card" : "flex"} />
    </div>
  )
}
// --- nested cn() calls, className properties, mergeProps, render props
import { mergeProps } from "@base-ui/react/merge-props"
import { useRender } from "@base-ui/react/use-render"
import { cn } from "@/lib/utils"

function C({ className, render, ...props }: any) {
  const state = { className: cn("cn-card-header gap-2", className) }
  const literal = { className: "cn-card-title flex" }
  const plain = cn("cn-card px-1", className)
  return useRender({
    defaultTagName: "div",
    props: mergeProps<"div">(
      { className: cn("cn-card-content", "cn-menu-translucent p-1", className) },
      { className: cn("cn-card-content", className) },
      props
    ),
    render,
    state,
    literal,
    plain,
  }) && (
    <div
      className={composeRenderProps(className, (className) =>
        cn("cn-card-footer flex", className, cn("cn-card-action"))
      )}
    />
  )
}
// --- sorting: concatenation, templates, tagged templates, object keys, ellipsis
import { cn, cva } from "@/lib/utils"

const v = cva("p-2 flex   m-1", {
  variants: { s: { a: "z-10 static block", b: `px-2 ${"py-1 m-0"} mx-1 flex` } },
})

function D({ className, x }: { className?: string; x: string }) {
  return (
    <div className={cn("p-1 " + x + " flex block", { "p-2 m-1": x, "flex z-1": !x }, String.raw`p-3 m-2`, cn.raw`z-1 flex`)}>
      <span className={`p-4 ${x} m-1 flex ${x}block`} />
      <span className={"m-2 p-1" + x} />
      <span className={cn("... flex m-1", "flex flex p-1 flex", "unknown-a unknown-a flex")} />
      <span className={cn(x ? "p-2 m-2" : "m-3 p-3", ["p-1 m-1"], v({ s: "a" }))} />
      <span className={cn("{{ m-1 p-1 }}", " ", "\tm-1\np-1 ")} />
      <span className={cn(`m-1 ${x} p-1`, `${"p-1 m-1"}`)} />
    </div>
  )
}
