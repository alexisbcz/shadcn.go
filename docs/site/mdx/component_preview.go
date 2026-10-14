package mdx

import (
	"strings"

	"shadcn.go/cn"
	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/popover"
	selectui "shadcn.go/components/ui/select"
	"shadcn.go/components/ui/separator"
	"shadcn.go/docs/content"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// ComponentPreview: ui/apps/v4/components/component-preview.tsx and
// component-preview-tabs.tsx. Props: name, styleName, type, className,
// previewClassName, align, hideCode, chromeLessOnMobile, direction, caption.
//
// The preview is the Go example rendered in radix-nova (the explicit
// "nova--<name>" registration if any, otherwise the example of the same name
// in style.Scope(style.Nova), like cmd/parity), whatever the page's base. The
// code panel shows the first 3 lines of the upstream source under a
// gradient with a "View Code" button; docs-mdx.js swaps in the full source
// (kept in a <template>) on click.
func ComponentPreview(env *Env, _ *content.Context, props content.Props, _ []h.Node) (h.Node, error) {
	name := props.String("name")
	styleName := props.String("styleName")
	if styleName == "" {
		styleName = "base-nova"
	}
	caption := props.String("caption")
	if props.String("type") == "block" {
		return blockPreview(name, styleName, caption), nil
	}
	hideCode, _ := props.Get("hideCode")
	hide := hideCode == true

	align := props.String("align")
	if align == "" {
		align = "center"
	}
	chromeless := "false"
	if v, _ := props.Get("chromeLessOnMobile"); v == true {
		chromeless = "true"
	}
	rtl := props.String("direction") == "rtl"

	var preview h.Node
	if ex, ok := LookupExample(name); ok {
		preview = ex(env)
	} else {
		preview = h.P(h.Class("text-sm text-muted-foreground"), "Example ", h.Code(name), " is not ported yet.")
	}
	dir := "ltr"
	if rtl {
		dir = "rtl"
	}
	previewDiv := h.Div(h.Data("slot", "preview"), h.Attr("dir", dir),
		h.Div(h.Data("align", align), h.Data("chromeless", chromeless),
			h.Class(cn.Cn("preview relative flex h-72 w-full justify-center p-10 data-[align=center]:items-center data-[align=end]:items-start data-[align=start]:items-start data-[chromeless=true]:h-auto data-[chromeless=true]:p-0 sm:data-[align=end]:items-end", props.String("previewClassName"))),
			preview,
		),
	)
	if rtl {
		previewDiv.Children(h.Data("lang", "ar"))
	}

	root := h.Div(h.Data("slot", "component-preview"), h.Data("not-typeset", "true"),
		h.Class(cn.Cn("group relative mt-4 mb-12 flex flex-col overflow-hidden rounded-2xl border", className(props))))
	if rtl {
		root.Children(rtlHeader())
	}
	root.Children(previewDiv)
	if !hide {
		root.Children(codePanel(name, styleName, rtl))
	}
	if caption == "" {
		return root, nil
	}
	hc := "false"
	if hide {
		hc = "true"
	}
	return h.El("figure").Data("hide-code", hc).Class("flex flex-col data-[hide-code=true]:gap-4").Children(
		root,
		h.El("figcaption").Class("-mt-8 text-center text-sm text-muted-foreground data-[hide-code=true]:mt-0").Children(caption),
	), nil
}

// codePanel is the data-slot="code" part of ComponentPreviewTabs.
func codePanel(name, styleName string, rtl bool) h.Node {
	full := []any{}
	if rtl {
		full = append(full, h.Div(h.Class("relative z-10 no-scrollbar overflow-x-auto border-t bg-code p-6 font-mono text-sm text-muted-foreground"),
			h.Pre("// You will notice this example uses dir and data-lang attributes.\n// This is because this site is not RTL by default.\n// In your application, you won't need these."),
			h.Span("// See the ", h.A(h.Href("/docs/rtl"), h.Class("underline underline-offset-4"), "RTL guide"), " for more information."),
		))
	}
	full = append(full, componentSource(name, "", styleName, "", "", "", false, 0))
	return h.Div(h.Data("slot", "code"), h.Data("mobile-code-visible", "false"),
		h.Class("relative overflow-hidden **:data-[slot=copy-button]:right-4 **:data-[slot=copy-button]:hidden data-[mobile-code-visible=true]:**:data-[slot=copy-button]:flex [&_[data-rehype-pretty-code-figure]]:m-0! [&_[data-rehype-pretty-code-figure]]:rounded-t-none [&_[data-rehype-pretty-code-figure]]:border-t [&_pre]:max-h-72"),
		h.Div(h.Class("relative"),
			componentSource(name, "", styleName, "", "", "", false, 3),
			h.Div(h.Class("absolute inset-0 flex items-center justify-center pb-4"),
				h.Div(h.Class("absolute inset-0"), h.Style("background:linear-gradient(to top, var(--color-code), color-mix(in oklab, var(--color-code) 60%, transparent), transparent)")),
				button.Button("View Code").Type("button").Size(button.SizeSm).Variant(button.VariantOutline).
					Class("relative z-10 rounded-lg bg-background text-foreground shadow-none hover:bg-muted dark:bg-background dark:text-foreground dark:hover:bg-muted").
					Data("docs-view-code", ""),
			),
		),
		h.El("template").Data("docs-full-code", "").Children(full...),
	)
}

// rtlHeader is the language bar of RTL previews (LanguageSelector and the
// translation notice popover). The language stays Arabic: the Go RTL
// examples are rendered in Arabic only.
func rtlHeader() h.Node {
	return h.Div(h.Class("flex h-16 items-center border-b px-4"),
		style.Scope(style.Nova,
			selectui.Select(selectui.DefaultValue("ar"),
				selectui.Trigger(selectui.Value()).Size(selectui.SizeSm).Class("w-36").
					Attr("dir", "ltr").Data("name", "language-selector"),
				selectui.Content(h.Attr("dir", "ltr"), h.Class("data-open:animate-none data-closed:animate-none"),
					selectui.Group(
						selectui.Item("English").Value("en"),
						selectui.Item("Arabic (العربية)").Value("ar"),
						selectui.Item("Hebrew (עברית)").Value("he"),
					),
				),
			),
			popover.Popover(
				popover.Trigger(lib.AsChild,
					button.Button(iconAlertCircle(), h.Span(h.Class("sr-only"), "Toggle")).
						Variant(button.VariantGhost).Size(button.SizeIconSm).Class("ml-auto size-7"),
				),
				popover.Content(popover.SideBottom, popover.AlignEnd, h.Class("w-56 text-xs"),
					h.Div("I used AI to translate the text for demonstration purposes. It's not perfect and may contain errors."),
					separator.Separator().Class("-mx-2.5 w-auto!"),
					h.Div(h.Data("lang", "ar"), "لقد استخدمت الذكاء الاصطناعي لترجمة النص للأغراض التجريبية فقط. قد لا تكون الترجمة دقيقة وقد تحتوي على أخطاء."),
					separator.Separator().Class("-mx-2.5 w-auto!"),
					h.Div(h.Data("lang", "he"), "השתמשתי בבינה מלאכותית כדי לתרגם את הטקסט למטרות הדגמה. זה לא מושלם ויכול להכיל שגיאות."),
				),
			),
		),
	)
}

// blockPreview is the type="block" branch: screenshots on mobile, the
// /view iframe from md up.
func blockPreview(name, styleName, caption string) h.Node {
	img := func(theme, class string) h.Node {
		return h.Img(h.Alt(name), h.Attr("width", "1440"), h.Attr("height", "900"),
			h.Attr("decoding", "async"), h.Data("nimg", "1"), h.Class(class), h.Attr("loading", "lazy"),
			h.Style("color:transparent"), h.Src("/r/styles/new-york/"+name+"-"+theme+".png"))
	}
	box := h.Div(h.Data("not-typeset", "true"),
		h.Class("relative mt-6 aspect-[4/2.5] w-full overflow-hidden rounded-2xl border md:-mx-1"),
		img("light", "absolute top-0 left-0 z-20 h-full w-[1600px] max-w-none bg-background object-cover object-left-top md:hidden dark:hidden md:dark:hidden"),
		img("dark", "absolute top-0 left-0 z-20 hidden h-full w-[1600px] max-w-none bg-background object-cover object-left-top md:hidden dark:block md:dark:hidden"),
		h.Div(h.Class("absolute inset-0 hidden w-[1600px] bg-background md:block"),
			h.El("iframe").Attr("src", "/view/"+styleName+"/"+name).Class("size-full")),
	)
	if caption == "" {
		return box
	}
	return h.El("figure").Class("flex flex-col gap-4").Children(box,
		h.El("figcaption").Class("text-center text-sm text-muted-foreground").Children(caption))
}

// LookupExample finds the Go port of a docs example by its upstream name
// (ui/apps/v4/examples/radix/<name>.tsx) and returns a function rendering
// it in radix-nova.
func LookupExample(name string) (func(env *Env) h.Node, bool) {
	if ex, ok := registry.Get("nova--" + name); ok {
		return func(env *Env) h.Node { return ex.Render(env.Request) }, true
	}
	ex, ok := registry.Get(strings.TrimPrefix(name, "nova--"))
	if !ok {
		return nil, false
	}
	return func(env *Env) h.Node { return style.Scope(style.Nova, ex.Render(env.Request)) }, true
}

func init() { Register("ComponentPreview", ComponentPreview) }
