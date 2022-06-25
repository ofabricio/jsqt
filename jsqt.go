package jsqt

import (
	"fmt"
	"sort"
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

func (q *Query) ParseFunOrKey(j Json) Json {
	if q.EqualByte('(') {
		return q.ParseFun(j)
	}
	return q.ParseKey(j)
}

func (q *Query) ParseFunOrRaw(j Json) Json {
	if q.EqualByte('(') {
		return q.ParseFun(j)
	}
	return q.ParseRaw()
}

func (q *Query) ParseFun(j Json) Json {
	if q.MatchByte('(') {
		fname := q.ParseRaw().String()
		j, _ = q.CallFun(fname, j), q.MatchByte(')')
		q.ws()
	}
	return j
}

func (q *Query) ParseKey(j Json) Json {
	key := ""
	if m := q.Mark(); q.UtilMatchString('"') {
		key = q.Token(m)
		key = key[1 : len(key)-1]
	} else if q.MatchUntilAnyByte3(' ', ')', 0) { // Anything.
		key = q.Token(m)
	}
	q.ws()
	return j.Get(key)
}

func (q *Query) ParseRaw() Json {
	raw := ""
	if m := q.Mark(); q.UtilMatchString('"') {
		raw = q.Token(m)
	} else if q.MatchUntilAnyByte3(' ', ')', 0) { // Anything.
		raw = q.Token(m)
	}
	q.ws()
	return New(raw)
}

func (q *Query) SkipArgs() {
	for !q.EqualByte(')') {
		q.SkipArg()
	}
}

func (q *Query) SkipArg() {
	q.MatchArg()
	q.ws()
}

func (q *Query) MatchArg() bool {
	return q.UtilMatchOpenCloseCount('(', ')', '"') || q.UtilMatchString('"') || q.MatchUntilAnyByte(' ', ')')
}

func (q *Query) GrabArg() Query {
	qq := *q
	qq.Scanner = Scanner(q.TokenFor(q.MatchArg))
	q.ws()
	return qq
}

func (q *Query) CallFun(fname string, j Json) Json {
	switch fname {
	case "get":
		return funcGet(q, j)
	case "obj":
		return funcObj(q, j)
	case "arr":
		return funcArr(q, j)
	case "raw":
		return q.ParseRaw()
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
	case "iterate-v":
		return funcIterateValues(q, j)
	case "iterate-k":
		return funcIterateKeys(q, j)
	case "iterate-kv":
		return funcIterateKeysValues(q, j)
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
	case "either":
		return funcEither(q, j)
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
	case "jsonify":
		return j.Jsonify()
	case "stringify":
		return j.Stringify()
	case "upper":
		return New(strings.ToUpper(j.String()))
	case "lower":
		return New(strings.ToLower(j.String()))
	case "replace":
		return funcReplace(q, j)
	case "concat":
		return funcConcat(q, j)
	case "sort":
		return funcSort(q, j)
	default:
		return New("")
	}
}

func (q *Query) IsEmpty() bool {
	return q.String() == ""
}

func (q *Query) ws() {
	q.MatchByte(' ')
}

// #region Functions

func funcGet(q *Query, j Json) Json {
	for !q.EqualByte(')') {
		j = q.ParseFunOrKey(j)
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
		v := q.ParseFunOrKey(j)
		o.WriteString(v.String())
	}
	o.WriteString("]")
	return New(o.String())
}

func funcObj(q *Query, j Json) Json {
	var o strings.Builder
	o.WriteString("{")
	for !q.EqualByte(')') {
		if k, v := q.ParseFunOrRaw(j), q.ParseFunOrKey(j); v.IsAnything() {
			if o.Len() > 1 {
				o.WriteString(",")
			}
			o.WriteByte('"')
			o.WriteString(k.TrimKey())
			o.WriteString(`":`)
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
			ini := *q
			j.ForEach(func(i, item Json) bool {
				sub := ini
				for !sub.EqualByte(')') {
					item = sub.ParseFunOrKey(item)
				}
				if item.IsAnything() {
					if o.Len() > 1 {
						o.WriteString(",")
					}
					o.WriteString(item.String())
				}
				*q = sub
				return false
			})
		} else {
			j = q.ParseFunOrKey(j)
		}
	}
	o.WriteString("]")
	return New(o.String())
}

func funcDefault(q *Query, j Json) Json {
	v := q.ParseRaw()
	if j.String() == "" {
		return v
	}
	return j
}

func funcIterate(q *Query, j Json) Json {
	keyArg := q.GrabArg()
	valArg := q.GrabArg()
	return j.Iterate(func(k, v Json) (Json, Json) {
		arr := New(`[` + k.String() + "," + v.String() + `]`)
		ka := keyArg
		va := valArg
		return ka.ParseFunOrKey(arr), va.ParseFunOrKey(arr)
	})
}

func funcIterateKeys(q *Query, j Json) Json {
	keyArg := q.GrabArg()
	return j.IterateKeys(func(k Json) Json {
		ka := keyArg
		return ka.ParseFun(k)
	})
}

func funcIterateValues(q *Query, j Json) Json {
	valArg := q.GrabArg()
	return j.IterateValues(func(v Json) Json {
		va := valArg
		return va.ParseFun(v)
	})
}

func funcIterateKeysValues(q *Query, j Json) Json {
	keyvalArg := q.GrabArg()
	return j.IterateKeysValues(func(v Json) Json {
		kv := keyvalArg
		return kv.ParseFun(v)
	})
}

func funcIf(q *Query, j Json) Json {
	cond := q.ParseFunOrKey(j)
	then := q.GrabArg()
	elze := q.GrabArg()
	if cond.IsAnything() {
		return then.ParseFunOrKey(j)
	}
	return elze.ParseFunOrKey(j)
}

func funcEither(q *Query, j Json) Json {
	v := q.ParseFunOrKey(j)
	for v.IsNully() && !q.EqualByte(')') {
		v = q.ParseFunOrKey(j)
	}
	q.SkipArgs()
	return v
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
	a := q.ParseFunOrKey(j)
	b := q.ParseFunOrRaw(j)
	if a.String() == b.String() {
		return j
	}
	return New("")
}

func funcNEQ(q *Query, j Json) Json {
	a := q.ParseFunOrKey(j)
	b := q.ParseFunOrRaw(j)
	if a.String() != b.String() {
		return j
	}
	return New("")
}

func funcGTE(q *Query, j Json) Json {
	a := q.ParseFunOrKey(j)
	b := q.ParseFunOrRaw(j)
	if a.String() >= b.String() {
		return j
	}
	return New("")
}

func funcLTE(q *Query, j Json) Json {
	a := q.ParseFunOrKey(j)
	b := q.ParseFunOrRaw(j)
	if a.String() <= b.String() {
		return j
	}
	return New("")
}

func funcGT(q *Query, j Json) Json {
	a := q.ParseFunOrKey(j)
	b := q.ParseFunOrRaw(j)
	if a.String() > b.String() {
		return j
	}
	return New("")
}

func funcLT(q *Query, j Json) Json {
	a := q.ParseFunOrKey(j)
	b := q.ParseFunOrRaw(j)
	if a.String() < b.String() {
		return j
	}
	return New("")
}

// #endregion Filters

func funcOr(q *Query, j Json) Json {
	a := q.ParseFun(j)
	b := q.ParseFun(j)
	if a.IsAnything() || b.IsAnything() {
		return j
	}
	return New("")
}

func funcAnd(q *Query, j Json) Json {
	a := q.ParseFun(j)
	b := q.ParseFun(j)
	if a.IsAnything() && b.IsAnything() {
		return j
	}
	return New("")
}

func funcNot(q *Query, j Json) Json {
	a := q.ParseFun(j)
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
		msg = q.ParseRaw().String()
	}
	fmt.Printf("[%s] %s\n", msg, j)
	return j
}

func funcReplace(q *Query, j Json) Json {
	old := q.ParseRaw()
	new := q.ParseRaw()
	if j.IsString() {
		return New(strings.ReplaceAll(j.String(), old.TrimKey(), new.TrimKey()))
	}
	return j
}

func funcConcat(q *Query, j Json) Json {
	var o strings.Builder

	v := q.ParseFunOrKey(j)

	if v.IsString() {
		o.WriteString(v.Jsonify().String())
		for !q.EqualByte(')') {
			v := q.ParseFunOrKey(j)
			o.WriteString(v.Jsonify().String())
		}
		return New(o.String()).Stringify()
	}

	o.WriteByte('[')
	o.WriteString(v.Flatten().String())
	for !q.EqualByte(')') {
		v := q.ParseFunOrKey(j)
		o.WriteByte(',')
		o.WriteString(v.Flatten().String())
	}
	o.WriteByte(']')
	return New(o.String())
}

func funcSort(q *Query, j Json) Json {
	asc := q.ParseRaw().String() == "asc"
	arg := q.GrabArg()
	var items []string
	j.ForEach(func(i, v Json) bool {
		items = append(items, v.String())
		return false
	})
	sort.SliceStable(items, func(i, j int) bool {
		ia := arg
		ja := ia
		va := ia.ParseFunOrKey(New(items[i])).String()
		vb := ja.ParseFunOrKey(New(items[j])).String()
		if asc {
			return va < vb
		}
		return va > vb
	})
	return New("[" + strings.Join(items, ",") + "]")
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

// Stringify converts a JSON value to a JSON string.
// Examples:
//   "Hello" -> "Hello"
//   3       -> "3"
//   {}      -> "{}"
//   "{ \"hello\": \"world\" }" -> { "hello": "world" }
// Stringify reverts Jsonify.
func (j Json) Stringify() Json {
	if j.IsString() {
		return j
	}
	return New(strconv.Quote(j.String()))
}

// Jsonify converts a JSON string to a JSON value.
// Examples:
//   "Hello" -> "Hello"
//   "3"     -> 3
//   "{}"    -> {}
//   "{ \"hello\": \"world\"}" -> { "hello": "world" }
// Jsonify reverts Stringify.
func (j Json) Jsonify() Json {
	if j.IsString() && !j.IsEmptyString() {
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
	return !j.IsFalsy() && j.IsAnything()
}

func (j Json) IsFalsy() bool {
	return j.IsEmptyObject() || j.IsEmptyArray() || j.IsEmptyString() ||
		j.IsFalse() || j.IsNull() || j.s.EqualByte('0')
}

func (j Json) IsSome() bool {
	return !j.IsNull() && j.IsAnything()
}

func (j Json) IsAnything() bool {
	return j.String() != ""
}

func (j Json) Iterator(o *strings.Builder, k Json, m func(o *strings.Builder, k, v Json)) {
	m(o, k, j)
	i := 0
	if j.IsObject() {
		o.WriteString("{")
		j.ForEachKeyVal(func(k, v Json) bool {
			if i > 0 {
				o.WriteString(",")
			}
			i++
			v.Iterator(o, k, m)
			return false
		})
		o.WriteString("}")
	} else if j.IsArray() {
		o.WriteString("[")
		j.ForEach(func(k, v Json) bool {
			if i > 0 {
				o.WriteString(",")
			}
			i++
			v.Iterator(o, New(""), m)
			return false
		})
		o.WriteString("]")
	}
}

// IterateKeys iterates over the keys (excluding values)
// of a valid Json and apply a map function to transform
// each emitted value.
func (j Json) IterateKeys(m func(Json) Json) Json {
	var o strings.Builder
	o.Grow(len(j.s))
	for i := 0; i < len(j.s); i++ {
		if j.s[i] > ' ' {
			if j.s[i] == '"' {
				// Scans through the string.
				ini := i
				for i = i + 1; i < len(j.s); i++ {
					if j.s[i] == '"' && j.s[i-1] != '\\' {
						end := i + 1
						// Skip spaces.
						for i = i + 1; i < len(j.s) && j.s[i] <= ' '; i++ {
						}
						// Emits if a key.
						if j.s[i] == ':' {
							o.WriteString(m(New(j.s[ini:end].String())).String())
						} else {
							o.WriteString(j.s[ini:end].String())
						}
						o.WriteByte(j.s[i])
						break
					}
				}
			} else {
				o.WriteByte(j.s[i])
			}
		}
	}
	return New(o.String())
}

// IterateValues iterates over the values (excluding the keys)
// of a valid Json and apply a map function to transform each
// emitted value.
func (j Json) IterateValues(m func(Json) Json) Json {
	var o strings.Builder
	o.Grow(len(j.s))
	for i := 0; i < len(j.s); i++ {
		if j.s[i] > ' ' {
			if j.s[i] == '"' {
				// Scans through the string.
				ini := i
				for i = i + 1; i < len(j.s); i++ {
					if j.s[i] == '"' && j.s[i-1] != '\\' {
						end := i + 1
						// Skip spaces.
						for i = i + 1; i < len(j.s) && j.s[i] <= ' '; i++ {
						}
						// Emits if not a key.
						if j.s[i] == ':' {
							o.WriteString(j.s[ini:end].String())
						} else {
							o.WriteString(m(New(j.s[ini:end].String())).String())
						}
						o.WriteByte(j.s[i])
						break
					}
				}
			} else if j.s[i] == '{' || j.s[i] == '}' || j.s[i] == ',' || j.s[i] == ':' || j.s[i] == '[' || j.s[i] == ']' {
				o.WriteByte(j.s[i])
			} else {
				// Scans through anything until these characters.
				ini := i
				for ; i < len(j.s); i++ {
					if j.s[i] == ',' || j.s[i] == '}' || j.s[i] == ']' {
						o.WriteString(m(New(j.s[ini:i].String())).String())
						o.WriteByte(j.s[i])
						break
					}
					if j.s[i] == ' ' {
						o.WriteString(m(New(j.s[ini:i].String())).String())
						break
					}
					if i == len(j.s)-1 {
						o.WriteString(m(New(j.s.String())).String())
					}
				}
			}
		}
	}
	return New(o.String())
}

// IterateKeysValues iterates over the keys and values of
// a valid Json consecutively and apply a map function to
// transform each emitted value.
func (j Json) IterateKeysValues(m func(Json) Json) Json {
	var o strings.Builder
	o.Grow(len(j.s))
	for i := 0; i < len(j.s); i++ {
		if j.s[i] > ' ' {
			if j.s[i] == '"' {
				// Scans through the string.
				ini := i
				for i = i + 1; i < len(j.s); i++ {
					if j.s[i] == '"' && j.s[i-1] != '\\' {
						o.WriteString(m(New(j.s[ini : i+1].String())).String())
						break
					}
				}
			} else if j.s[i] == '{' || j.s[i] == '}' || j.s[i] == ',' || j.s[i] == ':' || j.s[i] == '[' || j.s[i] == ']' {
				o.WriteByte(j.s[i])
			} else {
				// Scans through anything until these characters.
				ini := i
				for ; i < len(j.s); i++ {
					if j.s[i] == ',' || j.s[i] == '}' || j.s[i] == ']' {
						o.WriteString(m(New(j.s[ini:i].String())).String())
						o.WriteByte(j.s[i])
						break
					}
					if j.s[i] == ' ' {
						o.WriteString(m(New(j.s[ini:i].String())).String())
						break
					}
					if i == len(j.s)-1 {
						o.WriteString(m(New(j.s.String())).String())
					}
				}
			}
		}
	}
	return New(o.String())
}

// Iterate iterates over the keys and values of a valid Json
// and applies a map function to transform both at once.
func (j Json) Iterate(m func(k, v Json) (Json, Json)) Json {
	var o strings.Builder
	o.Grow(len(j.s))
	for j.s.More() {
		c := j.s.Curr()
		if c > ' ' {
			// Is a string?
			if ini := j.s.Mark(); j.s.UtilMatchString('"') {
				str := j.s.Token(ini)
				j.s.MatchWhileByteLTE(' ')
				// Is a key?
				if j.s.MatchByte(':') {
					j.s.MatchWhileByteLTE(' ')
					// Is a key of an object or array? Emit only the key.
					if j.s.EqualByte('{') || j.s.EqualByte('[') {
						k, _ := m(New(str), New(""))
						o.WriteString(k.String())
						o.WriteByte(':')
						continue
					}
					// Is a key of a value (string or anything else)? Emit both key and value.
					if ini := j.s.Mark(); j.s.UtilMatchString('"') || j.s.MatchUntilAnyByte4(',', '}', ']', ' ') {
						val := j.s.Token(ini)
						k, v := m(New(str), New(val))
						o.WriteString(k.String())
						o.WriteByte(':')
						o.WriteString(v.String())
					}
				} else {
					// Not a key. Emit as a value.
					_, v := m(New(""), New(str))
					o.WriteString(v.String())
				}
				continue
			}
			if c == '{' || c == '}' || c == ',' || c == ':' || c == '[' || c == ']' {
				o.WriteByte(c)
			} else {
				// Gets anything and emit it as a value.
				if ini := j.s.Mark(); j.s.MatchUntilAnyByte5(',', ' ', '}', ']', 0) {
					val := j.s.Token(ini)
					_, v := m(New(""), New(val))
					o.WriteString(v.String())
					continue
				}
			}
		}
		j.s.Next()
	}
	return New(o.String())
}

func (j Json) Get(keyOrIndex string) (r Json) {
	f := func(k, v Json) bool {
		if k.TrimKey() == keyOrIndex {
			r = v
			return true
		}
		return false
	}
	if j.IsObject() {
		j.ForEachKeyVal(f)
	} else {
		j.ForEach(f)
	}
	return r
}

func (j Json) ForEachKeyVal(f func(k, v Json) bool) {
	if j.s.MatchByte('{') {
		for !j.s.MatchByte('}') {
			j.ws()

			ini := j.s.Mark()
			j.s.UtilMatchString('"')
			key := j.s.Token(ini)

			j.ws()
			j.s.Next() // Skip ':' character.
			j.ws()

			ini = j.s.Mark()

			if c := j.s.Curr(); c == '{' || c == '[' {
				j.s.UtilMatchOpenCloseCount(c, c+2, '"')
			} else if c == '"' {
				j.s.UtilMatchString('"')
			} else {
				j.s.MatchUntilAnyByte4(',', '}', ']', ' ')
			}

			if f(New(key), New(j.s.Token(ini))) {
				return
			}

			j.ws()
			j.s.MatchByte(',')
		}
	}
}

func (j Json) ForEach(f func(i, v Json) bool) {
	if j.s.MatchByte('[') {
		for i := 0; j.ws() && !j.s.MatchByte(']'); i++ {
			ini := j.s.Mark()
			if c := j.s.Curr(); c == '{' || c == '[' {
				j.s.UtilMatchOpenCloseCount(c, c+2, '"')
			} else if c == '"' {
				j.s.UtilMatchString('"')
			} else {
				j.s.MatchUntilAnyByte4(',', '}', ']', ' ')
			}
			if f(New(strconv.Itoa(i)), New(j.s.Token(ini))) {
				return
			}
			j.ws()
			j.s.MatchByte(',')
		}
	}
}

func (j *Json) ws() bool {
	j.s.MatchWhileByteLTE(' ')
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
	if j.IsObject() {
		j.ForEachKeyVal(func(k, v Json) bool {
			c++
			return false
		})
	} else {
		j.ForEach(func(i, v Json) bool {
			c++
			return false
		})
	}
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
	var o strings.Builder
	o.Grow(len(s))
	for i := 0; i < len(s); i++ {
		if s[i] > ' ' {
			if s[i] == '"' {
				ini := i
				for i = i + 1; i < len(s); i++ {
					if s[i] == '"' && s[i-1] != '\\' {
						o.WriteString(s[ini : i+1])
						break
					}
				}
			} else {
				o.WriteByte(s[i])
			}
		}
	}
	return New(o.String())
}

func (j Json) Prettify() Json {
	pad := "    "
	s := j.String()
	var o strings.Builder
	o.Grow(len(s) << 2)
	depth := 0
	for i := 0; i < len(s); i++ {
		if s[i] > ' ' {
			switch s[i] {
			case '"':
				ini := i
				for i = i + 1; i < len(s); i++ {
					if s[i] == '"' && s[i-1] != '\\' {
						o.WriteString(s[ini : i+1])
						break
					}
				}
			case ',':
				o.WriteString(",\n")
				for d := 0; d < depth; d++ {
					o.WriteString(pad)
				}
			case '{', '[':
				ini := i
				open, clos := s[i], s[i]+2
				// Skip spaces.
				for i = i + 1; i < len(s) && s[i] <= ' '; i++ {
				}
				// Is empty object or array?
				if s[i] == clos {
					if open == '{' {
						o.WriteString("{}")
					} else {
						o.WriteString("[]")
					}
					continue
				}
				// Nop, go back.
				i = ini
				o.WriteByte(open)
				o.WriteByte('\n')
				depth++
				for d := 0; d < depth; d++ {
					o.WriteString(pad)
				}
			case '}', ']':
				o.WriteString("\n")
				depth--
				for d := 0; d < depth; d++ {
					o.WriteString(pad)
				}
				o.WriteByte(s[i])
			case ':':
				o.WriteString(": ")
			default:
				o.WriteByte(s[i])
			}
		}
	}
	return New(o.String())
}

// #endregion Json
