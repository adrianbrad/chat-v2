package chatservice

import (
	"fmt"
	"sync"

	"github.com/adrianbrad/chat-v2/internal/client"
	"github.com/adrianbrad/chat-v2/internal/room"
	"github.com/adrianbrad/chat-v2/internal/user"

	"github.com/gorilla/websocket"
)

type userRepository interface {
	GetOne(id string) (user *user.User, err error)
}

type roomRepository interface {
	GetAll() []*room.Room
}

type ChatService struct {
	userRepository userRepository
	roomRepository roomRepository

	clients map[client.Client]struct{}
	rooms   map[string]*room.Room

	addClient    chan client.Client
	removeClient chan client.Client

	createClient client.FactoryMethod

	stopChan chan struct{}
}

func NewChatService(
	userRepository userRepository,
	roomRepository roomRepository,
	createClientFactoryMethod client.FactoryMethod,
) *ChatService {

	repoRooms := roomRepository.GetAll()
	rooms := make(map[string]*room.Room, len(repoRooms))
	for _, room := range repoRooms {
		//! remember to init room channels
		rooms[room.ID] = room
	}

	cs := &ChatService{
		userRepository: userRepository,
		roomRepository: roomRepository,

		clients: make(map[client.Client]struct{}),
		rooms:   rooms,

		addClient:    make(chan client.Client),
		removeClient: make(chan client.Client),

		createClient: createClientFactoryMethod,

		stopChan: make(chan struct{}, 1),
	}
	go cs.run(nil)
	return cs
}

func (c *ChatService) run(wg *sync.WaitGroup) {
	for {
		select {
		case client := <-c.addClient:
			c.clients[client] = struct{}{}
			if wg != nil {
				wg.Done()
			}

		case client := <-c.removeClient:
			delete(c.clients, client)
			if wg != nil {
				wg.Done()
			}

		case <-c.stopChan:
			// if wg != nil {
			// 	wg.Done()
			// }
			return
		}
	}
}

func (c *ChatService) stop() {
	c.stopChan <- struct{}{}
}

func (c *ChatService) HandleWSConn(wsConn *websocket.Conn, data map[string]interface{}) (err error) {
	userID, ok := data["userID"].(string)
	if !ok {
		err = fmt.Errorf("User ID is not present in data or is not string, data: %+v", data)
		return
	}

	roomID, ok := data["roomID"].(string)
	if !ok {
		err = fmt.Errorf("Room ID is not present in data or is not string, data: %+v", data)
		return
	}

	user, err := c.userRepository.GetOne(userID)
	if err != nil {
		return
	}

	room, ok := c.rooms[roomID]
	if !ok {
		err = fmt.Errorf("Room with id: %s does not exist", roomID)
		return
	}

	client := c.createClient(wsConn, user, room.ID, room.MessageQueue)

	c.addClient <- client
	room.AddClient <- client
	defer func() {
		c.removeClient <- client
		c.rooms[roomID].RemoveClient <- client
	}()

	//We block execution until the websocket connection ended
	<-client.ConnectionEnded()
	return
}
