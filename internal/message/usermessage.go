package message

import "github.com/adrianbrad/chat-v2/internal/user"

type UserMessage struct {
	Content map[string]interface{}
	User    *user.User
}
