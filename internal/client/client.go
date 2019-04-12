package client

import (
	"github.com/adrianbrad/chat-v2/internal/user"
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
	AddToMessageQueue(message *ClientMessage)
	ConnectionEnded() chan struct{}
}

type client struct {
	wsConn

	user         *user.User
	MessageQueue chan *ClientMessage

	validWsConn     bool
	connectionEnded chan struct{}
}

func NewClient(wsConn wsConn, user *user.User) Client {
	return &client{
		wsConn:          wsConn,
		user:            user,
		validWsConn:     true,
		connectionEnded: make(chan struct{}, 1),
	}
}

// Proccess messages sent by the websocket connection and forward them to the channel given as parameter
func (client *client) Read(messageQueue chan *ClientMessage) {
	for client.validWsConn {
		var receivedMessage map[string]interface{}
		err := client.ReadJSON(&receivedMessage)
		//if reading from socket fails the for loop is broken and the socket is closed
		if err != nil {
			client.signalEndConnection()
			return
		}

		m := &ClientMessage{
			Content: receivedMessage,
			Client:  client,
		}
		messageQueue <- m
	}
}

// Send messages to the websocket connection
func (client *client) Write() {
	for client.validWsConn {
		for msg := range client.MessageQueue {
			err := client.WriteJSON(msg.Content)
			//if reading from socket fails the for loop is broken and the socket is closed
			if err != nil {
				client.signalEndConnection()
				return
			}
		}
	}
}

func (client *client) signalEndConnection() {
	client.validWsConn = false
	client.connectionEnded <- struct{}{}
}

func (client *client) ConnectionEnded() chan struct{} {
	return client.connectionEnded
}

func (client *client) AddToMessageQueue(message *ClientMessage) {
	client.MessageQueue <- message
}
