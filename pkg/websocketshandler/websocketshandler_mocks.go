package websocketshandler

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/mock"
)

type serviceMock struct {
	mock.Mock
}

func (s *serviceMock) HandleWSConn(wsConn *websocket.Conn, data map[string]interface{}) (err error) {
	args := s.Called(wsConn, data)
	err = args.Error(0)
	return
}

type upgraderMock struct {
	mock.Mock
}

func (s *upgraderMock) Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (wsConn *websocket.Conn, err error) {
	args := s.Called(w, r, responseHeader)
	err = args.Error(1)
	if err != nil {
		return nil, err
	}
	i := args.Get(0)
	result := i.(*websocket.Conn)
	return result, nil
}
