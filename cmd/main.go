package main

import (
	"context"
	"github.com/airchains-network/da-client/routes"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const defaultPort = ":5050"

func main() {
	r := routes.SetupRouter()

	srv := &http.Server{
		Addr:    getServerPort(),
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}

// getServerPort gets the server port from the environment or returns the default
func getServerPort() string {
	if port, exists := os.LookupEnv("PORT"); exists {
		return ":" + port
	}
	return defaultPort
}
