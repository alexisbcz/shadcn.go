package site

import "testing"

func TestFlattenNesting(t *testing.T) {
	for in, want := range map[string]string{
		`.a{color:red;&:hover{color:blue}}`:                         `.a{color:red}.a:hover{color:blue}`,
		`.a{@media (hover:hover){&:hover{color:blue}}}`:             `@media (hover:hover){.a:hover{color:blue}}`,
		`.s{--p:to bottom;@supports (x:y){--p:to bottom in oklab}}`: `.s{--p:to bottom}@supports (x:y){.s{--p:to bottom in oklab}}`,
		`.a,.b{&:is(:where(.g):hover *){x:y}}`:                      `.a:is(:where(.g):hover *),.b:is(:where(.g):hover *){x:y}`,
		`.a .b{.x &{c:d}}`:                                          `.x :is(.a .b){c:d}`,
		`@layer theme,base;@layer u{.\[\&\>\*\]\:x{&>*{a:b}}}`:      `@layer theme,base;@layer u{.\[\&\>\*\]\:x>*{a:b}}`,
		`@keyframes k{0%{a:b}to{a:c}}.a{b:c}`:                       `@keyframes k{0%{a:b}to{a:c}}.a{b:c}`,
		`@property --x{syntax:"*";inherits:false}`:                  `@property --x{syntax:"*";inherits:false}`,
		`.a{&:hover{@media (x){&:focus{c:d}}}}`:                     `@media (x){.a:hover:focus{c:d}}`,
		`.a{b:c;&:hover{d:e}f:g}`:                                   `.a{b:c}.a:hover{d:e}.a{f:g}`,
	} {
		if got := flattenNesting(in); got != want {
			t.Errorf("flattenNesting(%s)\n got %s\nwant %s", in, got, want)
		}
	}
}
