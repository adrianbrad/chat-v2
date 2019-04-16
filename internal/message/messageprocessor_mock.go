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

func (m *Mock) ProcessMessage(bareMessage BareMessage) (message Message, err error) {
	args := m.Called(bareMessage)
	err = args.Error(1)
	if err != nil {
		return
	}
	message, _ = args.Get(0).(Message)
	return
}
