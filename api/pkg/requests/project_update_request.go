package requests

import (
	"github.com/NdoleStudio/httpmock/pkg/entities"
	"github.com/NdoleStudio/httpmock/pkg/services"
)

// ProjectUpdateRequest is the payload for the /projects/create endpoint
type ProjectUpdateRequest struct {
	request
	ProjectID   string `json:"project_id" swaggerignore:"true"`
	Source      string `json:"source" swaggerignore:"true"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Sanitize the request by stripping whitespaces
func (request *ProjectUpdateRequest) Sanitize() *ProjectUpdateRequest {
	request.Name = request.sanitizeString(request.Name)
	request.Description = request.sanitizeString(request.Description)

	return request
}

// ToProjectUpdatePrams creates services.ProjectUpdateParams from ProjectUpdateRequest
func (request *ProjectUpdateRequest) ToProjectUpdatePrams(source string, userID entities.UserID) *services.ProjectUpdateParams {
	return &services.ProjectUpdateParams{
		Name:        request.Name,
		Description: request.Description,
		Source:      source,
		UserID:      userID,
	}
}
