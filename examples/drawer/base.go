package drawer

import (
	_ "embed"

	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/badge"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/dialog"
	bd "shadcn.go/components/ui/drawer/basedrawer"
	"shadcn.go/components/ui/field"
	radiogroup "shadcn.go/components/ui/radio-group"
	"shadcn.go/components/ui/runtime"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// Docs previews of the base and aria drawer pages
// (ui/apps/v4/examples/{base,aria}/drawer-*.tsx): Base UI drawers, in
// base-rhea for the base page (styleName="base-rhea", prefix base--) and
// aria-nova for the aria page.
func init() {
	runtime.Module("example-drawer-base-demo", baseDemoScript)
	for _, v := range []struct {
		prefix string
		st     style.Style
		aria   bool
	}{{"base--", style.Rhea, false}, {"aria--", style.Nova, true}} {
		st, aria := v.st, v.aria
		reg := func(name string, fn func() h.Node) {
			registry.Register(registry.Example{Name: v.prefix + name, Component: "drawer", Render: registry.Static(func() h.Node {
				return style.Scope(st, fn())
			})})
		}
		reg("drawer-demo", BaseDemo)
		reg("drawer-sides", BaseSides)
		reg("drawer-dialog", func() h.Node { return BaseDialog(aria) })
		reg("drawer-nested", BaseNested)
		reg("drawer-non-modal", BaseNonModal)
		reg("drawer-snap-points", BaseSnapPoints)
		reg("drawer-swipe-handle", BaseSwipeHandle)
	}
}

//go:embed base_demo.js
var baseDemoScript string

var deliveryTimes = []struct{ value, label, description, badge string }{
	{"asap", "Standard delivery", "25–35 min · Driver assigned now", "Fastest"},
	{"5-00", "5:00 PM – 5:15 PM", "Prep starts at 4:45 PM", ""},
	{"5-30", "5:30 PM – 5:45 PM", "Good if you're heading home", ""},
	{"6-00", "6:00 PM – 6:15 PM", "Most popular · High demand", ""},
	{"6-30", "6:30 PM – 6:45 PM", "Last slot before kitchen closes", ""},
}

// BaseDemo is drawer-demo: pick a delivery time (right on desktop, down
// with a swipe handle on mobile); Confirm closes the drawer and toasts the
// pick (base_demo.js).
func BaseDemo() h.Node {
	group := radiogroup.RadioGroup(radiogroup.Value("asap"), h.Class("gap-2"))
	for _, t := range deliveryTimes {
		id := "delivery-" + t.value
		title := []any{h.Class("flex items-center gap-2"), t.label}
		if t.badge != "" {
			title = append(title, badge.Badge(badge.VariantSecondary, t.badge))
		}
		group.Children(field.Label(h.For(id),
			field.Field(field.Horizontal,
				field.Content(field.Title(title...), field.Description(t.description)),
				radiogroup.Item(radiogroup.Value(t.value), h.ID(id), h.Data("delivery-label", t.label)),
			),
		))
	}
	return bd.Drawer(bd.SwipeDirection("right"), bd.MobileSwipeDirection("down"), bd.MobileShowSwipeHandle,
		bd.Trigger(button.Button(button.VariantSecondary, "Open Drawer")),
		bd.Content(
			bd.Header(bd.Title("Pick a delivery time"), bd.Description("We'll prepare your order as soon as possible.")),
			h.Div(h.Class("flex-1 scroll-fade overflow-y-auto p-4"), group),
			bd.Footer(
				bd.Close(button.Button(h.Class("h-[34px]"), h.Data("drawer-confirm-delivery", ""), "Confirm Delivery Time")),
				bd.Close(button.Button(button.VariantOutline, "Cancel")),
			),
		),
	)
}

func BaseSides() h.Node {
	return bd.Drawer(bd.SwipeDirection("left"),
		bd.Trigger(button.Button(button.VariantSecondary, "Open Left Drawer")),
		bd.Content(
			bd.Header(bd.Title("Move Goal"), bd.Description("Set your daily activity goal.")),
			placeholder("size-full rounded-2xl bg-muted"),
			bd.Footer(bd.Close(button.Button("Close"))),
		),
	)
}

// BaseDialog is drawer-dialog: a Dialog from 768px, a Base UI drawer below
// (swapped by media.js). The aria version pads the form horizontally and
// adds a Cancel footer to the drawer.
func BaseDialog(aria bool) h.Node {
	formClass := "p-4"
	var footer []any
	if aria {
		formClass = "px-4"
		footer = append(footer, bd.Footer(bd.Class("pt-2"), bd.Close(button.Button(button.VariantOutline, "Cancel"))))
	}
	return h.Group(
		dialog.Dialog(
			dialog.Trigger(lib.AsChild, button.Button(button.VariantOutline, "Edit Profile")),
			dialog.Content(h.Class("sm:max-w-[425px]"),
				dialog.Header(
					dialog.Title("Edit profile"),
					dialog.Description("Make changes to your profile here. Click save when you're done."),
				),
				profileForm(""),
			),
		),
		h.Template(h.Data("parity-ignore", ""), h.Data("example-media-swap", "(min-width: 768px)"),
			bd.Drawer(
				bd.Trigger(button.Button(button.VariantOutline, "Edit Profile")),
				bd.Content(append([]any{
					bd.Header(bd.Class("text-left"),
						bd.Title("Edit profile"),
						bd.Description("Make changes to your profile here. Click save when you're done."),
					),
					profileForm(formClass),
				}, footer...)...),
			),
		),
	)
}

func placeholder(class string) h.Node {
	return h.Div(h.Class("flex-1 p-4"), h.Div(h.Class(class)))
}

const nestedPlaceholder = "bg-muted group-data-[swipe-axis=x]/drawer-popup:size-full group-data-[swipe-axis=y]/drawer-popup:aspect-video group-data-[swipe-axis=y]/drawer-popup:w-full"

// nestedLevel is one drawer of drawer-nested: right on desktop, down with
// a swipe handle on mobile (useIsMobile).
func nestedLevel(trigger button.Variant, label, title, description string, inner h.Node) h.Node {
	footer := []any{}
	if inner != nil {
		footer = append(footer, inner)
	}
	footer = append(footer, bd.Close(button.Button(button.VariantOutline, "Close")))
	return bd.Drawer(bd.SwipeDirection("right"), bd.MobileSwipeDirection("down"), bd.MobileShowSwipeHandle,
		bd.Trigger(button.Button(trigger, label)),
		bd.Content(
			bd.Header(bd.Title(title), bd.Description(description)),
			placeholder(nestedPlaceholder),
			bd.Footer(footer...),
		),
	)
}

func BaseNested() h.Node {
	fourth := nestedLevel(button.VariantOutline, "Open Fourth Drawer", "Fourth Drawer", "This is the frontmost drawer in the stack.", nil)
	third := nestedLevel(button.VariantOutline, "Open Third Drawer", "Third Drawer", "Two drawers are stacked behind this one.", fourth)
	second := nestedLevel(button.VariantOutline, "Open Nested Drawer", "Nested Drawer", "The parent drawer stays mounted behind this one.", third)
	return nestedLevel(button.VariantSecondary, "Open Drawer", "Drawer", "Open another drawer from the same direction.", second)
}

const tallPlaceholder = "rounded-2xl bg-muted group-data-[swipe-axis=x]/drawer-popup:size-full group-data-[swipe-axis=y]/drawer-popup:h-80 group-data-[swipe-axis=y]/drawer-popup:w-full"

func BaseNonModal() h.Node {
	return bd.Drawer(bd.NonModal, bd.DisablePointerDismissal, bd.SwipeDirection("right"),
		bd.Trigger(button.Button(button.VariantOutline, "Non Modal")),
		bd.Content(
			bd.Header(bd.Title("Non Modal Drawer")),
			placeholder(tallPlaceholder),
			bd.Footer(bd.Close(button.Button("Close"))),
		),
	)
}

func BaseSnapPoints() h.Node {
	return bd.Drawer(bd.SnapPoints("31rem", 1), bd.ShowSwipeHandle,
		bd.Trigger(button.Button(button.VariantOutline, "Open Snap Drawer")),
		bd.Content(
			bd.Header(bd.Title("Snap points"), bd.Description("Drag the drawer to snap between a compact peek and a near full-height view.")),
			placeholder(tallPlaceholder),
			bd.Footer(bd.Close(button.Button("Close"))),
		),
	)
}

func BaseSwipeHandle() h.Node {
	return bd.Drawer(bd.ShowSwipeHandle,
		bd.Trigger(button.Button(button.VariantSecondary, "Open Drawer")),
		bd.Content(
			bd.Header(bd.Title("Drawer"), bd.Description("Drawer with a swipe handle.")),
			placeholder(tallPlaceholder),
			bd.Footer(bd.Close(button.Button("Close"))),
		),
	)
}
