package server

import (
	"net/http"
)

type PathHandler struct {
	Path    string
	Handler http.Handler
}

func NewMux(pathHandlers ...PathHandler) (mux *http.ServeMux) {
	mux = http.NewServeMux()
	for _, ph := range pathHandlers {
		mux.Handle(ph.Path, ph.Handler)
	}

	return
}
