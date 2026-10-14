package highlight

import (
	"shadcn.go/docs/hast"
)

// codePreClass is the class lib/highlight-code.ts's transformer gives <pre>.
const codePreClass = "no-scrollbar min-w-0 overflow-x-auto overflow-y-auto overscroll-x-contain overscroll-y-auto px-4 py-3.5 outline-none has-[[data-highlighted-line]]:px-0 has-[[data-line-numbers]]:px-0 has-[[data-slot=tabs]]:p-0 !bg-transparent"

// CodeTree highlights code like highlightCode in
// ui/apps/v4/lib/highlight-code.ts, which the docs components use for the
// source code they show (ComponentSource, the block and chart viewers):
// shiki's codeToHtml with the github-dark and github-light themes and the
// docs' transformer (a fixed class on <pre>, data-line-numbers on <code>,
// data-line on every line). It returns the <pre> element.
func (h *Highlighter) CodeTree(code, lang string) (*hast.Element, error) {
	if lang == "" {
		lang = "tsx"
	}
	pre, err := h.codeToHast(code, lang, true)
	if err != nil {
		return nil, err
	}
	pre.Set("class", codePreClass)
	if c, ok := firstChild(pre); ok && c.Tag == "code" {
		c.Set("data-line-numbers", "")
		for _, l := range c.Children {
			if e, ok := l.(*hast.Element); ok && e.Tag == "span" {
				if cls := e.Classes(); len(cls) > 0 && cls[0] == "line" {
					e.Set("data-line", "")
				}
			}
		}
	}
	return pre, nil
}

// Code returns the HTML string highlightCode returns (the docs insert it
// with dangerouslySetInnerHTML).
func (h *Highlighter) Code(code, lang string) (string, error) {
	pre, err := h.CodeTree(code, lang)
	if err != nil {
		return "", err
	}
	return hast.ToHTML(pre), nil
}
