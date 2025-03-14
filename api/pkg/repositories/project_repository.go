package repositories

import (
	"context"

	"github.com/google/uuid"

	"github.com/NdoleStudio/httpmock/pkg/entities"
)

// ProjectRepository loads and persists an entities.Project
type ProjectRepository interface {
	// Store a new entities.Project
	Store(ctx context.Context, project *entities.Project) error

	// Update a new entities.Project
	Update(ctx context.Context, user *entities.Project) error

	// Fetch all entities.Project for a user
	Fetch(ctx context.Context, userID entities.UserID) ([]*entities.Project, error)

	// Load an entities.Project by entities.UserID
	Load(ctx context.Context, userID entities.UserID, projectID uuid.UUID) (*entities.Project, error)

	// Delete an entities.Project by entities.UserID and projectID
	Delete(ctx context.Context, userID entities.UserID, projectID uuid.UUID) error

	// LoadWithSubdomain load an entities.Project by a subdomain
	LoadWithSubdomain(ctx context.Context, subdomain string) (*entities.Project, error)
}
