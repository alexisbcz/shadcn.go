package chart

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"shadcn.go/components/icons/lucide"
	h "shadcn.go/libhtml"
)

func render(t *testing.T, n h.Node) string {
	t.Helper()
	var b bytes.Buffer
	if err := n.Render(&b); err != nil {
		t.Fatal(err)
	}
	return b.String()
}

func TestStyleMatchesUpstream(t *testing.T) {
	css, err := Config{
		{Key: "views", Label: "Page Views"},
		{Key: "desktop", Label: "Desktop", Color: "var(--chart-1)"},
		{Key: "mobile", Label: "Mobile", Theme: &Theme{Light: "#fff", Dark: "#000"}},
	}.css("chart-_r_0_")
	if err != nil {
		t.Fatal(err)
	}
	want := "\n [data-chart=chart-_r_0_] {\n  --color-desktop: var(--chart-1);\n  --color-mobile: #fff;\n}\n" +
		"\n" +
		"\n.dark [data-chart=chart-_r_0_] {\n  --color-desktop: var(--chart-1);\n  --color-mobile: #000;\n}\n"
	if css != want {
		t.Fatalf("css = %q\nwant %q", css, want)
	}
	if css, _ := (Config{{Key: "visitors", Label: "Visitors"}}).css("x"); css != "" {
		t.Fatalf("no colors: css = %q", css)
	}
}

func TestStyleRejectsInjection(t *testing.T) {
	for _, c := range []Config{
		{{Key: "a", Color: "red;} body{display:none"}},
		{{Key: "a b", Color: "red"}},
		{{Key: "a", Color: "</style><script>"}},
	} {
		if err := Container(c, BarChart()).Render(&bytes.Buffer{}); err == nil {
			t.Fatalf("config %+v rendered", c)
		}
		if _, err := c.css("chart-x"); err == nil {
			t.Fatalf("config %+v accepted", c)
		}
	}
}

func TestContainerMarkup(t *testing.T) {
	out := render(t, Container(Config{{Key: "desktop", Label: "Desktop", Color: "#2563eb", Icon: lucide.Activity}}, Nonce("n0nce"), h.Class("min-h-[200px] w-full"),
		BarChart(Data([]Row{{"month": "January", "desktop": 186}}), Bar(DataKey("desktop"), Fill("var(--color-desktop)"), Radius(4))),
	))
	for _, want := range []string{
		`data-slot="chart"`,
		`data-chart="chart-_r_0_"`,
		`min-h-[200px] w-full`,
		`<style nonce="n0nce">`,
		`class="recharts-responsive-container" style="width: 100%; height: 100%; min-width: 0px;"`,
		`<script type="application/json" data-chart-spec="">`,
		`<template data-chart-icon="desktop" data-parity-ignore="">`,
		`lucide-activity`,
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %s in\n%s", want, out)
		}
	}
}

func TestSpecJSONEscapesScript(t *testing.T) {
	out := render(t, Container(Config{{Key: "a", Label: "</script><b>"}}, BarChart(Data([]Row{{"name": "</script>"}}))))
	if strings.Count(out, "</script>") != 1 {
		t.Fatalf("spec can end the script element:\n%s", out)
	}
}

func TestSpecTree(t *testing.T) {
	n := AreaChart(Data([]Row{{"month": "Jan"}}), Margin{Left: 12, Right: 12},
		XAxis(DataKey("month"), TickFormatter(Slice(0, 3))),
		Line(Dot(IconDot(lucide.GitCommitVertical, 24, Props{"stroke": "red", "fill": Field("fill")}))),
	)
	var templates int
	resolved := resolveIcons(n, func(Icon) string { templates++; return "fn-0" })
	b, err := json.Marshal(resolved)
	if err != nil {
		t.Fatal(err)
	}
	s := string(b)
	for _, want := range []string{
		`"type":"AreaChart"`,
		`"margin":{"bottom":0,"left":12,"right":12,"top":0}`,
		`"tickFormatter":{"$fn":"slice","args":[0,3]}`,
		`"dot":{"$fn":"iconDot","args":["fn-0",24,{"fill":{"$field":"fill"},"stroke":"red"}]}`,
	} {
		if !strings.Contains(s, want) {
			t.Errorf("missing %s in %s", want, s)
		}
	}
	if templates != 1 {
		t.Errorf("templates = %d", templates)
	}
}

func TestTooltipClasses(t *testing.T) {
	n := TooltipContent(Indicator("dashed"), ClassName("w-[150px]"))
	c := n.Props["classes"].(map[string]any)
	if got := c["root"]; got != "grid min-w-[8rem] items-start gap-1.5 rounded-lg border border-border/50 bg-background px-2.5 py-1.5 text-xs shadow-xl w-[150px]" {
		t.Errorf("root = %q", got)
	}
	if got := c["nest"].(map[string]string)["indicator"]; got != "shrink-0 rounded-[2px] border-(--color-border) w-0 border-[1.5px] border-dashed bg-transparent my-0.5" {
		t.Errorf("nested dashed indicator = %q", got)
	}
	if got := c["item"]; got != "flex w-full flex-wrap items-stretch gap-2 [&>svg]:h-2.5 [&>svg]:w-2.5 [&>svg]:text-muted-foreground" {
		t.Errorf("item = %q", got)
	}
}
