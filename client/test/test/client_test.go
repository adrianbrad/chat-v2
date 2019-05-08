package client_test

import (
	"fmt"
	"testing"
	"syscall/js"
	"net/http"
	"github.com/adrianbrad/chat-v2/client/client"
	"github.com/LinearZoetrope/testevents"
	"io/ioutil"

)

func TestWSConnection(t_ *testing.T) {
	t := testevents.Start(t_, "TestWSHandshakeSucess", true)
	defer t.Done()
	client.BindConnect()
	done := make(chan struct{}, 1)

	setOnMessageFunc := js.FuncOf(func(this js.Value, args []js.Value) interface {} {
		fmt.Println("Custom on mes func")
		fmt.Println(args[0].Get("data"))
		done <-struct{}{}
		return nil
	})


	callbackFunc := js.FuncOf(func(this js.Value, args[]js.Value) interface{}{
		args[0].Get("setOnMessageFunc").Invoke(setOnMessageFunc)
		args[0].Get("sendTextMessage").Invoke("hello")
		return nil
	})

	go func() {
		resp, _ := http.Get("http://localhost:3000/getToken?user=user_a")
		b,_ :=ioutil.ReadAll(resp.Body)
		js.Global().Get("connect").Invoke(string(b), "room_a", callbackFunc)
	}()


	<- done
}