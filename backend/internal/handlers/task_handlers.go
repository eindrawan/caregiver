package handlers

import (
	"caregiver-shift-tracker/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// getTaskByID retrieves a task by ID
// @Summary Get task by ID
// @Description Get a specific task by ID
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Success 200 {object} map[string]interface{} "success response with task details"
// @Failure 400 {object} map[string]interface{} "bad request"
// @Failure 404 {object} map[string]interface{} "task not found"
// @Failure 500 {object} map[string]interface{} "internal server error"
// @Router /api/v1/tasks/{id} [get]
func (h *Handler) getTaskByID(c *gin.Context) {
	id, err := h.parseIntParam(c, "id")
	if err != nil {
		h.errorResponse(c, http.StatusBadRequest, "Invalid task ID", err)
		return
	}

	task, err := h.taskService.GetTaskByID(id)
	if err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "Failed to get task", err)
		return
	}

	if task == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Task not found",
		})
		return
	}

	h.successResponse(c, task)
}

// updateTaskStatus updates the status of a task
// @Summary Update task status
// @Description Update the status of a task (completed or not_completed with reason)
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Param request body models.TaskUpdateRequest true "Task update request with status and optional reason"
// @Success 200 {object} map[string]interface{} "success response with updated task"
// @Failure 400 {object} map[string]interface{} "bad request"
// @Failure 404 {object} map[string]interface{} "task not found"
// @Failure 500 {object} map[string]interface{} "internal server error"
// @Router /api/v1/tasks/{id} [put]
func (h *Handler) updateTaskStatus(c *gin.Context) {
	id, err := h.parseIntParam(c, "id")
	if err != nil {
		h.errorResponse(c, http.StatusBadRequest, "Invalid task ID", err)
		return
	}

	var req models.TaskUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.errorResponse(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	updatedTask, err := h.taskService.UpdateTaskStatus(id, &req)
	if err != nil {
		if err.Error() == "task not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Task not found",
			})
			return
		}
		h.errorResponse(c, http.StatusInternalServerError, "Failed to update task status", err)
		return
	}

	h.successResponse(c, updatedTask)
}
