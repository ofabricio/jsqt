package jsqt

import (
	"fmt"
	"strconv"
	"strings"

	. "github.com/ofabricio/scanner"
)

func Get(jsn, qry string) Json {
	src := New(jsn)
	src.ws()
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
	return q.ParseFun(j)
}

func (q *Query) ParseFun(j Json) Json {
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

func (q *Query) ParseKey(j Json) Json {
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

func (q *Query) ParseRaw(j Json) Json {
	raw := ""
	m := q.Mark()
	if q.UtilMatchString('"') {
		raw = q.Token(m)
	} else if q.MatchUntilAnyByte(' ', ')') { // Anything.
		raw = q.Token(m)
	}
	return New(raw)
}

func (q *Query) ParseArgFunOrKey(j Json) Json {
	if q.MatchByte(' ') {
		if q.IsFun() {
			return q.ParseFun(j)
		}
		return q.ParseKey(j)
	}
	return New("")
}

func (q *Query) ParseArgFunOrRaw(j Json) Json {
	if q.MatchByte(' ') {
		if q.IsFun() {
			return q.ParseFun(j)
		}
		return q.ParseRaw(j)
	}
	return New("")
}

func (q *Query) ParseArgFun(j Json) Json {
	if q.MatchByte(' ') {
		return q.ParseFun(j)
	}
	return New("")
}

func (q *Query) ParseArgKey(j Json) Json {
	if q.MatchByte(' ') {
		return q.ParseKey(j)
	}
	return New("")
}

func (q *Query) ParseArgRaw(j Json) Json {
	if q.MatchByte(' ') {
		return q.ParseRaw(j)
	}
	return New("")
}

func (q Query) IsFun() bool {
	return q.EqualByte('(')
}

func (q *Query) CallFunc(fname string, j Json) Json {
	switch fname {
	case "get":
		return funcGet(q, j)
	case "obj":
		return funcObj(q, j)
	case "arr":
		return funcArr(q, j)
	case "raw":
		return funcRaw(q, j)
	case "collect":
		return funcCollect(q, j)
	case "flatten":
		return j.Flatten()
	case "size":
		return j.Size()
	case "default":
		return funcDefault(q, j)
	case "merge":
		return j.Merge()
	case "iterate":
		return funcIterate(q, j)
	case "is-num":
		return funcIsNum(q, j)
	case "is-obj":
		return funcIsObj(q, j)
	case "is-arr":
		return funcIsArr(q, j)
	case "is-str":
		return funcIsStr(q, j)
	case "is-bool":
		return funcIsBool(q, j)
	case "is-null":
		return funcIsNull(q, j)
	case "is-empty":
		return funcIsEmpty(q, j)
	case "is-empty-arr":
		return funcIsEmptyArr(q, j)
	case "is-empty-obj":
		return funcIsEmptyObj(q, j)
	case "is-empty-str":
		return funcIsEmptyStr(q, j)
	case "is-some":
		return funcIsSome(q, j)
	case "is-void":
		return funcIsVoid(q, j)
	case "is-blank":
		return funcIsBlank(q, j)
	case "is-nully":
		return funcIsNully(q, j)
	case "truthy":
		return funcIsTruthy(q, j)
	case "falsy":
		return funcIsFalsy(q, j)
	case "if":
		return funcIf(q, j)
	case "root":
		return q.Root
	case ".":
		return j
	case "==":
		return funcEQ(q, j)
	case "!=":
		return funcNEQ(q, j)
	case ">=":
		return funcGTE(q, j)
	case "<=":
		return funcLTE(q, j)
	case ">":
		return funcGT(q, j)
	case "<":
		return funcLT(q, j)
	case "or":
		return funcOr(q, j)
	case "and":
		return funcAnd(q, j)
	case "not":
		return funcNot(q, j)
	case "bool":
		return funcBool(q, j)
	case "debug":
		return funcDebug(q, j)
	case "keys":
		return j.Keys()
	case "values":
		return j.Values()
	case "entries":
		return j.Entries()
	case "ugly":
		return j.Uglify()
	case "pretty":
		return j.Prettify()
	case "to-json":
		return j.ToJson()
	case "to-str":
		return j.ToString()
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
	j.ForEach(func(i, item Json) bool {
		end := ini
		f(&end, item)
		*q = end
		return false
	})
}

// #region Functions

func funcRaw(q *Query, j Json) Json {
	return q.ParseArgRaw(j)
}

func funcGet(q *Query, j Json) Json {
	for !q.EqualByte(')') {
		j = q.ParseArgFunOrKey(j)
	}
	return j
}

func funcArr(q *Query, j Json) Json {
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

func funcObj(q *Query, j Json) Json {
	var o strings.Builder
	o.WriteString("{")
	for !q.EqualByte(')') {
		if k, v := q.ParseArgFunOrRaw(j), q.ParseArgFunOrKey(j); v.IsAnything() {
			if o.Len() > 1 {
				o.WriteString(",")
			}
			if k.IsString() {
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

func funcCollect(q *Query, j Json) Json {
	var o strings.Builder
	o.WriteString("[")
	for !q.EqualByte(')') {
		if j.IsArray() && !j.IsEmptyArray() {
			q.ForEach(j, func(sub *Query, item Json) {
				for !sub.EqualByte(')') {
					item = sub.ParseArgFunOrKey(item)
				}
				if item.IsAnything() {
					if o.Len() > 1 {
						o.WriteString(",")
					}
					o.WriteString(item.String())
				}
			})
		} else {
			j = q.ParseArgFunOrKey(j)
		}
	}
	o.WriteString("]")
	return New(o.String())
}

func funcDefault(q *Query, j Json) Json {
	d := q.ParseArgRaw(j) // Default value.
	if j.String() == "" {
		return d
	}
	return j
}

func funcIterate(q *Query, j Json) Json {
	ini := *q
	return j.Iterate(func(k, v Json) (Json, Json) {
		arr := New(`[` + k.String() + "," + v.String() + `]`)
		i := ini
		mk, mv := i.ParseArgFunOrKey(arr), i.ParseArgFunOrKey(arr)
		*q = i
		return mk, mv
	})
}

func funcIf(q *Query, j Json) Json {
	cond := q.ParseArgFunOrKey(j)
	then := q.ParseArgFunOrKey(j)
	elze := q.ParseArgFunOrKey(j)
	if cond.IsAnything() {
		return then
	}
	return elze
}

// #region funcIS

func funcIsNum(q *Query, j Json) Json {
	if j.IsNumber() {
		return j
	}
	return New("")
}

func funcIsObj(q *Query, j Json) Json {
	if j.IsObject() {
		return j
	}
	return New("")
}

func funcIsArr(q *Query, j Json) Json {
	if j.IsArray() {
		return j
	}
	return New("")
}

func funcIsStr(q *Query, j Json) Json {
	if j.IsString() {
		return j
	}
	return New("")
}

func funcIsBool(q *Query, j Json) Json {
	if j.IsBool() {
		return j
	}
	return New("")
}

func funcIsNull(q *Query, j Json) Json {
	if j.IsNull() {
		return j
	}
	return New("")
}

func funcIsEmpty(q *Query, j Json) Json {
	if j.IsEmpty() {
		return j
	}
	return New("")
}

func funcIsEmptyObj(q *Query, j Json) Json {
	if j.IsEmptyObject() {
		return j
	}
	return New("")
}

func funcIsEmptyArr(q *Query, j Json) Json {
	if j.IsEmptyArray() {
		return j
	}
	return New("")
}

func funcIsEmptyStr(q *Query, j Json) Json {
	if j.IsEmptyString() {
		return j
	}
	return New("")
}

func funcIsTruthy(q *Query, j Json) Json {
	if j.IsTruthy() {
		return j
	}
	return New("")
}

func funcIsFalsy(q *Query, j Json) Json {
	if j.IsFalsy() {
		return j
	}
	return New("")
}

func funcIsSome(q *Query, j Json) Json {
	if j.IsSome() {
		return j
	}
	return New("")
}

func funcIsVoid(q *Query, j Json) Json {
	if j.IsVoid() {
		return j
	}
	return New("")
}

func funcIsNully(q *Query, j Json) Json {
	if j.IsNully() {
		return j
	}
	return New("")
}

func funcIsBlank(q *Query, j Json) Json {
	if j.IsBlank() {
		return j
	}
	return New("")
}

// #endregion funcIS

// #region Filters

func funcEQ(q *Query, j Json) Json {
	f := q.ParseArgFunOrKey(j) // Field.
	v := q.ParseArgFunOrRaw(j) // Value.
	if f.String() == v.String() {
		return j
	}
	return New("")
}

func funcNEQ(q *Query, j Json) Json {
	f := q.ParseArgFunOrKey(j) // Field.
	v := q.ParseArgFunOrRaw(j) // Value.
	if f.String() != v.String() {
		return j
	}
	return New("")
}

func funcGTE(q *Query, j Json) Json {
	f := q.ParseArgFunOrKey(j) // Field.
	v := q.ParseArgFunOrRaw(j) // Value.
	if f.String() >= v.String() {
		return j
	}
	return New("")
}

func funcLTE(q *Query, j Json) Json {
	f := q.ParseArgFunOrKey(j) // Field.
	v := q.ParseArgFunOrRaw(j) // Value.
	if f.String() <= v.String() {
		return j
	}
	return New("")
}

func funcGT(q *Query, j Json) Json {
	f := q.ParseArgFunOrKey(j) // Field.
	v := q.ParseArgFunOrRaw(j) // Value.
	if f.String() > v.String() {
		return j
	}
	return New("")
}

func funcLT(q *Query, j Json) Json {
	f := q.ParseArgFunOrKey(j) // Field.
	v := q.ParseArgFunOrRaw(j) // Value.
	if f.String() < v.String() {
		return j
	}
	return New("")
}

// #endregion Filters

func funcOr(q *Query, j Json) Json {
	a := q.ParseArgFun(j)
	b := q.ParseArgFun(j)
	if a.IsAnything() || b.IsAnything() {
		return j
	}
	return New("")
}

func funcAnd(q *Query, j Json) Json {
	a := q.ParseArgFun(j)
	b := q.ParseArgFun(j)
	if a.IsAnything() && b.IsAnything() {
		return j
	}
	return New("")
}

func funcNot(q *Query, j Json) Json {
	a := q.ParseArgFun(j)
	if a.IsAnything() {
		return New("")
	}
	return j
}

func funcBool(q *Query, j Json) Json {
	if j.IsAnything() {
		return New("true")
	}
	return New("false")
}

func funcDebug(q *Query, j Json) Json {
	msg := "debug"
	if !q.EqualByte(')') {
		msg = q.ParseArgRaw(j).String()
	}
	fmt.Printf("[%s] %s\n", msg, j)
	return j
}

// #endregion Functions

// #endregion Query

// #region Json

type Json struct {
	s Scanner
}

// String returns the raw JSON data.
func (j Json) String() string {
	return j.s.String()
}

// ToString converts a JSON value to a JSON string.
func (j Json) ToString() Json {
	if j.IsString() {
		return j
	}
	return New(strconv.Quote(j.String()))
}

// ToString converts a JSON string to a JSON value.
func (j Json) ToJson() Json {
	if j.IsString() {
		v, _ := strconv.Unquote(j.String())
		return New(v)
	}
	return j
}

// Str returns a string value.
// Example: "Hello" -> Hello.
func (j Json) Str() string {
	v := j.String()
	if j.IsString() {
		v, _ = strconv.Unquote(v)
	}
	return v
}

// TrimKey removes the quotes from an object key.
// Example: "name" -> name.
func (j Json) TrimKey() string {
	v := j.String()
	if j.IsString() {
		return v[1 : len(v)-1]
	}
	return v
}

// Int returns an int value.
func (j Json) Int() int {
	v, _ := strconv.ParseInt(j.String(), 10, 0)
	return int(v)
}

// Float returns a float value.
func (j Json) Float() float64 {
	v, _ := strconv.ParseFloat(j.String(), 64)
	return v
}

// Bool returns a bool value.
func (j Json) Bool() bool {
	v, _ := strconv.ParseBool(j.String())
	return v
}

func (j Json) IsObject() bool {
	return j.s.EqualByte('{')
}

func (j Json) IsArray() bool {
	return j.s.EqualByte('[')
}

func (j Json) IsNumber() bool {
	return j.s.EqualByteRange('0', '9') || j.s.EqualByte('-')
}

func (j Json) IsString() bool {
	return j.s.EqualByte('"')
}

func (j Json) IsBool() bool {
	return j.IsTrue() || j.IsFalse()
}

func (j Json) IsTrue() bool {
	return j.s.EqualByte('t')
}

func (j Json) IsFalse() bool {
	return j.s.EqualByte('f')
}

func (j Json) IsNull() bool {
	return j.s.EqualByte('n')
}

func (j Json) IsEmptyString() bool {
	return j.String() == `""`
}

func (j Json) IsEmptyObject() bool {
	return j.s.MatchByte('{') && j.ws() && j.s.EqualByte('}')
}

func (j Json) IsEmptyArray() bool {
	return j.s.MatchByte('[') && j.ws() && j.s.EqualByte(']')
}

func (j Json) IsVoid() bool {
	return j.IsEmptyObject() || j.IsEmptyArray()
}

func (j Json) IsEmpty() bool {
	return j.IsEmptyObject() || j.IsEmptyArray() || j.IsEmptyString()
}

func (j Json) IsBlank() bool {
	return j.IsEmptyObject() || j.IsEmptyArray() || j.IsNull()
}

func (j Json) IsNully() bool {
	return j.IsEmptyObject() || j.IsEmptyArray() || j.IsEmptyString() || j.IsNull()
}

func (j Json) IsTruthy() bool {
	return !j.IsFalsy() && j.String() != ""
}

func (j Json) IsFalsy() bool {
	return j.IsEmptyObject() || j.IsEmptyArray() || j.IsEmptyString() ||
		j.IsFalse() || j.IsNull() || j.s.EqualByte('0')
}

func (j Json) IsSome() bool {
	return !j.IsNull() && j.String() != ""
}

func (j Json) IsAnything() bool {
	return j.String() != ""
}

// IterateValues iterates over the values of a valid Json.
func (j Json) IterateValues(m func(Json) Json) Json {
	var o strings.Builder
	if j.IsObject() {
		o.WriteString("{")
		j.ForEachKeyVal(func(k, v Json) bool {
			if o.Len() > 1 {
				o.WriteString(",")
			}
			o.WriteString(k.String())
			o.WriteByte(':')
			o.WriteString(v.IterateValues(m).String())
			return false
		})
		o.WriteString("}")
		return New(o.String())
	} else if j.IsArray() {
		o.WriteString("[")
		j.ForEach(func(i, v Json) bool {
			if o.Len() > 1 {
				o.WriteString(",")
			}
			o.WriteString(v.IterateValues(m).String())
			return false
		})
		o.WriteString("]")
		return New(o.String())
	}
	return m(j)
}

// Iterate iterates over a valid Json.
func (j Json) Iterate(m func(Json, Json) (Json, Json)) Json {
	_, mv := j.iterateInternal(New("null"), m)
	return mv
}

func (j Json) iterateInternal(k Json, m func(Json, Json) (Json, Json)) (Json, Json) {
	mk, mv := m(k, j)
	var o strings.Builder
	if mv.IsObject() {
		o.WriteString("{")
		mv.ForEachKeyVal(func(k, v Json) bool {
			if o.Len() > 1 {
				o.WriteString(",")
			}
			sk, sv := v.iterateInternal(k, m)
			o.WriteString(sk.String())
			o.WriteString(`:`)
			o.WriteString(sv.String())
			return false
		})
		o.WriteString("}")
		mv = New(o.String())
	} else if mv.IsArray() {
		o.WriteString("[")
		mv.ForEach(func(i, v Json) bool {
			if o.Len() > 1 {
				o.WriteString(",")
			}
			o.WriteString(v.Iterate(m).String())
			return false
		})
		o.WriteString("]")
		mv = New(o.String())
	}
	return mk, mv
}

func (j Json) Get(keyOrIndex string) (r Json) {
	f := func(k, v Json) bool {
		if k.TrimKey() == keyOrIndex {
			r = v
			return true
		}
		return false
	}
	j.ForEachKeyVal(f)
	j.ForEach(f)
	return r
}

func (j Json) Collect(keyOrIndex string) Json {
	if j.IsArray() && !(keyOrIndex[0] >= '0' && keyOrIndex[0] <= '9') {
		var o strings.Builder
		o.WriteString("[")
		j.ForEach(func(i, v Json) bool {
			if v = v.Collect(keyOrIndex); v.IsAnything() {
				if o.Len() > 1 {
					o.WriteString(",")
				}
				o.WriteString(v.String())
			}
			return false
		})
		o.WriteString("]")
		return New(o.String())
	}
	return j.Get(keyOrIndex)
}

func (j Json) ForEachKeyVal(f func(k, v Json) bool) {
	if j.s.MatchByte('{') {
		for j.ws() && !j.s.MatchByte('}') {
			k, _, _, _ := j.s.TokenFor(j.matchString), j.ws(), j.s.MatchByte(':'), j.ws()
			v, _, _ := j.getValue(), j.ws(), j.s.MatchByte(',')
			if f(New(k), New(v)) {
				return
			}
		}
	}
}

func (j Json) ForEach(f func(k, v Json) bool) {
	if j.s.MatchByte('[') {
		for i := 0; j.ws() && !j.s.MatchByte(']'); i++ {
			k := strconv.Itoa(i)
			v, _, _ := j.getValue(), j.ws(), j.s.MatchByte(',')
			if f(New(k), New(v)) {
				return
			}
		}
	}
}

func (j *Json) getValue() string {
	m := j.s.Mark()
	if j.s.UtilMatchOpenCloseCount('{', '}', '"') { // Match Object.
		return j.s.Token(m)
	}
	if j.s.UtilMatchOpenCloseCount('[', ']', '"') { // Match Array.
		return j.s.Token(m)
	}
	if j.matchString() {
		return j.s.Token(m)
	}
	if j.s.MatchUntilAnyByte4(',', '}', ']', ' ') { // Match Anything.
		return j.s.Token(m)
	}
	return ""
}

func (j *Json) matchString() bool {
	return j.s.UtilMatchString('"')
}

func (j *Json) ws() bool {
	j.s.MatchWhileAnyByte4(' ', '\t', '\n', '\r')
	return true
}

func (j Json) Flatten() Json {
	if j.IsArray() {
		v := j.String()
		return New(v[1 : len(v)-1])
	}
	return j
}

func (j Json) Size() Json {
	if j.IsString() {
		return New(strconv.Itoa(len(j.String()) - 2))
	}
	c := 0
	j.ForEach(func(i, v Json) bool {
		c++
		return false
	})
	return New(strconv.Itoa(c))
}

func (j Json) Keys() Json {
	var o strings.Builder
	o.WriteString("[")
	j.ForEachKeyVal(func(k, v Json) bool {
		if o.Len() > 1 {
			o.WriteString(",")
		}
		o.WriteString(k.String())
		return false
	})
	o.WriteString("]")
	return New(o.String())
}

func (j Json) Values() Json {
	var o strings.Builder
	o.WriteString("[")
	j.ForEachKeyVal(func(k, v Json) bool {
		if o.Len() > 1 {
			o.WriteString(",")
		}
		o.WriteString(v.String())
		return false
	})
	o.WriteString("]")
	return New(o.String())
}

func (j Json) Entries() Json {
	var o strings.Builder
	o.WriteString("[")
	j.ForEachKeyVal(func(k, v Json) bool {
		if o.Len() > 1 {
			o.WriteString(",")
		}
		o.WriteString("[")
		o.WriteString(k.String())
		o.WriteString(`,`)
		o.WriteString(v.String())
		o.WriteString("]")
		return false
	})
	o.WriteString("]")
	return New(o.String())
}

func (j Json) Merge() Json {
	done := make(map[string]bool)
	var o strings.Builder
	o.WriteString("{")
	j.ForEach(func(i, v Json) bool {
		v.ForEachKeyVal(func(k, v Json) bool {
			if !done[k.String()] {
				if o.Len() > 1 {
					o.WriteString(",")
				}
				o.WriteString(k.String())
				o.WriteString(`:`)
				o.WriteString(v.String())
			}
			done[k.String()] = true
			return false
		})
		return false
	})
	o.WriteString("}")
	return New(o.String())
}

func (j Json) Uglify() Json {
	s := j.String()
	var x strings.Builder
	x.Grow(len(s))
	for i := 0; i < len(s); i++ {
		if s[i] > ' ' {
			if s[i] == '"' {
				ini := i
				for i = i + 1; i < len(s); i++ {
					if s[i] == '"' && s[i-1] != '\\' {
						x.WriteString(s[ini : i+1])
						break
					}
				}
			} else {
				x.WriteByte(s[i])
			}
		}
	}
	return New(x.String())
}

func (j Json) Prettify() Json {
	pad := "    "
	s := j.String()
	var x strings.Builder
	x.Grow(len(s) << 1)
	depth := 0
top:
	for i := 0; i < len(s); i++ {
		if s[i] > ' ' {
			switch s[i] {
			case '"':
				ini := i
				for i = i + 1; i < len(s); i++ {
					if s[i] == '"' && s[i-1] != '\\' {
						x.WriteString(s[ini : i+1])
						break
					}
				}
			case ',':
				x.WriteString(",\n")
				for d := 0; d < depth; d++ {
					x.WriteString(pad)
				}
			case '{':
				// Check for empty object.
				for j := i + 1; j < len(s); j++ {
					if s[j] > ' ' {
						if s[j] == '}' {
							i = j
							x.WriteString("{}")
							continue top
						}
						break
					}
				}
				x.WriteString("{\n")
				depth++
				for d := 0; d < depth; d++ {
					x.WriteString(pad)
				}
			case '}':
				x.WriteString("\n")
				depth--
				for d := 0; d < depth; d++ {
					x.WriteString(pad)
				}
				x.WriteString("}")
			case '[':
				// Check for empty array.
				for j := i + 1; j < len(s); j++ {
					if s[j] > ' ' {
						if s[j] == ']' {
							i = j
							x.WriteString("[]")
							continue top
						}
						break
					}
				}
				x.WriteString("[\n")
				depth++
				for d := 0; d < depth; d++ {
					x.WriteString(pad)
				}
			case ']':
				x.WriteString("\n")
				depth--
				for d := 0; d < depth; d++ {
					x.WriteString(pad)
				}
				x.WriteString("]")
			case ':':
				x.WriteString(": ")
			default:
				x.WriteByte(s[i])
			}
		}
	}
	return New(x.String())
}

// #endregion Json
