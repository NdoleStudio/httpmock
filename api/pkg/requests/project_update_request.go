package requests

import (
	"github.com/NdoleStudio/httpmock/pkg/entities"
	"github.com/NdoleStudio/httpmock/pkg/services"
	"github.com/google/uuid"
)

// ProjectUpdateRequest is the payload for the /projects/create endpoint
type ProjectUpdateRequest struct {
	request
	ProjectID   string `json:"projectId" swaggerignore:"true"`
	Name        string `json:"name"`
	Subdomain   string `json:"subdomain"`
	Description string `json:"description"`
}

// Sanitize the request by stripping whitespaces
func (request *ProjectUpdateRequest) Sanitize() *ProjectUpdateRequest {
	request.Name = request.sanitizeString(request.Name)
	request.Subdomain = request.sanitizeString(request.Subdomain)
	request.Description = request.sanitizeString(request.Description)

	return request
}

// ToProjectUpdatePrams creates services.ProjectUpdateParams from ProjectUpdateRequest
func (request *ProjectUpdateRequest) ToProjectUpdatePrams(source string, userID entities.UserID) *services.ProjectUpdateParams {
	return &services.ProjectUpdateParams{
		Name:        request.Name,
		Subdomain:   request.Subdomain,
		Description: request.Description,
		ProjectID:   uuid.MustParse(request.ProjectID),
		Source:      source,
		UserID:      userID,
	}
}
