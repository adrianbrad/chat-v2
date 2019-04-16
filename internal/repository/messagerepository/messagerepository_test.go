package messagerepository

import (
	"testing"

	"github.com/adrianbrad/chat-v2/test"

	"github.com/adrianbrad/chat-v2/internal/user"

	"github.com/adrianbrad/chat-v2/internal/message"
	"github.com/stretchr/testify/assert"
)

func Test_MessageRepository(t *testing.T) {
	t.Run("a=Insert_Fail", func(t *testing.T) {
		db, err := test.SetupTestDB()
		if err != nil {
			t.Fatal(err.Error())
		}
		defer db.Close()

		messageRepo := MessageRepositoryDB{db}

		insert_fail_no_user(t, messageRepo)
		insert_fail(t, messageRepo)
		insert_success(t, messageRepo)
	})
}

func insert_fail_no_user(t *testing.T, messageRepo MessageRepositoryDB) {
	bareMessage := message.BareMessage{}

	_, err := messageRepo.Insert(bareMessage)

	assert.Equal(t, "User should not be nil", err.Error())
}

func insert_fail(t *testing.T, messageRepo MessageRepositoryDB) {
	bareMessage := message.BareMessage{User: &user.User{}}

	_, err := messageRepo.Insert(bareMessage)

	assert.Equal(t, `pq: insert or update on table "messages" violates foreign key constraint "fk_messages_rooms"`, err.Error())
}

func insert_success(t *testing.T, messageRepo MessageRepositoryDB) {
	bareMessage := message.BareMessage{User: &user.User{ID: "user_a"}, RoomID: "room_a", Body: message.MessageBody{Content: "a"}}

	message, err := messageRepo.Insert(bareMessage)

	assert.Nil(t, err)
	assert.Equal(t, bareMessage, message.BareMessage)
	assert.NotEmpty(t, message.ID)
	assert.NotEmpty(t, message.SentAt)
}
