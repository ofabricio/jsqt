// Package jsqt provides a language to query and transform JSON documents.
package jsqt

import (
	"encoding/json"
	"fmt"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unsafe"

	. "github.com/ofabricio/scanner" //lint:ignore ST1001 should not use dot imports
)

func Get(jsn, qry string) Json {
	return Compile(qry).Get(jsn)
}

func GetArgs(jsn, qry string, args []any) Json {
	return Compile(qry).GetArgs(jsn, args)
}

func (b ByteCode) Get(jsn string) Json {
	j := JSON(jsn)
	j.s.WS()
	q := Query{Root: j, code: b}
	return q.Parse(j)
}

func (b ByteCode) GetArgs(jsn string, args []any) Json {
	j := JSON(jsn)
	j.s.WS()
	q := Query{Root: j, code: b, args: args}
	return q.Parse(j)
}

func Compile(qry string) ByteCode {
	c := encoder{s: Scanner(qry), b: make(ByteCode, 0, len(qry))}
	c.Encode()
	return c.b
}

func JSON(jsn string) Json {
	return Json{Scanner(jsn)}
}

func Valid(jsn string) bool {
	return JSON(jsn).Valid()
}

// #region Query

// Query is the query language parser.
type Query struct {
	Root Json
	k, v Json
	save Json
	savs map[string]Json
	args []any
	defs map[string]ByteCode
	code ByteCode
}

func (q *Query) Parse(j Json) Json {
	return funcGet(q, j)
}

func (q *Query) ParseFunOrKey(j Json) Json {
	if q.code.isFun() {
		return q.ParseFun(j)
	}
	key := q.code.ident()
	if len(key) > 2 && key[0] == '"' {
		key = key[1 : len(key)-1]
	}
	return j.Get(key)
}

func (q *Query) ParseFunOrKeyOptional(j Json) Json {
	if q.MoreArg() {
		return q.ParseFunOrKey(j)
	}
	return j
}

func (q *Query) ParseFunOrRaw(j Json) Json {
	if q.code.isFun() {
		return q.ParseFun(j)
	}
	return JSON(q.code.ident())
}

func (q *Query) ParseFun(j Json) Json {
	qk, qv := q.k, q.v
	op := q.code.next()
	switch op {
	case opGet:
		j = funcGet(q, j)
	case opSet:
		j = funcSet(q, j)
	case opUpsert:
		j = funcUpsert(q, j)
	case opArr:
		j = funcArr(q, j)
	case opCollect:
		j = funcCollect(q, j)
	case opFirst:
		j = funcFirst(q, j)
	case opLast:
		j = funcLast(q, j)
	case opUnique:
		j = funcUnique(q, j)
	case opGroup:
		j = funcGroup(q, j)
	case opAt:
		j = funcAt(q, j)
	case opExpr:
		j = funcExpr(q, j)
	case opObj:
		j = funcObj(q, j)
	case opFlatten:
		j = funcFlatten(q, j)
	case opReverse:
		j = funcReverse(q, j)
	case opSort:
		j = funcSort(q, j)
	case opSlice:
		j = funcSlice(q, j)
	case opReduce:
		j = funcReduce(q, j)
	case opChunk:
		j = funcChunk(q, j)
	case opPartition:
		j = funcPartition(q, j)
	case opMin:
		j = funcMin(q, j)
	case opMax:
		j = funcMax(q, j)
	case opArg:
		j = funcArg(q, j)
	case opPick:
		j = funcPick(q, j)
	case opPluck:
		j = funcPluck(q, j)
	case opUnwind:
		j = funcUnwind(q, j)
	case opTranspose:
		j = funcTranspose(q, j)
	case opIf:
		j = funcIf(q, j)
	case opEither:
		j = funcEither(q, j)
	case opAnd:
		j = funcAnd(q, j)
	case opOr:
		j = funcOr(q, j)
	case opNot:
		j = funcNot(q, j)
	case opIn:
		j = funcIN(q, j)
	case opEQ:
		j = funcEQ(q, j)
	case opNEQ:
		j = funcNEQ(q, j)
	case opGT:
		j = funcGT(q, j)
	case opGTE:
		j = funcGTE(q, j)
	case opLT:
		j = funcLT(q, j)
	case opLTE:
		j = funcLTE(q, j)
	case opIsNum:
		j = funcIsNum(q, j)
	case opIsObj:
		j = funcIsObj(q, j)
	case opIsArr:
		j = funcIsArr(q, j)
	case opIsStr:
		j = funcIsStr(q, j)
	case opIsBool:
		j = funcIsBool(q, j)
	case opIsNull:
		j = funcIsNull(q, j)
	case opIsNully:
		j = funcIsNully(q, j)
	case opIsEmpty:
		j = funcIsEmpty(q, j)
	case opIsEmptyObj:
		j = funcIsEmptyObj(q, j)
	case opIsEmptyArr:
		j = funcIsEmptyArr(q, j)
	case opIsEmptyStr:
		j = funcIsEmptyStr(q, j)
	case opIsVoid:
		j = funcIsVoid(q, j)
	case opIsBlank:
		j = funcIsBlank(q, j)
	case opIsSome:
		j = funcIsSome(q, j)
	case opFalsy:
		j = funcIsFalsy(q, j)
	case opTruthy:
		j = funcIsTruthy(q, j)
	case opExists:
		j = funcExists(q, j)
	case opBool:
		j = funcBool(q, j)
	case opDefault:
		j = funcDefault(q, j)
	case opIterate:
		j = funcIterate(q, j)
	case opDebug:
		j = funcDebug(q, j)
	case opReplace:
		j = funcReplace(q, j)
	case opConcat:
		j = funcConcat(q, j)
	case opJoin:
		j = funcJoin(q, j)
	case opSplit:
		j = funcSplit(q, j)
	case opMatch:
		j = funcMatch(q, j)
	case opValid:
		j = funcValid(q, j)
	case opSave:
		j = funcSave(q, j)
	case opLoad:
		j = funcLoad(q, j)
	case opDef:
		j = funcDef(q, j)
	case opRoot:
		j = q.Root
	case opKey:
		j = q.k
	case opVal:
		j = q.v
	case opNothing:
		j = JSON("")
	case opMerge:
		j = j.Merge()
	case opSize:
		j = j.Size()
	case opRaw:
		j = q.ParseFunOrRaw(j)
	case opStringify:
		j = j.Stringify()
	case opJsonify:
		j = j.Jsonify()
	case opUgly:
		j = j.Uglify()
	case opPretty:
		j = j.Prettify()
	case opKeys:
		j = j.Keys()
	case opValues:
		j = j.Values()
	case opEntries:
		j = j.Entries()
	case opObjectify:
		j = j.Objectify()
	case opUpper:
		j = JSON(strings.ToUpper(j.String()))
	case opLower:
		j = JSON(strings.ToLower(j.String()))
	case opUnknown:
		j = funcUnknown(q, j)
	}
	if op > opThis {
		q.SkipArgs()
		q.code.match(opEOF)
	}
	q.k, q.v = qk, qv
	return j
}

func (q *Query) Match(flag string) bool {
	return q.code.matchIdent(flag)
}

func (q *Query) SkipArgs() {
	q.code.skipArgs()
}

func (q *Query) SkipArg() {
	q.code.skipArg()
}

func (q Query) MoreArg() bool {
	return q.code.more() && !q.code.equal(opEOF)
}

func callDefFun(q *Query, j Json, defFunMark ByteCode) Json {
	m := q.code
	q.code = defFunMark
	j = q.ParseFunOrKey(j)
	q.code = m
	return j
}

// #region Functions

func funcUnknown(q *Query, j Json) Json {
	name := q.code.ident()
	if defFunMark, ok := q.defs[name]; ok {
		j = callDefFun(q, j, defFunMark)
	}
	return j
}

func funcDef(q *Query, j Json) Json {
	fname := q.ParseFunOrRaw(j).String()
	m := q.code
	q.SkipArg()
	if q.defs == nil {
		q.defs = make(map[string]ByteCode)
	}
	q.defs[fname] = m
	return j
}

func funcGet(q *Query, j Json) Json {
	for q.MoreArg() {
		if q.Match("*") {
			j = funcCollect(q, j)
		} else {
			j = q.ParseFunOrKey(j)
		}
	}
	return j
}

func funcSet(q *Query, j Json) Json {
	insert := q.Match("-i")
	return funcSetInternal(q, j, insert)
}

func funcSetInternal(q *Query, j Json, insert bool) Json {
	if q.Match("-m") {
		j = q.ParseFunOrKey(j)
	}
	keyOrIndex := q.ParseFunOrRaw(j)
	if !q.MoreArg() {
		return keyOrIndex // The last item is the value.
	}
	if j.IsObject() {
		var o strings.Builder
		o.Grow(len(j.s) + 32)
		o.WriteString("{")
		found := false
		keyOrIdx := keyOrIndex.TrimQuote()
		j.ForEachKeyVal(func(k, v Json) bool {
			q.k, q.v = k, v
			if k.TrimQuote() == keyOrIdx {
				found = true
				if q.Match("-r") {
					k = q.ParseFunOrRaw(v)
				}
				if v = funcSetInternal(q, v, insert); v.Exists() {
					if o.Len() > 1 {
						o.WriteString(",")
					}
					o.WriteString(`"`)
					o.WriteString(k.TrimQuote())
					o.WriteString(`":`)
					o.WriteString(v.String())
				}
			} else {
				if o.Len() > 1 {
					o.WriteString(",")
				}
				o.WriteString(k.String())
				o.WriteString(":")
				o.WriteString(v.String())
			}
			return false
		})
		if !found && insert {
			if v := funcSetInternal(q, JSON("{}"), insert); v.Exists() {
				if o.Len() > 1 {
					o.WriteString(",")
				}
				if keyOrIndex.IsNumber() {
					o.Reset()
					o.WriteString(`[`)
					o.WriteString(v.String())
					o.WriteString(`]`)
					return JSON(o.String())
				} else {
					o.WriteString(`"`)
					o.WriteString(keyOrIndex.TrimQuote())
					o.WriteString(`":`)
					o.WriteString(v.String())
				}
			}
		}
		o.WriteString("}")
		return JSON(o.String())
	}
	if j.IsArray() {
		var o strings.Builder
		o.Grow(len(j.s) + 32)
		o.WriteString("[")
		found := false
		j.ForEach(func(i, v Json) bool {
			q.k, q.v = i, v
			if q.MoreArg() {
				if i.String() == keyOrIndex.String() {
					found = true
					v = funcSetInternal(q, v, insert)
				} else if keyOrIndex.s.EqualByte('*') {
					found = true
					m := q.code
					v = funcSetInternal(q, v, insert)
					q.code = m
				}
			}
			if v.Exists() {
				if o.Len() > 1 {
					o.WriteString(",")
				}
				o.WriteString(v.String())
			}
			return false
		})
		if !found && insert {
			if v := funcSetInternal(q, JSON("{}"), insert); v.Exists() {
				if o.Len() > 1 {
					o.WriteString(",")
				}
				o.WriteString(v.String())
			}
		}
		o.WriteString("]")
		return JSON(o.String())
	}
	return j
}

func funcArr(q *Query, j Json) Json {
	if q.Match("-t") {
		return funcArrTest(q, j)
	}
	var o strings.Builder
	o.Grow(64)
	o.WriteString("[")
	for q.MoreArg() {
		if v := q.ParseFunOrKey(j); v.Exists() {
			if o.Len() > 1 {
				o.WriteString(",")
			}
			o.WriteString(v.String())
		}
	}
	o.WriteString("]")
	return JSON(o.String())
}

func funcArrTest(q *Query, j Json) Json {
	if j.IsArray() {
		var ok bool
		m := q.code
		j.ForEach(func(i, v Json) bool {
			q.code = m
			ok = q.ParseFunOrKey(v).Exists()
			return !ok
		})
		if ok {
			return j
		}
	}
	return JSON("")
}

func funcObj(q *Query, j Json) Json {
	var o strings.Builder
	o.Grow(len(j.s) + 64)
	o.WriteString("{")
	writeKeyVals := func(j Json) {
		for q.MoreArg() {
			if k, v := q.ParseFunOrRaw(j), q.ParseFunOrKey(j); k.Exists() && v.Exists() {
				if o.Len() > 1 {
					o.WriteString(",")
				}
				o.WriteByte('"')
				o.WriteString(k.TrimQuote())
				o.WriteString(`":`)
				o.WriteString(v.String())
			}
		}
	}
	if q.Match("-i") {
		m := q.code
		j.ForEach(func(i, v Json) bool {
			q.k, q.v = i, v
			q.code = m
			writeKeyVals(v)
			return false
		})
		j.ForEachKeyVal(func(k, v Json) bool {
			q.k, q.v = k, v
			q.code = m
			writeKeyVals(v)
			return false
		})
	} else {
		writeKeyVals(j)
	}
	o.WriteString("}")
	return JSON(o.String())
}

func funcCollect(q *Query, j Json) Json {
	var o strings.Builder
	o.Grow(len(j.s))
	o.WriteString("[")
	m := q.code
	f := func(k, item Json) bool {
		q.k, q.v = k, item
		q.code = m
		if item = funcGet(q, item); item.Exists() {
			if o.Len() > 1 {
				o.WriteString(",")
			}
			o.WriteString(item.String())
		}
		return false
	}
	j.ForEachKeyVal(f)
	j.ForEach(f)
	o.WriteString("]")
	q.SkipArgs()
	return JSON(o.String())
}

func funcUnique(q *Query, j Json) Json {
	uniq := make(map[Json]bool)
	var o strings.Builder
	o.Grow(len(j.s))
	o.WriteString("[")
	m := q.code
	j.ForEach(func(i, item Json) bool {
		q.k, q.v = i, item
		q.code = m
		if item = funcGet(q, item); item.Exists() && !uniq[item] {
			uniq[item] = true
			if o.Len() > 1 {
				o.WriteString(",")
			}
			o.WriteString(item.String())
		}
		return false
	})
	o.WriteString("]")
	return JSON(o.String())
}

func funcFirst(q *Query, j Json) Json {
	var first Json
	m := q.code
	j.ForEach(func(i, item Json) bool {
		q.k, q.v = i, item
		q.code = m
		first = funcGet(q, item)
		return first.Exists()
	})
	return first
}

func funcLast(q *Query, j Json) Json {
	var last Json
	m := q.code
	j.ForEach(func(i, item Json) bool {
		q.k, q.v = i, item
		q.code = m
		if item = funcGet(q, item); item.Exists() {
			last = item
		}
		return false
	})
	return last
}

func funcFlatten(q *Query, j Json) Json {
	m := q.code
	if q.Match("-k") {
		if j.IsObject() {
			var o strings.Builder
			o.Grow(len(j.s))
			o.WriteString("{")
			m := q.code
			j.ForEachKeyVal(func(k, v Json) bool {
				if o.Len() > 1 {
					o.WriteString(",")
				}
				found := false
				q.code = m
				for q.MoreArg() {
					if q.ParseFunOrRaw(v).TrimQuote() == k.TrimQuote() {
						found = true
						break
					}
				}
				if vv := v.String(); found && v.IsObject() {
					o.WriteString(vv[1 : len(vv)-1])
				} else {
					o.WriteString(k.String())
					o.WriteString(":")
					o.WriteString(vv)
				}
				return false
			})
			o.WriteString("}")
			return JSON(o.String())
		} else if j.IsArray() {
			var o strings.Builder
			o.Grow(len(j.s))
			o.WriteString("[")
			j.ForEach(func(i, v Json) bool {
				q.code = m
				if v = funcFlatten(q, v); v.Exists() {
					if o.Len() > 1 {
						o.WriteString(",")
					}
					o.WriteString(v.String())
				}
				return false
			})
			o.WriteString("]")
			return JSON(o.String())
		}
	}
	depth := -1
	if v := q.ParseFunOrRaw(j); v.Exists() {
		depth = v.Int()
	}
	return j.Flatten(depth)
}

func funcSlice(q *Query, j Json) Json {
	if j.IsArray() {
		ini := q.ParseFunOrRaw(j).Int()
		end := q.ParseFunOrRaw(j).Int()
		if ini < 0 || end < 0 {
			size := j.Size().Int()
			if ini < 0 {
				ini = size + ini
			}
			if end < 0 {
				end = size + end
			}
		}
		var o strings.Builder
		o.Grow(len(j.s))
		o.WriteString("[")
		c := 0
		j.ForEach(func(i, v Json) bool {
			if c >= ini && (c < end || end == 0) {
				if o.Len() > 1 {
					o.WriteString(",")
				}
				o.WriteString(v.String())
			}
			c++
			return false
		})
		o.WriteString("]")
		return JSON(o.String())
	}
	return j
}

func funcAt(q *Query, j Json) Json {
	if j.IsArray() {
		at := q.ParseFunOrRaw(j)
		j.ForEach(func(i, v Json) bool {
			if i == at {
				j = v
				return true
			}
			return false
		})
	}
	return j
}

func funcReduce(q *Query, j Json) Json {
	acc := q.ParseFunOrRaw(j)
	m := q.code
	f := func(i, v Json) bool {
		q.k, q.v = i, acc
		q.code = m
		acc = q.ParseFunOrKey(v)
		return false
	}
	j.ForEachKeyVal(f)
	j.ForEach(f)
	return acc
}

func funcChunk(q *Query, j Json) Json {
	size := q.ParseFunOrRaw(j).Int()
	if size == 0 {
		size = 1
	}
	var o strings.Builder
	o.Grow(len(j.s) + 32)
	o.WriteString("[")
	c := 0
	j.ForEach(func(i, v Json) bool {
		q.k, q.v = i, v
		if o.Len() > 1 {
			o.WriteString(",")
		}
		mod := c % size
		if mod == 0 {
			o.WriteString("[")
		}
		o.WriteString(v.String())
		if mod == size-1 {
			o.WriteString("]")
		}
		c++
		return false
	})
	if c%size != 0 {
		o.WriteString("]")
	}
	o.WriteString("]")
	return JSON(o.String())
}

func funcPartition(q *Query, j Json) Json {
	falsy := make([]string, 0, 16)
	var o strings.Builder
	o.Grow(len(j.s) + 5)
	o.WriteString("[[")
	m := q.code
	j.ForEach(func(i, v Json) bool {
		q.k, q.v = i, v
		q.code = m
		if q.ParseFunOrKey(v).Exists() {
			if o.Len() > 2 {
				o.WriteString(",")
			}
			o.WriteString(v.String())
		} else {
			falsy = append(falsy, v.String())
		}
		return false
	})
	o.WriteString("],[")
	for i, v := range falsy {
		if i > 0 {
			o.WriteString(",")
		}
		o.WriteString(v)
	}
	o.WriteString("]]")
	return JSON(o.String())
}

func funcMin(q *Query, j Json) Json {
	var min Json
	ini := q.code
	j.ForEach(func(i, item Json) bool {
		q.k, q.v = i, item
		q.code = ini
		if item = funcGet(q, item); item.Exists() {
			if i.String() == "0" || item.LT(min) {
				min = item
			}
		}
		return false
	})
	return min
}

func funcMax(q *Query, j Json) Json {
	var max Json
	ini := q.code
	j.ForEach(func(i, item Json) bool {
		q.k, q.v = i, item
		q.code = ini
		if item = funcGet(q, item); item.Exists() {
			if i.String() == "0" || item.GT(max) {
				max = item
			}
		}
		return false
	})
	return max
}

func funcGroup(q *Query, j Json) Json {
	group := make(map[Json][]Json, 16)
	groupOrder := make([]Json, 0, len(group))
	m := q.code
	j.ForEach(func(i, item Json) bool {
		q.k, q.v = i, item
		q.code = m
		if g, v := q.ParseFunOrKey(item), q.ParseFunOrKey(item); g.Exists() && v.Exists() {
			if _, ok := group[g]; !ok {
				groupOrder = append(groupOrder, g)
			}
			group[g] = append(group[g], v)
		}
		return false
	})
	var o strings.Builder
	o.Grow(len(j.s))
	if q.Match("-a") {
		keyName := "key"
		valName := "values"
		if q.MoreArg() {
			keyName = q.ParseFunOrRaw(j).TrimQuote()
			valName = q.ParseFunOrRaw(j).TrimQuote()
		}
		o.WriteString("[")
		for _, gkey := range groupOrder {
			g := group[gkey]
			if o.Len() > 1 {
				o.WriteString(",")
			}
			o.WriteString(`{"`)
			o.WriteString(keyName)
			o.WriteString(`":`)
			o.WriteString(gkey.String())
			o.WriteString(`,"`)
			o.WriteString(valName)
			o.WriteString(`":[`)
			for i, v := range g {
				if i > 0 {
					o.WriteString(",")
				}
				o.WriteString(v.String())
			}
			o.WriteString("]")
			o.WriteString("}")
		}
		o.WriteString("]")
	} else {
		o.WriteString("{")
		for _, gkey := range groupOrder {
			g := group[gkey]
			if o.Len() > 1 {
				o.WriteString(",")
			}
			o.WriteString(`"`)
			o.WriteString(gkey.TrimQuote())
			o.WriteString(`":[`)
			for i, v := range g {
				if i > 0 {
					o.WriteString(",")
				}
				o.WriteString(v.String())
			}
			o.WriteString("]")
		}
		o.WriteString("}")
	}
	return JSON(o.String())
}

func funcUpsert(q *Query, j Json) Json {
	if j.IsObject() {
		done := make(map[string]bool)
		var o strings.Builder
		o.Grow(len(j.s))
		o.WriteString("{")
		for q.MoreArg() {
			if k, v := q.ParseFunOrRaw(j), q.ParseFunOrRaw(j); k.Exists() {
				key := k.TrimQuote()
				if v.Exists() {
					if o.Len() > 1 {
						o.WriteString(",")
					}
					o.WriteByte('"')
					o.WriteString(key)
					o.WriteString(`":`)
					o.WriteString(v.String())
				}
				done[key] = true
			}
		}
		j.ForEachKeyVal(func(k, v Json) bool {
			if key := k.TrimQuote(); !done[key] {
				if o.Len() > 1 {
					o.WriteString(",")
				}
				o.WriteByte('"')
				o.WriteString(key)
				o.WriteString(`":`)
				o.WriteString(v.String())
			}
			return false
		})
		o.WriteString("}")
		return JSON(o.String())
	}
	return j
}

func funcDefault(q *Query, j Json) Json {
	if j.Exists() {
		return j
	}
	return q.ParseFunOrRaw(j)
}

func funcIterate(q *Query, j Json) Json {
	if q.Match("-c") {
		return funcIterateCollect(q, j)
	}
	if q.Match("-f") {
		return funcIterateFast(q, j)
	}
	if q.Match("-kv") {
		return funcIterateKeysValues(q, j)
	}
	if q.Match("-k") {
		return funcIterateKeys(q, j)
	}
	if q.Match("-v") {
		return funcIterateValues(q, j)
	}
	return funcIterateAll(q, j)
}

func funcIterateCollect(q *Query, j Json) Json {
	includeRoot := q.Match("-r")
	depth := 0
	if q.Match("-d") {
		depth = q.ParseFunOrRaw(j).Int()
	}
	var o strings.Builder
	o.Grow(len(j.s))
	o.WriteString("[")
	m := q.code
	j.Iterator(depth, func(k, v Json) {
		q.k, q.v = k, v
		if !includeRoot && k.IsNull() {
			return
		}
		q.code = m
		if v = q.ParseFun(v); v.Exists() {
			if o.Len() > 1 {
				o.WriteString(",")
			}
			o.WriteString(v.String())
		}
	})
	o.WriteString("]")
	return JSON(o.String())
}

func funcIterateAll(q *Query, j Json) Json {
	includeRoot := q.Match("-r")
	depth := 0
	if q.Match("-d") {
		depth = q.ParseFunOrRaw(j).Int()
	}
	m := q.code
	return j.Iterate(depth, func(k, v Json) (Json, Json) {
		q.k, q.v = k, v
		if !includeRoot && k.IsNull() {
			return k, v
		}
		q.code = m
		k = q.ParseFunOrKey(k)
		v = q.ParseFunOrKey(v)
		return k, v
	})
}

func funcIterateFast(q *Query, j Json) Json {
	m := q.code
	return j.IterateFast(func(k, v Json) (Json, Json) {
		q.k, q.v = k, v
		q.code = m
		k = q.ParseFunOrKey(k)
		v = q.ParseFunOrKey(v)
		return k, v
	})
}

func funcIterateKeys(q *Query, j Json) Json {
	m := q.code
	return j.IterateKeys(func(k Json) Json {
		q.k = k
		q.code = m
		return q.ParseFunOrKey(k)
	})
}

func funcIterateValues(q *Query, j Json) Json {
	m := q.code
	return j.IterateValues(func(v Json) Json {
		q.v = v
		q.code = m
		return q.ParseFunOrKey(v)
	})
}

func funcIterateKeysValues(q *Query, j Json) Json {
	m := q.code
	return j.IterateKeysValues(func(kv Json) Json {
		q.k, q.v = kv, kv
		q.code = m
		return q.ParseFunOrKey(kv)
	})
}

func funcIf(q *Query, j Json) Json {
	for q.MoreArg() {
		not := q.Match("-n")
		condOrElse := q.ParseFunOrKey(j)
		isCond := q.MoreArg()
		if isCond {
			if condOrElse.Exists() != not {
				return q.ParseFunOrRaw(j)
			}
			q.SkipArg() // Skip "Then".
		} else {
			return condOrElse
		}
	}
	return j
}

func funcEither(q *Query, j Json) Json {
	v := q.ParseFunOrKey(j)
	for (v.IsNully() || !v.Exists()) && q.MoreArg() {
		v = q.ParseFunOrKey(j)
	}
	return v
}

func funcArg(q *Query, j Json) Json {
	arg := q.ParseFunOrRaw(j)
	val := q.args[arg.Int()]
	if f, ok := val.(func(Json) Json); ok {
		return f(j)
	}
	jsn, _ := json.Marshal(val) // I think this is cheating.
	return JSON(string(jsn))
}

func funcMatch(q *Query, j Json) Json {
	// Match a key. Returns the matched key.
	if q.Match("-kk") {
		if q.Match("-p") {
			return j.GetPrefixKey(q.ParseFunOrRaw(j).TrimQuote())
		}
		if q.Match("-s") {
			return j.GetSuffixKey(q.ParseFunOrRaw(j).TrimQuote())
		}
		if q.Match("-r") {
			return j.GetRegexKey(q.ParseFunOrRaw(j).TrimQuote())
		}
		return j.GetKey(q.ParseFunOrRaw(j).TrimQuote())
	}
	// Match a key. The context must be an object.
	if q.Match("-k") {
		if q.Match("-p") {
			return j.GetPrefix(q.ParseFunOrRaw(j).TrimQuote())
		}
		if q.Match("-s") {
			return j.GetSuffix(q.ParseFunOrRaw(j).TrimQuote())
		}
		if q.Match("-r") {
			return j.GetRegex(q.ParseFunOrRaw(j).TrimQuote())
		}
		return j.Get(q.ParseFunOrRaw(j).TrimQuote())
	}
	// Match a key value or a string.
	var v string
	if q.Match("-v") {
		v = q.ParseFunOrKey(j).TrimQuote()
	} else {
		v = j.TrimQuote()
	}
	switch {
	case q.Match("-p"):
		prefix := q.ParseFunOrRaw(j).TrimQuote()
		if strings.HasPrefix(v, prefix) {
			return j
		}
	case q.Match("-s"):
		suffix := q.ParseFunOrRaw(j).TrimQuote()
		if strings.HasSuffix(v, suffix) {
			return j
		}
	case q.Match("-r"):
		regex := q.ParseFunOrRaw(j).TrimQuote()
		if ok, _ := regexp.MatchString(regex, v); ok {
			return j
		}
	default:
		exact := q.ParseFunOrRaw(j).TrimQuote()
		if exact == v {
			return j
		}
	}
	return JSON("")
}

func funcExpr(q *Query, j Json) Json {
	v := funcTerm(q, j).Float()
	for q.MoreArg() {
		if q.Match("+") {
			v = v + funcTerm(q, j).Float()
			continue
		}
		if q.Match("-") {
			v = v - funcTerm(q, j).Float()
			continue
		}
		break
	}
	return JSON(strconv.FormatFloat(v, 'f', -1, 64))
}

func funcTerm(q *Query, j Json) Json {
	var v float64 = 1
	if q.Match("-") {
		v = -1
	}
	v = v * q.ParseFunOrRaw(j).Float()
	for q.MoreArg() {
		if q.Match("*") {
			v = v * funcTerm(q, j).Float()
			continue
		}
		if q.Match("/") {
			v = v / funcTerm(q, j).Float()
			continue
		}
		if q.Match("%") {
			v = math.Mod(v, funcTerm(q, j).Float())
			continue
		}
		break
	}
	return JSON(strconv.FormatFloat(v, 'f', -1, 64))
}

func funcUnwind(q *Query, j Json) Json {
	var o strings.Builder
	o.Grow(len(j.s) << 1)
	o.WriteString("[")
	if j.IsObject() {
		key := q.ParseFunOrRaw(j).TrimQuote()
		ren := key
		if q.Match("-r") {
			ren = q.ParseFunOrRaw(j).TrimQuote()
		}
		val := j.Get(key)
		val.ForEach(func(i, item Json) bool {
			if o.Len() > 1 {
				o.WriteString(",")
			}
			o.WriteString("{")
			idx := 0
			j.ForEachKeyVal(func(k, v Json) bool {
				if idx > 0 {
					o.WriteString(",")
				}
				idx++
				if k.TrimQuote() == key {
					o.WriteString(`"`)
					o.WriteString(ren)
					o.WriteString(`":`)
					o.WriteString(item.String())
				} else {
					o.WriteString(k.String())
					o.WriteString(":")
					o.WriteString(v.String())
				}
				return false
			})
			o.WriteString("}")
			return false
		})
	} else if j.IsArray() {
		m := q.code
		j.ForEach(func(i, v Json) bool {
			q.code = m
			if unwinded := funcUnwind(q, v).Flatten(-1); unwinded.Exists() {
				if o.Len() > 1 {
					o.WriteString(",")
				}
				o.WriteString(unwinded.String())
			}
			return false
		})
	}
	o.WriteString("]")
	return JSON(o.String())
}

func funcTranspose(q *Query, j Json) Json {
	if j.IsObject() {
		keyAndArr := make([]Json, 0, 16)
		j.ForEachKeyVal(func(k, v Json) bool {
			if v.IsArray() {
				keyAndArr = append(keyAndArr, k, v.Flatten(-1))
			}
			return false
		})
		var o strings.Builder
		o.Grow(len(j.s) << 1)
		o.WriteString("[")
		for {
			wrote := 0
			for i := 0; i < len(keyAndArr)-1; i += 2 {

				key, arr := keyAndArr[i], keyAndArr[i+1]

				arr.s.WS()
				val := arr.s.TokenFor(arr.matchValue)
				arr.s.WS()
				arr.s.MatchByte(',')
				keyAndArr[i+1] = arr

				if len(val) == 0 {
					continue
				}

				if wrote > 0 {
					o.WriteString(",")
				} else {
					if o.Len() > 1 {
						o.WriteString(",")
					}
					o.WriteString("{")
				}
				o.WriteString(key.String())
				o.WriteString(":")
				o.WriteString(val)
				wrote++
			}
			if wrote == 0 {
				break
			}
			o.WriteString("}")
		}
		o.WriteString("]")
		return JSON(o.String())
	}
	if j.IsArray() {
		var o strings.Builder
		o.Grow(len(j.s))
		o.WriteString("{")
		j.Get("0").ForEachKeyVal(func(k, v Json) bool {
			if o.Len() > 1 {
				o.WriteString(",")
			}
			o.WriteString(k.String())
			o.WriteString(":[")
			j.ForEach(func(i, item Json) bool {
				if item = item.Get(k.TrimQuote()); item.Exists() {
					if i.String() != "0" {
						o.WriteString(",")
					}
					o.WriteString(item.String())
				}
				return false
			})
			o.WriteString("]")
			return false
		})
		o.WriteString("}")
		return JSON(o.String())
	}
	return j
}

func funcValid(q *Query, j Json) Json {
	if j = q.ParseFunOrKeyOptional(j); j.Valid() {
		return j
	}
	return JSON("")
}

func funcIsNum(q *Query, j Json) Json {
	v := q.ParseFunOrKeyOptional(j)
	if v.IsNumber() {
		return j
	}
	return JSON("")
}

func funcIsObj(q *Query, j Json) Json {
	v := q.ParseFunOrKeyOptional(j)
	if v.IsObject() {
		return j
	}
	return JSON("")
}

func funcIsArr(q *Query, j Json) Json {
	v := q.ParseFunOrKeyOptional(j)
	if v.IsArray() {
		return j
	}
	return JSON("")
}

func funcIsStr(q *Query, j Json) Json {
	v := q.ParseFunOrKeyOptional(j)
	if v.IsString() {
		return j
	}
	return JSON("")
}

func funcIsBool(q *Query, j Json) Json {
	v := q.ParseFunOrKeyOptional(j)
	if v.IsBool() {
		return j
	}
	return JSON("")
}

func funcIsNull(q *Query, j Json) Json {
	v := q.ParseFunOrKeyOptional(j)
	if v.IsNull() {
		return j
	}
	return JSON("")
}

func funcIsEmpty(q *Query, j Json) Json {
	v := q.ParseFunOrKeyOptional(j)
	if v.IsEmpty() {
		return j
	}
	return JSON("")
}

func funcIsEmptyObj(q *Query, j Json) Json {
	v := q.ParseFunOrKeyOptional(j)
	if v.IsEmptyObject() {
		return j
	}
	return JSON("")
}

func funcIsEmptyArr(q *Query, j Json) Json {
	v := q.ParseFunOrKeyOptional(j)
	if v.IsEmptyArray() {
		return j
	}
	return JSON("")
}

func funcIsEmptyStr(q *Query, j Json) Json {
	v := q.ParseFunOrKeyOptional(j)
	if v.IsEmptyString() {
		return j
	}
	return JSON("")
}

func funcIsTruthy(q *Query, j Json) Json {
	v := q.ParseFunOrKeyOptional(j)
	if v.IsTruthy() {
		return j
	}
	return JSON("")
}

func funcIsFalsy(q *Query, j Json) Json {
	v := q.ParseFunOrKeyOptional(j)
	if v.IsFalsy() {
		return j
	}
	return JSON("")
}

func funcIsSome(q *Query, j Json) Json {
	v := q.ParseFunOrKeyOptional(j)
	if v.IsSome() {
		return j
	}
	return JSON("")
}

func funcIsVoid(q *Query, j Json) Json {
	v := q.ParseFunOrKeyOptional(j)
	if v.IsVoid() {
		return j
	}
	return JSON("")
}

func funcIsNully(q *Query, j Json) Json {
	v := q.ParseFunOrKeyOptional(j)
	if v.IsNully() {
		return j
	}
	return JSON("")
}

func funcIsBlank(q *Query, j Json) Json {
	v := q.ParseFunOrKeyOptional(j)
	if v.IsBlank() {
		return j
	}
	return JSON("")
}

func funcExists(q *Query, j Json) Json {
	v := q.ParseFunOrKeyOptional(j)
	if v.Exists() {
		return j
	}
	return JSON("")
}

func funcIN(q *Query, j Json) Json {
	b := q.ParseFunOrRaw(j)
	a := q.ParseFunOrKeyOptional(j)
	if a.IN(b) {
		return j
	}
	return JSON("")
}

func funcEQ(q *Query, j Json) Json {
	b := q.ParseFunOrRaw(j)
	a := q.ParseFunOrKeyOptional(j)
	if a.EQ(b) {
		return j
	}
	return JSON("")
}

func funcNEQ(q *Query, j Json) Json {
	b := q.ParseFunOrRaw(j)
	a := q.ParseFunOrKeyOptional(j)
	if a.NEQ(b) {
		return j
	}
	return JSON("")
}

func funcGTE(q *Query, j Json) Json {
	b := q.ParseFunOrRaw(j)
	a := q.ParseFunOrKeyOptional(j)
	if a.GTE(b) {
		return j
	}
	return JSON("")
}

func funcLTE(q *Query, j Json) Json {
	b := q.ParseFunOrRaw(j)
	a := q.ParseFunOrKeyOptional(j)
	if a.LTE(b) {
		return j
	}
	return JSON("")
}

func funcGT(q *Query, j Json) Json {
	b := q.ParseFunOrRaw(j)
	a := q.ParseFunOrKeyOptional(j)
	if a.GT(b) {
		return j
	}
	return JSON("")
}

func funcLT(q *Query, j Json) Json {
	b := q.ParseFunOrRaw(j)
	a := q.ParseFunOrKeyOptional(j)
	if a.LT(b) {
		return j
	}
	return JSON("")
}

func funcOr(q *Query, j Json) Json {
	for q.MoreArg() {
		if v := q.ParseFunOrKey(j); v.Exists() {
			return j
		}
	}
	return JSON("")
}

func funcAnd(q *Query, j Json) Json {
	for q.MoreArg() {
		if v := q.ParseFunOrKey(j); !v.Exists() {
			return v
		}
	}
	return j
}

func funcNot(q *Query, j Json) Json {
	a := q.ParseFunOrKey(j)
	if a.Exists() {
		return JSON("")
	}
	return j
}

func funcBool(q *Query, j Json) Json {
	j = q.ParseFunOrKeyOptional(j)
	if j.Exists() {
		return JSON("true")
	}
	return JSON("false")
}

func funcDebug(q *Query, j Json) Json {
	msg := "debug"
	if q.MoreArg() {
		msg = q.ParseFunOrRaw(j).String()
	}
	fmt.Printf("[%s] %s\n", msg, j.String())
	return j
}

func funcReplace(q *Query, j Json) Json {
	old := q.ParseFunOrRaw(j)
	new := q.ParseFunOrRaw(j)
	if j.IsString() {
		return JSON(strings.ReplaceAll(j.String(), old.TrimQuote(), new.TrimQuote()))
	}
	return j
}

func funcJoin(q *Query, j Json) Json {
	sep := q.ParseFunOrRaw(j).Str()
	j = q.ParseFunOrKeyOptional(j)
	var o strings.Builder
	o.Grow(len(j.s))
	j.ForEach(func(i, v Json) bool {
		if o.Len() > 0 {
			o.WriteString(sep)
		}
		o.WriteString(v.Str())
		return false
	})
	return JSON(o.String()).Stringify()
}

func funcSplit(q *Query, j Json) Json {
	sep := q.ParseFunOrRaw(j).Str()
	j = q.ParseFunOrKeyOptional(j)
	var o strings.Builder
	o.Grow(len(j.s) + 32)
	o.WriteString("[")
	for i, v := range strings.Split(j.TrimQuote(), sep) {
		if i > 0 {
			o.WriteString(",")
		}
		o.WriteString(`"`)
		o.WriteString(v)
		o.WriteString(`"`)
	}
	o.WriteString("]")
	return JSON(o.String())
}

func funcConcat(q *Query, j Json) Json {
	var o strings.Builder
	o.Grow(32)
	for q.MoreArg() {
		if v := q.ParseFunOrKey(j); v.IsString() {
			o.WriteString(v.Jsonify().String())
		} else if v.IsNumber() || v.IsBool() || v.IsNull() {
			o.WriteString(v.String())
		}
	}
	return JSON(o.String()).Stringify()
}

func funcSort(q *Query, j Json) Json {
	asc := !q.Match("desc")
	key := q.MoreArg()
	if j.IsObject() {
		if asc {
			return Get(j.String(), "(entries) (sort 0) (objectify)")
		}
		return Get(j.String(), "(entries) (sort desc 0) (objectify)")
	}
	if j.IsArray() {
		var items []string
		j.ForEach(func(i, v Json) bool {
			items = append(items, v.String())
			return false
		})
		m := q.code
		sort.SliceStable(items, func(i, j int) bool {
			var a, b string
			if key {
				q.code = m
				a = q.ParseFunOrKey(JSON(items[i])).String()
				q.code = m
				b = q.ParseFunOrKey(JSON(items[j])).String()
			} else {
				a = items[i]
				b = items[j]
			}
			if asc {
				return JSON(a).LT(JSON(b))
			}
			return JSON(a).GT(JSON(b))
		})
		return JSON("[" + strings.Join(items, ",") + "]")
	}
	return j
}

func funcReverse(q *Query, j Json) Json {
	if j.IsArray() {
		var items []string
		j.ForEach(func(i, v Json) bool {
			items = append(items, v.String())
			return false
		})
		for a, b := 0, len(items)-1; a < b; a, b = a+1, b-1 {
			items[a], items[b] = items[b], items[a]
		}
		return JSON("[" + strings.Join(items, ",") + "]")
	}
	return j
}

func funcPick(q *Query, j Json) Json {
	if j.IsObject() {
		var o strings.Builder
		o.Grow(len(j.s) >> 1)
		o.WriteByte('{')
		for q.MoreArg() {
			key := q.ParseFunOrRaw(j).TrimQuote()
			v := j.Get(key)
			if q.Match("-r") {
				key = q.ParseFunOrRaw(v).String()
			}
			if q.Match("-m") {
				v = q.ParseFunOrKey(v)
			}
			if v.Exists() {
				if o.Len() > 1 {
					o.WriteByte(',')
				}
				o.WriteByte('"')
				o.WriteString(key)
				o.WriteString(`":`)
				o.WriteString(v.String())
			}
		}
		o.WriteByte('}')
		return JSON(o.String())
	}
	return j
}

func funcPluck(q *Query, j Json) Json {
	if j.IsObject() {
		var o strings.Builder
		o.Grow(len(j.s))
		o.WriteByte('{')
		m := q.code
		j.ForEachKeyVal(func(k, v Json) bool {
			for q.code = m; q.MoreArg(); {
				key := q.ParseFunOrRaw(v)
				if key.TrimQuote() == k.TrimQuote() {
					return false
				}
			}
			if o.Len() > 1 {
				o.WriteByte(',')
			}
			o.WriteString(k.String())
			o.WriteString(`:`)
			o.WriteString(v.String())
			return false
		})
		o.WriteByte('}')
		return JSON(o.String())
	}
	return j
}

func funcSave(q *Query, j Json) Json {
	if q.Match("-k") {
		if q.savs == nil {
			q.savs = make(map[string]Json, 4)
		}
		for q.MoreArg() {
			k := q.ParseFunOrRaw(j).TrimQuote()
			var v Json
			if q.Match("-v") {
				v = q.ParseFunOrKey(j)
			} else {
				v = j.Get(k)
			}
			q.savs[k] = v
		}
	} else if q.MoreArg() {
		q.save = q.ParseFunOrKey(j)
	} else {
		q.save = j
	}
	return j
}

func funcLoad(q *Query, j Json) Json {
	if q.MoreArg() {
		id := q.ParseFunOrRaw(j).String()
		return q.savs[id]
	}
	return q.save
}

// #endregion Functions

// #endregion Query

// #region Json

// Json represents a JSON document.
type Json struct {
	s Scanner
}

// String returns the raw JSON data.
func (j Json) String() string {
	return j.s.String()
}

// Bytes returns the raw JSON data.
func (j Json) Bytes() []byte {
	return j.s.Bytes()
}

// Stringify converts a JSON to a JSON string.
// Examples:
//
//	"Hello" -> "\"Hello\""
//	""      -> "\"\""
//	3       -> "3"
//	{}      -> "{}"
//	{ "hello": "world" } -> "{ \"hello\": \"world\" }"
//
// Stringify reverts Jsonify.
func (j Json) Stringify() Json {
	return JSON(strconv.Quote(j.String()))
}

// Jsonify converts a JSON string to a JSON.
// Examples:
//
//	"true" -> true
//	"3"    -> 3
//	"{}"   -> {}
//	"{ \"hello\": \"world\"}" -> { "hello": "world" }
//
// Jsonify reverts Stringify.
func (j Json) Jsonify() Json {
	v, _ := strconv.Unquote(j.String())
	return JSON(v)
}

// TrimQuote removes the quotes from an object key.
// Example: "name" -> name.
func (j Json) TrimQuote() string {
	v := j.String()
	if j.IsString() && len(j.s) > 1 {
		return v[1 : len(v)-1]
	}
	return v
}

// Str converts a JSON value to string.
func (j Json) Str() string {
	v := j.String()
	if j.IsString() {
		v, _ = strconv.Unquote(v)
	}
	return v
}

// Int converts a JSON number to int.
func (j Json) Int() int {
	v, _ := strconv.ParseInt(j.String(), 10, 0)
	return int(v)
}

// Int64 converts a JSON number to int64.
func (j Json) Int64() int64 {
	v, _ := strconv.ParseInt(j.String(), 10, 64)
	return v
}

// Uint64 converts a JSON number to uint64.
func (j Json) Uint64() uint64 {
	v, _ := strconv.ParseUint(j.String(), 10, 64)
	return v
}

// Float converts a JSON number to float.
func (j Json) Float() float64 {
	v, _ := strconv.ParseFloat(j.String(), 64)
	return v
}

// Bool converts a JSON boolean to bool.
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
	return j.s.Equal(`""`)
}

func (j Json) IsEmptyObject() bool {
	return j.s.MatchByte('{') && j.s.WS() && j.s.EqualByte('}')
}

func (j Json) IsEmptyArray() bool {
	return j.s.MatchByte('[') && j.s.WS() && j.s.EqualByte(']')
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
	return !j.IsFalsy() && j.Exists()
}

func (j Json) IsFalsy() bool {
	return j.IsEmptyObject() || j.IsEmptyArray() || j.IsEmptyString() ||
		j.IsFalse() || j.IsNull() || j.s.EqualByte('0')
}

func (j Json) IsSome() bool {
	return !j.IsNull() && j.Exists()
}

func (j Json) Exists() bool {
	return j.String() != ""
}

func (j Json) IN(arr Json) (yes bool) {
	arr.ForEach(func(i, v Json) bool {
		yes = j.EQ(v)
		return yes
	})
	return
}

func (j Json) EQ(b Json) bool {
	if j.IsNumber() && b.IsNumber() {
		na, nb := toFloat(j, b)
		return na == nb
	}
	return j.String() == b.String()
}

func (j Json) NEQ(b Json) bool {
	if j.IsNumber() && b.IsNumber() {
		na, nb := toFloat(j, b)
		return na != nb
	}
	return j.String() != b.String()
}

func (j Json) GT(b Json) bool {
	if j.IsNumber() && b.IsNumber() {
		na, nb := toFloat(j, b)
		return na > nb
	}
	return j.String() > b.String()
}

func (j Json) GTE(b Json) bool {
	if j.IsNumber() && b.IsNumber() {
		na, nb := toFloat(j, b)
		return na >= nb
	}
	return j.String() >= b.String()
}

func (j Json) LT(b Json) bool {
	if j.IsNumber() && b.IsNumber() {
		na, nb := toFloat(j, b)
		return na < nb
	}
	return j.String() < b.String()
}

func (j Json) LTE(b Json) bool {
	if j.IsNumber() && b.IsNumber() {
		na, nb := toFloat(j, b)
		return na <= nb
	}
	return j.String() <= b.String()
}

func toFloat(a, b Json) (float64, float64) {
	na, _ := strconv.ParseFloat(a.String(), 64)
	nb, _ := strconv.ParseFloat(b.String(), 64)
	return na, nb
}

func (j Json) Iterator(depth int, m func(k, v Json)) {
	if depth == 0 {
		depth = -2
	}
	j.iterator(depth, JSON("null"), m)
}

func (j Json) iterator(depth int, key Json, m func(k, v Json)) {
	if depth == -1 {
		return
	}
	m(key, j)
	j.ForEachKeyVal(func(k, v Json) bool {
		v.iterator(depth-1, k, m)
		return false
	})
	j.ForEach(func(i, v Json) bool {
		v.iterator(depth-1, i, m)
		return false
	})
}

func (j Json) Iterate(depth int, m func(k, v Json) (Json, Json)) Json {
	if depth == 0 {
		depth = -2
	}
	_, v := j.iterate(depth, JSON("null"), m)
	return v
}

func (j Json) iterate(depth int, k Json, m func(k, v Json) (Json, Json)) (Json, Json) {
	if depth == -1 {
		return k, j
	}
	if j.IsObject() {
		var o strings.Builder
		o.Grow(len(j.s))
		o.WriteString("{")
		j.ForEachKeyVal(func(k, v Json) bool {
			if k, v = v.iterate(depth-1, k, m); k.Exists() && v.Exists() {
				if o.Len() > 1 {
					o.WriteString(",")
				}
				o.WriteString(`"`)
				o.WriteString(k.TrimQuote())
				o.WriteString(`":`)
				o.WriteString(v.String())
			}
			return false
		})
		o.WriteString("}")
		return m(k, JSON(o.String()))
	}
	if j.IsArray() {
		var o strings.Builder
		o.Grow(len(j.s))
		o.WriteString("[")
		j.ForEach(func(k, v Json) bool {
			if k, v = v.iterate(depth-1, k, m); k.Exists() && v.Exists() {
				if o.Len() > 1 {
					o.WriteString(",")
				}
				o.WriteString(v.String())
			}
			return false
		})
		o.WriteString("]")
		return m(k, JSON(o.String()))
	}
	return m(k, j)
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
							o.WriteString(m(JSON(j.s[ini:end].String())).String())
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
	return JSON(o.String())
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
							o.WriteString(m(JSON(j.s[ini:end].String())).String())
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
						o.WriteString(m(JSON(j.s[ini:i].String())).String())
						o.WriteByte(j.s[i])
						break
					}
					if j.s[i] <= ' ' {
						o.WriteString(m(JSON(j.s[ini:i].String())).String())
						break
					}
					if i == len(j.s)-1 {
						o.WriteString(m(JSON(j.s.String())).String())
					}
				}
			}
		}
	}
	return JSON(o.String())
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
						o.WriteString(m(JSON(j.s[ini : i+1].String())).String())
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
						o.WriteString(m(JSON(j.s[ini:i].String())).String())
						o.WriteByte(j.s[i])
						break
					}
					if j.s[i] <= ' ' {
						o.WriteString(m(JSON(j.s[ini:i].String())).String())
						break
					}
					if i == len(j.s)-1 {
						o.WriteString(m(JSON(j.s.String())).String())
					}
				}
			}
		}
	}
	return JSON(o.String())
}

// IterateFast iterates over the keys and values of a valid Json
// and applies a map function to transform both at once.
func (j Json) IterateFast(m func(k, v Json) (Json, Json)) Json {
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
						v := "{}"
						if j.s.EqualByte('[') {
							v = "[]"
						}
						k, _ := m(JSON(str), JSON(v))
						o.WriteString(k.String())
						o.WriteByte(':')
						continue
					}
					// Is a key of a value (string or anything else)? Emit both key and value.
					if ini := j.s.Mark(); j.s.UtilMatchString('"') || j.s.MatchUntilLTEOr4(' ', ',', '}', ']', 0) {
						val := j.s.Token(ini)
						k, v := m(JSON(str), JSON(val))
						o.WriteString(k.String())
						o.WriteByte(':')
						o.WriteString(v.String())
					}
				} else {
					// Not a key. Emit as a value.
					_, v := m(JSON("null"), JSON(str))
					o.WriteString(v.String())
				}
				continue
			}
			if c == '{' || c == '}' || c == ',' || c == ':' || c == '[' || c == ']' {
				o.WriteByte(c)
			} else {
				// Gets anything and emit it as a value.
				if ini := j.s.Mark(); j.s.MatchUntilLTEOr4(' ', ',', '}', ']', 0) {
					val := j.s.Token(ini)
					_, v := m(JSON("null"), JSON(val))
					o.WriteString(v.String())
					continue
				}
			}
		}
		j.s.Next()
	}
	return JSON(o.String())
}

func (j Json) Get(keyOrIndex string) (r Json) {
	f := func(k, v Json) bool {
		if k.TrimQuote() == keyOrIndex {
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

func (j Json) GetKey(key string) (r Json) {
	j.ForEachKeyVal(func(k, v Json) bool {
		if k.TrimQuote() == key {
			r = k
			return true
		}
		return false
	})
	return r
}

func (j Json) GetPrefix(prefix string) (r Json) {
	j.ForEachKeyVal(func(k, v Json) bool {
		if strings.HasPrefix(k.TrimQuote(), prefix) {
			r = v
			return true
		}
		return false
	})
	return r
}

func (j Json) GetPrefixKey(prefix string) (r Json) {
	j.ForEachKeyVal(func(k, v Json) bool {
		if strings.HasPrefix(k.TrimQuote(), prefix) {
			r = k
			return true
		}
		return false
	})
	return r
}

func (j Json) GetSuffix(suffix string) (r Json) {
	j.ForEachKeyVal(func(k, v Json) bool {
		if strings.HasSuffix(k.TrimQuote(), suffix) {
			r = v
			return true
		}
		return false
	})
	return r
}

func (j Json) GetSuffixKey(suffix string) (r Json) {
	j.ForEachKeyVal(func(k, v Json) bool {
		if strings.HasSuffix(k.TrimQuote(), suffix) {
			r = k
			return true
		}
		return false
	})
	return r
}

func (j Json) GetRegex(pattern string) (r Json) {
	if j.IsObject() {
		j.ForEachKeyVal(func(k, v Json) bool {
			if ok, _ := regexp.MatchString(pattern, k.TrimQuote()); ok {
				r = v
				return true
			}
			return false
		})
	}
	return r
}

func (j Json) GetRegexKey(pattern string) (r Json) {
	if j.IsObject() {
		j.ForEachKeyVal(func(k, v Json) bool {
			if ok, _ := regexp.MatchString(pattern, k.TrimQuote()); ok {
				r = k
				return true
			}
			return false
		})
	}
	return r
}

func (j Json) ForEachKeyVal(f func(k, v Json) bool) {
	if j.s.MatchByte('{') {
		for j.s.WS() && !j.s.MatchByte('}') {

			ini := j.s.Mark()
			j.s.UtilMatchString('"')
			key := j.s.Token(ini)

			j.s.WS()
			j.s.Next() // Skip ':' character.
			j.s.WS()

			ini = j.s.Mark()

			if c := j.s.Curr(); c == '"' {
				j.s.UtilMatchString('"')
			} else if c == '{' {
				j.s.UtilMatchOpenCloseCount('{', '}', '"')
			} else if c == '[' {
				j.s.UtilMatchOpenCloseCount('[', ']', '"')
			} else {
				j.s.MatchUntilLTEOr4(' ', ',', '}', ']', 0) // TODO: no need for 0. Create MatchUntilLTEOr3.
			}

			if t := j.s.Token(ini); len(t) == 0 || f(JSON(key), JSON(t)) {
				return
			}

			j.s.WS()
			j.s.MatchByte(',')
		}
	}
}

func (j Json) ForEach(f func(i, v Json) bool) {
	if j.s.MatchByte('[') {
		for i := 0; j.s.WS() && !j.s.MatchByte(']'); i++ {
			ini := j.s.Mark()
			if c := j.s.Curr(); c == '{' || c == '[' {
				j.s.UtilMatchOpenCloseCount(c, c+2, '"')
			} else if c == '"' {
				j.s.UtilMatchString('"')
			} else {
				j.s.MatchUntilLTEOr4(' ', ',', '}', ']', 0) // TODO: no need for 0. Create MatchUntilLTEOr3.
			}
			if t := j.s.Token(ini); len(t) == 0 || f(JSON(strconv.Itoa(i)), JSON(t)) {
				return
			}
			j.s.WS()
			j.s.MatchByte(',')
		}
	}
}

func (j *Json) matchValue() bool {
	return j.s.UtilMatchString('"') ||
		j.s.UtilMatchOpenCloseCount('{', '}', '"') ||
		j.s.UtilMatchOpenCloseCount('[', ']', '"') ||
		j.s.MatchUntilLTEOr4(' ', ',', '}', ']', 0)
}

// Flatten flattens a JSON array. Depth is the depth
// level to flatten. Use depth == 0 to deep flatten.
// If depth == -1 the array value is simply trimmed:
// `[3, 4]` becomes `3, 4`.
func (j Json) Flatten(depth int) Json {
	if j.s.MatchByte('[') {
		if depth == -1 {
			v := j.String()
			return JSON(v[0 : len(v)-1])
		}
		var o strings.Builder
		o.Grow(len(j.s))
		o.WriteString("[")
		d := 0
		for j.s.WS() && j.s.More() {
			if j.s.MatchByte(',') {
				continue
			}
			if j.s.MatchByte(']') {
				d--
				continue
			}
			if (d < depth || depth <= 0) && j.s.MatchByte('[') {
				d++
				continue
			}
			if o.Len() > 1 {
				o.WriteByte(',')
			}
			m := j.s.Mark()
			j.matchValue()
			v := j.s.Token(m)
			o.WriteString(v)
		}
		o.WriteString("]")
		return JSON(o.String())
	}
	if j.s.MatchByte('{') {
		var o strings.Builder
		o.Grow(len(j.s))
		o.WriteString("{")
		d := 0
		for j.s.WS() && j.s.More() {

			if j.s.MatchByte(',') {
				continue
			}
			if j.s.MatchByte('}') {
				d--
				continue
			}

			m := j.s.Mark()
			j.s.UtilMatchString('"')
			k := j.s.Token(m)

			j.s.WS()
			j.s.MatchByte(':')
			j.s.WS()

			if (d < depth || depth <= 0) && j.s.MatchByte('{') {
				d++
				continue
			}

			m = j.s.Mark()
			j.matchValue()
			v := j.s.Token(m)

			if o.Len() > 1 {
				o.WriteByte(',')
			}
			o.WriteString(k)
			o.WriteString(":")
			o.WriteString(v)
		}
		o.WriteString("}")
		return JSON(o.String())
	}
	return j
}

func (j Json) Size() Json {
	if j.IsString() {
		return JSON(strconv.Itoa(len(j.String()) - 2))
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
	return JSON(strconv.Itoa(c))
}

func (j Json) Keys() Json {
	var o strings.Builder
	o.Grow(len(j.s) >> 1)
	o.WriteString("[")
	j.ForEachKeyVal(func(k, v Json) bool {
		if o.Len() > 1 {
			o.WriteString(",")
		}
		o.WriteString(k.String())
		return false
	})
	o.WriteString("]")
	return JSON(o.String())
}

func (j Json) Values() Json {
	var o strings.Builder
	o.Grow(len(j.s))
	o.WriteString("[")
	j.ForEachKeyVal(func(k, v Json) bool {
		if o.Len() > 1 {
			o.WriteString(",")
		}
		o.WriteString(v.String())
		return false
	})
	o.WriteString("]")
	return JSON(o.String())
}

func (j Json) Entries() Json {
	var o strings.Builder
	o.Grow(len(j.s) + 64)
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
	return JSON(o.String())
}

func (j Json) Objectify() Json {
	var o strings.Builder
	o.Grow(len(j.s) + 128)
	o.WriteString("{")
	j.ForEach(func(i, v Json) bool {
		if o.Len() > 1 {
			o.WriteString(",")
		}
		o.WriteString(`"`)
		o.WriteString(v.Get("0").TrimQuote())
		o.WriteString(`":`)
		o.WriteString(v.Get("1").String())
		return false
	})
	o.WriteString("}")
	return JSON(o.String())
}

func (j Json) Merge() Json {
	done := make(map[string]bool)
	var o strings.Builder
	o.Grow(len(j.s))
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
				done[k.String()] = true
			}
			return false
		})
		return false
	})
	o.WriteString("}")
	return JSON(o.String())
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
	return JSON(o.String())
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
	return JSON(o.String())
}

func (j Json) Valid() bool {
	return j.valid() && !j.s.More()
}

func (j *Json) valid() bool {
	for j.s.More() {
		switch j.s.Curr() {
		case ' ', '\t', '\n', '\r':
		case '{':
			j.s.Next()
			for i := 0; j.s.WS() && j.s.More(); i++ {
				if j.s.MatchByte('}') {
					return true
				}
				if i > 0 && j.s.MatchByte(',') && j.s.WS() && !j.s.EqualByte('"') {
					return false
				}
				if !j.s.UtilMatchString('"') {
					return false
				}
				j.s.WS()
				if !j.s.MatchByte(':') {
					return false
				}
				if !j.valid() {
					return false
				}
			}
			return false
		case '[':
			j.s.Next()
			for i := 0; j.s.WS() && j.s.More(); i++ {
				if j.s.MatchByte(']') {
					return true
				}
				if i > 0 {
					if !(j.s.MatchByte(',') && j.valid()) {
						return false
					}
				} else if !j.valid() {
					return false
				}
			}
			return false
		case '"':
			return j.s.UtilMatchString('"')
		case 't':
			return j.s.Match("true")
		case 'f':
			return j.s.Match("false")
		case 'n':
			return j.s.Match("null")
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '-':
			return j.s.UtilMatchNumber()
		default:
			return false
		}
		j.s.Next()
	}
	return false
}

// #endregion Json

// #region Compiler

func (c *encoder) Encode() {
	for c.s.More() {
		c.encode()
	}
}

func (c *encoder) encode() {
	c.s.WS()
	if c.s.MatchByte('(') {
		name := c.Raw()
		op := c.FunOp(name)
		c.b = append(c.b, op)
		if op == opUnknown {
			c.b = append(c.b, opIdent, byte(len(name)))
			c.b = append(c.b, name...)
		}
		for c.s.WS() && c.s.More() && !c.s.MatchByte(')') {
			c.encode()
		}
		if op > opThis {
			c.b = append(c.b, opEOF)
		}
	} else if ident := c.Raw(); ident != "" {
		c.b = append(c.b, opIdent, byte(len(ident)))
		c.b = append(c.b, ident...)
	} else {
		c.s.Next()
	}
}

func (q *encoder) Raw() string {
	m := q.s.Mark()
	_ = q.s.UtilMatchString('"') ||
		q.s.UtilMatchOpenCloseCount('{', '}', '"') ||
		q.s.UtilMatchOpenCloseCount('[', ']', '"') ||
		q.s.MatchUntilLTEOr2(' ', ')', 0)
	return q.s.Token(m)
}

func (q *encoder) FunOp(name string) byte {
	switch name {
	case "get":
		return opGet
	case "set":
		return opSet
	case "obj":
		return opObj
	case "arr":
		return opArr
	case "raw":
		return opRaw
	case "collect":
		return opCollect
	case "unique":
		return opUnique
	case "first":
		return opFirst
	case "last":
		return opLast
	case "flatten":
		return opFlatten
	case "slice":
		return opSlice
	case "reduce":
		return opReduce
	case "chunk":
		return opChunk
	case "partition":
		return opPartition
	case "min":
		return opMin
	case "max":
		return opMax
	case "at":
		return opAt
	case "group":
		return opGroup
	case "upsert":
		return opUpsert
	case "size":
		return opSize
	case "default":
		return opDefault
	case "merge":
		return opMerge
	case "iterate":
		return opIterate
	case "is-num":
		return opIsNum
	case "is-obj":
		return opIsObj
	case "is-arr":
		return opIsArr
	case "is-str":
		return opIsStr
	case "is-bool":
		return opIsBool
	case "is-null":
		return opIsNull
	case "is-empty":
		return opIsEmpty
	case "is-empty-arr":
		return opIsEmptyArr
	case "is-empty-obj":
		return opIsEmptyObj
	case "is-empty-str":
		return opIsEmptyStr
	case "is-some":
		return opIsSome
	case "is-void":
		return opIsVoid
	case "is-blank":
		return opIsBlank
	case "is-nully":
		return opIsNully
	case "truthy":
		return opTruthy
	case "falsy":
		return opFalsy
	case "exists":
		return opExists
	case "if":
		return opIf
	case "either":
		return opEither
	case "root":
		return opRoot
	case "this":
		return opThis
	case "in":
		return opIn
	case "==":
		return opEQ
	case "!=":
		return opNEQ
	case ">=":
		return opGTE
	case "<=":
		return opLTE
	case ">":
		return opGT
	case "<":
		return opLT
	case "or":
		return opOr
	case "and":
		return opAnd
	case "not":
		return opNot
	case "bool":
		return opBool
	case "debug":
		return opDebug
	case "keys":
		return opKeys
	case "values":
		return opValues
	case "entries":
		return opEntries
	case "objectify":
		return opObjectify
	case "ugly":
		return opUgly
	case "pretty":
		return opPretty
	case "jsonify":
		return opJsonify
	case "stringify":
		return opStringify
	case "upper":
		return opUpper
	case "lower":
		return opLower
	case "replace":
		return opReplace
	case "join":
		return opJoin
	case "split":
		return opSplit
	case "concat":
		return opConcat
	case "sort":
		return opSort
	case "reverse":
		return opReverse
	case "pick":
		return opPick
	case "pluck":
		return opPluck
	case "def":
		return opDef
	case "save":
		return opSave
	case "load":
		return opLoad
	case "key":
		return opKey
	case "val":
		return opVal
	case "arg":
		return opArg
	case "match":
		return opMatch
	case "expr":
		return opExpr
	case "unwind":
		return opUnwind
	case "transpose":
		return opTranspose
	case "valid":
		return opValid
	case "nothing":
		return opNothing
	}
	return opUnknown
}

// encoder encodes the query language into bytecode.
type encoder struct {
	s Scanner
	b ByteCode
}

// ByteCode is the encoded query language.
type ByteCode []byte

func (b ByteCode) curr() byte {
	return b[0]
}

func (b *ByteCode) next() byte {
	c := b.curr()
	*b = (*b)[1:]
	return c
}

func (b ByteCode) equal(v byte) bool {
	return b.curr() == v
}

func (b *ByteCode) match(v byte) bool {
	if b.more() && b.equal(v) {
		*b = (*b)[1:]
		return true
	}
	return false
}

func (b *ByteCode) matchIdent(id string) bool {
	if b.equal(opIdent) {
		// 0 = op; 1 = size; 2 = identifier start; 2 + size = identifier end.
		size := (*b)[1]
		ident := (*b)[2 : 2+int(size)]
		if string(ident) == id {
			*b = (*b)[2+int(size):]
			return true
		}
	}
	return false
}

func (b ByteCode) more() bool {
	return len(b) > 0
}

func (b *ByteCode) ident() string {
	if b.match(opIdent) {
		size := b.next()
		ident := (*b)[:size]
		*b = (*b)[size:]
		return *(*string)(unsafe.Pointer(&ident))
	}
	return ""
}

func (b *ByteCode) skipArgs() {
	for b.more() && !b.equal(opEOF) {
		_ = b.skipIdent() ||
			b.skipFun()
	}
}

func (b *ByteCode) skipArg() {
	if b.more() && !b.equal(opEOF) {
		_ = b.skipIdent() ||
			b.skipFun()
	}
}

func (b *ByteCode) skipIdent() bool {
	if b.match(opIdent) {
		size := b.next()
		*b = (*b)[size:]
		return true
	}
	return false
}

func (b *ByteCode) skipFun() bool {
	if b.next() > opThis {
		b.skipArgs()
		return b.match(opEOF)
	}
	return true
}

func (b ByteCode) isFun() bool {
	return b.curr() > opIdent
}

const (
	opEOF byte = iota
	opIdent
	// Functions.
	opVal
	opKey
	opUgly
	opPretty
	opSize
	opTranspose
	opMerge
	opNothing
	opKeys
	opValues
	opEntries
	opObjectify
	opStringify
	opJsonify
	opUpper
	opLower
	opValid
	opReverse
	opJoin
	opRoot
	// If a function has no arguments put it above opThis, so that
	// it will not have an extra byte. Put it below otherwise.
	opThis
	opGet
	opDefault
	opReplace
	opSplit
	opDef
	opSlice
	opSet
	opObj
	opArr
	opRaw
	opAt
	opArg
	opCollect
	opReduce
	opChunk
	opPartition
	opUnique
	opFirst
	opLast
	opFlatten
	opMin
	opMax
	opGroup
	opUpsert
	opIterate
	opIsNum
	opIsObj
	opIsArr
	opIsStr
	opIsBool
	opIsNull
	opIsEmpty
	opIsEmptyArr
	opIsEmptyObj
	opIsEmptyStr
	opIsSome
	opIsVoid
	opIsBlank
	opIsNully
	opTruthy
	opFalsy
	opExists
	opIf
	opEither
	opIn
	opEQ
	opNEQ
	opGTE
	opLTE
	opGT
	opLT
	opOr
	opAnd
	opNot
	opBool
	opDebug
	opConcat
	opSort
	opPick
	opPluck
	opSave
	opLoad
	opMatch
	opExpr
	opUnwind
	opUnknown
)

// #endregion Compiler
