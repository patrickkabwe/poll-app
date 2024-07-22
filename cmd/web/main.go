package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"poll-app/cmd/api"
	"poll-app/config"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
)

func main() {
	// Create a new app
	app := &api.App{
		Router: chi.NewRouter(),
	}
	// The HTTP Server
	server := &http.Server{Addr: config.Env.PORT, Handler: app.Router}

	// Server run context
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, shutdownCtxRelease := context.WithTimeout(serverCtx, 30*time.Second)
		defer shutdownCtxRelease()

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		serverStopCtx()
	}()

	// Setup the app
	app.Setup()
	err := app.Start(server)
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
	log.Println("Server stopped")

	// Wait for server context to be stopped
	<-serverCtx.Done()
}
