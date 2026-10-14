// Package rtl ports applyRtlMapping of upstream's RTL transform
// (ui/packages/registry/src/transformers/transform-rtl.ts), which builds the
// styles/<base>-<style>/ui-rtl components the docs' *-rtl previews import:
// physical Tailwind classes become logical ones (ml- → ms-, right- → end-,
// ...), translate-x gets an rtl: counterpart and the cn-rtl-flip marker
// becomes rtl:rotate-180. Logical classes render the same in LTR, so a
// component may apply Classes whenever it renders right-to-left.
package rtl

import "strings"

var mappings = [][2]string{
	{"-ml-", "-ms-"}, {"-mr-", "-me-"}, {"ml-", "ms-"}, {"mr-", "me-"},
	{"pl-", "ps-"}, {"pr-", "pe-"},
	{"-left-", "-start-"}, {"-right-", "-end-"}, {"left-", "start-"}, {"right-", "end-"},
	{"inset-l-", "inset-inline-start-"}, {"inset-r-", "inset-inline-end-"},
	{"rounded-tl-", "rounded-ss-"}, {"rounded-tr-", "rounded-se-"}, {"rounded-bl-", "rounded-es-"}, {"rounded-br-", "rounded-ee-"},
	{"rounded-l-", "rounded-s-"}, {"rounded-r-", "rounded-e-"},
	{"border-l-", "border-s-"}, {"border-r-", "border-e-"}, {"border-l", "border-s"}, {"border-r", "border-e"},
	{"text-left", "text-start"}, {"text-right", "text-end"},
	{"scroll-ml-", "scroll-ms-"}, {"scroll-mr-", "scroll-me-"}, {"scroll-pl-", "scroll-ps-"}, {"scroll-pr-", "scroll-pe-"},
	{"float-left", "float-start"}, {"float-right", "float-end"},
	{"clear-left", "clear-start"}, {"clear-right", "clear-end"},
	{"origin-top-left", "origin-top-start"}, {"origin-top-right", "origin-top-end"},
	{"origin-bottom-left", "origin-bottom-start"}, {"origin-bottom-right", "origin-bottom-end"},
	{"origin-left", "origin-start"}, {"origin-right", "origin-end"},
}

var translateX = [][2]string{{"-translate-x-", "translate-x-"}, {"translate-x-", "-translate-x-"}}

var reverse = [][2]string{{"space-x-", "space-x-reverse"}, {"divide-x-", "divide-x-reverse"}}

var swap = [][2]string{{"cursor-w-resize", "cursor-e-resize"}, {"cursor-e-resize", "cursor-w-resize"}}

var slides = [][3]string{
	{"data-[side=inline-start]", "slide-in-from-right", "slide-in-from-end"},
	{"data-[side=inline-start]", "slide-out-to-right", "slide-out-to-end"},
	{"data-[side=inline-end]", "slide-in-from-left", "slide-in-from-start"},
	{"data-[side=inline-end]", "slide-out-to-left", "slide-out-to-start"},
}

var positioning = []string{"-left-", "-right-", "left-", "right-"}

// split is splitClassName of transform-css-vars.ts: variant, value and
// alpha modifier ("" when absent).
func split(c string) (variant, value, modifier string) {
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
			if depth == 0 && last < 0 {
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
	if i := strings.LastIndexByte(name, '/'); i >= 0 {
		return variant, name[:i], name[i+1:]
	}
	return variant, name, ""
}

func join(variant, value, modifier string) string {
	s := value
	if modifier != "" {
		s += "/" + modifier
	}
	if variant != "" {
		s = variant + ":" + s
	}
	return s
}

// Classes maps a class string as upstream's RTL transform maps a class
// literal.
func Classes(s string) string {
	var out []string
	for _, c := range strings.Split(s, " ") {
		out = append(out, class(c)...)
	}
	return strings.Join(out, " ")
}

func class(c string) []string {
	if strings.HasPrefix(c, "rtl:") || strings.HasPrefix(c, "ltr:") {
		return []string{c}
	}
	if c == "cn-rtl-flip" {
		return []string{"rtl:rotate-180"}
	}
	variant, value, modifier := split(c)
	if value == "" {
		return []string{c}
	}
	rtlVariant := "rtl:"
	if variant != "" {
		rtlVariant = "rtl:" + variant + ":"
	}
	mod := ""
	if modifier != "" {
		mod = "/" + modifier
	}
	for _, m := range translateX {
		if strings.HasPrefix(value, m[0]) {
			return []string{c, rtlVariant + strings.Replace(value, m[0], m[1], 1) + mod}
		}
	}
	for _, m := range reverse {
		if strings.HasPrefix(value, m[0]) {
			return []string{c, rtlVariant + m[1]}
		}
	}
	for _, m := range swap {
		if value == m[0] {
			return []string{c, rtlVariant + m[1]}
		}
	}
	for _, m := range slides {
		if strings.Contains(variant, m[0]) && strings.HasPrefix(value, m[1]) {
			return []string{join(variant, strings.Replace(value, m[1], m[2], 1), modifier)}
		}
	}
	physicalSide := strings.Contains(variant, "data-[side=left]") || strings.Contains(variant, "data-[side=right]")
	mapped := value
	for _, m := range mappings {
		if physicalSide && hasAnyPrefix(m[0], positioning) {
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
	return []string{join(variant, mapped, modifier)}
}

func hasAnyPrefix(s string, prefixes []string) bool {
	for _, p := range prefixes {
		if strings.HasPrefix(s, p) {
			return true
		}
	}
	return false
}
