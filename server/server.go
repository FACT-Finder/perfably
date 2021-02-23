package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

var notifySignal = signal.Notify
var serverShutdown = func(server *http.Server, ctx context.Context) error {
	return server.Shutdown(ctx)
}

// Start starts the http server
func Start(mux *mux.Router, addr string) error {
	server, shutdown := startServer(mux, addr)
	shutdownOnInterruptSignal(server, 2*time.Second, shutdown)
	return waitForServerToClose(shutdown)
}

func startServer(mux *mux.Router, addr string) (*http.Server, chan error) {
	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	shutdown := make(chan error)
	go func() {
		err := srv.ListenAndServe()
		shutdown <- err
	}()
	return srv, shutdown
}

func shutdownOnInterruptSignal(server *http.Server, timeout time.Duration, shutdown chan<- error) {
	interrupt := make(chan os.Signal, 1)
	notifySignal(interrupt, os.Interrupt)

	go func() {
		<-interrupt
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		if err := serverShutdown(server, ctx); err != nil {
			shutdown <- err
		}
	}()
}

func waitForServerToClose(shutdown <-chan error) error {
	err := <-shutdown
	if err == http.ErrServerClosed {
		return nil
	}
	return err
}
