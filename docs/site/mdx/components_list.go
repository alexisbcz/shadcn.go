package mdx

// ComponentsList: ui/apps/v4/components/components-list.tsx. Props: componentsFolder (content.PageTreeRef; resolve it with env.Tree.Find("components/"+base)), currentBase, variant.
//
// Placeholder: content.Placeholder (an <mdx-component> element with the
// props). Replace with the real rendering.
func init() {
	Register("ComponentsList", placeholder("ComponentsList"))
}
