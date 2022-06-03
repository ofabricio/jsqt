package jsqt

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {

	tt := []struct {
		give string
		when string
		then string
	}{
		{give: `{"a":{"b":{"c":[{"d":"one","e":{"f":[{"g":{"h":{"i":{"j":[{"k":{"l":"hi"}}]}}}}]}},{"d":"two","e":{"f":[{"g":{"h":{"i":{"j":[]}}}}]}}]}}}`, when: `a.b.c.*.{d:d,e:e.f.*.g.h.i.j.*.k.l}`, then: `[{"d":"one","e":[["hi"]]},{"d":"two","e":[[]]}]`},
		{give: `[{"b":3},{"c":4},{"b":5}]`, when: `*`, then: `[{"b":3},{"c":4},{"b":5}]`},
		{give: `[{"b":3},{"c":4},{"b":5}]`, when: `.*`, then: `[{"b":3},{"c":4},{"b":5}]`},
		{give: `[{"b":3},{"c":4},{"b":5}]`, when: `*.b`, then: `[3,5]`},
		{give: `[{"b":3},{"c":4},{"b":5}]`, when: `.*.b`, then: `[3,5]`},
		{give: `{"a":{"b":{"c":[{"d":1,"e":{"f":[{"g":{"h":{"i":{"j":[{"k":3}]}}}}]}}]}}}`, when: `a.b.c.*.{d,data:e.f.*.g.h.i.j.*.k}`, then: `[{"d":1,"data":[[3]]}]`},
		{give: `{"a":{"b":[3,4]}}`, when: `{a.b.*}`, then: `{"":[3,4]}`},
		{give: `{"a":[{"b":3}]}`, when: `a`, then: `[{"b":3}]`},
		{give: `{"a":[{"b":3}]}`, when: `a.*`, then: `[{"b":3}]`},
		{give: `{"a":[{"b":3}]}`, when: `a.*.b`, then: `[3]`},
		{give: `{"a":[{"b":[{"c":3}]}]}`, when: `a.*.b`, then: `[[{"c":3}]]`},
		{give: `{"a":[{"b":[{"c":3}]}]}`, when: `a.*.b.*`, then: `[[{"c":3}]]`},
		{give: `{"a":[{"b":[{"c":3}]}]}`, when: `a.*.b.*.c`, then: `[[3]]`},
		{give: `{"a":[{"b":[{"c":[{"d":3}]}]}]}`, when: `a.*.b.*.c`, then: `[[[{"d":3}]]]`},
		{give: `{"a":[{"b":[{"c":[{"d":3}]}]}]}`, when: `a.*.b.*.c.*`, then: `[[[{"d":3}]]]`},
		{give: `{"a":[{"b":[{"c":[{"d":3}]}]}]}`, when: `a.*.b.*.c.*.d`, then: `[[[3]]]`},
		{give: `{"a":3}`, when: `{a}`, then: `{"a":3}`},
		{give: `{"a":3}`, when: `{.a}`, then: `{"a":3}`},
		{give: `{"a":3}`, when: `{x:a}`, then: `{"x":3}`},
		{give: `{"a":{"b":3}}`, when: `{x:a.b}`, then: `{"x":3}`},
		{give: `{"a":{"b":3}}`, when: `{x:.a.b}`, then: `{"x":3}`},
		{give: `{"a":[{"b":1},{"c":2},{"b":3}]}`, when: `.a.*.{x:b}`, then: `[{"x":1},{},{"x":3}]`},
		{give: `{"a":[{"b":{"c":2}},{"b":{"d":3}}]}`, when: `.a.*.{b.c,b.d}`, then: `[{"c":2},{"d":3}]`},
		{give: `{"a":[{"b":{"c":2}},{"b":{"d":3}}]}`, when: `.a.*.{b.c}`, then: `[{"c":2},{}]`},
		{give: `{"a":[{"b":1},{"c":2},{"b":3}]}`, when: `.a.*.{b,c}`, then: `[{"b":1},{"c":2},{"b":3}]`},
		{give: `{"a":[{"b":1},{"c":2},{"b":3}]}`, when: `.a.*.{b}`, then: `[{"b":1},{},{"b":3}]`},
		{give: `{"a":[{"b":1},{"c":2},{"b":3}]}`, when: `.a.*.b`, then: `[1,3]`},
		{give: `{"a":[2,3]}`, when: `.a.1`, then: `3`},
		{give: `{"a":2}`, when: `.b`, then: ``},
		{give: `{"a":{"b":2}}`, when: `.a.b`, then: `2`},
		{give: `{"a":[]}`, when: `.a`, then: `[]`},
		{give: `{"a":{}}`, when: `.a`, then: `{}`},
		{give: `{"a":2}`, when: `.a`, then: `2`},
		{give: `2`, when: `.`, then: `2`},
		{give: `-2`, when: `.`, then: `-2`},
		{give: `false`, when: `.`, then: `false`},
		{give: `null`, when: `.`, then: `null`},
		{give: `{}`, when: `.`, then: `{}`},
		{give: `[]`, when: `.`, then: `[]`},
		{give: `[1,2]`, when: `.`, then: `[1,2]`},
		{give: `bla`, when: `.`, then: `bla`}, // Invalid Json though.
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
