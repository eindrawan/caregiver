package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"caregiver-shift-tracker/internal/config"
	"caregiver-shift-tracker/internal/database"
	"caregiver-shift-tracker/internal/handlers"
	"caregiver-shift-tracker/internal/repositories"
	"caregiver-shift-tracker/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	_ "caregiver-shift-tracker/docs" // Import generated docs
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Setup logger
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	if cfg.Environment == "development" {
		logger.SetLevel(logrus.DebugLevel)
		gin.SetMode(gin.DebugMode)
	} else {
		logger.SetLevel(logrus.InfoLevel)
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize database
	db, err := database.Initialize(cfg.DatabaseURL)
	if err != nil {
		logger.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := database.Migrate(db); err != nil {
		logger.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize repositories
	scheduleRepo := repositories.NewScheduleRepository(db)
	visitRepo := repositories.NewVisitRepository(db)
	taskRepo := repositories.NewTaskRepository(db)
	clientRepo := repositories.NewClientRepository(db)

	// Initialize services
	scheduleService := services.NewScheduleService(scheduleRepo, visitRepo, taskRepo, logger)
	visitService := services.NewVisitService(visitRepo, logger)
	taskService := services.NewTaskService(taskRepo, logger)
	clientService := services.NewClientService(clientRepo, logger)

	// Initialize handlers
	handler := handlers.NewHandler(scheduleService, visitService, taskService, clientService, logger)

	// Setup router
	router := handler.SetupRoutes()

	// Create server
	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		logger.Infof("Server starting on port %s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatalf("Server forced to shutdown: %v", err)
	}

	logger.Info("Server exited")
}
