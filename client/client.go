// +build js,wasm
package main

import (
	"syscall/js"

	"github.com/adrianbrad/websocketwasm"
)

func main() {
	var startChat js.Func
	startChat = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) != 3 {
			return "Invalid nr of args"
		}

		userID := args[0]
		roomID := args[1]
		callbackFunc := args[2]

		chatWS, err := websocketwasm.Dial(getWSBaseURL() + "echo")

		jsObj.Set("sendMessage", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			chatWS.WriteString(args[0].String())
		}))

		// output is an object with an '.addMessage' method
		jsObj.Set("setOutput", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			return 2
		}))

		startChat.Release()
		return jsObj
	})
	js.Global().Set("startChat", startChat)
}
