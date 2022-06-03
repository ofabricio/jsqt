package jsqt

import (
	"strconv"
	"strings"
	"unicode"

	. "github.com/ofabricio/scanner"
)

func Get(jsn, qry string) Json {
	src := New(jsn)
	q := query{Scanner(qry)}
	return New(q.Parse(src))
}

func New(jsn string) Json {
	return Json{Scanner(jsn)}
}

// #region Query

type query struct {
	Scanner
}

func (q *query) Parse(j Json) string {
	if q.Match(".") {
		return q.Parse(j)
	}
	if q.Match("*") {
		return q.ParseStar(j)
	}
	if q.Match("{") {
		return q.ParseObject(j)
	}
	if seg := q.TokenFor(q.MatchObjectKey); seg != "" {
		return q.Parse(j.Get(seg))
	}
	return j.String()
}

func (q *query) ParseObject(j Json) string {
	var obj strings.Builder
	obj.WriteString("{")
	for q.Match(",") || !q.Match("}") {
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

func (q *query) ParseObjectKey() string {
	m := q.Mark()
	if key := q.TokenFor(q.MatchObjectKey); q.Match(":") {
		return key
	}
	q.Back(m)
	return ""
}

func (q *query) GetLastPathSegment() string {
	m := q.Mark()
	s := q.TokenFor(q.MatchObjectKey)
	for q.Match(".") {
		s = q.TokenFor(q.MatchObjectKey)
	}
	q.Back(m)
	return s
}

func (q *query) ParseStar(j Json) string {
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

func (q *query) MatchObjectKey() bool {
	return q.MatchWhileBy(q.IsObjectKey)
}

func (q *query) IsObjectKey(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_'
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
	if j.Match("{") {
		for !j.Match("}") {
			k, _ := j.TokenFor(j.MatchString), j.Match(":")
			v, _ := j.GetValue(), j.Match(",")
			if f(strings.Trim(k, `"`), New(v)) {
				return
			}
		}
	}
}

func (j *Json) IterateArray(f func(string, Json) bool) {
	if j.Match("[") {
		for i := 0; !j.Match("]"); i++ {
			k := strconv.Itoa(i)
			v, _ := j.GetValue(), j.Match(",")
			if f(k, New(v)) {
				return
			}
		}
	}
}

func (j *Json) GetValue() string {
	if j.Equal("{") {
		return j.TokenFor(j.MatchObject)
	}
	if j.Equal("[") {
		return j.TokenFor(j.MatchArray)
	}
	if j.Equal(`"`) {
		return j.TokenFor(j.MatchString)
	}
	return j.TokenFor(j.MatchRest)
}

func (j *Json) MatchRest() bool {
	for !j.Equal(",") && !j.Equal("}") && !j.Equal("]") {
		j.Next()
	}
	return true
}

func (j *Json) MatchObject() bool {
	return j.MatchScope("{", "}")
}

func (j *Json) MatchArray() bool {
	return j.MatchScope("[", "]")
}

func (j *Json) MatchScope(open, clos string) bool {
	if j.Match(open) {
		c := 1
		for j.More() && c > 0 {
			if j.Match(open) {
				c++
				continue
			}
			if j.Match(clos) {
				c--
				continue
			}
			if j.MatchString() {
				continue
			}
			j.Next()
		}
		return c == 0
	}
	return false
}

func (j *Json) MatchString() bool {
	return j.Match(`"`) && j.MatchUntilEscape(`"`, `\`) && j.Match(`"`)
}

// #endregion Json
