package attachment

import (
	"io"

	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	"shadcn.go/cva"
)

// The registry/bases/radix/ui/attachment.tsx classes, resolved per style.
// The parts render them inside style.Scope; outside, the new-york-v4
// classes.
var file = style.NewFile("ui/attachment.tsx")

var (
	sAttachment = file.CVA(
		"cn-attachment group/attachment relative flex max-w-full min-w-0 shrink-0 flex-wrap border bg-card text-card-foreground transition-colors has-[>a,>button]:hover:bg-muted/50 data-[state=error]:border-destructive/30 data-[state=idle]:border-dashed",
		[]style.Variant{
			{Name: "size", Options: []style.Option{
				{Name: "default", Class: "cn-attachment-size-default"},
				{Name: "sm", Class: "cn-attachment-size-sm"},
				{Name: "xs", Class: "cn-attachment-size-xs"},
			}},
			{Name: "orientation", Options: []style.Option{
				{Name: "horizontal", Class: "cn-attachment-orientation-horizontal items-center"},
				{Name: "vertical", Class: "cn-attachment-orientation-vertical flex-col"},
			}},
		},
		nil,
	)
	sMedia = file.CVA(
		"cn-attachment-media relative flex aspect-square shrink-0 items-center justify-center overflow-hidden group-data-[state=error]/attachment:bg-destructive/10 group-data-[state=error]/attachment:text-destructive [&_svg]:pointer-events-none",
		[]style.Variant{{Name: "variant", Options: []style.Option{
			{Name: "icon", Class: "cn-attachment-media-variant-icon"},
			{Name: "image", Class: "cn-attachment-media-variant-image *:[img]:aspect-square *:[img]:w-full *:[img]:object-cover"},
		}}},
		cva.Selection{"variant": "icon"},
	)
	sContent     = file.Cn("cn-attachment-content max-w-full min-w-0 flex-1", style.Arg)
	sTitle       = file.Cn("cn-attachment-title block max-w-full min-w-0 truncate group-data-[state=processing]/attachment:shimmer group-data-[state=uploading]/attachment:shimmer", style.Arg)
	sDescription = file.Cn("cn-attachment-description block min-w-0 truncate text-muted-foreground group-data-[state=error]/attachment:text-destructive/80", "max-w-full", style.Arg)
	sActions     = file.Cn("cn-attachment-actions flex shrink-0 items-center", style.Arg)
	sTrigger     = file.Cn("cn-attachment-trigger absolute inset-0 z-10 outline-none", style.Arg)
	sGroup       = file.Cn("cn-attachment-group flex min-w-0 scroll-fade-x snap-x snap-mandatory scrollbar-none overflow-x-auto overscroll-x-contain *:data-[slot=attachment]:flex-none *:data-[slot=attachment]:snap-start", style.Arg)
)

func styleOf(w io.Writer) style.Style {
	if w == nil {
		return ""
	}
	st, _ := style.FromWriter(w)
	return st
}

func attachmentClass(st style.Style, s Size, o Orientation, className string) string {
	if st == "" {
		return Variants(s, o, className)
	}
	return lib.Cn(sAttachment.Classes(st, cva.Selection{"size": string(s), "orientation": string(o)}), className)
}

func mediaClass(st style.Style, v MediaVariant, className string) string {
	if st == "" {
		return MediaVariants(v, className)
	}
	return lib.Cn(sMedia.Classes(st, cva.Selection{"variant": string(v)}), className)
}
