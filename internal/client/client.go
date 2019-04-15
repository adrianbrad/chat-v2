package client

import (
	"github.com/adrianbrad/chat-v2/internal/message"
	"github.com/adrianbrad/chat-v2/internal/user"
	log "github.com/sirupsen/logrus"
)

type roomIdentifier struct {
	ID           string
	messageQueue chan map[string]interface{}
}

type wsConn interface {
	ReadJSON(v interface{}) error
	WriteJSON(v interface{}) error
	Close() error
}

type Client interface {
	AddToMessageQueue(message map[string]interface{})
	ConnectionEnded() chan error
	GetUser() *user.User
}

type client struct {
	wsConn
	messageProcessor

	user         *user.User
	MessageQueue chan map[string]interface{}

	connectionEnded chan error

	roomIdentifier roomIdentifier
}

func (client *client) run() (err error) {
	for {
		select {
		case err := <-client.connectionEnded:
			log.Info("Ws connection ended")
			return err
		default:
			client.read()
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

	messageToBeProcessed := &message.UserMessage{
		Content: receivedMessage,
		User:    client.user,
	}

	processedMessage, err := client.ProcessMessage(messageToBeProcessed)
	if err != nil {
		processedMessage = map[string]interface{}{
			"error": err.Error(),
		}
	}

	client.roomIdentifier.messageQueue <- processedMessage
}

// Send messages to the websocket connection
//another implementation is with for msg := range client.MessageQueue
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

func (client *client) AddToMessageQueue(message map[string]interface{}) {
	client.MessageQueue <- message
}

func (client *client) GetUser() *user.User {
	return client.user
}
