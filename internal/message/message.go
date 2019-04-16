package message

import (
	"time"

	"github.com/adrianbrad/chat-v2/internal/user"
)

type MessageBody struct {
	Content string `json:"content,omitempty"`
}

type BareMessage struct {
	Action string      `json:"action,omitempty"`
	Body   MessageBody `json:"body,omitempty"`
	User   *user.User  `json:"user,omitempt"`
	RoomID string      `json:"room_id,omitempty"`
}

type Message struct {
	BareMessage
	ID     int       `json:"id,omitempty"`
	SentAt time.Time `json:"sent_at,omitempty"`
	Error  string    `json:"error,omitempty"`
}

func NewBareMessage(message map[string]interface{}) (bareMessage BareMessage, err error) {
	return
}
