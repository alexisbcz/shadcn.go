package bubble

import (
	"io"

	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	"shadcn.go/cva"
)

// The registry/bases/radix/ui/bubble.tsx classes, resolved per style. The
// parts render them inside style.Scope; outside, the new-york-v4 classes.
var file = style.NewFile("ui/bubble.tsx")

var (
	sGroup  = file.Cn("cn-bubble-group flex min-w-0 flex-col", style.Arg)
	sBubble = file.CVA("cn-bubble group/bubble relative flex w-fit min-w-0 flex-col",
		[]style.Variant{{Name: "variant", Options: []style.Option{
			{Name: "default", Class: "cn-bubble-variant-default"},
			{Name: "secondary", Class: "cn-bubble-variant-secondary"},
			{Name: "muted", Class: "cn-bubble-variant-muted"},
			{Name: "tinted", Class: "cn-bubble-variant-tinted"},
			{Name: "outline", Class: "cn-bubble-variant-outline"},
			{Name: "ghost", Class: "cn-bubble-variant-ghost"},
			{Name: "destructive", Class: "cn-bubble-variant-destructive"},
		}}},
		cva.Selection{"variant": "default"},
	)
	sContent   = file.Cn("cn-bubble-content w-fit max-w-full min-w-0 overflow-hidden wrap-break-word [button]:text-left [button,a]:transition-colors", style.Arg)
	sReactions = file.CVA("cn-bubble-reactions absolute z-10 flex w-fit items-center justify-center",
		[]style.Variant{
			{Name: "side", Options: []style.Option{
				{Name: "top", Class: "cn-bubble-reactions-side-top"},
				{Name: "bottom", Class: "cn-bubble-reactions-side-bottom"},
			}},
			{Name: "align", Options: []style.Option{
				{Name: "start", Class: "cn-bubble-reactions-align-start"},
				{Name: "end", Class: "cn-bubble-reactions-align-end"},
			}},
		},
		cva.Selection{"side": "bottom", "align": "end"},
	)
)

func styleOf(w io.Writer) style.Style {
	if w == nil {
		return ""
	}
	st, _ := style.FromWriter(w)
	return st
}

func bubbleClass(st style.Style, v Variant, className string) string {
	if st == "" {
		return Variants(v, className)
	}
	return lib.Cn(sBubble.Classes(st, cva.Selection{"variant": string(v)}), className)
}

func reactionsClass(st style.Style, side Side, align Align, className string) string {
	if st == "" {
		return ReactionsVariants(side, align, className)
	}
	return lib.Cn(sReactions.Classes(st, cva.Selection{"side": string(side), "align": string(align)}), className)
}

func contentClassFor(st style.Style, className string) string {
	if st == "" {
		return lib.Cn(contentClass, className)
	}
	return sContent.Class(st, className)
}
