package contextmenu

import (
	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/menu"
	"shadcn.go/components/ui/popper"
	h "shadcn.go/libhtml"
)

// The classes of registry/bases/radix/ui/context-menu.tsx, used when a
// style is in scope (style.Scope).
var (
	baseFile         = style.NewFile("ui/context-menu.tsx")
	baseTrigger      = baseFile.Cn("cn-context-menu-trigger select-none", style.Arg)
	baseContent      = baseFile.Cn("cn-context-menu-content cn-menu-target cn-menu-translucent z-50 max-h-(--radix-context-menu-content-available-height) origin-(--radix-context-menu-content-transform-origin) overflow-x-hidden overflow-y-auto", style.Arg)
	baseItem         = baseFile.Cn("cn-context-menu-item group/context-menu-item relative flex cursor-default items-center outline-hidden select-none data-disabled:pointer-events-none data-disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:shrink-0", style.Arg)
	baseSubTrigger   = baseFile.Cn("cn-context-menu-sub-trigger flex cursor-default items-center outline-hidden select-none [&_svg]:pointer-events-none [&_svg]:shrink-0", style.Arg)
	baseChevron      = baseFile.Attr("cn-rtl-flip ml-auto")
	baseSubContent   = baseFile.Cn("cn-context-menu-sub-content cn-menu-target cn-menu-translucent z-50 origin-(--radix-context-menu-content-transform-origin) overflow-hidden", style.Arg)
	baseCheckboxItem = baseFile.Cn("cn-context-menu-checkbox-item relative flex cursor-default items-center outline-hidden select-none data-disabled:pointer-events-none data-disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:shrink-0", style.Arg)
	baseIndicator    = baseFile.Attr("cn-context-menu-item-indicator pointer-events-none")
	baseRadioItem    = baseFile.Cn("cn-context-menu-radio-item relative flex cursor-default items-center outline-hidden select-none data-disabled:pointer-events-none data-disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:shrink-0", style.Arg)
	baseLabel        = baseFile.Cn("cn-context-menu-label", style.Arg)
	baseSeparator    = baseFile.Cn("cn-context-menu-separator", style.Arg)
	baseShortcut     = baseFile.Cn("cn-context-menu-shortcut", style.Arg)
)

// styledKind is ContextMenu in style st: the markup of the bases source.
var styledKind = menu.Cached(func(st style.Style) *menu.Kind {
	ind := func() h.Node {
		return h.Span(h.Class(baseIndicator.Class(st)), menu.Indicator(lucide.Check()))
	}
	return &menu.Kind{
		Name:              "context-menu",
		Content:           baseContent.Class(st, ""),
		SubContent:        baseSubContent.Class(st, ""),
		Item:              baseItem.Class(st, ""),
		CheckboxItem:      baseCheckboxItem.Class(st, ""),
		RadioItem:         baseRadioItem.Class(st, ""),
		Label:             baseLabel.Class(st, ""),
		Separator:         baseSeparator.Class(st, ""),
		Shortcut:          baseShortcut.Class(st, ""),
		SubTrigger:        baseSubTrigger.Class(st, ""),
		CheckboxIndicator: ind,
		RadioIndicator:    ind,
		SubTriggerChevron: func() h.Node { return lucide.ChevronRight(h.Class(baseChevron.Class(st))) },
		ContentForced: func(o *popper.Options) {
			o.Side = popper.SideRight
			o.SideOffset = 2
			o.Align = popper.AlignStart
		},
	}
})
