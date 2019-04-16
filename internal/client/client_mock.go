package client

import (
	"github.com/adrianbrad/chat-v2/internal/message"
	"github.com/adrianbrad/chat-v2/internal/user"

	"github.com/stretchr/testify/mock"
)

var ClientMock *Mock

type Mock struct {
	mock.Mock
	connectionEnded chan error
}

func InitClientMock() {
	ClientMock = &Mock{
		connectionEnded: make(chan error, 1),
	}
}

func (m *Mock) Read() {}

//Write has to block execution so we have time to assert the clients map length in chatservice_test.go
func (m *Mock) Write() {
	_ = m.Called()
}

func (m *Mock) AddToMessageQueue(message message.Message) {
	_ = m.Called(message)
}

func (m *Mock) ConnectionEnded() chan error {
	return m.connectionEnded
}

func (m *Mock) GetUser() *user.User {
	args := m.Called()
	return args.Get(0).(*user.User)
}
