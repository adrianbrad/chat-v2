package room

import (
	"sync"

	"github.com/adrianbrad/chat-v2/internal/client"
)

type Room struct {
	ID string

	Clients      map[client.Client]struct{}
	AddClient    chan client.Client
	RemoveClient chan client.Client

	MessageQueue chan map[string]interface{}

	done chan struct{}
}

func New(ID string) *Room {

	room := &Room{
		ID: ID,

		Clients:      make(map[client.Client]struct{}),
		AddClient:    make(chan client.Client),
		RemoveClient: make(chan client.Client),

		MessageQueue: make(chan map[string]interface{}),

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
			r.broadcastMessage(message)
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
