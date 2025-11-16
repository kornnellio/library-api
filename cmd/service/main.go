package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"library-api/internal/database"
	"library-api/internal/repository"
	"library-api/internal/transport"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	if err := database.Connect(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.Close()

	// Create database tables
	if err := database.CreateTables(); err != nil {
		log.Fatal("Failed to create tables:", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(database.DB)
	bookRepo := repository.NewBookRepository(database.DB)

	// Initialize handlers
	handler := transport.NewHandler(userRepo, bookRepo)

	// Setup Gin router
	r := gin.Default()
	handler.SetupRoutes(r)

	// Setup HTTP server
	srv := &http.Server{
		Addr:    ":" + getPort(),
		Handler: r,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on :%s", getPort())
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server failed:", err)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server stopped")
}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return port
}
