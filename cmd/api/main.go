package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/zuyatna/shop-retail-employee-service/internal/app"
	"github.com/zuyatna/shop-retail-employee-service/internal/config"
)

func main() {
	if err := godotenv.Load(); err != nil {
		// Try loading from root directory if running from cmd/api/
		if err := godotenv.Load("../../.env"); err != nil {
			log.Println("No .env file found, fallback to system env")
		}
	}

	cfg := config.Load()

	application, err := app.New(cfg)
	if err != nil {
		log.Fatalf("failed to start application: %v", err)
	}

	server := &http.Server{
		Addr:    cfg.HTTPAddr,
		Handler: application.Router,
	}

	// Start server
	go func() {
		log.Printf("HTTP server running on %s\n", cfg.HTTPAddr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen error: %v", err)
		}
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	log.Println("shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("http shutdown error: %v", err)
	}

	application.Close()

	log.Println("server gracefully stopped")
}
