package jsqt

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetClosure_Isolated(t *testing.T) {

	tt := []struct {
		give string
		when string
		then string
	}{
		{give: `{"a":[{"b":{}}]}`, when: `(collect (get a) () (get b) (get c))`, then: `[]`},
		// {give: `{"a":[{"b":{"c":3}},{"b":{}}]}`, when: `(collect (get a) () (get b) (get c))`, then: `[3,4]`},
	}
	for _, tc := range tt {
		r := GetClosure(tc.give, tc.when)
		assert.Equal(t, tc.then, r.String(), "TC: %v", tc)
	}
}

func TestGetClosure(t *testing.T) {

	tt := []struct {
		give string
		when string
		then string
	}{
		// Join.
		{
			give: `{"a":{"b":[{"c":"one","d":"one-val"},{"c":"two","d":"two-val"}]}}`,
			when: `(join (get a b) () (get c) (get d))`,
			then: `{"one":"one-val","two":"two-val"}`,
		},
		// Collect.
		{
			give: `{"a":{"b":{"c":[{"d":"one","e":{"f":[{"g":{"h":{"i":{"j":[{"k":{"l":"hi"}}]}}}}]}},{"d":"two","e":{"f":[{"g":{"h":{"i":{"j":[]}}}}]}}]}}}`,
			when: `(collect (get a b c) () (obj d (get d) e (flatten (get e f g h i j k l))))`,
			then: `[{"d":"one","e":["hi"]},{"d":"two","e":[]}]`,
		},
		{give: `{"a":[{"b":{"c":3}},{"b":{}}]}`, when: `(collect (get a) () (get b c))`, then: `[3]`},
		{give: `{"a":[{"b":3},{"b":4}]}`, when: `(collect (get a) () (get b))`, then: `[3,4]`},
		{give: `[{"b":3},{"c":4},{"b":5}]`, when: `(collect (.) () (get b))`, then: `[3,5]`},
		// Array.
		{give: `{"a":3,"b":4}`, when: `(arr (get a) (get b) (get a))`, then: `[3,4,3]`},
		// Object.
		{give: `{"a":3,"b":4}`, when: `(obj x (get a) y (get b))`, then: `{"x":3,"y":4}`},
		{give: `{"a":3,"b":4}`, when: `(obj x (get a) y (get b))`, then: `{"x":3,"y":4}`},
		// Get.
		{give: `{"a":[{"b":3},{"b":4}]}`, when: `(get a b)`, then: `[3,4]`},
		{give: `[{"a":3},{"a":4}]`, when: `(get a)`, then: `[3,4]`},
		{give: `[{"a":3},{"a":4}]`, when: `(get 0 a)`, then: `3`},
		{give: `[2,3,4]`, when: `(get 3)`, then: ``},
		{give: `[2,3,4]`, when: `(get 2)`, then: `4`},
		{give: `[2,3,4]`, when: `(get 1)`, then: `3`},
		{give: `[2,3,4]`, when: `(get 0)`, then: `2`},
		{give: `{"aa":{"bb":{"cc":3}}}`, when: `(get aa bb cc)`, then: `3`},
		{give: `{"a b":{"c d":3}}`, when: `(get "a b" "c d")`, then: `3`},
		{give: `{"a b":3}`, when: `(get "a b")`, then: `3`},
		{give: `{"a":{"b":3}}`, when: `(get a b)`, then: `3`},
		{give: `{"a":3}`, when: `(get a)`, then: `3`},
		// Root.
		{give: `3`, when: `(root)`, then: `3`},
		{give: ``, when: `(root)`, then: ``},
		// Raw.
		{give: ``, when: `(raw {})`, then: `{}`},
		{give: ``, when: `(raw [])`, then: `[]`},
		{give: ``, when: `(raw null)`, then: `null`},
		{give: ``, when: `(raw true)`, then: `true`},
		{give: ``, when: `(raw false)`, then: `false`},
		{give: ``, when: `(raw 3e2)`, then: `3e2`},
		{give: ``, when: `(raw 3)`, then: `3`},
		{give: ``, when: `(raw "a")`, then: `"a"`},
		{give: ``, when: `(raw -3)`, then: `-3`},
		{give: ``, when: `(raw 1.2)`, then: `1.2`},
	}
	for _, tc := range tt {
		r := GetClosure(tc.give, tc.when)
		assert.Equal(t, tc.then, r.String(), "TC: %v", tc)
	}
}

func TestGet2(t *testing.T) {

	tt := []struct {
		give string
		when string
		then string
	}{
		{give: `[{"a":3},{"a":4}]`, when: `.0.a`, then: `3`},
		{give: `[2,3,4]`, when: `.3`, then: ``},
		{give: `[2,3,4]`, when: `.2`, then: `4`},
		{give: `[2,3,4]`, when: `.1`, then: `3`},
		{give: `[2,3,4]`, when: `.0`, then: `2`},
		{give: `{"aa":{"bb":{"cc":3}}}`, when: `.aa.bb.cc`, then: `3`},
		{give: `{"a":{"b":3}}`, when: `.a.b`, then: `3`},
		{give: `{"a b":3}`, when: `."a b"`, then: `3`},
		{give: `{"a":3}`, when: `.a`, then: `3`},
		{give: `3`, when: `.`, then: `3`},
		{give: ``, when: `.`, then: ``},
		{give: ``, when: `{}`, then: `{}`},
		{give: ``, when: `[]`, then: `[]`},
		{give: ``, when: `null`, then: `null`},
		{give: ``, when: `true`, then: `true`},
		{give: ``, when: `false`, then: `false`},
		{give: ``, when: `3e2`, then: `3e2`},
		{give: ``, when: `3`, then: `3`},
		{give: ``, when: `-3`, then: `-3`},
		{give: ``, when: `1.2`, then: `1.2`},
	}
	for _, tc := range tt {
		r := Get2(tc.give, tc.when)
		assert.Equal(t, tc.then, r.String(), "TC: %v", tc)
	}
}

func TestGet(t *testing.T) {

	tt := []struct {
		give string
		when string
		then string
	}{
		// Filter.
		{
			give: `[{"a":[1,5],"b":3},{"a":[3,6],"b":7}]`,
			when: `.(collect).|a b>!4|.(collect).|. .>!4|`,
			then: `[[6]]`,
		},
		{
			give: `[3,4,5,6]`,
			when: `.(collect).|. .>@root.1|`,
			then: `[5,6]`,
		},
		{
			give: `[3,4,5,6]`,
			when: `.(collect).|. .>!4|`,
			then: `[5,6]`,
		},
		{
			give: `{"a":[{"b":[{"c":[{"d":3},{"d":4}]},{"c":[{"d":5},{"d":6}]}]}]}`,
			when: `a.(collect).b.(collect).c.(omitempty).(collect).|d d>!4|`,
			then: `[[[5,6]]]`,
		},
		{
			give: `{"a":[{"b":[{"c":[{"d":3},{"d":4}]},{"c":[{"d":5},{"d":6}]}]}]}`,
			when: `a.(collect).b.(collect).c.(collect).|d d>!4|`,
			then: `[[[],[5,6]]]`,
		},
		{
			give: `{"a":[{"b":3,"c":4},{"b":5,"c":6}]}`,
			when: `a.(collect).|b c>!4|`,
			then: `[5]`,
		},
		{
			give: `{"a":[{"b":3,"c":4},{"b":5,"c":6}]}`,
			when: `a.(collect).|{b} c>!4|`,
			then: `[{"b":5}]`,
		},
		// Merge function.
		{give: `[{"a":3},{"b":4},{"c":5}]`, when: `.(merge)`, then: `{"a":3,"b":4,"c":5}`},
		// Default function.
		{give: `[{"b":3},{"c":4},{"b":5}]`, when: `.(collect).(default 0).b`, then: `[3,0,5]`},
		// OmitEmpty Function.
		{give: `{"a":[[3],[]]}`, when: `.a.(size).(omitempty)`, then: `1`},
		{give: `{"a":[[3],[]]}`, when: `.a.(collect).(omitempty)`, then: `[[3]]`},
		{give: `{"a":[{"d":3},{"e":4}]}`, when: `.a.(size).(omitempty).{d}`, then: `1`},
		{give: `{"a":[{"d":3},{"e":4}]}`, when: `.a.(collect).(omitempty).{d}`, then: `[{"d":3}]`},
		// Size function.
		{give: `{"a":[{"b":[{"c":[{"d":3}]}]}]}`, when: `.a.(flatten).b.(collect).c.(size).d`, then: `[1]`},
		{give: `{"a":[{"d":3},{"e":4}]}`, when: `.a.(size).d`, then: `1`},
		{give: `{"a":[3,4]}`, when: `.a.(size)`, then: `2`},
		{give: `{"a":[3]}`, when: `.a.(size)`, then: `1`},
		{give: `{"a":[]}`, when: `.a.(size)`, then: `0`},
		// Flatten function.
		{give: `{"a":{"b":{"c":[{"d":"one","e":{"f":[{"g":{"h":{"i":{"j":[{"k":{"l":"hi"}}]}}}}]}},{"d":"two","e":{"f":[{"g":{"h":{"i":{"j":[]}}}}]}}]}}}`, when: `.a.b.c.(collect).{d:.d,e:.e.f.(flatten).g.h.i.j.(collect).k.l}`, then: `[{"d":"one","e":["hi"]},{"d":"two","e":[]}]`},
		{give: `{"a":[{"b":[{"c":[{"d":3}]}]}]}`, when: `.a.(flatten).b.(flatten).c.(flatten).d`, then: `3`},
		{give: `{"a":[{"b":[{"c":[{"d":3}]}]}]}`, when: `.a.(flatten).b.(flatten).c.(collect).d`, then: `[3]`},
		{give: `{"a":[{"b":[{"c":[{"d":3}]}]}]}`, when: `.a.(flatten).b.(collect).c.(collect).d`, then: `[[3]]`},
		{give: `{"a":[{"b":[{"c":3}]}]}`, when: `.a.(flatten).b.(collect).c`, then: `[3]`},
		// First context (dot) is optional.
		{give: `{"a":3,"b":4}`, when: `{x:a}`, then: `{"x":3}`},
		{give: `{"a":3,"b":4}`, when: `.{x:a}`, then: `{"x":3}`},
		{give: `{"a":3,"b":4}`, when: `.{x:.a}`, then: `{"x":3}`},
		{give: `{"a":3,"b":4}`, when: `{x:.a}`, then: `{"x":3}`},
		{give: `{"a":3,"b":4}`, when: `.{.a}`, then: `{"a":3}`},
		{give: `{"a":3,"b":4}`, when: `{.a}`, then: `{"a":3}`},
		{give: `{"a":3,"b":4}`, when: `.{a}`, then: `{"a":3}`},
		// Root node access.
		{give: `{"a":{"b":{"c":{"d":3}}},"e":4}`, when: `{a.b.{c.d,x:@root.e}}`, then: `{"b":{"d":3,"x":4}}`},
		// Array generation.
		{give: `3`, when: `[.,.]`, then: `[3,3]`},
		{give: `3`, when: `[.,@root]`, then: `[3,3]`},
		{give: `3`, when: `[{x:.}]`, then: `[{"x":3}]`},
		{give: `{"a":[3,4]}`, when: `a.(collect).[.]`, then: `[[3],[4]]`},
		{give: `{"a":{"b":[3,4]}}`, when: `[a.b.1,a.b.0]`, then: `[4,3]`},
		// General.
		{give: `{"a":{"a b c":3}}`, when: `{a."a b c","x y z":a."a b c"}`, then: `{"a b c":3,"x y z":3}`},
		{give: `{"a":{"b":{"c":3}}}`, when: `a.{b.{c}}`, then: `{"b":{"c":3}}`},
		{give: `{"a":{"b":{"c":[{"d":"one","e":{"f":[{"g":{"h":{"i":{"j":[{"k":{"l":"hi"}}]}}}}]}},{"d":"two","e":{"f":[{"g":{"h":{"i":{"j":[]}}}}]}}]}}}`, when: `a.b.c.(collect).{d:d,e:e.f.(collect).g.h.i.j.(collect).k.l}`, then: `[{"d":"one","e":[["hi"]]},{"d":"two","e":[[]]}]`},
		{give: `[{"b":3},{"c":4},{"b":5}]`, when: `.(collect)`, then: `[{"b":3},{"c":4},{"b":5}]`},
		{give: `[{"b":3},{"c":4},{"b":5}]`, when: `.(collect).b`, then: `[3,5]`},
		{give: `{"a":{"b":{"c":[{"d":1,"e":{"f":[{"g":{"h":{"i":{"j":[{"k":3}]}}}}]}}]}}}`, when: `a.b.c.(collect).{d,data:e.f.(collect).g.h.i.j.(collect).k}`, then: `[{"d":1,"data":[[3]]}]`},
		{give: `{"a":{"b":[3,4]}}`, when: `{.a.b.(collect)}`, then: `{"b":[3,4]}`},
		{give: `{"a":[{"b":3}]}`, when: `a`, then: `[{"b":3}]`},
		{give: `{"a":[{"b":3}]}`, when: `a.(collect)`, then: `[{"b":3}]`},
		{give: `{"a":[{"b":3}]}`, when: `a.(collect).b`, then: `[3]`},
		{give: `{"a":[{"b":[{"c":3}]}]}`, when: `a.(collect).b`, then: `[[{"c":3}]]`},
		{give: `{"a":[{"b":[{"c":3}]}]}`, when: `a.(collect).b.(collect)`, then: `[[{"c":3}]]`},
		{give: `{"a":[{"b":[{"c":3}]}]}`, when: `a.(collect).b.(collect).c`, then: `[[3]]`},
		{give: `{"a":[{"b":[{"c":[{"d":3}]}]}]}`, when: `a.(collect).b.(collect).c`, then: `[[[{"d":3}]]]`},
		{give: `{"a":[{"b":[{"c":[{"d":3}]}]}]}`, when: `a.(collect).b.(collect).c.(collect)`, then: `[[[{"d":3}]]]`},
		{give: `{"a":[{"b":[{"c":[{"d":3}]}]}]}`, when: `a.(collect).b.(collect).c.(collect).d`, then: `[[[3]]]`},
		{give: `{"a":3}`, when: `{x:a}`, then: `{"x":3}`},
		{give: `{"a":{"b":3}}`, when: `{x:a.b}`, then: `{"x":3}`},
		{give: `{"a":[{"b":1},{"c":2},{"b":3}]}`, when: `a.(collect).{x:b}`, then: `[{"x":1},{},{"x":3}]`},
		{give: `{"a":[{"b":{"c":2}},{"b":{"d":3}}]}`, when: `a.(collect).{b.c,b.d}`, then: `[{"c":2},{"d":3}]`},
		{give: `{"a":[{"b":{"c":2}},{"b":{"d":3}}]}`, when: `a.(collect).{b.c}`, then: `[{"c":2},{}]`},
		{give: `{"a":[{"b":1},{"c":2},{"b":3}]}`, when: `a.(collect).{b,.c}`, then: `[{"b":1},{"c":2},{"b":3}]`},
		{give: `{"a":[{"b":1},{"c":2},{"b":3}]}`, when: `a.(collect).{b}`, then: `[{"b":1},{},{"b":3}]`},
		{give: `{"a":[{"b":1},{"c":2},{"b":3}]}`, when: `a.(collect).b`, then: `[1,3]`},
		{give: `{"a":[2,3]}`, when: `a.1`, then: `3`},
		{give: `{"a":2}`, when: `b`, then: ``},
		{give: `{"a":{"b":2}}`, when: `a.b`, then: `2`},
		{give: `{"a":[]}`, when: `a`, then: `[]`},
		{give: `{"a":{}}`, when: `a`, then: `{}`},
		{give: `{"a":2}`, when: `a`, then: `2`},
		{give: `{"a":"a\"b"}`, when: `a`, then: `"a\"b"`},
		{give: `2`, when: `.`, then: `2`},
		{give: `-2`, when: `.`, then: `-2`},
		{give: `false`, when: `.`, then: `false`},
		{give: `null`, when: `.`, then: `null`},
		{give: `{}`, when: `.`, then: `{}`},
		{give: `[]`, when: `.`, then: `[]`},
		{give: `[1,2]`, when: `.`, then: `[1,2]`},
		{give: `"a"`, when: `.`, then: `"a"`},
	}

	for _, tc := range tt {
		r := Get(tc.give, tc.when)
		assert.Equal(t, tc.then, r.String(), "TC: %v", tc)
	}
}

func TestJsonGet(t *testing.T) {

	tt := []struct {
		jsn string
		qry string
		exp string
	}{
		{jsn: `{"a":2,"b":3}`, qry: `a`, exp: `2`},
		{jsn: `{"a":2,"b":3}`, qry: `b`, exp: `3`},
		{jsn: `{"a":2,"b":3}`, qry: `c`, exp: ``},
		{jsn: `{"a":[2,3]}`, qry: `a`, exp: `[2,3]`},
		{jsn: `[2,3]`, qry: `0`, exp: `2`},
		{jsn: `[2,3]`, qry: `1`, exp: `3`},
		{jsn: `[2,3]`, qry: `2`, exp: ``},
		{jsn: `{"a":{"b":2}}`, qry: `a`, exp: `{"b":2}`},
	}

	for _, tc := range tt {
		j := New(tc.jsn)
		r := j.Get(tc.qry)
		assert.Equal(t, tc.exp, r.String(), "TC: %v", tc)
	}
}

func TestJsonCollect(t *testing.T) {

	tt := []struct {
		jsn string
		qry string
		exp string
	}{
		{jsn: `[{"a":3},{"a":4}]`, qry: `a`, exp: `[3,4]`},
		{jsn: `{"a":[2,3]}`, qry: `a`, exp: `[2,3]`},
		{jsn: `{"a":2,"b":3}`, qry: `a`, exp: `2`},
		{jsn: `{"a":2,"b":3}`, qry: `b`, exp: `3`},
		{jsn: `{"a":2,"b":3}`, qry: `c`, exp: ``},
		{jsn: `{"a":[2,3]}`, qry: `a`, exp: `[2,3]`},
		{jsn: `[2,3]`, qry: `0`, exp: `2`},
		{jsn: `[2,3]`, qry: `1`, exp: `3`},
		{jsn: `[2,3]`, qry: `2`, exp: ``},
		{jsn: `{"a":{"b":2}}`, qry: `a`, exp: `{"b":2}`},
	}

	for _, tc := range tt {
		j := New(tc.jsn)
		r := j.Collect(tc.qry)
		assert.Equal(t, tc.exp, r.String(), "TC: %v", tc)
	}
}

func ExampleJson_Iterate() {

	j := New(`{"a":"x","b":2,"c":{"a":3,"b":{"a":4,"b":[{"a":5},{"a":6,"b":7,"c":[8,9,0]}]}},"d":1}`)
	j.Iterate(func(k string, v Json) bool {
		fmt.Println(k, v.String())
		return false
	})

	// Output:
	// . {"a":"x","b":2,"c":{"a":3,"b":{"a":4,"b":[{"a":5},{"a":6,"b":7,"c":[8,9,0]}]}},"d":1}
	// .a "x"
	// .b 2
	// .c {"a":3,"b":{"a":4,"b":[{"a":5},{"a":6,"b":7,"c":[8,9,0]}]}}
	// .c.a 3
	// .c.b {"a":4,"b":[{"a":5},{"a":6,"b":7,"c":[8,9,0]}]}
	// .c.b.a 4
	// .c.b.b [{"a":5},{"a":6,"b":7,"c":[8,9,0]}]
	// .c.b.b.0 {"a":5}
	// .c.b.b.0.a 5
	// .c.b.b.1 {"a":6,"b":7,"c":[8,9,0]}
	// .c.b.b.1.a 6
	// .c.b.b.1.b 7
	// .c.b.b.1.c [8,9,0]
	// .c.b.b.1.c.0 8
	// .c.b.b.1.c.1 9
	// .c.b.b.1.c.2 0
	// .d 1
}

func TestJsonIterateObject(t *testing.T) {
	tt := []struct {
		inp string
		out []string
	}{
		{inp: `{}`, out: nil},
		{inp: `{"a":2}`, out: []string{"a", "2"}},
		{inp: `{"a":2,"b":3}`, out: []string{"a", "2", "b", "3"}},
		{inp: `{"a":{"b":2}}`, out: []string{"a", `{"b":2}`}},
		{inp: `{"a":[2]}`, out: []string{"a", "[2]"}},
	}
	for _, tc := range tt {
		var r []string
		j := New(tc.inp)
		j.IterateObject(func(k string, v Json) bool {
			r = append(r, k, v.String())
			return false
		})
		assert.Equal(t, tc.out, r, tc.inp)
	}
}

func TestJsonIterateArray(t *testing.T) {
	tt := []struct {
		inp string
		out []string
	}{
		{inp: "[]", out: nil},
		{inp: "[10]", out: []string{"0", "10"}},
		{inp: "[10,20]", out: []string{"0", "10", "1", "20"}},
		{inp: "[10,20,30]", out: []string{"0", "10", "1", "20", "2", "30"}},
		{inp: "[{},{},[10]]", out: []string{"0", "{}", "1", "{}", "2", "[10]"}},
	}
	for _, tc := range tt {
		var r []string
		j := New(tc.inp)
		j.IterateArray(func(k string, v Json) bool {
			r = append(r, k, v.String())
			return false
		})
		assert.Equal(t, tc.out, r, tc.inp)
	}
}

func TestJson_MatchObject(t *testing.T) {
	tt := []struct {
		inp string
		exp bool
	}{
		{inp: "", exp: false},
		{inp: "{", exp: false},
		{inp: "}", exp: false},
		{inp: "[{}]", exp: false},
		{inp: "{}", exp: true},
		{inp: `{"a":"b"}`, exp: true},
		{inp: `{"a":"b","c":"d"}`, exp: true},
		{inp: `{"a":{},"c":[]}`, exp: true},
		{inp: `{"a":{"b":1},"c":[2]}`, exp: true},
	}
	for _, tc := range tt {
		j := New(tc.inp)
		m := j.Mark()
		ok := j.MatchObject()
		assert.Equal(t, tc.exp, ok, tc.inp)
		if ok {
			assert.Equal(t, tc.inp, j.Token(m), tc.inp)
		}
	}
}

func TestJson_MatchArray(t *testing.T) {
	tt := []struct {
		inp string
		exp bool
	}{
		{inp: "", exp: false},
		{inp: "{[]}", exp: false},
		{inp: "[]", exp: true},
		{inp: "[true]", exp: true},
		{inp: "[false]", exp: true},
		{inp: "[null]", exp: true},
		{inp: "[123]", exp: true},
		{inp: "[{}]", exp: true},
		{inp: "[[]]", exp: true},
		{inp: `[true,1,false,null,[],{},"a"]`, exp: true},
	}
	for _, tc := range tt {
		j := New(tc.inp)
		m := j.Mark()
		ok := j.MatchArray()
		assert.Equal(t, tc.exp, ok, tc.inp)
		if ok {
			assert.Equal(t, tc.inp, j.Token(m), tc.inp)
		}
	}
}
