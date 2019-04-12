package chatservice

import (
	"fmt"
	"sync"
	"testing"

	"github.com/adrianbrad/chat-v2/internal/client"
	"github.com/adrianbrad/chat-v2/internal/room"
	"github.com/adrianbrad/chat-v2/internal/user"

	"github.com/adrianbrad/chat-v2/internal/repository/roomrepository"

	"github.com/adrianbrad/chat-v2/internal/repository/userrepository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setUp(delta int) (roomsSlice []*room.Room, rr *roomrepository.Mock, usr *user.User, ur *userrepository.Mock, service *ChatService, teardown func()) {
	roomsSlice = []*room.Room{
		&room.Room{
			ID:           "room1",
			AddClient:    make(chan client.Client),
			RemoveClient: make(chan client.Client),
		},
		&room.Room{
			ID: "room2",
		},
		&room.Room{
			ID: "room3",
			// AddClient
		},
	}

	go func() {
		<-roomsSlice[0].AddClient
	}()
	go func() {
		<-roomsSlice[0].RemoveClient
	}()

	rr = &roomrepository.Mock{}
	rr.On("GetAll").Return(roomsSlice)

	usr = &user.User{
		ID: "user 1",
	}
	ur = &userrepository.Mock{}
	ur.On("GetOne", mock.Anything).Return(usr, nil)

	rooms := make(map[string]*room.Room, len(roomsSlice))

	for _, room := range roomsSlice {
		//! remember to init channel
		rooms[room.ID] = room
	}

	service = &ChatService{
		userRepository: ur,

		clients: make(map[client.Client]struct{}),
		rooms:   rooms,

		addClient:    make(chan client.Client),
		removeClient: make(chan client.Client),

		createClient: client.NewMock,

		stopChan: make(chan struct{}, 1),
	}

	var wg sync.WaitGroup
	wg.Add(delta + 1)

	go service.run(&wg)

	teardown = func() {
		service.stop()
		wg.Wait()
		close(service.addClient)
		close(service.removeClient)
		close(service.stopChan)
	}

	return
}

func Test_HandleWSConn_InvalidUserID(t *testing.T) {
	service := &ChatService{}

	err := service.HandleWSConn(nil, nil)

	assert.Equal(t, "User ID is not present in data or is not string, data: map[]", err.Error())
}

func Test_HandleWSConn_InvalidRoomID(t *testing.T) {
	service := &ChatService{}

	err := service.HandleWSConn(nil, map[string]interface{}{"userID": "a"})

	assert.Equal(t, "Room ID is not present in data or is not string, data: map[userID:a]", err.Error())
}

func Test_HandleWSConn_ErrorRetrievingUser(t *testing.T) {
	_, _, _, urDefault, _, _ := setUp(0)
	ur := urDefault
	ur.ExpectedCalls = ur.ExpectedCalls[:0]
	service := &ChatService{
		userRepository: ur,
	}

	ur.On("GetOne", mock.Anything).Return(nil, fmt.Errorf("error given by test"))
	err := service.HandleWSConn(nil, map[string]interface{}{"userID": "a", "roomID": "b"})
	assert.Equal(t, "error given by test", err.Error())
}

func Test_HandleWSConn_Success(t *testing.T) {
	_, _, _, _, service, teardown := setUp(2)

	client.InitClientMock()
	client.ClientMock.On("Write").Return().Run(func(mock.Arguments) {
		client.ClientMock.ConnectionEnded() <- struct{}{}
	})
	//We have to run this in parallel and make sure that we have something that blocks during execution, in our case the mockClient.Write method
	service.HandleWSConn(nil, map[string]interface{}{"userID": "a", "roomID": "room1"})

	teardown()

	assert.Equal(t, 0, len(service.clients))
}
