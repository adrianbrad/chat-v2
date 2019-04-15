package client

import (
	"fmt"
	"testing"

	"github.com/adrianbrad/chat-v2/internal/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type wsConnMock struct {
	mock.Mock
}

func (c *wsConnMock) ReadJSON(v interface{}) error {
	args := c.Called(v)
	return args.Error(0)
}

func (c *wsConnMock) WriteJSON(v interface{}) error {
	args := c.Called(v)
	return args.Error(0)
}

func (c *wsConnMock) Close() error {
	args := c.Called()
	return args.Error(0)
}

func Test_Client_Read_Error(t *testing.T) {
	wsConn := &wsConnMock{}
	connectionEndedChan := make(chan error, 1)
	c := &client{
		wsConn:          wsConn,
		connectionEnded: connectionEndedChan,
	}

	wsConn.On("ReadJSON", mock.Anything).Return(fmt.Errorf("err"))

	c.read()

	connEnded := <-connectionEndedChan

	assert.Equal(t, "err", connEnded.Error())
}

func Test_Client_Read_ProcessMessageError(t *testing.T) {
	wsConn := &wsConnMock{}
	mp := message.NewMessageProcessorMock()
	messageQueue := make(chan map[string]interface{}, 1)

	connectionEndedChan := make(chan error, 1)
	c := &client{
		wsConn:           wsConn,
		messageProcessor: mp,

		connectionEnded: connectionEndedChan,
		roomIdentifier:  roomIdentifier{messageQueue: messageQueue},
	}

	wsConn.On("ReadJSON", mock.Anything).Return(nil)
	mp.On("ProcessMessage", mock.Anything).Return(nil, fmt.Errorf("err"))

	c.read()
	receivedMessage := <-messageQueue

	assert.Equal(t, map[string]interface{}{"error": "err"}, receivedMessage)
}

func Test_Client_Write_Error(t *testing.T) {
	wsConn := &wsConnMock{}
	messageQueue := make(chan map[string]interface{}, 1)
	connectionEndedChan := make(chan error, 1)

	c := &client{
		wsConn:          wsConn,
		MessageQueue:    messageQueue,
		connectionEnded: connectionEndedChan,
	}

	//we have to put a value in the channel in order to for loop over
	c.AddToMessageQueue(map[string]interface{}{})
	wsConn.On("WriteJSON", mock.Anything).Return(fmt.Errorf("err"))
	c.write()
	connEnded := <-c.ConnectionEnded()

	assert.Equal(t, "err", connEnded.Error())
}

func Test_Client_Write_Success(t *testing.T) {
	wsConn := &wsConnMock{}
	messageQueue := make(chan map[string]interface{}, 1)
	connectionEndedChan := make(chan error, 1)

	c := &client{
		wsConn:          wsConn,
		MessageQueue:    messageQueue,
		connectionEnded: connectionEndedChan,
	}

	//we have to put a value in the channel in order to for loop over
	c.AddToMessageQueue(map[string]interface{}{})
	wsConn.On("WriteJSON", mock.Anything).Return(nil)
	c.write()
	select {
	case <-c.ConnectionEnded():
		t.Error("Connection should not be eneded")
	default:
		fmt.Println("ok")
	}
}

func Test_Client_Run_Error(t *testing.T) {
	connectionEndedChan := make(chan error, 1)

	c := &client{
		connectionEnded: connectionEndedChan,
	}

	c.ConnectionEnded() <- fmt.Errorf("err")
	err := c.run()
	assert.Equal(t, "err", err.Error())
}

func Test_Client_Run_SuccessCycle(t *testing.T) {
	connectionEndedChan := make(chan error, 1)
	wsConn := &wsConnMock{}
	mp := message.NewMessageProcessorMock()
	mp.On("ProcessMessage", mock.Anything).Return(map[string]interface{}{}, nil)

	c := &client{
		messageProcessor: mp,
		wsConn:           wsConn,
		connectionEnded:  connectionEndedChan,
	}

	wsConn.On("WriteJSON", mock.Anything).Return(nil)
	wsConn.On("ReadJSON", mock.Anything).Return(nil)

	err := make(chan error, 1)
	go func() {
		e := c.run()
		err <- e
	}()

	c.stop(fmt.Errorf("Controlled error"))
	e := <-err
	assert.Equal(t, "Controlled error", e.Error())
}
