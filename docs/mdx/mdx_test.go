package mdx

import "testing"

// TestUnverified checks that constructs the docs content does not use are
// rejected.
func TestUnverified(t *testing.T) {
	for _, src := range []string{
		"~~gone~~\n",
		"line  \nbreak\n",
		"line\\\nbreak\n",
		"[a](b \"title\")\n",
		"see www.example.com\n",
		"mail me@example.com\n",
		"a &amp; b\n",
		"<!-- comment -->\n",
		"<https://example.com>\n",
		"Title\n---\n",
		"[a\nb]: /c\n",
		"[a\\]]: /b\n",
		"import { type A } from \"b\"\n",
		"- a\n<X\n/>\n",
		"a `b\n\tc` d\n",
		"-\ta\n",
	} {
		if _, err := Parse(src); err == nil {
			t.Errorf("expected an error for %q", src)
		}
	}
	for _, src := range []string{
		"see https://example.com/x.\n",
		"<Steps>\n\n### A\n\nText *em* **strong** `code` [l](/x)\n\n</Steps>\n",
		"| a | b |\n| - | :-: |\n| `x\\|y` | <br /> |\n",
		"text[^1]\n\n[^1]: note\n",
		"{/* comment */}\n\n<div>{[\"a\"].map((x) => (<b key={x}>{x}</b>))}</div>\n",
	} {
		if _, err := Parse(src); err != nil {
			t.Errorf("%q: %v", src, err)
		}
	}
}
