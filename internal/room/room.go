package room

import (
	"sync"

	"github.com/adrianbrad/chat-v2/internal/client"
	"github.com/adrianbrad/chat-v2/internal/message"
	log "github.com/sirupsen/logrus"
)

type Room struct {
	ID string

	Clients      map[client.Client]struct{}
	AddClient    chan client.Client
	RemoveClient chan client.Client

	MessageQueue chan message.Message

	done chan struct{}
}

func New(ID string) *Room {

	room := &Room{
		ID: ID,

		Clients:      make(map[client.Client]struct{}),
		AddClient:    make(chan client.Client),
		RemoveClient: make(chan client.Client),

		MessageQueue: make(chan message.Message),

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
			log.Infof("Added user with id: %s to room with id: %s", clientJoins.GetUser().ID, r.ID)

		case clientLeaves := <-r.RemoveClient:
			delete(r.Clients, clientLeaves)
			if wg != nil {
				wg.Done()
			}
			log.Infof("Removed user with id: %s from room with id: %s", clientLeaves.GetUser().ID, r.ID)

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

func (r *Room) broadcastMessage(message message.Message) {
	for client := range r.Clients {
		client.AddToMessageQueue(message)
	}
}
