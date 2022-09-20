//go:build js && wasm

package main

import (
	"syscall/js"

	"github.com/ofabricio/jsqt"
)

func main() {
	js.Global().Set("jsqtGet", jsqtGet())
	<-make(chan bool)
}

func jsqtGet() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		jsn := args[0].String()
		qry := args[1].String()
		return jsqt.Get(jsn, qry).String()
	})
}
