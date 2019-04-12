package client

import (
	"github.com/adrianbrad/chat-v2/internal/user"

	"github.com/stretchr/testify/mock"
)

var ClientMock *Mock

type Mock struct {
	mock.Mock
	connectionEnded chan struct{}
}

func InitClientMock() {
	ClientMock = &Mock{
		connectionEnded: make(chan struct{}, 1),
	}
}

func NewMock(wsConn wsConn, user *user.User) Client {
	return ClientMock
}

func (m *Mock) Read(messageQueue chan *ClientMessage) {}

//Write has to block execution so we have time to assert the clients map length in chatservice_test.go
func (m *Mock) Write() {
	_ = m.Called()
}

func (m *Mock) AddToMessageQueue(message *ClientMessage) {
	// _ = m.Called(message)
}

func (m *Mock) ConnectionEnded() chan struct{} {
	return m.connectionEnded
}
