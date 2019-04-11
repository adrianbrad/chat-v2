package roomrepository

import (
	"chat-v2/internal/chat/models/room"

	"github.com/stretchr/testify/mock"
)

type Mock struct {
	mock.Mock
}

func (m *Mock) GetAll() (rooms []*room.Room) {
	args := m.Called()
	return args.Get(0).([]*room.Room)
}
