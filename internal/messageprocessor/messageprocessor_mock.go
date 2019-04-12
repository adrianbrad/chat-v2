package messageprocessor

import (
	"github.com/adrianbrad/chat-v2/internal/client"

	"github.com/stretchr/testify/mock"
)

type Mock struct {
	mock.Mock
}

func (m *Mock) ProcessMessage(message *client.ClientMessage) (processedMessage map[string]interface{}) {
	args := m.Called(message)
	processedMessage, _ = args.Get(0).(map[string]interface{})
	return processedMessage
}
