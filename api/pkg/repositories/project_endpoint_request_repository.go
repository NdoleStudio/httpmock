package repositories

import (
	"context"

	"github.com/NdoleStudio/httpmock/pkg/entities"
)

// ProjectEndpointRequestRepository loads and persists an entities.ProjectEndpointRequests
type ProjectEndpointRequestRepository interface {
	// Store a new entities.ProjectEndpointRequest
	Store(ctx context.Context, request *entities.ProjectEndpointRequest) error
}
