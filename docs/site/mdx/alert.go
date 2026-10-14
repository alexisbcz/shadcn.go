package mdx

// Alert, AlertTitle, AlertDescription: new-york-v4 alert.
//
// Placeholder: content.Placeholder (an <mdx-component> element with the
// props). Replace with the real rendering.
func init() {
	Register("Alert", placeholder("Alert"))
	Register("AlertTitle", placeholder("AlertTitle"))
	Register("AlertDescription", placeholder("AlertDescription"))
}
