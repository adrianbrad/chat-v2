package client

import (
	"github.com/adrianbrad/chat-v2/internal/user"
	"time"

	"github.com/stretchr/testify/mock"
)

type Mock struct {
	mock.Mock
}

func NewMock(wsConn wsConn, user *user.User) Client {
	return &Mock{}
}

func (m *Mock) Read(messageQueue chan *ClientMessage) {}

//Write has to block execution so we have time to assert the clients map length in chatservice_test.go
func (m *Mock) Write() { time.Sleep(100 * time.Millisecond); return }

func (m *Mock) AddToMessageQueue(message *ClientMessage) {
	_ = m.Called(message)

}
