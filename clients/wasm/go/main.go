// +build js,wasm

package main

import "github.com/adrianbrad/chat-v2/clients/wasm/go/client"

func main() {
	client.BindConnect()
	select {}
}
