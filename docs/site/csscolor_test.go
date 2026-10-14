package site

import "testing"

// Values from the reference build's stylesheet (lightningcss output).
func TestLabColors(t *testing.T) {
	for in, want := range map[string]string{
		"--surface: oklch(0.98 0 0)":           "--surface: lab(97.68% -.0000298023 .0000119209)",
		"oklch(97.1% 0.013 17.38)":             "lab(96.5005% 4.18508 1.52329)",
		"oklch(1 0 0 / 10%)":                   "lab(100% 0 0/.1)",
		"oklch(from currentColor l c h / 0.2)": "oklch(from currentColor l c h / 0.2)",
	} {
		if got := labColors(in); got != want {
			t.Errorf("labColors(%q) = %q, want %q", in, got, want)
		}
	}
}
