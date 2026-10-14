package style_test

import (
	"fmt"

	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	"shadcn.go/cva"
)

// The Button of registry/bases/radix/ui/button.tsx, written once with its
// tokens.
var buttonFile = style.NewFile("ui/button.tsx")

var buttonVariants = buttonFile.CVA(
	"cn-button group/button inline-flex shrink-0 items-center justify-center whitespace-nowrap transition-all outline-none select-none disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:shrink-0",
	[]style.Variant{
		{Name: "variant", Options: []style.Option{
			{Name: "default", Class: "cn-button-variant-default"},
			{Name: "outline", Class: "cn-button-variant-outline"},
		}},
		{Name: "size", Options: []style.Option{
			{Name: "default", Class: "cn-button-size-default"},
			{Name: "sm", Class: "cn-button-size-sm"},
		}},
	},
	cva.Selection{"variant": "default", "size": "default"},
)

// The CardTitle of registry/bases/radix/ui/card.tsx:
// className={cn("cn-card-title cn-font-heading", className)}.
var cardFile = style.NewFile("ui/card.tsx")

var cardTitle = cardFile.Cn("cn-card-title cn-font-heading", style.Arg)

func ExampleFile() {
	for _, st := range []style.Style{style.Nova, style.Lyra} {
		// className={cn(buttonVariants({ variant, size, className }))}
		fmt.Println(st, lib.Cn(buttonVariants.Classes(st, cva.Selection{"variant": "outline", "size": "sm"}, "w-full")))
	}
	// Output:
	// nova group/button inline-flex shrink-0 items-center justify-center border bg-clip-padding font-medium whitespace-nowrap transition-all outline-none select-none focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/50 active:not-aria-[haspopup]:translate-y-px disabled:pointer-events-none disabled:opacity-50 aria-invalid:border-destructive aria-invalid:ring-3 aria-invalid:ring-destructive/20 dark:aria-invalid:border-destructive/50 dark:aria-invalid:ring-destructive/40 [&_svg]:pointer-events-none [&_svg]:shrink-0 border-border bg-background hover:bg-muted hover:text-foreground aria-expanded:bg-muted aria-expanded:text-foreground dark:border-input dark:bg-input/30 dark:hover:bg-input/50 h-7 gap-1 rounded-[min(var(--radius-md),12px)] px-2.5 text-[0.8rem] in-data-[slot=button-group]:rounded-lg has-data-[icon=inline-end]:pr-1.5 has-data-[icon=inline-start]:pl-1.5 [&_svg:not([class*='size-'])]:size-3.5 w-full
	// lyra group/button inline-flex shrink-0 items-center justify-center border bg-clip-padding text-xs font-medium whitespace-nowrap transition-all outline-none select-none focus-visible:border-ring focus-visible:ring-1 focus-visible:ring-ring/50 active:not-aria-[haspopup]:translate-y-px disabled:pointer-events-none disabled:opacity-50 aria-invalid:border-destructive aria-invalid:ring-1 aria-invalid:ring-destructive/20 dark:aria-invalid:border-destructive/50 dark:aria-invalid:ring-destructive/40 [&_svg]:pointer-events-none [&_svg]:shrink-0 border-border bg-background hover:bg-muted hover:text-foreground aria-expanded:bg-muted aria-expanded:text-foreground dark:border-input dark:bg-input/30 dark:hover:bg-input/50 h-7 gap-1 rounded-none px-2.5 has-data-[icon=inline-end]:pr-1.5 has-data-[icon=inline-start]:pl-1.5 [&_svg:not([class*='size-'])]:size-3.5 w-full
}

func ExampleFile_Cn() {
	for _, st := range []style.Style{style.Nova, style.Vega} {
		fmt.Println(st, cardTitle.Args(st))
		fmt.Println(st, cardTitle.Class(st, "text-lg"))
	}
	// Output:
	// nova [cn-font-heading text-base leading-snug font-medium group-data-[size=sm]/card:text-sm style.Arg]
	// nova cn-font-heading font-medium group-data-[size=sm]/card:text-sm text-lg
	// vega [cn-font-heading text-base leading-normal font-medium group-data-[size=sm]/card:text-sm style.Arg]
	// vega cn-font-heading font-medium group-data-[size=sm]/card:text-sm text-lg
}

func ExampleResolve() {
	fmt.Println(style.Resolve(style.Nova, "cn-card-description"))
	fmt.Println(style.Resolve(style.Mira, "cn-card-description"))
	// Output:
	// text-sm text-muted-foreground
	// text-xs/relaxed text-muted-foreground
}
