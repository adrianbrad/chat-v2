package chatservice

import (
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/adrianbrad/chat-v2/internal/client"
	"github.com/adrianbrad/chat-v2/internal/room"
	"github.com/adrianbrad/chat-v2/internal/user"

	"github.com/gorilla/websocket"
)

type userRepository interface {
	GetOne(id string) (user *user.User, err error)
	Create(user user.User) (err error)
}

type roomRepository interface {
	GetAll() (rooms []*room.Room, err error)
}

type ChatService struct {
	userRepository userRepository
	roomRepository roomRepository

	clientFactory client.Factory

	clients map[client.Client]struct{}
	rooms   map[string]*room.Room

	addClient    chan client.Client
	removeClient chan client.Client

	stopChan chan struct{}
}

func NewChatService(
	userRepository userRepository,
	roomRepository roomRepository,
	clientFactory client.Factory,
) *ChatService {

	repoRooms, err := roomRepository.GetAll()
	if err != nil {
		log.Fatal(err)
	}
	rooms := make(map[string]*room.Room, len(repoRooms))
	for _, room := range repoRooms {
		rooms[room.ID] = room
	}

	log.Infof("Retrieved following rooms from db: %+v", rooms)

	cs := &ChatService{
		userRepository: userRepository,
		roomRepository: roomRepository,

		clientFactory: clientFactory,

		clients: make(map[client.Client]struct{}),
		rooms:   rooms,

		addClient:    make(chan client.Client),
		removeClient: make(chan client.Client),

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
			log.Infof("User %+v joined", *client.GetUser())

			if wg != nil {
				wg.Done()
			}

		case client := <-c.removeClient:
			delete(c.clients, client)
			log.Infof("User %+v left", *client.GetUser())

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

func (c *ChatService) ProcessData(data map[string]interface{}) (processedData map[string]interface{}, err error) {
	userID, ok := data["userID"].(string)
	if !ok || userID == "" {
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

	processedData = make(map[string]interface{})
	processedData["user"] = user
	processedData["room"] = room
	return
}

func (c *ChatService) HandleWSConn(wsConn *websocket.Conn, processedData map[string]interface{}) (err error) {
	room, ok := processedData["room"].(*room.Room)
	if !ok {
		err = fmt.Errorf("Error while retrieving roomID from the message")
		return
	}
	user, ok := processedData["user"].(*user.User)
	if !ok {
		err = fmt.Errorf("Error while retrieving user from the message")
		return
	}

	client := c.clientFactory.Create(wsConn, user, room.ID, room.MessageQueue)

	c.addClient <- client
	room.AddClient <- client
	defer func() {
		c.removeClient <- client
		room.RemoveClient <- client
		client = nil
	}()

	//We block execution until the websocket connection ended
	err = client.Run()
	return
}
