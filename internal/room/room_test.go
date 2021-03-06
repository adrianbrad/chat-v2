package room

import (
	"sync"
	"testing"

	"github.com/adrianbrad/chat-v2/internal/user"

	"github.com/adrianbrad/chat-v2/internal/client"
	"github.com/adrianbrad/chat-v2/internal/message"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

//delta - number of times you expect values to be offloaded from channels in the room.run method
func setUp(delta int) (room *Room, teardown func()) {
	roomID := "testing"
	room = &Room{
		ID:           roomID,
		MessageQueue: make(chan message.Message),

		Clients:      make(map[client.Client]struct{}),
		AddClient:    make(chan client.Client),
		RemoveClient: make(chan client.Client),

		done: make(chan struct{}, 1),
	}

	var wg sync.WaitGroup
	wg.Add(delta + 1)

	go room.run(&wg)

	teardown = func() {
		room.stop()
		wg.Wait()
		close(room.AddClient)
		close(room.RemoveClient)
		close(room.MessageQueue)
		close(room.done)
	}

	return
}

//Using the wait group to bypass the race condition generated by reading the clients map in the testing function, map which was modified by the run() method in a goroutine
func Test_AddClientToRoom(t *testing.T) {
	room, teardown := setUp(1)

	c := &client.Mock{}
	c.On("GetUser").Return(&user.User{ID: "test_user"})
	room.AddClient <- c

	teardown()

	assert.Equal(t, 1, len(room.Clients))
}

func Test_RemoveClientFromRoom(t *testing.T) {
	room, teardown := setUp(1)

	c := &client.Mock{}
	c.On("GetUser").Return(&user.User{ID: "test_user"})
	room.Clients[c] = struct{}{}

	room.RemoveClient <- c

	teardown()

	assert.Equal(t, 0, len(room.Clients))
}

func Test_AddMessageToMessageQueue(t *testing.T) {
	room, teardown := setUp(1)
	client1 := &client.Mock{}
	client2 := &client.Mock{}
	client1.On("GetUser").Return(&user.User{ID: "test_user1"})
	client2.On("GetUser").Return(&user.User{ID: "test_user2"})

	senderClient := &client.Mock{}

	var messageCount int
	incrementMessageCount := func(args mock.Arguments) {
		messageCount++
	}

	client1.On("AddToMessageQueue", mock.Anything).Return().Run(incrementMessageCount)
	client2.On("AddToMessageQueue", mock.Anything).Return().Run(incrementMessageCount)
	senderClient.On("AddToMessageQueue", mock.Anything).Return().Run(incrementMessageCount)

	room.Clients[client1] = struct{}{}
	room.Clients[client2] = struct{}{}
	room.Clients[senderClient] = struct{}{}

	message := message.Message{}
	room.MessageQueue <- message

	teardown()

	assert.Equal(t, 3, messageCount)
}

func Test_SendMessageAfterUserLeavesRoom(t *testing.T) {
	room, teardown := setUp(2)
	client1 := &client.Mock{}
	client2 := &client.Mock{}
	client1.On("GetUser").Return(&user.User{ID: "test_user1"})
	client2.On("GetUser").Return(&user.User{ID: "test_user2"})

	senderClient := &client.Mock{}

	var messageCount int
	incrementMessageCount := func(args mock.Arguments) {
		// message := args.Get(0).(*client.ClientMessage)
		// assert.Equal(t, message.Client, senderClient)
		messageCount++
	}

	client1.On("AddToMessageQueue", mock.Anything).Return().Run(incrementMessageCount)
	client2.On("AddToMessageQueue", mock.Anything).Return().Run(incrementMessageCount)
	senderClient.On("AddToMessageQueue", mock.Anything).Return().Run(incrementMessageCount)

	room.Clients[client1] = struct{}{}
	room.Clients[client2] = struct{}{}
	room.Clients[senderClient] = struct{}{}

	message := message.Message{}

	room.RemoveClient <- client2

	room.MessageQueue <- message

	teardown()

	assert.Equal(t, 2, messageCount)
}
