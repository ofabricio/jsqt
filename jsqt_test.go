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
		// Test empty context input for many functions.
		{give: `["3 4","5 6","7 8"]`, when: `(collect (!= (this) "5 6") (replace " " "_"))`, then: `["3_4","7_8"]`},
		{give: `[["3","4"],["5","6"],["7","8"]]`, when: `(collect (!= 1 "6") (join "_"))`, then: `["3_4","7_8"]`},
		{give: `[["3","4"],["5","6"],["7","8"]]`, when: `(collect (!= 1 "6") (concat 0 1))`, then: `["34","78"]`},
		{give: `[[3,4],[5,6],[7,8]]`, when: `(collect (!= 1 6) (reverse))`, then: `[[4,3],[8,7]]`},
		{give: `[[3,4],[5,6],[7,8]]`, when: `(collect (!= 1 6) (sort desc))`, then: `[[4,3],[8,7]]`},
		// (this)
		{give: `[3,4]`, when: `(collect (> (this) 3))`, then: `[4]`},
		// (join)
		{give: `["a","b","c"]`, when: `(join "_")`, then: `"a_b_c"`},
		// (reverse)
		{give: `[{"a":3},{"b":4}]`, when: `(reverse)`, then: `[{"b":4},{"a":3}]`},
		{give: `[3,4,5,6,2,7]`, when: `(reverse)`, then: `[7,2,6,5,4,3]`},
		// (sort)
		{give: `{"c":5,"b":4,"a":3}`, when: `(sort asc)`, then: `{"a":3,"b":4,"c":5}`},
		{give: `{"a":5,"b":4,"c":3}`, when: `(sort desc)`, then: `{"c":3,"b":4,"a":5}`},
		{give: `[5,3,4]`, when: `(sort desc)`, then: `[5,4,3]`},
		{give: `[5,4,3]`, when: `(sort asc)`, then: `[3,4,5]`},
		{give: `[{"a":3},{"a":4},{"a":5}]`, when: `(sort desc a)`, then: `[{"a":5},{"a":4},{"a":3}]`},
		{give: `[{"a":5},{"a":4},{"a":3}]`, when: `(sort asc a)`, then: `[{"a":3},{"a":4},{"a":5}]`},
		// (either) should skip its arguments.
		{give: `{"a":"","b":"B","c":""}`, when: `(get (either a b c) (lower))`, then: `"b"`},
		// (either)
		{give: `{"a":"","b":"","c":""}`, when: `(either a b c)`, then: `""`},
		{give: `{"a":"A","b":"","c":"C"}`, when: `(either a b c)`, then: `"A"`},
		{give: `{"a":"A","b":"B","c":""}`, when: `(either a b c)`, then: `"A"`},
		{give: `{"a":"A","b":"","c":""}`, when: `(either a b c)`, then: `"A"`},
		{give: `{"a":"","b":"","c":"C"}`, when: `(either a b c)`, then: `"C"`},
		{give: `{"a":"","b":"B","c":""}`, when: `(either a b c)`, then: `"B"`},
		{give: `{"a":"","b":"B","c":"C"}`, when: `(either a b c)`, then: `"B"`},
		{give: `{"a":"A","b":"B","c":"C"}`, when: `(either a b c)`, then: `"A"`},
		{give: `{"a":"","b":""}`, when: `(either a b)`, then: `""`},
		{give: `{"a":"","b":"B"}`, when: `(either a b)`, then: `"B"`},
		{give: `{"a":"A","b":""}`, when: `(either a b)`, then: `"A"`},
		{give: `{"a":"A","b":"B"}`, when: `(either a b)`, then: `"A"`},
		// (concat)
		{give: `{ "one": "hello" }`, when: `(concat one (raw " \"world\""))`, then: `"hello \"world\""`},
		{give: `{ "one": "hello", "two": "world" }`, when: `(concat one (raw " ") two)`, then: `"hello world"`},
		//(replace)
		{give: `"a b"`, when: `(replace " " "_")`, then: `"a_b"`},
		// (upper)
		{give: `"a"`, when: `(upper)`, then: `"A"`},
		// (lower)
		{give: `"A"`, when: `(lower)`, then: `"a"`},
		// (stringify)
		{give: `3`, when: `(stringify)`, then: `"3"`},
		{give: `-3`, when: `(stringify)`, then: `"-3"`},
		{give: `[]`, when: `(stringify)`, then: `"[]"`},
		{give: `{}`, when: `(stringify)`, then: `"{}"`},
		{give: `""`, when: `(stringify)`, then: `"\"\""`},
		{give: `{"a":[{"b":3},4,"5"]}`, when: `(stringify)`, then: `"{\"a\":[{\"b\":3},4,\"5\"]}"`},
		// (jsonify)
		{give: `"{\"a\":[{\"b\":3},4,\"5\"]}"`, when: `(jsonify)`, then: `{"a":[{"b":3},4,"5"]}`},
		{give: `"3"`, when: `(jsonify)`, then: `3`},
		{give: `"{}"`, when: `(jsonify)`, then: `{}`},
		// (bool)
		{give: `3`, when: `(get (is-num) (bool))`, then: `true`},
		{give: `{}`, when: `(get (is-num) (bool))`, then: `false`},
		// (or) (and) (not)
		{give: `[3,"",4,"5"]`, when: `(collect (not (is-str)))`, then: `[3,4]`},
		{give: `[{"a":3},{"a":4},{"a":5},{"a":6}]`, when: `(collect (or (< a 4) (> a 5)) a)`, then: `[3,6]`},
		{give: `[{"a":3},{"a":4},{"a":5},{"a":6}]`, when: `(collect (and (>= a 4) (<= a 5)) a)`, then: `[4,5]`},
		// (objectify)
		{give: `[["a",3],["b",4]]`, when: `(objectify)`, then: `{"a":3,"b":4}`},
		// (entries)
		{give: `{"a":3,"b":4}`, when: `(collect (entries) (flatten))`, then: `["a",3,"b",4]`},
		{give: `{"a":3,"b":4}`, when: `(entries)`, then: `[["a",3],["b",4]]`},
		//(values)
		{give: `{"a":3,"b":4}`, when: `(values)`, then: `[3,4]`},
		// (key)
		{give: `{"a":3,"b":4}`, when: `(keys)`, then: `["a","b"]`},
		// (if)
		{give: `{"a":""}`, when: `(get a (if (is-str) (raw {}) (raw 3)))`, then: `{}`},          // Then.
		{give: `{"a":{"b":3}}`, when: `(get a (if (is-str) (raw {}) (this)))`, then: `{"b":3}`}, // Else.
		{give: `3`, when: `(if (is-num) (obj b (this)) (raw 3))`, then: `{"b":3}`},              // Then.
		{give: `{"b":3}`, when: `(if (is-num) (raw 3) (this))`, then: `{"b":3}`},                // Else.
		// (is-void)
		{give: `{}`, when: `(is-void)`, then: `{}`},
		{give: `[]`, when: `(is-void)`, then: `[]`},
		{give: `""`, when: `(is-void)`, then: ``},
		// (is-blank)
		{give: `{}`, when: `(is-blank)`, then: `{}`},
		{give: `[]`, when: `(is-blank)`, then: `[]`},
		{give: `null`, when: `(is-blank)`, then: `null`},
		{give: `""`, when: `(is-blank)`, then: ``},
		// (is-nully)
		{give: `{}`, when: `(is-nully)`, then: `{}`},
		{give: `[]`, when: `(is-nully)`, then: `[]`},
		{give: `null`, when: `(is-nully)`, then: `null`},
		{give: `""`, when: `(is-nully)`, then: `""`},
		// (is-some)
		{give: `3`, when: `(is-some)`, then: `3`},
		{give: `""`, when: `(is-some)`, then: `""`},
		{give: `null`, when: `(is-some)`, then: ``},
		// (nothing) - Undocumented.
		{give: `3`, when: `(nothing)`, then: ``},
		// (exists)
		{give: `3`, when: `(exists)`, then: `3`},
		{give: ``, when: `(exists)`, then: ``},
		// (truthy)
		{give: `{}`, when: `(truthy)`, then: ``},
		{give: `[]`, when: `(truthy)`, then: ``},
		{give: `0`, when: `(truthy)`, then: ``},
		{give: `""`, when: `(truthy)`, then: ``},
		{give: `null`, when: `(truthy)`, then: ``},
		{give: `[0]`, when: `(truthy)`, then: `[0]`},
		{give: `3`, when: `(truthy)`, then: `3`},
		{give: `{"a":3}`, when: `(truthy)`, then: `{"a":3}`},
		{give: `true`, when: `(truthy)`, then: `true`},
		// (falsy)
		{give: `{}`, when: `(falsy)`, then: `{}`},
		{give: `[]`, when: `(falsy)`, then: `[]`},
		{give: `0`, when: `(falsy)`, then: `0`},
		{give: `""`, when: `(falsy)`, then: `""`},
		{give: `false`, when: `(falsy)`, then: `false`},
		{give: `[0]`, when: `(falsy)`, then: ``},
		{give: `3`, when: `(falsy)`, then: ``},
		{give: `null`, when: `(falsy)`, then: `null`},
		// (is-empty-obj)
		{give: `{}`, when: `(is-empty-obj)`, then: `{}`},
		{give: `{"a":3}`, when: `(is-empty-obj)`, then: ``},
		{give: `[]`, when: `(is-empty-obj)`, then: ``},
		// (is-empty-arr)
		{give: `{}`, when: `(is-empty-arr)`, then: ``},
		{give: `[0]`, when: `(is-empty-arr)`, then: ``},
		{give: `[]`, when: `(is-empty-arr)`, then: `[]`},
		// (is-empty-str)
		{give: `3`, when: `(is-empty-str)`, then: ``},
		{give: `""`, when: `(is-empty-str)`, then: `""`},
		// (is-empty)
		{give: `3`, when: `(is-empty)`, then: ``},
		{give: `{}`, when: `(is-empty)`, then: `{}`},
		{give: `[]`, when: `(is-empty)`, then: `[]`},
		{give: `""`, when: `(is-empty)`, then: `""`},
		// (is-null)
		{give: `3`, when: `(is-null)`, then: ``},
		{give: `null`, when: `(is-null)`, then: `null`},
		// (is-bool)
		{give: `3`, when: `(is-bool)`, then: ``},
		{give: `true`, when: `(is-bool)`, then: `true`},
		{give: `false`, when: `(is-bool)`, then: `false`},
		// (is-str)
		{give: `3`, when: `(is-str)`, then: ``},
		{give: `"3"`, when: `(is-str)`, then: `"3"`},
		// (is-arr)
		{give: `3`, when: `(is-arr)`, then: ``},
		{give: `[]`, when: `(is-arr)`, then: `[]`},
		// (is-obj)
		{give: `3`, when: `(is-obj)`, then: ``},
		{give: `{}`, when: `(is-obj)`, then: `{}`},
		// (is-num)
		{give: `"3"`, when: `(is-num)`, then: ``},
		{give: `3`, when: `(is-num)`, then: `3`},
		// (ugly) (pretty)
		{give: `[ { "a" : 3 , "b" : [ 4 , { "c" : 5, "d": "e f" } ], "c": [ ], "d": { } } ]`, when: `(pretty)`, then: "[\n    {\n        \"a\": 3,\n        \"b\": [\n            4,\n            {\n                \"c\": 5,\n                \"d\": \"e f\"\n            }\n        ],\n        \"c\": [],\n        \"d\": {}\n    }\n]"},
		{give: `[ { "a" : 3 , "b" : [ 4 , { "c" : 5, "d": "e f" } ], "c": [ ], "d": { } } ]`, when: `(ugly)`, then: `[{"a":3,"b":[4,{"c":5,"d":"e f"}],"c":[],"d":{}}]`},
		// (filters)
		{give: `[{"a":3,"b":6},{"a":4,"b":7},{"a":5,"b":8}]`, when: `(collect (== b 7) a)`, then: `[4]`},
		{give: `[{"a":3,"b":6},{"a":4,"b":7},{"a":5,"b":8}]`, when: `(collect (!= b 7) a)`, then: `[3,5]`},
		{give: `[{"a":3,"b":6},{"a":4,"b":7},{"a":5,"b":8}]`, when: `(collect (>= b 7) a)`, then: `[4,5]`},
		{give: `[{"a":3,"b":6},{"a":4,"b":7},{"a":5,"b":8}]`, when: `(collect (<= b 7) a)`, then: `[3,4]`},
		{give: `[{"a":3,"b":6},{"a":4,"b":7},{"a":5,"b":8}]`, when: `(collect (> b 7) a)`, then: `[5]`},
		{give: `[{"a":3,"b":6},{"a":4,"b":7},{"a":5,"b":8}]`, when: `(collect (< b 7) a)`, then: `[3]`},
		// (iterate-kv)
		{give: "{ \"a\"\t:\t3\t}", when: `(iterate-kv (this))`, then: `{"a":3}`},
		{give: `{ "a" : 3, "b": [ 3 , { "c": "d" } ] }`, when: `(iterate-kv (upper))`, then: `{"A":3,"B":[3,{"C":"D"}]}`},
		{give: `{ "a": "b" }`, when: `(iterate-kv (upper))`, then: `{"A":"B"}`},
		{give: `3`, when: `(iterate-kv (stringify))`, then: `"3"`},
		{give: `3`, when: `(iterate-kv (this))`, then: `3`},
		// (iterate-v)
		{give: "{ \"a\"\t:\t3\t}", when: `(iterate-v (this))`, then: `{"a":3}`},
		{give: `{ "a" : 3, "b": [ 3 , { "a" : 3 } ] }`, when: `(iterate-v (if (== (this) 3) (raw 4) (this))))`, then: `{"a":4,"b":[4,{"a":4}]}`},
		{give: `{ "a": 3 }`, when: `(iterate-v (if (== (this) 3) (raw 4) (this)))`, then: `{"a":4}`},
		{give: `3`, when: `(iterate-v (stringify))`, then: `"3"`},
		{give: `3`, when: `(iterate-v (this))`, then: `3`},
		// (iterate-k)
		{give: "{ \"a\"\t:\t3\t}", when: `(iterate-k (this))`, then: `{"a":3}`},
		{give: `{ "a" : 3 , "b" : [ 3, { "a": 3 } ] }`, when: `(iterate-k (if (== (this) "a") (raw "x") (this))))`, then: `{"x":3,"b":[3,{"x":3}]}`},
		{give: `{ "a" : 3 }`, when: `(iterate-k (if (== (this) "a") (raw "x") (this)))`, then: `{"x":3}`},
		{give: `3`, when: `(iterate-k (stringify))`, then: `3`},
		{give: `3`, when: `(iterate-k (this))`, then: `3`},
		// (iterate)
		{give: "{ \"a\"\t:\t3\t}", when: `(iterate 0 1)`, then: `{"a":3}`},
		{give: `{ "a": "aaa", "b" : "bbb" }`, when: `(iterate (if (get 1 (is-str)) (get 1) (get 0)) (if (get 0 (is-str)) (get 0) (get 1)))`, then: `{"aaa":"a","bbb":"b"}`},
		{give: `{ "a" : 3 , "b" : [ { "c" : 4 } , { "c" : 5 } ] , "d" : [ 6 , true ] }`, when: `(iterate 0 (get 1 (if (is-num) (stringify) (this))))`, then: `{"a":"3","b":[{"c":"4"},{"c":"5"}],"d":["6",true]}`},
		{give: `{ "a" : 3, "b" : 4}`, when: `(iterate 0 (get 1 (if (is-num) (stringify) (this))))`, then: `{"a":"3","b":"4"}`},
		{give: `[3,4]`, when: `(iterate 0 (get 1 (if (is-num) (stringify) (this))))`, then: `["3","4"]`},
		{give: `3`, when: `(iterate 0 (get 1 (stringify)))`, then: `"3"`},
		{give: `3`, when: `(iterate 0 1)`, then: `3`},
		// (default)
		{give: `[{"b":3},{"c":4},{"b":5}]`, when: `(collect b (default 0))`, then: `[3,0,5]`},
		// (size)
		{give: `{"a":3,"b":4}`, when: `(size)`, then: `2`},
		{give: `"abc"`, when: `(size)`, then: `3`},
		{give: `[3,4]`, when: `(size)`, then: `2`},
		// (merge)
		{give: `[{"a":3},{"b":4}]`, when: `(merge)`, then: `{"a":3,"b":4}`},
		// (collect)
		{
			give: `{"a":{"b":{"c":[{"d":"one","e":{"f":[{"g":{"h":{"i":{"j":[{"k":{"l":"hi"}}]}}}}]}},{"d":"two","e":{"f":[{"g":{"h":{"i":{"j":[]}}}}]}}]}}}`,
			when: `(collect a b c (obj x d e (collect e f g h i j (flatten) k l)))`,
			then: `[{"x":"one","e":["hi"]},{"x":"two","e":[]}]`,
		},
		{give: `{"a":[3,4,5]}`, when: `(collect a (this))`, then: `[3,4,5]`},
		{give: `[{"a":3}]`, when: `(collect a)`, then: `[3]`},
		{give: ``, when: `(collect)`, then: `[]`},
		{give: `[]`, when: `(collect a)`, then: `[]`},
		{give: `{"a":[{"b":{"c":3}},{"b":{}}]}`, when: `(collect a b c)`, then: `[3]`},
		{give: `{"a":[{"b":3},{"b":4}]}`, when: `(collect a b)`, then: `[3,4]`},
		{give: `[{"a":3},{"b":4},{"a":5}]`, when: `(collect a)`, then: `[3,5]`},
		// (array)
		{give: `{"a":3,"b":4}`, when: `(arr a b a (raw "hi"))`, then: `[3,4,3,"hi"]`},
		{give: `{"a":3,"b":4}`, when: `(arr a b a)`, then: `[3,4,3]`},
		{give: `{"a":3,"b":4}`, when: `(arr (get a) (get b) (get a))`, then: `[3,4,3]`},
		{give: ``, when: `(arr)`, then: `[]`},
		// (object)
		{give: `{"a":"aaa","b":"bbb"}}`, when: `(obj (get a) (get b))`, then: `{"aaa":"bbb"}`},
		{give: `{"a":3,"b":4}`, when: `(obj "a b" a y b)`, then: `{"a b":3,"y":4}`},
		{give: `{"a":3,"b":4}`, when: `(obj "a b" (get a) y (get b))`, then: `{"a b":3,"y":4}`},
		{give: `{"a":3,"b":4}`, when: `(obj x (get a) y (get b))`, then: `{"x":3,"y":4}`},
		{give: `{"a":{"b":{"c":3}}}`, when: `(get a b (obj x c))`, then: `{"x":3}`},
		// (get) order should not matter.
		{give: `{"a":3,"b":4]}`, when: `(obj x b y a z a w b)`, then: `{"x":4,"y":3,"z":3,"w":4}`},
		// (get)
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
		// (root)
		{give: `3`, when: `(root)`, then: `3`},
		{give: ``, when: `(root)`, then: ``},
		// (raw)
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

func Test_Invalid_Query(t *testing.T) {

	tt := []struct {
		give string
		when string
		then string
	}{
		{give: `[3,4]`, when: `(collect (this)`, then: `[3,4]`},
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
		{give: `  {  "a"  :	3,  "b"  :  [  {  "c"  :  4  }  ]  }`, when: `(get b 0 c)`, then: `4`},
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
		give string
		when string
		then string
	}{
		{give: `{"a":2,"b":3}`, when: `a`, then: `2`},
		{give: `{"a":2,"b":3}`, when: `b`, then: `3`},
		{give: `{"a":2,"b":3}`, when: `c`, then: ``},
		{give: `{"a":[2,3]}`, when: `a`, then: `[2,3]`},
		{give: `[2,3]`, when: `0`, then: `2`},
		{give: `[2,3]`, when: `1`, then: `3`},
		{give: `[2,3]`, when: `2`, then: ``},
		{give: `{"a":{"b":2}}`, when: `a`, then: `{"b":2}`},
	}

	for _, tc := range tt {
		j := JSON(tc.give)
		r := j.Get(tc.when)
		assert.Equal(t, tc.then, r.String(), "TC: %v", tc)
	}
}

func BenchmarkJson_Get(b *testing.B) {
	j := JSON(TestData1)
	for i := 0; i < b.N; i++ {
		_ = j.Get("age")
	}
}

func TestJsonGet_Order(t *testing.T) {

	j := JSON(`{"a":3,"b":4}`)
	assert.Equal(t, "4", j.Get("b").String())
	assert.Equal(t, "3", j.Get("a").String())

	j = JSON(`[3,4]`)
	assert.Equal(t, "4", j.Get("1").String())
	assert.Equal(t, "3", j.Get("0").String())
}

func TestJsonForEachKeyVal(t *testing.T) {
	tt := []struct {
		give string
		then []string
	}{
		{give: `{}`, then: nil},
		{give: `{"a":2}`, then: []string{"a", "2"}},
		{give: `{"a":2,"b":3}`, then: []string{"a", "2", "b", "3"}},
		{give: `{"a":{"b":2}}`, then: []string{"a", `{"b":2}`}},
		{give: `{"a":[2]}`, then: []string{"a", "[2]"}},
	}
	for _, tc := range tt {
		var r []string
		j := JSON(tc.give)
		j.ForEachKeyVal(func(k, v Json) bool {
			r = append(r, k.Str(), v.String())
			return false
		})
		assert.Equal(t, tc.then, r, tc.give)
	}
}

func TestJsonForEach(t *testing.T) {
	tt := []struct {
		give string
		then []string
	}{
		{give: "[]", then: nil},
		{give: "[10]", then: []string{"0", "10"}},
		{give: "[10,20]", then: []string{"0", "10", "1", "20"}},
		{give: "[10,20,30]", then: []string{"0", "10", "1", "20", "2", "30"}},
		{give: "[{},{},[10]]", then: []string{"0", "{}", "1", "{}", "2", "[10]"}},
	}
	for _, tc := range tt {
		var r []string
		j := JSON(tc.give)
		j.ForEach(func(i, v Json) bool {
			r = append(r, i.String(), v.String())
			return false
		})
		assert.Equal(t, tc.then, r, tc.give)
	}
}

func TestJsonStr(t *testing.T) {
	tt := []struct {
		give string
		then string
	}{
		{give: ``, then: ``},
		{give: `""`, then: ``},
		{give: `"a"`, then: `a`},
		{give: `"\"a\""`, then: `"a"`},
	}
	for _, tc := range tt {
		j := JSON(tc.give)
		assert.Equal(t, tc.then, j.Str(), tc.give)
	}
}

func TestJsonInt(t *testing.T) {
	tt := []struct {
		give string
		then int
	}{
		{give: ``, then: 0},
		{give: `0`, then: 0},
		{give: `1`, then: 1},
		{give: `-2`, then: -2},
		{give: `100`, then: 100},
	}
	for _, tc := range tt {
		j := JSON(tc.give)
		assert.Equal(t, tc.then, j.Int(), tc.give)
	}
}

func TestJsonFloat(t *testing.T) {
	tt := []struct {
		give string
		then float64
	}{
		{give: ``, then: 0},
		{give: `0`, then: 0},
		{give: `1`, then: 1},
		{give: `-2.2`, then: -2.2},
		{give: `1e2`, then: 100},
	}
	for _, tc := range tt {
		j := JSON(tc.give)
		assert.Equal(t, tc.then, j.Float(), tc.give)
	}
}

func TestJsonBool(t *testing.T) {
	tt := []struct {
		give string
		when bool
	}{
		{give: ``, when: false},
		{give: `false`, when: false},
		{give: `true`, when: true},
	}
	for _, tc := range tt {
		j := JSON(tc.give)
		assert.Equal(t, tc.when, j.Bool(), tc.give)
	}
}

func TestJsonIsEmpty(t *testing.T) {
	tt := []struct {
		give string
		when bool
	}{
		{give: ``, when: false},
		{give: `{}`, when: true},
		{give: `[]`, when: true},
		{give: `""`, when: true},
		{give: `null`, when: false},
		{give: `0`, when: false},
		{give: `false`, when: false},
		{give: `true`, when: false},
	}
	for _, tc := range tt {
		j := JSON(tc.give)
		assert.Equal(t, tc.when, j.IsEmpty(), tc.give)
	}
}

func TestJsonIsEmptyString(t *testing.T) {
	tt := []struct {
		give string
		then bool
	}{
		{give: ``, then: false},
		{give: `""`, then: true},
		{give: `"a"`, then: false},
	}
	for _, tc := range tt {
		j := JSON(tc.give)
		assert.Equal(t, tc.then, j.IsEmptyString(), tc.give)
	}
}

func TestJsonIsEmptyObject(t *testing.T) {
	tt := []struct {
		give string
		then bool
	}{
		{give: ``, then: false},
		{give: `{}`, then: true},
		{give: `{   }`, then: true},
		{give: `[]`, then: false},
		{give: `""`, then: false},
		{give: `3`, then: false},
	}
	for _, tc := range tt {
		j := JSON(tc.give)
		assert.Equal(t, tc.then, j.IsEmptyObject(), tc.give)
	}
}

func TestJsonIsEmptyArray(t *testing.T) {
	tt := []struct {
		give string
		then bool
	}{
		{give: ``, then: false},
		{give: `[]`, then: true},
		{give: `[  ]`, then: true},
		{give: `{}`, then: false},
		{give: `""`, then: false},
		{give: `3`, then: false},
	}
	for _, tc := range tt {
		j := JSON(tc.give)
		assert.Equal(t, tc.then, j.IsEmptyArray(), tc.give)
	}
}

func TestJsonIsFalsy(t *testing.T) {
	tt := []struct {
		give string
		then bool
	}{
		{give: ``, then: false},
		{give: `[]`, then: true},
		{give: `{}`, then: true},
		{give: `""`, then: true},
		{give: `false`, then: true},
		{give: `null`, then: true},
		{give: `0`, then: true},
		//
		{give: `3`, then: false},
		{give: `[0]`, then: false},
		{give: `{"a":0}`, then: false},
		{give: `"a"`, then: false},
		{give: `true`, then: false},
	}
	for _, tc := range tt {
		j := JSON(tc.give)
		assert.Equal(t, tc.then, j.IsFalsy(), tc)
	}
}

func TestJsonIsTruthy(t *testing.T) {
	tt := []struct {
		give string
		then bool
	}{
		{give: ``, then: false},
		{give: `[]`, then: false},
		{give: `{}`, then: false},
		{give: `""`, then: false},
		{give: `false`, then: false},
		{give: `null`, then: false},
		{give: `0`, then: false},
		//
		{give: `3`, then: true},
		{give: `[0]`, then: true},
		{give: `{"a":0}`, then: true},
		{give: `"a"`, then: true},
		{give: `true`, then: true},
	}
	for _, tc := range tt {
		j := JSON(tc.give)
		assert.Equal(t, tc.then, j.IsTruthy(), tc)
	}
}

func TestJsonIsVoid(t *testing.T) {
	tt := []struct {
		give string
		then bool
	}{
		{give: ``, then: false},
		{give: `{}`, then: true},
		{give: `[]`, then: true},
		{give: `""`, then: false},
		{give: `null`, then: false},
		{give: `0`, then: false},
		{give: `false`, then: false},
		{give: `true`, then: false},
	}
	for _, tc := range tt {
		j := JSON(tc.give)
		assert.Equal(t, tc.then, j.IsVoid(), tc.give)
	}
}

func TestJsonIsBlank(t *testing.T) {
	tt := []struct {
		give string
		then bool
	}{
		{give: ``, then: false},
		{give: `{}`, then: true},
		{give: `[]`, then: true},
		{give: `""`, then: false},
		{give: `null`, then: true},
		{give: `0`, then: false},
		{give: `false`, then: false},
		{give: `true`, then: false},
	}
	for _, tc := range tt {
		j := JSON(tc.give)
		assert.Equal(t, tc.then, j.IsBlank(), tc.give)
	}
}

func TestJsonIsNully(t *testing.T) {
	tt := []struct {
		give string
		then bool
	}{
		{give: ``, then: false},
		{give: `{}`, then: true},
		{give: `[]`, then: true},
		{give: `""`, then: true},
		{give: `null`, then: true},
		{give: `0`, then: false},
		{give: `false`, then: false},
		{give: `true`, then: false},
	}
	for _, tc := range tt {
		j := JSON(tc.give)
		assert.Equal(t, tc.then, j.IsNully(), tc.give)
	}
}

func TestJsonIsSome(t *testing.T) {
	tt := []struct {
		give string
		then bool
	}{
		{give: ``, then: false},
		{give: `{}`, then: true},
		{give: `[]`, then: true},
		{give: `""`, then: true},
		{give: `null`, then: false},
		{give: `0`, then: true},
		{give: `false`, then: true},
		{give: `true`, then: true},
	}
	for _, tc := range tt {
		j := JSON(tc.give)
		assert.Equal(t, tc.then, j.IsSome(), tc)
	}
}

func TestJsonStringify(t *testing.T) {
	tt := []struct {
		give string
		then string
	}{
		{give: ``, then: `""`},
		{give: `{}`, then: `"{}"`},
		{give: `{ "hello": "wo\"rld" }`, then: `"{ \"hello\": \"wo\\\"rld\" }"`},
		{give: `[]`, then: `"[]"`},
		{give: `123`, then: `"123"`},
		{give: `null`, then: `"null"`},
		{give: `false`, then: `"false"`},
		{give: `true`, then: `"true"`},
		{give: `""`, then: `"\"\""`},
		{give: `"a b"`, then: `"\"a b\""`},
	}
	for _, tc := range tt {
		j := JSON(tc.give)
		assert.Equal(t, tc.then, j.Stringify().String(), tc)
	}
}

func BenchmarkJsonToStringify(b *testing.B) {
	j := JSON(`{ "hello": "wo\"rld" }`)
	for i := 0; i < b.N; i++ {
		_ = j.Stringify()
	}
}

func TestJsonJsonify(t *testing.T) {
	tt := []struct {
		give string
		then string
	}{
		{give: ``, then: ``},
		{give: `""`, then: ``},
		{give: `"{}"`, then: `{}`},
		{give: `"{ \"hello\": \"wo\\\"rld\" }"`, then: `{ "hello": "wo\"rld" }`},
		{give: `"[]"`, then: `[]`},
		{give: `"123"`, then: `123`},
		{give: `"null"`, then: `null`},
		{give: `"false"`, then: `false`},
		{give: `"true"`, then: `true`},
		{give: `"a b"`, then: `a b`}, // Invalid JSON.
	}
	for _, tc := range tt {
		j := JSON(tc.give)
		assert.Equal(t, tc.then, j.Jsonify().String(), tc)
	}
}

func BenchmarkJsonJsonify(b *testing.B) {
	j := JSON(`"{ \"hello\": \"wo\\\"rld\" }"`)
	for i := 0; i < b.N; i++ {
		_ = j.Jsonify()
	}
}

func Example_funcDebug() {

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

func ExampleJson_ForEachKeyVal() {
	m := func(k, v Json) bool {
		fmt.Println(k, v)
		return false
	}

	j := JSON(TestData1)
	j.ForEachKeyVal(m)

	// Output:
	// "name" "Mary"
	// "last" "Jane"
	// "token" null
	// "settings" {}
	// "posts" []
	// "address" {"city":"Place","country":"USA"}
	// "contacts" [{"name":"Karen"},{"name":"Michelle","last":"Jane"}]
	// "age" 33
	// "random" [3,null,{},[],"",false]
}

func BenchmarkJson_ForEachKeyVal(b *testing.B) {
	m := func(k, v Json) bool {
		return false
	}
	j := JSON(TestData1)
	for i := 0; i < b.N; i++ {
		j.ForEachKeyVal(m)
	}
}

func ExampleJson_ForEach() {
	m := func(k, v Json) bool {
		fmt.Println(k, v)
		return false
	}

	j := JSON(TestData2)
	j.ForEach(m)

	// Output:
	// 0 {"name":"Karen"}
	// 1 {"name":"Michelle","last":"Jane"}
}

func BenchmarkJson_ForEach(b *testing.B) {
	m := func(k, v Json) bool {
		return false
	}
	j := JSON(TestData2)
	for i := 0; i < b.N; i++ {
		j.ForEach(m)
	}
}

func ExampleJson_Iterate() {

	j := JSON(`{"a":1,"b":2,"c":{"a":3,"b":{"a":4,"b":[{"a":5},{"a":6,"b":7,"c":[8,9,0,{},[]]}]}},"d":1}`)
	v := j.Iterate(func(k, v Json) (Json, Json) {
		k = JSON(strings.ToUpper(k.String()))
		if v.IsNumber() {
			return k, v.Stringify()
		}
		return k, v
	})

	fmt.Println(v)

	// Output:
	// {"A":"1","B":"2","C":{"A":"3","B":{"A":"4","B":[{"A":"5"},{"A":"6","B":"7","C":["8","9","0",{},[]]}]}},"D":"1"}
}

func BenchmarkJson_Iterate(b *testing.B) {
	m := func(k, v Json) (Json, Json) {
		return k, v
	}
	j := JSON(TestData1)
	for i := 0; i < b.N; i++ {
		_ = j.Iterate(m)
	}
}

func Benchmark_QueryFunction_Iterate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Get(`{ "a": 3, "b": 4 }`, `(iterate 0 1)`)
	}
}

func ExampleJson_IterateKeysValues() {
	m := func(v Json) Json {
		fmt.Println(v)
		return v
	}

	j := JSON(TestData1)
	r := j.IterateKeysValues(m)

	fmt.Println(r)

	// Output:
	// "name"
	// "Mary"
	// "last"
	// "Jane"
	// "token"
	// null
	// "settings"
	// "posts"
	// "address"
	// "city"
	// "Place"
	// "country"
	// "USA"
	// "contacts"
	// "name"
	// "Karen"
	// "name"
	// "Michelle"
	// "last"
	// "Jane"
	// "age"
	// 33
	// "random"
	// 3
	// null
	// ""
	// false
	// {"name":"Mary","last":"Jane","token":null,"settings":{},"posts":[],"address":{"city":"Place","country":"USA"},"contacts":[{"name":"Karen"},{"name":"Michelle","last":"Jane"}],"age":33,"random":[3,null,{},[],"",false]}
}

func BenchmarkJson_IterateKeysValues(b *testing.B) {
	m := func(v Json) Json {
		return v
	}
	j := JSON(TestData1)
	for i := 0; i < b.N; i++ {
		_ = j.IterateKeysValues(m)
	}
}

func Benchmark_QueryFunction_IterateKV(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Get(`{ "a": 3, "b": 4 }`, `(iterate-kv (this) (this))`)
	}
}

func ExampleJson_IterateKeys() {
	m := func(v Json) Json {
		fmt.Println(v)
		return v
	}

	j := JSON(TestData1)
	r := j.IterateKeys(m)

	fmt.Println(r)

	// Output:
	// "name"
	// "last"
	// "token"
	// "settings"
	// "posts"
	// "address"
	// "city"
	// "country"
	// "contacts"
	// "name"
	// "name"
	// "last"
	// "age"
	// "random"
	// {"name":"Mary","last":"Jane","token":null,"settings":{},"posts":[],"address":{"city":"Place","country":"USA"},"contacts":[{"name":"Karen"},{"name":"Michelle","last":"Jane"}],"age":33,"random":[3,null,{},[],"",false]}
}

func BenchmarkJson_IterateKeys(b *testing.B) {
	m := func(v Json) Json {
		return v
	}
	j := JSON(TestData1)
	for i := 0; i < b.N; i++ {
		_ = j.IterateKeys(m)
	}
}

func Benchmark_QueryFunction_IterateK(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Get(`{ "a": 3, "b": 4 }`, `(iterate-k (this))`)
	}
}

func ExampleJson_IterateValues() {
	m := func(v Json) Json {
		fmt.Println(v)
		return v
	}

	j := JSON(TestData1)
	r := j.IterateValues(m)

	fmt.Println(r)

	// Output:
	// "Mary"
	// "Jane"
	// null
	// "Place"
	// "USA"
	// "Karen"
	// "Michelle"
	// "Jane"
	// 33
	// 3
	// null
	// ""
	// false
	// {"name":"Mary","last":"Jane","token":null,"settings":{},"posts":[],"address":{"city":"Place","country":"USA"},"contacts":[{"name":"Karen"},{"name":"Michelle","last":"Jane"}],"age":33,"random":[3,null,{},[],"",false]}
}

func BenchmarkJson_IterateValues(b *testing.B) {
	m := func(v Json) Json {
		return v
	}
	j := JSON(TestData1)
	for i := 0; i < b.N; i++ {
		_ = j.IterateValues(m)
	}
}

func Benchmark_QueryFunction_IterateV(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Get(`{ "a": 3, "b": 4 }`, `(iterate-v (this))`)
	}
}

func ExampleJson_Iterator() {

	m := func(o *strings.Builder, k, v Json) {
		if k.String() == "" {
			if !v.IsObject() && !v.IsArray() {
				o.WriteString(v.String())
			}
			return
		}
		o.WriteString(k.String())
		o.WriteString(":")
		if !v.IsObject() && !v.IsArray() {
			fmt.Println(k, v)
			o.WriteString(v.String())
		}
	}

	var o strings.Builder

	j := JSON(TestData1)
	j.Iterator(&o, JSON(""), m)

	fmt.Println(o.String())

	// Output:
	// "name" "Mary"
	// "last" "Jane"
	// "token" null
	// "city" "Place"
	// "country" "USA"
	// "name" "Karen"
	// "name" "Michelle"
	// "last" "Jane"
	// "age" 33
	// {"name":"Mary","last":"Jane","token":null,"settings":{},"posts":[],"address":{"city":"Place","country":"USA"},"contacts":[{"name":"Karen"},{"name":"Michelle","last":"Jane"}],"age":33,"random":[3,null,{},[],"",false]}
}

func BenchmarkJson_Iterator(b *testing.B) {
	m := func(o *strings.Builder, k, v Json) {
		o.WriteString(k.String())
		o.WriteString(":")
		if !v.IsObject() && !v.IsArray() {
			o.WriteString(v.String())
		}
	}
	var o strings.Builder
	o.Grow(2550)
	j := JSON(TestData1)
	for i := 0; i < b.N; i++ {
		o.Reset()
		j.Iterator(&o, JSON(""), m)
	}
}

func ExampleJson_Uglify() {

	j := JSON(`{ "a": "b" }`)
	j = j.Uglify()

	fmt.Println(j)

	// Output:
	// {"a":"b"}
}

func BenchmarkJson_Uglify(b *testing.B) {
	j := JSON(TestData1)
	for i := 0; i < b.N; i++ {
		_ = j.Uglify()
	}
}

func ExampleJson_Prettify() {

	j := JSON(`{ "a": "b" }`)
	j = j.Prettify()

	fmt.Println(j)

	// Output:
	// {
	//     "a": "b"
	// }
}

func BenchmarkJson_Prettify(b *testing.B) {
	j := JSON(TestData1)
	for i := 0; i < b.N; i++ {
		_ = j.Prettify()
	}
}

func Benchmark_QueryFunction_Sort(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Get(`[{ "a": 5 }, { "a": 4 }, { "a": 3 }]`, `(sort asc a)`)
	}
}

const TestData1 = `{"name":"Mary","last":"Jane","token":null,"settings":{},"posts":[],"address":{"city":"Place","country":"USA"},"contacts":[{"name":"Karen"},{"name":"Michelle","last":"Jane"}],"age":33,"random":[3,null,{},[],"",false]}`
const TestData2 = `[{"name":"Karen"},{"name":"Michelle","last":"Jane"}]`
