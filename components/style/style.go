// Package style is the Go port of the shadcn/ui style system.
//
// Upstream writes each component once, in registry/bases/<base>/ui, with
// semantic cn-* class tokens ("cn-button", "cn-card-header"), and its build
// turns that source into one component per style: transformStyle replaces
// each token with the Tailwind classes that registry/styles/style-<name>.css
// gives it (createStyleMap), then prettier-plugin-tailwindcss sorts the
// class strings. This package does the same at render time, so a Go
// component is also written once, with the tokens of the base source, and
// renders in any of the 8 styles exactly as upstream's generated
// styles/<base>-<style>/ui/*.tsx does:
//
//	var file = style.NewFile("ui/card.tsx")
//
//	var cardClass = file.Cn("cn-card group/card flex flex-col", style.Arg)
//
//	func Card(st style.Style, className string) h.Node {
//		return h.Div(h.Data("slot", "card"), h.Class(cardClass.Class(st, className)))
//	}
//
// See File for cva() definitions and the other kinds of class sites, and
// Default, SetDefault, NewContext and Middleware for choosing the style.
//
// The 8 style stylesheets are embedded verbatim (CSS), with the Tailwind
// class order prettier used (order.txt.gz). `go generate` refreshes both from
// the upstream clone in ./ui; the tests check they are current.
package style

//go:generate go run ./gen -root ../..

import (
	"context"
	"embed"
	"net/http"
	"sync"
	"sync/atomic"
)

// Style is one of upstream's 8 visual styles.
type Style string

// The styles, in the order of upstream's registry/styles.tsx.
const (
	Vega Style = "vega"
	Nova Style = "nova"
	Maia Style = "maia"
	Lyra Style = "lyra"
	Mira Style = "mira"
	Luma Style = "luma"
	Sera Style = "sera"
	Rhea Style = "rhea"
)

// All lists every style.
var All = []Style{Vega, Nova, Maia, Lyra, Mira, Luma, Sera, Rhea}

// index returns the position of s in All, or -1.
func (s Style) index() int {
	for i, x := range All {
		if x == s {
			return i
		}
	}
	return -1
}

// Valid reports whether s is one of the 8 styles.
func (s Style) Valid() bool { return s.index() >= 0 }

// Parse returns the style with the given name ("nova", "vega", ...).
func Parse(name string) (Style, bool) {
	s := Style(name)
	return s, s.Valid()
}

//go:embed css/*.css
var cssFiles embed.FS

// CSS returns the style's stylesheet, ui/apps/v4/registry/styles/style-<name>.css
// upstream, byte for byte. It returns "" for an invalid style.
func (s Style) CSS() string {
	if !s.Valid() {
		return ""
	}
	b, err := cssFiles.ReadFile("css/style-" + string(s) + ".css")
	if err != nil {
		panic(err)
	}
	return string(b)
}

var (
	mapsOnce  [8]sync.Once
	styleMaps [8]map[string]string
)

// Map returns the style's token map, CreateStyleMap of its stylesheet. The
// map is shared: do not modify it. It is nil for an invalid style.
func (s Style) Map() map[string]string {
	i := s.index()
	if i < 0 {
		return nil
	}
	mapsOnce[i].Do(func() {
		m, err := CreateStyleMap(s.CSS())
		if err != nil {
			panic(err)
		}
		styleMaps[i] = m
	})
	return styleMaps[i]
}

// String returns the style's name. Components accept a Style as a child,
// like their other typed props: button.Button(style.Vega, "Save").
func (s Style) String() string { return string(s) }

var defaultStyle atomic.Value // Style

// Default returns the style used when none is given: Nova, the style of the
// documentation (radix-nova), unless SetDefault changed it.
func Default() Style {
	if s, ok := defaultStyle.Load().(Style); ok {
		return s
	}
	return Nova
}

// SetDefault sets the process-wide default style. It panics on an invalid
// style.
func SetDefault(s Style) {
	if !s.Valid() {
		panic("style: invalid style " + string(s))
	}
	defaultStyle.Store(s)
}

// Or returns s, or the default style when s is empty.
func (s Style) Or() Style {
	if s == "" {
		return Default()
	}
	return s
}

type ctxKey struct{}

// NewContext returns a context carrying the style, for a request that
// renders in a style other than the default.
func NewContext(ctx context.Context, s Style) context.Context {
	return context.WithValue(ctx, ctxKey{}, s)
}

// FromContext returns the context's style, or the default.
func FromContext(ctx context.Context) Style {
	if s, ok := ctx.Value(ctxKey{}).(Style); ok && s.Valid() {
		return s
	}
	return Default()
}

// Middleware picks the style of each request with pick (for example from a
// cookie, a query parameter or the URL) and stores it in the request's
// context, where handlers read it with FromContext. Requests for which pick
// reports false use the default.
func Middleware(pick func(*http.Request) (Style, bool)) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if s, ok := pick(r); ok && s.Valid() {
				r = r.WithContext(NewContext(r.Context(), s))
			}
			next.ServeHTTP(w, r)
		})
	}
}
