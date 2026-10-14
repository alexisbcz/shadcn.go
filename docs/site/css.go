package site

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"shadcn.go/tailwind"
)

// cssPath is the URL of the site stylesheet.
const cssPath = "/assets/docs.css"

// stylesheet builds docs/site/styles.css with package tailwind in the
// background at startup (several seconds: it scans the upstream app and the
// Go sources), caching the result in the user cache directory keyed by the
// sources' sizes and modification times.
type stylesheet struct {
	done chan struct{}
	css  string
	hash string
	err  error
}

// cssSources are the directories whose files go into the cache key.
var cssSources = []string{
	"docs/site/styles.css",
	"ui/apps/v4/app",
	"ui/apps/v4/components",
	"ui/apps/v4/content",
	"ui/apps/v4/registry",
	"ui/apps/v4/styles",
	"ui/apps/v4/examples",
	"ui/apps/v4/lib",
	"ui/apps/v4/hooks",
	"docs",
	"components",
	"examples",
	"cmd/docs",
}

func newStylesheet(root string, logf func(string, ...any)) *stylesheet {
	s := &stylesheet{done: make(chan struct{})}
	go func() {
		defer close(s.done)
		start := time.Now()
		s.css, s.err = buildCSS(root)
		if s.err != nil {
			logf("docs: css: %v", s.err)
			return
		}
		sum := sha256.Sum256([]byte(s.css))
		s.hash = hex.EncodeToString(sum[:8])
		logf("docs: css %d bytes in %s", len(s.css), time.Since(start).Round(time.Millisecond))
	}()
	return s
}

// version waits for the build and returns the stylesheet's hash.
func (s *stylesheet) version() string {
	<-s.done
	return s.hash
}

func buildCSS(root string) (string, error) {
	sum := sha256.New()
	for _, dir := range cssSources {
		_ = filepath.WalkDir(filepath.Join(root, dir), func(p string, d fs.DirEntry, err error) error {
			if err != nil {
				return nil
			}
			if d.IsDir() {
				if n := d.Name(); n == "node_modules" || n == ".next" {
					return filepath.SkipDir
				}
				return nil
			}
			if info, err := d.Info(); err == nil {
				fmt.Fprintf(sum, "%s\x00%d\x00%d\n", p, info.Size(), info.ModTime().UnixNano())
			}
			return nil
		})
	}
	fmt.Fprintf(sum, "tailwind %s\n", tailwind.Version)
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		cacheDir = os.TempDir()
	}
	cacheFile := filepath.Join(cacheDir, "shadcn-docs", hex.EncodeToString(sum.Sum(nil))+".css")
	if b, err := os.ReadFile(cacheFile); err == nil {
		return labColors(string(b)), nil
	}
	res, err := tailwind.Build(tailwind.Options{
		Input:  filepath.Join(root, "docs/site/styles.css"),
		Dir:    filepath.Join(root, "ui/apps/v4"),
		Minify: true,
	})
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(filepath.Dir(cacheFile), 0o755); err == nil {
		tmp := cacheFile + fmt.Sprintf(".%d.tmp", os.Getpid())
		if os.WriteFile(tmp, []byte(res.CSS), 0o644) == nil {
			_ = os.Rename(tmp, cacheFile)
		}
	}
	return labColors(res.CSS), nil
}

func (s *Site) serveCSS(w http.ResponseWriter, r *http.Request) {
	<-s.css.done
	if s.css.err != nil {
		http.Error(w, s.css.err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/css; charset=utf-8")
	if r.URL.Query().Get("v") == s.css.hash {
		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
	} else {
		w.Header().Set("Cache-Control", "no-cache")
	}
	_, _ = w.Write([]byte(s.css.css))
}
