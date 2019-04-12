package client

import (
	"github.com/adrianbrad/chat-v2/internal/user"
)

type roomIdentifier struct {
	ID           string
	messageQueue chan *ClientMessage
}

type FactoryMethod func(wsConn wsConn, user *user.User, roomID string, roomMessageQueue chan *ClientMessage) Client

type wsConn interface {
	ReadJSON(v interface{}) error
	WriteJSON(v interface{}) error
	Close() error
}

type Client interface {
	AddToMessageQueue(message map[string]interface{})
	ConnectionEnded() chan struct{}
	GetUser() *user.User
}

type client struct {
	wsConn

	user         *user.User
	MessageQueue chan map[string]interface{}

	validWsConn     bool
	connectionEnded chan struct{}

	roomIdentifier roomIdentifier
}

func NewClient(wsConn wsConn, user *user.User, roomID string, roomMessageQueue chan *ClientMessage) Client {
	c := &client{
		wsConn:          wsConn,
		user:            user,
		validWsConn:     true,
		connectionEnded: make(chan struct{}, 1),
		roomIdentifier: roomIdentifier{
			ID:           roomID,
			messageQueue: roomMessageQueue,
		},
	}

	go c.Read()
	go c.Write()
	return c
}

// Proccess messages sent by the websocket connection and forward them to the channel given as parameter
func (client *client) Read() {
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
		client.roomIdentifier.messageQueue <- m
	}
}

// Send messages to the websocket connection
func (client *client) Write() {
	for client.validWsConn {
		for msg := range client.MessageQueue {
			err := client.WriteJSON(msg)
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

func (client *client) AddToMessageQueue(message map[string]interface{}) {
	client.MessageQueue <- message
}

func (client *client) GetUser() *user.User {
	return client.user
}
