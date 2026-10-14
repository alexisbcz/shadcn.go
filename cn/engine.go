package cn

// This file is a line-by-line port of the merge engine of the cn npm package
// (cn@0.2.2, dist/engine.js, src/engine.ts upstream). Names follow the
// JavaScript so the two can be read side by side. The differences are
// mechanical:
//
//   - JavaScript indexes strings by UTF-16 code unit, Go by byte. The engine
//     only ever compares input characters with ASCII (the tables are checked
//     to be ASCII by the generator), splits tokens on ASCII whitespace and
//     slices at the positions of ASCII characters, so both see the same
//     tokens and make the same decisions.
//   - The engine's mutable scratch state (token arrays, interning maps, claim
//     tables) lives in a per-call scratch value instead of closure variables,
//     so the Go engine is safe for concurrent use. Context and arbitrary
//     group ids are interned per call instead of per engine lifetime; only
//     their equality within one call matters, so results are identical.
//   - The per-token memo and the two-generation whole-string cache are
//     performance devices that never change a result; the Go package keeps a
//     bounded, concurrency-safe whole-string cache instead (see cn.go).

import (
	"regexp"
	"sort"
	"strings"
	"sync"
)

const (
	external = -1
	dead     = -1
)

// fnv is FNV-1a over bytes (code units in JavaScript); it only places
// entries in hash tables, so its value never affects a result.
func fnv(s string, start, end int) uint32 {
	h := uint32(2166136261)
	for p := start; p < end; p++ {
		h = (h ^ uint32(s[p])) * 16777619
	}
	return h
}

// engine holds the tables plus the read-only indexes createEngine derives
// from them.
type engine struct {
	adjRow           []int32
	vgStart          []int32
	postfixLookupSet []bool
	nodeHasLit       []bool
	litMaxLen        int
	litNoArb         bool
	litMask          uint32
	litTable         []int32
	orderSensitive   map[string]bool
	scratch          sync.Pool
}

func newEngine() *engine {
	e := &engine{}
	e.adjRow = make([]int32, groupCount)
	for i := range e.adjRow {
		e.adjRow[i] = -1
	}
	for i, g := range adjGid {
		e.adjRow[g] = int32(i)
	}
	e.vgStart = make([]int32, len(vlistRef)+1)
	for l := range vlistRef {
		e.vgStart[l+1] = e.vgStart[l] + vlistPat[vlistRef[l]+1] - vlistPat[vlistRef[l]]
	}
	e.postfixLookupSet = make([]bool, groupCount)
	for _, g := range postfixLookupGroups {
		e.postfixLookupSet[g] = true
	}
	nodeCount := len(edgeStart) - 1
	e.nodeHasLit = make([]bool, nodeCount)
	e.litNoArb = true
	for i := range litAnchor {
		e.nodeHasLit[litAnchor[i]] = true
		l := int(poolOffsets[litPool[i]*2+1])
		if l > e.litMaxLen {
			e.litMaxLen = l
		}
		c0 := poolText[poolOffsets[litPool[i]*2]]
		if c0 == '[' || c0 == '(' {
			e.litNoArb = false
		}
	}
	size := 1
	for size < len(litAnchor)*2 {
		size <<= 1
	}
	e.litMask = uint32(size - 1)
	e.litTable = make([]int32, size)
	for i := range e.litTable {
		e.litTable[i] = -1
	}
	for i := range litAnchor {
		off := int(poolOffsets[litPool[i]*2])
		idx := (fnv(poolText, off, off+int(poolOffsets[litPool[i]*2+1])) ^ uint32(litAnchor[i])*2654435761) & e.litMask
		for e.litTable[idx] != -1 {
			idx = (idx + 1) & e.litMask
		}
		e.litTable[idx] = int32(i)
	}
	e.orderSensitive = map[string]bool{}
	for _, m := range strings.Split(orderSensitiveModifiers, " ") {
		e.orderSensitive[m] = true
	}
	e.scratch.New = func() any { return newScratch() }
	return e
}

func (e *engine) litProbe(anchor int32, input string, s, end int) int32 {
	idx := (fnv(input, s, end) ^ uint32(anchor)*2654435761) & e.litMask
	n := end - s
	for {
		entry := e.litTable[idx]
		if entry == -1 {
			return -1
		}
		if litAnchor[entry] == anchor && int(poolOffsets[litPool[entry]*2+1]) == n {
			off := int(poolOffsets[litPool[entry]*2])
			if poolText[off:off+n] == input[s:end] {
				return litGroup[entry]
			}
		}
		idx = (idx + 1) & e.litMask
	}
}

// --- validators ----------------------------------------------------------------

// JavaScript's "." matches anything but line terminators (\n, \r, U+2028,
// U+2029); Go's matches anything but \n.
const jsDot = `[^\n\r\x{2028}\x{2029}]`

var (
	lengthUnitRegex    = regexp.MustCompile(`\d+(%|px|r?em|[sdl]?v([hwib]|min|max)|pt|pc|in|cm|mm|cap|ch|ex|r?lh|cq(w|h|i|b|min|max))|\b(calc|min|max|clamp)\(` + jsDot + `+\)|^0$`)
	colorFunctionRegex = regexp.MustCompile(`^(rgba?|hsla?|hwb|(ok)?(lab|lch)|color-mix)\(` + jsDot + `+\)$`)
	shadowRegex        = regexp.MustCompile(`^(inset_)?-?((\d+)?\.?(\d+)[a-z]+|0)_-?((\d+)?\.?(\d+)[a-z]+|0)`)
	imageRegex         = regexp.MustCompile(`^(url|image|image-set|cross-fade|element|(repeating-)?(linear|radial|conic)-gradient)\(` + jsDot + `+\)$`)
	fractionRegex      = regexp.MustCompile(`^\d+(?:\.\d+)?/\d+(?:\.\d+)?$`)
	tshirtRegex        = regexp.MustCompile(`^(\d+(\.\d+)?)?(xs|sm|md|lg|xl)$`)
)

func isWordCode(c byte) bool {
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' || c == '_'
}

// arb is the result of analyzeArb: the shape of an arbitrary value tail.
type arb struct {
	kind       int // 0 none, 1 [..], 2 (..)
	labelS     int
	labelE     int
	valS, valE int
}

func analyzeArb(a *arb, input string, s, e int) {
	a.kind = 0
	a.labelS = -1
	if e-s < 3 {
		return
	}
	c0 := input[s]
	cl := input[e-1]
	if c0 == '[' && cl == ']' {
		a.kind = 1
	} else if c0 == '(' && cl == ')' {
		a.kind = 2
	} else {
		return
	}
	a.valS = s + 1
	a.valE = e - 1
	p := s + 1
	if isWordCode(input[p]) {
		p++
		for p < e-1 {
			c := input[p]
			if !isWordCode(c) && c != '-' {
				break
			}
			p++
		}
		if p < e-2 && input[p] == ':' {
			a.labelS = s + 1
			a.labelE = p
			a.valS = p + 1
		}
	}
}

func spanEq(input string, s, e int, str string) bool {
	return e-s == len(str) && input[s:e] == str
}

func spanIsNamedContainerQuery(input string, s, e int) bool {
	if e-s < 11 || !spanEq(input, s, s+10, "@container") {
		return false
	}
	if input[s+10] == '/' {
		return e-s >= 12
	}
	c11 := input[s+11]
	return c11 == 's' && e-s >= 17 && spanEq(input, s+10, s+16, "-size/") ||
		c11 == 'n' && e-s >= 19 && spanEq(input, s+10, s+18, "-normal/")
}

var vKind = [...]int{1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2, 2, 2}

var vLabels = func() [][]string {
	var out [][]string
	for _, s := range strings.Split("length|number|number weight|family-name|position percentage|length size bg-size|image url|shadow|length|family-name|position percentage|length size bg-size|image url|shadow|number weight", "|") {
		out = append(out, strings.Split(s, " "))
	}
	return out
}()

var vFall = [...]int{2, 3, 1, 0, 0, 0, 4, 5, 0, 0, 0, 0, 0, 1, 1}

func runValidator(a *arb, op int32, input string, s, e int) bool {
	if op >= 10 {
		// op >= 25 are custom validators, which the default tables do not use.
		i := op - 10
		if a.kind != vKind[i] {
			return false
		}
		if a.labelS >= 0 {
			for _, l := range vLabels[i] {
				if spanEq(input, a.labelS, a.labelE, l) {
					return true
				}
			}
			return false
		}
		switch vFall[i] {
		case 0:
			return false
		case 1:
			return true
		case 2:
			v := input[a.valS:a.valE]
			return lengthUnitRegex.MatchString(v) && !colorFunctionRegex.MatchString(v)
		case 3:
			return isNumStr(input[a.valS:a.valE])
		case 4:
			return imageRegex.MatchString(input[a.valS:a.valE])
		default:
			return shadowRegex.MatchString(input[a.valS:a.valE])
		}
	}
	switch op {
	case 0:
		return true
	case 1:
		return a.kind == 0
	case 2:
		return a.kind == 1
	case 3:
		return a.kind == 2
	case 4:
		return fractionRegex.MatchString(input[s:e])
	case 5:
		return isNumStr(input[s:e])
	case 6:
		return isIntStr(input[s:e])
	case 7:
		return e > s && input[e-1] == '%' && isNumStr(input[s:e-1])
	case 8:
		return tshirtRegex.MatchString(input[s:e])
	default:
		return spanIsNamedContainerQuery(input, s, e)
	}
}

// --- per-call state -------------------------------------------------------------

type internKey struct {
	k   string
	imp bool
}

type claimKey struct{ ctx, gid int32 }

type scratch struct {
	tokStart, tokEnd, tokGid, tokCtx []int32
	tokFlags                         []uint8
	keep                             []bool
	ckptNode, ckptTail               []int32
	ctxByRaw                         map[internKey]int32
	ctxByCanon                       map[string]int32
	nextCtxID                        int32
	dynByName                        map[string]int32
	nextDynID                        int32
	claim0                           []int32
	claims                           map[claimKey]struct{}
	epoch                            int32
	arb                              arb
	mods, segment, result            []string
}

func newScratch() *scratch {
	return &scratch{
		ckptNode:   make([]int32, 64),
		ckptTail:   make([]int32, 64),
		ctxByRaw:   map[internKey]int32{},
		ctxByCanon: map[string]int32{},
		dynByName:  map[string]int32{},
		claim0:     make([]int32, groupCount),
		claims:     map[claimKey]struct{}{},
	}
}

func (sc *scratch) reset() {
	sc.tokStart = sc.tokStart[:0]
	sc.tokEnd = sc.tokEnd[:0]
	sc.tokGid = sc.tokGid[:0]
	sc.tokCtx = sc.tokCtx[:0]
	sc.tokFlags = sc.tokFlags[:0]
	clear(sc.ctxByRaw)
	clear(sc.ctxByCanon)
	sc.nextCtxID = 2
	clear(sc.dynByName)
	sc.nextDynID = groupCount
	clear(sc.claims)
	sc.epoch++
	if sc.epoch == 0 {
		clear(sc.claim0)
		sc.epoch = 1
	}
}

// claimTest marks (ctx, gid) as claimed and reports whether it already was.
func (sc *scratch) claimTest(ctx, gid int32) bool {
	if ctx == 0 && gid < groupCount {
		if sc.claim0[gid] == sc.epoch {
			return true
		}
		sc.claim0[gid] = sc.epoch
		return false
	}
	k := claimKey{ctx, gid}
	if _, ok := sc.claims[k]; ok {
		return true
	}
	sc.claims[k] = struct{}{}
	return false
}

func (sc *scratch) canonicalizeContext(raw string, important bool, orderSensitive map[string]bool) int32 {
	mods := sc.mods[:0]
	dB, dP, start := 0, 0, 0
	for i := 0; i < len(raw); i++ {
		c := raw[i]
		if dB == 0 && dP == 0 && c == ':' {
			mods = append(mods, raw[start:i])
			start = i + 1
		} else if c == '[' {
			dB++
		} else if c == ']' {
			dB--
		} else if c == '(' {
			dP++
		} else if c == ')' {
			dP--
		}
	}
	mods = append(mods, raw[start:])
	sc.mods = mods
	canonical := mods[0]
	if len(mods) > 1 {
		// The canonical string is only used as an identity key, so any
		// consistent order of the sortable segments gives the same classes
		// as JavaScript's code-unit sort.
		result := sc.result[:0]
		segment := sc.segment[:0]
		for _, mod := range mods {
			if len(mod) > 0 && mod[0] == '[' || orderSensitive[mod] {
				if len(segment) > 0 {
					sort.Strings(segment)
					result = append(result, segment...)
					segment = segment[:0]
				}
				result = append(result, mod)
			} else {
				segment = append(segment, mod)
			}
		}
		if len(segment) > 0 {
			sort.Strings(segment)
			result = append(result, segment...)
		}
		sc.result, sc.segment = result, segment
		canonical = strings.Join(result, ":")
	}
	key := canonical
	if important {
		key += " !"
	}
	id, ok := sc.ctxByCanon[key]
	if !ok {
		id = sc.nextCtxID
		sc.nextCtxID++
		sc.ctxByCanon[key] = id
	}
	return id
}

// --- merge ----------------------------------------------------------------------

func (e *engine) resolveAt(sc *scratch, input string, bs, endPos int, nodeAt int32, ckptAt int) int32 {
	if endPos-bs >= 2 && input[bs] == '[' && input[endPos-1] == ']' {
		colon := -1
		for p := bs + 1; p < endPos-1; p++ {
			if input[p] == ':' {
				colon = p
				break
			}
		}
		if colon == -1 || colon == bs+1 {
			return external
		}
		name := input[bs+1 : colon]
		id, ok := sc.dynByName[name]
		if !ok {
			id = sc.nextDynID
			sc.nextDynID++
			sc.dynByName[name] = id
		}
		return id
	}
	if nodeAt >= 0 && nodeGroup[nodeAt] >= 0 {
		return nodeGroup[nodeAt]
	}
	for k := ckptAt - 1; k >= 0; k-- {
		tailStart := int(sc.ckptTail[k])
		if tailStart > endPos {
			continue
		}
		nodeID := sc.ckptNode[k]
		tlen := endPos - tailStart
		if e.nodeHasLit[nodeID] && tlen > 0 && tlen <= e.litMaxLen {
			c0 := input[tailStart]
			if !e.litNoArb || c0 != '[' && c0 != '(' {
				if g := e.litProbe(nodeID, input, tailStart, endPos); g >= 0 {
					return g
				}
			}
		}
		vl := nodeVlist[nodeID]
		if vl < 0 {
			continue
		}
		pat := vlistRef[vl]
		vs := vlistPat[pat]
		ve := vlistPat[pat+1]
		if vs == ve {
			continue
		}
		analyzeArb(&sc.arb, input, tailStart, endPos)
		g0 := e.vgStart[vl] - vs
		for v := vs; v < ve; v++ {
			if runValidator(&sc.arb, vlistOps[v], input, tailStart, endPos) {
				return vlistGroup[g0+v]
			}
		}
	}
	return external
}

func (e *engine) pushCkpt(sc *scratch, top int, node int32, tail int) int {
	if top == len(sc.ckptNode) {
		sc.ckptNode = append(sc.ckptNode, make([]int32, len(sc.ckptNode))...)
		sc.ckptTail = append(sc.ckptTail, make([]int32, len(sc.ckptTail))...)
	}
	sc.ckptNode[top] = node
	sc.ckptTail[top] = int32(tail)
	return top + 1
}

// mergeClassList is mergeClassList (the uncached merge) of the engine.
func (e *engine) mergeClassList(input string) string {
	sc := e.scratch.Get().(*scratch)
	defer e.scratch.Put(sc)
	sc.reset()

	n := len(input)
	sawNonSpaceWS := false
	i := 0
	for i < n {
		c := input[i]
		if c == ' ' || c >= 9 && c <= 13 {
			if c != ' ' {
				sawNonSpaceWS = true
			}
			i++
			continue
		}
		ts := i
		for i < n {
			c = input[i]
			if c == ' ' {
				break
			}
			if c >= 9 && c <= 13 {
				sawNonSpaceWS = true
				break
			}
			i++
		}
		te := i
		sc.tokStart = append(sc.tokStart, int32(ts))
		sc.tokEnd = append(sc.tokEnd, int32(te))
		gid, ctxID, flags := e.classify(sc, input, ts, te)
		sc.tokGid = append(sc.tokGid, gid)
		sc.tokCtx = append(sc.tokCtx, ctxID)
		sc.tokFlags = append(sc.tokFlags, flags)
	}

	tokenCount := len(sc.tokStart)
	if tokenCount == 0 {
		return ""
	}
	if tokenCount == 1 {
		return input[sc.tokStart[0]:sc.tokEnd[0]]
	}
	if cap(sc.keep) < tokenCount {
		sc.keep = make([]bool, tokenCount)
	}
	keep := sc.keep[:tokenCount]
	didDrop := false
	for t := tokenCount - 1; t >= 0; t-- {
		gid := sc.tokGid[t]
		if gid == external {
			keep[t] = true
			continue
		}
		ctxID := sc.tokCtx[t]
		if sc.claimTest(ctxID, gid) {
			keep[t] = false
			didDrop = true
			continue
		}
		keep[t] = true
		if gid < groupCount {
			if r := e.adjRow[gid]; r >= 0 {
				for k := adjStart[r]; k < adjStart[r+1]; k++ {
					sc.claimTest(ctxID, adjTgt[k])
				}
			}
			if sc.tokFlags[t]&1 != 0 {
				for k := range patGid {
					if patGid[k] == gid {
						sc.claimTest(ctxID, patTgt[k])
					}
				}
			}
		}
	}
	if !didDrop && !sawNonSpaceWS {
		// Single spaces between tokens and none around them: the input is
		// already the output.
		total := 0
		for t := 0; t < tokenCount; t++ {
			total += int(sc.tokEnd[t] - sc.tokStart[t])
		}
		if n == total+tokenCount-1 {
			return input
		}
	}
	var out strings.Builder
	out.Grow(n)
	for t := 0; t < tokenCount; t++ {
		if !keep[t] {
			continue
		}
		if out.Len() > 0 {
			out.WriteByte(' ')
		}
		out.WriteString(input[sc.tokStart[t]:sc.tokEnd[t]])
	}
	return out.String()
}

// classify resolves one token to its conflict group (external when it takes
// part in no conflict), its variant context and its flags (1 = has a
// postfix modifier). It is the body of the tokenizer loop of
// mergeClassList, after the memo lookup.
func (e *engine) classify(sc *scratch, input string, ts, te int) (gid, ctxID int32, flags uint8) {
	pts := ts
	depthB, depthP := 0, 0
	lastColon, lastSlash := -1, -1
	for p := pts; p < te; p++ {
		pc := input[p]
		if depthB == 0 && depthP == 0 {
			if pc == ':' {
				lastColon = p
				continue
			}
			if pc == '/' {
				lastSlash = p
				continue
			}
		}
		switch pc {
		case '[':
			depthB++
		case ']':
			depthB--
		case '(':
			depthP++
		case ')':
			depthP--
		}
	}
	modStart := pts
	if lastColon >= pts {
		modStart = lastColon + 1
	}
	bs := modStart
	be := te
	important := false
	prefixShift := 0
	if be > bs && input[be-1] == '!' {
		important = true
		be--
	} else if be > bs && input[bs] == '!' {
		important = true
		bs++
		prefixShift = 1
	}
	postfixEnd := -1
	if lastSlash > modStart {
		postfixEnd = lastSlash + prefixShift
		if postfixEnd >= be {
			postfixEnd = -1
		}
	}
	feedStart := bs
	if be-bs > 1 && input[bs] == '-' {
		feedStart = bs + 1
	}
	node := int32(0)
	lp, le := 0, 0
	pending := int32(-1)
	ckptTop := 0
	if nodeVlist[0] >= 0 || e.nodeHasLit[0] {
		ckptTop = e.pushCkpt(sc, 0, 0, feedStart)
	}
	slashNode := int32(dead)
	slashCkpt := 0
	for p := feedStart; p < be; p++ {
		if p == postfixEnd {
			if lp < le {
				slashNode = dead
			} else {
				slashNode = node
			}
			slashCkpt = ckptTop
		}
		if node == dead {
			continue
		}
		cc := input[p]
		arrived := int32(-1)
		if lp < le {
			if labelText[lp] == cc {
				lp++
				if lp == le {
					node = pending
					arrived = node
				}
			} else {
				node = dead
			}
		} else {
			es := edgeStart[node]
			ee := edgeStart[node+1]
			next := int32(dead)
			for ed := es; ed < ee; ed++ {
				ls := labelStart[ed]
				if labelText[ls] == cc {
					if labelStart[ed+1]-ls == 1 {
						next = edgeTarget[ed]
						arrived = next
					} else {
						lp = int(ls) + 1
						le = int(labelStart[ed+1])
						pending = edgeTarget[ed]
						next = node
					}
					break
				}
			}
			node = next
		}
		if arrived >= 0 && (nodeVlist[arrived] >= 0 || e.nodeHasLit[arrived]) && p+1 < be && input[p+1] == '-' {
			ckptTop = e.pushCkpt(sc, ckptTop, arrived, p+2)
		}
	}
	if postfixEnd == be {
		if lp < le {
			slashNode = dead
		} else {
			slashNode = node
		}
		slashCkpt = ckptTop
	}
	endNode := node
	if lp < le {
		endNode = dead
	}
	hasPostfix := false
	if postfixEnd >= 0 {
		hasPostfix = true
		gid = e.resolveAt(sc, input, bs, postfixEnd, slashNode, slashCkpt)
		if gid != external && gid < groupCount && e.postfixLookupSet[gid] {
			gidFull := e.resolveAt(sc, input, bs, be, endNode, ckptTop)
			if gidFull != external && gidFull != gid {
				gid = gidFull
				hasPostfix = false
			}
		} else if gid == external {
			gid = e.resolveAt(sc, input, bs, be, endNode, ckptTop)
			hasPostfix = false
		}
	} else {
		gid = e.resolveAt(sc, input, bs, be, endNode, ckptTop)
	}
	if gid == external {
		return external, 0, 0
	}
	if hasPostfix {
		flags = 1
	}
	if pts >= lastColon {
		if important {
			ctxID = 1
		}
	} else {
		k := internKey{input[pts:lastColon], important}
		id, ok := sc.ctxByRaw[k]
		if !ok {
			id = sc.canonicalizeContext(k.k, important, e.orderSensitive)
			sc.ctxByRaw[k] = id
		}
		ctxID = id
	}
	return gid, ctxID, flags
}
