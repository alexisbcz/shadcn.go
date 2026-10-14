package sidebar

import (
	"io"

	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	"shadcn.go/cva"
)

// The registry/bases/radix/ui/sidebar.tsx classes, resolved per style. A
// part renders them when it renders inside style.Scope; outside, the
// new-york-v4 classes.
var file = style.NewFile("ui/sidebar.tsx")

var (
	sGap = file.Cn(
		"cn-sidebar-gap relative w-(--sidebar-width) bg-transparent",
		"group-data-[collapsible=offcanvas]:w-0",
		"group-data-[side=right]:rotate-180",
		style.Arg,
	)
	sContainer = file.Cn(
		"fixed inset-y-0 z-10 hidden h-svh w-(--sidebar-width) transition-[left,right,width] duration-200 ease-linear data-[side=left]:left-0 data-[side=left]:group-data-[collapsible=offcanvas]:left-[calc(var(--sidebar-width)*-1)] data-[side=right]:right-0 data-[side=right]:group-data-[collapsible=offcanvas]:right-[calc(var(--sidebar-width)*-1)] md:flex",
		style.Arg, style.Arg,
	)
	sInner   = file.Attr("cn-sidebar-inner flex size-full flex-col")
	sTrigger = file.Cn("cn-sidebar-trigger", style.Arg)
	sRail    = file.Cn(
		"cn-sidebar-rail absolute inset-y-0 z-20 hidden w-4 transition-all ease-linear group-data-[side=left]:-right-4 group-data-[side=right]:left-0 after:absolute after:inset-y-0 after:start-1/2 after:w-[2px] sm:flex ltr:-translate-x-1/2 rtl:-translate-x-1/2",
		"in-data-[side=left]:cursor-w-resize in-data-[side=right]:cursor-e-resize",
		"[[data-side=left][data-state=collapsed]_&]:cursor-e-resize [[data-side=right][data-state=collapsed]_&]:cursor-w-resize",
		"group-data-[collapsible=offcanvas]:translate-x-0 group-data-[collapsible=offcanvas]:after:left-full hover:group-data-[collapsible=offcanvas]:bg-sidebar",
		"[[data-side=left][data-collapsible=offcanvas]_&]:-right-2",
		"[[data-side=right][data-collapsible=offcanvas]_&]:-left-2",
		style.Arg,
	)
	sInset        = file.Cn("cn-sidebar-inset relative flex w-full flex-1 flex-col", style.Arg)
	sInput        = file.Cn("cn-sidebar-input", style.Arg)
	sHeader       = file.Cn("cn-sidebar-header flex flex-col", style.Arg)
	sFooter       = file.Cn("cn-sidebar-footer flex flex-col", style.Arg)
	sSeparator    = file.Cn("cn-sidebar-separator w-auto", style.Arg)
	sContent      = file.Cn("cn-sidebar-content flex min-h-0 flex-1 flex-col overflow-auto group-data-[collapsible=icon]:overflow-hidden", style.Arg)
	sGroup        = file.Cn("cn-sidebar-group relative flex w-full min-w-0 flex-col", style.Arg)
	sGroupLabel   = file.Cn("cn-sidebar-group-label flex shrink-0 items-center outline-hidden [&>svg]:shrink-0", style.Arg)
	sGroupAction  = file.Cn("cn-sidebar-group-action flex aspect-square items-center justify-center outline-hidden transition-transform group-data-[collapsible=icon]:hidden after:absolute after:-inset-2 md:after:hidden [&>svg]:shrink-0", style.Arg)
	sGroupContent = file.Cn("cn-sidebar-group-content w-full", style.Arg)
	sMenu         = file.Cn("cn-sidebar-menu flex w-full min-w-0 flex-col", style.Arg)
	sMenuItem     = file.Cn("group/menu-item relative", style.Arg)
	sMenuAction   = file.Cn(
		"cn-sidebar-menu-action flex items-center justify-center outline-hidden transition-transform group-data-[collapsible=icon]:hidden after:absolute after:-inset-2 md:after:hidden [&>svg]:shrink-0",
		style.Arg, style.Arg,
	)
	sMenuBadge        = file.Cn("cn-sidebar-menu-badge flex items-center justify-center tabular-nums select-none group-data-[collapsible=icon]:hidden", style.Arg)
	sMenuSkeleton     = file.Cn("cn-sidebar-menu-skeleton flex items-center", style.Arg)
	sMenuSkeletonIcon = file.Attr("cn-sidebar-menu-skeleton-icon")
	sMenuSkeletonText = file.Attr("cn-sidebar-menu-skeleton-text max-w-(--skeleton-width) flex-1")
	sMenuSub          = file.Cn("cn-sidebar-menu-sub flex min-w-0 flex-col", style.Arg)
	sMenuSubItem      = file.Cn("group/menu-sub-item relative", style.Arg)
	sMenuSubButton    = file.Cn("cn-sidebar-menu-sub-button flex min-w-0 -translate-x-px items-center overflow-hidden outline-hidden group-data-[collapsible=icon]:hidden disabled:pointer-events-none disabled:opacity-50 aria-disabled:pointer-events-none aria-disabled:opacity-50 [&>span:last-child]:truncate [&>svg]:shrink-0", style.Arg)

	sMenuButton = file.CVA(
		"cn-sidebar-menu-button peer/menu-button group/menu-button flex w-full items-center overflow-hidden outline-hidden disabled:pointer-events-none disabled:opacity-50 aria-disabled:pointer-events-none aria-disabled:opacity-50 [&_svg]:size-4 [&_svg]:shrink-0 [&>span:last-child]:truncate",
		[]style.Variant{
			{Name: "variant", Options: []style.Option{
				{Name: "default", Class: "cn-sidebar-menu-button-variant-default"},
				{Name: "outline", Class: "cn-sidebar-menu-button-variant-outline"},
			}},
			{Name: "size", Options: []style.Option{
				{Name: "default", Class: "cn-sidebar-menu-button-size-default"},
				{Name: "sm", Class: "cn-sidebar-menu-button-size-sm"},
				{Name: "lg", Class: "cn-sidebar-menu-button-size-lg"},
			}},
		},
		cva.Selection{"variant": "default", "size": "default"},
	)
)

const menuActionHoverClass = "group-focus-within/menu-item:opacity-100 group-hover/menu-item:opacity-100 peer-data-active/menu-button:text-sidebar-accent-foreground aria-expanded:opacity-100 md:opacity-0"

// styleOf returns the style in scope when rendering to w ("" outside a
// style.Scope, or for a nil w).
func styleOf(w io.Writer) style.Style {
	if w == nil {
		return ""
	}
	st, _ := style.FromWriter(w)
	return st
}


func menuButtonClass(st style.Style, v ButtonVariant, s ButtonSize, className string) string {
	if st == "" {
		return MenuButtonVariants(v, s, className)
	}
	return lib.Cn(sMenuButton.Classes(st, cva.Selection{"variant": string(v), "size": string(s)}), className)
}
