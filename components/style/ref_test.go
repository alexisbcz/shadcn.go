package style

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// refDir is the data of testdata/ref/<dir>.json.gz, written by
// testdata/refgen/extract.mjs from upstream's sources.
type refDir struct {
	Dir   string    `json:"dir"`
	Files []refFile `json:"files"`
}

type refFile struct {
	Path string    `json:"path"`
	Base []refRoot `json:"base"`
	Gen  []refRoot `json:"gen"`
}

type refRoot struct {
	Line int            `json:"line"`
	N    map[string]any `json:"n"`
}

func loadRef(name string) (*refDir, error) {
	f, err := os.Open(filepath.Join("testdata/ref", name+".json.gz"))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	zr, err := gzip.NewReader(f)
	if err != nil {
		return nil, err
	}
	var d refDir
	if err := json.NewDecoder(zr).Decode(&d); err != nil {
		return nil, err
	}
	return &d, nil
}

func decodeRoots(rs []refRoot) []*node {
	out := make([]*node, len(rs))
	for i, r := range rs {
		out[i] = decodeNode(r.N)
	}
	return out
}

func decodeNode(m map[string]any) *node {
	if m == nil {
		return nil
	}
	str := func(k string) string { s, _ := m[k].(string); return s }
	sub := func(k string) *node {
		v, _ := m[k].(map[string]any)
		return decodeNode(v)
	}
	list := func(k string) []*node {
		l, _ := m[k].([]any)
		out := make([]*node, 0, len(l))
		for _, x := range l {
			if mm, ok := x.(map[string]any); ok {
				out = append(out, decodeNode(mm))
			}
		}
		return out
	}
	switch m["k"] {
	case "str":
		return &node{kind: kStr, str: str("v"), tpl: m["t"] != nil}
	case "tpl":
		n := &node{kind: kTpl, exprs: list("e")}
		for _, q := range m["q"].([]any) {
			p := q.([]any)
			n.quasis = append(n.quasis, [2]string{p[0].(string), p[1].(string)})
		}
		return n
	case "call":
		return &node{kind: kCall, callee: str("f"), calleeRoot: str("r"), calleeNode: sub("c"), args: list("a")}
	case "obj":
		return &node{kind: kObject, children: list("p")}
	case "pa":
		return &node{kind: kProp, name: str("n"), ident: m["id"] == float64(1), key: sub("nk"), value: sub("v")}
	case "attr":
		return &node{kind: kAttr, name: str("n"), value: sub("v")}
	case "jx":
		return &node{kind: kJsxExpr, value: sub("e")}
	case "jsx":
		attrs := list("a")
		n := &node{kind: kJSX, children: append(attrs, list("c")...), nAttrs: len(attrs)}
		return n
	case "bin":
		return &node{kind: kBinary, op: str("op"), left: sub("l"), right: sub("r")}
	case "id":
		return &node{kind: kID, str: str("v")}
	case "g":
		return &node{kind: kGeneric, str: str("s"), children: list("c")}
	}
	panic(fmt.Sprintf("unknown node %v", m["k"]))
}

// compareTrees compares a transformed base tree with the generated one. It
// counts the strings compared and reports the first difference.
func compareTrees(got, want *node, path string, strings *int) string {
	if got == nil || want == nil {
		if got != want {
			return fmt.Sprintf("%s: got %s, want %s", path, describe(got), describe(want))
		}
		return ""
	}
	if got.kind != want.kind {
		return fmt.Sprintf("%s: got %s, want %s", path, describe(got), describe(want))
	}
	switch got.kind {
	case kStr:
		*strings++
		if got.str != want.str {
			return fmt.Sprintf("%s: got %q, want %q", path, got.str, want.str)
		}
		return ""
	case kTpl:
		if len(got.quasis) != len(want.quasis) {
			return fmt.Sprintf("%s: template parts differ", path)
		}
		for i := range got.quasis {
			*strings++
			if got.quasis[i] != want.quasis[i] {
				return fmt.Sprintf("%s.quasi[%d]: got %q, want %q", path, i, got.quasis[i], want.quasis[i])
			}
		}
		return compareLists(got.exprs, want.exprs, path+".expr", strings)
	case kCall:
		if got.callee != want.callee || got.calleeRoot != want.calleeRoot {
			return fmt.Sprintf("%s: callee %q, want %q", path, got.callee, want.callee)
		}
		if d := compareTrees(got.calleeNode, want.calleeNode, path+".callee", strings); d != "" {
			return d
		}
		return compareLists(got.args, want.args, path+"."+got.callee+"()", strings)
	case kProp:
		// A string key is compared as a string (prettier sorts it), others by name.
		if (got.key == nil && got.name != want.name) || got.ident != want.ident || (got.key == nil) != (want.key == nil) {
			return fmt.Sprintf("%s: property %q, want %q", path, got.name, want.name)
		}
		if d := compareTrees(got.key, want.key, path+".key", strings); d != "" {
			return d
		}
		return compareTrees(got.value, want.value, path+"."+got.name, strings)
	case kAttr:
		if got.name != want.name {
			return fmt.Sprintf("%s: attribute %q, want %q", path, got.name, want.name)
		}
		return compareTrees(got.value, want.value, path+"@"+got.name, strings)
	case kJsxExpr:
		return compareTrees(got.value, want.value, path+"{}", strings)
	case kBinary:
		if got.op != want.op {
			return fmt.Sprintf("%s: operator %q, want %q", path, got.op, want.op)
		}
		if d := compareTrees(got.left, want.left, path+".left", strings); d != "" {
			return d
		}
		return compareTrees(got.right, want.right, path+".right", strings)
	case kID:
		if got.str != want.str {
			return fmt.Sprintf("%s: identifier %q, want %q", path, got.str, want.str)
		}
		return ""
	case kObject, kJSX:
		return compareLists(got.children, want.children, path, strings)
	case kGeneric:
		if got.str != want.str {
			return fmt.Sprintf("%s: %s, want %s", path, got.str, want.str)
		}
		return compareLists(got.children, want.children, path+"."+got.str, strings)
	}
	return ""
}

func compareLists(got, want []*node, path string, strings *int) string {
	if len(got) != len(want) {
		var g, w []string
		for _, n := range got {
			g = append(g, describe(n))
		}
		for _, n := range want {
			w = append(w, describe(n))
		}
		return fmt.Sprintf("%s: got %d nodes [%s], want %d [%s]", path, len(got), strings2(g), len(want), strings2(w))
	}
	for i := range got {
		if d := compareTrees(got[i], want[i], fmt.Sprintf("%s[%d]", path, i), strings); d != "" {
			return d
		}
	}
	return ""
}

func strings2(s []string) string { return strings.Join(s, ", ") }

func describe(n *node) string {
	if n == nil {
		return "nothing"
	}
	switch n.kind {
	case kStr:
		return fmt.Sprintf("%q", n.str)
	case kID:
		return n.str
	case kCall:
		return n.callee + "()"
	case kGeneric:
		return n.str
	}
	return fmt.Sprintf("node(%d)", n.kind)
}
