// +build js,wasm
package main

import "github.com/adrianbrad/chat-v2/client/client"

func main() {
	client.BindConnect()
	select {}
}
