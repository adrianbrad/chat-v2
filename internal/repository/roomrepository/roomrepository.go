package roomrepository

import (
	"github.com/adrianbrad/chat-v2/internal/repository"
	"github.com/adrianbrad/chat-v2/internal/room"
)

type UserRepository struct {
	repository.DBRepository
}

func (r *UserRepository) GetAll() (rooms []*room.Room) {
	return nil
}
