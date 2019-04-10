package cmd

import (
	"chat-v2/internal/server"
	"net/http"
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
