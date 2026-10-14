package cva

import (
	"reflect"
	"testing"
)

var button = New("base", Config{
	Variants: Variants{
		"variant": {"default": "bg-primary", "outline": "border"},
		"size":    {"default": "h-9", "sm": "h-8"},
	},
	Defaults: Selection{"variant": "default", "size": "default"},
	Compound: []Compound{{When: Selection{"variant": "outline", "size": "sm"}, Class: "shadow-xs"}},
	Order:    []string{"variant", "size"},
})

func TestClasses(t *testing.T) {
	cases := []struct {
		sel   Selection
		extra []string
		want  string
	}{
		{nil, nil, "base bg-primary h-9"},
		{Selection{"size": "sm"}, nil, "base bg-primary h-8"},
		{Selection{"variant": "outline", "size": "sm"}, []string{"  w-full "}, "base border h-8 shadow-xs w-full"},
		{Selection{"variant": "nope"}, nil, "base h-9"},
		{Selection{"variant": ""}, nil, "base bg-primary h-9"},
	}
	for _, c := range cases {
		if got := button.Classes(c.sel, c.extra...); got != c.want {
			t.Errorf("%v: got %q want %q", c.sel, got, c.want)
		}
	}
}

func TestOptions(t *testing.T) {
	if got := button.Options("size"); !reflect.DeepEqual(got, []string{"default", "sm"}) {
		t.Fatal(got)
	}
}
