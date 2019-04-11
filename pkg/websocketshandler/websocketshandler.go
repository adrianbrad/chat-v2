package websocketshandler

import (
	"net/http"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type service interface {
	HandleWSConn(wsConn *websocket.Conn, data map[string]interface{}) (err error)
}

type upgrader interface {
	Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (wsConn *websocket.Conn, err error)
}

type WebsocketsHandler struct {
	service
	upgrader               upgrader
	getDataFromRequestFunc func(r *http.Request) (data map[string]interface{}, err error)
}

func NewWebsocketsHandler(
	serv service,
	upgrader upgrader,
	getDataFunc func(r *http.Request) (data map[string]interface{}, err error),
) *WebsocketsHandler {
	return &WebsocketsHandler{
		service:                serv,
		upgrader:               upgrader,
		getDataFromRequestFunc: getDataFunc,
	}
}

func (wh *WebsocketsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data, err := wh.getDataFromRequestFunc(r)
	if err != nil {
		http.Error(
			w,
			"Could not get Data from websocket request",
			http.StatusBadRequest,
		)
		log.Errorf("Error while data from websocket hadnshake request: %+v\n error: %s", r, err.Error())
		return
	}

	wsConn, err := wh.upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(
			w,
			"Error while upgrading to websocket",
			http.StatusInternalServerError,
		)
		log.Errorf("Error while upgrading to websocket: %s", err.Error())
		return
	}

	// ! this should be blocking as long as the websocket connection is valid
	err = wh.service.HandleWSConn(wsConn, data)
	if err != nil {
		http.Error(
			w,
			"Error while handling websocket session",
			http.StatusInternalServerError,
		)
		log.Errorf("Error while handling websocket session with data: %s, error: %s", data, err.Error())
		return
	}
}
