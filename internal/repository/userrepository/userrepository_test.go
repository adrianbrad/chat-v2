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

		userRepo := NewUserRepositoryDB(db)

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

		userRepo := &UserRepositoryDB{
			DBRepository: db,
			users:        make(map[string]*user.User),
		}

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

		userRepo := NewUserRepositoryDB(db)

		userRepo.users["user_a"] = &user.User{
			ID:       "user_a",
			Nickname: "someone",
			Permissions: map[string]struct{}{
				user.SendMoney.String():   struct{}{},
				user.SendMessage.String(): struct{}{},
			},
		}

		updatePermissionsSuccess(t, userRepo)

		updateOneSuccess(t, userRepo)
		updateOneFail(t, userRepo)
	})
}

// ! this test somehow hangs the execution
// func Test_UserRepository_Delete(t *testing.T) {
// 	t.Run("a=Delete", func(t *testing.T) {
// 		db, err := test.SetupTestDB()
// 		if err != nil {
// 			t.Fatal(err.Error())
// 		}
// 		defer db.Close()

// 		userRepo := NewUserRepositoryDB(db)

// 		deleteOneSuccess(t, userRepo)
// 		deleteOneFail(t, userRepo)
// 	})
// }

func getOneSuccess(t *testing.T, userRepo *UserRepositoryDB) {
	u, err := userRepo.GetOne("user_a")
	assert.Nil(t, err)

	expectedUser := user.User{
		ID:       "user_a",
		Nickname: "someone",
		Permissions: map[string]struct{}{
			"send_message": struct{}{},
			"send_money":   struct{}{},
		},
	}

	userRepo.stop()

	assert.Equal(t, expectedUser, *u)

	assert.Equal(t, expectedUser, *userRepo.users["user_a"])

}

func getOneFail(t *testing.T, userRepo *UserRepositoryDB) {
	u, err := userRepo.GetOne("inexistent_user")

	assert.Equal(t, "sql: no rows in result set", err.Error())

	var emptyUser user.User
	assert.Equal(t, emptyUser, *u)
}

func insertOneSuccess(t *testing.T, userRepo *UserRepositoryDB) {
	err := userRepo.Create(user.User{ID: "aaa", Nickname: "brad"})
	assert.Nil(t, err)
}

func insertOneFail(t *testing.T, userRepo *UserRepositoryDB) {
	err := userRepo.Create(user.User{ID: "user_a", Nickname: "brad"})
	assert.Equal(t, `pq: duplicate key value violates unique constraint "users_pkey"`, err.Error())
}

func updateOneSuccess(t *testing.T, userRepo *UserRepositoryDB) {
	go userRepo.run()

	updatedUser := user.User{
		ID:       "user_a",
		Nickname: "brad",
		Permissions: map[string]struct{}{
			user.SendMessage.String(): struct{}{},
		},
	}

	err := userRepo.Update(updatedUser)

	userRepo.stop()

	assert.Nil(t, err)

	assert.Equal(t, updatedUser, *userRepo.users["user_a"])
}

func updateOneFail(t *testing.T, userRepo *UserRepositoryDB) {
	err := userRepo.Update(user.User{ID: "aaa", Nickname: "brad"})
	assert.Equal(t, `Invalid number of rows affected by create: 0`, err.Error())
}

func deleteOneSuccess(t *testing.T, userRepo *UserRepositoryDB) {
	err := userRepo.Delete("user_b")
	assert.Nil(t, err)
}

func deleteOneFail(t *testing.T, userRepo *UserRepositoryDB) {
	err := userRepo.Delete("random")
	assert.Equal(t, `Invalid number of rows affected by delete: 0`, err.Error())
}

func updatePermissionsSuccess(t *testing.T, userRepo *UserRepositoryDB) {
	u, err := userRepo.GetOne("user_a")
	assert.Nil(t, err)
	assert.Equal(t, u.Permissions, map[string]struct{}{user.SendMessage.String(): struct{}{}, user.SendMoney.String(): struct{}{}})

	err = userRepo.UpdatePermissions("user_a", []string{user.MuteOthers.String()})

	userRepo.stop()

	assert.Nil(t, err)

	assert.Equal(t, userRepo.users["user_a"].Permissions, map[string]struct{}{user.MuteOthers.String(): struct{}{}})
}
