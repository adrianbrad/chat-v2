package room

import (
	"sync"
	"testing"

	"github.com/adrianbrad/chat-v2/internal/client"
	"github.com/adrianbrad/chat-v2/internal/messageprocessor"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

//delta - number of times you expect values to be offloaded from channels in the room.run method
func setUp(delta int) (messageProcessorMock *messageprocessor.Mock, room *Room, teardown func()) {
	messageProcessorMock = &messageprocessor.Mock{}
	roomID := "testing"
	room = &Room{
		messageProcessor: messageProcessorMock,

		ID: roomID,

		Clients:      make(map[client.Client]struct{}),
		AddClient:    make(chan client.Client),
		RemoveClient: make(chan client.Client),

		MessageQueue: make(chan *client.ClientMessage),

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
	_, room, teardown := setUp(1)

	c := client.NewClient(nil, nil)
	room.AddClient <- c

	teardown()

	assert.Equal(t, 1, len(room.Clients))
}

func Test_RemoveClientFromRoom(t *testing.T) {
	_, room, teardown := setUp(1)

	c := client.NewClient(nil, nil)
	room.Clients[c] = struct{}{}

	room.RemoveClient <- c

	teardown()

	assert.Equal(t, 0, len(room.Clients))
}

func Test_AddMessageToMessageQueue(t *testing.T) {
	mp, room, teardown := setUp(1)
	client1 := &client.Mock{}
	client2 := &client.Mock{}

	senderClient := &client.Mock{}

	var messageCount int
	incrementMessageCount := func(args mock.Arguments) {
		message := args.Get(0).(*client.ClientMessage)
		assert.Equal(t, message.Client, senderClient)
		messageCount++
	}

	client1.On("AddToMessageQueue", mock.Anything).Return().Run(incrementMessageCount)
	client2.On("AddToMessageQueue", mock.Anything).Return().Run(incrementMessageCount)
	senderClient.On("AddToMessageQueue", mock.Anything).Return().Run(incrementMessageCount)

	room.Clients[client1] = struct{}{}
	room.Clients[client2] = struct{}{}
	room.Clients[senderClient] = struct{}{}

	message := &client.ClientMessage{
		Client: senderClient,
	}

	mp.On("ProcessMessage", mock.Anything).Return(message)

	room.MessageQueue <- message

	teardown()

	assert.Equal(t, 3, messageCount)
}

func Test_SendMessageAfterUserLeavesRoom(t *testing.T) {
	mp, room, teardown := setUp(2)
	client1 := &client.Mock{}
	client2 := &client.Mock{}

	senderClient := &client.Mock{}

	var messageCount int
	incrementMessageCount := func(args mock.Arguments) {
		message := args.Get(0).(*client.ClientMessage)
		assert.Equal(t, message.Client, senderClient)
		messageCount++
	}

	client1.On("AddToMessageQueue", mock.Anything).Return().Run(incrementMessageCount)
	client2.On("AddToMessageQueue", mock.Anything).Return().Run(incrementMessageCount)
	senderClient.On("AddToMessageQueue", mock.Anything).Return().Run(incrementMessageCount)

	room.Clients[client1] = struct{}{}
	room.Clients[client2] = struct{}{}
	room.Clients[senderClient] = struct{}{}

	message := &client.ClientMessage{
		Client: senderClient,
	}

	mp.On("ProcessMessage", mock.Anything).Return(message)

	room.RemoveClient <- client2

	room.MessageQueue <- message

	teardown()

	assert.Equal(t, 2, messageCount)
}
