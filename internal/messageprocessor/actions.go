package messageprocessor

type MessageType int

const (
	Message MessageType = iota
	Action
)

var messageTypes = [...]string{"message", "action"}

func (a MessageType) String() string {
	return messageTypes[a]
}
