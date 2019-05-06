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

		key := args[0]
		roomID := args[1]
		//callback is a function that receives one parameter which is either an error or the chat object
		callbackFunc := args[2]

		chatWS, err := websocketwasm.Dial(getWSBaseURL() + "key=" + key + "&roomID=" + roomID)

		jsObj.Set("sendMessage", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			chatWS.WriteString(args[0].String())
		}))

		//handl the callback func
		startChat.Release()
		js.Global().Set("startChat", js.Undefined())
		return jsObj
	})
	js.Global().Set("startChat", startChat)
}
