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
		// Entries.
		{give: `{"a":3,"b":4}`, when: `(collect (entries) (flatten))`, then: `["a",3,"b",4]`},
		{give: `{"a":3,"b":4}`, when: `(entries)`, then: `[["a",3],["b",4]]`},
		// Values.
		{give: `{"a":3,"b":4}`, when: `(values)`, then: `[3,4]`},
		// Keys.
		{give: `{"a":3,"b":4}`, when: `(keys)`, then: `["a","b"]`},
		// If.
		{give: `{"a":""}`, when: `(get a (if (is-str) (raw {}) (raw 3)))`, then: `{}`},       // Then.
		{give: `{"a":{"b":3}}`, when: `(get a (if (is-str) (raw {}) (.)))`, then: `{"b":3}`}, // Else.
		{give: `3`, when: `(if (is-num) (obj b (.)) (raw 3))`, then: `{"b":3}`},              // Then.
		{give: `{"b":3}`, when: `(if (is-num) (raw 3) (.))`, then: `{"b":3}`},                // Else.
		// IsNull
		{give: `3`, when: `(is-null)`, then: ``},
		{give: `null`, when: `(is-null)`, then: `null`},
		// IsBool
		{give: `3`, when: `(is-bool)`, then: ``},
		{give: `true`, when: `(is-bool)`, then: `true`},
		{give: `false`, when: `(is-bool)`, then: `false`},
		// IsStr
		{give: `3`, when: `(is-str)`, then: ``},
		{give: `"3"`, when: `(is-str)`, then: `"3"`},
		// IsArr
		{give: `3`, when: `(is-arr)`, then: ``},
		{give: `[]`, when: `(is-arr)`, then: `[]`},
		// IsObj
		{give: `3`, when: `(is-obj)`, then: ``},
		{give: `{}`, when: `(is-obj)`, then: `{}`},
		// IsNum
		{give: `"3"`, when: `(is-num)`, then: ``},
		{give: `3`, when: `(is-num)`, then: `3`},
		// Ugly / Nice / Pretty.
		{give: `[ { "a" : 3 , "b" : [ 4 , { "c" : 5 } ], "c": [] } ]`, when: `(pretty)`, then: "[\n    {\n        \"a\": 3,\n        \"b\": [\n            4,\n            {\n                \"c\": 5\n            }\n        ],\n        \"c\": []\n    }\n]"},
		{give: `[ { "a" : 3 , "b" : [ 4 , { "c" : 5 } ], "c": [] } ]`, when: `(nice)`, then: `[{ "a": 3, "b": [4, { "c": 5 }], "c": [] }]`},
		{give: `[ { "a" : 3 , "b" : [ 4 , { "c" : 5 } ], "c": [] } ]`, when: `(ugly)`, then: `[{"a":3,"b":[4,{"c":5}],"c":[]}]`},
		// Filters.
		{give: `[{"a":3,"b":6},{"a":4,"b":7},{"a":5,"b":8}]`, when: `(collect (== b 7) a)`, then: `[4]`},
		{give: `[{"a":3,"b":6},{"a":4,"b":7},{"a":5,"b":8}]`, when: `(collect (!= b 7) a)`, then: `[3,5]`},
		{give: `[{"a":3,"b":6},{"a":4,"b":7},{"a":5,"b":8}]`, when: `(collect (>= b 7) a)`, then: `[4,5]`},
		{give: `[{"a":3,"b":6},{"a":4,"b":7},{"a":5,"b":8}]`, when: `(collect (<= b 7) a)`, then: `[3,4]`},
		{give: `[{"a":3,"b":6},{"a":4,"b":7},{"a":5,"b":8}]`, when: `(collect (> b 7) a)`, then: `[5]`},
		{give: `[{"a":3,"b":6},{"a":4,"b":7},{"a":5,"b":8}]`, when: `(collect (< b 7) a)`, then: `[3]`},
		// Iterate.
		{give: `{"a":3,"b":[{"c":4},{"c":5}],"d":[6,7]}`, when: `(iterate num2str)`, then: `{"a":"3","b":[{"c":"4"},{"c":"5"}],"d":["6","7"]}`},
		{give: `{"a":3,"b":4}`, when: `(iterate num2str)`, then: `{"a":"3","b":"4"}`},
		{give: `[3,4]`, when: `(iterate num2str)`, then: `["3","4"]`},
		{give: `3`, when: `(iterate num2str)`, then: `"3"`},
		// OmitEmpty.
		{give: `{"a":[[3],[]]}`, when: `(collect a (omitempty))`, then: `[[3]]`},
		{give: `{"a":[{"b":3},{"c":4}]}`, when: `(collect a (obj x b) (omitempty))`, then: `[{"x":3}]`},
		{give: `{"a":{}}`, when: `(get a (omitempty))`, then: ``},
		// Default.
		{give: `[{"b":3},{"c":4},{"b":5}]`, when: `(collect b (default 0))`, then: `[3,0,5]`},
		// Size.
		{give: `[3,4]`, when: `(size)`, then: `2`},
		// Merge.
		{give: `[{"a":3},{"b":4}]`, when: `(merge)`, then: `{"a":3,"b":4}`},
		// Join.
		{
			give: `[{"c":"one","d":"one-val"},{"c":"two","d":"two-val"}]`,
			when: `(join c d)`,
			then: `{"one":"one-val","two":"two-val"}`,
		},
		// Collect.
		{
			give: `{"a":{"b":{"c":[{"d":"one","e":{"f":[{"g":{"h":{"i":{"j":[{"k":{"l":"hi"}}]}}}}]}},{"d":"two","e":{"f":[{"g":{"h":{"i":{"j":[]}}}}]}}]}}}`,
			when: `(collect a b c (obj x d e (collect e f g h i j (flatten) k l)))`,
			then: `[{"x":"one","e":["hi"]},{"x":"two","e":[]}]`,
		},
		{give: `{"a":[{"b":{"c":3}},{"b":{}}]}`, when: `(collect a b c)`, then: `[3]`},
		{give: `{"a":[{"b":3},{"b":4}]}`, when: `(collect a b)`, then: `[3,4]`},
		{give: `[{"a":3},{"b":4},{"a":5}]`, when: `(collect a)`, then: `[3,5]`},
		// Array.
		{give: `{"a":3,"b":4}`, when: `(arr a b a (raw "hi"))`, then: `[3,4,3,"hi"]`},
		{give: `{"a":3,"b":4}`, when: `(arr a b a)`, then: `[3,4,3]`},
		{give: `{"a":3,"b":4}`, when: `(arr (get a) (get b) (get a))`, then: `[3,4,3]`},
		// Object.
		{give: `{"a":"aaa","b":"bbb"}}`, when: `(obj (get a) (get b))`, then: `{"aaa":"bbb"}`},
		{give: `{"a":3,"b":4}`, when: `(obj "a b" a y b)`, then: `{"a b":3,"y":4}`},
		{give: `{"a":3,"b":4}`, when: `(obj "a b" (get a) y (get b))`, then: `{"a b":3,"y":4}`},
		{give: `{"a":3,"b":4}`, when: `(obj x (get a) y (get b))`, then: `{"x":3,"y":4}`},
		{give: `{"a":{"b":{"c":3}}}`, when: `(get a b (obj x c))`, then: `{"x":3}`},
		// Get.
		{give: `{"a":[{"b":3},{"c":4}]}`, when: `(get a 1 c)`, then: `4`},
		{give: `{"a":[{"b":3},{"c":4}]}`, when: `(get a 0 b)`, then: `3`},
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
		{give: ` {  "a"  :	3,  "b"  :  [  {  "c"  :  4  }  ]  }`, when: `(get b 0 c)`, then: `4`},
		{give: `[3,4 ]`, when: `(get 1)`, then: `4`},
		{give: `[3, 4]`, when: `(get 1)`, then: `4`},
		{give: `[3 ,4]`, when: `(get 1)`, then: `4`},
		{give: `[ 3,4]`, when: `(get 1)`, then: `4`},
		{give: `{"a":3, "b":4}`, when: `(get b)`, then: `4`},
		{give: `{"a":3 ,"b":4}`, when: `(get a)`, then: `3`},
		{give: `{"a": 3}`, when: `(get a)`, then: `3`},
		{give: `{"a" :3}`, when: `(get a)`, then: `3`},
		{give: `{ "a":3}`, when: `(get a)`, then: `3`},
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

func TestJsonStr(t *testing.T) {
	tt := []struct {
		inp string
		out string
	}{
		{inp: ``, out: ``},
		{inp: `""`, out: ``},
		{inp: `"a"`, out: `a`},
		{inp: `""a""`, out: `"a"`},
	}
	for _, tc := range tt {
		j := New(tc.inp)
		assert.Equal(t, tc.out, j.Str(), tc.inp)
	}
}

func TestJsonInt(t *testing.T) {
	tt := []struct {
		inp string
		out int
	}{
		{inp: ``, out: 0},
		{inp: `0`, out: 0},
		{inp: `1`, out: 1},
		{inp: `-2`, out: -2},
		{inp: `100`, out: 100},
	}
	for _, tc := range tt {
		j := New(tc.inp)
		assert.Equal(t, tc.out, j.Int(), tc.inp)
	}
}

func TestJsonFloat(t *testing.T) {
	tt := []struct {
		inp string
		out float64
	}{
		{inp: ``, out: 0},
		{inp: `0`, out: 0},
		{inp: `1`, out: 1},
		{inp: `-2.2`, out: -2.2},
		{inp: `1e2`, out: 100},
	}
	for _, tc := range tt {
		j := New(tc.inp)
		assert.Equal(t, tc.out, j.Float(), tc.inp)
	}
}

func TestJsonBool(t *testing.T) {
	tt := []struct {
		inp string
		out bool
	}{
		{inp: ``, out: false},
		{inp: `false`, out: false},
		{inp: `true`, out: true},
	}
	for _, tc := range tt {
		j := New(tc.inp)
		assert.Equal(t, tc.out, j.Bool(), tc.inp)
	}
}

func ExampleFuncDebug() {

	j := `[{ "a": { "b": [3] } }, { "a": { "b": [4] } }]`

	v := Get(j, `(collect a (debug) b (debug b_val) (flatten) (debug flatd))`)

	fmt.Println("Result:", v)

	// Output:
	// [debug] { "b": [3] }
	// [b_val] [3]
	// [flatd] 3
	// [debug] { "b": [4] }
	// [b_val] [4]
	// [flatd] 4
	// Result: [3,4]
}
