package content

import (
	"sort"
	"strconv"
	"strings"
)

// Slugger generates unique heading ids exactly like github-slugger 2.0.0,
// which fumadocs' remarkHeading uses: lowercase, remove the characters in
// its regex, spaces to dashes, and "-1", "-2", ... for repeats.
type Slugger struct{ occurrences map[string]int }

// Slug returns the unique slug for value.
func (s *Slugger) Slug(value string) string {
	if s.occurrences == nil {
		s.occurrences = map[string]int{}
	}
	result := Slug(value)
	original := result
	for {
		if _, ok := s.occurrences[result]; !ok {
			break
		}
		s.occurrences[original]++
		result = original + "-" + strconv.Itoa(s.occurrences[original])
	}
	s.occurrences[result] = 0
	return result
}

// Reset forgets the slugs generated so far.
func (s *Slugger) Reset() { s.occurrences = nil }

// Slug is github-slugger's slug(value) without uniqueness.
func Slug(value string) string {
	var b strings.Builder
	for _, r := range jsLower(value) {
		if removed(r) {
			continue
		}
		if r == ' ' {
			b.WriteByte('-')
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}

func removed(r rune) bool {
	i := sort.Search(len(sluggerRemoved), func(i int) bool { return sluggerRemoved[i][1] >= r })
	return i < len(sluggerRemoved) && sluggerRemoved[i][0] <= r
}

// jsLower is String.prototype.toLowerCase for the text the docs use. Go's
// strings.ToLower differs only for a few special casings (İ, final sigma
// contexts), which are reported by the tests if they ever appear.
func jsLower(s string) string { return strings.ToLower(s) }
