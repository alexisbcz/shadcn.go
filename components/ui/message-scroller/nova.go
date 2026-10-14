package messagescroller

import (
	"io"

	"shadcn.go/components/style"
)

// The registry/bases/radix/ui/message-scroller.tsx classes, resolved per
// style. The parts render them inside style.Scope; outside, the
// new-york-v4 classes.
var file = style.NewFile("ui/message-scroller.tsx")

var (
	sRoot     = file.Cn("cn-message-scroller group/message-scroller relative flex size-full min-h-0 flex-col overflow-hidden", style.Arg)
	sViewport = file.Cn("cn-message-scroller-viewport size-full min-h-0 min-w-0 scroll-fade-b scrollbar-thin scrollbar-gutter-stable overflow-y-auto overscroll-contain contain-content data-autoscrolling:scrollbar-thumb-transparent data-autoscrolling:scrollbar-track-transparent data-pending-scroll:invisible", style.Arg)
	sContent  = file.Cn("cn-message-scroller-content flex h-max min-h-full flex-col", style.Arg)
	sItem     = file.Cn("cn-message-scroller-item min-w-0 shrink-0 [contain-intrinsic-size:auto_10rem] [content-visibility:auto]", style.Arg)
	sButton   = file.Cn("cn-message-scroller-button absolute inset-s-1/2 -translate-x-1/2 border-border bg-background text-foreground transition-[translate,scale,opacity] duration-200 hover:bg-muted hover:text-foreground data-[active=false]:pointer-events-none data-[active=false]:scale-95 data-[active=false]:opacity-0 data-[active=false]:duration-400 data-[active=false]:ease-[cubic-bezier(0.7,0,0.84,0)] data-[active=true]:translate-y-0 data-[active=true]:scale-100 data-[active=true]:opacity-100 data-[active=true]:ease-[cubic-bezier(0.23,1,0.32,1)] data-[direction=end]:bottom-4 data-[direction=end]:data-[active=false]:translate-y-full data-[direction=start]:top-4 data-[direction=start]:data-[active=false]:-translate-y-full rtl:translate-x-1/2 data-[direction=start]:[&_svg]:rotate-180", style.Arg)
)

// classIn is the new-york-v4 class ny, or site's class in the style in
// scope when rendering to w.
func classIn(w io.Writer, ny string, site *style.Site) string {
	if w != nil {
		if st, ok := style.FromWriter(w); ok {
			return site.Class(st, "")
		}
	}
	return ny
}
