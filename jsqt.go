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
		if v := q.TokenFor(func() bool {
			return q.UtilMatchString('"')
		}); v != "" {
			return v
		}
		return q.TokenAnything()
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
	if fname == "raw" {
		return q.ParseFuncArg(j)
	}
	if fname == "root" {
		return q.Root.String()
	}
	if fname == "." {
		return j.String()
	}
	if fname == "get" {
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
	if fname == "obj" {
		var out strings.Builder
		out.WriteString("{")
		for !q.EqualByte(')') {
			if out.Len() > 1 {
				out.WriteString(",")
			}
			k := q.ParseFuncArg(j)
			v := q.ParseFuncArg(j)
			if v != "" {
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
	if fname == "arr" {
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
	if fname == "flatten" {
		v := q.ParseFuncArg(j)
		return v[1 : len(v)-1]
	}
	if fname == "collect" {
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
	if fname == "join" {
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
	if fname == "size" {
		c := 0
		j = New(q.ParseFuncArg(j)) // Input.
		j.ForEach(func(i string, v Json) bool {
			c++
			return false
		})
		return strconv.Itoa(c)
	}
	if fname == "merge" {
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
	if fname == "default" {
		v := q.ParseFuncArg(j) // Input.
		d := q.ParseFuncArg(j) // Value.
		if v == "" {
			return d
		}
		return v
	}
	if fname == "omitempty" {
		v := q.ParseFuncArg(j)
		if v == "{}" || v == "[]" {
			return ""
		}
		return v
	}
	return ""
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
	if v := j.TokenFor(j.MatchObject); v != "" {
		return v
	}
	if v := j.TokenFor(j.MatchArray); v != "" {
		return v
	}
	if v := j.TokenFor(j.MatchString); v != "" {
		return v
	}
	return j.TokenFor(j.MatchRest)
}

func (j *Json) MatchRest() bool {
	return j.MatchUntilAnyByte3(',', '}', ']')
}

func (j *Json) MatchObject() bool {
	return j.UtilMatchOpenCloseCount('{', '}', '"')
}

func (j *Json) MatchArray() bool {
	return j.UtilMatchOpenCloseCount('[', ']', '"')
}

func (j *Json) MatchString() bool {
	return j.UtilMatchString('"')
}

// #endregion Json

// #region Functions

// #endregion Functions
