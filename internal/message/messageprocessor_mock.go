package message

import (
	"github.com/stretchr/testify/mock"
)

type Mock struct {
	mock.Mock
}

func NewMessageProcessorMock() *Mock {
	return &Mock{}
}

func (m *Mock) ProcessMessage(message *UserMessage) (processedMessage map[string]interface{}, err error) {
	args := m.Called(message)
	err = args.Error(1)
	if err != nil {
		return
	}
	processedMessage, _ = args.Get(0).(map[string]interface{})
	return
}
