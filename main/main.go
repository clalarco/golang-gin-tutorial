package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"clalarco.io/rest/albums"
	"clalarco.io/rest/common"
	"github.com/gin-gonic/gin"
)

// Copied from https://dev.to/jacobsngoodwin/full-stack-memory-app-01-setup-go-server-with-reload-in-docker-62n
func main() {
	log.Println("Starting server...")

	router := gin.Default()

	AddRoutes(router)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Graceful server shutdown - https://github.com/gin-gonic/examples/blob/master/graceful-shutdown/graceful-shutdown/server.go
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to initialize server: %v\n", err)
		}
	}()

	log.Printf("Listening on port %v\n", srv.Addr)

	// Wait for kill signal of channel
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// This blocks until a signal is passed into the quit channel
	<-quit

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown server
	log.Println("Shutting down server...")
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v\n", err)
	}
}

func AddRoutes(engine *gin.Engine) {
	// Add routes
	common.AddRoutes(engine, "/ping")
	albums.AddRoutes(engine, "/albums")
}
