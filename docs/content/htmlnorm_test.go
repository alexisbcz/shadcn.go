package content

import (
	"fmt"
	"html"
	"sort"
	"strings"
)

// A small HTML parser for the well-formed markup React and libhtml write,
// used to compare renderings as DOM trees: attribute order, character
// reference spelling, comments (React's <!-- --> text separators) and
// adjacent text node boundaries do not matter; everything else does.

type dnode struct {
	tag      string // "" for text
	text     string
	attrs    map[string]string
	children []*dnode
}

var voidTags = map[string]bool{
	"area": true, "base": true, "br": true, "col": true, "embed": true, "hr": true, "img": true,
	"input": true, "link": true, "meta": true, "source": true, "track": true, "wbr": true,
}

var rawTextTags = map[string]bool{"script": true, "style": true}

// parseHTML parses a fragment. It fails on mismatched tags.
func parseHTML(s string) ([]*dnode, error) {
	root := &dnode{tag: "#root"}
	stack := []*dnode{root}
	top := func() *dnode { return stack[len(stack)-1] }
	addText := func(t string) {
		if t == "" {
			return
		}
		p := top()
		if n := len(p.children); n > 0 && p.children[n-1].tag == "" {
			p.children[n-1].text += t
			return
		}
		p.children = append(p.children, &dnode{text: t})
	}
	i := 0
	for i < len(s) {
		lt := strings.IndexByte(s[i:], '<')
		if lt < 0 {
			addText(html.UnescapeString(s[i:]))
			break
		}
		addText(html.UnescapeString(s[i : i+lt]))
		i += lt
		switch {
		case strings.HasPrefix(s[i:], "<!--"):
			end := strings.Index(s[i+4:], "-->")
			if end < 0 {
				return nil, fmt.Errorf("unterminated comment")
			}
			i += 4 + end + 3
		case strings.HasPrefix(s[i:], "<!"):
			end := strings.IndexByte(s[i:], '>')
			i += end + 1
		case strings.HasPrefix(s[i:], "</"):
			end := strings.IndexByte(s[i:], '>')
			if end < 0 {
				return nil, fmt.Errorf("unterminated end tag")
			}
			name := strings.ToLower(strings.TrimSpace(s[i+2 : i+end]))
			i += end + 1
			if top().tag != name {
				return nil, fmt.Errorf("mismatched </%s>, open <%s>", name, top().tag)
			}
			stack = stack[:len(stack)-1]
		default:
			n, next, selfClosing, err := parseStartTag(s, i)
			if err != nil {
				return nil, err
			}
			i = next
			p := top()
			p.children = append(p.children, n)
			if rawTextTags[n.tag] {
				end := strings.Index(s[i:], "</"+n.tag)
				if end < 0 {
					return nil, fmt.Errorf("unterminated <%s>", n.tag)
				}
				if end > 0 {
					n.children = []*dnode{{text: s[i : i+end]}}
				}
				i += end
				close := strings.IndexByte(s[i:], '>')
				i += close + 1
				continue
			}
			if !voidTags[n.tag] && !selfClosing {
				stack = append(stack, n)
			}
		}
	}
	if len(stack) != 1 {
		return nil, fmt.Errorf("unclosed <%s>", top().tag)
	}
	return root.children, nil
}

func parseStartTag(s string, i int) (*dnode, int, bool, error) {
	j := i + 1
	for j < len(s) && !isSpaceByte(s[j]) && s[j] != '>' && s[j] != '/' {
		j++
	}
	n := &dnode{tag: strings.ToLower(s[i+1 : j]), attrs: map[string]string{}}
	for {
		for j < len(s) && isSpaceByte(s[j]) {
			j++
		}
		if j >= len(s) {
			return nil, 0, false, fmt.Errorf("unterminated tag <%s", n.tag)
		}
		if s[j] == '>' {
			return n, j + 1, false, nil
		}
		if strings.HasPrefix(s[j:], "/>") {
			return n, j + 2, true, nil
		}
		k := j
		for k < len(s) && !isSpaceByte(s[k]) && s[k] != '=' && s[k] != '>' && !strings.HasPrefix(s[k:], "/>") {
			k++
		}
		name := strings.ToLower(s[j:k])
		j = k
		for j < len(s) && isSpaceByte(s[j]) {
			j++
		}
		if j < len(s) && s[j] == '=' {
			j++
			for j < len(s) && isSpaceByte(s[j]) {
				j++
			}
			var v string
			if s[j] == '"' || s[j] == '\'' {
				q := s[j]
				end := strings.IndexByte(s[j+1:], q)
				if end < 0 {
					return nil, 0, false, fmt.Errorf("unterminated attribute")
				}
				v = s[j+1 : j+1+end]
				j += end + 2
			} else {
				k := j
				for k < len(s) && !isSpaceByte(s[k]) && s[k] != '>' {
					k++
				}
				v = s[j:k]
				j = k
			}
			v = html.UnescapeString(v)
			if name == "class" {
				// A class list is a set of tokens: React keeps the trailing
				// space @tabler/icons-react leaves in it, libhtml does not.
				v = strings.Join(strings.Fields(v), " ")
			}
			n.attrs[name] = v
		} else {
			n.attrs[name] = ""
		}
	}
}

func isSpaceByte(c byte) bool { return c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == '\f' }

// canon serializes nodes canonically (attributes sorted).
func canon(nodes []*dnode) string {
	var b strings.Builder
	for _, n := range nodes {
		writeCanon(&b, n)
	}
	return b.String()
}

func writeCanon(b *strings.Builder, n *dnode) {
	if n.tag == "" {
		b.WriteString(html.EscapeString(n.text))
		return
	}
	b.WriteByte('<')
	b.WriteString(n.tag)
	keys := make([]string, 0, len(n.attrs))
	for k := range n.attrs {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Fprintf(b, " %s=%q", k, n.attrs[k])
	}
	b.WriteByte('>')
	for _, c := range n.children {
		writeCanon(b, c)
	}
	b.WriteString("</")
	b.WriteString(n.tag)
	b.WriteByte('>')
}

// firstDiff locates the first difference between two canonical strings.
func firstDiff(a, b string) string {
	i := 0
	for i < len(a) && i < len(b) && a[i] == b[i] {
		i++
	}
	lo := max(i-160, 0)
	return fmt.Sprintf("\n got: …%s\nwant: …%s", a[lo:min(i+160, len(a))], b[lo:min(i+160, len(b))])
}
