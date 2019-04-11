package websocketshandler

import (
	testutils "chat-v2/test/utils"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func Test_ServeHTTP_InvalidID(t *testing.T) {
	var upgrader *websocket.Upgrader
	getDataFunc := func(r *http.Request) (data map[string]interface{}, err error) {
		return nil, fmt.Errorf("err")
	}

	wh := NewWebsocketsHandler(
		&serviceMock{},
		upgrader,
		getDataFunc,
	)

	r := testutils.NewTestRequest(t, "", "", nil)
	rr := httptest.NewRecorder()

	wh.ServeHTTP(rr, r)

	bodyBytes := testutils.ReadRequestBody(t, rr.Body)

	assert.Equal(t, rr.Code, http.StatusBadRequest)
	assert.Equal(t, "Could not get Data from websocket request\n", string(bodyBytes))
}

func Test_ServeHTTP_ErrorWhileUpgradingToWebsocketSession(t *testing.T) {
	service := &serviceMock{}

	upgrader := &upgraderMock{}
	upgrader.On("Upgrade",
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).Return(nil, fmt.Errorf("err"))

	getDataFunc := func(r *http.Request) (data map[string]interface{}, err error) {
		return make(map[string]interface{}), nil
	}

	wh := NewWebsocketsHandler(
		service,
		upgrader,
		getDataFunc,
	)

	r := testutils.NewTestRequest(t, "", "", nil)

	rr := httptest.NewRecorder()

	wh.ServeHTTP(rr, r)

	bodyBytes := testutils.ReadRequestBody(t, rr.Body)

	assert.Equal(t, rr.Code, http.StatusInternalServerError)
	assert.Equal(t, "Error while upgrading to websocket\n", string(bodyBytes))
}

func Test_ServeHTTP_ErrorWhileHandlingNewConnection(t *testing.T) {
	service := &serviceMock{}
	service.On("HandleWSConn", mock.Anything, mock.Anything).Return(fmt.Errorf("err"))

	upgrader := &upgraderMock{}
	upgrader.On("Upgrade",
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).Return(&websocket.Conn{}, nil)

	getDataFunc := func(r *http.Request) (data map[string]interface{}, err error) {
		return make(map[string]interface{}), nil
	}

	wh := NewWebsocketsHandler(
		service,
		upgrader,
		getDataFunc,
	)

	r := testutils.NewTestRequest(t, "", "", nil)

	rr := httptest.NewRecorder()

	wh.ServeHTTP(rr, r)

	bodyBytes := testutils.ReadRequestBody(t, rr.Body)

	assert.Equal(t, rr.Code, http.StatusInternalServerError)
	assert.Equal(t, "Error while handling websocket session\n", string(bodyBytes))
}

func Test_ServeHTTP_Succes(t *testing.T) {
	service := &serviceMock{}
	service.On("HandleWSConn", mock.Anything, mock.Anything).Return(nil)

	upgrader := &upgraderMock{}
	upgrader.On("Upgrade",
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).Return(&websocket.Conn{}, nil)

	getDataFunc := func(r *http.Request) (data map[string]interface{}, err error) {
		return make(map[string]interface{}), nil
	}

	wh := NewWebsocketsHandler(
		service,
		upgrader,
		getDataFunc,
	)

	r := testutils.NewTestRequest(t, "", "", nil)

	rr := httptest.NewRecorder()

	wh.ServeHTTP(rr, r)

	bodyBytes := testutils.ReadRequestBody(t, rr.Body)

	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, "", string(bodyBytes))
}
