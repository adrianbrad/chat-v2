package room

import (
	"testing"
	"time"

	"github.com/adrianbrad/chat-v2/internal/client"
	"github.com/adrianbrad/chat-v2/internal/messageprocessor"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setUp() (messageProcessorMock *messageprocessor.Mock, room *Room) {
	messageProcessorMock = &messageprocessor.Mock{}
	roomID := "testing"
	room = New(messageProcessorMock, roomID)
	return
}

func Test_AddClientToRoom(t *testing.T) {
	_, room := setUp()
	c := client.NewClient(nil, nil)
	room.AddClient <- c

	assert.Equal(t, 1, len(room.Clients))
}

func Test_RemoveClientFromRoom(t *testing.T) {
	_, room := setUp()
	c := client.NewClient(nil, nil)
	room.Clients[c] = struct{}{}

	room.RemoveClient <- c
	assert.Equal(t, 0, len(room.Clients))
}

func Test_AddMessageToMessageQueue(t *testing.T) {
	mp, room := setUp()
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
	time.Sleep(10 * time.Millisecond)
	assert.Equal(t, 3, messageCount)
}

func Test_SendMessageAfterUserLeavesRoom(t *testing.T) {
	mp, room := setUp()
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
	time.Sleep(10 * time.Millisecond)
	assert.Equal(t, 3, messageCount)
}
