package jsqt

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {

	tt := []struct {
		give string
		when string
		then string
	}{
		// Iterate.
		{give: `{"a":3,"b":[{"c":4},{"c":5}],"d":[6,7]}`, when: `(iterate (.) num2str)`, then: `{"a":"3","b":[{"c":"4"},{"c":"5"}],"d":["6","7"]}`},
		{give: `{"a":3,"b":4}`, when: `(iterate (.) num2str)`, then: `{"a":"3","b":"4"}`},
		{give: `[3,4]`, when: `(iterate (.) num2str)`, then: `["3","4"]`},
		{give: `3`, when: `(iterate (.) num2str)`, then: `"3"`},
		// OmitEmpty.
		{give: `{"a":[[3],[]]}`, when: `(collect (get a) () (omitempty (.)))`, then: `[[3]]`},
		{give: `{"a":[{"b":3},{"c":4}]}`, when: `(get a (omitempty (obj x (get b))))`, then: `[{"x":3}]`},
		{give: `{"a":{}}`, when: `(omitempty (get a))`, then: ``},
		// Default.
		{give: `[{"b":3},{"c":4},{"b":5}]`, when: `(collect (.) () (default (get b) 0))`, then: `[3,0,5]`},
		// Size.
		{give: `[3,4]`, when: `(size (.))`, then: `2`},
		// Merge.
		{give: `[{"a":3},{"b":4}]`, when: `(merge (.))`, then: `{"a":3,"b":4}`},
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
		{give: `{"a":3,"b":4}`, when: `(obj "a b" (get a) y (get b))`, then: `{"a b":3,"y":4}`},
		{give: `{"a":3,"b":4}`, when: `(obj x (get a) y (get b))`, then: `{"x":3,"y":4}`},
		// Get.
		{
			give: `{"a":{"b":{"c":[{"d":"one","e":{"f":[{"g":{"h":{"i":{"j":[{"k":{"l":"hi"}}]}}}}]}},{"d":"two","e":{"f":[{"g":{"h":{"i":{"j":[]}}}}]}}]}}}`,
			when: `(get a b c (obj d (get d) e (flatten (get e f g h i j k l))))`,
			then: `[{"d":"one","e":["hi"]},{"d":"two","e":[]}]`,
		},
		{give: `{"a":[{"b":3},{"c":4}]}`, when: `(get a b)`, then: `[3]`},
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
		r := Get(tc.give, tc.when)
		assert.Equal(t, tc.then, r.String(), "TC: %v", tc)
	}
}

func TestJsonWS(t *testing.T) {

	tt := []struct {
		give string
		when string
		then string
	}{
		{give: `{  "a"  :	3,  "b"  :  [  {  "c"  :  4  }  ]  }`, when: `(get b c)`, then: `[4]`},
		{give: `[3,4 ]`, when: `(get 1)`, then: `4`},
		{give: `[3, 4]`, when: `(get 1)`, then: `4`},
		{give: `[3 ,4]`, when: `(get 1)`, then: `4`},
		{give: `[ 3,4]`, when: `(get 1)`, then: `4`},
		{give: `[ ]`, when: `(get (.))`, then: `[ ]`},
		{give: `{"a":3, "b":4}`, when: `(get b)`, then: `4`},
		{give: `{"a":3 ,"b":4}`, when: `(get a)`, then: `3`},
		{give: `{"a": 3}`, when: `(get a)`, then: `3`},
		{give: `{"a" :3}`, when: `(get a)`, then: `3`},
		{give: `{ "a":3}`, when: `(get a)`, then: `3`},
		{give: `{ }`, when: `(get (.))`, then: `{ }`},
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

	j := New(`{"a":1,"b":2,"c":{"a":3,"b":{"a":4,"b":[{"a":5},{"a":6,"b":7,"c":[8,9,0]}]}},"d":1}`)
	v := j.Iterate(func(k string, v Json) (string, string) {
		if v.IsNumber() {
			return strings.ToUpper(k), `"` + v.String() + `"`
		}
		return strings.ToUpper(k), v.String()
	})

	fmt.Println(v)

	// Output:
	// {"A":"1","B":"2","C":{"A":"3","B":{"A":"4","B":[{"A":"5"},{"A":"6","B":"7","C":["8","9","0"]}]}},"D":"1"}
}

func TestJsonForEachKeyVal(t *testing.T) {
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
		j.ForEachKeyVal(func(k string, v Json) bool {
			r = append(r, k, v.String())
			return false
		})
		assert.Equal(t, tc.out, r, tc.inp)
	}
}

func TestJsonForEach(t *testing.T) {
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
		j.ForEach(func(k string, v Json) bool {
			r = append(r, k, v.String())
			return false
		})
		assert.Equal(t, tc.out, r, tc.inp)
	}
}
