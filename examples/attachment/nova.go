package attachment

import (
	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/attachment"
	"shadcn.go/components/ui/dialog"
	"shadcn.go/components/ui/spinner"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// Docs previews of the attachment page
// (ui/apps/v4/examples/radix/attachment-*.tsx). Upstream's examples import
// the radix-rhea components, so they render in style.Rhea.
func init() {
	docs := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: "nova--" + name, Component: "attachment", Render: registry.Static(func() h.Node {
			return style.Scope(style.Rhea, fn())
		})}
	}
	registry.Register(
		docs("attachment-demo", DocsDemo),
		docs("attachment-group", DocsGroup),
		docs("attachment-image", DocsImage),
		docs("attachment-sizes", DocsSizes),
		docs("attachment-states", DocsStates),
		docs("attachment-trigger", DocsTrigger),
	)
}

type docsImage struct{ name, meta, src, alt string }

var docsImages = []docsImage{
	{"workspace.png", "PNG · 820 KB", "https://images.unsplash.com/photo-1497366754035-f200968a6e72?w=900&auto=format&fit=crop&q=80", "Workspace"},
	{"desk-reference.jpg", "JPG · 1.1 MB", "https://images.unsplash.com/photo-1497215728101-856f4ea42174?w=900&auto=format&fit=crop&q=80", "Desk"},
	{"office-reference.jpg", "JPG · 940 KB", "https://images.unsplash.com/photo-1497366811353-6870744d04b2?w=900&auto=format&fit=crop&q=80", "Office"},
}

func DocsDemo() h.Node {
	return h.Div(h.Class("mx-auto flex w-full max-w-sm flex-col gap-3 py-12"),
		attachment.Group(h.Map(docsImages, func(im docsImage) h.Node {
			return attachment.Attachment(attachment.Vertical,
				attachment.Media(attachment.MediaImage, h.Img(h.Src(im.src), h.Alt(im.alt))),
				details(im.name, im.meta),
			)
		})),
		attachment.Attachment(attachment.StateUploading, h.Class("w-full"),
			attachment.Media(spinner.Spinner()),
			details("sales-dashboard.pdf", "Uploading · 64%"),
			attachment.Actions(remove("Cancel upload")),
		),
		attachment.Attachment(h.Class("w-full"),
			attachment.Media(lucide.FileCode()),
			details("message-renderer.tsx", "TypeScript · 12 KB"),
			attachment.Actions(remove("Remove message-renderer.tsx")),
		),
	)
}

func DocsGroup() h.Node {
	type item struct {
		name, meta string
		icon       func(...any) *h.Element
		src        string
	}
	items := []item{
		{"briefing-notes.pdf", "PDF · 1.4 MB", lucide.FileText, ""},
		{"workspace.png", "PNG · 820 KB", nil, "https://images.unsplash.com/photo-1497366754035-f200968a6e72?w=900&auto=format&fit=crop&q=80"},
		{"customers.csv", "CSV · 18 KB", lucide.Table, ""},
		{"renderer.tsx", "TSX · 12 KB", lucide.FileCode, ""},
	}
	return h.Div(h.Class("mx-auto w-full max-w-sm py-12"),
		attachment.Group(h.Class("w-full"), h.Map(items, func(it item) h.Node {
			var media h.Node
			if it.src != "" {
				media = attachment.Media(attachment.MediaImage, h.Img(h.Src(it.src), h.Alt(it.name)))
			} else {
				media = attachment.Media(it.icon())
			}
			return attachment.Attachment(h.Class("w-64"),
				media,
				details(it.name, it.meta),
				attachment.Actions(remove("Remove "+it.name)),
			)
		})),
	)
}

func DocsImage() h.Node {
	return h.Div(h.Class("mx-auto w-full max-w-sm py-12"),
		attachment.Group(h.Class("w-full"), h.Map(docsImages, func(im docsImage) h.Node {
			return attachment.Attachment(attachment.Vertical,
				attachment.Media(attachment.MediaImage, h.Img(h.Src(im.src), h.Alt(im.alt))),
				details(im.name, im.meta),
				attachment.Actions(remove("Remove "+im.name)),
				attachment.Trigger(lib.AsChild,
					h.A(h.Href(im.src), h.Target("_blank"), h.Rel("noreferrer"), h.Aria("label", "Open "+im.name)),
				),
			)
		})),
	)
}

func DocsSizes() h.Node {
	return h.Div(h.Class("mx-auto flex w-full max-w-sm flex-col gap-3 py-12"),
		attachment.Attachment(attachment.SizeDefault, h.Class("w-full"),
			attachment.Media(lucide.FileText()),
			details("Default attachment", "PDF · 2.4 MB"),
		),
		attachment.Attachment(attachment.SizeSm, h.Class("w-full"),
			attachment.Media(lucide.FileText()),
			details("Small attachment", "PDF · 2.4 MB"),
		),
		attachment.Attachment(attachment.SizeXs, h.Class("w-full"),
			attachment.Media(lucide.FileText()),
			attachment.Content(attachment.Title("Extra small attachment")),
		),
	)
}

func DocsStates() h.Node {
	return h.Div(h.Class("mx-auto flex w-full max-w-sm flex-col gap-2 py-12"),
		attachment.Attachment(attachment.StateIdle, h.Class("w-full"),
			attachment.Media(lucide.Clock()),
			details("selected-file.pdf", "Ready to upload"),
			attachment.Actions(remove("Remove selected-file.pdf")),
		),
		attachment.Attachment(attachment.StateUploading, h.Class("w-full"),
			attachment.Media(spinner.Spinner()),
			details("design-system.zip", "Uploading · 64%"),
			attachment.Actions(remove("Cancel upload")),
		),
		attachment.Attachment(attachment.StateProcessing, h.Class("w-full"),
			attachment.Media(lucide.FileText()),
			details("market-research.pdf", "Processing document"),
			attachment.Actions(remove("Remove market-research.pdf")),
		),
		attachment.Attachment(attachment.StateError, h.Class("w-full"),
			attachment.Media(lucide.FileWarning()),
			details("financial-model.xlsx", "Upload failed. Try again."),
			attachment.Actions(
				remove("Retry upload", lucide.RefreshCw()),
				remove("Remove financial-model.xlsx"),
			),
		),
		attachment.Attachment(attachment.StateDone, h.Class("w-full"),
			attachment.Media(lucide.Check()),
			details("uploaded-report.pdf", "Uploaded · 1.8 MB"),
			attachment.Actions(remove("Remove uploaded-report.pdf")),
		),
	)
}

func DocsTrigger() h.Node {
	return h.Div(h.Class("mx-auto w-full max-w-sm py-12"),
		dialog.Dialog(
			attachment.Attachment(h.Class("w-full"),
				attachment.Media(lucide.FileSearch()),
				details("research-summary.pdf", "Open preview dialog"),
				attachment.Actions(
					remove("Copy link", lucide.Copy()),
					remove("Remove research-summary.pdf"),
				),
				dialog.Trigger(lib.AsChild, attachment.Trigger(h.Aria("label", "Preview research-summary.pdf"))),
			),
			dialog.Content(h.Class("sm:max-w-md"),
				dialog.Header(
					dialog.Title("research-summary.pdf"),
					dialog.Description("The attachment trigger fills the card and opens the dialog, while the actions stay independently clickable above it."),
				),
			),
		),
	)
}
