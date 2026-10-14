package cn

import (
	"math"
	"math/big"
	"strconv"
	"unicode/utf8"
)

// jsNumber is JavaScript's Number(string) (ECMAScript StringToNumber). ok is
// false where JavaScript yields NaN.
func jsNumber(s string) (v float64, ok bool) {
	s = trimJSSpace(s)
	if s == "" {
		return 0, true
	}
	if len(s) > 2 && s[0] == '0' {
		base := 0
		switch s[1] {
		case 'x', 'X':
			base = 16
		case 'o', 'O':
			base = 8
		case 'b', 'B':
			base = 2
		}
		if base != 0 {
			for i := 2; i < len(s); i++ {
				if digitVal(s[i]) >= base {
					return 0, false
				}
			}
			n, _ := new(big.Int).SetString(s[2:], base)
			f, _ := new(big.Float).SetInt(n).Float64()
			return f, true
		}
	}
	body := s
	if body[0] == '+' || body[0] == '-' {
		body = body[1:]
	}
	if body == "Infinity" {
		if s[0] == '-' {
			return math.Inf(-1), true
		}
		return math.Inf(1), true
	}
	if !isStrUnsignedDecimalLiteral(body) {
		return 0, false
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		if ne, isNE := err.(*strconv.NumError); !isNE || ne.Err != strconv.ErrRange {
			return 0, false
		}
	}
	return f, true
}

func digitVal(c byte) int {
	switch {
	case c >= '0' && c <= '9':
		return int(c - '0')
	case c >= 'a' && c <= 'f':
		return int(c-'a') + 10
	case c >= 'A' && c <= 'F':
		return int(c-'A') + 10
	}
	return 99
}

// isStrUnsignedDecimalLiteral reports whether s is
// digits [. digits?] [exp] | . digits [exp].
func isStrUnsignedDecimalLiteral(s string) bool {
	i := 0
	digits := func() int {
		start := i
		for i < len(s) && s[i] >= '0' && s[i] <= '9' {
			i++
		}
		return i - start
	}
	intDigits := digits()
	fracDigits := 0
	if i < len(s) && s[i] == '.' {
		i++
		fracDigits = digits()
	}
	if intDigits == 0 && fracDigits == 0 {
		return false
	}
	if i < len(s) && (s[i] == 'e' || s[i] == 'E') {
		i++
		if i < len(s) && (s[i] == '+' || s[i] == '-') {
			i++
		}
		if digits() == 0 {
			return false
		}
	}
	return i == len(s)
}

// isJSSpace reports whether r is StrWhiteSpaceChar: WhiteSpace (TAB, VT,
// FF, ZWNBSP, Zs) or LineTerminator.
func isJSSpace(r rune) bool {
	switch r {
	case '\t', '\n', '\v', '\f', '\r', ' ', 0xA0, 0x1680, 0x2028, 0x2029, 0x202F, 0x205F, 0x3000, 0xFEFF:
		return true
	}
	return r >= 0x2000 && r <= 0x200A
}

func trimJSSpace(s string) string {
	for len(s) > 0 {
		r, n := utf8.DecodeRuneInString(s)
		if r == utf8.RuneError && n <= 1 || !isJSSpace(r) {
			break
		}
		s = s[n:]
	}
	for len(s) > 0 {
		r, n := utf8.DecodeLastRuneInString(s)
		if r == utf8.RuneError && n <= 1 || !isJSSpace(r) {
			break
		}
		s = s[:len(s)-n]
	}
	return s
}

// isNumStr is the engine's isNumStr: !!v && !Number.isNaN(Number(v)).
func isNumStr(v string) bool {
	if v == "" {
		return false
	}
	_, ok := jsNumber(v)
	return ok
}

// isIntStr is !!v && Number.isInteger(Number(v)).
func isIntStr(v string) bool {
	if v == "" {
		return false
	}
	f, ok := jsNumber(v)
	return ok && !math.IsInf(f, 0) && f == math.Trunc(f)
}
