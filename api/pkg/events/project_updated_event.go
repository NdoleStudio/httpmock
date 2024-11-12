package events

import (
	"time"

	"github.com/google/uuid"

	"github.com/NdoleStudio/httpmock/pkg/entities"
)

// ProjectUpdated is raised when a user is created
const ProjectUpdated = "project.updated"

// ProjectUpdatedPayload stores the data for the ProjectUpdated event
type ProjectUpdatedPayload struct {
	UserID             entities.UserID `json:"user_id"`
	ProjectID          uuid.UUID       `json:"project_id"`
	ProjectSubdomain   string          `json:"project_subdomain"`
	ProjectName        string          `json:"project_name"`
	ProjectDescription string          `json:"project_description"`
	ProjectUpdatedAt   time.Time       `json:"project_updated_at"`
}
