package server

import "net/http"

type service interface {
	HandleJoin()
	HandleMessage()
	HandleLeave()
}
type websocketsHandler struct {
	service
}

func NewWebsocketsHandler(serv service) *websocketsHandler {
	return &websocketsHandler{
		service: serv,
	}
}

func (wh *websocketsHandler) ServeHTTP(w http.ResponseWriter, r *http.Response) {
	//TODO upgrade to websockets
}
