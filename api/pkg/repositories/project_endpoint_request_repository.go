package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/oklog/ulid/v2"

	"github.com/NdoleStudio/httpmock/pkg/entities"
)

// ProjectEndpointRequestRepository loads and persists an entities.ProjectEndpointRequests
type ProjectEndpointRequestRepository interface {
	// Store a new entities.ProjectEndpointRequest
	Store(ctx context.Context, request *entities.ProjectEndpointRequest) error

	// Delete an entities.ProjectEndpointRequest
	Delete(ctx context.Context, request *entities.ProjectEndpointRequest) error

	// Load an entities.ProjectEndpointRequest by its ID
	Load(ctx context.Context, userID entities.UserID, requestID ulid.ULID) (*entities.ProjectEndpointRequest, error)

	// Index fetches the list of all project endpoint requests available to the currently authenticated user
	Index(ctx context.Context, userID entities.UserID, endpointID uuid.UUID, limit uint, previousID *ulid.ULID, nextID *ulid.ULID) ([]*entities.ProjectEndpointRequest, error)
}
