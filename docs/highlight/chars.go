package highlight

import (
	"strings"

	"shadcn.go/docs/hast"
)

// This file ports rehype-pretty-code 0.14.1's highlighted-characters
// support (meta strings like /word/ or "word"3-4#id): src/chars/*.ts. The
// port keeps the original control flow, including its JavaScript string
// semantics (String.prototype.replace replaces the first occurrence).

const visited = "rehype-pretty-code-visited"

type charsOptions struct {
	ranges  [][]int
	ids     map[string]string
	counter map[string]int
}

func jsReplace(s, old, new string) string { return strings.Replace(s, old, new, 1) }

func getContent(n hast.Node) (string, bool) {
	if n == nil {
		return "", false
	}
	return hast.ToString(n), true
}

func nodeAt(nodes []hast.Node, i int) hast.Node {
	if i < 0 || i >= len(nodes) {
		return nil
	}
	return nodes[i]
}

func asElement(n hast.Node) (*hast.Element, bool) {
	e, ok := n.(*hast.Element)
	return e, ok
}

func charsHighlighter(element *hast.Element, charsList []string, opts *charsOptions) {
	textContent := hast.ToString(element)
	for index, chars := range charsList {
		if chars == "" || !strings.Contains(textContent, chars) {
			continue
		}
		textContent2 := hast.ToString(element)
		startIndex := 0
		for strings.Contains(textContent2, chars) {
			var currentRange []int
			if index < len(opts.ranges) {
				currentRange = opts.ranges[index]
			}
			id := chars + "-" + itoa(index)
			opts.counter[id]++
			ignoreChars := len(currentRange) > 0 && !contains(currentRange, opts.counter[id])
			toWrap := getElementsToHighlight(element, chars, startIndex, ignoreChars)
			if len(toWrap) == 0 {
				break
			}
			wrapHighlightedChars(element, toWrap, opts, ignoreChars)
			startIndex = max(toWrap[len(toWrap)-1].index-2, 0)
			var b strings.Builder
			for _, c := range element.Children {
				if e, ok := c.(*hast.Element); ok {
					if e.HasKey(visited) || e.HasKey("data-highlighted-chars-mark") {
						continue
					}
				}
				b.WriteString(hast.ToString(c))
			}
			textContent2 = b.String()
		}
	}
	for _, c := range element.Children {
		if e, ok := c.(*hast.Element); ok && e.HasKey(visited) {
			e.Unset(visited)
		}
	}
}

type wrapItem struct {
	element hast.Node
	index   int
}

func getElementsToHighlight(element *hast.Element, chars string, startIndex int, ignoreChars bool) []wrapItem {
	var toWrap []wrapItem
	charsSoFar := ""
	for i := startIndex; i < len(element.Children); i++ {
		elements := element.Children
		remaining := chars
		if charsSoFar != "" {
			remaining = jsReplace(chars, charsSoFar, "")
		}
		if remaining == "" {
			return toWrap
		}
		maybe, ok := asElement(elements[i])
		if !ok || maybe.HasKey(visited) {
			continue
		}
		content, _ := getContent(maybe)
		if content == chars || charsSoFar+content == chars {
			toWrap = append(toWrap, wrapItem{maybe, i})
			return toWrap
		}
		if strings.HasPrefix(chars, charsSoFar+content) {
			if nextElementMaybeContinuesChars(elements, i+1, jsReplace(remaining, content, "")) {
				toWrap = append(toWrap, wrapItem{elements[i], i})
				charsSoFar += content
				continue
			}
		}
		overlap := findOverlap(content, remaining)
		partialMatch := overlap != "" && strings.HasPrefix(remaining, overlap)
		if !partialMatch {
			continue
		}
		nextPart := jsReplace(remaining, overlap, "")
		if nextContent, ok := getContent(nodeAt(elements, i+1)); nextPart != "" && ok && nextContent != "" &&
			!nextElementMaybeContinuesChars(elements, i+1, nextPart) {
			continue
		}
		splitParts := strings.Split(content, overlap)
		leftPart, rightPart := splitParts[0], ""
		if len(splitParts) > 1 {
			rightPart = splitParts[1]
		}
		var rest []string
		if len(splitParts) > 2 {
			rest = splitParts[2:]
		}
		if rightPart == "" && leftPart == "" && len(rest) == 0 {
			continue
		}
		withNext := content
		if nc, ok := getContent(nodeAt(elements, i+1)); ok {
			withNext += nc
		}
		nextNodeOverlap := findOverlap(withNext, remaining)
		splitIndex := strings.Index(withNext, nextNodeOverlap)
		if strings.HasSuffix(chars, overlap) || strings.HasPrefix(chars, overlap) {
			rightString := jsReplace(rightPart, overlap, "")
			leftString := jsSubstring(content, 0, splitIndex)
			nextContinues := nextElementMaybeContinuesChars(elements, i+1, nextPart)
			newEl, updated := splitElement(element, i, overlap, rightString, leftString, rest, nextContinues, ignoreChars)
			charsSoFar += overlap
			toWrap = append(toWrap, wrapItem{newEl, updated})
		}
	}
	return toWrap
}

// jsSubstring is String.prototype.substring on ASCII-safe byte offsets.
func jsSubstring(s string, a, b int) string {
	if a < 0 {
		a = 0
	}
	if b < 0 {
		b = 0
	}
	if a > b {
		a, b = b, a
	}
	if b > len(s) {
		b = len(s)
	}
	if a > len(s) {
		a = len(s)
	}
	return s[a:b]
}

func splitElement(parent *hast.Element, index int, inner, rightString, leftString string, rest []string, nextContinues, ignoreChars bool) (hast.Node, int) {
	toWrap := parent.Children[index]
	el, isEl := asElement(toWrap)
	if isEl && (len(el.Children) == 0 || !isTextNode(el.Children[0])) || ignoreChars {
		return toWrap, index
	}
	newIndex := index
	if t, ok := el.Children[0].(*hast.Text); ok {
		t.Value = inner
	}
	rightStr := rightString
	if len(rest) > 0 {
		for _, s := range rest {
			if s == "" {
				rightStr += inner
			} else {
				rightStr += inner + s
			}
		}
	}
	clone := func(text string) *hast.Element {
		c := &hast.Element{Tag: el.Tag, Props: append([]hast.Prop(nil), el.Props...), Data: el.Data}
		c.Children = []hast.Node{&hast.Text{Value: text}}
		return c
	}
	if leftString != "" {
		parent.Children = insert(parent.Children, newIndex, clone(leftString))
	}
	if rightStr != "" && !nextContinues {
		if leftString != "" {
			newIndex += 2
		} else {
			newIndex++
		}
		parent.Children = insert(parent.Children, newIndex, clone(rightStr))
	}
	return toWrap, index + 1
}

func isTextNode(n hast.Node) bool { _, ok := n.(*hast.Text); return ok }

func insert(list []hast.Node, i int, n hast.Node) []hast.Node {
	if i > len(list) {
		i = len(list)
	}
	list = append(list, nil)
	copy(list[i+1:], list[i:])
	list[i] = n
	return list
}

func nextElementMaybeContinuesChars(elements []hast.Node, nextIndex int, remainingPart string) bool {
	if remainingPart == "" {
		return false
	}
	content, ok := getContent(nodeAt(elements, nextIndex))
	if !ok || content == "" {
		return false
	}
	includesNext := strings.HasPrefix(content, remainingPart) || strings.HasPrefix(remainingPart, content)
	overlap := findOverlap(content, remainingPart)
	if overlap == remainingPart && strings.HasPrefix(content, remainingPart) {
		return true
	}
	if includesNext {
		return nextElementMaybeContinuesChars(elements, nextIndex+1, jsReplace(remainingPart, content, ""))
	}
	return false
}

func findOverlap(a, b string) string {
	for len(b) > 0 {
		if strings.HasSuffix(a, b) || strings.Contains(a, b) {
			return b
		}
		b = b[:len(b)-1]
	}
	return ""
}

func wrapHighlightedChars(parent *hast.Element, toWrap []wrapItem, opts *charsOptions, ignoreWord bool) {
	if len(toWrap) == 0 {
		return
	}
	first := toWrap[0].element
	if ignoreWord {
		if e, ok := asElement(first); ok {
			e.Set(visited, "")
		}
		return
	}
	if len(toWrap) > 1 {
		mark := hast.El("mark", "data-highlighted-chars-mark", "")
		for _, w := range toWrap {
			mark.Children = append(mark.Children, w.element)
		}
		start := toWrap[0].index
		end := start + len(toWrap)
		if end > len(parent.Children) {
			end = len(parent.Children)
		}
		kids := append([]hast.Node{}, parent.Children[:start]...)
		kids = append(kids, mark)
		kids = append(kids, parent.Children[end:]...)
		parent.Children = kids
		var word strings.Builder
		for _, n := range mark.Children {
			if e, ok := asElement(n); ok && len(e.Children) > 0 {
				if t, ok := e.Children[0].(*hast.Text); ok {
					word.WriteString(t.Value)
				}
			}
		}
		mark.Set("data-highlighted-chars", "")
		setID(mark, opts.ids, word.String())
		return
	}
	el, ok := asElement(first)
	if !ok || len(el.Children) == 0 {
		return
	}
	t, ok := el.Children[0].(*hast.Text)
	if !ok {
		return
	}
	el.Set(visited, "")
	el.Set("data-highlighted-chars", "")
	setID(el, opts.ids, t.Value)
	el.Tag = "mark"
	span := &hast.Element{Tag: "span"}
	if style, ok := el.Get("style"); ok {
		span.Set("style", style)
	} else {
		span.Props = []hast.Prop{{Name: "style", Unset: true}}
	}
	span.Children = el.Children
	el.Children = []hast.Node{span}
	if el.HasKey("style") {
		el.Unset("style")
	} else {
		el.Props = append(el.Props, hast.Prop{Name: "style", Unset: true})
	}
}

func setID(el *hast.Element, ids map[string]string, word string) {
	if id, ok := ids[word]; ok {
		el.Set("data-chars-id", id)
	} else if el.HasKey("data-chars-id") {
		el.Unset("data-chars-id")
	} else {
		el.Props = append(el.Props, hast.Prop{Name: "data-chars-id", Unset: true})
	}
}
