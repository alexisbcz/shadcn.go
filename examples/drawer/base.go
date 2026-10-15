package drawer

import (
	"shadcn.go/components/style"
	"shadcn.go/components/ui/button"
	bd "shadcn.go/components/ui/drawer/basedrawer"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// Docs previews of the base and aria drawer pages that only exist there
// (ui/apps/v4/examples/{base,aria}/drawer-{nested,non-modal,snap-points,
// swipe-handle}.tsx): Base UI drawers, in base-rhea for the base page and
// aria-nova for the aria page. The base page previews them with
// styleName="base-rhea", which has no prefix of its own: they register as
// nova-- names (no radix example has these names).
func init() {
	for _, v := range []struct {
		prefix string
		st     style.Style
	}{{"nova--", style.Rhea}, {"aria--", style.Nova}} {
		st := v.st
		reg := func(name string, fn func() h.Node) {
			registry.Register(registry.Example{Name: v.prefix + name, Component: "drawer", Render: registry.Static(func() h.Node {
				return style.Scope(st, fn())
			})})
		}
		reg("drawer-nested", BaseNested)
		reg("drawer-non-modal", BaseNonModal)
		reg("drawer-snap-points", BaseSnapPoints)
		reg("drawer-swipe-handle", BaseSwipeHandle)
	}
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
