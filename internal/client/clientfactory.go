package client

import (
	"github.com/adrianbrad/chat-v2/internal/message"
	"github.com/adrianbrad/chat-v2/internal/user"
)

type messageProcessor interface {
	ProcessMessage(bareMessage message.BareMessage) (message message.Message, err error)
}

type BareMessageFactoryFunc func(message map[string]interface{}) (bareMessage message.BareMessage, err error)

type CreateFunc func(wsConn wsConn, user *user.User, roomID string, roomMessageQueue chan message.Message) Client

func (c CreateFunc) Create(wsConn wsConn, user *user.User, roomID string, roomMessageQueue chan message.Message) Client {
	return c(wsConn, user, roomID, roomMessageQueue)
}

type Factory interface {
	Create(wsConn wsConn, user *user.User, roomID string, roomMessageQueue chan message.Message) Client
}

type factory struct {
	messageProcessor       messageProcessor
	bareMessageFactoryFunc BareMessageFactoryFunc
}

func NewFactory(messageProcessor messageProcessor, bareMessageFactoryFunc BareMessageFactoryFunc) Factory {
	return &factory{
		messageProcessor:       messageProcessor,
		bareMessageFactoryFunc: bareMessageFactoryFunc,
	}
}

func (f *factory) Create(wsConn wsConn, user *user.User, roomID string, roomMessageQueue chan message.Message) Client {
	c := &client{
		wsConn:          wsConn,
		user:            user,
		MessageQueue:    make(chan message.Message, 1),
		connectionEnded: make(chan error, 1),
		roomIdentifier: roomIdentifier{
			ID:           roomID,
			messageQueue: roomMessageQueue,
		},
		messageProcessor:       f.messageProcessor,
		bareMessageFactoryFunc: f.bareMessageFactoryFunc,
	}

	return c
}

func NewTestingFactory() Factory {
	return CreateFunc(func(wsConn wsConn, user *user.User, roomID string, roomMessageQueue chan message.Message) Client {
		return ClientMock
	})
}
