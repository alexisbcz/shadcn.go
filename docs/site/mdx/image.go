package mdx

// Image: next/image (props className, src, width, height, alt, ... as merged by package content).
//
// Placeholder: content.Placeholder (an <mdx-component> element with the
// props). Replace with the real rendering.
func init() {
	Register("Image", placeholder("Image"))
}
