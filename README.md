# JSQT

This Go package provides a language to query and transform JSON documents.

[![build](https://github.com/ofabricio/jsqt/actions/workflows/go.yml/badge.svg)](https://github.com/ofabricio/jsqt/actions/workflows/go.yml)
[![playground](https://img.shields.io/badge/play-ground-blueviolet.svg?logo=github&labelColor=2b3137)](https://ofabricio.github.io/jsqt/)

### Example

```go
package main

import "github.com/ofabricio/jsqt"

func main() {

    j := `{ "data": { "message": "Hello" } }`

    v := jsqt.Get(j, `(get data message)`)

    fmt.Println(v) // "Hello"
}
```

The `jsqt.Get(jsn, qry)` function applies a query to a JSON.
Note that it only works on a valid JSON.

### Notes

- ⚠ Many functions are not consolidated yet. Watch for updates if you are using them,
  because they can change anytime as there is no official release yet.
- Don't open PR.

# Install

```
go get github.com/ofabricio/jsqt
```

# Query functions

Query functions have a name and arguments and live inside `()`.
For example, in `(get a b)` function, `get` is the query function name and `a` and `b` are its arguments.
Make sure to write a valid query since no validation is done during the parsing.

There are three types of function arguments:

- **Function** - These are functions, for example: `(get name)`, `(root)`.
  When the parser finds a function it calls it and uses its result as argument.
  A function and its arguments receive the current JSON context as input.
- **Key** - These are object keys or array indexes, for example: `name`, `"full name"`, `0`.
  When the parser finds a key it gets the value of the key and uses it as argument.
- **Raw** - These are anything you type, for example: `name`, `3`, `true`.
  When the parser finds a raw value it uses that value as argument.
  There is also a `(raw)` function that can be used to pass raw values as argument
  when a function accepts functions but not raw values.

## (get)

This function gets values from a JSON.

```clj
(get arg ...)
```

The argument list can be keys, functions or the `*` symbol.

The arguments work like a pipeline: the output of the first one is the input of the next one and so on.

The `*` symbol is to iterate over arrays or objects;
when the context is an array it makes `(get)` emit each array item to the next argument and collect the results into an array;
when the context is an object it makes `(get)` emit each key value to the next argument and collect the results into an array;
When `*` is used, two other functions become available:
[(key)](#key-val) that returns the current array index or object key and
[(val)](#key-val) that returns the current array item or object value.

This is one of the most important functions as its pipeline behavior is what allows passing
a context to other functions. Because of that the root function is also a get function.

**Example**

```go
j := `{
    "data": { "store name": "Grocery" },
    "tags": [
        { "name": "Fruit", "items": [{ "name": "Apple" }] },
        { "name": "Snack", "items": [{ "name": "Chips" }] },
        { "name": "Drink", "items": [{ "name": "Beers" }, { "name": "Wine" }] }
    ]
}`

a := jsqt.Get(j, `(get data)`)
b := jsqt.Get(j, `(get data "store name")`)
c := jsqt.Get(j, `(get tags 1)`)
d := jsqt.Get(j, `(get tags 1 name)`)
e := jsqt.Get(j, `(get tags * name)`)
f := jsqt.Get(j, `(get tags (size))`)
g := jsqt.Get(j, `(get tags * items * name)`)
h := jsqt.Get(j, `(get tags * items (get * name) (flatten))`)
i := jsqt.Get(j, `tags * (== "Drink" name) items 0 name`) // Can omit root (get).

fmt.Println(a) // { "store name": "Grocery" }
fmt.Println(b) // "Grocery"
fmt.Println(c) // { "name": "Snack", "items": [{ "name": "Chips" }] }
fmt.Println(d) // "Snack"
fmt.Println(e) // ["Fruit","Snack","Drink"]
fmt.Println(f) // 3
fmt.Println(g) // [["Apple"],["Chips"],["Beers","Wine"]]
fmt.Println(h) // ["Apple","Chips","Beers","Wine"]
fmt.Println(i) // ["Beers"]
```

**Example**

```go
a := jsqt.Get(`[ 3, 4 ]`, `(get * (obj (key) (val)))`)
b := jsqt.Get(`{ "a": 3, "b": 4 }`, `(get * (obj (key) (val)))`)

fmt.Println(a) // [{"0":3},{"1":4}]
fmt.Println(b) // [{"a":3},{"b":4}]
```

It is possible to nest `*` and still access both indexes with the help of the [(save)](#save-load) function.

**Example**

```go
j := `[ [ 3, 4 ], [ 5 ] ]`

a := jsqt.Get(j, `(get * (save (key)) * (concat (load) (raw "-") (key) (raw "-") (val)))`)

fmt.Println(a) // [["0-0-3","0-1-4"],["1-0-5"]]
```

## (collect)

This function is just an alias for `(get *)` for readability.

```clj
(collect arg ...)
```

## (obj)

This function creates a JSON object.

```clj
(obj key val ...)
(obj -i key val ...)
```

The arguments are pairs of JSON keys and values.
`key` can be a function or a raw value and `val` can be a function or a key.

**Example**

```go
j := `{ "loc": [ 63.4682, -20.1754 ] }`

a := jsqt.Get(j, `(obj lat (get loc 0) lng (get loc 1))`)
b := jsqt.Get(j, `(get loc (obj lat 0 lng 1))`) // Same as above.

fmt.Println(a) // {"lat":63.465,"lng":-20.178}
fmt.Println(b) // {"lat":63.465,"lng":-20.178}
```

Use `-i` to iterate over an array or object and create a new object out of it.
When `-i` is used, two other functions become available:
[(key)](#key-val) that returns an array index or object key; and
[(val)](#key-val) that returns an array item or key value.

**Example**

```go
j := `[{ "id": 6, "name": "June" }, { "id": 7, "name": "July" }]`

a := jsqt.Get(j, `(obj -i (get id) (val))`)

fmt.Println(a) // {"6":{ "id": 6, "name": "June" },"7":{ "id": 7, "name": "July" }}
```

**Example**

```go
j := `{ "id": 4, "name": "April" }`

a := jsqt.Get(j, `(obj -i (concat (raw "key_") (key)) (val))`)

fmt.Println(a) // {"key_id":4,"key_name":"April"}
```

## (arr)

This function creates a JSON array.

```clj
(arr item ...)
(arr -t cond)
```

Each argument becomes an array item. The arguments can be functions or keys.

**Example**

```go
j := `{ "a": 3, "b": 4 }`

a := jsqt.Get(j, `(arr a (raw 5) b)`)

fmt.Println(a) // [3,5,4]
```

Use `-t cond` to test each array item against a condition.
If all items succeed it returns the current context;
if any item fails it returns an empty context.
This can be used to validate a JSON schema. See [(and)](#and-or-not) for another example.

**Example**

```go
j := `[ 3, 4 ]`

a := jsqt.Get(j, `(arr -t (is-num))`)
b := jsqt.Get(j, `(arr -t (is-str))`)

fmt.Println(a) // [ 3, 4 ]
fmt.Println(b) //
```

## (raw)

This function creates a raw JSON value.

```clj
(raw value)
```

The `value` argument can be any valid JSON value.
Note that the argument is not validated, make sure to use a valid JSON value.

**Example**

```go
j := `{ "message": "Hello" }`

a := jsqt.Get(j, `(arr message (raw "World"))`)

fmt.Println(a) // ["Hello","World"]
```

## (set)

This function sets or removes values. It can also rename keys.

```clj
(set key ... val)
(set -i key ... val)
(set key -r newkey ... val)
(set key -m map ... val)
```

`key ...` is a list of keys or functions.

`val` (the last item of the list) is the value to be set and can be a function or a raw value.

By default `(set)` does not insert a field it does not find. Use `-i` flag to insert.

Use `-r newkey` after any key to rename it.

Use `-m map` after any key to apply a map function to its value.

The `*` symbol is to iterate on each array item.

**Example**

```go
j := `{"data":{"name":"Market"},"fruits":[{"name":"apple"},{"name":"grape"}]}`

a := jsqt.Get(j, `(set data name "Grocery")`)
b := jsqt.Get(j, `(set fruits (nothing))`)
c := jsqt.Get(j, `(set fruits 1 (nothing))`)
d := jsqt.Get(j, `(set fruits 0 name (raw "banana"))`)
e := jsqt.Get(j, `(set fruits * name "banana")`)
f := jsqt.Get(j, `(set -i data open true)`)
g := jsqt.Get(j, `(set data open true)`)
h := jsqt.Get(j, `(set fruits -r items * name -r value (this))`)

fmt.Println(a) // {"data":{"name":"Grocery"},"fruits":[{"name":"apple"},{"name":"grape"}]}
fmt.Println(b) // {"data":{"name":"Market"}}
fmt.Println(c) // {"data":{"name":"Market"},"fruits":[{"name":"apple"}]}
fmt.Println(d) // {"data":{"name":"Market"},"fruits":[{"name":"banana"},{"name":"grape"}]}
fmt.Println(e) // {"data":{"name":"Market"},"fruits":[{"name":"banana"},{"name":"banana"}]}
fmt.Println(f) // {"data":{"name":"Market","open":true},"fruits":[{"name":"apple"},{"name":"grape"}]}
fmt.Println(g) // {"data":{"name":"Market"},"fruits":[{"name":"apple"},{"name":"grape"}]}
fmt.Println(h) // {"data":{"name":"Market"},"items":[{"value":"apple"},{"value":"grape"}]}
```

**Example**

```go
j := `[{ "a": 3, "b": [{ "c": 4, "d": 5 }] }, { "a": 6, "b": [{ "c": 7, "d": 8 }] }]`

a := jsqt.Get(j, `(set * -m (pick b) b * -m (upsert x 0) c 9)`)

fmt.Println(a) // [{"b":[{"x":0,"c":9,"d":5}]},{"b":[{"x":0,"c":9,"d":8}]}]
```

[(key)](#key-val) and [(val)](#key-val) are available.

## (upsert)

This function creates, updates or removes object fields.

```clj
(upsert key val ...)
```

The arguments are pairs of JSON keys and values.
Both `key` and `val` can be a function or a raw value.
When `val` is an empty context the key is removed.

**Example**

```go
j := `{ "msg": "Hello", "author": "May", "deleted": false }`

a := jsqt.Get(j, `(upsert id 3 msg "World" deleted (nothing))`)

fmt.Println(a) // {"id":3,"msg":"World","author":"May"}
```

## (flatten)

This function flattens a JSON array or object.

```clj
(flatten)
(flatten depth)
(flatten -k key ...)
```

`(flatten)` just trims the `[]` out of a value. In some contexts this avoids allocations. Use with care.

`(flatten depth)` applies a proper flatten.
The `depth` argument is the depth level to flatten. Use `0` for a deep flatten.

`(flatten -k key ...)` flattens only the given keys. It also work with objects inside array.

**Example**

```go
a := jsqt.Get(`[[3], [4], [5]]`, `(collect (flatten))`)
b := jsqt.Get(`[3, [4], [[5]]]`, `(flatten 1)`)
c := jsqt.Get(`[3, [4], [[5]]]`, `(flatten 2)`)
d := jsqt.Get(`[3, [4], [[5]]]`, `(flatten 0)`)
e := jsqt.Get(`{"a":3,"b":{"c":4},"d":{"e":5}}`, `(flatten)`)
f := jsqt.Get(`{"a":3,"b":{"c":4},"d":{"e":5}}`, `(flatten -k b)`)
g := jsqt.Get(`[{"a":{"b":3}},{"a":{"b":4}}]`, `(flatten -k a)`)

fmt.Println(a) // [3,4,5]
fmt.Println(b) // [3,4,[5]]
fmt.Println(c) // [3,4,5]
fmt.Println(d) // [3,4,5]
fmt.Println(e) // {"a":3,"c":4,"e":5}
fmt.Println(f) // {"a":3,"c":4,"e":{"e":5}}
fmt.Println(g) // [{"b":3},{"b":4}]
```

## (size)

This function returns the size of a JSON array or object or the number of bytes in a string.

```clj
(size)
```

**Example**

```go
a := jsqt.Get(`[ 3, 7 ]`, `(size)`)
b := jsqt.Get(`{ "a": 3 }`, `(size)`)
c := jsqt.Get(`"Wisdom"`, `(size)`)

fmt.Println(a) // 2
fmt.Println(b) // 1
fmt.Println(c) // 6
```

## (first) (last)

These functions return the first or last item of a JSON array.

```clj
(first)
(first arg ...)
(last)
(last arg ...)
```

The arguments are optional and they have the same behavior as in `(get)` function,
except that it returns only the first or last item found.

**Example**

```go
j := `[{ "a": 1, "b": 3 }, { "a": 2, "b": 4 }, { "a": 3, "b": 4 }, { "a": 4, "b": 5 }]`

a := jsqt.Get(j, `(first)`)
b := jsqt.Get(j, `(last)`)
c := jsqt.Get(j, `(first (== 4 b))`)
d := jsqt.Get(j, `(last  (== 4 b))`)
e := jsqt.Get(j, `(first (== 4 b) a)`)
f := jsqt.Get(j, `(last  (== 4 b) a)`)

fmt.Println(a) // { "a": 1, "b": 3 }
fmt.Println(b) // { "a": 4, "b": 5 }
fmt.Println(c) // { "a": 2, "b": 4 }
fmt.Println(d) // { "a": 3, "b": 4 }
fmt.Println(e) // 2
fmt.Println(f) // 3
```

## (unique)

This function collects unique values from an array.

```clj
(unique)
(unique arg ...)
```

The arguments are optional and they have the same behavior as in `(get)` function,
except that it collects unique values.

**Example**

```go
a := jsqt.Get(`[3,4,3,4,5]`, `(unique)`)
b := jsqt.Get(`[{"a":3},{"a":3},{"a":4}]`, `(unique a)`)

fmt.Println(a) // [3,4,5]
fmt.Println(b) // [3,4]
```

## (slice)

This function returns a slice of a JSON array selected from start (inclusive) to end (exclusive).

```clj
(slice start end)
```

Both `start` and `end` can be a function or raw value. They can be negative. `end` is optional.

**Example**

```go
j := `[ "ant", "bear", "camel", "duck", "elephant" ]`

a := jsqt.Get(j, `(slice 2)`)
b := jsqt.Get(j, `(slice 2 4)`)
c := jsqt.Get(j, `(slice -2)`)
d := jsqt.Get(j, `(slice 2 -1)`)

fmt.Println(a) // ["camel","duck","elephant"]
fmt.Println(b) // ["camel","duck"]
fmt.Println(c) // ["duck","elephant"]
fmt.Println(d) // ["camel","duck"]
```

## (reverse)

This function reverses a JSON array.

```clj
(reverse)
```

**Example**

```go
a := jsqt.Get(`[3,7,2,4]`, `(reverse)`)

fmt.Println(a) // [4,2,7,3]
```

## (at)

This function returns an array item at an index.

```clj
(at index)
```

`index` can be a function or a raw value.

Note that there is no difference between `(at 0)` and `(get 0)`,
but the same is not true for `(at (key))` and `(get (key))`:
`at` will return the array item at the index returned by the function,
whereas `get` will return the value returned by `(key)`.

**Example**

```go
j := `[ "ant", "bear" ]`

a := jsqt.Get(j, `(at 1)`)
b := jsqt.Get(j, `(at  (raw 1))`)
c := jsqt.Get(j, `(get (raw 1))`) // Note the difference.

fmt.Println(a) // "bear"
fmt.Println(b) // "bear"
fmt.Println(c) // 1
```

## (reduce)

This function reduces an array or object to a single value.

```clj
(reduce ini map)
```

`ini` can be a function or a raw value and it is the initial value.

`map` must be a function; in this function [(val)](#key-val) is the accumulator (previous value),
[(key)](#key-val) is the array index or object key; and [(this)](#this) is the array item or key value.

**Example**

```go
a := jsqt.Get(`[ 3, 4, 5 ]`, `(reduce 2 (expr (val) + (this)))`)
b := jsqt.Get(`{ "a": 3, "b": 4, "c": 5 }`, `(reduce 2 (expr (val) + (this)))`)

fmt.Println(a) // 14
fmt.Println(b) // 14
```

## (chunk)

This function splits an array into groups the length of size.

```clj
(chunk size)
```

`size` can be a function or a raw value.

**Example**

```go
a := jsqt.Get(`[ 3, 4, 5, 6 ]`, `(chunk 2)`)
b := jsqt.Get(`[ 3, 4, 5, 6 ]`, `(chunk 3)`)

fmt.Println(a) // [[3,4],[5,6]]
fmt.Println(b) // [[3,4,5],[6]]
```

## (partition)

This function splits an array into two groups given a condition.

```clj
(partition cond)
```

`cond` can be a function or a key.

**Example**

```go
a := jsqt.Get(`[ 3, 30, 4, 40 ]`, `(partition (> 10))`)

fmt.Println(a) // [[30,40],[3,4]]
```

## (min) (max)

These functions return the min or max value.

```clj
(min)
(min arg ...)
(max)
(max arg ...)
```

The arguments are optional and they behave as in `(get)` function on each array item.

**Example**

```go
j := `[{ "a": 3 }, { "a": 1 }, { "a": 2 }]`

a := jsqt.Get(j, `(min a)`)
b := jsqt.Get(j, `(max a)`)

fmt.Println(a) // 1
fmt.Println(b) // 3
```

## (root)

This function returns the root JSON document.

```clj
(root)
```

Use it to access the root document from anywhere.

**Example**

```go
j := `[3,4]`

a := jsqt.Get(j, `(obj data (root))`)

fmt.Println(a) // {"data":[3,4]}
```

## (this)

This function returns the current JSON context value.

```clj
(this)
```

**Example**

```go
j := `[3,4]`

a := jsqt.Get(j, `(collect (obj value (this)))`)

fmt.Println(a) // [{"value":3},{"value":4}]
```

## (comparison)

These comparison functions return the current context if true or an empty context if false.

```clj
(== a b)
(!= a b)
(>= a b)
(<= a b)
(> a b)
(< a b)
(in a b)
```

`a` can be a function or a raw value.

`b` is optional and can be a function or a key.
When this argument is used, a comparison like `(> 33 age)` reads "is age greater than 33?".

**Example**

```go
j := `[{ "a": 3, "b": true }, { "a": 4, "b": false }, { "a": 5, "b": true }]`

a := jsqt.Get(j, `(collect (== true b) a)`)
b := jsqt.Get(j, `(collect (in [3, 4] a) b)`)

fmt.Println(a) // [3,5]
fmt.Println(b) // [true,false]
```

## (is-x)

These functions test a context for the corresponding value and
return the context if true or an empty context if false.

```clj
(is-obj)
(is-arr)
(is-num)
(is-str)
(is-bool)
(is-null)
(is-empty-arr)
(is-empty-obj)
(is-empty-str)
(is-empty)
(is-blank)
(is-nully)
(is-void)
(is-some)
(exists)
(truthy)
(falsy)
```

See [Truth Table](#truth-table) for the values that some functions above match.

All these functions accept an optional argument that can be a key or a function.

**Example**

```go
a := jsqt.Get(`3`, `(is-num)`)
b := jsqt.Get(`3`, `(is-str)`)
c := jsqt.Get(`{ "a": 3 }`, `(is-num a)`)
d := jsqt.Get(`{ "a": 3 }`, `(is-str a)`)

fmt.Println(a) // 3
fmt.Println(b) //
fmt.Println(c) // { "a": 3 }
fmt.Println(d) //
```

## (bool)

This function returns `true` if it gets a value or `false` if it gets an empty context.

```clj
(bool)
(bool arg)
```

`arg` is optional and can be a function or a key.

**Example**

```go
a := jsqt.Get(`[]`, `(is-arr) (bool)`)
b := jsqt.Get(`{}`, `(is-arr) (bool)`)
c := jsqt.Get(`{ "a": 3 }`, `(bool (is-arr a))`)
d := jsqt.Get(`{ "a": 3 }`, `(bool a)`)

fmt.Println(a) // true
fmt.Println(b) // false
fmt.Println(c) // false
fmt.Println(d) // true
```

## (and) (or) (not)

These functions apply AND, OR and NOT logic to its arguments.
They return the current context when true or an empty context otherwise.

```clj
(and a b ...)
(or  a b ...)
(not a)
```

The arguments can be functions or keys. Both `(and)` and `(or)` have variadic arguments.

**Example**

```go
j := `[{"a":3},{"a":4},{"a":5},{"a":6}]`

a := jsqt.Get(j, `(collect (or (< 4 a) (> 5 a)) a)`)
b := jsqt.Get(j, `(collect (and (>= 4 a) (<= 5 a)) a)`)
c := jsqt.Get(j, `(collect (not (<= 4 a)) a)`)

fmt.Println(a) // [3,6]
fmt.Println(b) // [4,5]
fmt.Println(c) // [5,6]
```

Note that `(and)` can be used to validate a JSON schema.

**Example**

```go
j := `
    {
        "name": "Chesterton",
        "age": 62,
        "books": [
            { "name": "Orthodoxy", "pages": 166 },
            { "name": "The Everlasting Man", "pages": 269 }
        ]
    }`

// Invalid schema because (size) is not 2, but 3.
a := jsqt.Get(j, `
    (and
        (is-str name)
        (== 62  age)
        (== 2 (size))
    )`)

// Valid schema.
b := jsqt.Get(j, `
    (and
        (is-str name)
        (is-num age)
        (get books (arr -t (and
            (is-str name)
            (is-num pages)
        )))
    )`)

fmt.Println(a) //
fmt.Println(b) // The result here is the input (b == j).
```

Note the help of [(arr -t)](#arr) function.

## (if)

This function works like a regular `if` or `switch/case`.

```clj
(if cond then ... else)
```

If `cond` is true `then` is executed, otherwise `else` is;
`else` is optional and returns `(this)` when omited.

`(if)` accepts variadic pairs of cond/then and in that case it works like a switch/case.

`cond` and `else` can be a function or a key;
`then` can be a function or a raw value.

Use `-n` to negate a condition.

**Example**

```go
j := `[ 3, {}, 4, [], 5 ]`

a := jsqt.Get(j, `(collect (if (is-obj) "obj"))`)
b := jsqt.Get(j, `(collect (if (is-obj) "obj" (is-arr) "arr"))`)
c := jsqt.Get(j, `(collect (if (is-obj) "obj" (is-arr) "arr" (raw "nop")))`)
d := jsqt.Get(j, `(collect (if -n (is-obj) "nop")`)

fmt.Println(a) // [3,"obj",4,[],5]
fmt.Println(b) // [3,"obj",4,"arr",5]
fmt.Println(c) // ["nop","obj","nop","arr","nop"]
fmt.Println(d) // ["nop",{},"nop","nop","nop"]
```

## (either)

This function returns the first argument value that is not [nully](#truth-table) nor empty context.

```clj
(either a b ...)
```

The arguments are keys or functions.

**Example**

```go
a := jsqt.Get(`{ "a": "A", "b": "" }`, `(either a b)`)
b := jsqt.Get(`{ "a": "", "b": "B" }`, `(either a b)`)

fmt.Println(a) // "A"
fmt.Println(b) // "B"
```

## (default)

This function returns a default value if it receives an empty context;
returns the received value otherwise.

```clj
(default val)
```

The `val` argument can be a function or a raw value.

**Example**

```go
j := `[{ "a": 3 }, { "b": 4 }, { "a": 5 }]`

a := jsqt.Get(j, `(collect a (default 0))`)

fmt.Println(a) // [3,0,5]
```

## (valid)

This function validates a JSON document.
If valid it returns the value it receives; if invalid it returns an empty context.

```clj
(valid)
(valid val)
```

`val` is optional and can be a function or a key.

**Example**

```go
a := jsqt.Get(`{ "a": 3 }`, `(valid)`)
b := jsqt.Get(`{ "a": 3, }`, `(valid)`)

fmt.Println(a) // { "a": 3 }
fmt.Println(b) //
```

Note that there is also the `jsqt.Valid(jsn)` function.

## (pick) (pluck)

These functions pick or pluck fields from a JSON object.

```clj
(pick key ...)
(pick key -m map ...)
(pick key -r newkey ...)
(pluck key ...)
```

The arguments are a list of object keys (raw values) or functions that return a key.

Use `-m` to enable a map function for a key.
This function receives the value of the key as context so that you can do something with it.
It can be used to deep pick fields.

Use `-r` to rename a key.

**Example**

```go
j := `{ "three": 3, "four": 4, "five": 5, "group": { "six": 6, "seven": 7 } }`

a := jsqt.Get(j, `(pick three five)`)
b := jsqt.Get(j, `(pluck three five group)`)
c := jsqt.Get(j, `(pick four group -m (pick seven))`)
d := jsqt.Get(j, `(pick three -r third five)`)

fmt.Println(a) // {"three":3,"five":5}
fmt.Println(b) // {"four":4}
fmt.Println(c) // {"four":4,"group":{"seven":7}}
fmt.Println(d) // {"third":3,"five":5}
```

## (merge)

This function merges an array of objects into one object.

```clj
(merge)
```

**Example**

```go
j := `[{ "a": 3 }, { "b": 4 }, { "c": 5 }]`

a := jsqt.Get(j, `(merge)`)

fmt.Println(a) // {"a":3,"b":4,"c":5}
```

## (upper) (lower)

These functions make string values uppercase or lowercase.

```clj
(upper)
(lower)
```

**Example**

```go
a := jsqt.Get(`"hello"`, `(upper)`)
b := jsqt.Get(`"WORLD"`, `(lower)`)

fmt.Println(a) // "HELLO"
fmt.Println(b) // "world"
```

## (stringify) (jsonify)

These functions stringify or jsonify JSON values.

```clj
(stringify)
(jsonify)
```

**Example**

```go
a := jsqt.Get(`{"a":3}`, `(stringify)`)
b := jsqt.Get(`"{\"a\":3}"`, `(jsonify)`)

fmt.Println(a) // "{\"a\":3}"
fmt.Println(b) // {"a":3}
```

## (replace)

This function replaces all occurrences of `old` with `new` in a string.

```clj
(replace old new)
```

The arguments must be a string.

**Example**

```go
a := jsqt.Get(`"hello world"`, `(replace " " "_")`)

fmt.Println(a) // "hello_world"
```

## (join)

This function joins an array of strings given a separator.

```clj
(join sep)
(join sep arg)
```

`sep` argument must be a string; `arg` is optional and can be a key or function.

**Example**

```go
a := jsqt.Get(`["a","b","c"]`, `(join "_")`)
b := jsqt.Get(`{"x":["a","b","c"]}`, `(join "_" x)`)

fmt.Println(a) // "a_b_c"
fmt.Println(b) // "a_b_c"
```

## (split)

This function splits a string given a separator.

```clj
(split sep)
(split sep arg)
```

`sep` argument must be a string; `arg` is optional and can be a key or function.

**Example**

```go
a := jsqt.Get(`"one,two"`, `(split ",")`)
b := jsqt.Get(`{ "a": "one,two" }`, `(split "," a)`)

fmt.Println(a) // ["one","two"]
fmt.Println(b) // ["one","two"]
```

## (concat)

This function concats values into a string.

```clj
(concat a b ...)
```

The argument list must be keys or functions that return strings, numbers, booleans or nulls.

**Example**

```go
a := jsqt.Get(`{ "one": "Hello", "two": "World" }`, `(concat one (raw " ") two)`)

fmt.Println(a) // "Hello World"
```

## (sort)

This function sorts a JSON array or object keys.

```clj
(sort)
(sort key)
(sort desc)
(sort desc key)
```

`key` is for sorting an array of objects by a key.

Use `desc` to sort descending.

**Example**

```go
a := jsqt.Get(`[5,4,3]`, `(sort)`)
b := jsqt.Get(`["c","b","a"]`, `(sort desc)`)
c := jsqt.Get(`{ "b": 3, "a": 4 }`, `(sort)`)
d := jsqt.Get(`[{ "a": 4 }, { "a": 3 }]`, `(sort a)`)
e := jsqt.Get(`[{ "a": 3 }, { "a": 4 }]`, `(sort desc a)`)

fmt.Println(a) // [3,4,5]
fmt.Println(b) // ["c","b","a"]
fmt.Println(c) // {"a":4,"b":3}
fmt.Println(d) // [{ "a": 3 },{ "a": 4 }]
fmt.Println(e) // [{ "a": 4 },{ "a": 3 }]
```

## (keys) (values) (entries) (objectify)

The `keys` function collects all keys of an object into an array.

The `values` function collects all key values of an object into an array.

The `entries` function collects all keys and values of an object into an array.

The `objectify` function reverts `entries`.

```clj
(keys)
(values)
(entries)
(objectify)
```

**Example**

```go
j := `{ "a": 3, "b": 4 }`

a := jsqt.Get(j, `(keys)`)
b := jsqt.Get(j, `(values)`)
c := jsqt.Get(j, `(entries)`)
d := jsqt.Get(j, `(entries) (objectify)`)

fmt.Println(a) // ["a","b"]
fmt.Println(b) // [3,4]
fmt.Println(c) // [["a",3],["b",4]]
fmt.Println(d) // {"a":3,"b":4}
```

## (ugly) (pretty)

These functions uglify or prettify a JSON.

```clj
(ugly)
(pretty)
```

**Example**

```go
j := `
  {
    "id": 1,
    "name": "Bret",
    "address": { "city": "Gwen" }
  }`

a := jsqt.Get(j, `(ugly)`)
b := jsqt.Get(j, `(pretty)`)

fmt.Println(a) // {"id":1,"name":"Bret","address":{"city":"Gwen"}}
fmt.Println(b)
/*
{
    "id": 1,
    "name": "Bret",
    "address": {
        "city": "Gwen"
    }
}
*/
```

## (iterate)

This function iterates over keys and values of a valid JSON
and applies a map function to transform them.

```clj
(iterate key val)
(iterate -r key val)
(iterate -c val)

(iterate -f key val)
(iterate -k key)
(iterate -v val)
(iterate -kv keyval)
```

Arguments must be functions.

The `iterate` iterates over all keys and values; values include objects and arrays.
Use `-r` flag to tell iterate to emit the root value. When the root value is emitted `(key)` is `null`.
If either `key` or `val` functions return an empty context the field is removed from the result.
Note that this version of iterate is a depth-first post-order traversal.

The `iterate -f` is just a faster version of iterate, but it does not emit objects and arrays.

The `iterate -k` is just a faster version that iterates over all keys. The `key` argument receives the key string.

The `iterate -v` is just a faster version that iterates over values, but values do not include objects and arrays.
The `val` argument receives the value.

The `iterate -kv` is just a faster version that iterates over all keys and values consecutively,
but values do not include objects and arrays.
The `keyval` argument receives a key or a value consecutively.

The `iterate -c` can be used to recursivelly collect values into an array.
Note that this version of iterate is a depth-first pre-order traversal.

**Example**

```go
j := `{ "Month": "May", "Next": { "Month": "June", "Other": "" }, "Other": "" }`

// Preserve non-empty fields (aka: remove all empty fields).
a := jsqt.Get(j, `(iterate (key) (not (is-empty)))`)

// Convert all keys to lowercase and replace the values with the size of the string.
b := jsqt.Get(j, `(iterate -f (lower) (size))`)

// Convert all keys to uppercase.
c := jsqt.Get(j, `(iterate -k (upper))`)

// Convert all values to uppercase.
d := jsqt.Get(j, `(iterate -v (upper))`)

// Convert all keys and values to uppercase.
e := jsqt.Get(j, `(iterate -kv (upper))`)

// Collect all values in which the keys equal "Month".
f := jsqt.Get(j, `(iterate -c (== "Month" (key)))`)

fmt.Println(a) // {"Month":"May","Next":{"Month":"June"}}
fmt.Println(b) // {"month":3,"next":{"month":4,"other":0},"other":0}
fmt.Println(c) // {"MONTH":"May","NEXT":{"MONTH":"June","OTHER":""},"OTHER":""}
fmt.Println(d) // {"Month":"MAY","Next":{"Month":"JUNE","Other":""},"Other":""}
fmt.Println(e) // {"MONTH":"MAY","NEXT":{"MONTH":"JUNE","OTHER":""},"OTHER":""}
fmt.Println(f) // ["May","June"]
```

It is also possible to use [(key)](#key-val) and [(val)](#key-val) functions with iterate.

## (debug)

This function prints JSON values to the stdout for debugging.

```clj
(debug)
(debug label)
```

The `label` argument is optional and can be used to label a debug step.

**Example**

```go
j := `[{ "a": { "b": [3] } }, { "a": { "b": [4] } }]`

v := Get(j, `(collect a (debug) b (debug b_val) (flatten) (debug flatn))`)

fmt.Println("Result:", v)

// Output:
// [debug] { "b": [3] }
// [b_val] [3]
// [flatn] 3
// [debug] { "b": [4] }
// [b_val] [4]
// [flatn] 4
// Result: [3,4]
```

## (def)

This function allows defining custom functions.
This is like giving a function an alias, so that the alias can be used instead of a long function.
Also useful to avoid code repetition.

```clj
(def name fun)
```

The `name` argument is the custom function name.
The `fun` argument is the defined function.

**Example**

```go
j := `[3, 4]`

a := jsqt.Get(j, `(def num2str (collect (stringify))) (obj a (num2str) b (num2str))`)

fmt.Println(a) // {"a":["3","4"],"b":["3","4"]}
```

## (save) (load)

These functions save and load a context.

```clj
(save)
(save val)
(save -k key ...)
(save -k key -v val ...)

(load)
(load key)
```

`(save)` saves the value it receives;
`(save val)` saves the value from `val` (can be a key or a function); these two forms of save make no allocation.
Note that each call to save overrides the previous value.
Use `(load)` to load a value saved by these two methods.

`(save -k key ...)` saves the value of a key under an id of the same name;
Use `-v val` (can be a key or a function) to save a value under an id.
Use `(load key)` to load a value saved with `-k`.

Save function returns the context it receives.

**Example**

```go
j := `{ "a": 3, "b": 4 }`

a := jsqt.Get(j, `a (save) (arr (load))`)
b := jsqt.Get(j, `(save a) (arr (load))`)
c := jsqt.Get(j, `(save (get a)) (arr (load))`)
d := jsqt.Get(j, `(save -k a b x -v (raw 7)) (arr (load a) (load b) (load x))`)

fmt.Println(a) // [3]
fmt.Println(b) // [3]
fmt.Println(c) // [3]
fmt.Println(d) // [3,4,7]
```

## (key) (val)

These functions are available only inside a few functions.

```clj
(key)
(val)
```

In functions that iterate over arrays `(key)` is the array index and `(val)` is the array item.

In functions that deal with objects `(key)` is the object key and `(val)` is the key value.

**Example**

```go
a := jsqt.Get(`[ 33, 44 ]`, `(collect (key))`)
b := jsqt.Get(`[ 33, 44 ]`, `(collect (val))`)
c := jsqt.Get(`{ "a": 3 }`, `(iterate (concat (key) (val)) (concat (val) (key)))`)

fmt.Println(a) // [0,1]
fmt.Println(b) // [33,44]
fmt.Println(c) // {"a3":"3a"}
```

Note that `(key)` and `(val)` are overridden inside nested functions and this might lead to confusion.
The solution in this case is to [(save)](#save-load) the previous value and [(load)](#save-load) it.

## (arg)

This function returns the value of an argument provided with `GetWith(jsn, qry, args)`.

```clj
(arg index)
```

`index` is a function or a raw value.

**Example**

```go
a := jsqt.GetWith(``, `(obj msg (arg 0) val (arg 1))`, []any{"hello", 3})

fmt.Println(a) // {"msg":"hello","val":3}
```

Use a function in the format `func (Json) Json` as argument
to apply a custom logic to the current context.

**Example**

```go
j := `{ "date": "2022-09-07T12:30:00Z" }`

f := func(date jsqt.Json) jsqt.Json {
    curDate, _ := time.Parse(time.RFC3339, date.TrimQuote())
    newDate := curDate.Format(time.RFC1123)
    return jsqt.JSON(newDate).Stringify()
}

a := jsqt.GetWith(j, `(set date (arg 0))`, []any{f})

fmt.Println(a) // {"date":"Wed, 07 Sep 2022 12:30:00 UTC"}
```

## (expr)

This function calculates math expressions.

```clj
(expr a op b ...)
```

`a` and `b` are the operands and can be a function or a raw value.

`op` is the operator and can be any of these: `+ - * / %`.

**Example**

```go
j := `{ "a": 3 }`

a := jsqt.Get(j, `(expr 4 * 5 + (get a))`)
b := jsqt.Get(j, `(expr 4 * (expr 5 + (get a)))`)

fmt.Println(a) // 23
fmt.Println(b) // 32
```

## (group)

This function groups values.

```clj
(group key val)
(group key val -a)
(group key val -a newkey newval)
```

`key` can be a function or a key and it is the value that becomes the group key.

`val` can be a function or a key and it is the value that is added to a group.

`(group)` groups into this format: `{ "group1": [], "group2": [] }`.
Use `-a` to group into an array in the format: `[{ "key": "group1", "values": [] }, { "key": "group2", "values": [] }]`.
To rename `"key"` and `"values"` use `-a newkey newval` (both must be provided).

[(key)](#key-val) and [(val)](#key-val) can be used to access the array index and value.

**Example**

```go
j := `[{ "g": "dog", "v": 15 }, { "g": "dog", "v": 12 }, { "g": "cat", "v": 10 }]`

a := jsqt.Get(j, `(group g (pluck g))`)
b := jsqt.Get(j, `(group g (pluck g) -a)`)
c := jsqt.Get(j, `(group g (pluck g) -a group vals)`)

fmt.Println(a) // {"dog":[{"v":15},{"v":12}],"cat":[{"v":10}]}
fmt.Println(b) // [{"key":"dog","values":[{"v":15},{"v":12}]},{"key":"cat","values":[{"v":10}]}]
fmt.Println(c) // [{"group":"dog","vals":[{"v":15},{"v":12}]},{"group":"cat","vals":[{"v":10}]}]
```

**Example**

```go
j := `[3,4,3,4,5]`

a := jsqt.Get(j, `(group (val) (key))`)

fmt.Println(a) // {"3":[0,2],"4":[1,3],"5":[4]}
```

## (unwind)

This function deconstructs an array field.

```clj
(unwind key)
(unwind key -r newkey)
```

`key` can be a key or a function. Use `-r` to rename a key.

**Example**

```go
j := `{ "a": 3, "b": [ 4, 5 ] }`

a := jsqt.Get(j, `(unwind b)`)
b := jsqt.Get(j, `(unwind b -r x)`)

fmt.Println(a) // [{"a":3,"b":4},{"a":3,"b":5}]
fmt.Println(b) // [{"a":3,"x":4},{"a":3,"x":5}]
```

Unwind also work with objects inside arrays:

**Example**

```go
j := `[{ "a": 3, "b": [ 4, 5 ] }, { "a": 6, "b": [ 7, 8 ] }]`

a := jsqt.Get(j, `(unwind b)`)

fmt.Println(a) // [{"a":3,"b":4},{"a":3,"b":5},{"a":6,"b":7},{"a":6,"b":8}]
```

## (transpose)

This function is easier to understand with an example.

```clj
(transpose)
```

It will convert this input:

```json
{
    "a": [ 3, 5 ],
    "b": [ 4, 6 ]
}
```

Into this output:

```json
[
    {
        "a": 3,
        "b": 4
    },
    {
        "a": 5,
        "b": 6
    }
]
```

And vice-versa. It will convert the output above back to the input again if you transpose it.

Note that in general transpose is reversible only when fields have the same number of items.

## (match)

This function matches a value against a prefix, suffix or regular expression.

```clj
(match pattern)
(match -r pattern)
(match -p pattern)
(match -s pattern)
(match -k pattern)
(match -v key pattern)
(match -kk pattern)
```

The `pattern` argument can be a function or a raw value.

Use `-r` to match a regular expression; `-p` to match a prefix; `-s` to match a suffix;
and no flag for an exact match.

Use `-k` to match an object key by a pattern. It returns the matched key value or an empty context.

Use `-v key` to match the value of a key by a pattern. It returns the current context or an empty context.

Use `-kk` to match a key and return the matched key.

**Example**

```go
j := `{ "first_name": "Jim", "last_name": "May", "first_class": "July" }`

a := jsqt.Get(j, `(match -k -r name$)`)
b := jsqt.Get(j, `(match -k -s class)`)
c := jsqt.Get(j, `(match -k -p last)`)
d := jsqt.Get(j, `(match -k first_name)`)
e := jsqt.Get(j, `(iterate (match -r name) (val))`)
f := jsqt.Get(j, `(keys) * (match -s name)`)
g := jsqt.Get(j, `(values) * (obj name (this)) (match -v name -p J)`)
h := jsqt.Get(j, `(match -kk -s name)`)

fmt.Println(a) // "Jim"
fmt.Println(b) // "July"
fmt.Println(c) // "May"
fmt.Println(d) // "Jim"
fmt.Println(e) // {"first_name":"Jim","last_name":"May"}
fmt.Println(f) // ["first_name","last_name"]
fmt.Println(g) // [{"name":"Jim"},{"name":"July"}]
fmt.Println(h) // "first_name"
```

# Truth Table

|       | void | empty | blank | nully | some | falsy | truthy |
| ----- | ---- | ----- | ----- | ----- | ---- | ----- | ------ |
| []    |  T   |   T   |   T   |   T   |  T   |   T   |        |
| {}    |  T   |   T   |   T   |   T   |  T   |   T   |        |
| ""    |      |   T   |       |   T   |  T   |   T   |        |
| null  |      |       |   T   |   T   |      |   T   |        |
| 0     |      |       |       |       |  T   |   T   |        |
| false |      |       |       |       |  T   |   T   |        |
| true¹ |      |       |       |       |  T   |       |   T    |

¹ This value is the same for all other values (`3`, `"a"`, `[3]`, `{ "a": 3 }`, etc).

---

<p align="center">
    <i><a href="https://www.youtube.com/watch?v=AWb0Z5yA1XA">non nobis Domine sed nomine tuo da gloriam</a></i>
</p>
