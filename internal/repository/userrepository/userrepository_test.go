package userrepository

import (
	"testing"

	"github.com/adrianbrad/chat-v2/internal/user"
	"github.com/adrianbrad/chat-v2/test"
	"github.com/stretchr/testify/assert"
)

func Test_UserRepository(t *testing.T) {
	t.Run("a=GetOne", func(t *testing.T) {
		db, err := test.SetupTestDB()
		if err != nil {
			t.Fatal(err.Error())
		}
		defer db.Close()

		userRepo := UserRepositoryDB{db}

		getOneSuccess(t, userRepo)
		getOneFail(t, userRepo)
	})
}

func getOneSuccess(t *testing.T, userRepo UserRepositoryDB) {
	u, err := userRepo.GetOne("user_a")
	assert.Nil(t, err)

	expectedUser := user.User{
		ID:       "user_a",
		Nickname: "someone",
		Permissions: map[string]struct{}{
			"talk": struct{}{},
			"run":  struct{}{},
		},
	}
	assert.Equal(t, expectedUser, *u)
}

func getOneFail(t *testing.T, userRepo UserRepositoryDB) {
	u, err := userRepo.GetOne("user_b")

	assert.Equal(t, "sql: no rows in result set", err.Error())

	var emptyUser user.User
	assert.Equal(t, emptyUser, *u)
}
