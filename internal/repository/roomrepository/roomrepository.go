package roomrepository

import (
	"chat-v2/internal/repository"
	"chat-v2/internal/room"
)

type UserRepository struct {
	repository.DBRepository
}

func (r *UserRepository) GetAll() (rooms []*room.Room) {
	return nil
}
