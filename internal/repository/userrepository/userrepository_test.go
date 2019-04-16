package userrepository

import (
	"testing"

	"github.com/adrianbrad/chat-v2/internal/user"
	"github.com/adrianbrad/chat-v2/test"
	"github.com/stretchr/testify/assert"
)

func Test_UserRepository_GetOne(t *testing.T) {
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

func Test_UserRepository_Create(t *testing.T) {
	t.Run("a=Create", func(t *testing.T) {
		db, err := test.SetupTestDB()
		if err != nil {
			t.Fatal(err.Error())
		}
		defer db.Close()

		userRepo := UserRepositoryDB{db}

		insertOneSuccess(t, userRepo)
		insertOneFail(t, userRepo)
	})
}

func Test_UserRepository_Update(t *testing.T) {
	t.Run("a=Update", func(t *testing.T) {
		db, err := test.SetupTestDB()
		if err != nil {
			t.Fatal(err.Error())
		}
		defer db.Close()

		userRepo := UserRepositoryDB{db}

		updateOneSuccess(t, userRepo)
		updateOneFail(t, userRepo)
	})
}

func Test_UserRepository_Delete(t *testing.T) {
	t.Run("a=Delete", func(t *testing.T) {
		db, err := test.SetupTestDB()
		if err != nil {
			t.Fatal(err.Error())
		}
		defer db.Close()

		userRepo := UserRepositoryDB{db}

		deleteOneSuccess(t, userRepo)
		deleteOneFail(t, userRepo)
	})
}

func getOneSuccess(t *testing.T, userRepo UserRepositoryDB) {
	u, err := userRepo.GetOne("user_a")
	assert.Nil(t, err)

	expectedUser := user.User{
		ID:       "user_a",
		Nickname: "someone",
		Permissions: map[string]struct{}{
			"send_message": struct{}{},
			"run":          struct{}{},
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

func insertOneSuccess(t *testing.T, userRepo UserRepositoryDB) {
	err := userRepo.Create(user.User{ID: "aaa", Nickname: "brad"})
	assert.Nil(t, err)
}

func insertOneFail(t *testing.T, userRepo UserRepositoryDB) {
	err := userRepo.Create(user.User{ID: "user_a", Nickname: "brad"})
	assert.Equal(t, `pq: duplicate key value violates unique constraint "users_pkey"`, err.Error())
}

func updateOneSuccess(t *testing.T, userRepo UserRepositoryDB) {
	err := userRepo.Update(
		user.User{
			ID:       "user_a",
			Nickname: "brad",
			Permissions: map[string]struct{}{
				user.SendMessage.String(): struct{}{},
			}})

	assert.Nil(t, err)
}

func updateOneFail(t *testing.T, userRepo UserRepositoryDB) {
	err := userRepo.Update(user.User{ID: "aaa", Nickname: "brad"})
	assert.Equal(t, `Invalid number of rows affected by create: 0`, err.Error())
}

func deleteOneSuccess(t *testing.T, userRepo UserRepositoryDB) {
	err := userRepo.Delete("user_b")
	assert.Nil(t, err)
}

func deleteOneFail(t *testing.T, userRepo UserRepositoryDB) {
	err := userRepo.Delete("random")
	assert.Equal(t, `Invalid number of rows affected by delete: 0`, err.Error())
}
