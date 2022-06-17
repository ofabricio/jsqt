package jsqt

import (
	"strconv"
	"strings"

	. "github.com/ofabricio/scanner"
)

func Get(jsn, qry string) Json {
	src := New(jsn)
	q := Query{Scanner: Scanner(qry), Root: src}
	return New(q.Parse(src))
}

func Get2(jsn, qry string) Json {
	src := New(jsn)
	q := Query{Scanner: Scanner(qry), Root: src}
	return q.Parse2(src)
}

func New(jsn string) Json {
	return Json{Scanner(jsn)}
}

// #region Query

type Query struct {
	Scanner
	Root Json
}

func (q *Query) Parse2(j Json) Json {
	return q.ParseFunc2(j)
}

func (q *Query) ParseFunc2(j Json) Json {
	if q.MatchByte('(') {
		if fname := q.TokenAnything(); fname != "" {
			v, _ := q.CallFunc2(fname, j), q.MatchByte(')')
			return v
		}
		if q.MatchByte(')') {
			return j
		}
	}
	return New("")
}

func (q *Query) CallFunc2(fname string, j Json) Json {
	switch fname {
	case "get":
		return FuncGet2(q, j)
	case "obj":
		return FuncObj2(q, j)
	case "collect":
		return FuncCollect2(q, j)
	case "arr":
		return FuncArr2(q, j)
	case "raw":
		return FuncRaw2(q, j)
	case "root":
		return q.Root
	case "size":
		return FuncSize2(q, j)
	case "default":
		return FuncDefault2(q, j)
	case "omitempty":
		return FuncOmitempty2(q, j)
	case "merge":
		return FuncMerge2(q, j)
	case "join":
		return FuncJoin2(q, j)
	case "iterate":
		return FuncIterate2(q, j)
	case "flatten":
		v := j.String()
		return New(v[1 : len(v)-1])
	default:
		return New("")
	}
}

func (q *Query) ParseFuncArg2(j Json) Json {
	if q.MatchByte(' ') {
		if v := q.ParseFunc2(j); v.String() != "" {
			return v
		}
		return q.ParseFuncArgKey(j)
	}
	return New("")
}

func (q *Query) ParseFuncArgFunc(j Json) Json {
	q.MatchByte(' ')
	return q.ParseFunc2(j)
}

func (q *Query) ParseFuncArgKey(j Json) Json {
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

func (q *Query) ParseFuncArgRaw(j Json) Json {
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

func FuncRaw2(q *Query, j Json) Json {
	return q.ParseFuncArgRaw(j)
}

func FuncGet2(q *Query, j Json) Json {
	for !q.EqualByte(')') {
		j = q.ParseFuncArg2(j)
	}
	return j
}

func FuncArr2(q *Query, j Json) Json {
	var o strings.Builder
	o.WriteString("[")
	for !q.EqualByte(')') {
		if o.Len() > 1 {
			o.WriteString(",")
		}
		v := q.ParseFuncArg2(j)
		o.WriteString(v.String())
	}
	o.WriteString("]")
	return New(o.String())
}

func FuncObj2(q *Query, j Json) Json {
	var o strings.Builder
	o.WriteString("{")
	for !q.EqualByte(')') {
		if k, v := q.ParseFuncArgRaw(j), q.ParseFuncArg2(j); v.String() != "" {
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

func FuncCollect2(q *Query, j Json) Json {
	var o strings.Builder
	o.WriteString("[")
	for !q.EqualByte(')') {
		if j.IsArray() {
			q.ForEach(j, func(sub *Query, item Json) {
				for !sub.EqualByte(')') {
					item = sub.ParseFuncArg2(item)
				}
				if item.String() != "" {
					if o.Len() > 1 {
						o.WriteString(",")
					}
					o.WriteString(item.String())
				}
			})
		} else {
			j = q.ParseFuncArgKey(j)
		}
	}
	o.WriteString("]")
	return New(o.String())
}

func FuncDefault2(q *Query, j Json) Json {
	d := q.ParseFuncArgRaw(j) // Default value.
	if j.String() == "" {
		return d
	}
	return j
}

func FuncOmitempty2(q *Query, j Json) Json {
	// v := q.ParseFuncArg2(j)
	if j.String() == "{}" || j.String() == "[]" {
		return New("")
	}
	return j
}

func FuncSize2(q *Query, j Json) Json {
	c := 0
	j.ForEach(func(i string, v Json) bool {
		c++
		return false
	})
	return New(strconv.Itoa(c))
}

func FuncMerge2(q *Query, j Json) Json {
	done := make(map[string]bool)
	var b strings.Builder
	b.WriteString("{")
	// j = New(q.ParseFuncArg(j)) // Input.
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

func FuncJoin2(q *Query, j Json) Json {
	var o strings.Builder
	o.WriteString("{")
	// j = New(q.ParseFuncArg(j)) // Input.
	q.ForEach(j, func(sub *Query, item Json) {
		for !sub.EqualByte(')') {
			if o.Len() > 1 {
				o.WriteString(",")
			}
			k := sub.ParseFuncArg2(item) // Key.
			v := sub.ParseFuncArg2(item) // Value.
			o.WriteString(k.String())
			o.WriteString(":")
			o.WriteString(v.String())
		}
	})
	o.WriteString("}")
	return New(o.String())
}

func FuncIterate2(q *Query, j Json) Json {
	m := q.TokenAnything() // Map function.
	_ = m
	// TODO: create functions map.
	return New(j.Iterate(num2str))
}

func (q *Query) Parse(j Json) string {
	if v := q.ParseFunc(j); v != "" {
		return v
	}
	return ""
}

func (q *Query) ParseFuncArg(j Json) string {
	if q.MatchByte(' ') {
		if v := q.ParseFunc(j); v != "" {
			return v
		}
		m := q.Mark()
		if q.UtilMatchString('"') {
			return q.Token(m)
		}
		if q.MatchUntilAnyByte(' ', ')') { // Anything.
			return q.Token(m)
		}
	}
	return ""
}

func (q *Query) ParseFunc(j Json) string {
	if q.MatchByte('(') {
		if fname := q.TokenAnything(); fname != "" {
			v, _ := q.CallFunc(fname, j), q.MatchByte(')')
			return v
		}
		if q.MatchByte(')') {
			return j.String()
		}
	}
	return ""
}

func (q *Query) TokenAnything() string {
	return q.TokenFor(func() bool {
		return q.MatchUntilAnyByte(' ', ')')
	})
}

func (q *Query) CallFunc(fname string, j Json) string {
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
	case "iterate":
		return FuncIterate(q, j)
	default:
		return ""
	}
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

// #region Functions

func FuncRaw(q *Query, j Json) string {
	return q.ParseFuncArg(j)
}

func FuncRoot(q *Query, j Json) string {
	return q.Root.String()
}

func FuncCurrent(q *Query, j Json) string {
	return j.String()
}

func FuncGet(q *Query, j Json) string {
	for !q.EqualByte(')') {
		key := q.ParseFuncArg(j)
		if key == "" {
			return j.String()
		}
		if key[0] == '"' {
			key = key[1 : len(key)-1]
		}
		j = j.Get(key)
	}
	return j.String()
}

func FuncObj(q *Query, j Json) string {
	var o strings.Builder
	o.WriteString("{")
	for !q.EqualByte(')') {
		if k, v := q.ParseFuncArg(j), q.ParseFuncArg(j); v != "" {
			if o.Len() > 1 {
				o.WriteString(",")
			}
			if k[0] == '"' {
				o.WriteString(k)
			} else {
				o.WriteString(`"`)
				o.WriteString(k)
				o.WriteString(`"`)
			}
			o.WriteString(`:`)
			o.WriteString(v)
		}
	}
	o.WriteString("}")
	return o.String()
}

func FuncArr(q *Query, j Json) string {
	var o strings.Builder
	o.WriteString("[")
	for !q.EqualByte(')') {
		if o.Len() > 1 {
			o.WriteString(",")
		}
		v := q.ParseFuncArg(j)
		o.WriteString(v)
	}
	o.WriteString("]")
	return o.String()
}

func FuncFlatten(q *Query, j Json) string {
	if v := q.ParseFuncArg(j); len(v) > 1 {
		return v[1 : len(v)-1]
	}
	return j.String()
}

func FuncCollect(q *Query, j Json) string {
	var o strings.Builder
	o.WriteString("[")
	j = New(q.ParseFuncArg(j)) // Input.
	q.ForEach(j, func(sub *Query, item Json) {
		if v := sub.ParseFuncArg(item); v != "" { // Filter.
			if v = sub.ParseFuncArg(item); v != "" { // Mapper.
				if o.Len() > 1 {
					o.WriteString(",")
				}
				o.WriteString(v)
			}
		}
	})
	o.WriteString("]")
	return o.String()
}

func FuncJoin(q *Query, j Json) string {
	var o strings.Builder
	o.WriteString("{")
	j = New(q.ParseFuncArg(j)) // Input.
	q.ForEach(j, func(sub *Query, item Json) {
		if f := sub.ParseFuncArg(item); f == "" { // Filter.
			return
		}
		for !sub.EqualByte(')') {
			if o.Len() > 1 {
				o.WriteString(",")
			}
			k := sub.ParseFuncArg(item) // Key.
			v := sub.ParseFuncArg(item) // Value.
			o.WriteString(k)
			o.WriteString(":")
			o.WriteString(v)
		}
	})
	o.WriteString("}")
	return o.String()
}

func FuncSize(q *Query, j Json) string {
	c := 0
	j = New(q.ParseFuncArg(j)) // Input.
	j.ForEach(func(i string, v Json) bool {
		c++
		return false
	})
	return strconv.Itoa(c)
}

func FuncMerge(q *Query, j Json) string {
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

func FuncDefault(q *Query, j Json) string {
	v := q.ParseFuncArg(j) // Input.
	d := q.ParseFuncArg(j) // Value.
	if v == "" {
		return d
	}
	return v
}

func FuncOmitEmpty(q *Query, j Json) string {
	v := q.ParseFuncArg(j)
	if v == "{}" || v == "[]" {
		return ""
	}
	return v
}

func FuncIterate(q *Query, j Json) string {
	v := New(q.ParseFuncArg(j)) // Input.
	m := q.TokenAnything()      // Map function.
	_ = m
	// TODO: create functions map.
	return v.Iterate(num2str)
}

// #endregion Functions

func num2str(k string, v Json) (string, string) {
	if v.IsNumber() {
		return k, `"` + v.String() + `"`
	}
	return k, v.String()
}
