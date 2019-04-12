package messageprocessor

import (
	"time"

	"github.com/adrianbrad/chat-v2/internal/client"
)

type messageRepository interface {
	Insert(message map[string]interface{}) (messageID int, sentAt time.Time)
}

type MessageProcessor struct {
	messageRepository messageRepository
}

func (m *MessageProcessor) ProcessMessage(message *client.ClientMessage) (processedMessage map[string]interface{}) {
	return nil
}
