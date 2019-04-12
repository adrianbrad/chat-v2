package userrepository

import (
	"github.com/adrianbrad/chat-v2/internal/user"

	"github.com/stretchr/testify/mock"
)

type Mock struct {
	mock.Mock
}

func (m *Mock) GetOne(id string) (usr *user.User, err error) {
	args := m.Called(id)
	err = args.Error(1)
	if err != nil {
		return
	}
	usr = args.Get(0).(*user.User)
	return
}
