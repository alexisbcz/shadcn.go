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
	methods    map[string][]*ast.FuncDecl // receiver type -> methods
	methodFile map[*ast.FuncDecl]*ast.File
	fset       *token.FileSet
	files      map[string]*ast.File // by file name
	funcs      map[string]*ast.FuncDecl
	file       map[string]*ast.File // func name -> its file
	src        map[string][]byte
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
		for name, e := range p.registrations() {
			if code, err := p.source(e); err != nil {
				log.Printf("%s: %v", name, err)
			} else {
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
	p := &pkg{methods: map[string][]*ast.FuncDecl{}, methodFile: map[*ast.FuncDecl]*ast.File{}, fset: token.NewFileSet(), files: map[string]*ast.File{}, funcs: map[string]*ast.FuncDecl{}, file: map[string]*ast.File{}, src: map[string][]byte{}}
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
			fd, ok := d.(*ast.FuncDecl)
			if !ok {
				continue
			}
			if fd.Recv == nil {
				p.funcs[fd.Name.Name] = fd
				p.file[fd.Name.Name] = af
			} else if t := recvType(fd.Recv.List[0].Type); t != "" {
				p.methods[t] = append(p.methods[t], fd)
				p.methodFile[fd] = af
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
func (p *pkg) registrations() map[string]entry {
	out := map[string]entry{}
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
				if fn.fn == "" && fn.lit == nil {
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
				var name string
				var fn entry
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
				if name != "" && (fn.lit != nil || p.funcs[fn.fn] != nil) {
					out[name] = fn
				}
			}
			return true
		})
	}
	return out
}

// entry is what an example renders: a package function, or a function
// literal (the registration's own closure), which the snippet turns into a
// function named Example.
type entry struct {
	fn     string
	lit    *ast.FuncLit
	file   string
	consts []*ast.ValueSpec // local constants the literal uses
}

// resolve returns the entry of a registration's render argument: a package
// function, or a function literal (inline or a local closure).
func (p *pkg) resolve(f *ast.File, e ast.Expr) entry {
	fname := p.fset.File(f.Pos()).Name()
	switch x := e.(type) {
	case *ast.Ident:
		if p.funcs[x.Name] != nil {
			return entry{fn: x.Name}
		}
		// A local closure: name := func() h.Node { ... }
		var lit *ast.FuncLit
		ast.Inspect(f, func(n ast.Node) bool {
			if as, ok := n.(*ast.AssignStmt); ok && len(as.Lhs) == 1 && len(as.Rhs) == 1 {
				if id, ok := as.Lhs[0].(*ast.Ident); ok && id.Name == x.Name {
					if fl, ok := as.Rhs[0].(*ast.FuncLit); ok {
						lit = fl
					}
				}
			}
			return lit == nil
		})
		if lit != nil {
			return p.literal(lit, fname)
		}
	case *ast.FuncLit:
		return p.literal(x, fname)
	}
	return entry{}
}

// literal is the entry of a function literal; one that only calls a
// package function without arguments is that function.
func (p *pkg) literal(lit *ast.FuncLit, file string) entry {
	if len(lit.Body.List) == 1 {
		if ret, ok := lit.Body.List[0].(*ast.ReturnStmt); ok && len(ret.Results) == 1 {
			if call, ok := ret.Results[0].(*ast.CallExpr); ok && len(call.Args) == 0 {
				if id, ok := call.Fun.(*ast.Ident); ok && p.funcs[id.Name] != nil {
					return entry{fn: id.Name}
				}
			}
		}
	}
	// Closures referring to local variables of the registering function
	// cannot stand alone; local constants are copied into them.
	consts, ok := p.locals(lit)
	if !ok {
		return entry{}
	}
	return entry{lit: lit, file: file, consts: consts}
}

// locals returns the local constants lit uses, and false when it uses
// local variables (identifiers that are neither its own, package-level
// declarations nor imports).
func (p *pkg) locals(lit *ast.FuncLit) ([]*ast.ValueSpec, bool) {
	var consts []*ast.ValueSpec
	seenConst := map[*ast.ValueSpec]bool{}
	decls := p.topLevel()
	own := map[string]bool{}
	ast.Inspect(lit, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.AssignStmt:
			if x.Tok == token.DEFINE {
				for _, l := range x.Lhs {
					if id, ok := l.(*ast.Ident); ok {
						own[id.Name] = true
					}
				}
			}
		case *ast.Field:
			for _, n := range x.Names {
				own[n.Name] = true
			}
		case *ast.ValueSpec:
			for _, n := range x.Names {
				own[n.Name] = true
			}
		case *ast.RangeStmt:
			for _, e := range []ast.Expr{x.Key, x.Value} {
				if id, ok := e.(*ast.Ident); ok {
					own[id.Name] = true
				}
			}
		}
		return true
	})
	local := false
	ast.Inspect(lit, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.SelectorExpr:
			// Only the receiver of a selector can be a local.
			ast.Inspect(x.X, func(m ast.Node) bool {
				if id, ok := m.(*ast.Ident); ok && id == x.X {
					if id.Obj != nil && id.Obj.Kind != ast.Pkg && !own[id.Name] {
						if _, ok := decls[id.Name]; !ok {
							local = true
						}
					}
				}
				return true
			})
			return false
		case *ast.Ident:
			if x.Obj != nil && x.Obj.Kind == ast.Con && !own[x.Name] {
				if vs, ok := x.Obj.Decl.(*ast.ValueSpec); ok && !seenConst[vs] {
					if _, pkgLevel := decls[x.Name]; !pkgLevel || decls[x.Name].node.Pos() > vs.Pos() || decls[x.Name].node.End() < vs.End() {
						seenConst[vs] = true
						consts = append(consts, vs)
					}
				}
			}
			if x.Obj != nil && x.Obj.Kind == ast.Var && !own[x.Name] {
				if _, ok := decls[x.Name]; !ok {
					local = true
				}
			}
		}
		return true
	})
	return consts, !local
}

func funcName(e ast.Expr) string {
	if id, ok := e.(*ast.Ident); ok {
		return id.Name
	}
	return ""
}

// source returns a standalone snippet: the imports it uses, the function
// (with its doc comment), then the package-level declarations it depends
// on, recursively (helpers, types with their methods, data).
func (p *pkg) source(e entry) (string, error) {
	decls := p.topLevel()
	fn := e.fn
	if e.lit != nil {
		fn = "Example"
		src := p.src[e.file]
		off := func(pos token.Pos) int { return p.fset.Position(pos).Offset }
		body := string(src[off(e.lit.Body.Pos()):off(e.lit.Body.End())])
		var consts strings.Builder
		for _, vs := range e.consts {
			consts.WriteString("\n\tconst " + string(src[off(vs.Pos()):off(vs.End())]))
		}
		if consts.Len() > 0 {
			body = "{" + consts.String() + "\n" + body[1:]
		}
		text := "func Example" + string(src[off(e.lit.Type.Params.Pos()):off(e.lit.Type.End())]) + " " + body
		var af *ast.File
		for name, f := range p.files {
			if name == e.file {
				af = f
			}
		}
		decls[fn] = decl{e.lit, af, text}
	}
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
			if id, ok := n.(*ast.Ident); ok {
				if d, ok := decls[id.Name]; ok && id.Name != name && !seen[id.Name] && d.node != nil {
					visit(id.Name)
				}
			}
			return true
		})
		// A type's methods come with it.
		for _, m := range p.methods[name] {
			key := name + "." + m.Name.Name
			if !seen[key] {
				decls[key] = decl{m, p.methodFile[m], p.declText(m, m.Doc)}
				visit(key)
			}
		}
	}
	visit(fn)
	emitted := map[ast.Node]bool{}
	var kept []string
	for _, name := range order {
		if n := decls[name].node; !emitted[n] {
			emitted[n] = true
			kept = append(kept, name)
		}
	}
	order = kept
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
			if im.Name != nil && im.Name.Name != pkgName(path) {
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
	if len(base) > 1 && base[0] == 'v' && strings.Trim(base[1:], "0123456789") == "" {
		return filepath.Base(filepath.Dir(path))
	}
	if strings.Contains(base, "-") {
		return strings.ReplaceAll(base, "-", "")
	}
	return ""
}

// recvType is the type name of a method receiver (T, *T, T[P]).
func recvType(e ast.Expr) string {
	switch x := e.(type) {
	case *ast.StarExpr:
		return recvType(x.X)
	case *ast.IndexExpr:
		return recvType(x.X)
	case *ast.Ident:
		return x.Name
	}
	return ""
}

// declText is the source of a declaration with its doc comment.
func (p *pkg) declText(n ast.Node, doc *ast.CommentGroup) string {
	f := p.fset.File(n.Pos())
	src := p.src[f.Name()]
	start := n.Pos()
	if doc != nil {
		start = doc.Pos()
	}
	return string(src[p.fset.Position(start).Offset:p.fset.Position(n.End()).Offset])
}
