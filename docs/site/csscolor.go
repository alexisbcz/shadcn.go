package site

import (
	"math"
	"regexp"
	"strconv"
	"strings"
)

// The reference stylesheet goes through lightningcss (Next.js), which
// rewrites every oklch() color as lab() (behind a hex fallback, in an
// @supports block that every current browser takes). The lab values differ
// from the oklch ones by float32 rounding, which is invisible on solid
// fills but shifts the dithering of large gradients (section-soft), so the
// site stylesheet gets the same lab() values. Relative colors (oklch(from
// ...)) are left alone, as lightningcss does.

var reOKLCHColor = regexp.MustCompile(`oklch\(\s*([\d.]+%?)\s+([\d.]+%?)\s+([\d.]+)(?:deg)?\s*(?:/\s*([\d.]+%?)\s*)?\)`)

// labColors rewrites the oklch() colors of css as lightningcss does.
func labColors(css string) string {
	return reOKLCHColor.ReplaceAllStringFunc(css, func(m string) string {
		g := reOKLCHColor.FindStringSubmatch(m)
		l := cssNumber(g[1], 1)
		c := cssNumber(g[2], 0.4)
		h, err := strconv.ParseFloat(g[3], 32)
		if err != nil {
			return m
		}
		L, a, b := oklchToLab(float32(l), float32(c), float32(h))
		out := "lab(" + labNumber(L) + "% " + labNumber(a) + " " + labNumber(b)
		if g[4] != "" {
			out += "/" + labNumber(float32(cssNumber(g[4], 1)))
		}
		return out + ")"
	})
}

// cssNumber parses a number or a percentage of full (100% = full).
func cssNumber(s string, full float64) float64 {
	if p, ok := strings.CutSuffix(s, "%"); ok {
		v, _ := strconv.ParseFloat(p, 64)
		return v / 100 * full
	}
	v, _ := strconv.ParseFloat(s, 64)
	return v
}

// labNumber formats like lightningcss: 6 significant digits, no exponent,
// no leading zero.
func labNumber(v float32) string {
	r, _ := strconv.ParseFloat(strconv.FormatFloat(float64(v), 'g', 6, 32), 64)
	s := strconv.FormatFloat(r, 'f', -1, 64)
	if s == "-0" {
		s = "0"
	}
	if strings.HasPrefix(s, "0.") {
		s = s[1:]
	} else if strings.HasPrefix(s, "-0.") {
		s = "-" + s[2:]
	}
	return s
}

// oklchToLab converts in float32 (as lightningcss): OKLCh -> OKLab -> LMS
// -> XYZ D65 -> XYZ D50 (Bradford) -> CIE Lab.
func oklchToLab(l, c, hdeg float32) (float32, float32, float32) {
	h := float64(hdeg) * math.Pi / 180
	a := c * float32(math.Cos(h))
	b := c * float32(math.Sin(h))
	l1 := l + 0.3963377773761749*a + 0.21580375730991364*b
	m1 := l - 0.10556134581565857*a - 0.0638541728258133*b
	s1 := l - 0.08948418209496575*a - 1.2914855480194092*b
	L, M, S := l1*l1*l1, m1*m1*m1, s1*s1*s1
	x := 1.2268798758459243*L - 0.5578149944602171*M + 0.2813910456659647*S
	y := -0.0405757452148008*L + 1.1122868032803170*M - 0.0717110580655164*S
	z := -0.0763729366746601*L - 0.4214933324022432*M + 1.5869240198367816*S
	x50 := 1.0479298208405488*x + 0.022946793341019088*y - 0.05019222954313557*z
	y50 := 0.029627815688159344*x + 0.990434484573249*y - 0.01707382502938514*z
	z50 := -0.009243058152591178*x + 0.015055144896577895*y + 0.7518742899580008*z
	const e, k = float32(216.0 / 24389.0), float32(24389.0 / 27.0)
	f := func(v float32) float32 {
		if v > e {
			return float32(math.Cbrt(float64(v)))
		}
		return (k*v + 16) / 116
	}
	fx, fy, fz := f(x50/float32(0.3457/0.3585)), f(y50), f(z50/float32((1.0-0.3457-0.3585)/0.3585))
	return 116*fy - 16, 500 * (fx - fy), 200 * (fy - fz)
}
