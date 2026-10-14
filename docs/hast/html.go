package hast

import (
	"strconv"
	"strings"
	"unicode/utf8"
)

// ToHTML serializes nodes exactly like hast-util-to-html 9.0.5 with its
// default options: double quotes, hexadecimal character references for the
// characters that need escaping, no omitted tags, void elements without a
// closing slash.
//
// MDX nodes, which hast-util-to-html cannot serialize, are written as
// placeholder elements so trees can be compared with upstream's (the
// oracle in docs/oracle writes them the same way):
//
//	<mdx-jsx data-name="Name" data-flow="true" data-attrs="[...]">...</mdx-jsx>
//	<mdx-expr data-flow="true" data-value="..."></mdx-expr>
//	<mdx-esm data-value="..."></mdx-esm>
func ToHTML(nodes ...Node) string {
	var b strings.Builder
	for _, n := range nodes {
		writeNode(&b, n, "")
	}
	return b.String()
}

var voidElements = map[string]bool{
	"area": true, "base": true, "basefont": true, "bgsound": true, "br": true, "col": true,
	"command": true, "embed": true, "frame": true, "hr": true, "image": true, "img": true,
	"input": true, "keygen": true, "link": true, "meta": true, "param": true, "source": true,
	"track": true, "wbr": true,
}

func writeNode(b *strings.Builder, n Node, parentTag string) {
	switch x := n.(type) {
	case *Text:
		if parentTag == "script" || parentTag == "style" {
			b.WriteString(x.Value)
		} else {
			escape(b, x.Value, "<&")
		}
	case *Element:
		b.WriteByte('<')
		b.WriteString(x.Tag)
		first := true
		for _, p := range x.Props {
			if p.Unset {
				continue
			}
			if first {
				b.WriteByte(' ')
				first = false
			} else {
				b.WriteByte(' ')
			}
			b.WriteString(p.Name)
			if p.Bool {
				continue
			}
			b.WriteString(`="`)
			escape(b, p.Value, "\x00\"&'`")
			b.WriteByte('"')
		}
		b.WriteByte('>')
		if voidElements[x.Tag] && len(x.Children) == 0 {
			return
		}
		for _, c := range x.Children {
			writeNode(b, c, x.Tag)
		}
		b.WriteString("</")
		b.WriteString(x.Tag)
		b.WriteByte('>')
	case *MdxJsx:
		e := El("mdx-jsx", "data-name", x.Name, "data-flow", strconv.FormatBool(x.Flow), "data-attrs", AttrsJSON(x.Attrs))
		e.Children = x.Children
		writeNode(b, e, parentTag)
	case *MdxExpression:
		writeNode(b, El("mdx-expr", "data-flow", strconv.FormatBool(x.Flow), "data-value", x.Value), parentTag)
	case *MdxEsm:
		writeNode(b, El("mdx-esm", "data-value", x.Value), parentTag)
	}
}

// escape writes s with the characters in subset replaced by hexadecimal
// character references (stringify-entities' default format).
func escape(b *strings.Builder, s, subset string) {
	start := 0
	for i := 0; i < len(s); i++ {
		if strings.IndexByte(subset, s[i]) >= 0 {
			b.WriteString(s[start:i])
			b.WriteString("&#x")
			b.WriteString(strings.ToUpper(strconv.FormatInt(int64(s[i]), 16)))
			b.WriteByte(';')
			start = i + 1
		}
	}
	b.WriteString(s[start:])
}

// AttrsJSON encodes JSX attributes as the oracle does:
// [{"name":"a","value":"x"},{"name":"b"},{"name":"c","expr":"1 + 1"}].
func AttrsJSON(attrs []JsxAttr) string {
	var b strings.Builder
	b.WriteByte('[')
	for i, a := range attrs {
		if i > 0 {
			b.WriteByte(',')
		}
		b.WriteString(`{"name":`)
		JSONString(&b, a.Name)
		switch a.Kind {
		case AttrString:
			b.WriteString(`,"value":`)
			JSONString(&b, a.Value)
		case AttrExpression:
			b.WriteString(`,"expr":`)
			JSONString(&b, a.Value)
		}
		b.WriteByte('}')
	}
	b.WriteByte(']')
	return b.String()
}

// JSONString writes s as JSON.stringify does.
func JSONString(b *strings.Builder, s string) {
	b.WriteByte('"')
	for i := 0; i < len(s); {
		r, size := utf8.DecodeRuneInString(s[i:])
		switch {
		case r == '"':
			b.WriteString(`\"`)
		case r == '\\':
			b.WriteString(`\\`)
		case r == '\n':
			b.WriteString(`\n`)
		case r == '\r':
			b.WriteString(`\r`)
		case r == '\t':
			b.WriteString(`\t`)
		case r == '\b':
			b.WriteString(`\b`)
		case r == '\f':
			b.WriteString(`\f`)
		case r < 0x20:
			b.WriteString(`\u00`)
			b.WriteString(strconv.FormatInt(int64(r)>>4, 16))
			b.WriteString(strconv.FormatInt(int64(r)&15, 16))
		default:
			b.WriteString(s[i : i+size])
		}
		i += size
	}
	b.WriteByte('"')
}
