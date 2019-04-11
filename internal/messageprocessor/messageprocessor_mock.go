package messageprocessor

import (
	"chat-v2/internal/client"

	"github.com/stretchr/testify/mock"
)

type Mock struct {
	mock.Mock
}

func (m *Mock) ProcessMessage(message *client.ClientMessage) (processedMessage *client.ClientMessage) {
	args := m.Called(message)
	processedMessage, _ = args.Get(0).(*client.ClientMessage)
	return processedMessage
}
