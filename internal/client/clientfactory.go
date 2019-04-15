package client

import (
	"github.com/adrianbrad/chat-v2/internal/user"
)

type messageProcessor interface {
	ProcessMessage(message *ClientMessage) (processedMessage map[string]interface{}, err error)
}

type CreateFunc func(wsConn wsConn, user *user.User, roomID string, roomMessageQueue chan *ClientMessage) Client

func (c CreateFunc) Create(wsConn wsConn, user *user.User, roomID string, roomMessageQueue chan *ClientMessage) Client {
	return c(wsConn, user, roomID, roomMessageQueue)
}

type Factory interface {
	Create(wsConn wsConn, user *user.User, roomID string, roomMessageQueue chan *ClientMessage) Client
}

type factory struct {
	messageProcessor messageProcessor
}

func NewFactory(messageProcessor messageProcessor) Factory {
	return &factory{
		messageProcessor: messageProcessor,
	}
}

func (f *factory) Create(wsConn wsConn, user *user.User, roomID string, roomMessageQueue chan *ClientMessage) Client {
	c := &client{
		wsConn: wsConn,
		user:   user,
		roomIdentifier: roomIdentifier{
			ID:           roomID,
			messageQueue: roomMessageQueue,
		},
		messageProcessor: f.messageProcessor,
	}

	go c.Read()
	go c.Write()

	return c
}

func NewTestingFactory() Factory {
	return CreateFunc(func(wsConn wsConn, user *user.User, roomID string, roomMessageQueue chan *ClientMessage) Client {
		return ClientMock
	})
}
