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
	q.ByteMatch('.')
	if q.ByteMatch('@') && q.Match("root") {
		return q.Parse(q.Root)
	}
	if path := q.MatchPath(); path != "" {
		return q.Parse(j.Get(path))
	}
	if name, args := q.MatchFunc(); name != "" {
		return q.CallFunc(name, args, j)
	}
	if obj := q.ParseObject(j); obj != "" {
		return obj
	}
	if arr := q.ParseArray(j); arr != "" {
		return arr
	}
	return j.String()
}

func (q *query) MatchPath() string {
	if path := q.ByteTokenBy(q.IsObjectKey); path != "" {
		return path
	}
	if m := q.Mark(); q.ByteMatchString('"') {
		str := q.Token(m)
		return str[1 : len(str)-1] // Remove "".
	}
	return ""
}

func (q *query) MatchFunc() (string, string) {
	if q.ByteMatch('(') {
		name, _ := q.ByteTokenBy(q.IsFuncName), q.ByteMatch(' ')
		args, _ := q.ByteTokenBy(q.IsFuncName), q.ByteMatch(')')
		return name, args
	}
	return "", ""
}

func (q *query) CallFunc(name, args string, j Json) string {
	// TODO: add global function map.
	if name == "flatten" {
		return flatten(q, j, args)
	}
	if name == "size" {
		return size(q, j, args)
	}
	if name == "omitempty" {
		return omitempty(q, j, args)
	}
	if name == "default" {
		return defaultValue(q, j, args)
	}
	if name == "collect" {
		return collect(q, j, args)
	}
	if name == "merge" {
		return merge(q, j, args)
	}
	return j.String()
}

func (q *query) ParseObject(j Json) string {
	if q.ByteMatch('{') {
		var obj strings.Builder
		obj.WriteString("{")
		for q.ByteMatch(',') || !q.ByteMatch('}') {
			key := q.ParseObjectKey()
			if key == "" {
				key = q.GetLastPathSegment()
			}
			if v := q.Parse(j); v != "" {
				if obj.Len() > 1 {
					obj.WriteString(",")
				}
				obj.WriteString(`"`)
				obj.WriteString(key)
				obj.WriteString(`":`)
				obj.WriteString(v)
			}
		}
		obj.WriteString("}")
		return obj.String()
	}
	return ""
}

func (q *query) ParseObjectKey() string {
	m := q.Mark()
	if key := q.MatchPath(); q.ByteMatch(':') {
		return key
	}
	q.Back(m)
	return ""
}

func (q *query) GetLastPathSegment() string {
	m := q.Mark()
	key := q.MatchPath()
	for q.ByteMatch('.') {
		if k := q.MatchPath(); k != "" {
			key = k
		}
	}
	q.Back(m)
	return key
}

func (q *query) ParseArray(j Json) string {
	if q.ByteMatch('[') {
		var obj strings.Builder
		obj.WriteString("[")
		for q.ByteMatch(',') || !q.ByteMatch(']') {
			if v := q.Parse(j); v != "" {
				if obj.Len() > 1 {
					obj.WriteString(",")
				}
				obj.WriteString(v)
			}
		}
		obj.WriteString("]")
		return obj.String()
	}
	return ""
}

func (q *query) IterateArray(j Json) string {
	var arr strings.Builder
	arr.WriteString("[")
	end := *q
	j.IterateArray(func(i string, v Json) bool {
		sub := *q
		if s := sub.Parse(v); s != "" {
			end = sub
			if arr.Len() > 1 {
				arr.WriteString(",")
			}
			arr.WriteString(s)
		}
		return false
	})
	*q = end
	arr.WriteString("]")
	return arr.String()
}

func (q *query) IsObjectKey(r byte) bool {
	return r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || r == '_'
}

func (q *query) IsFuncName(r byte) bool {
	return r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || r == '_'
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
	j.IterateObject(fn)
	j.IterateArray(fn)
}

func (j *Json) Get(keyOrIndex string) (r Json) {
	f := func(k string, v Json) bool {
		if k == keyOrIndex {
			r = v
			return true
		}
		return false
	}
	j.IterateObject(f)
	j.IterateArray(f)
	return r
}

func (j *Json) IterateObject(f func(string, Json) bool) {
	if j.ByteMatch('{') {
		for !j.ByteMatch('}') {
			k, _ := j.TokenFor(j.MatchString), j.ByteMatch(':')
			v, _ := j.GetValue(), j.ByteMatch(',')
			if f(strings.Trim(k, `"`), New(v)) {
				return
			}
		}
	}
}

func (j *Json) IterateArray(f func(string, Json) bool) {
	if j.ByteMatch('[') {
		for i := 0; !j.ByteMatch(']'); i++ {
			k := strconv.Itoa(i)
			v, _ := j.GetValue(), j.ByteMatch(',')
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
	return j.ByteMatchUntil(',', '}', ']')
}

func (j *Json) MatchObject() bool {
	return j.MatchCounting('{', '}')
}

func (j *Json) MatchArray() bool {
	return j.MatchCounting('[', ']')
}

func (j *Json) MatchString() bool {
	return j.ByteMatchString('"')
}

// #endregion Json

// #region Functions

func defaultValue(q *query, j Json, arg string) string {
	v := q.Parse(j)
	if v == "" {
		return arg
	}
	return v
}

func flatten(q *query, j Json, arg string) string {
	v := q.IterateArray(j)
	v = strings.TrimPrefix(v, "[")
	v = strings.TrimSuffix(v, "]")
	return v
}

func size(q *query, j Json, arg string) string {
	c := 0
	j = New(q.IterateArray(j))
	j.IterateArray(func(i string, v Json) bool {
		c++
		return false
	})
	return strconv.Itoa(c)
}

func omitempty(q *query, j Json, arg string) string {
	v := q.Parse(j)
	if v == "{}" {
		return ""
	}
	if v == "[]" {
		return ""
	}
	return j.String()
}

func collect(q *query, j Json, arg string) string {
	return q.IterateArray(j)
}

func merge(q *query, j Json, arg string) string {
	var b strings.Builder
	b.WriteString("{")
	done := make(map[string]bool)
	j = New(q.IterateArray(j))
	j.IterateArray(func(i string, v Json) bool {
		v.IterateObject(func(k string, v Json) bool {
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

// #endregion Functions
