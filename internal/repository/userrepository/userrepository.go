package userrepository

import (
	"chat-v2/internal/repository"
	"chat-v2/internal/user"
)

type UserRepository struct {
	repository.DBRepository
}

func (r *UserRepository) GetOne(id string) (user *user.User, err error) {
	return nil, nil
}
