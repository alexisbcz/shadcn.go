package mdx

import (
	_ "embed"
	"encoding/json"
	"strings"
	"sync"

	h "shadcn.go/libhtml"
)

// go-examples.json is the Go source of every docs example, generated from
// examples/ by gen/goexamples: the code panels show it in their Go tab.
//
//go:embed go-examples.json
var goExamplesJSON []byte

var goExamples = sync.OnceValue(func() map[string]string {
	m := map[string]string{}
	_ = json.Unmarshal(goExamplesJSON, &m)
	return m
})

// GoExample returns the Go source of an example as a docs style shows it,
// looked up like the preview (base--/aria--, then nova--, then the plain
// example rendered in the style).
func GoExample(name, styleName string) (string, bool) {
	m := goExamples()
	var prefixes []string
	switch {
	case strings.HasPrefix(styleName, "aria-"):
		prefixes = []string{"aria--", "nova--"}
	case strings.HasPrefix(styleName, "base-"):
		prefixes = []string{"base--", "nova--"}
	case styleName != "new-york-v4" && styleName != "":
		prefixes = []string{"nova--"}
	}
	for _, p := range prefixes {
		if code, ok := m[p+name]; ok {
			return code, true
		}
	}
	code, ok := m[name]
	if !ok {
		return "", false
	}
	if len(prefixes) > 0 {
		// The plain example, which the docs render in the nova style.
		code += "\n\n// The docs render this in the nova style:\n// wrap it in style.Scope(style.Nova, ...) (github.com/alexisbcz/components/style)."
	}
	return code, true
}

// goSource is ComponentSource for Go code.
func goSource(code string, maxLines int) h.Node {
	if maxLines > 0 {
		code = firstLines(code, maxLines)
	}
	return h.Div(h.Class("relative"), componentCode(code, "go", ""))
}

// langSwitch toggles the code panels between React and Go (site.js keeps
// the choice for every panel and page).
func langSwitch() h.Node {
	btn := func(lang, label string) h.Node {
		return h.Button(h.Type("button"), h.Data("code-lang-switch", lang), h.Aria("label", "Show "+label+" code"),
			h.Class("rounded-md px-2 py-0.5 text-xs font-medium text-muted-foreground transition-colors hover:text-foreground"),
			label)
	}
	return h.Div(h.Data("slot", "code-lang"),
		h.Class("absolute top-2 right-14 z-20 flex items-center gap-0.5 rounded-lg border bg-background p-0.5"),
		btn("react", "React"), btn("go", "Go"))
}
