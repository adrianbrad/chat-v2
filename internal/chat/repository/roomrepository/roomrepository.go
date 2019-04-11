package roomrepository

import (
	"chat-v2/internal/chat/models/room"
	"chat-v2/internal/chat/repository"
)

type UserRepository struct {
	repository.DBRepository
}

func (r *UserRepository) GetAll() (rooms []*room.Room) {
	return nil
}
