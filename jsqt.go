package jsqt

import (
	"strconv"
	"strings"

	. "github.com/ofabricio/scanner"
)

func Get(jsn, qry string) Json {
	src := New(jsn)
	src.WS()
	q := Query{Scanner: Scanner(qry), Root: src}
	return q.Parse(src)
}

func New(jsn string) Json {
	return Json{Scanner(jsn)}
}

// #region Query

type Query struct {
	Scanner
	Root Json
}

func (q *Query) Parse(j Json) Json {
	return q.ParseFunc(j)
}

func (q *Query) ParseFunc(j Json) Json {
	if q.MatchByte('(') {
		if fname := q.TokenAnything(); fname != "" {
			v, _ := q.CallFunc(fname, j), q.MatchByte(')')
			return v
		}
		if q.MatchByte(')') {
			return j
		}
	}
	return New("")
}

func (q *Query) ParseArgFunOrKey(j Json) Json {
	if q.MatchByte(' ') {
		if q.EqualByte('(') {
			if v := q.ParseFunc(j); v.String() != "" {
				return v
			}
			return New("")
		}
		return q.ParseArgKey(j)
	}
	return New("")
}

func (q *Query) ParseArgFunOrRaw(j Json) Json {
	if q.MatchByte(' ') {
		if q.EqualByte('(') {
			if v := q.ParseFunc(j); v.String() != "" {
				return v
			}
			return New("")
		}
		return q.ParseArgRaw(j)
	}
	return New("")
}

func (q *Query) ParseArgKey(j Json) Json {
	q.MatchByte(' ')
	key := ""
	m := q.Mark()
	if q.UtilMatchString('"') {
		key = q.Token(m)
		key = key[1 : len(key)-1]
	} else if q.MatchUntilAnyByte(' ', ')') { // Anything.
		key = q.Token(m)
	}
	return j.Get(key)
}

func (q *Query) ParseArgRaw(j Json) Json {
	q.MatchByte(' ')
	v := ""
	m := q.Mark()
	if q.UtilMatchString('"') {
		v = q.Token(m)
	} else if q.MatchUntilAnyByte(' ', ')') { // Anything.
		v = q.Token(m)
	}
	return New(v)
}

func (q *Query) CallFunc(fname string, j Json) Json {
	switch fname {
	case "get":
		return FuncGet(q, j)
	case "obj":
		return FuncObj(q, j)
	case "arr":
		return FuncArr(q, j)
	case "raw":
		return FuncRaw(q, j)
	case "collect":
		return FuncCollect(q, j)
	case "flatten":
		return FuncFlatten(q, j)
	case "size":
		return FuncSize(q, j)
	case "default":
		return FuncDefault(q, j)
	case "omitempty":
		return FuncOmitempty(q, j)
	case "merge":
		return FuncMerge(q, j)
	case "join":
		return FuncJoin(q, j)
	case "iterate":
		return FuncIterate(q, j)
	case "==":
		return FuncEq(q, j)
	case "root":
		return q.Root
	case ".":
		return j
	default:
		return New("")
	}
}

func (q *Query) TokenAnything() string {
	return q.TokenFor(func() bool {
		return q.MatchUntilAnyByte(' ', ')')
	})
}

func (q *Query) ForEach(j Json, f func(sub *Query, item Json)) {
	ini := *q
	j.ForEach(func(i string, item Json) bool {
		end := ini
		f(&end, item)
		*q = end
		return false
	})
}

// #region Functions

func FuncRaw(q *Query, j Json) Json {
	return q.ParseArgRaw(j)
}

func FuncGet(q *Query, j Json) Json {
	for !q.EqualByte(')') {
		j = q.ParseArgFunOrKey(j)
	}
	return j
}

func FuncArr(q *Query, j Json) Json {
	var o strings.Builder
	o.WriteString("[")
	for !q.EqualByte(')') {
		if o.Len() > 1 {
			o.WriteString(",")
		}
		v := q.ParseArgFunOrKey(j)
		o.WriteString(v.String())
	}
	o.WriteString("]")
	return New(o.String())
}

func FuncObj(q *Query, j Json) Json {
	var o strings.Builder
	o.WriteString("{")
	for !q.EqualByte(')') {
		if k, v := q.ParseArgFunOrRaw(j), q.ParseArgFunOrKey(j); v.String() != "" {
			if o.Len() > 1 {
				o.WriteString(",")
			}
			if k.String()[0] == '"' {
				o.WriteString(k.String())
			} else {
				o.WriteString(`"`)
				o.WriteString(k.String())
				o.WriteString(`"`)
			}
			o.WriteString(`:`)
			o.WriteString(v.String())
		}
	}
	o.WriteString("}")
	return New(o.String())
}

func FuncCollect(q *Query, j Json) Json {
	var o strings.Builder
	o.WriteString("[")
	for !q.EqualByte(')') {
		if j.IsArray() {
			q.ForEach(j, func(sub *Query, item Json) {
				for !sub.EqualByte(')') {
					item = sub.ParseArgFunOrKey(item)
				}
				if item.String() != "" {
					if o.Len() > 1 {
						o.WriteString(",")
					}
					o.WriteString(item.String())
				}
			})
		} else {
			j = q.ParseArgKey(j)
		}
	}
	o.WriteString("]")
	return New(o.String())
}

func FuncDefault(q *Query, j Json) Json {
	d := q.ParseArgRaw(j) // Default value.
	if j.String() == "" {
		return d
	}
	return j
}

func FuncOmitempty(q *Query, j Json) Json {
	if j.String() == "{}" || j.String() == "[]" {
		return New("")
	}
	return j
}

func FuncSize(q *Query, j Json) Json {
	c := 0
	j.ForEach(func(i string, v Json) bool {
		c++
		return false
	})
	return New(strconv.Itoa(c))
}

func FuncMerge(q *Query, j Json) Json {
	done := make(map[string]bool)
	var b strings.Builder
	b.WriteString("{")
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
	return New(b.String())
}

func FuncJoin(q *Query, j Json) Json {
	var o strings.Builder
	o.WriteString("{")
	q.ForEach(j, func(sub *Query, item Json) {
		for !sub.EqualByte(')') {
			if o.Len() > 1 {
				o.WriteString(",")
			}
			k := sub.ParseArgFunOrKey(item) // Key.
			v := sub.ParseArgFunOrKey(item) // Value.
			o.WriteString(k.String())
			o.WriteString(":")
			o.WriteString(v.String())
		}
	})
	o.WriteString("}")
	return New(o.String())
}

func FuncIterate(q *Query, j Json) Json {
	m := q.TokenAnything() // Map function.
	_ = m
	// TODO: create functions map.
	return New(j.Iterate(num2str))
}

func FuncEq(q *Query, j Json) Json {
	f := q.ParseArgFunOrKey(j) // Field.
	v := q.ParseArgFunOrRaw(j) // Value.
	if f.String() == v.String() {
		return j
	}
	return New("")
}

func FuncFlatten(q *Query, j Json) Json {
	v := j.String()
	return New(v[1 : len(v)-1])
}

// #endregion Functions

// #endregion Query

// #region Json

type Json struct {
	Scanner
}

func (j Json) String() string {
	return j.Scanner.String()
}

func (j *Json) IsObject() bool {
	return j.EqualByte('{')
}

func (j *Json) IsArray() bool {
	return j.EqualByte('[')
}

func (j *Json) IsNumber() bool {
	return j.EqualByteRange('0', '9') || j.EqualByte('-')
}

func (j *Json) IsString() bool {
	return j.EqualByte('"')
}

func (j *Json) IsBool() bool {
	return j.EqualByte('t') && j.EqualByte('f')
}

func (j *Json) IsNull() bool {
	return j.EqualByte('n')
}

// Iterate iterates over a valid Json.
func (j *Json) Iterate(m func(string, Json) (string, string)) string {
	var o strings.Builder
	if j.IsObject() {
		o.WriteString("{")
		j.ForEachKeyVal(func(k string, v Json) bool {
			if o.Len() > 1 {
				o.WriteString(",")
			}
			nk, _ := m(k, v)
			o.WriteString(`"`)
			o.WriteString(nk)
			o.WriteString(`":`)
			o.WriteString(v.Iterate(m))
			return false
		})
		o.WriteString("}")
		return o.String()
	}
	if j.IsArray() {
		o.WriteString("[")
		j.ForEach(func(i string, v Json) bool {
			if o.Len() > 1 {
				o.WriteString(",")
			}
			o.WriteString(v.Iterate(m))
			return false
		})
		o.WriteString("]")
		return o.String()
	}
	_, nv := m("", *j)
	return nv
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

func (j *Json) Collect(keyOrIndex string) Json {
	if j.IsArray() && !(keyOrIndex[0] >= '0' && keyOrIndex[0] <= '9') {
		var o strings.Builder
		o.WriteString("[")
		j.ForEach(func(i string, v Json) bool {
			if s := v.Collect(keyOrIndex).String(); s != "" {
				if o.Len() > 1 {
					o.WriteString(",")
				}
				o.WriteString(s)
			}
			return false
		})
		o.WriteString("]")
		return New(o.String())
	}
	return j.Get(keyOrIndex)
}

func (j *Json) ForEachKeyVal(f func(string, Json) bool) {
	if j.MatchByte('{') {
		for j.WS() && !j.MatchByte('}') {
			k, _, _, _ := j.TokenFor(j.MatchString), j.WS(), j.MatchByte(':'), j.WS()
			v, _, _ := j.GetValue(), j.WS(), j.MatchByte(',')
			if f(k[1:len(k)-1], New(v)) {
				return
			}
		}
	}
}

func (j *Json) ForEach(f func(string, Json) bool) {
	if j.MatchByte('[') {
		for i := 0; j.WS() && !j.MatchByte(']'); i++ {
			k := strconv.Itoa(i)
			v, _, _ := j.GetValue(), j.WS(), j.MatchByte(',')
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
	if j.MatchUntilAnyByte4(',', '}', ']', ' ') { // Match Anything.
		return j.Token(m)
	}
	return ""
}

func (j *Json) MatchString() bool {
	return j.UtilMatchString('"')
}

func (j *Json) WS() bool {
	j.MatchWhileAnyByte4(' ', '\t', '\n', '\r')
	return true
}

// #endregion Json

func num2str(k string, v Json) (string, string) {
	if v.IsNumber() {
		return k, `"` + v.String() + `"`
	}
	return k, v.String()
}
