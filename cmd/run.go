package cmd

import (
	"net/http"

	"github.com/adrianbrad/chat-v2/internal/server"
)

func Run() {
	test := server.PathHandler{
		Path: "/",
		Handler: func() http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("hello"))
			})
		}(),
	}

	mux := server.NewMux(test)

	server.Run(":8080", mux)
}
