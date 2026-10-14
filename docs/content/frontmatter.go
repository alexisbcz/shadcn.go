package content

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Frontmatter is a page's YAML frontmatter, read like fumadocs (js-yaml
// load). The docs use the keys below; Values keeps every key in order.
type Frontmatter struct {
	Title       string
	Description string
	// Links holds links.doc and links.api.
	Links *Links
	// Base and Component are set on component pages.
	Base, Component string
	// Date is set on changelog pages (a YAML timestamp, UTC).
	Date     time.Time
	HasDate  bool
	Featured bool
	// Values holds all keys in source order.
	Values []KV
}

// Links are the documentation and API reference links of a component page.
type Links struct{ Doc, API string }

// KV is one frontmatter key. Value is a string, bool, time.Time or []KV
// (a nested mapping).
type KV struct {
	Key   string
	Value any
}

var frontmatterRe = regexp.MustCompile(`(?s)^---\r?\n(.+?)\r?\n---\r?\n?`)

// SplitFrontmatter separates the frontmatter from the content, like
// fumadocs-core/content/md/frontmatter.
func SplitFrontmatter(src string) (matter, content string) {
	m := frontmatterRe.FindStringSubmatchIndex(src)
	if m == nil {
		return "", src
	}
	return src[m[2]:m[3]], src[m[1]:]
}

// ParseFrontmatter parses the YAML subset the docs use: mappings (nested by
// indentation), plain and quoted strings, booleans and dates.
func ParseFrontmatter(yaml string) (*Frontmatter, error) {
	values, err := parseMapping(strings.Split(yaml, "\n"), 0)
	if err != nil {
		return nil, err
	}
	fm := &Frontmatter{Values: values}
	for _, kv := range values {
		switch kv.Key {
		case "title":
			fm.Title, _ = kv.Value.(string)
		case "description":
			fm.Description, _ = kv.Value.(string)
		case "base":
			fm.Base, _ = kv.Value.(string)
		case "component":
			fm.Component, _ = kv.Value.(string)
		case "featured":
			fm.Featured, _ = kv.Value.(bool)
		case "date":
			if t, ok := kv.Value.(time.Time); ok {
				fm.Date, fm.HasDate = t, true
			}
		case "links":
			if m, ok := kv.Value.([]KV); ok {
				fm.Links = &Links{}
				for _, l := range m {
					s, _ := l.Value.(string)
					switch l.Key {
					case "doc":
						fm.Links.Doc = s
					case "api":
						fm.Links.API = s
					}
				}
			}
		}
	}
	return fm, nil
}

func indentOf(s string) int { return len(s) - len(strings.TrimLeft(s, " ")) }

func parseMapping(lines []string, indent int) ([]KV, error) {
	var out []KV
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		if strings.TrimSpace(line) == "" || strings.HasPrefix(strings.TrimSpace(line), "#") {
			continue
		}
		if indentOf(line) != indent {
			return nil, fmt.Errorf("frontmatter: unexpected indentation in %q", line)
		}
		key, rest, ok := strings.Cut(line[indent:], ":")
		if !ok || strings.ContainsAny(key, " \"'") {
			return nil, fmt.Errorf("frontmatter: unsupported line %q", line)
		}
		rest = strings.TrimSpace(rest)
		if rest == "" {
			// Nested mapping.
			j := i + 1
			for j < len(lines) && (strings.TrimSpace(lines[j]) == "" || indentOf(lines[j]) > indent) {
				j++
			}
			if j == i+1 {
				out = append(out, KV{Key: key, Value: nil})
				continue
			}
			child, err := parseMapping(lines[i+1:j], indentOf(lines[i+1]))
			if err != nil {
				return nil, err
			}
			out = append(out, KV{Key: key, Value: child})
			i = j - 1
			continue
		}
		v, err := parseScalar(rest)
		if err != nil {
			return nil, err
		}
		out = append(out, KV{Key: key, Value: v})
	}
	return out, nil
}

var dateRe = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

func parseScalar(s string) (any, error) {
	switch {
	case strings.HasPrefix(s, `"`):
		if !strings.HasSuffix(s, `"`) || len(s) < 2 {
			return nil, fmt.Errorf("frontmatter: unterminated string %s", s)
		}
		return strconv.Unquote(s)
	case strings.HasPrefix(s, "'"):
		if !strings.HasSuffix(s, "'") || len(s) < 2 {
			return nil, fmt.Errorf("frontmatter: unterminated string %s", s)
		}
		return strings.ReplaceAll(s[1:len(s)-1], "''", "'"), nil
	case s == "true" || s == "True" || s == "TRUE":
		return true, nil
	case s == "false" || s == "False" || s == "FALSE":
		return false, nil
	case dateRe.MatchString(s):
		return time.Parse("2006-01-02", s)
	}
	if strings.ContainsAny(s[:1], "[{&*!|>%@`") || strings.Contains(s, " #") || strings.Contains(s, ": ") {
		return nil, fmt.Errorf("frontmatter: unsupported value %q", s)
	}
	return s, nil
}
