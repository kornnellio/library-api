package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"library-api/internal/config"
	"library-api/internal/database"
	"library-api/internal/logger"
	"library-api/internal/middleware"
	"library-api/internal/repository"
	"library-api/internal/service"
	"library-api/internal/transport"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize logger
	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}
	if err := logger.Init(env); err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}
	defer logger.Sync()

	// Set Gin mode based on environment
	if env == "production" || env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize database
	db, err := database.Connect(&cfg.Database, logger.Logger)
	if err != nil {
		logger.Logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	// Create database tables
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := db.CreateTables(ctx); err != nil {
		logger.Logger.Fatal("Failed to create tables", zap.Error(err))
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db.DB)
	bookRepo := repository.NewBookRepository(db.DB)

	// Initialize services
	userService := service.NewUserService(userRepo)
	bookService := service.NewBookService(bookRepo)

	// Initialize handlers
	handler := transport.NewHandler(userService, bookService)

	// Setup Gin router
	r := gin.New()

	// Middleware
	r.Use(gin.Recovery())
	r.Use(logger.RequestIDMiddleware())
	r.Use(logger.LoggerMiddleware())
	r.Use(middleware.CORS())
	r.Use(middleware.SecurityHeaders())
	r.Use(middleware.RequestSizeLimit(1 << 20)) // 1MB limit

	// Setup routes
	handler.SetupRoutes(r)

	// Setup HTTP server
	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Start server in a goroutine
	go func() {
		logger.Logger.Info("Server starting", zap.String("port", cfg.Server.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Logger.Fatal("Server failed", zap.Error(err))
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	logger.Logger.Info("Shutting down server...")

	// Graceful shutdown with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Logger.Info("Server stopped")
}
