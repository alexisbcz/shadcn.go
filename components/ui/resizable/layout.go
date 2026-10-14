package resizable

import (
	"math"
	"strconv"
	"strings"
)

// The server renders the layout the client computes when it mounts, for the
// sizes it can resolve without measuring (percentages). This is a port of
// the layout functions of react-resizable-panels; resizable.js holds the
// complete client version.

type constraint struct {
	collapsedSize float64
	collapsible   bool
	defaultSize   float64
	hasDefault    bool
	minSize       float64
	maxSize       float64
}

// percent parses a size the server can resolve: "N%" or a bare number.
func percent(s string) (float64, bool) {
	s = strings.TrimSpace(s)
	if strings.HasSuffix(s, "%") {
		s = strings.TrimSuffix(s, "%")
	} else if strings.TrimLeft(s, "0123456789.+-eE") != "" {
		return 0, false
	}
	f, err := strconv.ParseFloat(s, 64)
	return f, err == nil
}

// round3 is parseFloat(x.toFixed(3)).
func round3(x float64) float64 {
	v, _ := strconv.ParseFloat(strconv.FormatFloat(x, 'f', 3, 64), 64)
	return v
}

func fuzzyEqual(a, b, tolerance float64) bool { return math.Abs(round3(a)-round3(b)) <= tolerance }

func fuzzyCompare(a, b float64) int {
	if fuzzyEqual(a, b, 0) {
		return 0
	}
	if a > b {
		return 1
	}
	return -1
}

func initialLayout(panels []*PanelBuilder, defaultLayout DefaultLayout) ([]float64, []constraint, bool) {
	if len(panels) == 0 {
		return nil, nil, false
	}
	cons := make([]constraint, len(panels))
	for i, p := range panels {
		c := constraint{maxSize: 100}
		var ok bool
		if p.collapsedSize != "" {
			if c.collapsedSize, ok = percent(p.collapsedSize); !ok {
				return nil, nil, false
			}
			c.collapsedSize = round3(c.collapsedSize)
		}
		if p.defaultSize != "" {
			if c.defaultSize, ok = percent(p.defaultSize); !ok {
				return nil, nil, false
			}
			c.defaultSize = round3(c.defaultSize)
			c.hasDefault = true
		}
		if p.minSize != "" {
			if c.minSize, ok = percent(p.minSize); !ok {
				return nil, nil, false
			}
			c.minSize = round3(c.minSize)
		}
		if p.maxSize != "" {
			if c.maxSize, ok = percent(p.maxSize); !ok {
				return nil, nil, false
			}
			c.maxSize = round3(c.maxSize)
		}
		c.collapsible = p.collapsible
		cons[i] = c
	}
	var layout []float64
	if defaultLayout != nil {
		// Panel ids are generated; a default layout can only name panels
		// given explicit ids, which the client checks. Let it decide.
		return nil, nil, false
	}
	// Default layout from the default sizes.
	layout = make([]float64, len(cons))
	count, total := 0, 0.0
	for i, c := range cons {
		if c.hasDefault {
			count++
			layout[i] = round3(c.defaultSize)
			total += layout[i]
		}
	}
	if rest := len(cons) - count; rest != 0 {
		size := round3((100 - total) / float64(rest))
		for i, c := range cons {
			if !c.hasDefault {
				layout[i] = size
			}
		}
	}
	return validateLayout(layout, cons), cons, true
}

func resizePanel(c constraint, size float64) float64 {
	if fuzzyCompare(size, c.minSize) < 0 {
		if c.collapsible {
			half := (c.collapsedSize + c.minSize) / 2
			if fuzzyCompare(size, half) < 0 {
				size = c.collapsedSize
			} else {
				size = c.minSize
			}
		} else {
			size = c.minSize
		}
	}
	size = math.Min(c.maxSize, size)
	return round3(size)
}

func validateLayout(layout []float64, cons []constraint) []float64 {
	out := append([]float64(nil), layout...)
	sum := 0.0
	for _, v := range out {
		sum += v
	}
	if !fuzzyEqual(sum, 100, 0) && len(out) > 0 {
		for i := range out {
			out[i] = 100 / sum * out[i]
		}
	}
	remaining := 0.0
	for i := range cons {
		v := out[i]
		s := resizePanel(cons[i], v)
		if v != s {
			remaining += v - s
			out[i] = s
		}
	}
	if !fuzzyEqual(remaining, 0, 0) {
		for i := range cons {
			v := out[i]
			s := resizePanel(cons[i], v+remaining)
			if v != s {
				remaining -= s - v
				out[i] = s
				if fuzzyEqual(remaining, 0, 0) {
					break
				}
			}
		}
	}
	return out
}

func equalLayouts(a, b []float64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// adjustLayoutByDelta ports react-resizable-panels' adjustLayoutByDelta for
// a non-keyboard trigger.
func adjustLayoutByDelta(delta float64, initial []float64, cons []constraint, pivot [2]int, prev []float64) []float64 {
	if fuzzyEqual(delta, 0, 0) {
		return initial
	}
	next := append([]float64(nil), initial...)
	first, second := pivot[0], pivot[1]
	deltaApplied := 0.0
	{
		inc := -1
		idx := first
		if delta < 0 {
			inc = 1
			idx = second
		}
		maxAvailable := 0.0
		for {
			maxAvailable += resizePanel(cons[idx], 100) - initial[idx]
			idx += inc
			if idx < 0 || idx >= len(cons) {
				break
			}
		}
		m := math.Min(math.Abs(delta), math.Abs(maxAvailable))
		if delta < 0 {
			delta = -m
		} else {
			delta = m
		}
	}
	{
		idx := second
		if delta < 0 {
			idx = first
		}
		for idx >= 0 && idx < len(cons) {
			remaining := math.Abs(delta) - math.Abs(deltaApplied)
			prevSize := initial[idx]
			unsafe := prevSize - remaining
			safe := resizePanel(cons[idx], unsafe)
			if !fuzzyEqual(prevSize, safe, 0) {
				deltaApplied += prevSize - safe
				next[idx] = safe
				if compareFixed3(deltaApplied, math.Abs(delta)) >= 0 {
					break
				}
			}
			if delta < 0 {
				idx--
			} else {
				idx++
			}
		}
	}
	if equalLayouts(prev, next) {
		return prev
	}
	{
		idx := first
		if delta < 0 {
			idx = second
		}
		unsafe := initial[idx] + deltaApplied
		safe := resizePanel(cons[idx], unsafe)
		next[idx] = safe
		if !fuzzyEqual(safe, unsafe, 0) {
			remaining := unsafe - safe
			j := first
			if delta < 0 {
				j = second
			}
			for j >= 0 && j < len(cons) {
				prevSize := next[j]
				s := resizePanel(cons[j], prevSize+remaining)
				if !fuzzyEqual(prevSize, s, 0) {
					remaining -= s - prevSize
					next[j] = s
				}
				if fuzzyEqual(remaining, 0, 0) {
					break
				}
				if delta > 0 {
					j--
				} else {
					j++
				}
			}
		}
	}
	total := 0.0
	for _, v := range next {
		total += v
	}
	if !fuzzyEqual(total, 100, 0.1) {
		return prev
	}
	return next
}

// compareFixed3 compares a.toFixed(3) and b.toFixed(3) numerically, like
// localeCompare with numeric collation.
func compareFixed3(a, b float64) int {
	x, y := round3(a), round3(b)
	switch {
	case x < y:
		return -1
	case x > y:
		return 1
	}
	return 0
}

// separatorAria computes a handle's aria-value* like the client: the
// handle controls the panel before it.
func separatorAria(kids []any, handle *resolvedHandle, layout []float64, cons []constraint) *handleAria {
	// The handle controls the closest panel before it.
	panelIndex := -1
	var id string
	seen := 0
	for _, k := range kids {
		if p, ok := k.(*resolvedPanel); ok {
			seen++
			id = p.id
			continue
		}
		if k == any(handle) {
			panelIndex = seen - 1
			break
		}
	}
	if panelIndex < 0 || panelIndex+1 >= len(layout) {
		return nil
	}
	now := layout[panelIndex]
	c := cons[panelIndex]
	minTarget := c.minSize
	if c.collapsible {
		minTarget = c.collapsedSize
	}
	pivot := [2]int{panelIndex, panelIndex + 1}
	minLayout := validateLayout(adjustLayoutByDelta(minTarget-now, layout, cons, pivot, layout), cons)
	maxLayout := validateLayout(adjustLayoutByDelta(c.maxSize-now, layout, cons, pivot, layout), cons)
	return &handleAria{controls: id, now: now, min: minLayout[panelIndex], max: maxLayout[panelIndex], hasMin: true, hasMax: true}
}
