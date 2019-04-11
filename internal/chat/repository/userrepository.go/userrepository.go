package userrepository

import (
	"chat-v2/internal/chat/models/user"
	"chat-v2/internal/chat/repository"
)

type UserRepository struct {
	repository.DBRepository
}

func (r *UserRepository) GetOne(id string) (user *user.User, err error) {
	return nil, nil
}
