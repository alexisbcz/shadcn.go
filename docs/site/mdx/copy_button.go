package mdx

import (
	"shadcn.go/cn"
	"shadcn.go/components/icons/tabler"
	"shadcn.go/components/ui/button"
	"shadcn.go/docs/content"
	h "shadcn.go/libhtml"
)

// CopyButton: ui/apps/v4/components/copy-button.tsx (new-york-v4 Button),
// rendered by package content before each highlighted code block's <code>.
// Props: value, src, className, variant. The copy itself is done by
// docs-mdx.js, which copies the text of the code block the button belongs
// to (the value upstream passes is that same text).
func CopyButton(_ *Env, _ *content.Context, props content.Props, _ []h.Node) (h.Node, error) {
	return copyButton(props.String("variant"), className(props)), nil
}

// copyButton renders CopyButton with a variant ("" for ghost) and extra
// classes.
func copyButton(variant, class string) h.Node {
	if variant == "" {
		variant = "ghost"
	}
	return button.Button(
		h.Span(h.Class("sr-only"), "Copy"),
		tabler.IconCopy(),
	).Variant(button.Variant(variant)).Size(button.SizeIcon).
		Data("slot", "copy-button").
		Class(cn.Cn("absolute top-3 right-2 z-10 size-7 bg-code hover:opacity-100 focus-visible:opacity-100", class)).
		Data("copied", "false")
}

func init() { Register("CopyButton", CopyButton) }
