package message

import (
	"fmt"
)

type messageRepository interface {
	Insert(bareMessage BareMessage) (message Message, err error)
}

type MessageProcessor struct {
	messageRepository messageRepository
}

func (m *MessageProcessor) ProcessMessage(bareMessage BareMessage) (message Message, err error) {
	switch bareMessage.Action {
	case Text.String():
		return m.processTextMessage(bareMessage)
	case Action.String():
		// return m.processActionMessage(messageContent, message.User.ID)
	default:
		err = fmt.Errorf("Message type is invalid")
		return
	}

	return
}

func (m *MessageProcessor) processTextMessage(message BareMessage) (processedMessage Message, err error) {
	return m.messageRepository.Insert(message)
}

// func (m *MessageProcessor) processActionMessage(message map[string]interface{}, userID string) (processedMessage map[string]interface{}, err error) {
// 	return
// }
