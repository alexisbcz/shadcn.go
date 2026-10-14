package style

import (
	"io"
	"strings"

	h "shadcn.go/libhtml"
)

// Upstream generates, for radix-nova and base-nova, a right-to-left variant
// of every component (styles/<base>-<style>/ui-rtl) with transformDirection
// (packages/registry/src/transformers/transform-rtl.ts): physical classes
// (pl-2, text-left, right-0, ...) become logical ones (ps-2, text-start,
// end-0, ...). The docs' *-rtl examples import those components. ScopeRTL
// selects them for a subtree; style-aware components map their own classes
// with RTLClass when RTLFromWriter reports true.

type rtlKey struct{}

// ScopeRTL renders children with the ui-rtl variant of the components.
func ScopeRTL(children ...any) h.Node { return h.WithValue(rtlKey{}, true, children...) }

// RTLFromWriter reports whether an enclosing ScopeRTL is in effect when
// rendering to w (the writer libhtml passes to Render).
func RTLFromWriter(w io.Writer) bool {
	v, ok := h.Lookup(w, rtlKey{})
	return ok && v == true
}

// rtlMappings is RTL_MAPPINGS, in upstream's order.
var rtlMappings = [][2]string{
	{"-ml-", "-ms-"},
	{"-mr-", "-me-"},
	{"ml-", "ms-"},
	{"mr-", "me-"},
	{"pl-", "ps-"},
	{"pr-", "pe-"},
	{"-left-", "-start-"},
	{"-right-", "-end-"},
	{"left-", "start-"},
	{"right-", "end-"},
	{"inset-l-", "inset-inline-start-"},
	{"inset-r-", "inset-inline-end-"},
	{"rounded-tl-", "rounded-ss-"},
	{"rounded-tr-", "rounded-se-"},
	{"rounded-bl-", "rounded-es-"},
	{"rounded-br-", "rounded-ee-"},
	{"rounded-l-", "rounded-s-"},
	{"rounded-r-", "rounded-e-"},
	{"border-l-", "border-s-"},
	{"border-r-", "border-e-"},
	{"border-l", "border-s"},
	{"border-r", "border-e"},
	{"text-left", "text-start"},
	{"text-right", "text-end"},
	{"scroll-ml-", "scroll-ms-"},
	{"scroll-mr-", "scroll-me-"},
	{"scroll-pl-", "scroll-ps-"},
	{"scroll-pr-", "scroll-pe-"},
	{"float-left", "float-start"},
	{"float-right", "float-end"},
	{"clear-left", "clear-start"},
	{"clear-right", "clear-end"},
	{"origin-top-left", "origin-top-start"},
	{"origin-top-right", "origin-top-end"},
	{"origin-bottom-left", "origin-bottom-start"},
	{"origin-bottom-right", "origin-bottom-end"},
	{"origin-left", "origin-start"},
	{"origin-right", "origin-end"},
}

var rtlSlideMappings = [][3]string{
	{"data-[side=inline-start]", "slide-in-from-right", "slide-in-from-end"},
	{"data-[side=inline-start]", "slide-out-to-right", "slide-out-to-end"},
	{"data-[side=inline-end]", "slide-in-from-left", "slide-in-from-start"},
	{"data-[side=inline-end]", "slide-out-to-left", "slide-out-to-start"},
}

// splitClassName is splitClassName of transform-css-vars.ts: variant, name
// and alpha modifier.
func splitClassName(c string) (variant, value, modifier string) {
	if !strings.ContainsAny(c, "/:") {
		return "", c, ""
	}
	last, depth := -1, 0
	for i := len(c) - 1; i >= 0; i-- {
		switch c[i] {
		case ']':
			depth++
		case '[':
			depth--
		case ':':
			if depth == 0 {
				last = i
			}
		}
		if last >= 0 {
			break
		}
	}
	name := c
	if last >= 0 {
		variant, name = c[:last], c[last+1:]
	}
	if s := strings.LastIndex(name, "/"); s >= 0 {
		return variant, name[:s], name[s+1:]
	}
	return variant, name, ""
}

func joinClass(variant, value, modifier string) string {
	if modifier != "" {
		value += "/" + modifier
	}
	if variant != "" {
		return variant + ":" + value
	}
	return value
}

// RTLClass is upstream's applyRtlMapping: the class string of the ui-rtl
// variant of a component, for the classes of the component's own source.
// (Upstream then sorts the classes again; the order does not change the
// rendering.)
func RTLClass(class string) string {
	var out []string
	for _, c := range strings.Split(class, " ") {
		if strings.HasPrefix(c, "rtl:") || strings.HasPrefix(c, "ltr:") {
			out = append(out, c)
			continue
		}
		if c == "cn-rtl-flip" {
			out = append(out, "rtl:rotate-180")
			continue
		}
		variant, value, modifier := splitClassName(c)
		if value == "" {
			out = append(out, c)
			continue
		}
		rtlVariant := "rtl:"
		if variant != "" {
			rtlVariant = "rtl:" + variant + ":"
		}
		mod := ""
		if modifier != "" {
			mod = "/" + modifier
		}
		if v, ok := strings.CutPrefix(value, "-translate-x-"); ok {
			out = append(out, c, rtlVariant+"translate-x-"+v+mod)
			continue
		}
		if v, ok := strings.CutPrefix(value, "translate-x-"); ok {
			out = append(out, c, rtlVariant+"-translate-x-"+v+mod)
			continue
		}
		if strings.HasPrefix(value, "space-x-") {
			out = append(out, c, rtlVariant+"space-x-reverse")
			continue
		}
		if strings.HasPrefix(value, "divide-x-") {
			out = append(out, c, rtlVariant+"divide-x-reverse")
			continue
		}
		if value == "cursor-w-resize" {
			out = append(out, c, rtlVariant+"cursor-e-resize")
			continue
		}
		if value == "cursor-e-resize" {
			out = append(out, c, rtlVariant+"cursor-w-resize")
			continue
		}
		slid := false
		for _, m := range rtlSlideMappings {
			if strings.Contains(variant, m[0]) && strings.HasPrefix(value, m[1]) {
				out = append(out, joinClass(variant, strings.Replace(value, m[1], m[2], 1), modifier))
				slid = true
				break
			}
		}
		if slid {
			continue
		}
		physicalSide := strings.Contains(variant, "data-[side=left]") || strings.Contains(variant, "data-[side=right]")
		mapped := value
		for _, m := range rtlMappings {
			if physicalSide && (strings.HasPrefix(m[0], "-left-") || strings.HasPrefix(m[0], "-right-") || strings.HasPrefix(m[0], "left-") || strings.HasPrefix(m[0], "right-")) {
				continue
			}
			if strings.HasPrefix(value, m[0]) {
				if !strings.HasSuffix(m[0], "-") && value != m[0] {
					continue
				}
				mapped = strings.Replace(value, m[0], m[1], 1)
				break
			}
		}
		out = append(out, joinClass(variant, mapped, modifier))
	}
	// prettier-plugin-tailwindcss removes duplicates.
	seen := make(map[string]bool, len(out))
	uniq := out[:0]
	for _, c := range out {
		if c == "" || !seen[c] {
			seen[c] = true
			uniq = append(uniq, c)
		}
	}
	return strings.Join(uniq, " ")
}
