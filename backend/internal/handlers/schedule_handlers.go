package handlers

import (
	"caregiver-shift-tracker/internal/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// getSchedules retrieves all schedules with optional filtering
// @Summary Get all schedules
// @Description Get all schedules with optional filtering by caregiver_id, date, status
// @Tags schedules
// @Accept json
// @Produce json
// @Param caregiver_id query int false "Filter by caregiver ID"
// @Param date query string false "Filter by date (YYYY-MM-DD format)"
// @Param status query string false "Filter by status (scheduled, in_progress, completed, missed)"
// @Param limit query int false "Limit number of results"
// @Param offset query int false "Offset for pagination"
// @Success 200 {object} map[string]interface{} "success response with schedules data"
// @Failure 400 {object} map[string]interface{} "bad request"
// @Failure 500 {object} map[string]interface{} "internal server error"
// @Router /api/v1/schedules [get]
func (h *Handler) getSchedules(c *gin.Context) {
	// Parse query parameters
	filter := &models.ScheduleFilter{}

	if caregiverID, err := h.parseIntQuery(c, "caregiver_id"); err != nil {
		h.errorResponse(c, http.StatusBadRequest, "Invalid caregiver_id", err)
		return
	} else {
		filter.CaregiverID = caregiverID
	}

	if dateStr := c.Query("date"); dateStr != "" {
		if date, err := time.Parse("2006-01-02", dateStr); err != nil {
			h.errorResponse(c, http.StatusBadRequest, "Invalid date format, use YYYY-MM-DD", err)
			return
		} else {
			filter.Date = &date
		}
	}

	if status := c.Query("status"); status != "" {
		filter.Status = &status
	}

	if limit, err := h.parseIntQuery(c, "limit"); err != nil {
		h.errorResponse(c, http.StatusBadRequest, "Invalid limit", err)
		return
	} else {
		filter.Limit = limit
	}

	if offset, err := h.parseIntQuery(c, "offset"); err != nil {
		h.errorResponse(c, http.StatusBadRequest, "Invalid offset", err)
		return
	} else {
		filter.Offset = offset
	}

	// Get schedules
	schedules, err := h.scheduleService.GetAllSchedules(filter)
	if err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "Failed to get schedules", err)
		return
	}

	h.successResponse(c, schedules)
}

// getTodaySchedules retrieves today's schedules for a caregiver
// @Summary Get today's schedules
// @Description Get today's schedules for a specific caregiver
// @Tags schedules
// @Accept json
// @Produce json
// @Param caregiver_id query int true "Caregiver ID"
// @Success 200 {object} map[string]interface{} "success response with today's schedules"
// @Failure 400 {object} map[string]interface{} "bad request"
// @Failure 500 {object} map[string]interface{} "internal server error"
// @Router /api/v1/schedules/today [get]
func (h *Handler) getTodaySchedules(c *gin.Context) {
	caregiverID, err := h.parseIntQuery(c, "caregiver_id")
	if err != nil || caregiverID == nil {
		h.errorResponse(c, http.StatusBadRequest, "caregiver_id is required", err)
		return
	}

	schedules, err := h.scheduleService.GetTodaySchedules(*caregiverID)
	if err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "Failed to get today's schedules", err)
		return
	}

	h.successResponse(c, schedules)
}

// getScheduleStats retrieves schedule statistics for a caregiver
// @Summary Get schedule statistics
// @Description Get schedule statistics for a specific caregiver
// @Tags schedules
// @Accept json
// @Produce json
// @Param caregiver_id query int true "Caregiver ID"
// @Success 200 {object} map[string]interface{} "success response with schedule statistics"
// @Failure 400 {object} map[string]interface{} "bad request"
// @Failure 500 {object} map[string]interface{} "internal server error"
// @Router /api/v1/schedules/stats [get]
func (h *Handler) getScheduleStats(c *gin.Context) {
	caregiverID, err := h.parseIntQuery(c, "caregiver_id")
	if err != nil || caregiverID == nil {
		h.errorResponse(c, http.StatusBadRequest, "caregiver_id is required", err)
		return
	}

	stats, err := h.scheduleService.GetScheduleStats(*caregiverID)
	if err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "Failed to get schedule stats", err)
		return
	}

	h.successResponse(c, stats)
}

// getScheduleByID retrieves a schedule by ID with full details
// @Summary Get schedule by ID
// @Description Get a specific schedule by ID with full details including visit and tasks
// @Tags schedules
// @Accept json
// @Produce json
// @Param id path int true "Schedule ID"
// @Success 200 {object} map[string]interface{} "success response with schedule details"
// @Failure 400 {object} map[string]interface{} "bad request"
// @Failure 404 {object} map[string]interface{} "schedule not found"
// @Failure 500 {object} map[string]interface{} "internal server error"
// @Router /api/v1/schedules/{id} [get]
func (h *Handler) getScheduleByID(c *gin.Context) {
	id, err := h.parseIntParam(c, "id")
	if err != nil {
		h.errorResponse(c, http.StatusBadRequest, "Invalid schedule ID", err)
		return
	}

	schedule, err := h.scheduleService.GetScheduleByID(id)
	if err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "Failed to get schedule", err)
		return
	}

	if schedule == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Schedule not found",
		})
		return
	}

	h.successResponse(c, schedule)
}

// startVisit starts a visit for a schedule
// @Summary Start a visit
// @Description Start a visit for a specific schedule with geolocation
// @Tags schedules
// @Accept json
// @Produce json
// @Param id path int true "Schedule ID"
// @Param request body models.VisitStartRequest true "Visit start request with geolocation"
// @Success 200 {object} map[string]interface{} "success response"
// @Failure 400 {object} map[string]interface{} "bad request"
// @Failure 404 {object} map[string]interface{} "schedule not found"
// @Failure 500 {object} map[string]interface{} "internal server error"
// @Router /api/v1/schedules/{id}/start [post]
func (h *Handler) startVisit(c *gin.Context) {
	id, err := h.parseIntParam(c, "id")
	if err != nil {
		h.errorResponse(c, http.StatusBadRequest, "Invalid schedule ID", err)
		return
	}

	var req models.VisitStartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.errorResponse(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if err := h.scheduleService.StartVisit(id, &req); err != nil {
		if err.Error() == "schedule not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Schedule not found",
			})
			return
		}
		h.errorResponse(c, http.StatusInternalServerError, "Failed to start visit", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Visit started successfully",
	})
}

// endVisit ends a visit for a schedule
// @Summary End a visit
// @Description End a visit for a specific schedule with geolocation
// @Tags schedules
// @Accept json
// @Produce json
// @Param id path int true "Schedule ID"
// @Param request body models.VisitEndRequest true "Visit end request with geolocation and optional notes"
// @Success 200 {object} map[string]interface{} "success response"
// @Failure 400 {object} map[string]interface{} "bad request"
// @Failure 404 {object} map[string]interface{} "schedule not found"
// @Failure 500 {object} map[string]interface{} "internal server error"
// @Router /api/v1/schedules/{id}/end [post]
func (h *Handler) endVisit(c *gin.Context) {
	id, err := h.parseIntParam(c, "id")
	if err != nil {
		h.errorResponse(c, http.StatusBadRequest, "Invalid schedule ID", err)
		return
	}

	var req models.VisitEndRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.errorResponse(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if err := h.scheduleService.EndVisit(id, &req); err != nil {
		if err.Error() == "schedule not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Schedule not found",
			})
			return
		}
		h.errorResponse(c, http.StatusInternalServerError, "Failed to end visit", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Visit ended successfully",
	})
}

// cancelVisit cancels an in-progress visit for a schedule
// @Summary Cancel a visit
// @Description Cancel an in-progress visit for a specific schedule
// @Tags schedules
// @Accept json
// @Produce json
// @Param id path int true "Schedule ID"
// @Success 200 {object} map[string]interface{} "success response"
// @Failure 400 {object} map[string]interface{} "bad request"
// @Failure 404 {object} map[string]interface{} "schedule not found"
// @Failure 500 {object} map[string]interface{} "internal server error"
// @Router /api/v1/schedules/{id}/cancel [post]
func (h *Handler) cancelVisit(c *gin.Context) {
	h.logger.Info("Cancelling visit")
	id, err := h.parseIntParam(c, "id")
	if err != nil {
		h.errorResponse(c, http.StatusBadRequest, "Invalid schedule ID", err)
		return
	}

	if err := h.scheduleService.CancelVisit(id); err != nil {
		if err.Error() == "schedule not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Schedule not found",
			})
			return
		}
		if err.Error() == "visit not in progress" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Visit is not in progress",
			})
			return
		}
		h.errorResponse(c, http.StatusInternalServerError, "Failed to cancel visit", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Visit cancelled successfully",
	})
}
