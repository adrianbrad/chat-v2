package message

type MessageType int

const (
	Text MessageType = iota
	Action
)

var messageTypes = [...]string{"text", "action"}

func (a MessageType) String() string {
	return messageTypes[a]
}
