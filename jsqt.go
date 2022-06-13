package jsqt

import (
	"strconv"
	"strings"

	. "github.com/ofabricio/scanner"
)

func Get(jsn, qry string) Json {
	src := New(jsn)
	q := query{Scanner: Scanner(qry), Root: src}
	return New(q.Parse(src))
}

func New(jsn string) Json {
	return Json{Scanner(jsn)}
}

// #region Query

type query struct {
	Scanner
	Root Json
}

func (q *query) Parse(j Json) string {
	if v := q.ParseFunc(j); v != "" {
		return v
	}
	return ""
}

func (q *query) ParseFuncArg(j Json) string {
	if q.MatchByte(' ') {
		if v := q.ParseFunc(j); v != "" {
			return v
		}
		m := q.Mark()
		if q.UtilMatchString('"') {
			return q.Token(m)
		}
		if q.MatchUntilAnyByte(' ', ')') {
			return q.Token(m)
		}
	}
	return ""
}

func (q *query) ParseFunc(j Json) string {
	if q.MatchByte('(') {
		if q.MatchByte(')') {
			return j.String()
		}
		if fname := q.TokenAnything(); fname != "" {
			v, _ := q.CallFunc(fname, j), q.MatchByte(')')
			return v
		}
	}
	return ""
}

func (q *query) TokenAnything() string {
	return q.TokenFor(func() bool {
		return q.MatchUntilAnyByte(' ', ')')
	})
}

func (q *query) CallFunc(fname string, j Json) string {
	switch fname {
	case "get":
		return FuncGet(q, j)
	case "obj":
		return FuncObj(q, j)
	case "arr":
		return FuncArr(q, j)
	case "collect":
		return FuncCollect(q, j)
	case "flatten":
		return FuncFlatten(q, j)
	case "raw":
		return FuncRaw(q, j)
	case ".":
		return FuncCurrent(q, j)
	case "join":
		return FuncJoin(q, j)
	case "size":
		return FuncSize(q, j)
	case "merge":
		return FuncMerge(q, j)
	case "default":
		return FuncDefault(q, j)
	case "omitempty":
		return FuncOmitEmpty(q, j)
	case "root":
		return FuncRoot(q, j)
	default:
		return ""
	}
}

func (q *query) ForEach(j Json, f func(sub *query, item Json)) {
	ini := *q
	j.ForEach(func(i string, item Json) bool {
		end := ini
		f(&end, item)
		*q = end
		return false
	})
}

// #endregion Query

// #region Json

type Json struct {
	Scanner
}

func (j Json) String() string {
	return j.Scanner.String()
}

// Iterate iterates over a valid Json.
func (j *Json) Iterate(f func(string, Json) bool) {
	j.iterate(".", f)
}

func (j *Json) iterate(path string, f func(string, Json) bool) {
	f(path, *j)
	if path != "." {
		path += "."
	}
	fn := func(k string, v Json) bool {
		v.iterate(path+k, f)
		return false
	}
	j.ForEachKeyVal(fn)
	j.ForEach(fn)
}

func (j *Json) Get(keyOrIndex string) (r Json) {
	f := func(k string, v Json) bool {
		if k == keyOrIndex {
			r = v
			return true
		}
		return false
	}
	j.ForEachKeyVal(f)
	j.ForEach(f)
	return r
}

func (j *Json) Collect(keyOrIndex string) (r Json) {
	if j.IsArray() && !(keyOrIndex[0] >= '0' && keyOrIndex[0] <= '9') {
		var out strings.Builder
		out.WriteString("[")
		j.ForEach(func(i string, v Json) bool {
			if s := v.Collect(keyOrIndex).String(); s != "" {
				if out.Len() > 1 {
					out.WriteString(",")
				}
				out.WriteString(s)
			}
			return false
		})
		out.WriteString("]")
		return New(out.String())
	}
	return j.Get(keyOrIndex)
}

func (j *Json) IsObject() bool {
	return j.EqualByte('{')
}

func (j *Json) IsArray() bool {
	return j.EqualByte('[')
}

func (j *Json) ForEachKeyVal(f func(string, Json) bool) {
	if j.MatchByte('{') {
		for !j.MatchByte('}') {
			k, _ := j.TokenFor(j.MatchString), j.MatchByte(':')
			v, _ := j.GetValue(), j.MatchByte(',')
			if f(strings.Trim(k, `"`), New(v)) {
				return
			}
		}
	}
}

func (j *Json) ForEach(f func(string, Json) bool) {
	if j.MatchByte('[') {
		for i := 0; !j.MatchByte(']'); i++ {
			k := strconv.Itoa(i)
			v, _ := j.GetValue(), j.MatchByte(',')
			if f(k, New(v)) {
				return
			}
		}
	}
}

func (j *Json) GetValue() string {
	m := j.Mark()
	if j.UtilMatchOpenCloseCount('{', '}', '"') { // Match Object.
		return j.Token(m)
	}
	if j.UtilMatchOpenCloseCount('[', ']', '"') { // Match Array.
		return j.Token(m)
	}
	if j.MatchString() {
		return j.Token(m)
	}
	if j.MatchUntilAnyByte3(',', '}', ']') { // Match Anything.
		return j.Token(m)
	}
	return ""
}

func (j *Json) MatchString() bool {
	return j.UtilMatchString('"')
}

// #endregion Json

// #region Functions

func FuncRaw(q *query, j Json) string {
	return q.ParseFuncArg(j)
}

func FuncRoot(q *query, j Json) string {
	return q.Root.String()
}

func FuncCurrent(q *query, j Json) string {
	return j.String()
}

func FuncGet(q *query, j Json) string {
	for {
		key := q.ParseFuncArg(j)
		if key == "" {
			return j.String()
		}
		if key[0] == '"' {
			key = key[1 : len(key)-1]
		}
		j = j.Collect(key)
	}
}

func FuncObj(q *query, j Json) string {
	var out strings.Builder
	out.WriteString("{")
	for !q.EqualByte(')') {
		if k, v := q.ParseFuncArg(j), q.ParseFuncArg(j); v != "" {
			if out.Len() > 1 {
				out.WriteString(",")
			}
			if k[0] == '"' {
				out.WriteString(k)
			} else {
				out.WriteString(`"`)
				out.WriteString(k)
				out.WriteString(`"`)
			}
			out.WriteString(`:`)
			out.WriteString(v)
		}
	}
	out.WriteString("}")
	return out.String()
}

func FuncArr(q *query, j Json) string {
	var out strings.Builder
	out.WriteString("[")
	for !q.EqualByte(')') {
		if out.Len() > 1 {
			out.WriteString(",")
		}
		v := q.ParseFuncArg(j)
		out.WriteString(v)
	}
	out.WriteString("]")
	return out.String()
}

func FuncFlatten(q *query, j Json) string {
	v := q.ParseFuncArg(j)
	return v[1 : len(v)-1]
}

func FuncCollect(q *query, j Json) string {
	var out strings.Builder
	out.WriteString("[")
	j = New(q.ParseFuncArg(j)) // Input.
	q.ForEach(j, func(sub *query, item Json) {
		for !sub.EqualByte(')') {
			item = New(sub.ParseFuncArg(item))
		}
		if item.String() != "" {
			if out.Len() > 1 {
				out.WriteString(",")
			}
			out.WriteString(item.String())
		}
	})
	out.WriteString("]")
	return out.String()
}

func FuncJoin(q *query, j Json) string {
	var out strings.Builder
	out.WriteString("{")
	j = New(q.ParseFuncArg(j)) // Input.
	q.ForEach(j, func(sub *query, item Json) {
		if f := sub.ParseFuncArg(item); f == "" { // Filter.
			return
		}
		for !sub.EqualByte(')') {
			if out.Len() > 1 {
				out.WriteString(",")
			}
			k := sub.ParseFuncArg(item) // Key.
			v := sub.ParseFuncArg(item) // Value.
			out.WriteString(k)
			out.WriteString(":")
			out.WriteString(v)
		}
	})
	out.WriteString("}")
	return out.String()
}

func FuncSize(q *query, j Json) string {
	c := 0
	j = New(q.ParseFuncArg(j)) // Input.
	j.ForEach(func(i string, v Json) bool {
		c++
		return false
	})
	return strconv.Itoa(c)
}

func FuncMerge(q *query, j Json) string {
	done := make(map[string]bool)
	var b strings.Builder
	b.WriteString("{")
	j = New(q.ParseFuncArg(j)) // Input.
	j.ForEach(func(i string, v Json) bool {
		v.ForEachKeyVal(func(k string, v Json) bool {
			if !done[k] {
				if b.Len() > 1 {
					b.WriteString(",")
				}
				b.WriteString(`"`)
				b.WriteString(k)
				b.WriteString(`":`)
				b.WriteString(v.String())
			}
			done[k] = true
			return false
		})
		return false
	})
	b.WriteString("}")
	return b.String()
}

func FuncDefault(q *query, j Json) string {
	v := q.ParseFuncArg(j) // Input.
	d := q.ParseFuncArg(j) // Value.
	if v == "" {
		return d
	}
	return v
}

func FuncOmitEmpty(q *query, j Json) string {
	v := q.ParseFuncArg(j)
	if v == "{}" || v == "[]" {
		return ""
	}
	return v
}

// #endregion Functions
