package chatservice

import (
	"github.com/adrianbrad/chat-v2/internal/client"
	"github.com/adrianbrad/chat-v2/internal/room"
	"github.com/adrianbrad/chat-v2/internal/user"
	"fmt"
	"testing"
	"time"

	"github.com/adrianbrad/chat-v2/internal/repository/roomrepository"

	"github.com/adrianbrad/chat-v2/internal/repository/userrepository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setUp() (rooms []*room.Room, rr *roomrepository.Mock, usr *user.User, ur *userrepository.Mock) {
	rooms = []*room.Room{
		&room.Room{
			ID: "room1",
		},
		&room.Room{
			ID: "room2",
		},
		&room.Room{
			ID: "room3",
			// AddClient
		},
	}
	rr = &roomrepository.Mock{}
	rr.On("GetAll").Return(rooms)

	usr = &user.User{
		ID: "user 1",
	}
	ur = &userrepository.Mock{}
	ur.On("GetOne", mock.Anything).Return(usr, nil)
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
	_, _, _, urDefault := setUp()
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
	rooms, rr, _, ur := setUp()
	service := NewChatService(
		ur,
		rr,
		client.NewMock,
	)
	handleJoinRoomChan := make(chan client.Client)
	handleLeaveRoomChan := make(chan client.Client)
	rooms[0].AddClient = handleJoinRoomChan
	rooms[0].RemoveClient = handleLeaveRoomChan

	//We have to offload the room channels, otherwise the test will be blocked
	go func() {
		<-handleJoinRoomChan
		return
	}()
	go func() {
		<-handleLeaveRoomChan
		return
	}()

	//We have to run this in parallel and make sure that we have something that blocks during execution, in our case the mockClient.Write method
	go service.HandleWSConn(nil, map[string]interface{}{"userID": "a", "roomID": "room1"})
	time.Sleep(50 * time.Millisecond)
	assert.Equal(t, 1, len(service.clients))
	time.Sleep(60 * time.Millisecond)
	assert.Equal(t, 0, len(service.clients))

}
