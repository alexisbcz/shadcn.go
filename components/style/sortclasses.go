package style

import (
	"sort"
	"strings"
)

// This file ports the class sorting of prettier-plugin-tailwindcss 0.7.2
// (sortClasses and sortClassList), which upstream's build runs on every
// generated file: it formats them with prettier, and the plugin is
// configured with tailwindFunctions ["cn", "cva"].

// collapse is the plugin's collapseWhitespace option (always enabled in
// upstream's configuration, which does not set tailwindPreserveWhitespace).
type collapse struct{ start, end bool }

var collapseAll = collapse{true, true}

type sortOpts struct {
	ignoreFirst, ignoreLast bool
	collapse                collapse
}

// isPluginSpace matches the plugin's [\t\r\f\n ] class separators.
func isPluginSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\r' || c == '\f' || c == '\n'
}

// sortClasses is the plugin's sortClasses with removeDuplicates enabled.
func sortClasses(s string, o sortOpts) string {
	if s == "" || strings.Contains(s, "{{") {
		return s
	}
	allSpace := true
	for i := 0; i < len(s); i++ {
		if !isPluginSpace(s[i]) {
			allSpace = false
			break
		}
	}
	if allSpace {
		return " "
	}

	// classStr.split(/([\t\r\f\n ]+)/)
	var classes, spaces []string
	start := 0
	for i := 0; i < len(s); {
		if !isPluginSpace(s[i]) {
			i++
			continue
		}
		j := i
		for j < len(s) && isPluginSpace(s[j]) {
			j++
		}
		classes = append(classes, s[start:i])
		spaces = append(spaces, s[i:j])
		start, i = j, j
	}
	classes = append(classes, s[start:])
	if classes[len(classes)-1] == "" {
		classes = classes[:len(classes)-1]
	}
	for i := range spaces {
		spaces[i] = " "
	}

	prefix := ""
	if o.ignoreFirst {
		prefix = shift(&classes) + shift(&spaces)
	}
	suffix := ""
	if o.ignoreLast {
		suffix = pop(&spaces) + pop(&classes)
	}

	list, removed := sortClassList(classes)
	kept := spaces[:0:0]
	for i, w := range spaces {
		if !removed[i+1] {
			kept = append(kept, w)
		}
	}
	var b strings.Builder
	for i, c := range list {
		b.WriteString(c)
		if i < len(kept) {
			b.WriteString(kept[i])
		}
	}
	result := b.String()

	// prefix.replace(/\s+$/g, " "), suffix.replace(/^\s+/g, " ")
	if t := trimRightJSSpace(prefix); t != prefix {
		prefix = t + " "
	}
	if t := trimLeftJSSpace(suffix); t != suffix {
		suffix = " " + t
	}
	if t := trimLeftJSSpace(result); t != result {
		if o.collapse.start {
			result = t
		} else {
			result = " " + t
		}
	}
	if t := trimRightJSSpace(result); t != result {
		if o.collapse.end {
			result = t
		} else {
			result = t + " "
		}
	}
	return prefix + result + suffix
}

func shift(s *[]string) string {
	if len(*s) == 0 {
		return ""
	}
	v := (*s)[0]
	*s = (*s)[1:]
	return v
}

func pop(s *[]string) string {
	if len(*s) == 0 {
		return ""
	}
	v := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return v
}

// sortClassList orders classes like the plugin: Tailwind's getClassOrder,
// unknown classes first in their original order ("..." last), then removes
// repeated known classes. removed holds the positions, in the sorted list,
// of the removed duplicates (the plugin drops the separators at those
// positions plus one).
func sortClassList(classes []string) ([]string, map[int]bool) {
	type entry struct {
		name string
		rank int32
	}
	es := make([]entry, len(classes))
	for i, c := range classes {
		es[i] = entry{c, classRank(c)}
	}
	isEllipsis := func(s string) bool { return s == "..." || s == "…" }
	cmp := func(a, z entry) int {
		switch {
		case isEllipsis(a.name):
			return 1
		case isEllipsis(z.name):
			return -1
		case a.rank == z.rank:
			return 0
		case a.rank < 0:
			return -1
		case z.rank < 0:
			return 1
		case a.rank < z.rank:
			return -1
		}
		return 1
	}
	sort.SliceStable(es, func(i, j int) bool { return cmp(es[i], es[j]) < 0 })

	removed := map[int]bool{}
	seen := map[string]bool{}
	out := make([]string, 0, len(es))
	for i, e := range es {
		if seen[e.name] {
			removed[i] = true
			continue
		}
		if e.rank >= 0 {
			seen[e.name] = true
		}
		out = append(out, e.name)
	}
	return out, removed
}
