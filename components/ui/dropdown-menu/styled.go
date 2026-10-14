package dropdownmenu

import (
	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/menu"
	"shadcn.go/components/ui/popper"
	h "shadcn.go/libhtml"
)

// The classes of registry/bases/radix/ui/dropdown-menu.tsx, used when a
// style is in scope (style.Scope).
var (
	baseFile         = style.NewFile("ui/dropdown-menu.tsx")
	baseContent      = baseFile.Cn("cn-dropdown-menu-content cn-menu-target cn-menu-translucent z-50 max-h-(--radix-dropdown-menu-content-available-height) w-(--radix-dropdown-menu-trigger-width) origin-(--radix-dropdown-menu-content-transform-origin) overflow-x-hidden overflow-y-auto data-[state=closed]:overflow-hidden", style.Arg)
	baseItem         = baseFile.Cn("cn-dropdown-menu-item group/dropdown-menu-item relative flex cursor-default items-center outline-hidden select-none data-disabled:pointer-events-none data-disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:shrink-0", style.Arg)
	baseCheckboxItem = baseFile.Cn("cn-dropdown-menu-checkbox-item relative flex cursor-default items-center outline-hidden select-none data-disabled:pointer-events-none data-disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:shrink-0", style.Arg)
	baseIndicator    = baseFile.Attr("cn-dropdown-menu-item-indicator pointer-events-none")
	baseRadioItem    = baseFile.Cn("cn-dropdown-menu-radio-item relative flex cursor-default items-center outline-hidden select-none data-disabled:pointer-events-none data-disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:shrink-0", style.Arg)
	baseLabel        = baseFile.Cn("cn-dropdown-menu-label", style.Arg)
	baseSeparator    = baseFile.Cn("cn-dropdown-menu-separator", style.Arg)
	baseShortcut     = baseFile.Cn("cn-dropdown-menu-shortcut", style.Arg)
	baseSubTrigger   = baseFile.Cn("cn-dropdown-menu-sub-trigger flex cursor-default items-center outline-hidden select-none [&_svg]:pointer-events-none [&_svg]:shrink-0", style.Arg)
	baseChevron      = baseFile.Attr("cn-rtl-flip ml-auto")
	baseSubContent   = baseFile.Cn("cn-dropdown-menu-sub-content cn-menu-target cn-menu-translucent z-50 origin-(--radix-dropdown-menu-content-transform-origin) overflow-hidden", style.Arg)
)

// styledKind is DropdownMenu in style st: the markup of the bases source.
var styledKind = menu.Cached(func(st style.Style, cls func(string) string) *menu.Kind {
	ind := func(slot string) func() h.Node {
		return func() h.Node {
			return h.Span(h.Class(cls(baseIndicator.Class(st))), h.Data("slot", slot), menu.Indicator(lucide.Check()))
		}
	}
	return &menu.Kind{
		Name:              "dropdown-menu",
		Content:           cls(baseContent.Class(st, "")),
		SubContent:        cls(baseSubContent.Class(st, "")),
		Item:              cls(baseItem.Class(st, "")),
		CheckboxItem:      cls(baseCheckboxItem.Class(st, "")),
		RadioItem:         cls(baseRadioItem.Class(st, "")),
		Label:             cls(baseLabel.Class(st, "")),
		Separator:         cls(baseSeparator.Class(st, "")),
		Shortcut:          cls(baseShortcut.Class(st, "")),
		SubTrigger:        cls(baseSubTrigger.Class(st, "")),
		CheckboxIndicator: ind("dropdown-menu-checkbox-item-indicator"),
		RadioIndicator:    ind("dropdown-menu-radio-item-indicator"),
		SubTriggerChevron: func() h.Node { return lucide.ChevronRight(h.Class(cls(baseChevron.Class(st)))) },
		ContentDefaults: func(o *popper.Options) {
			o.Align = popper.AlignStart
			o.SideOffset = 4
		},
		ContentIDs: true,
	}
})
