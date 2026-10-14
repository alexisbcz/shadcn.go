package site

import (
	_ "embed"
	"encoding/json"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/button"
	selectui "shadcn.go/components/ui/select"
	h "shadcn.go/libhtml"
)

// The colors page: app/(app)/colors/{layout,page}.tsx, lib/colors.ts,
// components/{color-palette,color,color-format-selector}.tsx. The palette
// is registry/_legacy-colors.ts, extracted by gen/colors.mjs.

//go:embed data/colors.json
var colorsJSON []byte

const colorsTitle = "Tailwind Colors in Every Format"
const colorsDescription = "The complete Tailwind color palette in HEX, RGB, HSL, CSS variables, and classes. Ready to copy and paste into your project."

// defaultColorFormat is useColors' initial format.
const defaultColorFormat = "hsl"

type colorSwatch struct {
	Name, ID, ClassName, Hex, RGB, HSL, Foreground, OKLCH, Var string
	Scale                                                      int
}

type colorPalette struct {
	Name   string
	Colors []colorSwatch
}

var (
	reRGB   = regexp.MustCompile(`^rgb\((\d+),(\d+),(\d+)\)$`)
	reHSL   = regexp.MustCompile(`^hsl\(([\d.]+),([\d.]+%),([\d.]+%)\)$`)
	reOKLCH = regexp.MustCompile(`^oklch\(([\d.]+)\s*,\s*([\d.]+)\s*,\s*([\d.]+)\)$`)
)

// getColors is lib/colors.ts getColors.
func getColors() []colorPalette {
	var raw []struct {
		Name   string `json:"name"`
		Colors []struct {
			Scale int    `json:"scale"`
			Hex   string `json:"hex"`
			RGB   string `json:"rgb"`
			HSL   string `json:"hsl"`
			OKLCH string `json:"oklch"`
		} `json:"colors"`
	}
	_ = json.Unmarshal(colorsJSON, &raw)
	out := make([]colorPalette, 0, len(raw))
	for _, p := range raw {
		cp := colorPalette{Name: p.Name}
		for _, c := range p.Colors {
			id := p.Name + "-" + strconv.Itoa(c.Scale)
			rgb := reRGB.ReplaceAllString(c.RGB, "$1 $2 $3")
			cp.Colors = append(cp.Colors, colorSwatch{
				Name: p.Name, ID: id, ClassName: id, Scale: c.Scale, Hex: c.Hex,
				RGB:        rgb,
				HSL:        reHSL.ReplaceAllString(c.HSL, "$1 $2 $3"),
				OKLCH:      "oklch(" + reOKLCH.ReplaceAllString(c.OKLCH, "$1 $2 $3") + ")",
				Var:        "--color-" + id,
				Foreground: foregroundFor(rgb),
			})
		}
		out = append(out, cp)
	}
	return out
}

func foregroundFor(rgb string) string {
	parts := strings.Fields(rgb)
	lin := func(s string) float64 {
		n, _ := strconv.ParseFloat(s, 64)
		b := n / 255
		if b <= 0.04045 {
			return b / 12.92
		}
		return math.Pow((b+0.055)/1.055, 2.4)
	}
	if len(parts) != 3 {
		return "#fff"
	}
	if 0.2126*lin(parts[0])+0.7152*lin(parts[1])+0.0722*lin(parts[2]) > 0.179 {
		return "#000"
	}
	return "#fff"
}

// colorFormats is getColorFormat, in its key order.
func colorFormats(c colorSwatch) [][2]string {
	return [][2]string{
		{"className", "bg-" + c.Name + "-100"},
		{"hex", c.Hex},
		{"rgb", c.RGB},
		{"hsl", c.HSL},
		{"oklch", c.OKLCH},
		{"var", "--color-" + c.Name + "-" + strconv.Itoa(c.Scale)},
	}
}

func colorValue(c colorSwatch, format string) string {
	for _, f := range colorFormats(c) {
		if f[0] == format {
			if format == "className" {
				return c.ClassName
			}
			if format == "var" {
				return c.Var
			}
			return f[1]
		}
	}
	return ""
}

// Colors serves /colors.
func (s *Site) Colors(r *http.Request) (h.Node, error) {
	var palettes []any
	for _, p := range getColors() {
		palettes = append(palettes, colorPaletteNode(p))
	}
	return s.Document(r, Meta{Title: colorsTitle, Description: colorsDescription, Canonical: "/colors"},
		s.AppLayout(r,
			h.Div(
				pageHeader(
					announcement(),
					pageHeaderHeading(colorsTitle),
					pageHeaderDescription(colorsDescription),
					// The colors layout's Button is styles/radix-luma/ui/button.
					pageActions(style.Scope(style.Luma,
						button.Button(lib.AsChild, h.Class("h-[35px]"), h.A(h.Href("#colors"), "Browse Colors")),
						button.Button(lib.AsChild, button.VariantSecondary, h.A(h.Href("/docs/theming"), "Documentation")),
					)),
				),
				h.Div(h.Class("hidden"),
					h.Div(h.Class("container-wrapper"),
						h.Div(h.Class("container flex items-center justify-between gap-8 py-4"), colorsNav()),
					),
				),
				h.Div(h.Class("container-wrapper"),
					h.Div(h.Class("container py-6"),
						h.Section(h.ID("colors"), h.Class("scroll-mt-20"),
							h.Div(h.Class("grid gap-8 lg:gap-16 xl:gap-20"), h.Group(palettes...)),
						),
					),
				),
			),
		),
	), nil
}

func colorsNav() h.Node {
	var links []any
	for i, p := range getColors() {
		links = append(links, h.A(h.Href("/colors#"+p.Name), h.Data("active", boolString(i == 0)),
			h.Class("flex h-7 items-center justify-center px-4 text-center text-base font-medium text-muted-foreground capitalize transition-colors hover:text-primary data-[active=true]:text-primary"),
			p.Name))
	}
	return h.Div(h.Class("flex items-center flex-1 overflow-hidden [&>a:first-child]:text-primary"),
		h.Div(h.Data("slot", "scroll-area"), h.Class("relative max-w-full"),
			h.Div(h.Data("slot", "scroll-area-viewport"), h.Class("size-full rounded-[inherit] transition-[color,box-shadow] outline-none focus-visible:ring-[3px] focus-visible:ring-ring/50 focus-visible:outline-1"), h.Style("overflow: scroll"),
				h.Div(h.Style("min-width: 100%; display: table"), h.Div(h.Class("flex items-center"), h.Group(links...))),
			),
		),
	)
}

func colorPaletteNode(p colorPalette) h.Node {
	var swatches []any
	for _, c := range p.Colors {
		swatches = append(swatches, colorNode(c))
	}
	return h.Div(h.ID(p.Name), h.Class("scroll-mt-20 rounded-lg"),
		h.Div(h.Class("flex items-center px-4"),
			h.Div(h.Class("flex-1 pl-1 text-sm font-medium"), h.H2(h.Class("capitalize"), p.Name)),
			colorFormatSelector(p.Colors[0]),
		),
		h.Div(h.Class("flex flex-col gap-4 py-4 sm:flex-row sm:gap-2"), h.Group(swatches...)),
	)
}

func colorFormatSelector(c colorSwatch) h.Node {
	var items []any
	for _, f := range colorFormats(c) {
		items = append(items, selectui.Item(selectui.ItemValue(f[0]),
			h.Class("gap-2 rounded-lg [&>span]:flex [&>span]:items-center [&>span]:gap-2"),
			h.Span(h.Class("font-medium"), f[0]),
			h.Span(h.Class("font-mono text-xs text-muted-foreground"), f[1]),
		))
	}
	return h.Div(h.Class("contents"), h.Data("color-format", ""), selectui.Select(selectui.DefaultValue(defaultColorFormat),
		selectui.Trigger(selectui.SizeSm, h.Class("border-secondary bg-secondary text-secondary-foreground shadow-none ml-auto"),
			h.Span(h.Class("font-medium"), "Format: "),
			h.Span(h.Class("font-mono text-muted-foreground"), h.Data("color-format-value", ""), defaultColorFormat),
		),
		selectui.Content(selectui.PositionPopper, selectui.AlignEnd, h.Class("rounded-xl"), h.Group(items...)),
	))
}

// colorNode is components/color.tsx; site.js keeps the format, copying and
// data-last-copied in sync on the client.
func colorNode(c colorSwatch) h.Node {
	values := map[string]string{}
	for _, f := range colorFormats(c) {
		values[f[0]] = colorValue(c, f[0])
	}
	b, _ := json.Marshal(values)
	return h.Button(h.Data("last-copied", "false"), h.Data("color-values", string(b)),
		h.Class("group relative flex aspect-[3/1] w-full flex-1 cursor-pointer flex-col gap-2 text-(--text) sm:aspect-[2/3] sm:h-auto sm:w-auto [&>svg]:absolute [&>svg]:top-4 [&>svg]:right-4 [&>svg]:z-10 [&>svg]:h-3.5 [&>svg]:w-3.5 [&>svg]:opacity-0 [&>svg]:transition-opacity"),
		h.Style("--bg: "+c.OKLCH+"; --text: "+c.Foreground),
		lucide.Clipboard(h.Class("group-hover:opacity-100")),
		h.Div(h.Class("border-ghost w-full flex-1 rounded-md bg-(--bg) after:rounded-lg after:border-input md:rounded-lg")),
		h.Div(h.Class("flex w-full flex-col items-center justify-center gap-1"),
			h.Span(h.Class("font-mono text-xs text-muted-foreground tabular-nums transition-colors group-hover:text-foreground group-data-[last-copied=true]:text-primary sm:hidden xl:flex"), c.ClassName),
			h.Span(h.Class("hidden font-mono text-xs text-muted-foreground tabular-nums transition-colors group-hover:text-foreground group-data-[last-copied=true]:text-primary sm:flex xl:hidden"), strconv.Itoa(c.Scale)),
		),
	)
}
