package websocketshandler

import (
	"net/http"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type chatService interface {
	HandleWSConn(wsConn *websocket.Conn, data map[string]interface{}) (err error)
	ProcessData(data map[string]interface{}) (processedData map[string]interface{}, err error)
}

type upgrader interface {
	Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (wsConn *websocket.Conn, err error)
}

type WebsocketsHandler struct {
	chatService            chatService
	upgrader               upgrader
	getDataFromRequestFunc func(r *http.Request) (data map[string]interface{}, err error)
}

func NewWebsocketsHandler(
	serv chatService,
	upgrader upgrader,
	getDataFunc func(r *http.Request) (data map[string]interface{}, err error),
) *WebsocketsHandler {
	return &WebsocketsHandler{
		chatService:            serv,
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
		log.Errorf("Error while retrieving data from websocket hadnshake request: %+v\n error: %s", r, err.Error())
		return
	}

	processedData, err := wh.chatService.ProcessData(data)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		log.Errorf("Error while processing data from websocket handshake request: %+v\n error: %s", r, err.Error())
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
	//!FIXME in case of an error this will try to write to a hijacked WriteHeader as it was modified by the upgrader
	err = wh.chatService.HandleWSConn(wsConn, processedData)
	if err != nil {
		log.Errorf("Received following exit error from the wsConn with data: %s, %s", data, err.Error())
		return
	}
}
