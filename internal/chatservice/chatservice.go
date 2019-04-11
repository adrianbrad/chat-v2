package chatservice

import (
	"chat-v2/internal/client"
	"chat-v2/internal/room"
	"chat-v2/internal/user"
	"fmt"

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

	createClient client.ClientFactoryMethod
}

func NewChatService(
	userRepository userRepository,
	roomRepository roomRepository,
	createClientFactoryMethod client.ClientFactoryMethod,
) *ChatService {

	repoRooms := roomRepository.GetAll()
	rooms := make(map[string]*room.Room, len(repoRooms))
	for _, room := range repoRooms {
		//! remember to init channel
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
	}
	go cs.run()
	return cs
}

func (c *ChatService) run() {
	for {
		select {
		case client := <-c.addClient:
			c.clients[client] = struct{}{}
		case client := <-c.removeClient:
			delete(c.clients, client)
		}
	}
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

	client := c.createClient(wsConn, user)

	c.addClient <- client
	room, ok := c.rooms[roomID]
	if !ok {
		err = fmt.Errorf("Room with id: %s does not exist", roomID)
		return
	}

	room.AddClient <- client
	defer func() {
		c.removeClient <- client
		c.rooms[roomID].RemoveClient <- client
	}()

	//we read the messages from the socket and forward them to the room message queue
	go client.Read(c.rooms[roomID].MessageQueue)
	//we retrieve the messages sent by the room to the client message queue and send them to the client
	// ! the for loop in this function blocks as long as the websocket connection is valid, this satisfies the websockethandler.ServeHTTP condition
	client.Write()
	return
}
