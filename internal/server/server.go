package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	log "github.com/sirupsen/logrus"
)

func Run(port string, mux *http.ServeMux) {
	server := &http.Server{
		Addr:    port,
		Handler: mux,
	}

	startServer(server)

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)

	<-stop

	log.Infof("Server shutting down")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	if err := server.Shutdown(ctx); err != nil {
		log.Errorf("Error while shutting down server: %s", err.Error())
	}
}

func startServer(server *http.Server) {
	log.Infof("Starting server on %s", server.Addr)
	go func() {
		log.Fatal(server.ListenAndServe())
	}()
}
