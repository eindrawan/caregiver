package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// getVisitByScheduleID retrieves a visit by schedule ID
// @Summary Get visit by schedule ID
// @Description Get visit details for a specific schedule
// @Tags visits
// @Accept json
// @Produce json
// @Param scheduleId path int true "Schedule ID"
// @Success 200 {object} map[string]interface{} "success response with visit details"
// @Failure 400 {object} map[string]interface{} "bad request"
// @Failure 404 {object} map[string]interface{} "visit not found"
// @Failure 500 {object} map[string]interface{} "internal server error"
// @Router /api/v1/visits/schedule/{scheduleId} [get]
func (h *Handler) getVisitByScheduleID(c *gin.Context) {
	scheduleID, err := h.parseIntParam(c, "scheduleId")
	if err != nil {
		h.errorResponse(c, http.StatusBadRequest, "Invalid schedule ID", err)
		return
	}

	visit, err := h.visitService.GetVisitByScheduleID(scheduleID)
	if err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "Failed to get visit", err)
		return
	}

	if visit == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Visit not found",
		})
		return
	}

	h.successResponse(c, visit)
}
