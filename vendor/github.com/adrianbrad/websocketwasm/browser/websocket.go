package browser

import (
	"syscall/js"
)

type WebSocket struct {
	js.Value
}

func NewWebsocket(url string) (ws *WebSocket, err error) {
	defer handleErrorIfRaised(err)
	webSocket := js.Global().Get("WebSocket").New(url)
	ws = &WebSocket{
		Value: webSocket,
	}

	return
}

func (w *WebSocket) OnOpen(callback js.Func) (err error) {
	defer handleErrorIfRaised(err)
	w.Call("addEventListener", "open", callback)
	return
}

func (w *WebSocket) OnClose(callback js.Func) (err error) {
	defer handleErrorIfRaised(err)
	w.Call("addEventListener", "close", callback)
	return
}

func (w *WebSocket) Send(data js.Value) (err error) {
	defer handleErrorIfRaised(err)
	w.Value.Call("send", data)
	return
}

func (w *WebSocket) OnMessage(callback js.Func) (err error) {
	defer handleErrorIfRaised(err)
	w.Call("addEventListener", "message", callback)
	return
}

func (w *WebSocket) Close() (err error) {
	defer handleErrorIfRaised(err)
	w.Value.Call("close", 1000)
	return
}

func handleErrorIfRaised(err error) {
	e := recover()
	if e == nil {
		return
	}
	if jsErr, ok := e.(*js.Error); ok && jsErr != nil {
		err = jsErr
	} else {
		panic(e)
	}
	return
}
