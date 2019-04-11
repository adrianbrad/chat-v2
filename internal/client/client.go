package client

import (
	"chat-v2/internal/user"
)

type ClientFactoryMethod func(wsConn wsConn, user *user.User) Client

type wsConn interface {
	ReadJSON(v interface{}) error
	WriteJSON(v interface{}) error
	Close() error
}

type Client interface {
	Read(messageQueue chan *ClientMessage)
	Write()
}

type client struct {
	wsConn
	user         *user.User
	MessageQueue chan *ClientMessage
}

func NewClient(wsConn wsConn, user *user.User) Client {
	return &client{
		wsConn: wsConn,
		user:   user,
	}
}

// Proccess messages sent by the websocket connection and forward them to the channel given as parameter
func (client *client) Read(messageQueue chan *ClientMessage) {
	for {
		var receivedMessage map[string]interface{}
		err := client.ReadJSON(&receivedMessage)
		m := &ClientMessage{
			Content: receivedMessage,
			Client:  client,
		}
		//if reading from socket fails the for loop is broken and the socket is closed
		if err != nil {
			return
		}

		messageQueue <- m
	}
}

// Send messages to the websocket connection
func (client *client) Write() {
	for {
		for msg := range client.MessageQueue {
			err := client.WriteJSON(msg.Content)
			//if reading from socket fails the for loop is broken and the socket is closed
			if err != nil {
				return
			}
		}
	}
}
