// Command goexamples extracts the Go source of every docs example: for each
// example registered in examples/ (registry.Register), the function that
// renders it, with the imports it uses rewritten to the published module
// paths. It writes docs/site/data/go-examples.json ({"nova--button-demo":
// "package main\n\nimport (...)\n\nfunc Demo() ..."}), which the docs site
// shows in its code panels' Go tab.
//
//	go run ./docs/site/gen/goexamples
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// published maps monorepo import paths to the public modules.
func published(path string) string {
	switch {
	case strings.HasPrefix(path, "shadcn.go/components/icons"):
		return "github.com/alexisbcz/icons" + strings.TrimPrefix(path, "shadcn.go/components/icons")
	case strings.HasPrefix(path, "shadcn.go/components"):
		return "github.com/alexisbcz/components" + strings.TrimPrefix(path, "shadcn.go/components")
	case strings.HasPrefix(path, "shadcn.go/"):
		return "github.com/alexisbcz/" + strings.TrimPrefix(path, "shadcn.go/")
	}
	return path
}

type pkg struct {
	fset  *token.FileSet
	files map[string]*ast.File // by file name
	funcs map[string]*ast.FuncDecl
	file  map[string]*ast.File // func name -> its file
	src   map[string][]byte
}

func main() {
	root := "."
	if len(os.Args) > 1 {
		root = os.Args[1]
	}
	out := map[string]string{}
	dirs, _ := filepath.Glob(filepath.Join(root, "examples", "*"))
	for _, dir := range dirs {
		if st, err := os.Stat(dir); err != nil || !st.IsDir() {
			continue
		}
		p := load(dir)
		if p == nil {
			continue
		}
		for name, fn := range p.registrations() {
			if code, err := p.source(fn); err == nil {
				out[name] = code
			}
		}
	}
	b, _ := json.MarshalIndent(out, "", " ")
	dst := filepath.Join(root, "docs/site/mdx/go-examples.json")
	if err := os.WriteFile(dst, b, 0o644); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s: %d examples\n", dst, len(out))
}

func load(dir string) *pkg {
	p := &pkg{fset: token.NewFileSet(), files: map[string]*ast.File{}, funcs: map[string]*ast.FuncDecl{}, file: map[string]*ast.File{}, src: map[string][]byte{}}
	matches, _ := filepath.Glob(filepath.Join(dir, "*.go"))
	for _, f := range matches {
		if strings.HasSuffix(f, "_test.go") {
			continue
		}
		b, err := os.ReadFile(f)
		if err != nil {
			continue
		}
		af, err := parser.ParseFile(p.fset, f, b, parser.ParseComments)
		if err != nil {
			continue
		}
		p.files[f] = af
		p.src[f] = b
		for _, d := range af.Decls {
			if fd, ok := d.(*ast.FuncDecl); ok && fd.Recv == nil {
				p.funcs[fd.Name.Name] = fd
				p.file[fd.Name.Name] = af
			}
		}
	}
	if len(p.files) == 0 {
		return nil
	}
	return p
}

// registrations finds name -> function for the package's examples: calls
// of a helper taking ("name", Func), with the prefix the helper adds
// ("nova--", "base--", "aria--"), and registry.Example literals.
func (p *pkg) registrations() map[string]string {
	out := map[string]string{}
	// Helpers: local closures or functions whose body builds the name with
	// a "<prefix>--" + name concatenation.
	prefixOf := map[string]string{}
	for _, f := range p.files {
		ast.Inspect(f, func(n ast.Node) bool {
			var name string
			var body ast.Node
			switch x := n.(type) {
			case *ast.AssignStmt:
				if len(x.Lhs) == 1 && len(x.Rhs) == 1 {
					if id, ok := x.Lhs[0].(*ast.Ident); ok {
						if fl, ok := x.Rhs[0].(*ast.FuncLit); ok {
							name, body = id.Name, fl.Body
						}
					}
				}
			case *ast.FuncDecl:
				name, body = x.Name.Name, x.Body
			}
			if body == nil {
				return true
			}
			ast.Inspect(body, func(m ast.Node) bool {
				if be, ok := m.(*ast.BinaryExpr); ok && be.Op == token.ADD {
					if lit, ok := be.X.(*ast.BasicLit); ok && lit.Kind == token.STRING {
						s, _ := strconv.Unquote(lit.Value)
						if strings.HasSuffix(s, "--") {
							prefixOf[name] = s
						}
					}
				}
				return true
			})
			return true
		})
	}
	for _, f := range p.files {
		ast.Inspect(f, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.CallExpr:
				if len(x.Args) != 2 {
					return true
				}
				lit, ok := x.Args[0].(*ast.BasicLit)
				if !ok || lit.Kind != token.STRING {
					return true
				}
				fn := p.resolve(f, x.Args[1])
				if fn == "" {
					return true
				}
				s, _ := strconv.Unquote(lit.Value)
				helper := ""
				if id, ok := x.Fun.(*ast.Ident); ok {
					helper = id.Name
				}
				out[prefixOf[helper]+s] = fn
			case *ast.CompositeLit:
				// registry.Example{Name: "x", Render: registry.Static(Fn)}
				var name, fn string
				for _, e := range x.Elts {
					kv, ok := e.(*ast.KeyValueExpr)
					if !ok {
						continue
					}
					key, _ := kv.Key.(*ast.Ident)
					if key == nil {
						continue
					}
					switch key.Name {
					case "Name":
						if lit, ok := kv.Value.(*ast.BasicLit); ok && lit.Kind == token.STRING {
							name, _ = strconv.Unquote(lit.Value)
						}
					case "Render":
						if call, ok := kv.Value.(*ast.CallExpr); ok && len(call.Args) == 1 {
							fn = p.resolve(f, call.Args[0])
						}
					}
				}
				if name != "" && fn != "" && p.funcs[fn] != nil {
					out[name] = fn
				}
			}
			return true
		})
	}
	return out
}

// resolve returns the package function an example renders: the function
// itself, or the first package function a local closure (or function
// literal) calls.
func (p *pkg) resolve(f *ast.File, e ast.Expr) string {
	var body ast.Node
	switch x := e.(type) {
	case *ast.Ident:
		if p.funcs[x.Name] != nil {
			return x.Name
		}
		// A local closure: name := func() h.Node { ... }
		ast.Inspect(f, func(n ast.Node) bool {
			if as, ok := n.(*ast.AssignStmt); ok && len(as.Lhs) == 1 && len(as.Rhs) == 1 {
				if id, ok := as.Lhs[0].(*ast.Ident); ok && id.Name == x.Name {
					if fl, ok := as.Rhs[0].(*ast.FuncLit); ok {
						body = fl.Body
					}
				}
			}
			return body == nil
		})
	case *ast.FuncLit:
		body = x.Body
	}
	if body == nil {
		return ""
	}
	found := ""
	ast.Inspect(body, func(n ast.Node) bool {
		if call, ok := n.(*ast.CallExpr); ok && found == "" {
			if id, ok := call.Fun.(*ast.Ident); ok && p.funcs[id.Name] != nil {
				found = id.Name
			}
		}
		return found == ""
	})
	return found
}

func funcName(e ast.Expr) string {
	if id, ok := e.(*ast.Ident); ok {
		return id.Name
	}
	return ""
}

// source returns a standalone snippet: the imports it uses, the function
// (with its doc comment), then the package-level declarations it depends
// on, recursively (helpers, types, data).
func (p *pkg) source(fn string) (string, error) {
	decls := p.topLevel()
	var order []string
	seen := map[string]bool{}
	var visit func(name string)
	visit = func(name string) {
		if seen[name] {
			return
		}
		seen[name] = true
		order = append(order, name)
		ast.Inspect(decls[name].node, func(n ast.Node) bool {
			if id, ok := n.(*ast.Ident); ok && id.Obj == nil {
				// Unresolved in the file scope: package level, or another
				// file of the package.
			}
			if id, ok := n.(*ast.Ident); ok {
				if d, ok := decls[id.Name]; ok && id.Name != name && !seen[id.Name] && d.node != nil {
					visit(id.Name)
				}
			}
			return true
		})
	}
	visit(fn)
	used := map[string]bool{}
	var parts []string
	files := map[*ast.File]bool{}
	for _, name := range order {
		d := decls[name]
		files[d.file] = true
		ast.Inspect(d.node, func(n ast.Node) bool {
			if se, ok := n.(*ast.SelectorExpr); ok {
				if id, ok := se.X.(*ast.Ident); ok {
					used[id.Name] = true
				}
			}
			return true
		})
		parts = append(parts, d.text)
	}
	imports := map[string]bool{}
	for af := range files {
		for _, im := range af.Imports {
			path, _ := strconv.Unquote(im.Path.Value)
			name := filepath.Base(path)
			if im.Name != nil {
				name = im.Name.Name
			} else if alias := pkgName(path); alias != "" {
				name = alias
			}
			if !used[name] || strings.Contains(path, "/examples/") {
				continue
			}
			line := strconv.Quote(published(path))
			if im.Name != nil {
				line = im.Name.Name + " " + line
			}
			imports[line] = true
		}
	}
	lines := make([]string, 0, len(imports))
	for l := range imports {
		lines = append(lines, l)
	}
	sort.Strings(lines)
	var b bytes.Buffer
	if len(lines) > 0 {
		b.WriteString("import (\n")
		for _, l := range lines {
			b.WriteString("\t" + l + "\n")
		}
		b.WriteString(")\n\n")
	}
	b.WriteString(strings.Join(parts, "\n\n"))
	b.WriteString("\n")
	out, err := format.Source(append([]byte("package example\n\n"), b.Bytes()...))
	if err != nil {
		return "", err
	}
	return strings.TrimPrefix(string(out), "package example\n\n"), nil
}

type decl struct {
	node ast.Node
	file *ast.File
	text string
}

// topLevel indexes the package's declarations (functions, types, vars,
// consts) by name, with their source (doc comment included). init and
// registration code is left out.
func (p *pkg) topLevel() map[string]decl {
	out := map[string]decl{}
	for fname, af := range p.files {
		src := p.src[fname]
		text := func(n ast.Node, doc *ast.CommentGroup) string {
			start := n.Pos()
			if doc != nil {
				start = doc.Pos()
			}
			return string(src[p.fset.Position(start).Offset:p.fset.Position(n.End()).Offset])
		}
		for _, d := range af.Decls {
			switch x := d.(type) {
			case *ast.FuncDecl:
				if x.Recv == nil && x.Name.Name != "init" {
					out[x.Name.Name] = decl{x, af, text(x, x.Doc)}
				}
			case *ast.GenDecl:
				if x.Tok == token.IMPORT {
					continue
				}
				for _, sp := range x.Specs {
					var names []string
					switch y := sp.(type) {
					case *ast.TypeSpec:
						names = []string{y.Name.Name}
					case *ast.ValueSpec:
						for _, n := range y.Names {
							names = append(names, n.Name)
						}
					}
					// The whole declaration (a var block declares several).
					t := text(x, x.Doc)
					if !x.Lparen.IsValid() {
						t = text(x, x.Doc)
					}
					for _, n := range names {
						if n != "_" {
							out[n] = decl{x, af, t}
						}
					}
				}
			}
		}
	}
	return out
}

// pkgName is the package name of a path whose last element differs from
// it ("shadcn.go/components/ui/dropdown-menu" -> "dropdownmenu").
func pkgName(path string) string {
	base := filepath.Base(path)
	if strings.Contains(base, "-") {
		return strings.ReplaceAll(base, "-", "")
	}
	return ""
}
