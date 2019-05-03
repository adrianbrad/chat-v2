// +build js,wasm
package main

import "syscall/js"

func main() {
	js.Global().Set("connect", js.FuncOf(nil))
	js.Global().Set("sendMessage", js.FuncOf(nil))
}
