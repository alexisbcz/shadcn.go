package highlight

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// rewritePattern replaces the Unicode property escapes of a u-flag pattern
// with explicit character classes (see uprops.go), and rejects modifier
// groups "(?i:...)", which ECMAScript 2018 does not have but goja's
// regexp2 accepts: oniguruma-to-es probes for them to choose its output
// target, and must see what an ES2018 engine reports.
func rewritePattern(src, flags string) (string, error) {
	unicode := strings.ContainsRune(flags, 'u')
	var b strings.Builder
	inClass := false
	for i := 0; i < len(src); i++ {
		c := src[i]
		if c == '\\' {
			if i+1 >= len(src) {
				b.WriteByte(c)
				continue
			}
			n := src[i+1]
			if unicode && (n == 'p' || n == 'P') && i+2 < len(src) && src[i+2] == '{' {
				end := strings.IndexByte(src[i+3:], '}')
				if end < 0 {
					return "", fmt.Errorf("SyntaxError: Invalid regular expression: /%s/: Invalid property name", src)
				}
				body, err := propClass(src[i+3:i+3+end], n == 'P')
				if err != nil {
					return "", err
				}
				if inClass {
					b.WriteString(body)
				} else {
					b.WriteString("[" + body + "]")
				}
				i += 3 + end
				continue
			}
			b.WriteByte(c)
			b.WriteByte(n)
			i++
			continue
		}
		if inClass {
			if c == ']' {
				inClass = false
			}
			b.WriteByte(c)
			continue
		}
		switch {
		case c == '[':
			inClass = true
			b.WriteByte(c)
			if i+1 < len(src) && src[i+1] == '^' {
				b.WriteByte('^')
				i++
			}
			if i+1 < len(src) && src[i+1] == ']' {
				b.WriteByte(']')
				i++
				inClass = false
			}
			continue
		case c == '(' && strings.HasPrefix(src[i+1:], "?") && modifierGroup.MatchString(src[i+2:]):
			return "", errors.New("SyntaxError: Invalid regular expression: /" + src + "/: Invalid group")
		}
		b.WriteByte(c)
	}
	return b.String(), nil
}

var modifierGroup = regexp.MustCompile(`^[imsx]+(?:-[imsx]*)?:|^-[imsx]+:`)

// propertyLiterals is the number of regex literals in dist/shiki.js that
// use Unicode property escapes (in oniguruma-parser and oniguruma-to-es).
// A different count means the bundle changed: fail rather than guess.
const propertyLiterals = 7

// literalRe finds regex literals on one line: /body/flags preceded by a
// token after which a regex (not a division) starts.
var literalRe = regexp.MustCompile(`(^|[(,=:\[!&|?{};]|return)(\s*)/((?:\\.|\[(?:\\.|[^\]\\\n])*\]|[^/\\\[\n])+)/([dgimsuvy]*)`)

// rewriteLiterals rewrites the regex literals that use \p{...} with the
// u flag, which goja parses but cannot match.
func rewriteLiterals(src string, want int) (string, error) {
	count := 0
	var firstErr error
	out := literalRe.ReplaceAllStringFunc(src, func(m string) string {
		sub := literalRe.FindStringSubmatch(m)
		body, flags := sub[3], sub[4]
		if !strings.Contains(flags, "u") || !(strings.Contains(body, `\p{`) || strings.Contains(body, `\P{`)) {
			return m
		}
		r, err := rewritePattern(body, flags)
		if err != nil {
			if firstErr == nil {
				firstErr = err
			}
			return m
		}
		count++
		return sub[1] + sub[2] + "/" + r + "/" + flags
	})
	if firstErr != nil {
		return "", firstErr
	}
	if count != want {
		return "", fmt.Errorf("highlight: bundle changed: %d regex literals with Unicode property escapes, want %d", count, want)
	}
	return out, nil
}
