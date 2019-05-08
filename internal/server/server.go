package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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
	serverShutdown := make(chan struct{}, 1)
	startServer(server, serverShutdown)

	signal.Notify(stop, syscall.SIGTERM)
	signal.Notify(stop, syscall.SIGINT)

	a := <-stop

	fmt.Println(a.String())

	log.Infof("Server shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Errorf("Error while shutting down server: %s", err.Error())
	}
	<-serverShutdown
}

func startServer(server *http.Server, serverShutdown chan struct{}) {
	log.Infof("Starting server on %s", server.Addr)
	go func() {
		e := server.ListenAndServe()

		if e.Error() != "http: Server closed" {
			log.Fatal(e)
		} else {
			log.Info(e)
		}
		serverShutdown <- struct{}{}
	}()
}
