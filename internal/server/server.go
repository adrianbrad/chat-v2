package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	log "github.com/sirupsen/logrus"
)

func Run(port string, mux *http.ServeMux, stop chan os.Signal, timeout time.Duration) {
	if stop == nil {
		log.Fatal("Nil stop channel passed to server.Run")
	}
	server := &http.Server{
		Addr:    port,
		Handler: mux,
	}

	startServer(server)

	signal.Notify(stop, os.Interrupt)

	<-stop

	log.Infof("Server shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Errorf("Error while shutting down server: %s", err.Error())
	}
}

func startServer(server *http.Server) {
	log.Infof("Starting server on %s", server.Addr)
	go func() {
		e := server.ListenAndServe()
		if e.Error() != "http: Server closed" {
			log.Fatal(e)
		} else {
			log.Info(e)
		}
	}()
}
