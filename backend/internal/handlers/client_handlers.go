package handlers

import (
	"caregiver-shift-tracker/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// getClients retrieves all clients with optional filtering
// @Summary Get all clients
// @Description Get all clients with optional filtering
// @Tags clients
// @Accept json
// @Produce json
// @Param is_active query boolean false "Filter by active status"
// @Param city query string false "Filter by city"
// @Param state query string false "Filter by state"
// @Param search query string false "Search by name, email, or phone"
// @Param limit query int false "Limit number of results"
// @Param offset query int false "Offset for pagination"
// @Success 200 {object} map[string]interface{} "success response with clients list"
// @Failure 400 {object} map[string]interface{} "bad request"
// @Failure 500 {object} map[string]interface{} "internal server error"
// @Router /api/v1/clients [get]
func (h *Handler) getClients(c *gin.Context) {
	filter := &models.ClientFilter{}

	// Parse query parameters
	if isActiveStr := c.Query("is_active"); isActiveStr != "" {
		if isActive, err := strconv.ParseBool(isActiveStr); err != nil {
			h.errorResponse(c, http.StatusBadRequest, "Invalid is_active parameter", err)
			return
		} else {
			filter.IsActive = &isActive
		}
	}

	if city := c.Query("city"); city != "" {
		filter.City = &city
	}

	if state := c.Query("state"); state != "" {
		filter.State = &state
	}

	if search := c.Query("search"); search != "" {
		filter.Search = &search
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err != nil {
			h.errorResponse(c, http.StatusBadRequest, "Invalid limit parameter", err)
			return
		} else {
			filter.Limit = &limit
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err != nil {
			h.errorResponse(c, http.StatusBadRequest, "Invalid offset parameter", err)
			return
		} else {
			filter.Offset = &offset
		}
	}

	clients, err := h.clientService.GetAllClients(filter)
	if err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "Failed to get clients", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Clients retrieved successfully",
		"data": gin.H{
			"clients": clients,
			"count":   len(clients),
		},
	})
}

// getClientByID retrieves a client by ID
// @Summary Get client by ID
// @Description Get a specific client by ID
// @Tags clients
// @Accept json
// @Produce json
// @Param id path int true "Client ID"
// @Success 200 {object} map[string]interface{} "success response with client details"
// @Failure 400 {object} map[string]interface{} "bad request"
// @Failure 404 {object} map[string]interface{} "client not found"
// @Failure 500 {object} map[string]interface{} "internal server error"
// @Router /api/v1/clients/{id} [get]
func (h *Handler) getClientByID(c *gin.Context) {
	id, err := h.parseIntParam(c, "id")
	if err != nil {
		h.errorResponse(c, http.StatusBadRequest, "Invalid client ID", err)
		return
	}

	client, err := h.clientService.GetClientByID(id)
	if err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "Failed to get client", err)
		return
	}

	if client == nil {
		h.errorResponse(c, http.StatusNotFound, "Client not found", nil)
		return
	}

	h.successResponse(c, gin.H{
		"client": client,
	})
}

// createClient creates a new client
// @Summary Create a new client
// @Description Create a new client with the provided information
// @Tags clients
// @Accept json
// @Produce json
// @Param client body models.ClientCreateRequest true "Client creation data"
// @Success 201 {object} map[string]interface{} "success response with created client"
// @Failure 400 {object} map[string]interface{} "bad request"
// @Failure 500 {object} map[string]interface{} "internal server error"
// @Router /api/v1/clients [post]
func (h *Handler) createClient(c *gin.Context) {
	var req models.ClientCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.errorResponse(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	client, err := h.clientService.CreateClient(&req)
	if err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "Failed to create client", err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Client created successfully",
		"data": gin.H{
			"client": client,
		},
	})
}

// updateClient updates an existing client
// @Summary Update a client
// @Description Update an existing client with the provided information
// @Tags clients
// @Accept json
// @Produce json
// @Param id path int true "Client ID"
// @Param client body models.ClientUpdateRequest true "Client update data"
// @Success 200 {object} map[string]interface{} "success response with updated client"
// @Failure 400 {object} map[string]interface{} "bad request"
// @Failure 404 {object} map[string]interface{} "client not found"
// @Failure 500 {object} map[string]interface{} "internal server error"
// @Router /api/v1/clients/{id} [put]
func (h *Handler) updateClient(c *gin.Context) {
	id, err := h.parseIntParam(c, "id")
	if err != nil {
		h.errorResponse(c, http.StatusBadRequest, "Invalid client ID", err)
		return
	}

	var req models.ClientUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.errorResponse(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	client, err := h.clientService.UpdateClient(id, &req)
	if err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "Failed to update client", err)
		return
	}

	if client == nil {
		h.errorResponse(c, http.StatusNotFound, "Client not found", nil)
		return
	}

	h.successResponse(c, gin.H{
		"client": client,
	})
}

// deleteClient deletes a client
// @Summary Delete a client
// @Description Delete a client by ID
// @Tags clients
// @Accept json
// @Produce json
// @Param id path int true "Client ID"
// @Success 200 {object} map[string]interface{} "success response"
// @Failure 400 {object} map[string]interface{} "bad request"
// @Failure 404 {object} map[string]interface{} "client not found"
// @Failure 409 {object} map[string]interface{} "conflict - client has associated schedules"
// @Failure 500 {object} map[string]interface{} "internal server error"
// @Router /api/v1/clients/{id} [delete]
func (h *Handler) deleteClient(c *gin.Context) {
	id, err := h.parseIntParam(c, "id")
	if err != nil {
		h.errorResponse(c, http.StatusBadRequest, "Invalid client ID", err)
		return
	}

	err = h.clientService.DeleteClient(id)
	if err != nil {
		if err.Error() == "client not found" {
			h.errorResponse(c, http.StatusNotFound, "Client not found", err)
			return
		}
		if err.Error() == "cannot delete client: client has associated schedules" {
			h.errorResponse(c, http.StatusConflict, "Cannot delete client with associated schedules", err)
			return
		}
		h.errorResponse(c, http.StatusInternalServerError, "Failed to delete client", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Client deleted successfully",
	})
}

// searchClients searches for clients
// @Summary Search clients
// @Description Search for clients by name, email, or phone
// @Tags clients
// @Accept json
// @Produce json
// @Param q query string true "Search query"
// @Success 200 {object} map[string]interface{} "success response with search results"
// @Failure 400 {object} map[string]interface{} "bad request"
// @Failure 500 {object} map[string]interface{} "internal server error"
// @Router /api/v1/clients/search [get]
func (h *Handler) searchClients(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		h.errorResponse(c, http.StatusBadRequest, "Search query is required", nil)
		return
	}

	clients, err := h.clientService.SearchClients(query)
	if err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "Failed to search clients", err)
		return
	}

	h.successResponse(c, gin.H{
		"clients": clients,
		"count":   len(clients),
		"query":   query,
	})
}
