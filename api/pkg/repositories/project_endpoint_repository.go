package repositories

import (
	"context"

	"github.com/google/uuid"

	"github.com/NdoleStudio/httpmock/pkg/entities"
)

// ProjectEndpointRepository loads and persists an entities.ProjectEndpoint
type ProjectEndpointRepository interface {
	// Store a new entities.ProjectEndpoint
	Store(ctx context.Context, project *entities.ProjectEndpoint) error

	// Update a new entities.ProjectEndpoint
	Update(ctx context.Context, user *entities.ProjectEndpoint) error

	// RegisterRequest register a request for an entities.ProjectEndpoint
	RegisterRequest(ctx context.Context, projectEndpointID uuid.UUID) error

	// UpdateSubdomain for an entities.ProjectEndpoint after the project was updated
	UpdateSubdomain(ctx context.Context, subdomain string, projectEndpointID uuid.UUID) error

	// Fetch all entities.ProjectEndpoint for a user
	Fetch(ctx context.Context, userID entities.UserID, projectID uuid.UUID) ([]*entities.ProjectEndpoint, error)

	// Load an entities.ProjectEndpoint by entities.UserID
	Load(ctx context.Context, userID entities.UserID, projectID uuid.UUID, projectEndpointID uuid.UUID) (*entities.ProjectEndpoint, error)

	// Delete an entities.ProjectEndpoint
	Delete(ctx context.Context, endpoint *entities.ProjectEndpoint) error

	// LoadByRequestForUser load an entities.ProjectEndpoint by a request path and method.
	LoadByRequestForUser(ctx context.Context, userID entities.UserID, projectID uuid.UUID, requestMethod, requestPath string) (*entities.ProjectEndpoint, error)

	// LoadByRequest load an entities.ProjectEndpoint by a http path and method.
	LoadByRequest(ctx context.Context, subdomain, requestMethod, requestPath string) (*entities.ProjectEndpoint, error)
}
