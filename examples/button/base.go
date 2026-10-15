package button

import (
	"shadcn.go/components/icons/tabler"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/button"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// button-render (ui/apps/v4/examples/{base,aria}/button-render.tsx): a link
// with className={buttonVariants({ variant: "secondary", size: "sm" })} of
// base-nova / aria-nova. buttonVariants is plain cva (no tailwind-merge),
// so its classes are kept as they are, conflicts included.
func init() {
	for _, prefix := range []string{"base--", "aria--"} {
		registry.Register(registry.Example{Name: prefix + "button-render", Component: "button", Render: registry.Static(Render)})
		registry.Register(registry.Example{Name: prefix + "button-with-icon", Component: "button", Render: registry.Static(func() h.Node {
			return style.Scope(style.Nova, BaseWithIcon())
		})})
	}
}

const renderClass = "group/button inline-flex shrink-0 items-center justify-center rounded-lg border border-transparent bg-clip-padding text-sm font-medium whitespace-nowrap transition-all outline-none select-none focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/50 active:not-aria-[haspopup]:translate-y-px disabled:pointer-events-none disabled:opacity-50 aria-invalid:border-destructive aria-invalid:ring-3 aria-invalid:ring-destructive/20 dark:aria-invalid:border-destructive/50 dark:aria-invalid:ring-destructive/40 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4 bg-secondary text-secondary-foreground hover:bg-[color-mix(in_oklch,var(--secondary),var(--foreground)_5%)] aria-expanded:bg-secondary aria-expanded:text-secondary-foreground h-7 gap-1 rounded-[min(var(--radius-md),12px)] px-2.5 text-[0.8rem] in-data-[slot=button-group]:rounded-lg has-data-[icon=inline-end]:pr-1.5 has-data-[icon=inline-start]:pl-1.5 [&_svg:not([class*='size-'])]:size-3.5"

func Render() h.Node {
	return h.A(h.Href("#"), h.Class(renderClass), "Login")
}

// BaseWithIcon is the base/aria button-with-icon: icons on either side
// (data-icon="inline-start" / "inline-end"); the radix one has one icon.
func BaseWithIcon() h.Node {
	return h.Div(h.Class("flex gap-2"),
		button.Button(button.VariantOutline, tabler.IconGitBranch(h.Data("icon", "inline-start")), " New Branch"),
		button.Button(button.VariantOutline, "Fork", iconGitFork(h.Data("icon", "inline-end"))),
	)
}

// iconGitFork is tabler's git-fork (outline, @tabler/icons-react 3.34.1),
// not in the generated tabler set.
func iconGitFork(children ...any) h.Node {
	return h.Svg(h.Xmlns("http://www.w3.org/2000/svg"), h.Width("24"), h.Height("24"), h.ViewBox("0 0 24 24"),
		h.Fill("none"), h.Stroke("currentColor"), h.StrokeWidth("2"), h.StrokeLinecap("round"), h.StrokeLinejoin("round"),
		h.Class("tabler-icon tabler-icon-git-fork"),
		h.UnsafeHTML(`<path d="M12 18m-2 0a2 2 0 1 0 4 0a2 2 0 1 0 -4 0"></path><path d="M7 6m-2 0a2 2 0 1 0 4 0a2 2 0 1 0 -4 0"></path><path d="M17 6m-2 0a2 2 0 1 0 4 0a2 2 0 1 0 -4 0"></path><path d="M7 8v2a2 2 0 0 0 2 2h6a2 2 0 0 0 2 -2v-2"></path><path d="M12 12l0 4"></path>`),
		h.Group(children...))
}
