package room

import (
	"chat-v2/internal/client"
)

type messageProcessor interface {
	ProcessMessage(*client.ClientMessage) *client.ClientMessage
}

type Room struct {
	messageProcessor

	ID string

	Clients      map[client.Client]struct{}
	AddClient    chan client.Client
	RemoveClient chan client.Client

	MessageQueue chan *client.ClientMessage
}

func New(
	messageProcessor messageProcessor,
	ID string,
) *Room {

	room := &Room{
		messageProcessor: messageProcessor,

		ID: ID,

		Clients:      make(map[client.Client]struct{}),
		AddClient:    make(chan client.Client),
		RemoveClient: make(chan client.Client),

		MessageQueue: make(chan *client.ClientMessage),
	}
	go room.run()

	return room
}

func (r *Room) run() {
	for {
		select {
		case clientJoins := <-r.AddClient:
			r.Clients[clientJoins] = struct{}{}

		case clientLeaves := <-r.RemoveClient:
			delete(r.Clients, clientLeaves)

		case message := <-r.MessageQueue:
			processedMessage := r.ProcessMessage(message)

			r.broadcastMessage(processedMessage)
		}
	}
}

func (r *Room) broadcastMessage(message *client.ClientMessage) {

	for client := range r.Clients {
		client.AddToMessageQueue(message)
	}
}
