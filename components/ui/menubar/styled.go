package menubar

import (
	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/menu"
	"shadcn.go/components/ui/popper"
	h "shadcn.go/libhtml"
)

// The classes of registry/bases/radix/ui/menubar.tsx, used when a style is
// in scope (style.Scope).
var (
	baseFile         = style.NewFile("ui/menubar.tsx")
	baseRoot         = baseFile.Cn("cn-menubar flex items-center", style.Arg)
	baseTrigger      = baseFile.Cn("cn-menubar-trigger flex items-center outline-hidden select-none", style.Arg)
	baseContent      = baseFile.Cn("cn-menubar-content cn-menu-target cn-menu-translucent z-50 origin-(--radix-menubar-content-transform-origin) overflow-hidden", style.Arg)
	baseItem         = baseFile.Cn("cn-menubar-item group/menubar-item relative flex cursor-default items-center outline-hidden select-none data-disabled:pointer-events-none [&_svg]:pointer-events-none [&_svg]:shrink-0", style.Arg)
	baseCheckboxItem = baseFile.Cn("cn-menubar-checkbox-item relative flex cursor-default items-center outline-hidden select-none data-disabled:pointer-events-none [&_svg]:pointer-events-none [&_svg]:shrink-0", style.Arg)
	baseCheckboxInd  = baseFile.Attr("cn-menubar-checkbox-item-indicator pointer-events-none absolute flex items-center justify-center")
	baseRadioItem    = baseFile.Cn("cn-menubar-radio-item relative flex cursor-default items-center outline-hidden select-none data-disabled:pointer-events-none [&_svg]:pointer-events-none [&_svg]:shrink-0", style.Arg)
	baseRadioInd     = baseFile.Attr("cn-menubar-radio-item-indicator pointer-events-none absolute flex items-center justify-center")
	baseLabel        = baseFile.Cn("cn-menubar-label", style.Arg)
	baseSeparator    = baseFile.Cn("cn-menubar-separator -mx-1 my-1 h-px", style.Arg)
	baseShortcut     = baseFile.Cn("cn-menubar-shortcut ml-auto", style.Arg)
	baseSubTrigger   = baseFile.Cn("cn-menubar-sub-trigger flex cursor-default items-center outline-none select-none", style.Arg)
	baseChevron      = baseFile.Attr("cn-rtl-flip ml-auto size-4")
	baseSubContent   = baseFile.Cn("cn-menubar-sub-content cn-menu-target cn-menu-translucent z-50 origin-(--radix-menubar-content-transform-origin) overflow-hidden", style.Arg)
)

// styledKind is Menubar in style st: the markup of the bases source.
var styledKind = menu.Cached(func(st style.Style) *menu.Kind {
	return &menu.Kind{
		Name:         "menubar",
		Content:      baseContent.Class(st, ""),
		SubContent:   baseSubContent.Class(st, ""),
		Item:         baseItem.Class(st, ""),
		CheckboxItem: baseCheckboxItem.Class(st, ""),
		RadioItem:    baseRadioItem.Class(st, ""),
		Label:        baseLabel.Class(st, ""),
		Separator:    baseSeparator.Class(st, ""),
		Shortcut:     baseShortcut.Class(st, ""),
		SubTrigger:   baseSubTrigger.Class(st, ""),
		CheckboxIndicator: func() h.Node {
			return h.Span(h.Class(baseCheckboxInd.Class(st)), menu.Indicator(lucide.Check()))
		},
		RadioIndicator: func() h.Node {
			return h.Span(h.Class(baseRadioInd.Class(st)), menu.Indicator(lucide.Check()))
		},
		SubTriggerChevron: func() h.Node { return lucide.ChevronRight(h.Class(baseChevron.Class(st))) },
		ContentDefaults: func(o *popper.Options) {
			o.Align = popper.AlignStart
			o.AlignOffset = -4
			o.SideOffset = 8
		},
		ContentAttrs:    []h.Prop{h.Data("radix-menubar-content", "")},
		SubTriggerAttrs: []h.Prop{h.Data("radix-menubar-subtrigger", "")},
		ContentIDs:      true,
	}
})
