package highlight

import (
	"strings"
	"unicode/utf16"
	"unicode/utf8"

	"github.com/dlclark/regexp2/v2"
)

// goja does not implement the RegExp d flag (hasIndices), which shiki's
// JavaScript engine needs: it reads match.indices to place every capture.
// The prelude routes exec of regexes created with d to execD, which matches
// with regexp2 in ECMAScript mode exactly as goja itself does for patterns
// it cannot hand to Go's regexp (same options, same pattern conversion,
// same UTF-16 or code-point matching), and reports the group positions.

type dRegexp struct {
	re      *regexp2.Regexp
	unicode bool
	sticky  bool
	names   []string // group names by number ("" when unnamed)
}

type dRegistry struct{ list []*dRegexp }

// compile returns the id of the compiled pattern.
func (r *dRegistry) compile(source, flags string) (int, error) {
	opts := regexp2.ECMAScript
	d := &dRegexp{}
	for _, f := range flags {
		switch f {
		case 'm':
			opts |= regexp2.Multiline
		case 's':
			opts |= regexp2.Singleline
		case 'i':
			opts |= regexp2.IgnoreCase
		case 'u':
			opts |= regexp2.Unicode
			d.unicode = true
		case 'y':
			d.sticky = true
		}
	}
	if d.unicode {
		source = surrogateEscapesToRunes(source)
	} else {
		source = astralToSurrogateEscapes(source)
	}
	re, err := regexp2.Compile(source, opts)
	if err != nil {
		return 0, err
	}
	d.re = re
	nums := re.GetGroupNumbers()
	d.names = make([]string, len(nums))
	for i, n := range nums {
		if name := re.GroupNameFromNumber(n); name != "" && name != itoa(n) {
			d.names[i] = name
		}
	}
	r.list = append(r.list, d)
	return len(r.list) - 1, nil
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var b [20]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	return string(b[i:])
}

// exec matches from the UTF-16 index start and returns the start and end
// of every group in UTF-16 units (-1, -1 for groups that did not
// participate), or nil when there is no match.
func (r *dRegistry) exec(id int, s string, start int) ([]int, error) {
	d := r.list[id]
	var runes []rune
	var pos []int // rune index -> UTF-16 index (len(runes)+1 entries)
	if d.unicode {
		runes = make([]rune, 0, len(s))
		pos = make([]int, 0, len(s)+1)
		u := 0
		for _, c := range s {
			runes = append(runes, c)
			pos = append(pos, u)
			if c >= 0x10000 {
				u += 2
			} else {
				u++
			}
		}
		pos = append(pos, u)
	} else {
		for _, u := range utf16.Encode([]rune(s)) {
			runes = append(runes, rune(u))
		}
	}
	rstart := start
	if d.unicode {
		rstart = len(runes)
		for i, p := range pos {
			if p >= start {
				rstart = i
				break
			}
		}
	}
	m, err := d.re.FindRunesMatchStartingAt(runes, rstart)
	if err != nil || m == nil {
		return nil, err
	}
	if d.sticky && m.RuneIndex != rstart {
		return nil, nil
	}
	groups := m.Groups()
	out := make([]int, 0, 2*len(groups))
	for _, g := range groups {
		if len(g.Captures) == 0 {
			out = append(out, -1, -1)
			continue
		}
		a, b := g.RuneIndex, g.RuneIndex+g.RuneLength
		if d.unicode {
			a, b = pos[a], pos[b]
		}
		out = append(out, a, b)
	}
	return out, nil
}

// surrogateEscapesToRunes replaces 😀-style escaped surrogate
// pairs with the literal code point, as goja does for u-flag patterns.
func surrogateEscapesToRunes(p string) string {
	var b strings.Builder
	pos := 0
	for i := 0; i < len(p)-11; {
		if p[i] != '\\' {
			_, size := utf8.DecodeRuneInString(p[i:])
			i += size
			continue
		}
		i++
		if p[i] == 'u' && p[i+5] == '\\' && p[i+6] == 'u' {
			first, ok1 := hex4(p[i+1 : i+5])
			second, ok2 := hex4(p[i+7 : i+11])
			if ok1 && ok2 && utf16.IsSurrogate(rune(first)) && first < 0xdc00 && second >= 0xdc00 && second < 0xe000 {
				b.WriteString(p[pos : i-1])
				b.WriteRune(utf16.DecodeRune(rune(first), rune(second)))
				i += 11
				pos = i
				continue
			}
		}
		i++
	}
	if pos == 0 {
		return p
	}
	b.WriteString(p[pos:])
	return b.String()
}

// astralToSurrogateEscapes writes literal astral characters as escaped
// surrogate pairs, as goja does for patterns without the u flag (they then
// match UTF-16 code units).
func astralToSurrogateEscapes(p string) string {
	if !strings.ContainsFunc(p, func(r rune) bool { return r >= 0x10000 }) {
		return p
	}
	var b strings.Builder
	for _, r := range p {
		if r >= 0x10000 {
			a, c := utf16.EncodeRune(r)
			b.WriteString(`\u`)
			b.WriteString(hexU(a))
			b.WriteString(`\u`)
			b.WriteString(hexU(c))
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}

func hexU(r rune) string {
	const digits = "0123456789ABCDEF"
	return string([]byte{digits[r>>12&15], digits[r>>8&15], digits[r>>4&15], digits[r&15]})
}

func hex4(s string) (int, bool) {
	v := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c >= '0' && c <= '9':
			v = v*16 + int(c-'0')
		case c >= 'a' && c <= 'f':
			v = v*16 + int(c-'a'+10)
		case c >= 'A' && c <= 'F':
			v = v*16 + int(c-'A'+10)
		default:
			return 0, false
		}
	}
	return v, true
}
