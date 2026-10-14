// Command docs-stub is a throwaway reverse proxy used only to test the docs
// parity harness (parity/docs/run.mjs). It serves the reference docs site
// unchanged, or with a controlled change injected into chosen pages, so the
// harness can be shown to detect and locate each kind of difference.
//
//	go run ./cmd/docs-stub -upstream http://localhost:4000 -addr 127.0.0.1:4100
//	go run ./cmd/docs-stub -paths /docs/components/radix/button \
//	    -replace 'class="scroll-m-24 text-3xl=>class="scroll-m-24 text-2xl'
//	go run ./cmd/docs-stub -paths /docs/installation -css 'h1{padding-left:1px}'
//
// -replace OLD=>NEW replaces the first occurrence of OLD in the HTML of the
// selected pages (repeatable, applied in order; the proxy exits with an
// error at the first request where OLD is absent, so a typo cannot pass as
// "no difference"). -css injects a <style> element before </head>. -paths
// restricts both to these request paths (comma-separated; default: every
// HTML page). Everything else is proxied byte for byte.
package main

import (
	"bytes"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type replacement struct{ old, new string }

type replaceFlags []replacement

func (r *replaceFlags) String() string { return fmt.Sprint(*r) }

func (r *replaceFlags) Set(v string) error {
	old, new, ok := strings.Cut(v, "=>")
	if !ok || old == "" {
		return fmt.Errorf("want OLD=>NEW, got %q", v)
	}
	*r = append(*r, replacement{old, new})
	return nil
}

func main() {
	upstream := flag.String("upstream", "http://localhost:4000", "reference server")
	addr := flag.String("addr", "127.0.0.1:4100", "listen address (port 0 picks a free one)")
	paths := flag.String("paths", "", "comma-separated request paths to modify (default: all HTML pages)")
	css := flag.String("css", "", "CSS injected in a <style> before </head>")
	var repl replaceFlags
	flag.Var(&repl, "replace", "OLD=>NEW, first occurrence in the HTML (repeatable)")
	flag.Parse()

	target, err := url.Parse(*upstream)
	if err != nil {
		log.Fatal(err)
	}
	only := map[string]bool{}
	for _, p := range strings.Split(*paths, ",") {
		if p = strings.TrimSpace(p); p != "" {
			only[p] = true
		}
	}
	modify := func(r *http.Request) bool {
		return (len(repl) > 0 || *css != "") && (len(only) == 0 || only[r.URL.Path])
	}

	proxy := httputil.NewSingleHostReverseProxy(target)
	director := proxy.Director
	proxy.Director = func(r *http.Request) {
		director(r)
		r.Host = target.Host
		if modify(r) {
			// The body is rewritten: ask for it uncompressed.
			r.Header.Del("Accept-Encoding")
		}
	}
	proxy.ModifyResponse = func(res *http.Response) error {
		req := res.Request
		if !modify(req) || !strings.HasPrefix(res.Header.Get("Content-Type"), "text/html") {
			return nil
		}
		// Only full HTML documents, not RSC payloads of client navigations.
		if req.Header.Get("RSC") != "" {
			return nil
		}
		var body []byte
		var err error
		if res.Header.Get("Content-Encoding") == "gzip" {
			zr, zerr := gzip.NewReader(res.Body)
			if zerr != nil {
				return zerr
			}
			body, err = io.ReadAll(zr)
			res.Header.Del("Content-Encoding")
		} else {
			body, err = io.ReadAll(res.Body)
		}
		res.Body.Close()
		if err != nil {
			return err
		}
		for _, rp := range repl {
			i := bytes.Index(body, []byte(rp.old))
			if i < 0 {
				log.Printf("docs-stub: %q not found in %s", rp.old, req.URL.Path)
				os.Exit(3)
			}
			body = append(body[:i:i], append([]byte(rp.new), body[i+len(rp.old):]...)...)
		}
		if *css != "" {
			body = bytes.Replace(body, []byte("</head>"), []byte("<style>"+*css+"</style></head>"), 1)
		}
		res.Body = io.NopCloser(bytes.NewReader(body))
		res.ContentLength = int64(len(body))
		res.Header.Set("Content-Length", strconv.Itoa(len(body)))
		res.Header.Del("ETag")
		log.Printf("docs-stub: modified %s", req.URL.Path)
		return nil
	}

	ln, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("docs-stub listening on http://%s (upstream %s)\n", ln.Addr(), target)
	log.Fatal(http.Serve(ln, proxy))
}
