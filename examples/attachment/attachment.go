// Package attachment ports the attachment parity fixtures (upstream ships
// Attachment without new-york-v4 examples).
package attachment

import (
	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/lib"
	"shadcn.go/components/ui/attachment"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/dialog"
	"shadcn.go/components/ui/spinner"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

func init() {
	ex := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: name, Component: "attachment", Render: registry.Static(fn)}
	}
	registry.Register(
		ex("attachment-states", States),
		ex("attachment-sizes", Sizes),
		ex("attachment-trigger", Trigger),
		ex("attachment-group", Group),
	)
}

// pixel is a 4x4 orange PNG: an image that loads without the network.
const pixel = "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAQAAAAECAIAAAAmkwkpAAAAEElEQVR4nGP4WSwGRwzEcQCBwhghZfPo6AAAAABJRU5ErkJggg=="

func remove(label string, icon ...any) h.Node {
	if len(icon) == 0 {
		icon = []any{lucide.X()}
	}
	return attachment.Action(append([]any{h.Aria("label", label)}, icon...)...)
}

func details(title, description string) h.Node {
	return attachment.Content(attachment.Title(title), attachment.Description(description))
}

// States is the parity fixture attachment-states.
func States() h.Node {
	return h.Div(h.Class("flex w-full max-w-2xl flex-col gap-6 py-6"),
		h.Div(h.Class("flex flex-col gap-2"),
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
			),
			attachment.Attachment(attachment.StateError, h.Class("w-full"),
				attachment.Media(lucide.FileWarning()),
				details("financial-model.xlsx", "Upload failed. Try again."),
				attachment.Actions(
					remove("Retry upload", lucide.RefreshCw()),
					attachment.Action(button.VariantDestructive, h.Aria("label", "Remove financial-model.xlsx"), lucide.X()),
				),
			),
			attachment.Attachment(
				attachment.Media(lucide.Check()),
				details("uploaded-report.pdf", "Uploaded · 1.8 MB"),
				attachment.Actions(attachment.Action(button.SizeIconSm, h.Aria("label", "Remove uploaded-report.pdf"), lucide.X())),
			),
		),
		h.Div(h.Class("flex flex-wrap gap-3"),
			attachment.Attachment(attachment.StateIdle, attachment.Vertical,
				attachment.Media(lucide.Clock()),
				details("selected.pdf", "Ready"),
			),
			attachment.Attachment(attachment.StateUploading, attachment.Vertical,
				attachment.Media(spinner.Spinner()),
				details("design.zip", "64%"),
			),
			attachment.Attachment(attachment.StateProcessing, attachment.Vertical,
				attachment.Media(lucide.FileText()),
			),
			attachment.Attachment(attachment.StateError, attachment.Vertical,
				attachment.Media(lucide.FileWarning()),
				details("model.xlsx", "Failed"),
				attachment.Actions(remove("Retry", lucide.RefreshCw())),
			),
			attachment.Attachment(attachment.StateDone, attachment.Vertical,
				attachment.Media(lucide.Check()),
				details("uploaded-report-with-a-long-name.pdf", "1.8 MB"),
			),
		),
	)
}

// Sizes is the parity fixture attachment-sizes.
func Sizes() h.Node {
	img := func(alt string) h.Node { return h.Img(h.Src(pixel), h.Alt(alt)) }
	return h.Div(h.Class("flex w-full max-w-3xl flex-col gap-4 py-4"),
		h.Map(attachment.SizeValues, func(s attachment.Size) h.Node {
			return h.Div(h.Class("flex flex-wrap items-start gap-3"),
				attachment.Attachment(s,
					attachment.Media(lucide.FileText()),
					attachment.Content(attachment.Title(string(s), h.Comment(""), " attachment"), attachment.Description("PDF · 2.4 MB")),
					attachment.Actions(remove("Remove")),
				),
				attachment.Attachment(s, attachment.Content(attachment.Title("No media"))),
				attachment.Attachment(s, attachment.Media(attachment.MediaImage, img("Pixel"))),
				attachment.Attachment(s, attachment.Vertical,
					attachment.Media(lucide.FileText()),
					attachment.Content(attachment.Title(string(s), h.Comment(""), ".pdf"), attachment.Description("2.4 MB")),
				),
				attachment.Attachment(s, attachment.Vertical, attachment.StateUploading,
					attachment.Media(spinner.Spinner()),
				),
			)
		}),
	)
}

// Group is the parity fixture attachment-group.
func Group() h.Node {
	img := func(alt string) h.Node { return h.Img(h.Src(pixel), h.Alt(alt)) }
	return h.Div(h.Class("flex w-full max-w-3xl flex-col gap-4 py-4"),
		attachment.Group(h.Class("w-full max-w-md"),
			h.Map([]attachment.State{attachment.StateIdle, attachment.StateUploading, attachment.StateDone}, func(st attachment.State) h.Node {
				return attachment.Attachment(st, attachment.Vertical,
					attachment.Media(attachment.MediaImage, img(string(st))),
					attachment.Content(attachment.Title(string(st), h.Comment(""), ".png"), attachment.Description("PNG · 820 KB")),
					attachment.Actions(remove("Remove "+string(st)+".png")),
				)
			}),
			attachment.Attachment(attachment.StateError, attachment.Vertical, attachment.Media(attachment.MediaImage, img("error"))),
			attachment.Attachment(attachment.StateProcessing, attachment.Vertical, attachment.Media(attachment.MediaImage, img("processing"))),
		),
	)
}

// Trigger is the parity fixture attachment-trigger.
func Trigger() h.Node {
	return h.Div(h.Class("flex w-full max-w-sm flex-col gap-3 py-12"),
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
		attachment.Attachment(attachment.SizeSm, h.Class("w-full"),
			attachment.Media(lucide.FileText()),
			details("notes.txt", "Open as link"),
			attachment.Trigger(lib.AsChild, h.A(h.Href("#notes"), h.Aria("label", "Open notes.txt"))),
		),
		attachment.Attachment(attachment.SizeXs, attachment.StateIdle, h.Class("w-full"),
			attachment.Content(attachment.Title("draft.md")),
			attachment.Trigger(h.Type("submit"), h.Class("rounded-lg"), h.Aria("label", "Submit draft.md")),
		),
	)
}
