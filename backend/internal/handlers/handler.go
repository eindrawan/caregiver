package handlers

import (
	"caregiver-shift-tracker/internal/middleware"
	"caregiver-shift-tracker/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// ScheduleServiceInterface defines the interface for schedule service
type ScheduleServiceInterface interface {
	GetAllSchedules(filter *models.ScheduleFilter) ([]models.Schedule, error)
	GetScheduleByID(id int) (*models.Schedule, error)
	GetTodaySchedules(caregiverID int) ([]models.Schedule, error)
	GetScheduleStats(caregiverID int) (*models.ScheduleStats, error)
	StartVisit(scheduleID int, req *models.VisitStartRequest) error
	EndVisit(scheduleID int, req *models.VisitEndRequest) error
	CancelVisit(scheduleID int) error
}

// VisitServiceInterface defines the interface for visit service
type VisitServiceInterface interface {
	GetVisitByScheduleID(scheduleID int) (*models.Visit, error)
}

// TaskServiceInterface defines the interface for task service
type TaskServiceInterface interface {
	GetTaskByID(id int) (*models.Task, error)
	UpdateTaskStatus(id int, req *models.TaskUpdateRequest) (*models.Task, error)
}

// ClientServiceInterface defines the interface for client service
type ClientServiceInterface interface {
	GetAllClients(filter *models.ClientFilter) ([]models.Client, error)
	GetClientByID(id int) (*models.Client, error)
	CreateClient(req *models.ClientCreateRequest) (*models.Client, error)
	UpdateClient(id int, req *models.ClientUpdateRequest) (*models.Client, error)
	DeleteClient(id int) error
	SearchClients(query string) ([]models.Client, error)
}

// Handler contains all HTTP handlers
type Handler struct {
	scheduleService ScheduleServiceInterface
	visitService    VisitServiceInterface
	taskService     TaskServiceInterface
	clientService   ClientServiceInterface
	logger          *logrus.Logger
}

// NewHandler creates a new handler
func NewHandler(
	scheduleService ScheduleServiceInterface,
	visitService VisitServiceInterface,
	taskService TaskServiceInterface,
	clientService ClientServiceInterface,
	logger *logrus.Logger,
) *Handler {
	return &Handler{
		scheduleService: scheduleService,
		visitService:    visitService,
		taskService:     taskService,
		clientService:   clientService,
		logger:          logger,
	}
}

// SetupRoutes sets up all routes
func (h *Handler) SetupRoutes() *gin.Engine {
	router := gin.New()

	// Middleware
	router.Use(gin.Logger())
	router.Use(middleware.ErrorHandlingMiddleware(h.logger))
	router.Use(middleware.RequestIDMiddleware())
	router.Use(middleware.SecurityHeadersMiddleware())
	router.Use(middleware.ValidationMiddleware())
	router.Use(middleware.RateLimitMiddleware(h.logger))
	router.Use(h.corsMiddleware())
	router.Use(h.loggingMiddleware())

	// Health check
	router.GET("/health", h.healthCheck)

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API routes
	api := router.Group("/api/v1")
	{
		// Schedule routes
		schedules := api.Group("/schedules")
		{
			schedules.GET("", h.getSchedules)
			schedules.GET("/today", h.getTodaySchedules)
			schedules.GET("/stats", h.getScheduleStats)
			schedules.GET("/:id", h.getScheduleByID)
			schedules.POST("/:id/start", h.startVisit)
			schedules.POST("/:id/end", h.endVisit)
			schedules.POST("/:id/cancel", h.cancelVisit)
		}

		// Task routes
		tasks := api.Group("/tasks")
		{
			tasks.GET("/:id", h.getTaskByID)
			tasks.PUT("/:id", h.updateTaskStatus)
		}

		// Visit routes (for additional visit operations if needed)
		visits := api.Group("/visits")
		{
			visits.GET("/schedule/:scheduleId", h.getVisitByScheduleID)
		}

		// Client routes
		clients := api.Group("/clients")
		{
			clients.GET("", h.getClients)
			clients.GET("/search", h.searchClients)
			clients.GET("/:id", h.getClientByID)
			clients.POST("", h.createClient)
			clients.PUT("/:id", h.updateClient)
			clients.DELETE("/:id", h.deleteClient)
		}
	}

	return router
}

// corsMiddleware handles CORS
func (h *Handler) corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set CORS headers for all requests
		c.Header("Access-Control-Allow-Origin", "http://localhost:8081")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Requested-With")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// loggingMiddleware logs requests
func (h *Handler) loggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		h.logger.WithFields(logrus.Fields{
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
			"ip":     c.ClientIP(),
		}).Info("Request received")

		c.Next()

		h.logger.WithFields(logrus.Fields{
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
			"status": c.Writer.Status(),
		}).Info("Request completed")
	}
}

// healthCheck returns the health status
func (h *Handler) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "caregiver-shift-tracker",
	})
}

// parseIntParam parses an integer parameter from the URL
func (h *Handler) parseIntParam(c *gin.Context, param string) (int, error) {
	str := c.Param(param)
	return strconv.Atoi(str)
}

// parseIntQuery parses an integer query parameter
func (h *Handler) parseIntQuery(c *gin.Context, param string) (*int, error) {
	str := c.Query(param)
	if str == "" {
		return nil, nil
	}
	val, err := strconv.Atoi(str)
	if err != nil {
		return nil, err
	}
	return &val, nil
}

// errorResponse sends an error response
func (h *Handler) errorResponse(c *gin.Context, statusCode int, message string, err error) {
	h.logger.WithError(err).WithFields(logrus.Fields{
		"status_code": statusCode,
		"message":     message,
		"path":        c.Request.URL.Path,
	}).Error("Request failed")

	c.JSON(statusCode, gin.H{
		"error":   message,
		"details": err.Error(),
	})
}

// successResponse sends a success response
func (h *Handler) successResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    data,
	})
}
