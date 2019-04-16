package client

import (
	"github.com/adrianbrad/chat-v2/internal/message"
	"github.com/adrianbrad/chat-v2/internal/user"
	log "github.com/sirupsen/logrus"
)

type roomIdentifier struct {
	ID           string
	messageQueue chan message.Message
}

type wsConn interface {
	ReadJSON(v interface{}) error
	WriteJSON(v interface{}) error
	Close() error
}

type Client interface {
	AddToMessageQueue(message message.Message)
	ConnectionEnded() chan error
	GetUser() *user.User
}

type client struct {
	wsConn
	messageProcessor

	user         *user.User
	MessageQueue chan message.Message

	connectionEnded chan error

	roomIdentifier roomIdentifier

	bareMessageFactoryFunc func(message map[string]interface{}) (bareMessage message.BareMessage, err error)
}

func (client *client) run() (err error) {
	for {
		select {
		case err := <-client.connectionEnded:
			log.Info("Ws connection ended")
			return err
		default:

			if _, canSendMessage := client.user.Permissions[user.SendMessage.String()]; canSendMessage {
				client.read()
			}

			client.write()
		}
	}
}

// Proccess messages sent by the websocket connection and forward them to the channel given as parameter
func (client *client) read() {
	var receivedMessage map[string]interface{}
	err := client.ReadJSON(&receivedMessage)
	//if reading from socket fails the for loop is broken and the socket is closed
	if err != nil {
		client.stop(err)
		return
	}

	var processedMessage message.Message
	defer func() {
		client.roomIdentifier.messageQueue <- processedMessage
	}()

	receivedMessage["room_id"] = client.roomIdentifier.ID

	bareMessage, err := client.bareMessageFactoryFunc(receivedMessage)
	if err != nil {
		processedMessage.Error = err.Error()
		return
	}

	processedMessage, err = client.ProcessMessage(bareMessage)
	if err != nil {
		processedMessage.Error = err.Error()
		return
	}
	return
}

// Send messages to the websocket connection
// another implementation is with for msg := range client.MessageQueue
func (client *client) write() {
	select {
	case msg := <-client.MessageQueue:

		err := client.WriteJSON(msg)
		//if writing from socket fails the for loop is broken and the socket is closed
		if err != nil {
			client.stop(err)
		}
	default:
	}
}

func (client *client) stop(err error) {
	client.connectionEnded <- err
}

func (client *client) ConnectionEnded() chan error {
	return client.connectionEnded
}

func (client *client) AddToMessageQueue(message message.Message) {
	client.MessageQueue <- message
}

func (client *client) GetUser() *user.User {
	return client.user
}
