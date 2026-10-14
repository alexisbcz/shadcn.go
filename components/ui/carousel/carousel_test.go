package carousel_test

import (
	"strings"
	"testing"

	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/carousel"
	h "shadcn.go/libhtml"
)

func render(t *testing.T, n h.Node) string {
	t.Helper()
	s, err := h.RenderString(n)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func TestHorizontal(t *testing.T) {
	got := render(t, carousel.Carousel(h.Class("w-full max-w-xs"),
		carousel.Content(carousel.Item("1")),
		carousel.Previous(),
		carousel.Next(),
	))
	for _, want := range []string{
		`class="relative w-full max-w-xs"`,
		`role="region"`,
		`aria-roledescription="carousel"`,
		`data-carousel="{&#34;axis&#34;:&#34;x&#34;}"`,
		`<div class="overflow-hidden" data-slot="carousel-content"><div class="flex -ml-4">`,
		`<div class="min-w-0 shrink-0 grow-0 basis-full pl-4" role="group" aria-roledescription="slide" data-slot="carousel-item">1</div>`,
		`data-slot="carousel-previous"`,
		`top-1/2 -left-12 -translate-y-1/2`,
		`<span class="sr-only">Previous slide</span>`,
		`<span class="sr-only">Next slide</span>`,
		`lucide-arrow-left`, `lucide-arrow-right`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q in\n%s", want, got)
		}
	}
	if strings.Count(got, " disabled ") != 2 {
		t.Errorf("both buttons should render disabled (React's first render):\n%s", got)
	}
}

func TestVerticalAndOptions(t *testing.T) {
	got := render(t, carousel.Carousel(carousel.OrientationVertical, carousel.Opts{Align: "start", Loop: true},
		carousel.Autoplay{Delay: 2000, StopOnHover: true},
		carousel.Content(h.Class("-mt-1 h-[200px]"), carousel.Item(h.Class("pt-1 md:basis-1/2"), "1")),
		carousel.Previous(button.VariantGhost, h.Class("x-1")),
		carousel.Next(),
	))
	for _, want := range []string{
		`&#34;axis&#34;:&#34;y&#34;`,
		`&#34;align&#34;:&#34;start&#34;`,
		`&#34;loop&#34;:true`,
		`&#34;autoplay&#34;:{&#34;delay&#34;:2000}`,
		`&#34;stopOnHover&#34;:true`,
		`<div class="flex flex-col -mt-1 h-[200px]">`,
		`class="min-w-0 shrink-0 grow-0 basis-full pt-1 md:basis-1/2"`,
		`-top-12 left-1/2 -translate-x-1/2 rotate-90`,
		`-bottom-12 left-1/2 -translate-x-1/2 rotate-90`,
		`data-variant="ghost"`,
		` x-1`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q in\n%s", want, got)
		}
	}
}

func TestNavChildrenReplaceIcon(t *testing.T) {
	got := render(t, carousel.Carousel(carousel.Content(), carousel.Next("Go")))
	if strings.Contains(got, "Next slide") || !strings.Contains(got, ">Go</button>") {
		t.Errorf("children should replace the icon and label:\n%s", got)
	}
}
