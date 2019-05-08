package websocketwasm

import (
	"fmt"
	"syscall/js"

	"github.com/adrianbrad/websocketwasm/browser"
)

type WebSocket struct {
	*browser.WebSocket
	Send     chan []byte
	Received chan js.Value

	onMessageFunc js.Func
}

func handshakeOnOpen(ch chan error, removeHandlers func()) js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		removeHandlers()
		close(ch)
		return nil
	})
}

func handshakeOnClose(ch chan error, removeHandlers func()) js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		removeHandlers()
		ch <- fmt.Errorf("%s", args[0].Get("code").String())
		close(ch)
		return nil
	})
}

func Dial(url string) (ws *WebSocket, err error) {
	wsb, err := browser.NewWebsocket(url)
	if err != nil {
		return
	}
	ws = &WebSocket{
		WebSocket: wsb,
		Send:      make(chan []byte),
		Received:  make(chan js.Value),
	}

	var (
		openHandler  js.Func
		closeHandler js.Func
	)

	removeHandlers := func() {
		ws.Call("removeEventListener", "open", openHandler)
		ws.Call("removeEventListener", "close", closeHandler)
		openHandler.Release()
		closeHandler.Release()
	}

	openCh := make(chan error, 1)

	// onOpenHandler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
	// 	removeHandlers()
	// 	close(openCh)
	// 	return nil
	// })
	// onCloseHandler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
	// 	removeHandlers()
	// 	openCh <- fmt.Errorf("%s", args[0].Get("code").String())
	// 	close(openCh)
	// 	return nil
	// })
	openHandler = handshakeOnOpen(openCh, removeHandlers)
	closeHandler = handshakeOnClose(openCh, removeHandlers)

	wsb.OnOpen(openHandler)
	wsb.OnClose(closeHandler)

	err, ok := <-openCh
	if ok && err != nil {
		return nil, err
	}

	ws.Set("binaryType", "arraybuffer")
	ws.onMessageFunc = js.FuncOf(ws.onMessageListener)
	wsb.OnMessage(ws.onMessageFunc)
	wsb.OnClose(js.FuncOf(ws.onCloseListener))
	return
}

func (w *WebSocket) onMessageListener(this js.Value, args []js.Value) interface{} {
	fmt.Println("Message")
	go func() {
		w.Received <- args[0]
	}()
	return nil
}

func (w *WebSocket) SetOnMessageFunc() interface{} {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		w.WebSocket.Call("removeEventListener", "message", w.onMessageFunc)
		w.onMessageFunc.Release()
		w.WebSocket.Call("addEventListener", "message", args[0])
		// TODO close receive channel
		return nil
	})
}

func (w *WebSocket) onCloseListener(this js.Value, args []js.Value) interface{} {
	fmt.Println("Close")
	return nil
}

func (w *WebSocket) Read(b []byte) (n int, err error) {
	m := <-w.Received
	receivedBytes := getFrameData(m.Get("data"))
	n = copy(b, receivedBytes)
	return
}

func (w *WebSocket) Write(b []byte) (n int, err error) {
	byteArray := js.TypedArrayOf(b)
	defer byteArray.Release()
	err = w.WebSocket.Send(js.ValueOf(byteArray))
	if err != nil {
		n = 0
		return
	}
	n = len(b)
	return
}

func (w *WebSocket) WriteString(s string) (n int, err error) {
	err = w.WebSocket.Send(js.ValueOf(s))
	if err != nil {
		n = 0
		return
	}
	n = len(s)
	return
}

func getFrameData(obj js.Value) []byte {
	if obj.InstanceOf(js.Global().Get("ArrayBuffer")) {
		uint8Array := js.Global().Get("Uint8Array").New(obj)
		data := make([]byte, uint8Array.Length())
		for i, arrayLen := 0, uint8Array.Length(); i < arrayLen; i++ {
			data[i] = byte(uint8Array.Index(i).Int())
		}
		return data
	}
	return []byte(obj.String())
}
