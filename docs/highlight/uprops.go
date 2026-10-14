package highlight

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"unicode"
)

// goja's regular expressions do not implement ECMAScript Unicode property
// escapes (\p{...} with the u flag) beyond a few general categories, and
// regexp2 rejects binary properties such as \p{Alpha} or \p{AHex}. The
// shiki JavaScript engine (oniguruma-to-es) emits them for POSIX classes
// and Oniguruma properties. The prelude rewrites each escape into an
// explicit character class built here from Go's Unicode tables, which
// matches exactly the same code points.

var (
	propMu    sync.Mutex
	propCache = map[string]string{}
)

// propClass returns the body of a character class (without brackets)
// matching the code points of the property, or its complement.
func propClass(name string, negate bool) (string, error) {
	key := name
	if negate {
		key = "^" + name
	}
	propMu.Lock()
	defer propMu.Unlock()
	if s, ok := propCache[key]; ok {
		return s, nil
	}
	set, err := propSet(name)
	if err != nil {
		return "", err
	}
	if negate {
		set = complement(set)
	}
	var b strings.Builder
	for _, r := range set {
		writeClassRune(&b, r[0])
		if r[1] != r[0] {
			if r[1] > r[0]+1 {
				b.WriteByte('-')
			}
			writeClassRune(&b, r[1])
		}
	}
	s := b.String()
	propCache[key] = s
	return s, nil
}

func writeClassRune(b *strings.Builder, r rune) {
	if r < 0x80 && (r >= '0' && r <= '9' || r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z') {
		b.WriteRune(r)
		return
	}
	fmt.Fprintf(b, `\u{%x}`, r)
}

type runeRange [2]rune

func fromTables(tables ...*unicode.RangeTable) []runeRange {
	var out []runeRange
	for _, t := range tables {
		for _, r := range t.R16 {
			for c := rune(r.Lo); c <= rune(r.Hi); c += rune(r.Stride) {
				out = append(out, runeRange{c, c})
				if r.Stride == 1 {
					out[len(out)-1][1] = rune(r.Hi)
					break
				}
			}
		}
		for _, r := range t.R32 {
			for c := rune(r.Lo); c <= rune(r.Hi); c += rune(r.Stride) {
				out = append(out, runeRange{c, c})
				if r.Stride == 1 {
					out[len(out)-1][1] = rune(r.Hi)
					break
				}
			}
		}
	}
	return normalize(out)
}

func normalize(rs []runeRange) []runeRange {
	sort.Slice(rs, func(i, j int) bool { return rs[i][0] < rs[j][0] })
	var out []runeRange
	for _, r := range rs {
		if n := len(out); n > 0 && r[0] <= out[n-1][1]+1 {
			if r[1] > out[n-1][1] {
				out[n-1][1] = r[1]
			}
			continue
		}
		out = append(out, r)
	}
	return out
}

func union(sets ...[]runeRange) []runeRange {
	var all []runeRange
	for _, s := range sets {
		all = append(all, s...)
	}
	return normalize(all)
}

func complement(s []runeRange) []runeRange {
	var out []runeRange
	next := rune(0)
	for _, r := range s {
		if r[0] > next {
			out = append(out, runeRange{next, r[0] - 1})
		}
		next = r[1] + 1
	}
	if next <= unicode.MaxRune {
		out = append(out, runeRange{next, unicode.MaxRune})
	}
	return out
}

func subtract(a, b []runeRange) []runeRange {
	// a ∩ ¬b
	nb := complement(b)
	var out []runeRange
	i, j := 0, 0
	for i < len(a) && j < len(nb) {
		lo := max(a[i][0], nb[j][0])
		hi := min(a[i][1], nb[j][1])
		if lo <= hi {
			out = append(out, runeRange{lo, hi})
		}
		if a[i][1] < nb[j][1] {
			i++
		} else {
			j++
		}
	}
	return out
}

var generalCategories = map[string]string{
	"Letter": "L", "Cased_Letter": "LC", "Uppercase_Letter": "Lu", "Lowercase_Letter": "Ll",
	"Titlecase_Letter": "Lt", "Modifier_Letter": "Lm", "Other_Letter": "Lo", "Mark": "M",
	"Combining_Mark": "M", "Nonspacing_Mark": "Mn", "Spacing_Mark": "Mc", "Enclosing_Mark": "Me",
	"Number": "N", "Decimal_Number": "Nd", "digit": "Nd", "Letter_Number": "Nl", "Other_Number": "No",
	"Punctuation": "P", "punct": "P", "Connector_Punctuation": "Pc", "Dash_Punctuation": "Pd",
	"Open_Punctuation": "Ps", "Close_Punctuation": "Pe", "Initial_Punctuation": "Pi",
	"Final_Punctuation": "Pf", "Other_Punctuation": "Po", "Symbol": "S", "Math_Symbol": "Sm",
	"Currency_Symbol": "Sc", "Modifier_Symbol": "Sk", "Other_Symbol": "So", "Separator": "Z",
	"Space_Separator": "Zs", "Line_Separator": "Zl", "Paragraph_Separator": "Zp", "Other": "C",
	"Control": "Cc", "cntrl": "Cc", "Format": "Cf", "Surrogate": "Cs", "Private_Use": "Co", "Unassigned": "Cn",
}

func propSet(name string) ([]runeRange, error) {
	if i := strings.IndexByte(name, '='); i >= 0 {
		key, val := name[:i], name[i+1:]
		switch key {
		case "General_Category", "gc":
			return propSet(val)
		case "Script", "sc", "Script_Extensions", "scx":
			if t, ok := unicode.Scripts[val]; ok {
				return fromTables(t), nil
			}
		}
		return nil, fmt.Errorf("unsupported Unicode property %q", name)
	}
	if short, ok := generalCategories[name]; ok {
		name = short
	}
	switch name {
	case "Any":
		return []runeRange{{0, unicode.MaxRune}}, nil
	case "ASCII":
		return []runeRange{{0, 0x7f}}, nil
	case "AHex", "ASCII_Hex_Digit":
		return []runeRange{{'0', '9'}, {'A', 'F'}, {'a', 'f'}}, nil
	case "LC":
		return fromTables(unicode.Lu, unicode.Ll, unicode.Lt), nil
	case "Cn":
		return complement(fromTables(allAssigned()...)), nil
	case "C":
		return union(fromTables(unicode.Cc, unicode.Cf, unicode.Cs, unicode.Co), complement(fromTables(allAssigned()...))), nil
	case "Alpha", "Alphabetic":
		return fromTables(unicode.Lu, unicode.Ll, unicode.Lt, unicode.Lm, unicode.Lo, unicode.Nl, unicode.Other_Alphabetic), nil
	case "Upper", "Uppercase":
		return fromTables(unicode.Lu, unicode.Other_Uppercase), nil
	case "Lower", "Lowercase":
		return fromTables(unicode.Ll, unicode.Other_Lowercase), nil
	case "Cased":
		return fromTables(unicode.Lu, unicode.Ll, unicode.Lt, unicode.Other_Uppercase, unicode.Other_Lowercase), nil
	case "White_Space", "space":
		return fromTables(unicode.White_Space), nil
	case "IDS", "ID_Start":
		return subtract(fromTables(unicode.L, unicode.Nl, unicode.Other_ID_Start), fromTables(unicode.Pattern_Syntax, unicode.Pattern_White_Space)), nil
	case "IDC", "ID_Continue":
		return subtract(fromTables(unicode.L, unicode.Nl, unicode.Other_ID_Start, unicode.Mn, unicode.Mc, unicode.Nd, unicode.Pc, unicode.Other_ID_Continue), fromTables(unicode.Pattern_Syntax, unicode.Pattern_White_Space)), nil
	}
	if t, ok := unicode.Categories[name]; ok {
		return fromTables(t), nil
	}
	if t, ok := unicode.Properties[name]; ok {
		return fromTables(t), nil
	}
	if t, ok := unicode.Scripts[name]; ok {
		return fromTables(t), nil
	}
	return nil, fmt.Errorf("unsupported Unicode property %q", name)
}

func allAssigned() []*unicode.RangeTable {
	var ts []*unicode.RangeTable
	for _, k := range []string{"L", "M", "N", "P", "S", "Z", "Cc", "Cf", "Cs", "Co"} {
		ts = append(ts, unicode.Categories[k])
	}
	return ts
}
