package requests

import (
	"strings"

	"github.com/NdoleStudio/httpmock/pkg/entities"
	"github.com/NdoleStudio/httpmock/pkg/services"
)

// ProjectCreateRequest is the payload for the /projects/create endpoint
type ProjectCreateRequest struct {
	request
	Name        string `json:"name"`
	Description string `json:"description"`
	Subdomain   string `json:"subdomain"`
}

// Sanitize the request by stripping whitespaces
func (request *ProjectCreateRequest) Sanitize() *ProjectCreateRequest {
	request.Name = request.sanitizeString(request.Name)
	request.Description = request.sanitizeString(request.Description)
	request.Subdomain = strings.TrimSuffix(request.sanitizeString(request.Subdomain), ".httpmock.dev")
	return request
}

// ToProjectCreateParams creates services.ProjectCreateParams from ProjectCreateRequest
func (request *ProjectCreateRequest) ToProjectCreateParams(source string, userID entities.UserID) *services.ProjectCreateParams {
	return &services.ProjectCreateParams{
		Name:        request.Name,
		Description: request.Description,
		Subdomain:   request.Subdomain,
		UserID:      userID,
		Source:      source,
	}
}
