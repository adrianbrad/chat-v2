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

func (m *Mock) Read() {}

//Write has to block execution so we have time to assert the clients map length in chatservice_test.go
func (m *Mock) Write() {
	_ = m.Called()
}

func (m *Mock) AddToMessageQueue(message map[string]interface{}) {
	_ = m.Called(message)
}

func (m *Mock) ConnectionEnded() chan struct{} {
	return m.connectionEnded
}

func (m *Mock) GetUser() *user.User {
	return nil
}
