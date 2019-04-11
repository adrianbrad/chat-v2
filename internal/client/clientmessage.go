package client

type ClientMessage struct {
	Content map[string]interface{}
	Client  Client
}
