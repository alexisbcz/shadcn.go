// Package runtime serves the client-side behavior of the shadcn/ui
// components: the vanilla JavaScript that replaces Radix UI and the other
// React libraries upstream components rely on (focus management, dismissal,
// positioning, keyboard navigation, presence animations, ...).
//
// Each component package that needs behavior embeds its own module and
// registers it from init:
//
//	//go:embed dialog.js
//	var script string
//
//	func init() { runtime.Module("dialog", script) }
//
// Pages then include the runtime once, and only the behaviors of the
// components linked into the binary are served:
//
//	mux.Handle(runtime.Prefix, runtime.Handler())
//	h.Head(..., runtime.Scripts(r))
//
// Modules are ES modules served from Prefix; they import the shared
// primitives in js/ by relative path ("./core.js", "./presence.js", ...).
package runtime

import (
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"fmt"
	"io/fs"
	"net/http"
	"path"
	"slices"
	"strings"
	"sync"

	h "shadcn.go/libhtml"
	"shadcn.go/libhtml/secure"
)

// Prefix is the URL path the runtime is served under.
const Prefix = "/_shadcn/"

//go:embed js
var primitives embed.FS

var (
	mu      sync.RWMutex
	modules = map[string]string{} // component name -> module source
	etag    string
)

// Module registers a component's behavior module. It is served as
// Prefix+name+".js" and imported by the runtime's index module. Call it from
// the component package's init.
func Module(name, source string) {
	mu.Lock()
	defer mu.Unlock()
	if _, dup := modules[name]; dup {
		panic(fmt.Sprintf("runtime: module %q registered twice", name))
	}
	if _, err := primitives.Open("js/" + name + ".js"); err == nil {
		panic(fmt.Sprintf("runtime: module %q clashes with a primitive", name))
	}
	modules[name] = source
	etag = ""
}

// Names returns the registered component modules, sorted.
func Names() []string {
	mu.RLock()
	defer mu.RUnlock()
	names := make([]string, 0, len(modules))
	for n := range modules {
		names = append(names, n)
	}
	slices.Sort(names)
	return names
}

// index is the entry module: it starts the core and imports every
// registered component module.
func index() string {
	var b strings.Builder
	b.WriteString("import { start } from \"./core.js\"\n")
	for _, n := range Names() {
		fmt.Fprintf(&b, "import %q\n", "./"+n+".js")
	}
	b.WriteString("start()\n")
	return b.String()
}

func version() string {
	mu.Lock()
	defer mu.Unlock()
	if etag != "" {
		return etag
	}
	sum := sha256.New()
	_ = fs.WalkDir(primitives, "js", func(p string, d fs.DirEntry, err error) error {
		if err == nil && !d.IsDir() {
			b, _ := primitives.ReadFile(p)
			sum.Write([]byte(p))
			sum.Write(b)
		}
		return nil
	})
	names := make([]string, 0, len(modules))
	for n := range modules {
		names = append(names, n)
	}
	slices.Sort(names)
	for _, n := range names {
		sum.Write([]byte(n))
		sum.Write([]byte(modules[n]))
	}
	etag = `"` + hex.EncodeToString(sum.Sum(nil))[:16] + `"`
	return etag
}

// Handler serves the runtime under Prefix.
func Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := strings.TrimPrefix(r.URL.Path, Prefix)
		if !strings.HasSuffix(name, ".js") || strings.Contains(name, "..") {
			http.NotFound(w, r)
			return
		}
		var src string
		switch base := strings.TrimSuffix(name, ".js"); {
		case base == "index":
			src = index()
		default:
			mu.RLock()
			s, ok := modules[base]
			mu.RUnlock()
			if ok {
				src = s
			} else if b, err := primitives.ReadFile(path.Join("js", name)); err == nil {
				src = string(b)
			} else {
				http.NotFound(w, r)
				return
			}
		}
		tag := version()
		w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("ETag", tag)
		if r.Header.Get("If-None-Match") == tag {
			w.WriteHeader(http.StatusNotModified)
			return
		}
		_, _ = w.Write([]byte(src))
	})
}

// Scripts returns the <script> tag that loads the runtime, carrying the
// request's CSP nonce (see secure.Headers). Put it in <head>.
func Scripts(r *http.Request) h.Node {
	return secure.Script(r, h.Type("module"), h.Src(Prefix+"index.js"))
}
