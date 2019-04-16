package messagerepository

import (
	"database/sql"
	"fmt"

	"github.com/adrianbrad/chat-v2/internal/message"
	"github.com/adrianbrad/chat-v2/internal/repository"
)

type MessageRepositoryDB struct {
	repository.DBRepository
}

func NewMessageRepositoryDB(db *sql.DB) *MessageRepositoryDB {
	return &MessageRepositoryDB{db}
}

func (r *MessageRepositoryDB) Insert(bareMesssage message.BareMessage) (message message.Message, err error) {
	if bareMesssage.User == nil {
		err = fmt.Errorf("User should not be nil")
		return
	}
	err = r.QueryRow(`
	INSERT INTO "messages" 
		( "content", "user_id", "room_id")
	VALUES($1, $2, $3)
	RETURNING "message_id", "created_at"
	`,
		bareMesssage.Body.Content,
		bareMesssage.User.ID,
		bareMesssage.RoomID).
		Scan(
			&message.ID,
			&message.SentAt,
		)
	if err != nil {
		return
	}

	message.BareMessage = bareMesssage
	return
}
