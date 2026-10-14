package mdx

import (
	"strconv"
	"strings"

	"shadcn.go/cn"
	"shadcn.go/docs/content"
	h "shadcn.go/libhtml"
)

// ComponentSource: ui/apps/v4/components/component-source.tsx. The
// highlighted source of name (the demo of that name, or else the registry
// item's first file) or src in styleName, formatted like upstream (see
// sources.go). Props: name, src, title, language, collapsible (default
// true), className, styleName (default "base-nova"), maxLines.
func ComponentSource(_ *Env, _ *content.Context, props content.Props, _ []h.Node) (h.Node, error) {
	collapsible := true
	if v, ok := props.Get("collapsible"); ok && v == false {
		collapsible = false
	}
	maxLines := 0
	if v, ok := props.Get("maxLines"); ok {
		switch n := v.(type) {
		case float64:
			maxLines = int(n)
		case int:
			maxLines = n
		case string:
			maxLines, _ = strconv.Atoi(n)
		}
	}
	return componentSource(props.String("name"), props.String("src"), props.String("styleName"),
		props.String("title"), props.String("language"), className(props), collapsible, maxLines), nil
}

// componentSource renders ComponentSource; nil when there is no code (like
// upstream, which renders nothing).
func componentSource(name, src, styleName, title, language, class string, collapsible bool, maxLines int) h.Node {
	var code string
	var ok bool
	switch {
	case src != "":
		code, ok = SourceFile(styleName, src)
	case name != "":
		code, ok = Source(styleName, name)
	}
	if !ok || code == "" {
		return h.Group()
	}
	if maxLines > 0 {
		code = firstLines(code, maxLines)
	}
	lang := language
	if lang == "" && title != "" {
		lang = title[strings.LastIndex(title, ".")+1:]
	}
	if lang == "" {
		lang = "tsx"
	}
	figure := componentCode(code, lang, title)
	if !collapsible {
		return h.Div(h.Class(cn.Cn("relative", class)), figure)
	}
	return codeCollapsibleWrapper(class, figure)
}

// componentCode is ComponentCode in component-source.tsx.
func componentCode(code, lang, title string) h.Node {
	fig := h.El("figure").Attr("data-rehype-pretty-code-figure", "").Class("[&>pre]:max-h-96")
	if title != "" {
		fig.Children(h.El("figcaption").Attr("data-rehype-pretty-code-title", "").
			Class("flex items-center gap-2 text-code-foreground [&_svg]:size-4 [&_svg]:text-code-foreground [&_svg]:opacity-70").
			Attr("data-language", lang).
			Children(content.LanguageIcon(lang), title))
	}
	return fig.Children(
		copyButton("", ""),
		h.Div(h.Data("not-typeset", "true"), highlightCode(code, lang)),
	)
}

func init() { Register("ComponentSource", ComponentSource) }
