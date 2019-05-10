// +build js,wasm

package client

import (
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/adrianbrad/websocketwasm"
)

func BindConnect() {
	var connect js.Func
	connect = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) != 3 {
			return "Invalid nr of args"
		}

		if args[0].Type() != js.TypeString {
			return "key is not string"
		}

		if args[1].Type() != js.TypeString {
			return "room id is not string"
		}

		if args[2].Type() != js.TypeFunction {
			return "callback func is not a func"
		}

		key := args[0]
		roomID := args[1]
		//callback is a function that receives one parameter which is either an error or the chat object
		callbackFunc := args[2]
		go func() {
			goObj := map[string]interface{}{}
			jsObj := js.ValueOf(goObj)
			chatWS, err := websocketwasm.Dial("ws://localhost:8080/chat?" + "key=" + key.String() + "&roomID=" + roomID.String())
			if err != nil {
				callbackFunc.Invoke(err.Error())
				return
			}
			jsObj.Set("sendTextMessage", composeTextMessage(chatWS))

			userInfo := make([]byte, 1024)

			n, err := chatWS.Read(userInfo)
			if err != nil {
				callbackFunc.Invoke(err.Error())
				return
			}
			userInfo = userInfo[:n]
			var uI map[string]interface{}
			json.Unmarshal(userInfo, &uI)
			jsObj.Set("user", js.ValueOf(uI))
			jsObj.Set("setOnMessageFunc", chatWS.SetOnMessageFunc())
			callbackFunc.Invoke(jsObj)
		}()

		connect.Release()
		js.Global().Set("connect", js.Undefined())
		return nil
	})
	js.Global().Set("connect", connect)
}

func composeTextMessage(chatWS *websocketwasm.WebSocket) js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		chatWS.WriteString(fmt.Sprintf(`
	{
		"action": "text",
		"body":{
			"content": "%s"
		}
	}
	`, args[0].String()))
		return nil
	})
}
