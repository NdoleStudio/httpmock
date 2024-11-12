package requests

import (
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
	request.Subdomain = request.sanitizeString(request.Subdomain)
	return request
}

// ToProjectCreateParams creates services.ProjectCreateParams from ProjectCreateRequest
func (request *ProjectCreateRequest) ToProjectCreateParams(source string, userID entities.UserID) *services.ProjectCreateParams {
	return &services.ProjectCreateParams{
		Name:      request.Name,
		Subdomain: request.Subdomain,
		UserID:    userID,
		Source:    source,
	}
}
