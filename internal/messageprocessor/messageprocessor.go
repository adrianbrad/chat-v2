package messageprocessor

import (
	"fmt"
	"time"

	"github.com/adrianbrad/chat-v2/internal/client"
)

type messageRepository interface {
	Insert(message map[string]interface{}) (messageID int, sentAt time.Time)
}

type MessageProcessor struct {
	messageRepository messageRepository
}

func (m *MessageProcessor) ProcessMessage(message *client.ClientMessage) (processedMessage map[string]interface{}, err error) {
	messageContent := message.Content
	messageType, ok := messageContent["action"].(string)
	if !ok {
		err = fmt.Errorf("Action not present or not string")
		return
	}

	switch messageType {
	case Message.String():
	case Action.String():
	}
	return
}
