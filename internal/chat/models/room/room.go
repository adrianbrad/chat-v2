package room

import (
	"chat-v2/internal/chat/models/client"
)

type Room struct {
	ID           string
	clients      map[client.Client]struct{}
	AddClient    chan client.Client
	RemoveClient chan client.Client
	MessageQueue chan *client.ClientMessage
}
