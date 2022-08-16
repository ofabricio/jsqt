# JSQT

This Go package provides a closure-like syntax to query and transform JSON documents.

[![build](https://github.com/ofabricio/jsqt/actions/workflows/go.yml/badge.svg)](https://github.com/ofabricio/jsqt/actions/workflows/go.yml)

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

- ⚠ In the current state of development many functions are not consolidated yet.
  Watch for updates if you are using them, because they can change anytime as there is no official release yet.
  To name a few ones that especially bother me: `iterateX`, `is-X` and filters `==, !=, etc`.
- Don't open PR.

# Install

```
go get github.com/ofabricio/jsqt
```

# Query functions

Query functions have a name and arguments and live inside `()`.
The function name is always after a `(`.
For example, in `(get a b)` function, `get` is the query function name and `a` and `b` are its arguments.
Make sure to write a valid query since no validation is done during the parsing.

There are three types of function arguments:

- **Function** - These are functions, for example: `(get name)`, `(root)`.
  When the parser finds a function it calls it and uses its result as argument.
  A function always receives the current JSON context as input;
  its arguments also receive the current context.
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

The list of arguments can be keys or functions.
The arguments work like a pipeline: the output of the first one is the input of the next one and so on.

This is one of the most important functions as its pipeline behavior is what allows passing
a context to other functions. Because of that the root function is also a get function.

**Example**

```go
j := `{ "a": { "b c": 3 }, "d": [4, 5] }`

a := jsqt.Get(j, `(get a)`)
b := jsqt.Get(j, `(get a "b c")`)
c := jsqt.Get(j, `(get d 1)`)
d := jsqt.Get(j, `(get d (size))`)
e := jsqt.Get(j, `d (size)`) // The root function is a (get) so you can omit it.

fmt.Println(a) // { "b c": 3 }
fmt.Println(b) // 3
fmt.Println(c) // 5
fmt.Println(d) // 2
fmt.Println(e) // 2
```

## (collect)

This function collects values into a JSON array.

```clj
(collect arg ...)
```

The list of arguments can be keys or functions.

The `(collect)` function arguments works like `(get)` function, but if some argument is an array,
each individual array item is emitted to the next argument and collected into an array.
This allows to filter and change each item.

Individual items must be handled or an empty array is returned.
For example `(collect myarr)` returns nothing, but `(collect myarr (this))` returns an identity array.

**Example**

```go
j := `{ "a": [{ "b": 3, "c": [5] }, { "b": 4, "c": [6] }] }`

a := jsqt.Get(j, `(collect a b)`)
b := jsqt.Get(j, `(collect a c (flatten))`)

fmt.Println(a) // [3, 4]
fmt.Println(b) // [5, 6]
```

## (obj)

This function creates a JSON object.

```clj
(obj field value ...)
```

The arguments are pairs of JSON keys and values. The `field` can be a function or a raw value.
The `value` can be a function or a key.

**Example**

```go
j := `{ "loc": [ 63.4682, -20.1754 ] }`

a := jsqt.Get(j, `(obj lat (get loc 0) lng (get loc 1))`)
b := jsqt.Get(j, `(get loc (obj lat 0 lng 1))`) // Same as above.

fmt.Println(a) // {"lat":63.465,"lng":-20.178}
fmt.Println(b) // {"lat":63.465,"lng":-20.178}
```

## (arr)

This function creates a JSON array.

```clj
(arr item ...)
```

Each argument becomes an array item. The arguments can be functions or keys.

**Example**

```go
j := `{ "author": "Mary", "comment": "Hello" }`

a := jsqt.Get(j, `(arr (raw "author") author (raw "comment") comment)`)

fmt.Println(a) // ["author","Mary","comment","Hello"]
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

## (upsert)

This function creates a new object field or updates an existing one.

```clj
(upsert field value ...)
```

The arguments are pairs of JSON keys and values.
Both `field` and `value` can be a function or a raw value.

**Example**

```go
j := `{ "message": "Hello" }`

a := jsqt.Get(j, `(upsert message "World")`)
b := jsqt.Get(j, `(upsert id 123)`)

fmt.Println(a) // {"message":"World"}
fmt.Println(b) // {"id":123,"message":"Hello"}
```

## (flatten)

This function flattens a JSON array.

```clj
(flatten)
```

**Example**

```go
a := jsqt.Get(`[ [3], [4] ]`, `(collect (flatten))`)
b := jsqt.Get(`{ "a": [ [{ "b": 3 }], [{ "b": 4 }] ]`, `(collect a (flatten) b)`)

fmt.Println(a) // [3,4]
fmt.Println(b) // [3,4]
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

## (filters)

These are filter functions.
They return the current context if true or an empty context if false.

```clj
(== field value)
(!= field value)
(>= field value)
(<= field value)
(> field value)
(< field value)
```

The `field` argument is the field to filter on, and can be a key or a function.
The `value` argument is the expected value, and can be a raw value or a function.

**Example**

```go
j := `[{ "a": 3, "b": true }, { "a": 4, "b": false }, { "a": 5, "b": true }]`

a := jsqt.Get(j, `(collect (== b true) a)`)

fmt.Println(a) // [3,5]
```

## (is-x)

These functions test a JSON for the corresponding value and return it if true or an empty context if false.

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

**Example**

```go
j := `{}`

a := jsqt.Get(j, `(is-obj)`)
b := jsqt.Get(j, `(is-num)`)
 
fmt.Println(a) // {}
fmt.Println(b) //
```

## (bool)

This function returns `true` if it gets a value or `false` if it gets an empty context.

```clj
(bool)
```

**Example**

```go
a := jsqt.Get(`[]`, `(get (is-arr) (bool))`)
b := jsqt.Get(`{}`, `(get (is-arr) (bool))`)

fmt.Println(a) // true
fmt.Println(b) // false
```

## (or) (and) (not)

These functions apply OR, AND, NOT logic to its arguments.

```clj
(or a b)
(and a b)
(not a)
```

The arguments must be functions.

**Example**

```go
j := `[{"a":3},{"a":4},{"a":5},{"a":6}]`

a := jsqt.Get(j, `(collect (or (< a 4) (> a 5)) a)`)
b := jsqt.Get(j, `(collect (and (>= a 4) (<= a 5)) a)`)
c := jsqt.Get(j, `(collect (not (<= a 4)) a)`)

fmt.Println(a) // [3,6]
fmt.Println(b) // [4,5]
fmt.Println(c) // [5,6]
```

## (if)

This function works like a regular `if`.

```clj
(if cond then else)
(if cond then)
```

if `cond` is true `then` is executed, otherwise `else` is. The arguments can be functions or keys.
The `else` argument is optional and returns `(this)` when omited.

**Example**

```go
j := `{}`

a := jsqt.Get(j, `(if (is-obj) (raw "It's an object!") (raw "Not an object"))`)

fmt.Println(a) // "It's an object!"
```

## (either)

This function returns the first argument value that is not [nully](#truth-table).

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
(default value)
```

The `value` argument can be a function or a raw value.

**Example**

```go
j := `[{ "a": 3 }, { "b": 4 }, { "a": 5 }]`

a := jsqt.Get(j, `(collect a (default 0))`)

fmt.Println(a) // [3,0,5]
```

## (pick) (pluck)

These functions pick or pluck fields from a JSON object.

```clj
(pick key ...)
(pluck key ...)
```

The arguments are a list of object keys.

**Example**

```go
j := `{ "a": 3, "b": 4, "c": 5 }`

a := jsqt.Get(j, `(pick a c)`)
b := jsqt.Get(j, `(pluck a c)`)

fmt.Println(a) // {"a":3,"c":5} 
fmt.Println(b) // {"b":4} 
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
```

The `sep` argument must be a string.

**Example**

```go
a := jsqt.Get(`["a","b","c"]`, `(join "_")`)

fmt.Println(a) // "a_b_c"
```

## (concat)

This function concats strings or arrays.

```clj
(concat a b ...)
```

The argument list must be keys or functions that return strings or arrays.

**Example**

```go
a := jsqt.Get(`{ "one": "Hello", "two": "World" }`, `(concat one (raw " ") two)`)
b := jsqt.Get(`{ "a": [3,4], "b": [5,6] }`, `(concat a b)`)

fmt.Println(a) // "Hello World"
fmt.Println(b) // [3,4,5,6]
```

## (sort)

This function sorts a JSON array or object keys.

```clj
(sort dir)
(sort dir key)
```

The `dir` argument is the sorting direction.
The `key` argument is for sorting an array of objects by a key.

**Example**

```go
a := jsqt.Get(`[{ "a": 3 }, { "a": 4 }]`, `(sort desc a)`)
b := jsqt.Get(`[5,4,3]`, `(sort asc)`)
c := jsqt.Get(`["c","b","a"]`, `(sort asc)`)
d := jsqt.Get(`{ "b": 3, "a": 4 }`, `(sort asc)`)

fmt.Println(a) // [{ "a": 4 },{ "a": 3 }]
fmt.Println(b) // [3,4,5]
fmt.Println(c) // ["a","b","c"]
fmt.Println(d) // {"a":4,"b":3}
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

These functions iterate over keys and values of a valid JSON
and apply a map function to transform them.

```clj
(iterate key val)
(iterate-pair key val)
(iterate-k key)
(iterate-v val)
(iterate-kv keyval)

(iterate-all key val)
(iterate-all-pair key val)
```

The `iterate` iterates over all keys and values, but values do not include objects and arrays.
Both `key` and `val` arguments must be a function.

The `iterate-pair` iterates over all keys and values, but values do not include objects and arrays.
Both `key` and `val` arguments can be a function or an index key as they receive an array in the format `[key, val]`.

The `iterate-k` iterates over all keys. The `key` argument must be a function and receives the key string.

The `iterate-v` iterates over all values, but values do not include objects and arrays.
The `val` argument must be a function and receives the value.

The `iterate-kv` iterates over all keys and values consecutively, but values do not include objects and arrays.
The `keyval` argument must be a function and receives a key or a value consecutively.

**Example**

```go
j := `{ "One": "Two", "Three": "Four" }`

a := jsqt.Get(j, `(iterate (lower) (upper))`)
b := jsqt.Get(j, `(iterate-pair (concat 0 (raw "-") 1) (concat 1 (raw "-") 0))`)
c := jsqt.Get(j, `(iterate-k (upper))`)
d := jsqt.Get(j, `(iterate-v (upper))`)
e := jsqt.Get(j, `(iterate-kv (upper))`)

fmt.Println(a) // {"one":"TWO","three":"FOUR"}
fmt.Println(b) // {"One-Two":"Two-One","Three-Four":"Four-Three"}
fmt.Println(c) // {"ONE":"Two","THREE":"Four"}
fmt.Println(d) // {"One":"TWO","Three":"FOUR"}
fmt.Println(e) // {"ONE":"TWO","THREE":"FOUR"}
```

The `iterate-all` iterates over all keys and values, but values also include objects and arrays.
Both `key` and `val` arguments must be functions.
If either key or val functions return an empty context the field is removed from the result.

The `iterate-all-pair` iterates over all keys and values, but values also include objects and arrays.
Both `key` and `val` arguments can be a function or an index key as they receive an array in the format `[key, val]`.
If either key or val functions return an empty context the field is removed from the result.

**Example**

```go
j := `{ "a": 3, "b": {}, "c": { "d": {} }, "e": { "f": 4 } }`

a := jsqt.Get(j, `(iterate-all (upper) (not (is-empty)))`)
b := jsqt.Get(j, `(iterate-all-pair (get 0 (upper)) (get 1 (not (is-empty))))`)

fmt.Println(a) // {"A":3,"E":{"F":4}}
fmt.Println(b) // {"A":3,"E":{"F":4}}
```

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
