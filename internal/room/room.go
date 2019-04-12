package room

import (
	"sync"

	"github.com/adrianbrad/chat-v2/internal/client"
)

type messageProcessor interface {
	ProcessMessage(*client.ClientMessage) map[string]interface{}
}

type Room struct {
	messageProcessor

	ID string

	Clients      map[client.Client]struct{}
	AddClient    chan client.Client
	RemoveClient chan client.Client

	MessageQueue chan *client.ClientMessage

	done chan struct{}
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

		done: make(chan struct{}, 1),
	}
	go room.run(nil)

	return room
}

func (r *Room) run(wg *sync.WaitGroup) {
	for {
		select {
		case clientJoins := <-r.AddClient:
			r.Clients[clientJoins] = struct{}{}
			if wg != nil {
				wg.Done()
			}

		case clientLeaves := <-r.RemoveClient:
			delete(r.Clients, clientLeaves)
			if wg != nil {
				wg.Done()
			}

		case message := <-r.MessageQueue:
			processedMessage := r.ProcessMessage(message)

			r.broadcastMessage(processedMessage)
			if wg != nil {
				wg.Done()
			}

		case <-r.done:
			if wg != nil {
				wg.Done()
			}
			return
		}
	}
}

func (r *Room) stop() {
	r.done <- struct{}{}
}

func (r *Room) broadcastMessage(message map[string]interface{}) {
	for client := range r.Clients {
		client.AddToMessageQueue(message)
	}
}
