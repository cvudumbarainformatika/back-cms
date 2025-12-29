package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/cvudumbarainformatika/backend/bootstrap"
)

func main() {
	// Initialize application
	app, err := bootstrap.NewApplication()
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	// Setup graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Run server in a goroutine
	go func() {
		if err := app.Run(); err != nil {
			log.Fatalf("Failed to run server: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-quit
	log.Println("Received shutdown signal")

	// Shutdown application
	if err := app.Shutdown(); err != nil {
		log.Fatalf("Failed to shutdown application: %v", err)
	}
}
