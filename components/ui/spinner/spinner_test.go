package spinner_test

import (
	"testing"

	"shadcn.go/components/ui/spinner"
	h "shadcn.go/libhtml"
)

// The expected markup is what React renders for the upstream component
// (react-dom/server, lucide-react 0.474.0), attributes in the same order.
func TestSpinnerMatchesReact(t *testing.T) {
	got := h.MustRenderString(spinner.Spinner(h.Class("size-8 text-blue-500")))
	want := `<svg class="lucide lucide-loader-circle animate-spin size-8 text-blue-500" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" role="status" aria-label="Loading"><path d="M21 12a9 9 0 1 1-6.219-8.56"></path></svg>`
	if got != want {
		t.Fatalf("\n got %s\nwant %s", got, want)
	}
}
