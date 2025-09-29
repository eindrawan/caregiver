package services

import (
	"caregiver-shift-tracker/internal/models"
	"caregiver-shift-tracker/internal/repositories"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

// ClientService handles business logic for clients
type ClientService struct {
	clientRepo repositories.ClientRepository
	logger     *logrus.Logger
}

// NewClientService creates a new client service
func NewClientService(clientRepo repositories.ClientRepository, logger *logrus.Logger) *ClientService {
	return &ClientService{
		clientRepo: clientRepo,
		logger:     logger,
	}
}

// GetAllClients retrieves all clients with optional filtering
func (s *ClientService) GetAllClients(filter *models.ClientFilter) ([]models.Client, error) {
	s.logger.Debug("Getting all clients")

	clients, err := s.clientRepo.GetAll(filter)
	if err != nil {
		s.logger.WithError(err).Error("Failed to get clients")
		return nil, fmt.Errorf("failed to get clients: %w", err)
	}

	s.logger.WithField("count", len(clients)).Debug("Successfully retrieved clients")
	return clients, nil
}

// GetClientByID retrieves a client by ID
func (s *ClientService) GetClientByID(id int) (*models.Client, error) {
	s.logger.WithField("client_id", id).Debug("Getting client by ID")

	if id <= 0 {
		return nil, fmt.Errorf("invalid client ID: %d", id)
	}

	client, err := s.clientRepo.GetByID(id)
	if err != nil {
		s.logger.WithError(err).WithField("client_id", id).Error("Failed to get client")
		return nil, fmt.Errorf("failed to get client: %w", err)
	}

	if client == nil {
		s.logger.WithField("client_id", id).Debug("Client not found")
		return nil, nil
	}

	s.logger.WithField("client_id", id).Debug("Successfully retrieved client")
	return client, nil
}

// CreateClient creates a new client
func (s *ClientService) CreateClient(req *models.ClientCreateRequest) (*models.Client, error) {
	s.logger.WithField("client_name", req.Name).Debug("Creating new client")

	if err := s.validateClientCreateRequest(req); err != nil {
		s.logger.WithError(err).WithField("client_name", req.Name).Error("Client validation failed")
		return nil, fmt.Errorf("client validation failed: %w", err)
	}

	client := &models.Client{
		Name:      req.Name,
		Email:     req.Email,
		Phone:     req.Phone,
		Address:   req.Address,
		City:      req.City,
		State:     req.State,
		ZipCode:   req.ZipCode,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		Notes:     req.Notes,
		IsActive:  true, // New clients are active by default
	}

	if err := s.clientRepo.Create(client); err != nil {
		s.logger.WithError(err).WithField("client_name", req.Name).Error("Failed to create client")
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	s.logger.WithField("client_id", client.ID).Info("Successfully created client")
	return client, nil
}

// UpdateClient updates an existing client
func (s *ClientService) UpdateClient(id int, req *models.ClientUpdateRequest) (*models.Client, error) {
	s.logger.WithField("client_id", id).Debug("Updating client")

	if id <= 0 {
		return nil, fmt.Errorf("invalid client ID: %d", id)
	}

	// Get existing client
	client, err := s.clientRepo.GetByID(id)
	if err != nil {
		s.logger.WithError(err).WithField("client_id", id).Error("Failed to get client for update")
		return nil, fmt.Errorf("failed to get client: %w", err)
	}

	if client == nil {
		s.logger.WithField("client_id", id).Debug("Client not found for update")
		return nil, nil
	}

	// Update fields if provided
	if req.Name != nil {
		client.Name = *req.Name
	}
	if req.Email != nil {
		client.Email = *req.Email
	}
	if req.Phone != nil {
		client.Phone = *req.Phone
	}
	if req.Address != nil {
		client.Address = *req.Address
	}
	if req.City != nil {
		client.City = *req.City
	}
	if req.State != nil {
		client.State = *req.State
	}
	if req.ZipCode != nil {
		client.ZipCode = *req.ZipCode
	}
	if req.Latitude != nil {
		client.Latitude = *req.Latitude
	}
	if req.Longitude != nil {
		client.Longitude = *req.Longitude
	}
	if req.Notes != nil {
		client.Notes = *req.Notes
	}
	if req.IsActive != nil {
		client.IsActive = *req.IsActive
	}

	if err := s.validateClient(client); err != nil {
		s.logger.WithError(err).WithField("client_id", id).Error("Client validation failed")
		return nil, fmt.Errorf("client validation failed: %w", err)
	}

	if err := s.clientRepo.Update(client); err != nil {
		s.logger.WithError(err).WithField("client_id", id).Error("Failed to update client")
		return nil, fmt.Errorf("failed to update client: %w", err)
	}

	s.logger.WithField("client_id", id).Info("Successfully updated client")
	return client, nil
}

// DeleteClient deletes a client
func (s *ClientService) DeleteClient(id int) error {
	s.logger.WithField("client_id", id).Debug("Deleting client")

	if id <= 0 {
		return fmt.Errorf("invalid client ID: %d", id)
	}

	// Check if client exists
	client, err := s.clientRepo.GetByID(id)
	if err != nil {
		s.logger.WithError(err).WithField("client_id", id).Error("Failed to get client for deletion")
		return fmt.Errorf("failed to get client: %w", err)
	}

	if client == nil {
		s.logger.WithField("client_id", id).Debug("Client not found for deletion")
		return fmt.Errorf("client not found")
	}

	if err := s.clientRepo.Delete(id); err != nil {
		s.logger.WithError(err).WithField("client_id", id).Error("Failed to delete client")
		return fmt.Errorf("failed to delete client: %w", err)
	}

	s.logger.WithField("client_id", id).Info("Successfully deleted client")
	return nil
}

// SearchClients searches for clients by name, email, or phone
func (s *ClientService) SearchClients(query string) ([]models.Client, error) {
	s.logger.WithField("query", query).Debug("Searching clients")

	if strings.TrimSpace(query) == "" {
		return []models.Client{}, nil
	}

	clients, err := s.clientRepo.Search(query)
	if err != nil {
		s.logger.WithError(err).WithField("query", query).Error("Failed to search clients")
		return nil, fmt.Errorf("failed to search clients: %w", err)
	}

	s.logger.WithField("query", query).WithField("count", len(clients)).Debug("Successfully searched clients")
	return clients, nil
}

// validateClientCreateRequest validates a client create request
func (s *ClientService) validateClientCreateRequest(req *models.ClientCreateRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return fmt.Errorf("client name is required")
	}

	if strings.TrimSpace(req.Address) == "" {
		return fmt.Errorf("address is required")
	}

	if strings.TrimSpace(req.City) == "" {
		return fmt.Errorf("city is required")
	}

	if strings.TrimSpace(req.State) == "" {
		return fmt.Errorf("state is required")
	}

	if strings.TrimSpace(req.ZipCode) == "" {
		return fmt.Errorf("zip code is required")
	}

	return nil
}

// validateClient validates a client
func (s *ClientService) validateClient(client *models.Client) error {
	if strings.TrimSpace(client.Name) == "" {
		return fmt.Errorf("client name is required")
	}

	if strings.TrimSpace(client.Address) == "" {
		return fmt.Errorf("address is required")
	}

	if strings.TrimSpace(client.City) == "" {
		return fmt.Errorf("city is required")
	}

	if strings.TrimSpace(client.State) == "" {
		return fmt.Errorf("state is required")
	}

	if strings.TrimSpace(client.ZipCode) == "" {
		return fmt.Errorf("zip code is required")
	}

	return nil
}
