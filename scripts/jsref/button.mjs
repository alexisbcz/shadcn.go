import fs from "node:fs"
import { cva } from "class-variance-authority"
import { clsx } from "clsx"
import { twMerge } from "tailwind-merge"
const cn = (...a) => twMerge(clsx(a))
const src = fs.readFileSync(process.argv[2], "utf8")
const start = src.indexOf("cva(")
const end = src.indexOf("\n)\n", start)
const buttonVariants = eval(src.slice(start, end + 2))
const variants = ["default","destructive","outline","secondary","ghost","link"]
const sizes = ["default","xs","sm","lg","icon","icon-xs","icon-sm","icon-lg"]
const extras = ["", "rounded-full", "w-full px-8 text-base", "gap-4 h-12 bg-red-500 hover:bg-red-600", "size-12 [&_svg:not([class*='size-'])]:size-6", "border-2 shadow-none focus-visible:ring-2"]
const out = []
for (const variant of variants) for (const size of sizes) for (const className of extras)
  out.push({ variant, size, className, want: cn(buttonVariants({ variant, size, className })) })
console.log(JSON.stringify(out))
